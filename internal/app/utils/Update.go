package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// UpdateUtil provides atomic binary replacement functionality for cross-platform self-updates
type UpdateUtil struct {
	platformUtil *PlatformUtil
}

// NewUpdateUtil creates a new UpdateUtil instance
func NewUpdateUtil() *UpdateUtil {
	return &UpdateUtil{
		platformUtil: &PlatformUtil{},
	}
}

// BackupCurrentBinary creates a backup of the current running binary
// Returns the backup file path or an error if the backup fails
func (u *UpdateUtil) BackupCurrentBinary(execPath string) (string, error) {
	if execPath == "" {
		return "", fmt.Errorf("executable path cannot be empty")
	}

	fmt.Printf("[DEBUG] Backup: Validating current binary at %s\n", execPath)
	// Verify the current binary exists and is accessible
	if _, err := os.Stat(execPath); os.IsNotExist(err) {
		fmt.Printf("[ERROR] Backup: Current binary does not exist: %s\n", execPath)
		return "", fmt.Errorf("current binary does not exist: %s", execPath)
	}

	// Generate backup filename with timestamp
	timestamp := time.Now().Format("20060102-150405")
	backupPath := fmt.Sprintf("%s.backup.%s", execPath, timestamp)
	fmt.Printf("[DEBUG] Backup: Creating backup at %s\n", backupPath)

	// On Windows, we need to handle potential file locking issues
	if runtime.GOOS == "windows" {
		fmt.Printf("[DEBUG] Backup: Using Windows-specific backup method\n")
		return u.backupBinaryWindows(execPath, backupPath)
	}

	// Unix-like systems: standard copy operation
	fmt.Printf("[DEBUG] Backup: Using Unix-style backup method\n")
	return u.backupBinaryUnix(execPath, backupPath)
}

// backupBinaryUnix performs backup on Unix-like systems
func (u *UpdateUtil) backupBinaryUnix(execPath, backupPath string) (string, error) {
	// Open source file
	src, err := os.Open(execPath)
	if err != nil {
		return "", fmt.Errorf("failed to open current binary for backup: %w", err)
	}
	defer src.Close()

	// Get source file info to preserve permissions
	srcInfo, err := src.Stat()
	if err != nil {
		return "", fmt.Errorf("failed to get current binary info: %w", err)
	}

	// Create backup file
	dst, err := os.Create(backupPath)
	if err != nil {
		return "", fmt.Errorf("failed to create backup file: %w", err)
	}
	defer dst.Close()

	// Copy file content
	_, err = io.Copy(dst, src)
	if err != nil {
		os.Remove(backupPath) // Clean up on failure
		return "", fmt.Errorf("failed to copy binary content to backup: %w", err)
	}

	// Preserve original permissions
	if err := os.Chmod(backupPath, srcInfo.Mode()); err != nil {
		os.Remove(backupPath) // Clean up on failure
		return "", fmt.Errorf("failed to set backup file permissions: %w", err)
	}

	return backupPath, nil
}

// backupBinaryWindows performs backup on Windows with special handling for file locking
func (u *UpdateUtil) backupBinaryWindows(execPath, backupPath string) (string, error) {
	// On Windows, the running executable might be locked
	// We'll try a few different approaches

	// First, try direct copy (works if binary isn't locked)
	if err := u.copyFile(execPath, backupPath); err == nil {
		return backupPath, nil
	}

	// If direct copy fails, try using Windows-specific methods
	// Create a temporary copy first
	tempPath := backupPath + ".temp"
	defer os.Remove(tempPath)

	// Try copying to temp location first
	if err := u.copyFile(execPath, tempPath); err != nil {
		return "", fmt.Errorf("failed to create backup on Windows (file may be locked): %w", err)
	}

	// Move temp to final backup location
	if err := os.Rename(tempPath, backupPath); err != nil {
		return "", fmt.Errorf("failed to finalize backup on Windows: %w", err)
	}

	return backupPath, nil
}

// ReplaceCurrentBinary atomically replaces the current binary with a new one
// This operation is designed to be atomic - either it succeeds completely or fails without corruption
func (u *UpdateUtil) ReplaceCurrentBinary(newBinaryPath, execPath string) error {
	if newBinaryPath == "" {
		return fmt.Errorf("new binary path cannot be empty")
	}
	if execPath == "" {
		return fmt.Errorf("executable path cannot be empty")
	}

	fmt.Printf("[DEBUG] Replace: Validating new binary at %s\n", newBinaryPath)
	// Verify new binary exists
	if _, err := os.Stat(newBinaryPath); os.IsNotExist(err) {
		fmt.Printf("[ERROR] Replace: New binary does not exist: %s\n", newBinaryPath)
		return fmt.Errorf("new binary does not exist: %s", newBinaryPath)
	}

	fmt.Printf("[DEBUG] Replace: Validating current binary at %s\n", execPath)
	// Verify current binary exists
	if _, err := os.Stat(execPath); os.IsNotExist(err) {
		fmt.Printf("[ERROR] Replace: Current binary does not exist: %s\n", execPath)
		return fmt.Errorf("current binary does not exist: %s", execPath)
	}

	// Platform-specific replacement
	if runtime.GOOS == "windows" {
		fmt.Printf("[DEBUG] Replace: Using Windows-specific replacement method\n")
		return u.replaceBinaryWindows(newBinaryPath, execPath)
	}

	fmt.Printf("[DEBUG] Replace: Using Unix-style replacement method\n")
	return u.replaceBinaryUnix(newBinaryPath, execPath)
}

// replaceBinaryUnix performs atomic replacement on Unix-like systems using rename
func (u *UpdateUtil) replaceBinaryUnix(newBinaryPath, execPath string) error {
	// Get current binary info to preserve permissions
	currentInfo, err := os.Stat(execPath)
	if err != nil {
		return fmt.Errorf("failed to get current binary info: %w", err)
	}

	// Ensure new binary has executable permissions
	if err := os.Chmod(newBinaryPath, currentInfo.Mode()); err != nil {
		return fmt.Errorf("failed to set permissions on new binary: %w", err)
	}

	// Create temporary file in the same directory as target (required for atomic rename)
	tempPath := execPath + ".new"
	
	// Copy new binary to temp location in target directory
	if err := u.copyFile(newBinaryPath, tempPath); err != nil {
		return fmt.Errorf("failed to copy new binary to target directory: %w", err)
	}
	defer os.Remove(tempPath) // Clean up temp file on any error

	// Ensure temp file has correct permissions
	if err := os.Chmod(tempPath, currentInfo.Mode()); err != nil {
		return fmt.Errorf("failed to set permissions on temporary binary: %w", err)
	}

	// Atomic rename operation
	if err := os.Rename(tempPath, execPath); err != nil {
		return fmt.Errorf("failed to atomically replace binary: %w", err)
	}

	return nil
}

// replaceBinaryWindows handles binary replacement on Windows with file locking considerations
func (u *UpdateUtil) replaceBinaryWindows(newBinaryPath, execPath string) error {
	// On Windows, we can't replace a running executable directly
	// We need to use a different strategy

	// Method 1: Try direct replacement (works if process exits immediately after)
	tempPath := execPath + ".new"
	if err := u.copyFile(newBinaryPath, tempPath); err != nil {
		return fmt.Errorf("failed to copy new binary: %w", err)
	}
	defer os.Remove(tempPath) // Clean up on failure

	// Try to replace directly
	if err := os.Rename(tempPath, execPath); err == nil {
		return nil // Success
	}

	// Method 2: Use Windows batch script approach for locked files
	return u.replaceBinaryWindowsBatch(newBinaryPath, execPath)
}

// replaceBinaryWindowsBatch uses a batch script to replace a locked binary on Windows
func (u *UpdateUtil) replaceBinaryWindowsBatch(newBinaryPath, execPath string) error {
	// Create a batch script that will:
	// 1. Wait for current process to exit
	// 2. Replace the binary
	// 3. Restart the new binary
	// 4. Clean up

	batchPath := execPath + ".update.bat"
	batchContent := fmt.Sprintf(`@echo off
timeout /t 2 /nobreak >nul
move /y "%s" "%s"
del "%%~f0"
`, newBinaryPath, execPath)

	// Write batch file
	if err := os.WriteFile(batchPath, []byte(batchContent), 0644); err != nil {
		return fmt.Errorf("failed to create update batch script: %w", err)
	}

	// The batch script will run after this process exits
	// For now, we just copy the file and let the caller handle process restart
	tempPath := execPath + ".pending"
	if err := u.copyFile(newBinaryPath, tempPath); err != nil {
		os.Remove(batchPath)
		return fmt.Errorf("failed to stage new binary: %w", err)
	}

	return fmt.Errorf("Windows update requires process restart - new binary staged at %s", tempPath)
}

// RestoreFromBackup restores the binary from a backup file
func (u *UpdateUtil) RestoreFromBackup(backupPath, execPath string) error {
	if backupPath == "" {
		return fmt.Errorf("backup path cannot be empty")
	}
	if execPath == "" {
		return fmt.Errorf("executable path cannot be empty")
	}

	// Verify backup exists
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("backup file does not exist: %s", backupPath)
	}

	// Use the same replacement logic as for updates
	if runtime.GOOS == "windows" {
		return u.replaceBinaryWindows(backupPath, execPath)
	}

	return u.replaceBinaryUnix(backupPath, execPath)
}

// CleanupBackup removes a backup file after successful update
func (u *UpdateUtil) CleanupBackup(backupPath string) error {
	if backupPath == "" {
		return fmt.Errorf("backup path cannot be empty")
	}

	// Check if backup exists
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return nil // Already cleaned up or never existed
	}

	// Remove backup file
	if err := os.Remove(backupPath); err != nil {
		return fmt.Errorf("failed to remove backup file %s: %w", backupPath, err)
	}

	return nil
}

// copyFile copies a file from src to dst with error handling
func (u *UpdateUtil) copyFile(src, dst string) error {
	// Open source file
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	// Get source file info
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to get source file info: %w", err)
	}

	// Create destination file
	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	// Copy file content
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		os.Remove(dst) // Clean up on failure
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	// Preserve file permissions
	if err := os.Chmod(dst, srcInfo.Mode()); err != nil {
		os.Remove(dst) // Clean up on failure
		return fmt.Errorf("failed to set file permissions: %w", err)
	}

	return nil
}

// PerformAtomicUpdate performs a complete atomic update operation with rollback capability
// This is a high-level function that orchestrates the entire update process
func (u *UpdateUtil) PerformAtomicUpdate(newBinaryPath string) error {
	if newBinaryPath == "" {
		return fmt.Errorf("new binary path cannot be empty")
	}

	fmt.Printf("[DEBUG] Update: Starting atomic update process\n")
	fmt.Printf("[DEBUG] Update: New binary path: %s\n", newBinaryPath)

	// Get current executable path
	fmt.Printf("[DEBUG] Update: Detecting current executable path\n")
	execPath, err := u.platformUtil.GetExecutablePath()
	if err != nil {
		fmt.Printf("[ERROR] Update: Failed to get executable path: %v\n", err)
		return fmt.Errorf("failed to get current executable path: %w", err)
	}
	fmt.Printf("[DEBUG] Update: Current executable: %s\n", execPath)

	// Step 1: Create backup
	fmt.Printf("[DEBUG] Update: Creating backup of current binary\n")
	backupPath, err := u.BackupCurrentBinary(execPath)
	if err != nil {
		fmt.Printf("[ERROR] Update: Backup failed: %v\n", err)
		return fmt.Errorf("failed to backup current binary: %w", err)
	}
	fmt.Printf("[DEBUG] Update: Backup created at: %s\n", backupPath)

	// Step 2: Attempt replacement
	fmt.Printf("[DEBUG] Update: Replacing current binary with new version\n")
	err = u.ReplaceCurrentBinary(newBinaryPath, execPath)
	if err != nil {
		fmt.Printf("[ERROR] Update: Binary replacement failed: %v\n", err)
		fmt.Printf("[DEBUG] Update: Attempting rollback to previous version\n")
		// Rollback on failure
		if rollbackErr := u.RestoreFromBackup(backupPath, execPath); rollbackErr != nil {
			fmt.Printf("[ERROR] Update: Rollback also failed: %v\n", rollbackErr)
			return fmt.Errorf("update failed and rollback also failed: update error: %w, rollback error: %w", err, rollbackErr)
		}
		fmt.Printf("[DEBUG] Update: Rollback completed successfully\n")
		// Clean up backup after successful rollback
		u.CleanupBackup(backupPath)
		return fmt.Errorf("update failed but was rolled back successfully: %w", err)
	}
	fmt.Printf("[DEBUG] Update: Binary replacement completed\n")

	// Step 3: Verify new binary works (basic check)
	fmt.Printf("[DEBUG] Update: Verifying new binary functionality\n")
	if err := u.verifyNewBinary(execPath); err != nil {
		fmt.Printf("[ERROR] Update: Binary verification failed: %v\n", err)
		fmt.Printf("[DEBUG] Update: Attempting rollback due to verification failure\n")
		// Rollback on verification failure
		if rollbackErr := u.RestoreFromBackup(backupPath, execPath); rollbackErr != nil {
			fmt.Printf("[ERROR] Update: Verification rollback also failed: %v\n", rollbackErr)
			return fmt.Errorf("new binary verification failed and rollback also failed: verification error: %w, rollback error: %w", err, rollbackErr)
		}
		fmt.Printf("[DEBUG] Update: Verification rollback completed successfully\n")
		// Clean up backup after successful rollback
		u.CleanupBackup(backupPath)
		return fmt.Errorf("new binary verification failed but was rolled back successfully: %w", err)
	}
	fmt.Printf("[DEBUG] Update: Binary verification passed\n")

	// Step 4: Clean up backup after successful update
	fmt.Printf("[DEBUG] Update: Cleaning up backup file\n")
	if err := u.CleanupBackup(backupPath); err != nil {
		// This is not critical - update succeeded, just log the cleanup failure
		fmt.Printf("[WARNING] Update: Backup cleanup failed (non-critical): %v\n", err)
		return fmt.Errorf("update succeeded but failed to clean up backup %s: %w", backupPath, err)
	}
	fmt.Printf("[DEBUG] Update: Backup cleanup completed\n")

	fmt.Printf("[DEBUG] Update: Atomic update process completed successfully\n")
	return nil
}

// verifyNewBinary performs basic verification that the new binary is functional
func (u *UpdateUtil) verifyNewBinary(execPath string) error {
	// Check that file exists and is executable
	info, err := os.Stat(execPath)
	if err != nil {
		return fmt.Errorf("new binary is not accessible: %w", err)
	}

	// Check that it's not empty
	if info.Size() == 0 {
		return fmt.Errorf("new binary is empty")
	}

	// On Unix, check that it has executable permissions
	if runtime.GOOS != "windows" {
		if info.Mode()&0111 == 0 {
			return fmt.Errorf("new binary is not executable")
		}
	}

	// TODO: Could add more sophisticated verification like:
	// - Running `./binary --version` to check it starts
	// - Verifying binary signature
	// - Checking that it's the correct architecture

	return nil
}

// CleanupTempFiles removes temporary files created during update process
func (u *UpdateUtil) CleanupTempFiles(execPath string) error {
	if execPath == "" {
		return fmt.Errorf("executable path cannot be empty")
	}

	// List of temporary file patterns to clean up
	patterns := []string{
		execPath + ".new",
		execPath + ".old",
		execPath + ".temp",
		execPath + ".pending",
		execPath + ".update.bat",
	}

	for _, pattern := range patterns {
		if _, err := os.Stat(pattern); err == nil {
			if err := os.Remove(pattern); err != nil {
				// Log but don't fail on cleanup errors
				fmt.Printf("Warning: failed to remove temporary file %s: %v\n", pattern, err)
			}
		}
	}

	// Also clean up old backup files (older than 7 days)
	return u.cleanupOldBackups(execPath, 7*24*time.Hour)
}

// cleanupOldBackups removes backup files older than the specified duration
func (u *UpdateUtil) cleanupOldBackups(execPath string, maxAge time.Duration) error {
	dir := filepath.Dir(execPath)
	baseName := filepath.Base(execPath)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read directory for backup cleanup: %w", err)
	}

	now := time.Now()
	cleaned := 0

	for _, entry := range entries {
		name := entry.Name()
		// Look for backup files matching our pattern
		if filepath.HasPrefix(name, baseName+".backup.") {
			fullPath := filepath.Join(dir, name)
			info, err := entry.Info()
			if err != nil {
				continue
			}

			// Remove if older than maxAge
			if now.Sub(info.ModTime()) > maxAge {
				if err := os.Remove(fullPath); err == nil {
					cleaned++
				}
			}
		}
	}

	if cleaned > 0 {
		fmt.Printf("Cleaned up %d old backup files\n", cleaned)
	}

	return nil
}