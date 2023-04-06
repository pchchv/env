package env

import (
	"bytes"
	"io"
	"os"
	"strings"
)

// Parse reads the env file from io.Reader,
// returning a map of keys and values.
func Parse(r io.Reader) (map[string]string, error) {
	var buf bytes.Buffer

	_, err := io.Copy(&buf, r)
	if err != nil {
		return nil, err
	}

	return UnmarshalBytes(buf.Bytes())
}

// Load reads the env file(s) and loads them into ENV for this process.
// Call this function as close as possible to the beginning of your program (ideally in main).
// If you call Load without any args, it will load the .env at the current path by default.
// Otherwise you can tell it which files to load (there can be more than one), for example:
//
//	env.Load("fileone", "filetwo")
//
// It is important to note that it DOES NOT DELETE env variables that already exist -
// use the .env file to set dev vars or reasonable defaults.
func Load(filenames ...string) (err error) {
	filenames = filenamesOrDefault(filenames)
	for _, filename := range filenames {
		err = loadFile(filename, false)
		if err != nil {
			return
		}
	}

	return
}

// Overload reads your env file(s) and loads them into ENV for this process.
// Call this function as close as possible to the beginning of program (ideally in main).
// If you call Overload without any args, it will load the .env at the current path by default.
// Otherwise you can tell it which files to load (there can be several), for example:
//
//	godotenv.Overload("fileone", "filetwo")
//
// It is important to note that this OVERRIDE an env variable that already exists -
// think of the .env file as forcibly setting all variables.
func Overload(filenames ...string) (err error) {
	filenames = filenamesOrDefault(filenames)
	for _, filename := range filenames {
		err = loadFile(filename, true)
		if err != nil {
			return
		}
	}

	return
}

// Read reads all envs (with the same load semantics as Load),
// but returns the values as a map instead of automatically writing them to the env.
func Read(filenames ...string) (envMap map[string]string, err error) {
	filenames = filenamesOrDefault(filenames)
	envMap = make(map[string]string)

	for _, filename := range filenames {
		individualEnvMap, individualErr := readFile(filename)

		if individualErr != nil {
			err = individualErr
			return
		}

		for key, value := range individualEnvMap {
			envMap[key] = value
		}
	}

	return
}

// Unmarshal reads the env file from the string,
// returning a map of keys and values.
func Unmarshal(str string) (envMap map[string]string, err error) {
	return UnmarshalBytes([]byte(str))
}

// UnmarshalBytes parses env file from byte slices of characters,
// returning a map of keys and values.
func UnmarshalBytes(src []byte) (map[string]string, error) {
	out := make(map[string]string)
	err := parseBytes(src, out)

	return out, err
}

func filenamesOrDefault(filenames []string) []string {
	if len(filenames) == 0 {
		return []string{".env"}
	}

	return filenames
}

func readFile(filename string) (envMap map[string]string, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	return Parse(file)
}

func loadFile(filename string, overload bool) error {
	envMap, err := readFile(filename)
	if err != nil {
		return err
	}

	currentEnv := map[string]bool{}
	rawEnv := os.Environ()
	for _, rawEnvLine := range rawEnv {
		key := strings.Split(rawEnvLine, "=")[0]
		currentEnv[key] = true
	}

	for key, value := range envMap {
		if !currentEnv[key] || overload {
			_ = os.Setenv(key, value)
		}
	}

	return nil
}
