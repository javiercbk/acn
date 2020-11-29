package acn

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"

	"golang.org/x/tools/go/packages"
)

var (
	defaultBlacklist = []*regexp.Regexp{
		// ignore tests
		regexp.MustCompile(`.*_test\.go`),
		// ignore anything inside testdata
		regexp.MustCompile(".*" + string(os.PathSeparator) + "testdata" + string(os.PathSeparator) + ".*"),
		// ignore anything inside vendor
		regexp.MustCompile(".*" + string(os.PathSeparator) + "vendor" + string(os.PathSeparator) + ".*"),
		// ignore anything starting with dot
		regexp.MustCompile(".*" + string(os.PathSeparator) + "\\..*"),
		// ignore anything starting with underscore
		regexp.MustCompile(".*" + string(os.PathSeparator) + "\\_.*"),
	}
)

// ASTFuncDeclContext is a function declaration in the ast
type ASTFuncDeclContext struct {
	fset *token.FileSet
	file *ast.File
	decl *ast.FuncDecl
}

// Config is the Navigator config
type Config struct {
	// Folder is the folder of the Go project
	Folder       string
	MainFileName *regexp.Regexp
	Blacklist    []*regexp.Regexp
}

// Navigator is a stateful code navigator that will feed matching AST lines to the caller
type Navigator struct {
	logger   *log.Logger
	conf     Config
	pkgs     []*packages.Package
	mainFunc *ASTFuncDeclContext
}

// NewNavigator creates a new Navigator
func NewNavigator(logger *log.Logger, conf Config) Navigator {
	return Navigator{
		logger: logger,
		conf:   conf,
	}
}

// FindStructCompositeLiteral finds all struct values
func (n *Navigator) FindStructCompositeLiteral() ([]*ast.CompositeLit, error) {
	// use go-guru to search
	return nil, nil
}

// FindFunctionCallFromFunc funds a function call using a function declaration as node starting point.
// If funcDecl is nil, this function will assume that the main function is the starting point
func (n *Navigator) FindFunctionCallFromFunc(funcDecl *ast.FuncDecl) (*ast.CallExpr, error) {
	decl := funcDecl
	// if func declaration is nil then assume main
	if decl == nil {
		err := n.findMainFunc()
		if err != nil {
			return nil, err
		}
		decl = n.mainFunc.decl
	}
	// use go-guru to search for function calls in decl
	// start exploring depth first in every function call
	return nil, nil
}

// GoProject
type GoProject struct {
	pkgs []*packages.Package
}

// AnalizeGoProject analizes a go project and extract all selected packages metadata
func AnalizeGoProject(rootDir string, explorePackages []string) (GoProject, error) {
	proj := GoProject{}
	var fset = token.NewFileSet()
	cfg := &packages.Config{Fset: fset, Mode: packages.LoadAllSyntax, Dir: rootDir}
	pkgs, err := packages.Load(cfg, explorePackages...)
	if err == nil {
		proj.pkgs = pkgs
	}
	return proj, err
}

func (g *GoProject) FindFunction(funcName string, funcDeclFound *ast.FuncDecl) bool {
	found := false
	for _, pkg := range g.pkgs {
		for _, fileAst := range pkg.Syntax {
			ast.Inspect(fileAst, func(n ast.Node) bool {
				if funcDecl, ok := n.(*ast.FuncDecl); ok {
					if funcDecl.Name.Name == funcName {
						funcDeclFound = funcDecl
						found = true
					}
				}
				return !found
			})
			if found {
				return true
			}
		}
	}
	return false
}
func (n *Navigator) findMainFunc() error {
	if n.mainFunc != nil {
		return nil
	}
	blacklist := n.conf.Blacklist
	if blacklist == nil {
		blacklist = defaultBlacklist
	}
	goFiles, err := listGoFilesRecursively(n.conf.Folder, blacklist)
	if err != nil {
		n.logger.Printf("error listing files recursively: %v", err)
		return err
	}
	fileNameMatcher := fileNameMatcherFactory(n.conf.MainFileName)
	for _, goFile := range goFiles {
		if !fileNameMatcher(goFile) {
			continue
		}
		fset := token.NewFileSet()
		file, err := astForFile(goFile, fset)
		if err != nil {
			n.logger.Printf("error reading ast for file %s: %v", goFile, err)
			return err
		}
		funcDecl, err := findFunction(file, "main", paramCountOpt(0), isFunctionOpt())
		if err == nil {
			// cache main function
			n.mainFunc = &ASTFuncDeclContext{
				fset: fset,
				file: file,
				decl: funcDecl,
			}
			return nil
		}
	}
	return ErrNotFound
}

func fileNameMatcherFactory(fileRegexp *regexp.Regexp) func(goFile string) bool {
	if fileRegexp == nil {
		return func(goFile string) bool {
			return true
		}
	}
	return func(goFile string) bool {
		return fileRegexp.MatchString(goFile)
	}
}

// findFuncOption is a function that analyzes a func declaration and returns true if it matches or false if it does not.
type findFuncOption func(*ast.FuncDecl) bool

func paramCountOpt(paramsCount int) findFuncOption {
	return func(decl *ast.FuncDecl) bool {
		return paramsCount == len(decl.Type.Params.List)
	}
}

func isFunctionOpt() findFuncOption {
	return func(decl *ast.FuncDecl) bool {
		return decl.Recv == nil
	}
}

func findFunction(file *ast.File, name string, options ...findFuncOption) (*ast.FuncDecl, error) {
	for _, d := range file.Decls {
		switch x := d.(type) {
		case *ast.FuncDecl:
			if x.Name.Name == name {
				matches := true
				for i := range options {
					if !options[i](x) {
						matches = false
						break
					}
				}
				if matches {
					return x, nil
				}
			}
		}
	}
	return nil, ErrNotFound
}

func astForFile(filePath string, fset *token.FileSet) (*ast.File, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return astForReader(filePath, f, fset)
}

func astForReader(filePath string, r io.Reader, fset *token.FileSet) (*ast.File, error) {
	src, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return parser.ParseFile(fset, filePath, src, parser.ParseComments)
}
