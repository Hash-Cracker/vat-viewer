package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

// nordTheme defines the Nord color theme
var nordTheme = styles.Register(chroma.MustNewStyle("nord", chroma.StyleEntries{
	// Base colors
	chroma.Text:                "#D8DEE9",
	chroma.Error:               "#BF616A",
	chroma.Background:          "#2E3440",
	
	// Comments
	chroma.Comment:             "#4C566A",
	chroma.CommentHashbang:     "#4C566A",
	chroma.CommentMultiline:    "#4C566A",
	chroma.CommentPreproc:      "#5E81AC",
	chroma.CommentSingle:       "#4C566A",
	chroma.CommentSpecial:      "#5E81AC",
	
	// Keywords
	chroma.Keyword:             "#81A1C1",
	chroma.KeywordConstant:     "#81A1C1",
	chroma.KeywordDeclaration:  "#81A1C1",
	chroma.KeywordNamespace:    "#81A1C1",
	chroma.KeywordPseudo:       "#81A1C1",
	chroma.KeywordReserved:     "#81A1C1",
	chroma.KeywordType:         "#81A1C1",
	
	// Operators and Punctuation
	chroma.Operator:            "#81A1C1",
	chroma.OperatorWord:        "#81A1C1",
	chroma.Punctuation:         "#D8DEE9",
	
	// Names and Identifiers
	chroma.Name:                "#D8DEE9",
	chroma.NameAttribute:       "#8FBCBB",
	chroma.NameBuiltin:         "#88C0D0",
	chroma.NameBuiltinPseudo:   "#88C0D0",
	chroma.NameClass:           "#8FBCBB",
	chroma.NameConstant:        "#8FBCBB",
	chroma.NameDecorator:       "#88C0D0",
	chroma.NameEntity:          "#8FBCBB",
	chroma.NameException:       "#BF616A",
	chroma.NameFunction:        "#88C0D0",
	chroma.NameFunctionMagic:   "#88C0D0",
	chroma.NameLabel:           "#8FBCBB",
	chroma.NameNamespace:       "#8FBCBB",
	chroma.NameOther:          "#D8DEE9",
	chroma.NameTag:            "#81A1C1",
	chroma.NameVariable:       "#D8DEE9",
	chroma.NameVariableClass:  "#D8DEE9",
	chroma.NameVariableGlobal: "#D8DEE9",
	chroma.NameVariableInstance: "#D8DEE9",
	
	// Literals
	chroma.Literal:             "#81A1C1",
	chroma.LiteralDate:         "#B48EAD",
	chroma.LiteralString:       "#A3BE8C",
	chroma.LiteralStringBacktick: "#A3BE8C",
	chroma.LiteralStringChar:   "#A3BE8C",
	chroma.LiteralStringDoc:    "#A3BE8C",
	chroma.LiteralStringDouble: "#A3BE8C",
	chroma.LiteralStringEscape: "#EBCB8B",
	chroma.LiteralStringHeredoc: "#A3BE8C",
	chroma.LiteralStringInterpol: "#A3BE8C",
	chroma.LiteralStringOther:  "#A3BE8C",
	chroma.LiteralStringRegex:  "#EBCB8B",
	chroma.LiteralStringSingle: "#A3BE8C",
	chroma.LiteralStringSymbol: "#A3BE8C",
	chroma.LiteralNumber:       "#B48EAD",
	chroma.LiteralNumberBin:    "#B48EAD",
	chroma.LiteralNumberFloat:  "#B48EAD",
	chroma.LiteralNumberHex:    "#B48EAD",
	chroma.LiteralNumberInteger: "#B48EAD",
	chroma.LiteralNumberOct:    "#B48EAD",
	
	// Generic
	chroma.Generic:             "#D8DEE9",
	chroma.GenericDeleted:      "#BF616A",
	chroma.GenericEmph:         "#D8DEE9 italic",
	chroma.GenericError:        "#BF616A",
	chroma.GenericHeading:      "#88C0D0 bold",
	chroma.GenericInserted:     "#A3BE8C",
	chroma.GenericOutput:       "#D8DEE9",
	chroma.GenericPrompt:       "#4C566A",
	chroma.GenericStrong:       "#D8DEE9 bold",
	chroma.GenericSubheading:   "#88C0D0",
	chroma.GenericTraceback:    "#BF616A",
	chroma.GenericUnderline:    "underline",
}))

func highlightCode(content string, filename string) string {
	// Get the lexer based on filename
	lexer := lexers.Match(filename)
	if lexer == nil {
		lexer = lexers.Fallback
	}

	// Create a custom analyzer for improved token handling
	analyzer := &tokenAnalyzer{
		importedPackages: make(map[string]bool),
		packageFuncs:    make(map[string]bool),
	}
	analyzer.scanImports(content)

	formatter := formatters.Get("terminal256")
	if formatter == nil {
		formatter = formatters.Fallback
	}

	style := styles.Get("nord")
	if style == nil {
		style = styles.Fallback
	}

	iterator, err := lexer.Tokenise(nil, content)
	if err != nil {
		return content
	}

	// Enhance the token stream with package function recognition
	enhancedIterator := analyzer.enhanceTokens(iterator)

	var buf strings.Builder
	err = formatter.Format(&buf, style, enhancedIterator)
	if err != nil {
		return content
	}

	return buf.String()
}

// tokenAnalyzer helps analyze and enhance token recognition
type tokenAnalyzer struct {
	importedPackages map[string]bool
	packageFuncs     map[string]bool
	currentPackage   string
}

// scanImports scans the content for import statements and builds a map of imported packages
func (a *tokenAnalyzer) scanImports(content string) {
	scanner := bufio.NewScanner(strings.NewReader(content))
	inImportBlock := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "package ") {
			a.currentPackage = strings.TrimPrefix(line, "package ")
			continue
		}

		if line == "import (" {
			inImportBlock = true
			continue
		}

		if line == ")" {
			inImportBlock = false
			continue
		}

		if inImportBlock || strings.HasPrefix(line, "import ") {
			pkgPath := strings.Trim(strings.TrimPrefix(line, "import "), "\"` ")
			if strings.Contains(pkgPath, " ") {
				parts := strings.Fields(pkgPath)
				pkgPath = strings.Trim(parts[len(parts)-1], "\"` ")
			}
			a.importedPackages[pkgPath] = true
		}
	}
}

// enhanceTokens processes the token stream to improve function recognition
func (a *tokenAnalyzer) enhanceTokens(iterator chroma.Iterator) chroma.Iterator {
	var tokens []chroma.Token
	var lastDot bool
	var lastIdent string

	for _, token := range iterator.Tokens() {
		// Track package-qualified function calls
		if token.Type == chroma.NameOther {
			if lastDot {
				// This is likely a function call from an imported package
				if a.importedPackages[lastIdent] {
					token.Type = chroma.NameFunctionMagic
				}
			}
			lastIdent = token.Value
			lastDot = false
		} else if token.Value == "." {
			lastDot = true
		} else {
			lastDot = false
		}

		tokens = append(tokens, token)
	}

	return chroma.Literator(tokens...)
}

func printFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		return
	}

	highlighted := highlightCode(string(content), filename)
	fmt.Print(highlighted)
}

func main() {
	if len(os.Args) < 2 {
		scanner := bufio.NewScanner(os.Stdin)
		var content strings.Builder
		for scanner.Scan() {
			content.WriteString(scanner.Text() + "\n")
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
			return
		}
		
		highlighted := highlightCode(content.String(), "input.go")
		fmt.Print(highlighted)
	} else {
		for _, filename := range os.Args[1:] {
			printFile(filename)
		}
	}
}
