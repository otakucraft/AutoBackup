package utils

import "os"

func DirExists(path string) bool {
	dir, err := os.Stat(path)
	return err == nil && dir.IsDir()
}

func FileExists(path string) bool {
	dir, err := os.Stat(path)
	return err == nil && !dir.IsDir()
}

func TouchDir(path string) bool {
	err := os.Mkdir(path, 0755)
	return err == nil
}

func TouchFile(path string) bool {
	file, err := os.Create(path)
	if err != nil {
		return false
	}
	err = file.Close()
	return err == nil
}
