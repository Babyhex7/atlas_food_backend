package groq

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
)

const okBody = "{\"model\":\"m-1\",\"choices\":[{\"message\":{\"content\":\"```json\\n{\\\"a\\\":1}\\n```\"}}],\"usage\":{\"prompt_tokens\":3,\"completion_tokens\":4,\"total_tokens\":7}}"

func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return NewClient("key", "m-1", srv.URL, 5, 100)
}

func TestCompleteJSONSuccess(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer key" {
			t.Errorf("header Authorization salah: %q", r.Header.Get("Authorization"))
		}
		w.Write([]byte(okBody))
	})
	res, err := c.CompleteJSON(context.Background(), "sys", "user")
	if err != nil {
		t.Fatal(err)
	}
	if res.Content != `{"a":1}` {
		t.Errorf("pagar markdown harus dibuang, dapat %q", res.Content)
	}
	if res.TotalTokens != 7 || res.PromptTokens != 3 || res.ModelUsed != "m-1" {
		t.Errorf("metadata salah: %+v", res)
	}
}

func TestNotConfigured(t *testing.T) {
	c := NewClient("  ", "", "", 0, 0)
	if c.Configured() {
		t.Fatal("kunci kosong harus dianggap belum dikonfigurasi")
	}
	if _, err := c.CompleteJSON(context.Background(), "s", "u"); KindOf(err) != ErrNotConfigured {
		t.Errorf("mau ErrNotConfigured, dapat %v", err)
	}
}

func TestRetriesTransientThenSucceeds(t *testing.T) {
	var calls int32
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		if atomic.AddInt32(&calls, 1) == 1 {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.Write([]byte(okBody))
	})
	if _, err := c.CompleteJSON(context.Background(), "s", "u"); err != nil {
		t.Fatalf("429 sekali harus pulih lewat percobaan ulang: %v", err)
	}
	if calls != 2 {
		t.Errorf("mau 2 pemanggilan, dapat %d", calls)
	}
}

func TestRejectedIsNotRetriedAndHidesBody(t *testing.T) {
	var calls int32
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":{"message":"Invalid API Key gsk_rahasia"}}`))
	})
	_, err := c.CompleteJSON(context.Background(), "s", "u")
	if KindOf(err) != ErrRejected || calls != 1 {
		t.Errorf("401 tidak boleh diulang: kind=%s calls=%d", KindOf(err), calls)
	}
	if strings.Contains(err.Error(), "gsk_rahasia") {
		t.Error("isi balasan penyedia tidak boleh ikut di pesan error")
	}
}

func TestEmptyChoices(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(`{"choices":[]}`))
	})
	if _, err := c.CompleteJSON(context.Background(), "s", "u"); KindOf(err) != ErrEmptyResponse {
		t.Errorf("mau ErrEmptyResponse, dapat %v", err)
	}
}

func TestCancelledContextIsTimeout(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) { w.Write([]byte(okBody)) })
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := c.CompleteJSON(ctx, "s", "u"); KindOf(err) != ErrTimeout {
		t.Errorf("mau ErrTimeout, dapat %v", err)
	}
}
