package pkg

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	fileSeparator = "_"
	folderName    = "sections"
	fileExtention = ".txt"
)

type FileSections = map[string][]Section

func GetFileSections(chunkSize int) (*FileSections, error) {
	goModPath, err := getRootFolderPath()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(filepath.Dir(goModPath), folderName)

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	fs := FileSections{}

	for {
		entries, err := f.ReadDir(chunkSize)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		if err := parseFileMetaChunk(entries, fs, path); err != nil {
			return nil, err
		}
	}
	return &fs, nil
}

func parseFileMetaChunk(de []os.DirEntry, fs FileSections, path string) error {
	for _, e := range de {
		s := strings.Split(strings.TrimSuffix(e.Name(), fileExtention), fileSeparator)
		content, err := readContentIntoFileMeta(filepath.Join(path, e.Name()))
		if err != nil {
			return err
		}
		fs[s[0]] = append(fs[s[0]], Section{SectionName: s[1], Content: string(content)})
	}
	return nil
}

func getRootFolderPath() (string, error) {
	curDir, _ := os.Getwd()
	for {
		goModPath := filepath.Join(curDir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return goModPath, nil
		}
		parentDir := filepath.Dir(goModPath)
		if parentDir == curDir {
			return "", errors.New("reached root dir without go.mod")
		}
		curDir = parentDir
	}
}

func readContentIntoFileMeta(filepath string) ([]byte, error) {
	buf, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	return buf, err
}
