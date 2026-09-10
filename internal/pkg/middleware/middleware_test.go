package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func newCORSEngine() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(CORS([]string{"https://atlas.example.com", "http://localhost:3000"}))
	r.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "pong") })
	return r
}

// TestCORSAllowsListedOrigin - origin yang terdaftar harus dapat header lengkap.
func TestCORSAllowsListedOrigin(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set("Origin", "https://atlas.example.com")
	newCORSEngine().ServeHTTP(w, req)

	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "https://atlas.example.com" {
		t.Fatalf("Allow-Origin = %q, mau origin yang terdaftar", got)
	}
	if got := w.Header().Get("Access-Control-Allow-Credentials"); got != "true" {
		t.Fatalf("Allow-Credentials = %q, mau true", got)
	}
	if got := w.Header().Get("Vary"); got != "Origin" {
		t.Fatalf("Vary = %q, mau Origin", got)
	}
}

// TestCORSRejectsUnknownOrigin - regresi: middleware ini dulu memantulkan
// Origin APA PUN sambil menyalakan Allow-Credentials, sehingga situs mana pun
// bisa memanggil API atas nama user yang sedang login.
func TestCORSRejectsUnknownOrigin(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set("Origin", "https://situs-jahat.example")
	newCORSEngine().ServeHTTP(w, req)

	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("origin asing tidak boleh dipantulkan, dapat %q", got)
	}
	if got := w.Header().Get("Access-Control-Allow-Credentials"); got != "" {
		t.Fatalf("Allow-Credentials tidak boleh di-set untuk origin asing, dapat %q", got)
	}
}

// TestCORSPreflightFromUnknownOriginIsBlocked - preflight dari origin asing ditolak.
func TestCORSPreflightFromUnknownOriginIsBlocked(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodOptions, "/ping", nil)
	req.Header.Set("Origin", "https://situs-jahat.example")
	newCORSEngine().ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("preflight origin asing = %d, mau %d", w.Code, http.StatusForbidden)
	}
}

// TestCORSWithoutOriginPassesThrough - curl/health check tanpa Origin tetap jalan.
func TestCORSWithoutOriginPassesThrough(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	newCORSEngine().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("request tanpa Origin = %d, mau 200", w.Code)
	}
}

// TestRedactSensitiveQuery - JWT lewat query pada handshake WebSocket, jadi
// nilainya tidak boleh mendarat utuh di log server.
func TestRedactSensitiveQuery(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"token=rahasia.jwt.value", "token=REDACTED"},
		{"token=abc&invite=xyz", "token=REDACTED&invite=REDACTED"},
		{"page=2&limit=10", "page=2&limit=10"},
		{"Token=abc", "Token=REDACTED"},
		{"refresh_token=abc&q=nasi", "refresh_token=REDACTED&q=nasi"},
	}

	for _, c := range cases {
		if got := redactSensitiveQuery(c.in); got != c.want {
			t.Errorf("redactSensitiveQuery(%q) = %q, mau %q", c.in, got, c.want)
		}
	}
}
