package cryptox

import "testing"

func TestUtils(t *testing.T) {
	md5 := GenerateMd5("/tmp/test/test.txt")
	t.Log(md5)
}
