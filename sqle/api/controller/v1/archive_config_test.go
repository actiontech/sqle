package v1

import "testing"

func TestArchiveConfig_CheckSize(t *testing.T) {
	cfg := archiveConfig{
		MaxTotalSize:    10 * 1024 * 1024, // 10MB
		MaxFileCount:    1000,
		MaxNestingDepth: 1,
	}

	cases := map[string]struct {
		currentTotal int64
		entrySize    int64
		expectErr    bool
	}{
		"size within limit: 5MB + 3MB": {
			currentTotal: 5 * 1024 * 1024,
			entrySize:    3 * 1024 * 1024,
			expectErr:    false,
		},
		"size exceeds limit: 9MB + 2MB": {
			currentTotal: 9 * 1024 * 1024,
			entrySize:    2 * 1024 * 1024,
			expectErr:    true,
		},
		"size exactly at limit: 7MB + 3MB": {
			currentTotal: 7 * 1024 * 1024,
			entrySize:    3 * 1024 * 1024,
			expectErr:    false,
		},
		"size one byte over limit": {
			currentTotal: 10 * 1024 * 1024,
			entrySize:    1,
			expectErr:    true,
		},
		"both zero": {
			currentTotal: 0,
			entrySize:    0,
			expectErr:    false,
		},
		"entry size equals max total size": {
			currentTotal: 0,
			entrySize:    10 * 1024 * 1024,
			expectErr:    false,
		},
		"current total equals max total size with zero entry": {
			currentTotal: 10 * 1024 * 1024,
			entrySize:    0,
			expectErr:    false,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			err := cfg.checkSize(tc.currentTotal, tc.entrySize)
			if tc.expectErr && err == nil {
				t.Errorf("expected error but got nil")
			}
			if !tc.expectErr && err != nil {
				t.Errorf("expected no error but got: %v", err)
			}
		})
	}
}

func TestArchiveConfig_CheckFileCount(t *testing.T) {
	cfg := archiveConfig{
		MaxTotalSize:    10 * 1024 * 1024,
		MaxFileCount:    1000,
		MaxNestingDepth: 1,
	}

	cases := map[string]struct {
		count     int
		expectErr bool
	}{
		"file count within limit: 999": {
			count:     999,
			expectErr: false,
		},
		"file count exceeds limit: 1001": {
			count:     1001,
			expectErr: true,
		},
		"file count exactly at limit: 1000": {
			count:     1000,
			expectErr: false,
		},
		"file count one over limit: 1001": {
			count:     1001,
			expectErr: true,
		},
		"file count zero": {
			count:     0,
			expectErr: false,
		},
		"file count is 1": {
			count:     1,
			expectErr: false,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			err := cfg.checkFileCount(tc.count)
			if tc.expectErr && err == nil {
				t.Errorf("expected error but got nil")
			}
			if !tc.expectErr && err != nil {
				t.Errorf("expected no error but got: %v", err)
			}
		})
	}
}

func TestDefaultArchiveConfig(t *testing.T) {
	if defaultArchiveConfig.MaxTotalSize != 10*1024*1024 {
		t.Errorf("expected MaxTotalSize to be 10MB, got %d", defaultArchiveConfig.MaxTotalSize)
	}
	if defaultArchiveConfig.MaxFileCount != 1000 {
		t.Errorf("expected MaxFileCount to be 1000, got %d", defaultArchiveConfig.MaxFileCount)
	}
	if defaultArchiveConfig.MaxNestingDepth != 1 {
		t.Errorf("expected MaxNestingDepth to be 1, got %d", defaultArchiveConfig.MaxNestingDepth)
	}
}

func TestSupportedArchiveExts(t *testing.T) {
	expectedExts := []string{".zip", ".rar", ".7z"}
	for _, ext := range expectedExts {
		if !supportedArchiveExts[ext] {
			t.Errorf("expected %s to be in supportedArchiveExts", ext)
		}
	}

	unexpectedExts := []string{".tar", ".gz", ".bz2"}
	for _, ext := range unexpectedExts {
		if supportedArchiveExts[ext] {
			t.Errorf("expected %s to NOT be in supportedArchiveExts", ext)
		}
	}
}

func TestSupportedTextExts(t *testing.T) {
	expectedExts := []string{".sql", ".txt", ".java", ".xlsx"}
	for _, ext := range expectedExts {
		if !supportedTextExts[ext] {
			t.Errorf("expected %s to be in supportedTextExts", ext)
		}
	}

	unexpectedExts := []string{".png", ".jpg", ".pdf"}
	for _, ext := range unexpectedExts {
		if supportedTextExts[ext] {
			t.Errorf("expected %s to NOT be in supportedTextExts", ext)
		}
	}
}
