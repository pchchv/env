package env

import (
	"os"
	"testing"
)

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
