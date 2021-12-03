package utils

import "os"

// DirExists check if path is a dir and exists.
func DirExists(path string) bool {
	dir, err := os.Stat(path)
	return err == nil && dir.IsDir()
}

// FileExists check if path is a file and exists.
func FileExists(path string) bool {
	dir, err := os.Stat(path)
	return err == nil && !dir.IsDir()
}

// TouchDir create dir.
func TouchDir(path string) bool {
	err := os.Mkdir(path, 0755)
	return err == nil
}

// TouchFile create file.
func TouchFile(path string) bool {
	file, err := os.Create(path)
	if err != nil {
		return false
	}
	err = file.Close()
	return err == nil
}
