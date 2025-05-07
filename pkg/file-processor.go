package pkg

import (
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

type fileSections = map[string][]Section

func getFolderEntriesCount(chunkSize int) (*fileSections, error) {
	wd, _ := os.Getwd()
	path := filepath.Join(wd)
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	fs := fileSections{}

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

func parseFileMetaChunk(de []os.DirEntry, fs fileSections, path string) error {
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

func readContentIntoFileMeta(filepath string) ([]byte, error) {
	buf, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	return buf, err
}
