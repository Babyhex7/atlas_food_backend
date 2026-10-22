package storage

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLocalStorageUploadAndDelete(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "storage_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	svc := NewLocalStorage(tempDir)
	if svc.Driver() != "local" {
		t.Errorf("Expected driver 'local', got %s", svc.Driver())
	}

	content := "hello test world"
	reader := strings.NewReader(content)
	ctx := context.Background()

	url, err := svc.Upload(ctx, "testfolder", "hello.txt", reader, int64(len(content)), "text/plain")
	if err != nil {
		t.Fatalf("Upload failed: %v", err)
	}

	expectedURL := "/uploads/testfolder/hello.txt"
	if url != expectedURL {
		t.Errorf("Expected URL %s, got %s", expectedURL, url)
	}

	savedFile := filepath.Join(tempDir, "testfolder", "hello.txt")
	if _, err := os.Stat(savedFile); os.IsNotExist(err) {
		t.Errorf("File was not saved to disk at %s", savedFile)
	}

	// Test Delete
	if err := svc.Delete(ctx, url); err != nil {
		t.Errorf("Delete failed: %v", err)
	}

	if _, err := os.Stat(savedFile); !os.IsNotExist(err) {
		t.Errorf("File was not deleted from disk")
	}
}
