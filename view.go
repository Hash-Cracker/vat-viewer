package main

import (
	"bufio"
	"fmt"
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
	Nord15 = "\033[38;2;180;142;173m" // Purple
	Reset  = "\033[0m"
)

type Token struct {
	Type  string
	Value string
}

type Language struct {
	Keywords        []string
	Types          []string
	Constants      []string
	Operators      []string
	Extensions     []string
	BuiltInFuncs   map[string]bool
	PackageFuncs   map[string]map[string]bool
	SpecialVars    []string
}

var languages = map[string]Language{
	"go": {
		Keywords:   []string{"package", "import", "func", "return", "if", "else", "for", "range", "var", "const", "defer", "type", "struct", "interface", "map", "chan", "go", "select", "case", "default", "switch"},
		Types:     []string{"string", "int", "int64", "int32", "int16", "int8", "uint", "uint64", "uint32", "uint16", "uint8", "byte", "rune", "float64", "float32", "complex64", "complex128", "error", "bool", "uintptr"},
		Constants: []string{"nil", "true", "false", "iota"},
		Operators: []string{":=", "!=", "==", ">=", "<=", "->", "&&", "||", "...", "<-", "&^"},
		Extensions: []string{".go"},
		BuiltInFuncs: map[string]bool{
			"make": true, "len": true, "cap": true, "new": true, "append": true, "copy": true, "close": true, "delete": true, 
			"complex": true, "real": true, "imag": true, "panic": true, "recover": true,
		},
		PackageFuncs: map[string]map[string]bool{
			"fmt": {
				"Printf": true, "Println": true, "Print": true, "Sprintf": true, "Fprintf": true, "Errorf": true,
			},
			"os": {
				"Open": true, "Create": true, "Remove": true, "Mkdir": true, "Exit": true, "Getenv": true,
			},
			"bufio": {
				"NewScanner": true, "NewReader": true, "NewWriter": true,
			},
			"strings": {
				"Split": true, "Join": true, "Contains": true, "HasPrefix": true, "HasSuffix": true, "Replace": true,
			},
			"regexp": {
				"Compile": true, "MustCompile": true, "MatchString": true,
			},
			"filepath": {
				"Join": true, "Base": true, "Dir": true, "Ext": true, "Clean": true,
			},
		},
		SpecialVars: []string{"err", "ctx"},
	},
}

func isFunction(word string) bool {
	return regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*\(`).MatchString(word)
}

func isVariable(word string) bool {
	return regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`).MatchString(word)
}

func tokenize(content string, lang Language) []Token {
	var tokens []Token
	
	// Split by whitespace but preserve important characters
	r := regexp.MustCompile(`[\s\(\)\[\]{},;.]`)
	parts := r.Split(content, -1)
	
	currentPackage := ""
	
	for _, part := range parts {
		if part == "" {
			continue
		}
		
		token := Token{Type: "text", Value: part}
		
		// Check if it's a package reference
		if strings.Contains(part, ".") {
			pkg := strings.Split(part, ".")[0]
			fn := strings.Split(part, ".")[1]
			if funcs, ok := lang.PackageFuncs[pkg]; ok {
				if funcs[fn] {
					token.Type = "package_func"
					currentPackage = pkg
					goto tokenFound
				}
			}
		}
		
		// Check if it's a built-in function
		if lang.BuiltInFuncs[part] {
			token.Type = "builtin_func"
			goto tokenFound
		}
		
		// Check if it's a function declaration or call
		if isFunction(part) {
			token.Type = "function"
			goto tokenFound
		}
		
		// Check keywords
		for _, keyword := range lang.Keywords {
			if part == keyword {
				token.Type = "keyword"
				goto tokenFound
			}
		}
		
		// Check types
		for _, typ := range lang.Types {
			if part == typ {
				token.Type = "type"
				goto tokenFound
			}
		}
		
		// Check constants
		for _, constant := range lang.Constants {
			if part == constant {
				token.Type = "constant"
				goto tokenFound
			}
		}
		
		// Check special variables
		for _, special := range lang.SpecialVars {
			if part == special {
				token.Type = "special_var"
				goto tokenFound
			}
		}
		
		// Check if it's a variable
		if isVariable(part) {
			token.Type = "variable"
			goto tokenFound
		}
		
		// Check operators
		for _, operator := range lang.Operators {
			if part == operator {
				token.Type = "operator"
				goto tokenFound
			}
		}
		
		// Check for strings
		if strings.HasPrefix(part, "\"") || strings.HasPrefix(part, "'") {
			token.Type = "string"
			goto tokenFound
		}
		
		// Check for numbers
		if regexp.MustCompile(`^[0-9]+$`).MatchString(part) {
			token.Type = "number"
			goto tokenFound
		}
		
	tokenFound:
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
	case "function":
		return Nord10 + token.Value + Reset
	case "builtin_func":
		return Nord15 + token.Value + Reset
	case "package_func":
		return Nord8 + token.Value + Reset
	case "variable":
		return Nord4 + token.Value + Reset
	case "special_var":
		return Nord11 + token.Value + Reset
	default:
		return Nord1 + token.Value + Reset
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

func main() {
	args := os.Args
	if len(args) < 2 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			tokens := tokenize(scanner.Text(), languages["go"])
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
