package utils

import "testing"

func TestUtils(t *testing.T) {

	t.Run("IgnoreFiles", func(t *testing.T) {
		ignoreFiles := []string{".git", ".idea", ".swp", ".swx"}
		fileName := "/tmp/test.swp"
		file := IgnoreFile(ignoreFiles, fileName)
		t.Log(file)
	})

	t.Run("CurrentUser", func(t *testing.T) {
		user := CurrentUser()
		t.Log(user)
	})

	t.Run("UserHome", func(t *testing.T) {
		home, err := UserHome()
		if err != nil {
			t.Error(err)
		}
		t.Log(home)
	})

	t.Run("RandomString", func(t *testing.T) {
		randomString := RandomString(12)
		t.Log(randomString)
	})

	t.Run("GenRandomString", func(t *testing.T) {
		randomString := GenRandomString(12, true)
		t.Log(randomString)
	})
}
