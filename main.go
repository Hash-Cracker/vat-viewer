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
	chroma.Text:                "#D8DEE9",
	chroma.Error:               "#BF616A",
	chroma.Comment:            "#4C566A",
	chroma.Keyword:            "#81A1C1",
	chroma.KeywordConstant:    "#81A1C1",
	chroma.KeywordDeclaration: "#81A1C1",
	chroma.KeywordNamespace:   "#81A1C1",
	chroma.KeywordType:        "#81A1C1",
	chroma.Operator:           "#81A1C1",
	chroma.Punctuation:        "#D8DEE9",
	chroma.Name:               "#D8DEE9",
	chroma.NameBuiltin:        "#88C0D0",
	chroma.NameFunction:       "#88C0D0",
	chroma.NameClass:          "#8FBCBB",
	chroma.NameVariable:       "#D8DEE9",
	chroma.LiteralString:      "#A3BE8C",
	chroma.LiteralNumber:      "#B48EAD",
	chroma.GenericHeading:     "#88C0D0",
	chroma.GenericSubheading:  "#88C0D0",
	chroma.Background:         "#2E3440",
}))

func highlightCode(content string, filename string) string {
	// Determine lexer based on filename extension
	lexer := lexers.Match(filename)
	if lexer == nil {
		lexer = lexers.Fallback
	}

	// Create formatter for terminal output
	formatter := formatters.Get("terminal256")
	if formatter == nil {
		formatter = formatters.Fallback
	}

	// Use Nord theme
	style := styles.Get("nord")
	if style == nil {
		style = styles.Fallback
	}

	iterator, err := lexer.Tokenise(nil, content)
	if err != nil {
		return content
	}

	var buf strings.Builder
	err = formatter.Format(&buf, style, iterator)
	if err != nil {
		return content
	}

	return buf.String()
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
		// No files provided; read from standard input
		scanner := bufio.NewScanner(os.Stdin)
		var content strings.Builder
		for scanner.Scan() {
			content.WriteString(scanner.Text() + "\n")
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
			return
		}
		
		// Highlight content from stdin as Go code by default
		highlighted := highlightCode(content.String(), "input.go")
		fmt.Print(highlighted)
	} else {
		// Loop through all provided files
		for _, filename := range os.Args[1:] {
			printFile(filename)
		}
	}
}
