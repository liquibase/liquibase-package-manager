package utils

import (
	"archive/zip"
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/go-github/v39/github"
)

// TestNewDownloadUtil tests the creation of a new DownloadUtil
func TestNewDownloadUtil(t *testing.T) {
	du := NewDownloadUtil()
	if du == nil {
		t.Fatal("NewDownloadUtil returned nil")
	}
	if du.client == nil {
		t.Fatal("DownloadUtil client is nil")
	}
	if du.client.Timeout == 0 {
		t.Fatal("DownloadUtil client timeout not set")
	}
}

// TestProgressWriter tests the progress writer functionality
func TestProgressWriter(t *testing.T) {
	var progressCalls []struct {
		current, total int64
	}

	onProgress := func(current, total int64) {
		progressCalls = append(progressCalls, struct {
			current, total int64
		}{current, total})
	}

	// Create a test writer (string builder)
	var output strings.Builder
	pw := NewProgressWriter(&output, 100, onProgress)

	// Write some data
	data := []byte("Hello, World!")
	n, err := pw.Write(data)
	if err != nil {
		t.Fatalf("ProgressWriter.Write failed: %v", err)
	}
	if n != len(data) {
		t.Fatalf("Expected to write %d bytes, wrote %d", len(data), n)
	}

	// Check that progress callback was called
	if len(progressCalls) != 1 {
		t.Fatalf("Expected 1 progress call, got %d", len(progressCalls))
	}
	if progressCalls[0].current != int64(len(data)) {
		t.Fatalf("Expected current to be %d, got %d", len(data), progressCalls[0].current)
	}
	if progressCalls[0].total != 100 {
		t.Fatalf("Expected total to be 100, got %d", progressCalls[0].total)
	}

	// Check that data was written correctly
	if output.String() != string(data) {
		t.Fatalf("Expected output to be %q, got %q", string(data), output.String())
	}
}

// TestDownloadWithProgress tests downloading with progress indication
func TestDownloadWithProgress(t *testing.T) {
	// Create test data
	testData := "This is test content for download"
	
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check User-Agent header
		if !strings.Contains(r.Header.Get("User-Agent"), "liquibase-package-manager") {
			t.Errorf("Expected User-Agent to contain 'liquibase-package-manager', got %q", r.Header.Get("User-Agent"))
		}
		
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(testData)))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(testData))
	}))
	defer server.Close()

	// Create temporary directory
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test-download.txt")

	// Create download util and download
	du := NewDownloadUtil()
	err := du.DownloadWithProgress(server.URL, testFile)
	if err != nil {
		t.Fatalf("DownloadWithProgress failed: %v", err)
	}

	// Verify file was created and has correct content
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read downloaded file: %v", err)
	}
	if string(content) != testData {
		t.Fatalf("Expected content %q, got %q", testData, string(content))
	}
}

// TestDownloadWithProgressErrors tests error conditions
func TestDownloadWithProgressErrors(t *testing.T) {
	du := NewDownloadUtil()

	tests := []struct {
		name     string
		url      string
		filepath string
		wantErr  bool
	}{
		{"empty URL", "", "/tmp/test", true},
		{"empty filepath", "http://example.com", "", true},
		{"invalid URL", "not-a-url", "/tmp/test", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := du.DownloadWithProgress(tt.url, tt.filepath)
			if (err != nil) != tt.wantErr {
				t.Errorf("DownloadWithProgress() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestDownloadWithProgressHTTPError tests HTTP error responses
func TestDownloadWithProgressHTTPError(t *testing.T) {
	// Create test server that returns 404
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
	}))
	defer server.Close()

	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test-download.txt")

	du := NewDownloadUtil()
	err := du.DownloadWithProgress(server.URL, testFile)
	if err == nil {
		t.Fatal("Expected error for 404 response, got nil")
	}
	if !strings.Contains(err.Error(), "404") {
		t.Fatalf("Expected error to mention 404, got: %v", err)
	}
}

// TestVerifyChecksum tests checksum verification
func TestVerifyChecksum(t *testing.T) {
	// Create test file with known content
	testData := "Hello, World!"
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	
	err := os.WriteFile(testFile, []byte(testData), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Calculate expected SHA256 checksum
	hash := sha256.Sum256([]byte(testData))
	expectedSHA256 := fmt.Sprintf("%x", hash)

	du := NewDownloadUtil()

	// Test successful verification
	err = du.VerifyChecksum(testFile, expectedSHA256, "sha256")
	if err != nil {
		t.Fatalf("VerifyChecksum failed: %v", err)
	}

	// Test failed verification
	err = du.VerifyChecksum(testFile, "wrong-checksum", "sha256")
	if err == nil {
		t.Fatal("Expected error for wrong checksum, got nil")
	}
}

// TestVerifyChecksumErrors tests error conditions for checksum verification
func TestVerifyChecksumErrors(t *testing.T) {
	du := NewDownloadUtil()

	tests := []struct {
		name             string
		filepath         string
		expectedChecksum string
		algorithm        string
		wantErr          bool
	}{
		{"empty filepath", "", "checksum", "sha256", true},
		{"empty checksum", "/tmp/test", "", "sha256", true},
		{"empty algorithm", "/tmp/test", "checksum", "", true},
		{"unsupported algorithm", "/tmp/test", "checksum", "md5", true},
		{"nonexistent file", "/nonexistent/file", "checksum", "sha256", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := du.VerifyChecksum(tt.filepath, tt.expectedChecksum, tt.algorithm)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyChecksum() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestDownloadToTemp tests downloading to temporary directory
func TestDownloadToTemp(t *testing.T) {
	testData := "Temporary download test content"
	
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(testData))
	}))
	defer server.Close()

	du := NewDownloadUtil()
	tempFile, err := du.DownloadToTemp(server.URL, "test-temp.txt")
	if err != nil {
		t.Fatalf("DownloadToTemp failed: %v", err)
	}
	defer os.Remove(tempFile)

	// Verify file exists and has correct content
	content, err := os.ReadFile(tempFile)
	if err != nil {
		t.Fatalf("Failed to read temp file: %v", err)
	}
	if string(content) != testData {
		t.Fatalf("Expected content %q, got %q", testData, string(content))
	}

	// Verify file is in temp directory
	if !strings.Contains(tempFile, "lpm-downloads") {
		t.Fatalf("Expected temp file to be in lpm-downloads directory, got %s", tempFile)
	}
}

// TestDownloadToTempErrors tests error conditions for temp downloads
func TestDownloadToTempErrors(t *testing.T) {
	du := NewDownloadUtil()

	tests := []struct {
		name     string
		url      string
		filename string
		wantErr  bool
	}{
		{"empty URL", "", "test.txt", true},
		{"empty filename", "http://example.com", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := du.DownloadToTemp(tt.url, tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("DownloadToTemp() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestDownloadBinaryUpdate tests binary update download functionality
func TestDownloadBinaryUpdate(t *testing.T) {
	// Create a test ZIP file with a binary
	tempDir := t.TempDir()
	zipPath := filepath.Join(tempDir, "test-release.zip")
	
	// Create ZIP file with test binary
	err := createTestZipWithBinary(zipPath, "lpm", "fake binary content")
	if err != nil {
		t.Fatalf("Failed to create test ZIP: %v", err)
	}

	// Read ZIP content for server response
	zipData, err := os.ReadFile(zipPath)
	if err != nil {
		t.Fatalf("Failed to read test ZIP: %v", err)
	}

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/zip")
		w.WriteHeader(http.StatusOK)
		w.Write(zipData)
	}))
	defer server.Close()

	// Create mock GitHub release with asset
	downloadURL := server.URL
	assetName := "lpm-linux.zip"
	asset := &github.ReleaseAsset{
		Name:               &assetName,
		BrowserDownloadURL: &downloadURL,
	}
	
	release := &github.RepositoryRelease{
		Assets: []*github.ReleaseAsset{asset},
	}

	du := NewDownloadUtil()
	extractedPath, err := du.DownloadBinaryUpdate(release, "linux")
	if err != nil {
		t.Fatalf("DownloadBinaryUpdate failed: %v", err)
	}
	defer os.Remove(extractedPath)

	// Verify extracted binary exists
	if _, err := os.Stat(extractedPath); os.IsNotExist(err) {
		t.Fatalf("Extracted binary does not exist: %s", extractedPath)
	}

	// Verify binary content
	content, err := os.ReadFile(extractedPath)
	if err != nil {
		t.Fatalf("Failed to read extracted binary: %v", err)
	}
	if string(content) != "fake binary content" {
		t.Fatalf("Expected binary content 'fake binary content', got %q", string(content))
	}
}

// TestDownloadBinaryUpdateErrors tests error conditions
func TestDownloadBinaryUpdateErrors(t *testing.T) {
	du := NewDownloadUtil()

	// Test nil release
	_, err := du.DownloadBinaryUpdate(nil, "linux")
	if err == nil {
		t.Fatal("Expected error for nil release, got nil")
	}

	// Test release with no suitable assets
	release := &github.RepositoryRelease{
		Assets: []*github.ReleaseAsset{},
	}
	_, err = du.DownloadBinaryUpdate(release, "linux")
	if err == nil {
		t.Fatal("Expected error for release with no assets, got nil")
	}
}

// TestDownloadWithChecksum tests combined download and checksum verification
func TestDownloadWithChecksum(t *testing.T) {
	testData := "Test content for checksum verification"
	hash := sha256.Sum256([]byte(testData))
	expectedChecksum := fmt.Sprintf("%x", hash)

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(testData))
	}))
	defer server.Close()

	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test-checksum.txt")

	du := NewDownloadUtil()
	err := du.DownloadWithChecksum(server.URL, testFile, expectedChecksum, "sha256")
	if err != nil {
		t.Fatalf("DownloadWithChecksum failed: %v", err)
	}

	// Verify file exists and has correct content
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read downloaded file: %v", err)
	}
	if string(content) != testData {
		t.Fatalf("Expected content %q, got %q", testData, string(content))
	}
}

// TestDownloadWithChecksumFailure tests checksum verification failure
func TestDownloadWithChecksumFailure(t *testing.T) {
	testData := "Test content"
	wrongChecksum := "wrong-checksum-value"

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(testData))
	}))
	defer server.Close()

	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test-checksum-fail.txt")

	du := NewDownloadUtil()
	err := du.DownloadWithChecksum(server.URL, testFile, wrongChecksum, "sha256")
	if err == nil {
		t.Fatal("Expected error for wrong checksum, got nil")
	}

	// Verify file was cleaned up on checksum failure
	if _, err := os.Stat(testFile); !os.IsNotExist(err) {
		t.Fatal("Expected file to be cleaned up on checksum failure")
	}
}

// TestCleanupTempDownloads tests temporary file cleanup
func TestCleanupTempDownloads(t *testing.T) {
	// Create temp download directory with test files
	tempDir := filepath.Join(os.TempDir(), "lpm-downloads")
	err := os.MkdirAll(tempDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create old and new files
	oldFile := filepath.Join(tempDir, "old-file.txt")
	newFile := filepath.Join(tempDir, "new-file.txt")
	
	err = os.WriteFile(oldFile, []byte("old"), 0644)
	if err != nil {
		t.Fatalf("Failed to create old file: %v", err)
	}
	
	err = os.WriteFile(newFile, []byte("new"), 0644)
	if err != nil {
		t.Fatalf("Failed to create new file: %v", err)
	}

	// Make old file actually old
	oldTime := time.Now().Add(-2 * time.Hour)
	err = os.Chtimes(oldFile, oldTime, oldTime)
	if err != nil {
		t.Fatalf("Failed to change old file time: %v", err)
	}

	// Clean up files older than 1 hour
	err = CleanupTempDownloads(1 * time.Hour)
	if err != nil {
		t.Fatalf("CleanupTempDownloads failed: %v", err)
	}

	// Verify old file was removed and new file remains
	if _, err := os.Stat(oldFile); !os.IsNotExist(err) {
		t.Fatal("Expected old file to be removed")
	}
	if _, err := os.Stat(newFile); os.IsNotExist(err) {
		t.Fatal("Expected new file to remain")
	}
}

// TestFormatBytes tests the byte formatting function
func TestFormatBytes(t *testing.T) {
	tests := []struct {
		bytes    int64
		expected string
	}{
		{512, "512 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
		{1073741824, "1.0 GB"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d bytes", tt.bytes), func(t *testing.T) {
			result := formatBytes(tt.bytes)
			if result != tt.expected {
				t.Errorf("formatBytes(%d) = %q, expected %q", tt.bytes, result, tt.expected)
			}
		})
	}
}

// Helper function to create a test ZIP file with a binary
func createTestZipWithBinary(zipPath, binaryName, content string) error {
	file, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer file.Close()

	w := zip.NewWriter(file)
	defer w.Close()

	f, err := w.Create(binaryName)
	if err != nil {
		return err
	}

	_, err = f.Write([]byte(content))
	return err
}

// TestExtractBinaryFromZip tests ZIP extraction functionality
func TestExtractBinaryFromZip(t *testing.T) {
	tempDir := t.TempDir()
	zipPath := filepath.Join(tempDir, "test.zip")
	
	// Create test ZIP with binary
	err := createTestZipWithBinary(zipPath, "lpm", "test binary content")
	if err != nil {
		t.Fatalf("Failed to create test ZIP: %v", err)
	}

	du := NewDownloadUtil()
	extractedPath, err := du.extractBinaryFromZip(zipPath, "linux")
	if err != nil {
		t.Fatalf("extractBinaryFromZip failed: %v", err)
	}
	defer os.Remove(extractedPath)

	// Verify extracted file
	content, err := os.ReadFile(extractedPath)
	if err != nil {
		t.Fatalf("Failed to read extracted file: %v", err)
	}
	if string(content) != "test binary content" {
		t.Fatalf("Expected content 'test binary content', got %q", string(content))
	}

	// Verify file is executable (on Unix-like systems)
	info, err := os.Stat(extractedPath)
	if err != nil {
		t.Fatalf("Failed to stat extracted file: %v", err)
	}
	if info.Mode()&0111 == 0 {
		t.Fatal("Expected extracted binary to be executable")
	}
}

// TestExtractBinaryFromZipWindows tests Windows binary extraction
func TestExtractBinaryFromZipWindows(t *testing.T) {
	tempDir := t.TempDir()
	zipPath := filepath.Join(tempDir, "test-windows.zip")
	
	// Create test ZIP with Windows binary
	err := createTestZipWithBinary(zipPath, "lpm.exe", "windows binary content")
	if err != nil {
		t.Fatalf("Failed to create test ZIP: %v", err)
	}

	du := NewDownloadUtil()
	extractedPath, err := du.extractBinaryFromZip(zipPath, "windows")
	if err != nil {
		t.Fatalf("extractBinaryFromZip failed: %v", err)
	}
	defer os.Remove(extractedPath)

	// Verify extracted file has correct name
	if !strings.HasSuffix(extractedPath, "lpm.exe") {
		t.Fatalf("Expected extracted path to end with 'lpm.exe', got %s", extractedPath)
	}

	// Verify content
	content, err := os.ReadFile(extractedPath)
	if err != nil {
		t.Fatalf("Failed to read extracted file: %v", err)
	}
	if string(content) != "windows binary content" {
		t.Fatalf("Expected content 'windows binary content', got %q", string(content))
	}
}