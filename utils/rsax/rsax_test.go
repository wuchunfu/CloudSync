package rsax

import "testing"

func TestRSA(t *testing.T) {
	t.Run("NewRSAFile", func(t *testing.T) {
		pubKeyFileName := "id_rsa.pub"
		priKeyFileName := "id_rsa"
		keyLength := 4096
		err := NewRSAFile(pubKeyFileName, priKeyFileName, keyLength)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("", func(t *testing.T) {
		keyLength := 4096
		pubKeyString, priKeyString, err := NewRSAString(keyLength)
		if err != nil {
			t.Error(err)
		}
		t.Log("\n" + pubKeyString)
		t.Log("\n" + priKeyString)
	})
}
