package log

import (
	"compress/gzip"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"
	"time"

	"github.com/StoneTrench/CovarexEngine/engine/platform"
)

// Helper function for log file management.
func compressAndRemove(fs *platform.VFS, srcPath, dstPath string) error {
	srcFile, err := fs.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open old file: %w", err)
	}
	defer srcFile.Close()

	dstFile, err := fs.Create(dstPath)
	if err != nil {
		return fmt.Errorf("failed to create new file: %w", err)
	}
	defer dstFile.Close()

	gw := gzip.NewWriter(dstFile)
	if _, err := io.Copy(gw, srcFile); err != nil {
		return fmt.Errorf("failed to compress old file: %w", err)
	}
	defer gw.Close()
	srcFile.Close()

	if err := fs.Remove(srcPath); err != nil {
		return fmt.Errorf("failed to remove old file: %w", err)
	}

	return nil
}

// Init initializes the main logger.
// Should be called once at startup.
func Init(fs *platform.VFS, logDir string, noFile bool) (err error) {
	if !fs.HasWriteAccess() {
		noFile = true
	}

	// By the way, these comments are human made and are here cause when I was half asleep debugging the code, I kept thinking parts of this were mistakes lol
	var latestPath string
	if !noFile {
		latestPath = path.Join(logDir, "latest.log")
	}

	var file io.Writer
	var old_log_error error
	if !noFile {
		// Ignore the mkdir error intentionally
		fs.MkdirAll(logDir, platform.MkdirAll_DEF_PERM)

		// Check if latest exists, if yes move it
		if _, err := fs.Stat(latestPath); err == nil {
			archiveName := path.Join(logDir, fmt.Sprintf("log-%s.log.gz", time.Now().Format("2006-01-02-150405")))

			if err := compressAndRemove(fs, latestPath, archiveName); err != nil {
				old_log_error = fmt.Errorf("failed to archive old log: %w", err)
			}
		}

		// No defer file.Close() because this stays open for the entire time the process is running
		file, err = fs.Create(latestPath)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
	}

	if file == nil {
		file = os.Stdout
	} else {
		file = io.MultiWriter(os.Stdout, file)
	}

	handler := NewBracketHandler(file)
	slog.SetDefault(slog.New(handler))

	if old_log_error != nil {
		slog.Warn(old_log_error.Error())
	}

	return nil
}
