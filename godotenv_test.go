package env

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

var noopPresets = make(map[string]string)

func loadEnvAndCompareValues(
	t *testing.T,
	loader func(files ...string) error,
	envFileName string,
	expectedValues map[string]string,
	presets map[string]string) {
	// first up, clear the env
	os.Clearenv()

	for k, v := range presets {
		os.Setenv(k, v)
	}

	err := loader(envFileName)
	if err != nil {
		t.Fatalf("Error loading %v", envFileName)
	}

	for k := range expectedValues {
		envValue := os.Getenv(k)
		v := expectedValues[k]
		if envValue != v {
			t.Errorf("Mismatch for key '%v': expected '%#v' got '%#v'", k, v, envValue)
		}
	}
}

func TestLoadWithNoArgsLoadsDotEnv(t *testing.T) {
	err := Load()
	pathError := err.(*os.PathError)
	if pathError == nil || pathError.Op != "open" || pathError.Path != ".env" {
		t.Errorf("Didn't try and open .env by default")
	}
}

func TestLoadFileNotFound(t *testing.T) {
	if err := Load("somefilethatwillneverexistever.env"); err == nil {
		t.Error("File wasn't found but Load didn't return an error")
	}
}

func TestLoadDoesNotOverride(t *testing.T) {
	envFileName := "fixtures/plain.env"

	// ensure NO overload
	presets := map[string]string{
		"OPTION_A": "do_not_override",
		"OPTION_B": "",
	}

	expectedValues := map[string]string{
		"OPTION_A": "do_not_override",
		"OPTION_B": "",
	}
	loadEnvAndCompareValues(t, Load, envFileName, expectedValues, presets)
}

func TestLoadPlainEnv(t *testing.T) {
	envFileName := "fixtures/plain.env"
	expectedValues := map[string]string{
		"OPTION_A": "1",
		"OPTION_B": "2",
		"OPTION_C": "3",
		"OPTION_D": "4",
		"OPTION_E": "5",
		"OPTION_H": "1 2",
	}

	loadEnvAndCompareValues(t, Load, envFileName, expectedValues, noopPresets)
}

func TestLoadExportedEnv(t *testing.T) {
	envFileName := "fixtures/exported.env"
	expectedValues := map[string]string{
		"OPTION_A": "2",
		"OPTION_B": "\\n",
	}

	loadEnvAndCompareValues(t, Load, envFileName, expectedValues, noopPresets)
}

func TestLoadEqualsEnv(t *testing.T) {
	envFileName := "fixtures/equals.env"
	expectedValues := map[string]string{
		"OPTION_A": "postgres://localhost:5432/database?sslmode=disable",
	}

	loadEnvAndCompareValues(t, Load, envFileName, expectedValues, noopPresets)
}

func TestLoadQuotedEnv(t *testing.T) {
	envFileName := "fixtures/quoted.env"
	expectedValues := map[string]string{
		"OPTION_A": "1",
		"OPTION_B": "2",
		"OPTION_C": "",
		"OPTION_D": "\\n",
		"OPTION_E": "1",
		"OPTION_F": "2",
		"OPTION_G": "",
		"OPTION_H": "\n",
		"OPTION_I": "echo 'asd'",
		"OPTION_J": "line 1\nline 2",
		"OPTION_K": "line one\nthis is \\'quoted\\'\none more line",
		"OPTION_L": "line 1\nline 2",
		"OPTION_M": "line one\nthis is \"quoted\"\none more line",
	}

	loadEnvAndCompareValues(t, Load, envFileName, expectedValues, noopPresets)
}

func TestOverloadWithNoArgsOverloadsDotEnv(t *testing.T) {
	err := Overload()
	pathError := err.(*os.PathError)
	if pathError == nil || pathError.Op != "open" || pathError.Path != ".env" {
		t.Errorf("Didn't try and open .env by default")
	}
}

func TestOverloadFileNotFound(t *testing.T) {
	if err := Overload("somefilethatwillneverexistever.env"); err == nil {
		t.Error("File wasn't found but Overload didn't return an error")
	}
}

func TestOverloadDoesOverride(t *testing.T) {
	envFileName := "fixtures/plain.env"

	// ensure NO overload
	presets := map[string]string{
		"OPTION_A": "do_not_override",
	}

	expectedValues := map[string]string{
		"OPTION_A": "1",
	}

	loadEnvAndCompareValues(t, Overload, envFileName, expectedValues, presets)
}

func TestReadPlainEnv(t *testing.T) {
	envFileName := "fixtures/plain.env"
	expectedValues := map[string]string{
		"OPTION_A": "1",
		"OPTION_B": "2",
		"OPTION_C": "3",
		"OPTION_D": "4",
		"OPTION_E": "5",
		"OPTION_F": "",
		"OPTION_G": "",
		"OPTION_H": "1 2",
	}

	envMap, err := Read(envFileName)
	if err != nil {
		t.Error("Error reading file")
	}

	if len(envMap) != len(expectedValues) {
		t.Error("Didn't get the right size map back")
	}

	for key, value := range expectedValues {
		if envMap[key] != value {
			t.Error("Read got one of the keys wrong")
		}
	}
}

func TestParse(t *testing.T) {
	envMap, err := Parse(bytes.NewReader([]byte("ONE=1\nTWO='2'\nTHREE = \"3\"")))
	expectedValues := map[string]string{
		"ONE":   "1",
		"TWO":   "2",
		"THREE": "3",
	}
	if err != nil {
		t.Fatalf("error parsing env: %v", err)
	}
	for key, value := range expectedValues {
		if envMap[key] != value {
			t.Errorf("expected %s to be %s, got %s", key, value, envMap[key])
		}
	}
}

func TestSubstitutions(t *testing.T) {
	envFileName := "fixtures/substitutions.env"
	expectedValues := map[string]string{
		"OPTION_A": "1",
		"OPTION_B": "1",
		"OPTION_C": "1",
		"OPTION_D": "11",
		"OPTION_E": "",
	}

	loadEnvAndCompareValues(t, Load, envFileName, expectedValues, noopPresets)
}

func TestExpanding(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]string
	}{
		{
			"expands variables found in values",
			"FOO=test\nBAR=$FOO",
			map[string]string{"FOO": "test", "BAR": "test"},
		},
		{
			"parses variables wrapped in brackets",
			"FOO=test\nBAR=${FOO}bar",
			map[string]string{"FOO": "test", "BAR": "testbar"},
		},
		{
			"expands undefined variables to an empty string",
			"BAR=$FOO",
			map[string]string{"BAR": ""},
		},
		{
			"expands variables in double quoted strings",
			"FOO=test\nBAR=\"quote $FOO\"",
			map[string]string{"FOO": "test", "BAR": "quote test"},
		},
		{
			"does not expand variables in single quoted strings",
			"BAR='quote $FOO'",
			map[string]string{"BAR": "quote $FOO"},
		},
		{
			"does not expand escaped variables",
			`FOO="foo\$BAR"`,
			map[string]string{"FOO": "foo$BAR"},
		},
		{
			"does not expand escaped variables",
			`FOO="foo\${BAR}"`,
			map[string]string{"FOO": "foo${BAR}"},
		},
		{
			"does not expand escaped variables",
			"FOO=test\nBAR=\"foo\\${FOO} ${FOO}\"",
			map[string]string{"FOO": "test", "BAR": "foo${FOO} test"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env, err := Parse(strings.NewReader(tt.input))
			if err != nil {
				t.Errorf("Error: %s", err.Error())
			}
			for k, v := range tt.expected {
				if strings.Compare(env[k], v) != 0 {
					t.Errorf("Expected: %s, Actual: %s", v, env[k])
				}
			}
		})
	}
}

func TestVariableStringValueSeparator(t *testing.T) {
	input := "TEST_URLS=\"stratum+tcp://stratum.antpool.com:3333\nstratum+tcp://stratum.antpool.com:443\""
	want := map[string]string{
		"TEST_URLS": "stratum+tcp://stratum.antpool.com:3333\nstratum+tcp://stratum.antpool.com:443",
	}
	got, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Error(err)
	}

	if len(got) != len(want) {
		t.Fatalf(
			"unexpected value:\nwant:\n\t%#v\n\ngot:\n\t%#v", want, got)
	}

	for k, wantVal := range want {
		gotVal, ok := got[k]
		if !ok {
			t.Fatalf("key %q doesn't present in result", k)
		}
		if wantVal != gotVal {
			t.Fatalf(
				"mismatch in %q value:\nwant:\n\t%s\n\ngot:\n\t%s", k,
				wantVal, gotVal)
		}
	}
}

func TestTrailingNewlines(t *testing.T) {
	cases := map[string]struct {
		input string
		key   string
		value string
	}{
		"Simple value without trailing newline": {
			input: "KEY=value",
			key:   "KEY",
			value: "value",
		},
		"Value with internal whitespace without trailing newline": {
			input: "KEY=value value",
			key:   "KEY",
			value: "value value",
		},
		"Value with internal whitespace with trailing newline": {
			input: "KEY=value value\n",
			key:   "KEY",
			value: "value value",
		},
		"YAML style - value with internal whitespace without trailing newline": {
			input: "KEY: value value",
			key:   "KEY",
			value: "value value",
		},
		"YAML style - value with internal whitespace with trailing newline": {
			input: "KEY: value value\n",
			key:   "KEY",
			value: "value value",
		},
	}

	for n, c := range cases {
		t.Run(n, func(t *testing.T) {
			result, err := Unmarshal(c.input)
			if err != nil {
				t.Errorf("Input: %q Unexpected error:\t%q", c.input, err)
			}
			if result[c.key] != c.value {
				t.Errorf("Input %q Expected:\t %q/%q\nGot:\t %q", c.input, c.key, c.value, result)
			}
		})
	}
}

func TestErrorReadDirectory(t *testing.T) {
	envFileName := "fixtures/"
	envMap, err := Read(envFileName)
	if err == nil {
		t.Errorf("Expected error, got %v", envMap)
	}
}

func TestErrorParsing(t *testing.T) {
	envFileName := "fixtures/invalid1.env"
	envMap, err := Read(envFileName)
	if err == nil {
		t.Errorf("Expected error, got %v", envMap)
	}
}

func TestLinesToIgnore(t *testing.T) {
	cases := map[string]struct {
		input string
		want  string
	}{
		"Line with nothing but line break": {
			input: "\n",
		},
		"Line with nothing but windows-style line break": {
			input: "\r\n",
		},
		"Line full of whitespace": {
			input: "\t\t ",
		},
		"Comment": {
			input: "# Comment",
		},
		"Indented comment": {
			input: "\t # comment",
		},
		"non-ignored value": {
			input: `export OPTION_B='\n'`,
			want:  `export OPTION_B='\n'`,
		},
	}

	for n, c := range cases {
		t.Run(n, func(t *testing.T) {
			got := string(getStatementStart([]byte(c.input)))
			if got != c.want {
				t.Errorf("Expected:\t %q\nGot:\t %q", c.want, got)
			}
		})
	}
}
