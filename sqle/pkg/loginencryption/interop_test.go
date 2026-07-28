package loginencryption

import (
	"crypto/rand"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tjfoc/gmsm/sm2"
)

func TestDecryptSmCryptoStyleCipherWithout04(t *testing.T) {
	path, priv := writeTempPrivateKey(t)
	require.NoError(t, Init(ModeRequired, path))

	plain := "SqleLoginPass#123"
	cipher, err := sm2.Encrypt(&priv.PublicKey, []byte(plain), rand.Reader, sm2.C1C3C2)
	require.NoError(t, err)
	require.Equal(t, byte(0x04), cipher[0])

	// Frontend sm-crypto.doEncrypt returns hex without leading 04.
	smCryptoStyle := hex.EncodeToString(cipher[1:])
	got, err := GetManager().ResolvePassword("", smCryptoStyle, GetManager().PublicInfo().KeyID)
	require.NoError(t, err)
	require.Equal(t, plain, got)
}
