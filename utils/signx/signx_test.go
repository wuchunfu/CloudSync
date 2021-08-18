package signx

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSign(t *testing.T) {
	as := assert.New(t)

	t.Run("GetMd5Sign", func(t *testing.T) {
		sign, err := GetMd5Sign("params")
		as.Nil(err)
		t.Log(sign)
	})

	t.Run("GetHmacMd5Sign", func(t *testing.T) {
		sign, err := GetHmacMd5Sign("secret", "params")
		as.Nil(err)
		t.Log(sign)
	})

	t.Run("GetSha1Sign", func(t *testing.T) {
		sign, err := GetSha1Sign("params")
		as.Nil(err)
		t.Log(sign)
	})

	t.Run("GetSha1B64Sign", func(t *testing.T) {
		sign, err := GetSha1B64Sign("params")
		as.Nil(err)
		t.Log(sign)
	})

	t.Run("GetHmacSha1Sign", func(t *testing.T) {
		sign, err := GetHmacSha1Sign("secret", "params")
		as.Nil(err)
		t.Log(sign)
	})

	t.Run("GetHmacSha1B64Sign", func(t *testing.T) {
		sign, err := GetHmacSha1B64Sign("secret", "params")
		as.Nil(err)
		t.Log(sign)
	})
}
