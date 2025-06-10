package utils

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

// Test helper functions

// createTestBinary creates a test binary file with specified content
func createTestBinary(t *testing.T, dir, name, content string) string {
	path := filepath.Join(dir, name)
	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("Failed to create test binary: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		t.Fatalf("Failed to write test binary content: %v", err)
	}

	// Make executable on Unix systems
	if runtime.GOOS != "windows" {
		if err := os.Chmod(path, 0755); err != nil {
			t.Fatalf("Failed to make test binary executable: %v", err)
		}
	}

	return path
}

// readFileContent reads the entire content of a file
func readFileContent(t *testing.T, path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read file content: %v", err)
	}
	return string(content)
}

// assertFileExists checks that a file exists
func assertFileExists(t *testing.T, path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("Expected file to exist: %s", path)
	}
}

// assertFileNotExists checks that a file does not exist
func assertFileNotExists(t *testing.T, path string) {
	if _, err := os.Stat(path); err == nil {
		t.Fatalf("Expected file to not exist: %s", path)
	}
}

// assertFileExecutable checks that a file has executable permissions (Unix only)
func assertFileExecutable(t *testing.T, path string) {
	if runtime.GOOS == "windows" {
		return // Skip on Windows
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}

	if info.Mode()&0111 == 0 {
		t.Fatalf("Expected file to be executable: %s", path)
	}
}

func TestNewUpdateUtil(t *testing.T) {
	util := NewUpdateUtil()
	if util == nil {
		t.Fatal("NewUpdateUtil returned nil")
	}
	if util.platformUtil == nil {
		t.Fatal("UpdateUtil.platformUtil is nil")
	}
}

func TestBackupCurrentBinary(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "update_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	util := NewUpdateUtil()

	t.Run("successful backup", func(t *testing.T) {
		// Create test binary
		testBinary := createTestBinary(t, tempDir, "test_binary", "original content")

		// Create backup
		backupPath, err := util.BackupCurrentBinary(testBinary)
		if err != nil {
			t.Fatalf("BackupCurrentBinary failed: %v", err)
		}

		// Verify backup was created
		assertFileExists(t, backupPath)

		// Verify backup content matches original
		originalContent := readFileContent(t, testBinary)
		backupContent := readFileContent(t, backupPath)
		if originalContent != backupContent {
			t.Fatalf("Backup content mismatch. Original: %s, Backup: %s", originalContent, backupContent)
		}

		// Verify backup has correct permissions
		assertFileExecutable(t, backupPath)

		// Verify backup filename format
		if !strings.Contains(backupPath, ".backup.") {
			t.Fatalf("Backup path doesn't contain expected pattern: %s", backupPath)
		}
	})

	t.Run("empty executable path", func(t *testing.T) {
		_, err := util.BackupCurrentBinary("")
		if err == nil {
			t.Fatal("Expected error for empty executable path")
		}
		if !strings.Contains(err.Error(), "executable path cannot be empty") {
			t.Fatalf("Unexpected error message: %v", err)
		}
	})

	t.Run("non-existent binary", func(t *testing.T) {
		nonExistentPath := filepath.Join(tempDir, "non_existent_binary")
		_, err := util.BackupCurrentBinary(nonExistentPath)
		if err == nil {
			t.Fatal("Expected error for non-existent binary")
		}
		if !strings.Contains(err.Error(), "does not exist") {
			t.Fatalf("Unexpected error message: %v", err)
		}
	})
}

func TestReplaceCurrentBinary(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "update_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	util := NewUpdateUtil()

	t.Run("successful replacement", func(t *testing.T) {
		// Create original binary
		originalBinary := createTestBinary(t, tempDir, "original_binary", "original content")
		
		// Create new binary
		newBinary := createTestBinary(t, tempDir, "new_binary", "new content")

		// Replace binary
		err := util.ReplaceCurrentBinary(newBinary, originalBinary)
		if err != nil {
			t.Fatalf("ReplaceCurrentBinary failed: %v", err)
		}

		// Verify replacement succeeded
		content := readFileContent(t, originalBinary)
		if content != "new content" {
			t.Fatalf("Binary replacement failed. Expected 'new content', got: %s", content)
		}

		// Verify permissions are preserved
		assertFileExecutable(t, originalBinary)
	})

	t.Run("empty new binary path", func(t *testing.T) {
		originalBinary := createTestBinary(t, tempDir, "test_binary2", "content")
		err := util.ReplaceCurrentBinary("", originalBinary)
		if err == nil {
			t.Fatal("Expected error for empty new binary path")
		}
	})

	t.Run("empty executable path", func(t *testing.T) {
		newBinary := createTestBinary(t, tempDir, "new_binary2", "content")
		err := util.ReplaceCurrentBinary(newBinary, "")
		if err == nil {
			t.Fatal("Expected error for empty executable path")
		}
	})

	t.Run("non-existent new binary", func(t *testing.T) {
		originalBinary := createTestBinary(t, tempDir, "test_binary3", "content")
		nonExistentPath := filepath.Join(tempDir, "non_existent")
		err := util.ReplaceCurrentBinary(nonExistentPath, originalBinary)
		if err == nil {
			t.Fatal("Expected error for non-existent new binary")
		}
	})

	t.Run("non-existent current binary", func(t *testing.T) {
		newBinary := createTestBinary(t, tempDir, "new_binary3", "content")
		nonExistentPath := filepath.Join(tempDir, "non_existent_current")
		err := util.ReplaceCurrentBinary(newBinary, nonExistentPath)
		if err == nil {
			t.Fatal("Expected error for non-existent current binary")
		}
	})
}

func TestRestoreFromBackup(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "update_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	util := NewUpdateUtil()

	t.Run("successful restore", func(t *testing.T) {
		// Create original binary
		originalBinary := createTestBinary(t, tempDir, "test_binary", "original content")
		
		// Create backup
		backupPath, err := util.BackupCurrentBinary(originalBinary)
		if err != nil {
			t.Fatalf("Failed to create backup: %v", err)
		}

		// Modify original (simulate corruption)
		err = os.WriteFile(originalBinary, []byte("corrupted content"), 0755)
		if err != nil {
			t.Fatalf("Failed to simulate corruption: %v", err)
		}

		// Restore from backup
		err = util.RestoreFromBackup(backupPath, originalBinary)
		if err != nil {
			t.Fatalf("RestoreFromBackup failed: %v", err)
		}

		// Verify restoration
		content := readFileContent(t, originalBinary)
		if content != "original content" {
			t.Fatalf("Restore failed. Expected 'original content', got: %s", content)
		}

		// Verify permissions
		assertFileExecutable(t, originalBinary)
	})

	t.Run("empty backup path", func(t *testing.T) {
		originalBinary := createTestBinary(t, tempDir, "test_binary2", "content")
		err := util.RestoreFromBackup("", originalBinary)
		if err == nil {
			t.Fatal("Expected error for empty backup path")
		}
	})

	t.Run("empty executable path", func(t *testing.T) {
		backupBinary := createTestBinary(t, tempDir, "backup", "content")
		err := util.RestoreFromBackup(backupBinary, "")
		if err == nil {
			t.Fatal("Expected error for empty executable path")
		}
	})

	t.Run("non-existent backup", func(t *testing.T) {
		originalBinary := createTestBinary(t, tempDir, "test_binary3", "content")
		nonExistentBackup := filepath.Join(tempDir, "non_existent_backup")
		err := util.RestoreFromBackup(nonExistentBackup, originalBinary)
		if err == nil {
			t.Fatal("Expected error for non-existent backup")
		}
	})
}

func TestCleanupBackup(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "update_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	util := NewUpdateUtil()

	t.Run("successful cleanup", func(t *testing.T) {
		// Create test binary and backup
		testBinary := createTestBinary(t, tempDir, "test_binary", "content")
		backupPath, err := util.BackupCurrentBinary(testBinary)
		if err != nil {
			t.Fatalf("Failed to create backup: %v", err)
		}

		// Verify backup exists
		assertFileExists(t, backupPath)

		// Clean up backup
		err = util.CleanupBackup(backupPath)
		if err != nil {
			t.Fatalf("CleanupBackup failed: %v", err)
		}

		// Verify backup is removed
		assertFileNotExists(t, backupPath)
	})

	t.Run("cleanup non-existent backup", func(t *testing.T) {
		nonExistentBackup := filepath.Join(tempDir, "non_existent_backup")
		err := util.CleanupBackup(nonExistentBackup)
		if err != nil {
			t.Fatalf("CleanupBackup should not fail for non-existent backup: %v", err)
		}
	})

	t.Run("empty backup path", func(t *testing.T) {
		err := util.CleanupBackup("")
		if err == nil {
			t.Fatal("Expected error for empty backup path")
		}
	})
}

func TestCopyFile(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "update_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	util := NewUpdateUtil()

	t.Run("successful copy", func(t *testing.T) {
		// Create source file
		srcPath := createTestBinary(t, tempDir, "source", "test content")
		dstPath := filepath.Join(tempDir, "destination")

		// Copy file
		err := util.copyFile(srcPath, dstPath)
		if err != nil {
			t.Fatalf("copyFile failed: %v", err)
		}

		// Verify copy
		assertFileExists(t, dstPath)
		srcContent := readFileContent(t, srcPath)
		dstContent := readFileContent(t, dstPath)
		if srcContent != dstContent {
			t.Fatalf("File content mismatch after copy. Source: %s, Destination: %s", srcContent, dstContent)
		}

		// Verify permissions
		assertFileExecutable(t, dstPath)
	})

	t.Run("non-existent source", func(t *testing.T) {
		nonExistentSrc := filepath.Join(tempDir, "non_existent_source")
		dstPath := filepath.Join(tempDir, "destination2")
		err := util.copyFile(nonExistentSrc, dstPath)
		if err == nil {
			t.Fatal("Expected error for non-existent source file")
		}
	})
}

func TestPerformAtomicUpdate(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "update_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	util := NewUpdateUtil()

	t.Run("successful atomic update", func(t *testing.T) {
		// Create original binary
		originalBinary := createTestBinary(t, tempDir, "test_binary", "original content")
		
		// Create new binary
		newBinary := createTestBinary(t, tempDir, "new_binary", "new content")

		// We can't easily mock the platform util in the current structure, 
		// so we'll test the individual components of the atomic update manually

		// Test the backup -> replace -> verify -> cleanup cycle manually
		// Step 1: Backup
		backupPath, err := util.BackupCurrentBinary(originalBinary)
		if err != nil {
			t.Fatalf("Backup failed: %v", err)
		}

		// Step 2: Replace
		err = util.ReplaceCurrentBinary(newBinary, originalBinary)
		if err != nil {
			t.Fatalf("Replace failed: %v", err)
		}

		// Step 3: Verify content changed
		content := readFileContent(t, originalBinary)
		if content != "new content" {
			t.Fatalf("Update failed. Expected 'new content', got: %s", content)
		}

		// Step 4: Cleanup
		err = util.CleanupBackup(backupPath)
		if err != nil {
			t.Fatalf("Cleanup failed: %v", err)
		}

		assertFileNotExists(t, backupPath)
	})

	t.Run("empty new binary path", func(t *testing.T) {
		err := util.PerformAtomicUpdate("")
		if err == nil {
			t.Fatal("Expected error for empty new binary path")
		}
	})
}

func TestVerifyNewBinary(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "update_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	util := NewUpdateUtil()

	t.Run("valid binary", func(t *testing.T) {
		validBinary := createTestBinary(t, tempDir, "valid_binary", "valid content")
		err := util.verifyNewBinary(validBinary)
		if err != nil {
			t.Fatalf("verifyNewBinary failed for valid binary: %v", err)
		}
	})

	t.Run("empty binary", func(t *testing.T) {
		emptyBinary := createTestBinary(t, tempDir, "empty_binary", "")
		err := util.verifyNewBinary(emptyBinary)
		if err == nil {
			t.Fatal("Expected error for empty binary")
		}
		if !strings.Contains(err.Error(), "empty") {
			t.Fatalf("Unexpected error message: %v", err)
		}
	})

	t.Run("non-existent binary", func(t *testing.T) {
		nonExistentBinary := filepath.Join(tempDir, "non_existent")
		err := util.verifyNewBinary(nonExistentBinary)
		if err == nil {
			t.Fatal("Expected error for non-existent binary")
		}
	})

	t.Run("non-executable binary on Unix", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("Skipping executable test on Windows")
		}

		nonExecBinary := createTestBinary(t, tempDir, "non_exec_binary", "content")
		// Remove executable permissions
		err := os.Chmod(nonExecBinary, 0644)
		if err != nil {
			t.Fatalf("Failed to remove executable permissions: %v", err)
		}

		err = util.verifyNewBinary(nonExecBinary)
		if err == nil {
			t.Fatal("Expected error for non-executable binary")
		}
		if !strings.Contains(err.Error(), "not executable") {
			t.Fatalf("Unexpected error message: %v", err)
		}
	})
}

func TestCleanupTempFiles(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "update_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	util := NewUpdateUtil()

	t.Run("cleanup temporary files", func(t *testing.T) {
		// Create main binary
		mainBinary := createTestBinary(t, tempDir, "main_binary", "content")

		// Create temporary files
		tempFiles := []string{
			mainBinary + ".new",
			mainBinary + ".old",
			mainBinary + ".temp",
			mainBinary + ".pending",
			mainBinary + ".update.bat",
		}

		for _, tempFile := range tempFiles {
			createTestBinary(t, filepath.Dir(tempFile), filepath.Base(tempFile), "temp content")
		}

		// Verify temp files exist
		for _, tempFile := range tempFiles {
			assertFileExists(t, tempFile)
		}

		// Cleanup
		err := util.CleanupTempFiles(mainBinary)
		if err != nil {
			t.Fatalf("CleanupTempFiles failed: %v", err)
		}

		// Verify temp files are removed
		for _, tempFile := range tempFiles {
			assertFileNotExists(t, tempFile)
		}
	})

	t.Run("empty executable path", func(t *testing.T) {
		err := util.CleanupTempFiles("")
		if err == nil {
			t.Fatal("Expected error for empty executable path")
		}
	})
}

func TestCleanupOldBackups(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "update_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	util := NewUpdateUtil()

	t.Run("cleanup old backups", func(t *testing.T) {
		mainBinary := filepath.Join(tempDir, "main_binary")
		
		// Create some backup files with different ages
		recentBackup := mainBinary + ".backup.20240101-120000"
		oldBackup := mainBinary + ".backup.20220101-120000"
		
		createTestBinary(t, filepath.Dir(recentBackup), filepath.Base(recentBackup), "recent backup")
		createTestBinary(t, filepath.Dir(oldBackup), filepath.Base(oldBackup), "old backup")

		// Make old backup actually old by changing its modification time
		oldTime := time.Now().Add(-30 * 24 * time.Hour) // 30 days ago
		err := os.Chtimes(oldBackup, oldTime, oldTime)
		if err != nil {
			t.Fatalf("Failed to set old time on backup: %v", err)
		}

		// Cleanup backups older than 7 days
		err = util.cleanupOldBackups(mainBinary, 7*24*time.Hour)
		if err != nil {
			t.Fatalf("cleanupOldBackups failed: %v", err)
		}

		// Recent backup should still exist, old backup should be removed
		assertFileExists(t, recentBackup)
		assertFileNotExists(t, oldBackup)
	})
}

// Benchmark tests

func BenchmarkBackupCurrentBinary(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "benchmark_*")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	util := NewUpdateUtil()
	testBinary := createBenchmarkBinary(b, tempDir, "benchmark_binary", strings.Repeat("x", 1024*1024)) // 1MB file

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		backupPath, err := util.BackupCurrentBinary(testBinary)
		if err != nil {
			b.Fatalf("Backup failed: %v", err)
		}
		os.Remove(backupPath) // Clean up for next iteration
	}
}

func BenchmarkReplaceCurrentBinary(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "benchmark_*")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	util := NewUpdateUtil()
	content := strings.Repeat("x", 1024*1024) // 1MB file

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		originalBinary := createBenchmarkBinary(b, tempDir, "original_"+string(rune(i)), content)
		newBinary := createBenchmarkBinary(b, tempDir, "new_"+string(rune(i)), content+"_new")
		b.StartTimer()

		err := util.ReplaceCurrentBinary(newBinary, originalBinary)
		if err != nil {
			b.Fatalf("Replace failed: %v", err)
		}
	}
}

// Helper function for benchmarks that need to create test binaries
func createBenchmarkBinary(tb testing.TB, dir, name, content string) string {
	path := filepath.Join(dir, name)
	file, err := os.Create(path)
	if err != nil {
		tb.Fatalf("Failed to create test binary: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		tb.Fatalf("Failed to write test binary content: %v", err)
	}

	// Make executable on Unix systems
	if runtime.GOOS != "windows" {
		if err := os.Chmod(path, 0755); err != nil {
			tb.Fatalf("Failed to make test binary executable: %v", err)
		}
	}

	return path
}