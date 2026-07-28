package loginencryption

import (
	"crypto/rand"
	"encoding/hex"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/x509"
)

func writeTempPrivateKey(t *testing.T) (string, *sm2.PrivateKey) {
	t.Helper()
	priv, err := sm2.GenerateKey(rand.Reader)
	require.NoError(t, err)
	pemBytes, err := x509.WritePrivateKeyToPem(priv, nil)
	require.NoError(t, err)
	dir := t.TempDir()
	path := filepath.Join(dir, "sm2_priv.pem")
	require.NoError(t, ioutil.WriteFile(path, pemBytes, 0600))
	return path, priv
}

func TestInitDisabled(t *testing.T) {
	require.NoError(t, Init(ModeDisabled, ""))
	m := GetManager()
	info := m.PublicInfo()
	assert.False(t, info.Enable)
	assert.Empty(t, info.PublicKey)
	assert.Empty(t, info.KeyID)

	pwd, err := m.ResolvePassword("plain-password", "", "")
	require.NoError(t, err)
	assert.Equal(t, "plain-password", pwd)
}

func TestInitRequiresKeyWhenEnabled(t *testing.T) {
	err := Init(ModeCompatible, "")
	require.Error(t, err)
	err = Init(ModeRequired, "/tmp/not-exist-sm2-key.pem")
	require.Error(t, err)
}

func TestCompatibleAndRequiredModes(t *testing.T) {
	path, priv := writeTempPrivateKey(t)

	require.NoError(t, Init(ModeCompatible, path))
	m := GetManager()
	info := m.PublicInfo()
	assert.True(t, info.Enable)
	assert.Equal(t, AlgorithmSM2, info.Algorithm)
	assert.Equal(t, CipherModeC1C3C2, info.CipherMode)
	assert.True(t, strings.HasPrefix(info.PublicKey, "04"))
	assert.NotEmpty(t, info.KeyID)

	// compatible: plaintext still accepted
	pwd, err := m.ResolvePassword("admin123", "", "")
	require.NoError(t, err)
	assert.Equal(t, "admin123", pwd)

	cipher, err := sm2.Encrypt(&priv.PublicKey, []byte("secret-pass"), rand.Reader, sm2.C1C3C2)
	require.NoError(t, err)
	cipherHex := hex.EncodeToString(cipher)
	pwd, err = m.ResolvePassword("", cipherHex, info.KeyID)
	require.NoError(t, err)
	assert.Equal(t, "secret-pass", pwd)

	// sm-crypto style: ciphertext without leading 04
	cipherWithout04 := hex.EncodeToString(cipher[1:])
	pwd, err = m.ResolvePassword("", cipherWithout04, info.KeyID)
	require.NoError(t, err)
	assert.Equal(t, "secret-pass", pwd)

	require.NoError(t, Init(ModeRequired, path))
	m = GetManager()
	_, err = m.ResolvePassword("admin123", "", "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "encrypted password is required")

	pwd, err = m.ResolvePassword("ignored", cipherHex, info.KeyID)
	require.NoError(t, err)
	assert.Equal(t, "secret-pass", pwd)
}

func TestResolvePasswordRejectsInvalidCipher(t *testing.T) {
	path, _ := writeTempPrivateKey(t)
	require.NoError(t, Init(ModeCompatible, path))
	m := GetManager()
	info := m.PublicInfo()

	_, err := m.ResolvePassword("", "zz-not-hex", info.KeyID)
	require.Error(t, err)

	_, err = m.ResolvePassword("", hex.EncodeToString([]byte("short")), info.KeyID)
	require.Error(t, err)

	_, err = m.ResolvePassword("", strings.Repeat("ab", maxCipherHexLen/2+1), info.KeyID)
	require.Error(t, err)

	_, err = m.ResolvePassword("", "04"+strings.Repeat("00", 100), "wrong-key-id")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "key_id")

	_, err = m.ResolvePassword("", "", "")
	require.Error(t, err)
}

func TestPublicKeyStableKeyID(t *testing.T) {
	path, _ := writeTempPrivateKey(t)
	require.NoError(t, Init(ModeCompatible, path))
	info1 := GetManager().PublicInfo()
	require.NoError(t, Init(ModeCompatible, path))
	info2 := GetManager().PublicInfo()
	assert.Equal(t, info1.KeyID, info2.KeyID)
	assert.Equal(t, info1.PublicKey, info2.PublicKey)
}

func TestLoadPrivateKeyFromHexFile(t *testing.T) {
	priv, err := sm2.GenerateKey(rand.Reader)
	require.NoError(t, err)
	path := filepath.Join(t.TempDir(), "sm2.hex")
	require.NoError(t, ioutil.WriteFile(path, []byte(hex.EncodeToString(priv.D.Bytes())), 0600))
	require.NoError(t, Init(ModeRequired, path))
	info := GetManager().PublicInfo()
	assert.True(t, info.Enable)
	assert.True(t, strings.HasPrefix(info.PublicKey, "04"))
}

func TestInvalidMode(t *testing.T) {
	err := Init("unknown", "")
	require.Error(t, err)
}

func TestDecryptDoesNotLeakSecrets(t *testing.T) {
	path, _ := writeTempPrivateKey(t)
	require.NoError(t, Init(ModeCompatible, path))
	m := GetManager()
	info := m.PublicInfo()
	_, err := m.ResolvePassword("", "deadbeef", info.KeyID)
	require.Error(t, err)
	assert.NotContains(t, err.Error(), path)
	privPEM, _ := ioutil.ReadFile(path)
	assert.NotContains(t, err.Error(), string(privPEM))
}
