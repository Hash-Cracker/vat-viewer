package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Nord theme colors
const (
	Nord0  = "\033[38;2;46;52;64m"    // Polar Night (dark)
	Nord1  = "\033[38;2;59;66;82m"    // Polar Night
	Nord4  = "\033[38;2;216;222;233m" // Snow Storm (light)
	Nord7  = "\033[38;2;143;188;187m" // Frost
	Nord8  = "\033[38;2;136;192;208m" // Frost
	Nord9  = "\033[38;2;129;161;193m" // Frost
	Nord10 = "\033[38;2;94;129;172m"  // Frost
	Nord11 = "\033[38;2;191;97;106m"  // Aurora (red)
	Nord13 = "\033[38;2;235;203;139m" // Aurora (yellow)
	Nord14 = "\033[38;2;163;190;140m" // Aurora (green)
	Reset  = "\033[0m"
)

type Token struct {
	Type  string
	Value string
}

type Language struct {
	Keywords    []string
	Types      []string
	Constants   []string
	Operators   []string
	Extensions  []string
}

var languages = map[string]Language{
	"go": {
		Keywords:   []string{"package", "import", "func", "return", "if", "else", "for", "range", "var", "const", "defer"},
		Types:     []string{"string", "int", "byte", "error", "bool"},
		Constants: []string{"nil", "true", "false"},
		Operators: []string{":=", "!=", "==", ">=", "<=", "->", "&&", "||"},
		Extensions: []string{".go"},
	},
	"rust": {
		Keywords:   []string{"fn", "let", "mut", "pub", "use", "mod", "struct", "impl", "trait", "enum"},
		Types:     []string{"String", "i32", "u32", "Vec", "Option", "Result"},
		Constants: []string{"None", "Some", "Ok", "Err"},
		Operators: []string{"::", "=>", "->", "&&", "||"},
		Extensions: []string{".rs"},
	},
	"cpp": {
		Keywords:   []string{"class", "namespace", "template", "public", "private", "protected", "virtual"},
		Types:     []string{"int", "char", "bool", "void", "auto", "string", "vector"},
		Constants: []string{"nullptr", "true", "false"},
		Operators: []string{"::", "->", "<<", ">>", "&&", "||"},
		Extensions: []string{".cpp", ".hpp", ".h"},
	},
	"c": {
		Keywords:   []string{"if", "else", "while", "for", "return", "struct", "typedef", "enum"},
		Types:     []string{"int", "char", "void", "float", "double", "size_t"},
		Constants: []string{"NULL", "EOF"},
		Operators: []string{"->", "&&", "||", "<<", ">>"},
		Extensions: []string{".c", ".h"},
	},
	"javascript": {
		Keywords:   []string{"function", "const", "let", "var", "if", "else", "return", "async", "await"},
		Types:     []string{"Array", "Object", "String", "Number", "Boolean"},
		Constants: []string{"null", "undefined", "true", "false"},
		Operators: []string{"=>", "===", "!==", "&&", "||"},
		Extensions: []string{".js", ".jsx"},
	},
}

func detectLanguage(filename string) string {
	ext := filepath.Ext(filename)
	for lang, info := range languages {
		for _, langExt := range info.Extensions {
			if ext == langExt {
				return lang
			}
		}
	}
	return ""
}

func tokenize(content string, lang Language) []Token {
	var tokens []Token
	
	// Convert content to string slice
	words := strings.Fields(content)
	
	for _, word := range words {
		token := Token{Type: "text", Value: word}
		
		// Check keywords
		for _, keyword := range lang.Keywords {
			if word == keyword {
				token.Type = "keyword"
				break
			}
		}
		
		// Check types
		for _, typ := range lang.Types {
			if word == typ {
				token.Type = "type"
				break
			}
		}
		
		// Check constants
		for _, constant := range lang.Constants {
			if word == constant {
				token.Type = "constant"
				break
			}
		}
		
		// Check operators
		for _, operator := range lang.Operators {
			if word == operator {
				token.Type = "operator"
				break
			}
		}
		
		// Check for strings
		if strings.HasPrefix(word, "\"") || strings.HasPrefix(word, "'") {
			token.Type = "string"
		}
		
		// Check for numbers
		if regexp.MustCompile(`^[0-9]+$`).MatchString(word) {
			token.Type = "number"
		}
		
		tokens = append(tokens, token)
	}
	
	return tokens
}

func highlight(token Token) string {
	switch token.Type {
	case "keyword":
		return Nord9 + token.Value + Reset
	case "type":
		return Nord7 + token.Value + Reset
	case "constant":
		return Nord13 + token.Value + Reset
	case "operator":
		return Nord8 + token.Value + Reset
	case "string":
		return Nord14 + token.Value + Reset
	case "number":
		return Nord11 + token.Value + Reset
	default:
		return Nord4 + token.Value + Reset
	}
}

func printFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	lang := detectLanguage(filename)
	if lang == "" {
		fmt.Fprintf(os.Stderr, "Unsupported file type\n")
		return
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		tokens := tokenize(scanner.Text(), languages[lang])
		for _, token := range tokens {
			fmt.Printf("%s ", highlight(token))
		}
		fmt.Println()
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
	}
}

func main() {
	args := os.Args
	if len(args) < 2 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			tokens := tokenize(scanner.Text(), languages["go"]) // Default to Go for stdin
			for _, token := range tokens {
				fmt.Printf("%s ", highlight(token))
			}
			fmt.Println()
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
		}
	} else {
		for _, filename := range args[1:] {
			printFile(filename)
		}
	}
}
