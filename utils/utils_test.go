package utils

import "testing"

func TestUtils(t *testing.T) {
	ignoreFiles := []string{".git", ".idea", ".swp", ".swx"}
	fileName := "/tmp/test.swp"
	file := IgnoreFile(ignoreFiles, fileName)
	t.Log(file)
}
