package loginencryption

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/tjfoc/gmsm/sm2"
)

// encryptForTest encrypts plaintext with current manager public key (test helper; was EncryptForTest).
func encryptForTest(plain string) (cipherHex string, keyID string, err error) {
	m := GetManager()
	if !m.Enabled() {
		return "", "", fmt.Errorf("encryption is disabled")
	}
	cipher, err := sm2.Encrypt(&m.privateKey.PublicKey, []byte(plain), rand.Reader, sm2.C1C3C2)
	if err != nil {
		return "", "", err
	}
	return hex.EncodeToString(cipher), m.keyID, nil
}

func resetManager(t *testing.T) {
	t.Helper()
	setGlobal(&Manager{mode: ModeDisabled})
	t.Cleanup(func() {
		setGlobal(&Manager{mode: ModeDisabled})
	})
}

func mustWriteKey(t *testing.T, path string) []byte {
	t.Helper()
	if err := generateAndSavePrivateKey(path); err != nil {
		t.Fatalf("generate key: %v", err)
	}
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("read key: %v", err)
	}
	return raw
}

func TestInit(t *testing.T) {
	t.Run("empty_mode_equals_required", func(t *testing.T) {
		resetManager(t)
		dir := t.TempDir()
		oldWD, err := os.Getwd()
		if err != nil {
			t.Fatal(err)
		}
		if err := os.Chdir(dir); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = os.Chdir(oldWD) })

		if err := Init("", ""); err != nil {
			t.Fatalf("Init: %v", err)
		}
		m := GetManager()
		if m.Mode() != ModeRequired {
			t.Fatalf("mode=%q, want %q", m.Mode(), ModeRequired)
		}
		if !m.Enabled() {
			t.Fatal("expected enabled")
		}
		info := m.PublicInfo()
		if !info.Enable || info.PublicKey == "" || info.KeyID == "" {
			t.Fatalf("PublicInfo incomplete: %+v", info)
		}
		if _, err := os.Stat(DefaultPrivateKeyPath); err != nil {
			t.Fatalf("default key not created: %v", err)
		}
	})

	t.Run("disabled", func(t *testing.T) {
		resetManager(t)
		if err := Init(ModeDisabled, "/nonexistent/should-be-ignored.pem"); err != nil {
			t.Fatalf("Init: %v", err)
		}
		m := GetManager()
		if m.Mode() != ModeDisabled || m.Enabled() {
			t.Fatalf("mode=%q enabled=%v", m.Mode(), m.Enabled())
		}
		info := m.PublicInfo()
		if info.Enable {
			t.Fatalf("PublicInfo.Enable=true in disabled mode")
		}
	})

	t.Run("explicit_path_missing_fails", func(t *testing.T) {
		resetManager(t)
		missing := filepath.Join(t.TempDir(), "missing", "key.pem")
		err := Init(ModeRequired, missing)
		if err == nil {
			t.Fatal("expected error for missing explicit path")
		}
		if _, statErr := os.Stat(missing); !os.IsNotExist(statErr) {
			t.Fatalf("must not create explicit missing path; stat=%v", statErr)
		}
		if GetManager().Enabled() {
			t.Fatal("manager must stay disabled after failed Init")
		}
	})

	t.Run("o_excl_does_not_overwrite_existing_key", func(t *testing.T) {
		resetManager(t)
		dir := t.TempDir()
		keyPath := filepath.Join(dir, "login_sm2_private.pem")
		before := mustWriteKey(t, keyPath)

		// Second generate must fail with O_EXCL (file exists).
		if err := generateAndSavePrivateKey(keyPath); err == nil {
			t.Fatal("expected O_EXCL failure on existing key")
		}
		afterGen, err := ioutil.ReadFile(keyPath)
		if err != nil {
			t.Fatal(err)
		}
		if string(afterGen) != string(before) {
			t.Fatal("generateAndSavePrivateKey overwrote existing key")
		}

		if err := Init(ModeCompatible, keyPath); err != nil {
			t.Fatalf("Init: %v", err)
		}
		afterInit, err := ioutil.ReadFile(keyPath)
		if err != nil {
			t.Fatal(err)
		}
		if string(afterInit) != string(before) {
			t.Fatal("Init overwrote existing key")
		}
		if GetManager().Mode() != ModeCompatible || !GetManager().Enabled() {
			t.Fatalf("mode=%q enabled=%v", GetManager().Mode(), GetManager().Enabled())
		}
	})
}

func TestResolveUsername(t *testing.T) {
	resetManager(t)
	keyPath := filepath.Join(t.TempDir(), "key.pem")
	mustWriteKey(t, keyPath)
	if err := Init(ModeRequired, keyPath); err != nil {
		t.Fatalf("Init required: %v", err)
	}
	cipher, keyID, err := encryptForTest("admin")
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}

	tests := []struct {
		name      string
		setup     func(t *testing.T)
		plain     string
		encrypted string
		keyID     string
		want      string
		wantErr   string
	}{
		{
			name:    "required_rejects_plaintext",
			plain:   "admin",
			wantErr: "plaintext username is not allowed",
		},
		{
			name:    "required_missing_ciphertext",
			wantErr: "encrypted username is required",
		},
		{
			name:      "required_key_id_mismatch",
			encrypted: cipher,
			keyID:     "deadbeefdeadbeef",
			wantErr:   "invalid key_id",
		},
		{
			name:      "required_ok",
			encrypted: cipher,
			keyID:     keyID,
			want:      "admin",
		},
		{
			name: "compatible_plaintext_escape",
			setup: func(t *testing.T) {
				if err := Init(ModeCompatible, keyPath); err != nil {
					t.Fatalf("Init compatible: %v", err)
				}
			},
			plain: "legacy-user",
			want:  "legacy-user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t)
			} else if err := Init(ModeRequired, keyPath); err != nil {
				t.Fatalf("Init required: %v", err)
			}
			got, err := GetManager().ResolveUsername(tt.plain, tt.encrypted, tt.keyID)
			if tt.wantErr != "" {
				if err == nil || err.Error() != tt.wantErr {
					t.Fatalf("err=%v, want %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestResolvePassword(t *testing.T) {
	resetManager(t)
	keyPath := filepath.Join(t.TempDir(), "key.pem")
	mustWriteKey(t, keyPath)
	if err := Init(ModeRequired, keyPath); err != nil {
		t.Fatalf("Init required: %v", err)
	}
	cipher, keyID, err := encryptForTest("secret")
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}

	tests := []struct {
		name      string
		setup     func(t *testing.T)
		plain     string
		encrypted string
		keyID     string
		want      string
		wantErr   string
	}{
		{
			name:    "required_rejects_plaintext",
			plain:   "secret",
			wantErr: "plaintext password is not allowed",
		},
		{
			name:    "required_missing_ciphertext",
			wantErr: "encrypted password is required",
		},
		{
			name:      "required_key_id_mismatch",
			encrypted: cipher,
			keyID:     "deadbeefdeadbeef",
			wantErr:   "invalid key_id",
		},
		{
			name:      "required_ok",
			encrypted: cipher,
			keyID:     keyID,
			want:      "secret",
		},
		{
			name: "compatible_plaintext_escape",
			setup: func(t *testing.T) {
				if err := Init(ModeCompatible, keyPath); err != nil {
					t.Fatalf("Init compatible: %v", err)
				}
			},
			plain: "legacy-pass",
			want:  "legacy-pass",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t)
			} else if err := Init(ModeRequired, keyPath); err != nil {
				t.Fatalf("Init required: %v", err)
			}
			got, err := GetManager().ResolvePassword(tt.plain, tt.encrypted, tt.keyID)
			if tt.wantErr != "" {
				if err == nil || err.Error() != tt.wantErr {
					t.Fatalf("err=%v, want %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}
