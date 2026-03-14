package analyzer

import (
	"go/ast"
	"regexp"
	"go/types"
	"strings"
	"go/token"
	"strconv"
	"unicode"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// regexp for checks
var (
	// allow only eng letters, digits, spaсes, basic punctuation
	allowedCharsRx = regexp.MustCompile(`^[a-zA-Z0-9\s\-_:\.,='"]+$`)
	specialEndRx = regexp.MustCompile(`[!?.]$`)

	// key words
	sensitiveWords =[]string{"password", "token", "api_key", "secret"}
)

var Analyzer = &analysis.Analyzer{
	Name:     "loglint",
	Doc:      "Checks logs for formatting and sensitive data",
	Run:      run,
	Requires:[]*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// taking only call expressions
	nodeFilter :=[]ast.Node{
		(*ast.CallExpr)(nil),
	}

	// iterating through nodes
	inspect.Preorder(nodeFilter, func(n ast.Node) {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return
		}

		// check if function call
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return
		}

		// get info about function
		obj := pass.TypesInfo.Uses[sel.Sel]
		if obj == nil {
			return
		}
		
		funcDecl, ok := obj.(*types.Func)
		if !ok {
			return
		}

		// check if package == slog / zap
		pkg := funcDecl.Pkg()
		if pkg == nil {
			return
		}
		pkgPath := pkg.Path()
		isSlog := pkgPath == "log/slog"
		isZap := strings.HasPrefix(pkgPath, "go.uber.org/zap")

		// return if smth else
		if !isSlog && !isZap {
			return
		}

		// check if func name is {"Debug", "Info", "Warn" etc}
		methodName := funcDecl.Name()
		if !isLogMethod(methodName) {
			return
		}

		// get log body
		if len(call.Args) == 0 {
			return
		}

		// trying to find string literal
		basicLit, ok := call.Args[0].(*ast.BasicLit)
		if !ok || basicLit.Kind != token.STRING {
			return
		}

		// remove quotes from log body
		msg, err := strconv.Unquote(basicLit.Value)
		if err != nil || len(msg) == 0 {
			return
		}

		// --- RULES ---

		// 1. must start with lowercase letter
		firstRune := []rune(msg)[0]
		if unicode.IsUpper(firstRune) {
			pass.Reportf(basicLit.Pos(), "log message must start with lowercase letter")
		}

		// 2-3. must include only eng letters, nums and basic puntuation
		if !allowedCharsRx.MatchString(msg) {
			pass.Reportf(basicLit.Pos(), "log message must contain only english letters, numbers and basic punctuation (no emojis or special chars)")
		}
		if specialEndRx.MatchString(msg) || strings.HasSuffix(msg, "...") {
			pass.Reportf(basicLit.Pos(), "log message should not end with punctuation marks like '!', '?' or '...'")
		}

		// 4. log message must not contain sensitive data
		lowerMsg := strings.ToLower(msg)
		for _, word := range sensitiveWords {
			if strings.Contains(lowerMsg, word) {
				pass.Reportf(basicLit.Pos(), "log message contains potentially sensitive data: %s", word)
				break
			}
		}
	})

	return nil, nil
}

// checks if func is used for logging
func isLogMethod(name string) bool {
	switch name {
	case "Debug", "Info", "Warn", "Error", "Fatal", "Panic":
		return true
	}
	return false
}
