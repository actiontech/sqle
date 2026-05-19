package v1

import (
	"testing"

	"golang.org/x/text/encoding/simplifiedchinese"
)

func TestProcessArchiveEntry(t *testing.T) {
	// Prepare GBK encoded content for GBK test case.
	// "SELECT * FROM 用户表;" in GBK encoding
	gbkSQL := "SELECT * FROM 用户表;"
	gbkBytes, err := simplifiedchinese.GBK.NewEncoder().Bytes([]byte(gbkSQL))
	if err != nil {
		t.Fatalf("failed to encode GBK test data: %v", err)
	}

	// Prepare Java content that contains SQL.
	// Keep it close to real-world usage so java-sql-extractor can capture it reliably.
	javaWithSQL := []byte(`import java.sql.Connection;
import java.sql.PreparedStatement;

public class Test {
    public void query(Connection conn) throws Exception {
        String sql = "SELECT * FROM users";
        PreparedStatement ps = conn.prepareStatement(sql);
        ps.executeQuery();
        ps.close();
    }
}`)

	// Prepare Java content without SQL
	javaNoSQL := []byte(`public class Test {
    public void hello() {
        System.out.println("hello");
    }
}`)

	// Prepare XML content
	xmlBytes := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE mapper PUBLIC "-//mybatis.org//DTD Mapper 3.0//EN" "http://mybatis.org/dtd/mybatis-3-mapper.dtd">
<mapper namespace="com.example.UserMapper">
    <select id="getUser" resultType="User">
        SELECT * FROM users WHERE id = #{id}
    </select>
</mapper>`)

	cases := map[string]struct {
		filename         string
		content          []byte
		expectSQL        string
		expectXMLNotNil  bool
		expectSupported  bool
		expectErr        bool
		checkSQLContains string // if non-empty, check sqlContent contains this substring instead of exact match
	}{
		"process .sql file (UTF-8)": {
			filename:        "test.sql",
			content:         []byte("SELECT 1; SELECT 2;"),
			expectSQL:       "SELECT 1; SELECT 2;",
			expectXMLNotNil: false,
			expectSupported: true,
			expectErr:       false,
		},
		"process .sql file (GBK)": {
			filename:        "test.sql",
			content:         gbkBytes,
			expectSQL:       gbkSQL,
			expectXMLNotNil: false,
			expectSupported: true,
			expectErr:       false,
		},
		"process .txt file": {
			filename:        "test.txt",
			content:         []byte("SELECT * FROM orders;"),
			expectSQL:       "SELECT * FROM orders;",
			expectXMLNotNil: false,
			expectSupported: true,
			expectErr:       false,
		},
		"process .java file (with SQL)": {
			filename:         "Test.java",
			content:          javaWithSQL,
			expectXMLNotNil:  false,
			expectSupported:  true,
			expectErr:        false,
			checkSQLContains: "SELECT * FROM users",
		},
		"process .java file (no SQL)": {
			filename:        "Test.java",
			content:         javaNoSQL,
			expectSQL:       "",
			expectXMLNotNil: false,
			expectSupported: true,
			expectErr:       false,
		},
		"process .xml file": {
			filename:        "mapper.xml",
			content:         xmlBytes,
			expectXMLNotNil: true,
			expectSupported: true,
			expectErr:       false,
		},
		"process unsupported format (.png)": {
			filename:        "image.png",
			content:         []byte{0x89, 0x50, 0x4E, 0x47},
			expectSQL:       "",
			expectXMLNotNil: false,
			expectSupported: false,
			expectErr:       false,
		},
		"process empty file (empty.sql)": {
			filename:        "empty.sql",
			content:         []byte{},
			expectSQL:       "",
			expectXMLNotNil: false,
			expectSupported: true,
			expectErr:       false,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			sqlContent, xmlContent, isSupported, err := processArchiveEntry(tc.filename, tc.content)

			// Check error
			if tc.expectErr && err == nil {
				t.Errorf("expected error but got nil")
				return
			}
			if !tc.expectErr && err != nil {
				t.Errorf("expected no error but got: %v", err)
				return
			}

			// Check isSupported
			if isSupported != tc.expectSupported {
				t.Errorf("expected isSupported=%v, got %v", tc.expectSupported, isSupported)
			}

			// Check xmlContent
			if tc.expectXMLNotNil && xmlContent == nil {
				t.Errorf("expected xmlContent to be non-nil, got nil")
			}
			if !tc.expectXMLNotNil && xmlContent != nil {
				t.Errorf("expected xmlContent to be nil, got non-nil")
			}

			// Check sqlContent
			if tc.checkSQLContains != "" {
				if len(sqlContent) == 0 {
					t.Errorf("expected sqlContent to contain %q, but sqlContent is empty", tc.checkSQLContains)
				} else if !containsSubstring(sqlContent, tc.checkSQLContains) {
					t.Errorf("expected sqlContent to contain %q, got %q", tc.checkSQLContains, sqlContent)
				}
			} else if sqlContent != tc.expectSQL {
				t.Errorf("expected sqlContent=%q, got %q", tc.expectSQL, sqlContent)
			}

			// Additional check: for .xml, verify filename is preserved
			if tc.expectXMLNotNil && xmlContent != nil {
				if xmlContent.FilePath != tc.filename {
					t.Errorf("expected xmlContent.FilePath=%q, got %q", tc.filename, xmlContent.FilePath)
				}
				if xmlContent.Content == "" {
					t.Errorf("expected xmlContent.Content to be non-empty")
				}
			}
		})
	}
}

// containsSubstring checks if s contains substr (simple helper to avoid importing strings in test).
func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
