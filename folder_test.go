package acn

import (
	"os"
	"regexp"
	"strings"
	"syscall"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestListGoFilesRecursively(t *testing.T) {
	tests := []struct {
		name          string
		folder        string
		blacklist     []*regexp.Regexp
		expectedFiles []string
		expectedErr   error
	}{
		{
			name:   "non existing directory",
			folder: "./this-folder-does-not-exists",
			expectedErr: &os.PathError{
				Op:   "lstat",
				Path: "./this-folder-does-not-exists",
				Err:  syscall.ENOENT,
			},
		}, {
			name:   "empty directory",
			folder: "test/mod/db",
		}, {
			name:          "project directory",
			folder:        "test/mod",
			expectedFiles: []string{"mod.go", "mod_test.go"},
		}, {
			name:          "project directory with backlist",
			folder:        "test/mod",
			expectedFiles: []string{"mod.go"},
			blacklist:     []*regexp.Regexp{regexp.MustCompile(`.*_test\.go`)},
		}, {
			name:          "project directory with backlist matching nothing",
			folder:        "test/mod",
			expectedFiles: []string{"mod.go", "mod_test.go"},
			blacklist:     []*regexp.Regexp{regexp.MustCompile(`.*\.yaml`)},
		}, {
			name:          "cmd directory",
			expectedFiles: []string{"mod.go", "mod_test.go"},
			folder:        "test/mod/cmd",
		}, {
			name:          "cmd directory with blacklist",
			folder:        "test/mod/cmd",
			expectedFiles: []string{"mod.go"},
			blacklist:     []*regexp.Regexp{regexp.MustCompile(`.*_test\.go`)},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			files, err := listGoFilesRecursively(test.folder, test.blacklist)
			if !cmp.Equal(test.expectedErr, err) {
				t.Fatalf("expected error to match\n%s", cmp.Diff(test.expectedErr, err))
			}
			for _, f := range test.expectedFiles {
				found := false
				for _, foundFile := range files {
					if strings.Contains(foundFile, f) {
						found = true
						break
					}
				}
				if !found {
					t.Fatalf("expected file %s to be found but was not", f)
				}
			}
		})
	}
}

func TestListGoFiles(t *testing.T) {
	tests := []struct {
		name          string
		folder        string
		blacklist     []*regexp.Regexp
		expectedFiles []string
		expectedErr   error
	}{
		{
			name:   "non existing directory",
			folder: "./this-folder-does-not-exists",
			expectedErr: &os.PathError{
				Op:   "open",
				Path: "./this-folder-does-not-exists",
				Err:  syscall.ENOENT,
			},
		}, {
			name:   "empty directory",
			folder: "test/mod/db",
		}, {
			name:   "project directory",
			folder: "test/mod",
		}, {
			name:      "project directory with backlist",
			folder:    "test/mod",
			blacklist: []*regexp.Regexp{regexp.MustCompile(`.*_test\.go`)},
		}, {
			name:          "cmd directory",
			expectedFiles: []string{"mod.go", "mod_test.go"},
			folder:        "test/mod/cmd",
		}, {
			name:          "cmd directory with blacklist",
			folder:        "test/mod/cmd",
			expectedFiles: []string{"mod.go"},
			blacklist:     []*regexp.Regexp{regexp.MustCompile(`.*_test\.go`)},
		}, {
			name:          "cmd directory with blacklist matching nothing",
			folder:        "test/mod/cmd",
			expectedFiles: []string{"mod.go", "mod_test.go"},
			blacklist:     []*regexp.Regexp{regexp.MustCompile(`.*\.yaml`)},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			files, err := listGoFiles(test.folder, test.blacklist)
			if !cmp.Equal(test.expectedErr, err) {
				t.Fatalf("expected error to match\n%s", cmp.Diff(test.expectedErr, err))
			}
			for _, f := range test.expectedFiles {
				found := false
				for _, foundFile := range files {
					if strings.Contains(foundFile, f) {
						found = true
						break
					}
				}
				if !found {
					t.Fatalf("expected file %s to be found but was not", f)
				}
			}
		})
	}
}
