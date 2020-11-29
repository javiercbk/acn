package acn

import (
	"regexp"
	"testing"
)

func TestFileNameMatcher(t *testing.T) {
	tests := []struct {
		name        string
		fileToMatch *regexp.Regexp
		filePath    string
		expected    bool
	}{
		{
			name:        "empty match string #1",
			fileToMatch: nil,
			filePath:    "path/to/file.go",
			expected:    true,
		}, {
			name:        "empty match string #2",
			fileToMatch: nil,
			filePath:    "file.go",
			expected:    true,
		}, {
			name:        "match file.go #1",
			fileToMatch: regexp.MustCompile(`.*file.go`),
			filePath:    "path/to/file.go",
			expected:    true,
		}, {
			name:        "match file.go #2",
			fileToMatch: regexp.MustCompile(`.*file.go`),
			filePath:    "file.go",
			expected:    true,
		}, {
			name:        "match file.go match #3",
			fileToMatch: regexp.MustCompile(`file.go`),
			filePath:    "path/to/file.go",
			expected:    true,
		}, {
			name:        "should not match #1",
			fileToMatch: regexp.MustCompile(`.*/file.go`),
			filePath:    "file.go",
			expected:    false,
		}, {
			name:        "should not match #2",
			fileToMatch: regexp.MustCompile(`.*other/file.go`),
			filePath:    "file.go",
			expected:    false,
		}, {
			name:        "should not match #2",
			fileToMatch: regexp.MustCompile(`.*other/file.go`),
			filePath:    "path/to/file.go",
			expected:    false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			matcher := fileNameMatcherFactory(test.fileToMatch)
			result := matcher(test.filePath)
			if test.expected != result {
				t.Fatalf("expected match to be %v but was %v", test.expected, result)
			}
		})
	}
}

func TestFindMain(t *testing.T) {
	// logger := log.New(ioutil.Discard, "", log.Llongfile)
	// project := NewGoProject(logger, Config{
	// 	Folder: "testdata/mod-project",
	// })
	// err := n.project()
	// if err != nil {
	// 	t.Fatalf("error finding main function %v", err)
	// }
	// if n.mainFunc == nil {
	// 	t.Fatal("main function was not found")
	// }
	// if n.mainFunc.decl.Name.Name != "main" {
	// 	t.Fatalf("expected function name to be main but was %s", n.mainFunc.decl.Name.Name)
	// }
	// if len(n.mainFunc.decl.Type.Params.List) != 0 {
	// 	t.Fatalf("expected function to not have parameters but had %d", len(n.mainFunc.decl.Type.Params.List))
	// }
	// if n.mainFunc.decl.Recv != nil {
	// 	t.Fatalf("expected function to not have receiver but it had %v, meaning it is a method", n.mainFunc.decl.Recv)
	// }
}

func TestFindFunction(t *testing.T) {
	// f := ast.FuncDecl{}
	// goProject := NewGoProject("./testdata/mod-project", []string{"modproj/http"})
	// goProject.Analize()
	// functionFound := goProject("./testdata/mod-project", []string{"modproj/http"}, "customHTTPErrorHandler", &f)
	// if !functionFound {
	// 	t.Fatal("expected function to be found")
	// }
}
