package vspb

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Installer struct {
	FromDir  string
	Filename string
	ToDir    string
}

func (i Installer) Install() error {
	fileList := make([]string, 0)
	if strings.Contains(i.Filename, ",") {
		fileList = strings.Split(i.Filename, ",")
	} else {
		fileList = append(fileList, i.Filename)
	}

	needInstall := len(fileList)
	installed := 0

	err := filepath.Walk(i.FromDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("find file failed in path %q: %w", path, err)
		}

		if info.IsDir() {
			return nil
		}

		for _, filename := range fileList {
			if info.Name() == filename {
				_, err := copyFile(path, i.ToDir+"/"+filename)
				if err != nil {
					return err
				}

				installed++
				if installed == needInstall {
					return filepath.SkipAll
				}
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	// TODO:
	// which file not exists
	if installed != needInstall {
		return fmt.Errorf("file %s not exists", i.Filename)
	}

	return nil
}

func copyFile(source, target string) (int64, error) {
	sourceFile, err := os.Open(source)
	if err != nil {
		return 0, err
	}
	defer sourceFile.Close()

	// if target exists, remove it.
	if _, err := os.Stat(target); err == nil {
		if err := os.Remove(target); err != nil {
			return 0, err
		}
	}

	targetFile, err := os.OpenFile(target, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0755)
	if err != nil {
		return 0, err
	}
	defer targetFile.Close()

	return io.Copy(targetFile, sourceFile)
}
