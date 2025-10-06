package audit

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

func SaveToFile(path string, data []byte) error {
	logFileName := time.Now().Format("060102")

	name := filepath.Join(path, logFileName) + ".txt"

	if err := os.MkdirAll(path, 0o755); err != nil {
		if errors.Is(err, syscall.EROFS) || errors.Is(err, syscall.EPERM) {
			return fmt.Errorf("cannot create dir %s: read-only filesystem or no permissions: %w", path, err)
		}
		return fmt.Errorf("mkdir %s: %w", path, err)
	}

	file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o644) // creates file if not exists
	if err != nil {
		return fmt.Errorf("open file %s: %w", name, err)
	}
	defer file.Close()
	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("write file %s: %w", name, err)
	}
	return nil
}
