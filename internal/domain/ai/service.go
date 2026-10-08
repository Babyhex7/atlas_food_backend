package ai

import (
	"atlas_food/internal/domain/submission"
	"atlas_food/internal/pkg/groq"
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

// refreshCooldown - jeda minimal antara dua analisis baru untuk submission yang sama.
// Permintaan force_refresh di dalam jeda ini dilayani dari hasil tersimpan,
// supaya tombol "analisis ulang" tidak bisa dipakai menguras kuota LLM.
const refreshCooldown = 30 * time.Second

// maxGenerateAttempts - berapa kali LLM boleh diminta ulang bila jawabannya tidak dapat diproses
const maxGenerateAttempts = 2

// generateBudget - batas waktu total satu analisis, mencakup seluruh percobaan.
// Harus di bawah batas waktu klien (60 detik): tanpa ini dua percobaan yang
// masing-masing nyaris habis waktu membuat klien menyerah lebih dulu dan
// menampilkan galat jaringan, padahal server masih bekerja.
const generateBudget = 50 * time.Second

// Completer - penyedia LLM yang dibutuhkan service. Dipisah sebagai interface
// supaya service dapat diuji tanpa jaringan dan penyedia dapat diganti.
type Completer interface {
	Configured() bool
	CompleteJSON(ctx context.Context, systemPrompt, userPrompt string) (*groq.Result, error)
}

// Service - business logic AI analysis
type Service interface {
	// AnalyzeNutrition - hasilkan (atau ambil dari simpanan) analisis untuk satu submission
	AnalyzeNutrition(ctx context.Context, userID string, req NutritionAnalysisRequest) (*NutritionAnalysisResult, error)
	// GetNutritionAnalysis - ambil hasil tersimpan tanpa pernah memanggil LLM
	GetNutritionAnalysis(userID, submissionRef string) (*NutritionAnalysisResult, error)
}

type service struct {
	repo  Repository
	llm   Completer
	locks keyedMutex
	now   func() time.Time
}

// NewService - buat service AI
func NewService(repo Repository, llm Completer) Service {
	return &service{repo: repo, llm: llm, now: time.Now}
}

// GetNutritionAnalysis - dipakai frontend saat halaman dibuka untuk menampilkan
// hasil yang sudah ada, tanpa memicu pemanggilan LLM yang berbiaya.
func (s *service) GetNutritionAnalysis(userID, submissionRef string) (*NutritionAnalysisResult, error) {
	sub, err := s.findSubmission(userID, submissionRef)
	if err != nil {
		return nil, err
	}

	stored, err := s.repo.FindBySubmissionID(sub.ID)
	if err != nil {
		return nil, err
	}
	if stored == nil {
		return nil, errNotAnalyzed()
	}

	result, ok := s.fromStored(sub, userID, stored)
	if !ok {
		return nil, errNotAnalyzed()
	}
	return result, nil
}

// AnalyzeNutrition - analisis gizi sebuah submission.
// Urutan: cek kepemilikan → pakai hasil tersimpan bila ada → hitung fakta
// deterministik → minta LLM menarasikan → validasi → simpan.
func (s *service) AnalyzeNutrition(ctx context.Context, userID string, req NutritionAnalysisRequest) (*NutritionAnalysisResult, error) {
	sub, err := s.findSubmission(userID, req.SubmissionID)
	if err != nil {
		return nil, err
	}

	// Satu submission hanya dianalisis satu kali pada satu waktu: klik ganda atau
	// dua tab tidak menghasilkan dua pemanggilan LLM — yang kedua menunggu lalu
	// menerima hasil yang baru disimpan.
	unlock := s.locks.Lock(sub.ID)
	defer unlock()

	stored, err := s.repo.FindBySubmissionID(sub.ID)
	if err != nil {
		return nil, err
	}
	if stored != nil {
		withinCooldown := s.now().Sub(generatedAt(stored)) < refreshCooldown
		if !req.ForceRefresh || withinCooldown {
			if result, ok := s.fromStored(sub, userID, stored); ok {
				return result, nil
			}
			// Baris tersimpan rusak: lanjut membuat analisis baru.
		}
	}

	actx := buildContext(sub, s.repo.GetProfile(userID, sub.SubmittedAt))
	if actx.Coverage.FoodCount == 0 {
		return nil, errNoNutritionData()
	}
	if !s.llm.Configured() {
		return nil, mapLLMError(&groq.Error{Kind: groq.ErrNotConfigured})
	}

	ctx, cancel := context.WithTimeout(ctx, generateBudget)
	defer cancel()

	userPrompt := actx.userPrompt()
	var (
		completion *groq.Result
		parsed     narrative
		valid      bool
	)
	for attempt := 1; attempt <= maxGenerateAttempts && !valid; attempt++ {
		completion, err = s.llm.CompleteJSON(ctx, systemPrompt, userPrompt)
		if err != nil {
			return nil, mapLLMError(err)
		}
		parsed, valid = parseNarrative(completion.Content)
		if !valid {
			log.Printf("[ai] jawaban LLM tidak valid (percobaan %d) untuk submission %s", attempt, sub.ID)
		}
	}
	if !valid {
		return nil, errInvalidResponse()
	}

	data := compose(actx, parsed)
	generated := s.now()

	// Kegagalan menyimpan tidak boleh membuang hasil yang sudah dibayar dan
	// sudah di tangan: user tetap menerimanya, hanya saja tidak tersimpan.
	if err := s.repo.Save(newLog(sub.ID, userPrompt, completion, data)); err != nil {
		log.Printf("[ai] gagal menyimpan hasil analisis submission %s: %v", sub.ID, err)
	}

	return &NutritionAnalysisResult{
		Source: "groq",
		Data:   data,
		Meta:   AnalysisMeta{Model: completion.ModelUsed, PromptVersion: PromptVersion, GeneratedAt: generated},
	}, nil
}

// findSubmission - ambil submission milik user; error seragam bila tidak ada atau bukan miliknya
func (s *service) findSubmission(userID, ref string) (*submission.SurveySubmission, error) {
	if userID == "" || ref == "" {
		return nil, errNotFound()
	}
	sub, err := s.repo.FindOwnedSubmission(ref, userID)
	if err != nil || sub == nil {
		return nil, errNotFound()
	}
	return sub, nil
}

// fromStored - bangun response dari baris tersimpan.
// Baris lama (sebelum kolom result_data ada) hanya menyimpan narasi mentah LLM,
// jadi fakta deterministiknya dihitung ulang dari submission.
func (s *service) fromStored(sub *submission.SurveySubmission, userID string, stored *AIResultLog) (*NutritionAnalysisResult, bool) {
	var data NutritionAnalysisData

	if stored.ResultData != nil && *stored.ResultData != "" {
		if err := json.Unmarshal([]byte(*stored.ResultData), &data); err != nil {
			return nil, false
		}
	} else {
		parsed, ok := parseNarrative(stored.RawResponse)
		if !ok {
			return nil, false
		}
		data = compose(buildContext(sub, s.repo.GetProfile(userID, sub.SubmittedAt)), parsed)
	}

	return &NutritionAnalysisResult{
		Source: "cache",
		Data:   data,
		Meta: AnalysisMeta{
			Model:         stored.ModelUsed,
			PromptVersion: stored.PromptVersion,
			GeneratedAt:   generatedAt(stored),
		},
	}, true
}

// generatedAt - waktu hasil terakhir dibuat; baris lama belum punya updated_at
func generatedAt(stored *AIResultLog) time.Time {
	if !stored.UpdatedAt.IsZero() {
		return stored.UpdatedAt
	}
	return stored.CreatedAt
}

// newLog - susun baris ai_result_logs dari satu hasil analisis
func newLog(submissionID, userPrompt string, completion *groq.Result, data NutritionAnalysisData) *AIResultLog {
	resultBytes, _ := json.Marshal(data)
	resultJSON := string(resultBytes)

	return &AIResultLog{
		ID:               uuid.New().String(),
		SubmissionID:     submissionID,
		InputPayload:     userPrompt,
		RawResponse:      completion.Content,
		ResultData:       &resultJSON,
		OverallStatus:    data.OverallStatus,
		ModelUsed:        completion.ModelUsed,
		PromptVersion:    PromptVersion,
		TokenUsed:        &completion.TotalTokens,
		PromptTokens:     &completion.PromptTokens,
		CompletionTokens: &completion.CompletionTokens,
		LatencyMs:        &completion.LatencyMs,
	}
}

// keyedMutex - mutex per kunci; entri dibuang begitu tidak ada lagi yang memakainya
type keyedMutex struct {
	mu      sync.Mutex
	entries map[string]*keyedEntry
}

type keyedEntry struct {
	mu   sync.Mutex
	refs int
}

// Lock - kunci berdasarkan key dan kembalikan fungsi pembukanya
func (k *keyedMutex) Lock(key string) func() {
	k.mu.Lock()
	if k.entries == nil {
		k.entries = map[string]*keyedEntry{}
	}
	entry, ok := k.entries[key]
	if !ok {
		entry = &keyedEntry{}
		k.entries[key] = entry
	}
	entry.refs++
	k.mu.Unlock()

	entry.mu.Lock()
	return func() {
		entry.mu.Unlock()
		k.mu.Lock()
		entry.refs--
		if entry.refs == 0 {
			delete(k.entries, key)
		}
		k.mu.Unlock()
	}
}
