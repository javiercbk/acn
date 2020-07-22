package acn

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// isGoFile returns true if a file name is a go file
func isGoFile(fileName string) bool {
	return strings.HasSuffix(fileName, ".go")
}

// listGoFilesRecursively returns a list of go files in a directory recursively
func listGoFilesRecursively(dir string, blacklist []*regexp.Regexp) ([]string, error) {
	files := make([]string, 0)
	err := filepath.Walk(dir, func(filePath string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && isGoFile(info.Name()) && !shouldIgnore(filePath, blacklist) {
			files = append(files, filePath)
		}
		return err
	})
	return files, err
}

// listGoFiles list all files in directory
func listGoFiles(dir string, blacklist []*regexp.Regexp) ([]string, error) {
	files := make([]string, 0)
	filesInDir, err := ioutil.ReadDir(dir)
	if err != nil {
		return files, err
	}
	for _, info := range filesInDir {
		if !info.IsDir() && isGoFile(info.Name()) && !shouldIgnore(info.Name(), blacklist) {
			files = append(files, info.Name())
		}
	}
	return files, err
}

func shouldIgnore(filePath string, blacklist []*regexp.Regexp) bool {
	for _, r := range blacklist {
		if r.MatchString(filePath) {
			return true
		}
	}
	return false
}
