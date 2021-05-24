package acn

import (
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseGoImports(t *testing.T) {
	type args struct {
		src string
	}
	type expected struct {
		goImports []GoImport
	}
	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "should return empty array if no imports are present",
			args: args{
				src: `
				package main
				const c = 1.0
				var X = f(3.14)*2 + c
				`,
			},
			expected: expected{
				goImports: []GoImport{},
			},
		},
		{
			name: "should parse one liner import",
			args: args{
				src: `
				package main
				import "fmt"
				const c = 1.0
				var X = f(3.14)*2 + c
				`,
			},
			expected: expected{
				goImports: []GoImport{
					{
						Path:        "fmt",
						PackageName: "fmt",
					},
				},
			},
		},
		{
			name: "should parse one liner import with alias",
			args: args{
				src: `
				package main
				import . "fmt"
				const c = 1.0
				var X = f(3.14)*2 + c
				`,
			},
			expected: expected{
				goImports: []GoImport{
					{
						Path:        "fmt",
						PackageName: "fmt",
						Alias:       ".",
					},
				},
			},
		},
		{
			name: "should parse multiple imports",
			args: args{
				src: `
				package main
				import (
					"fmt"
					_ "net/http"
					. "os"
					"github.com/javiercbk/acn"
					imperfectpack "github.com/javiercbk/impack"
					"gopkg.in/alexcesaro/statsd.v2"

				)
				const c = 1.0
				var X = f(3.14)*2 + c
				`,
			},
			expected: expected{
				goImports: []GoImport{
					{
						Path:        "fmt",
						PackageName: "fmt",
						Alias:       "",
					},
					{
						Path:        "net/http",
						PackageName: "http",
						Alias:       "_",
					},
					{
						Path:        "os",
						PackageName: "os",
						Alias:       ".",
					},
					{
						Path:        "github.com/javiercbk/acn",
						PackageName: "acn",
						Alias:       "",
					},
					{
						Path:        "github.com/javiercbk/impack",
						PackageName: "impack",
						Alias:       "imperfectpack",
					},
					{
						Path:        "gopkg.in/alexcesaro/statsd.v2",
						PackageName: "statsd",
						Alias:       "",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet() // positions are relative to fset
			f, err := parser.ParseFile(fset, "src.go", tt.args.src, 0)
			assert.Nil(t, err)
			goImports := ParseGoImports(f)
			assert.Equal(t, tt.expected.goImports, goImports)
		})
	}
}
