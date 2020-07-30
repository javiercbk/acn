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
			folder: "testdata/mod-project/empty",
		}, {
			name:          "project directory",
			folder:        "testdata/mod-project",
			expectedFiles: []string{"markers.go", "security.go", "security_test.go"},
		}, {
			name:          "project directory with backlist",
			folder:        "testdata/mod-project",
			expectedFiles: []string{"markers.go", "security.go"},
			blacklist:     []*regexp.Regexp{regexp.MustCompile(`.*_test\.go`)},
		}, {
			name:          "project directory with backlist matching nothing",
			folder:        "testdata/mod-project",
			expectedFiles: []string{"markers.go", "security.go", "security_test.go"},
			blacklist:     []*regexp.Regexp{regexp.MustCompile(`.*\.yaml`)},
		}, {
			name:          "security directory",
			folder:        "testdata/mod-project/security",
			expectedFiles: []string{"security.go", "security_test.go"},
		}, {
			name:          "security directory with blacklist",
			folder:        "testdata/mod-project/security",
			expectedFiles: []string{"security.go"},
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
			folder: "testdata/mod-project/empty",
		}, {
			name:          "project directory",
			folder:        "testdata/mod-project",
			expectedFiles: []string{"markers.go"},
		}, {
			name:      "project directory with backlist",
			folder:    "testdata/mod-project",
			blacklist: []*regexp.Regexp{regexp.MustCompile(`.*_test\.go`)},
		}, {
			name:          "security directory",
			expectedFiles: []string{"security.go", "security_test.go"},
			folder:        "testdata/mod-project/security",
		}, {
			name:          "security directory with blacklist",
			folder:        "testdata/mod-project/security",
			expectedFiles: []string{"security.go"},
			blacklist:     []*regexp.Regexp{regexp.MustCompile(`.*_test\.go`)},
		}, {
			name:          "security directory with blacklist matching nothing",
			folder:        "testdata/mod-project/security",
			expectedFiles: []string{"security.go", "security_test.go"},
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
