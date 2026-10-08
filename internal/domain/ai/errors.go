package ai

import (
	"atlas_food/internal/pkg/groq"
	"atlas_food/internal/pkg/utils"
	"net/http"
)

// Kode error endpoint AI. Frontend memilih tampilan dan aksi berdasarkan kode ini,
// jadi perubahan di sini adalah perubahan kontrak API.
const (
	CodeNotFound        = "NOT_FOUND"
	CodeNotAnalyzed     = "AI_NOT_ANALYZED"
	CodeNoNutritionData = "AI_NO_NUTRITION_DATA"
	CodeNotConfigured   = "AI_NOT_CONFIGURED"
	CodeRateLimited     = "AI_RATE_LIMITED"
	CodeTimeout         = "AI_TIMEOUT"
	CodeUnavailable     = "AI_UNAVAILABLE"
	CodeInvalidResponse = "AI_INVALID_RESPONSE"
)

// errNotFound - pesan seragam untuk "tidak ada" dan "bukan milik Anda",
// agar keberadaan submission orang lain tidak bisa ditebak dari response.
func errNotFound() *utils.AppError {
	return utils.NewAppError(http.StatusNotFound, CodeNotFound, "Laporan tidak ditemukan atau bukan milik Anda")
}

func errNotAnalyzed() *utils.AppError {
	return utils.NewAppError(http.StatusNotFound, CodeNotAnalyzed, "Laporan ini belum pernah dianalisis")
}

func errNoNutritionData() *utils.AppError {
	return utils.NewAppError(http.StatusUnprocessableEntity, CodeNoNutritionData,
		"Laporan ini belum memuat makanan dengan nilai gizi, sehingga belum dapat dianalisis")
}

func errInvalidResponse() *utils.AppError {
	return utils.NewAppError(http.StatusBadGateway, CodeInvalidResponse,
		"AI memberikan jawaban yang tidak dapat diproses. Silakan coba lagi")
}

// mapLLMError - terjemahkan kegagalan penyedia LLM menjadi error API.
// Detail mentah dari penyedia tidak pernah ikut dikirim ke klien.
func mapLLMError(err error) *utils.AppError {
	switch groq.KindOf(err) {
	case groq.ErrNotConfigured, groq.ErrRejected:
		// Kunci salah atau model tidak tersedia sama-sama masalah konfigurasi
		// server: mencoba lagi dari sisi user tidak akan membantu.
		return utils.NewAppError(http.StatusServiceUnavailable, CodeNotConfigured,
			"Layanan analisis AI belum tersedia. Hubungi pengelola survei")
	case groq.ErrRateLimited:
		return utils.NewAppError(http.StatusTooManyRequests, CodeRateLimited,
			"Layanan AI sedang ramai. Coba lagi dalam beberapa saat")
	case groq.ErrTimeout:
		return utils.NewAppError(http.StatusGatewayTimeout, CodeTimeout,
			"Analisis AI memakan waktu terlalu lama. Silakan coba lagi")
	case groq.ErrEmptyResponse:
		return errInvalidResponse()
	default:
		return utils.NewAppError(http.StatusServiceUnavailable, CodeUnavailable,
			"Layanan AI sedang tidak dapat dihubungi. Silakan coba lagi")
	}
}
