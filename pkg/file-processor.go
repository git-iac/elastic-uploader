package pkg

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

const fileSeparator = "_" // files seperators are always '_'

type fileSections = map[string][]Section

func getFolderEntriesCount(chunkSize int) (*fileSections, error) {
	wd, _ := os.Getwd()
	path := filepath.Join(wd, "sections")
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	fMetaMap := fileSections{}

	for {
		entries, err := f.ReadDir(chunkSize)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		if err := parseFileMetaChunk(entries, fMetaMap, path); err != nil {
			return nil, err
		}
	}
	return &fMetaMap, nil
}

func parseFileMetaChunk(de []os.DirEntry, fs fileSections, path string) error {
	for _, e := range de {
		s := strings.Split(strings.TrimSuffix(e.Name(), ".txt"), fileSeparator)

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
