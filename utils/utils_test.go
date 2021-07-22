package utils

import "testing"

func TestUtils(t *testing.T) {
	md5 := GetFileMd5("/tmp/test/test.txt")
	t.Log(md5)
}
