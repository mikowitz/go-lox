package lox

import (
	"strings"
	"testing"
)

func BenchmarkScanner_ScanTokens(b *testing.B) {
	tests := []struct {
		name   string
		source string
	}{
		{
			name:   "empty",
			source: "",
		},
		{
			name:   "single_token",
			source: "+",
		},
		{
			name:   "simple_expression",
			source: "1 + 2 * 3",
		},
		{
			name:   "variable_declaration",
			source: "var x = 42;",
		},
		{
			name:   "string_literal",
			source: `"hello world"`,
		},
		{
			name:   "number_literal",
			source: "123.456",
		},
		{
			name:   "identifier",
			source: "myVariable",
		},
		{
			name:   "keyword",
			source: "while",
		},
		{
			name:   "comment",
			source: "// this is a comment\n",
		},
		{
			name:   "multiline_string",
			source: "\"line1\nline2\nline3\"",
		},
		{
			name: "simple_function",
			source: `fun fibonacci(n) {
  if (n <= 1) return n;
  return fibonacci(n - 2) + fibonacci(n - 1);
}`,
		},
		{
			name: "class_declaration",
			source: `class Circle {
  init(radius) {
    this.radius = radius;
  }

  area() {
    return 3.14159 * this.radius * this.radius;
  }
}`,
		},
		{
			name:   "complex_expression",
			source: `var result = (x + y) * z / 2 - foo(bar, baz) >= 10 and flag == true or value != nil;`,
		},
		{
			name: "mixed_content",
			source: `// Calculate factorial
var factorial = 1;
for (var i = 1; i <= 10; i = i + 1) {
  factorial = factorial * i;
}
print factorial; // Should print 3628800`,
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			lox := &Lox{}
			b.ResetTimer()
			for b.Loop() {
				scanner := NewScanner(tt.source, lox)
				scanner.ScanTokens()
			}
		})
	}
}

func BenchmarkScanner_LargeFile(b *testing.B) {
	// Simulate a larger source file
	var sb strings.Builder

	// Add some imports/setup
	sb.WriteString("// Large test file\n")

	// Add 100 function definitions
	for i := range 100 {
		sb.WriteString("fun function")
		sb.WriteString(strings.Repeat("X", i%10)) // varying lengths
		sb.WriteString("(a, b, c) {\n")
		sb.WriteString("  var x = a + b * c;\n")
		sb.WriteString("  if (x > 100) {\n")
		sb.WriteString("    return x / 2;\n")
		sb.WriteString("  } else {\n")
		sb.WriteString("    return x * 2;\n")
		sb.WriteString("  }\n")
		sb.WriteString("}\n\n")
	}

	source := sb.String()
	lox := &Lox{}

	b.ResetTimer()
	for b.Loop() {
		scanner := NewScanner(source, lox)
		scanner.ScanTokens()
	}
}

func BenchmarkScanner_TokenTypes(b *testing.B) {
	tests := []struct {
		name   string
		source string
	}{
		{
			name:   "operators",
			source: "+ - * / ! != = == < <= > >=",
		},
		{
			name:   "punctuation",
			source: "( ) { } , . ;",
		},
		{
			name:   "keywords",
			source: "and class else false fun for if nil or print return super this true var while",
		},
		{
			name:   "numbers",
			source: "0 1 42 123.456 0.5 999.999",
		},
		{
			name:   "strings",
			source: `"" "a" "hello" "hello world" "multi\nline"`,
		},
		{
			name:   "identifiers",
			source: "a abc myVar my_var _private var123",
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			lox := &Lox{}
			b.ResetTimer()
			for b.Loop() {
				scanner := NewScanner(tt.source, lox)
				scanner.ScanTokens()
			}
		})
	}
}

func BenchmarkScanner_Unicode(b *testing.B) {
	tests := []struct {
		name   string
		source string
	}{
		{
			name:   "ascii_only",
			source: "var x = 42;",
		},
		{
			name:   "unicode_string",
			source: `"Hello 世界 🌍"`,
		},
		{
			name:   "mixed_unicode",
			source: `var message = "Hello 世界"; // Comment with émojis 🚀`,
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			lox := &Lox{}
			b.ResetTimer()
			for b.Loop() {
				scanner := NewScanner(tt.source, lox)
				scanner.ScanTokens()
			}
		})
	}
}

func BenchmarkScanner_ErrorCases(b *testing.B) {
	tests := []struct {
		name   string
		source string
	}{
		{
			name:   "unexpected_character",
			source: "@#$%",
		},
		{
			name:   "unterminated_string",
			source: `"unterminated`,
		},
		{
			name:   "mixed_valid_invalid",
			source: "var x = @ + 1;",
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ResetTimer()
			for b.Loop() {
				lox := &Lox{} // Fresh instance to reset error state
				scanner := NewScanner(tt.source, lox)
				scanner.ScanTokens()
			}
		})
	}
}

func BenchmarkScanner_MemoryAllocation(b *testing.B) {
	source := `fun example(a, b) {
  var result = a + b;
  return result * 2;
}`

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		lox := &Lox{}
		scanner := NewScanner(source, lox)
		scanner.ScanTokens()
	}
}
