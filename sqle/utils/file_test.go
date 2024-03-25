package utils

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// TestMkdirWithPermission 测试不同的权限设置
func TestMkdirWithPermission(t *testing.T) {
	// 临时目录前缀
	testDirPrefix := "test_dir_" + t.Name()

	// 不同的权限测试用例
	permTestCases := []os.FileMode{
		0755, // 常见的目录权限，用户可读写执行，组和其他可读执行
		0700, // 用户可读写执行，组和其他无权限
		0777, // 所有用户可读写执行（不推荐，因为安全性问题）
		0600, // 用户可读写，组和其他无权限
	}

	for _, perm := range permTestCases {
		// 创建临时目录
		testDir := fmt.Sprintf("%s_%o", testDirPrefix, perm)
		defer os.RemoveAll(testDir)

		// 调用函数创建目录并设置权限
		err := MkdirWithPermission(testDir, perm)
		if err != nil {
			t.Errorf("MkdirWithPermission(%s, %o) error = %v", testDir, perm, err)
			continue
		}

		// 检查目录是否存在
		if _, err := os.Stat(testDir); os.IsNotExist(err) {
			t.Errorf("MkdirWithPermission(%s, %o) failed to create directory", testDir, perm)
			continue
		}

		// 检查目录权限是否正确
		info, err := os.Stat(testDir)
		if err != nil {
			t.Errorf("MkdirWithPermission(%s, %o) error getting directory info: %v", testDir, perm, err)
			continue
		}
		if info.Mode().Perm() != perm {
			t.Errorf("MkdirWithPermission(%s, %o) incorrect permissions set: %v", testDir, perm, info.Mode())
		}
	}
}

// TestEnsureFilePermission 是测试函数，包含了多个测试用例
func TestEnsureFilePermission(t *testing.T) {
	// 创建临时文件
	file, err := os.CreateTemp("", "test_ensure_file_permission")
	if err != nil {
		t.Fatal("Failed to create temporary file:", err)
	}
	defer os.Remove(file.Name()) // 确保在测试结束后删除临时文件

	// 测试用例1: 文件不存在，期望创建文件并设置权限
	testCase1 := struct {
		filePath string
		perm     os.FileMode
	}{
		filePath: file.Name(),
		perm:     0755, // 期望的文件权限
	}
	testEnsureFilePermissionCase(t, testCase1)

	// 测试用例2: 文件已存在，但权限不匹配，期望更改权限
	testCase2 := struct {
		filePath string
		perm     os.FileMode
	}{
		filePath: file.Name(),
		perm:     0644, // 期望的文件权限，与创建时的权限不同
	}
	testEnsureFilePermissionCase(t, testCase2)

	// 测试用例3: 文件已存在，且权限已匹配，期望不更改权限
	testCase3 := struct {
		filePath string
		perm     os.FileMode
	}{
		filePath: file.Name(),
		perm:     0755, // 文件创建时的权限，期望不更改
	}
	testEnsureFilePermissionCase(t, testCase3)

}

// testEnsureFilePermissionCase 是一个辅助函数，用于测试单个用例
func testEnsureFilePermissionCase(t *testing.T, testCase struct {
	filePath string
	perm     os.FileMode
}) {
	t.Run("Test case", func(t *testing.T) {
		// 首先，设置文件的初始权限
		err := os.Chmod(testCase.filePath, 0666) // 设置一个不同的权限
		if err != nil {
			t.Fatal("Failed to change file permission before test:", err)
		}

		// 调用EnsureFilePermission函数
		err = EnsureFilePermission(testCase.filePath, testCase.perm)
		if err != nil {
			t.Errorf("EnsureFilePermission returned an error: %v", err)
		}

		// 检查文件权限是否正确
		info, err := os.Stat(testCase.filePath)
		if err != nil {
			t.Error("Failed to stat file after EnsureFilePermission:", err)
		}
		if info.Mode().Perm()&testCase.perm != testCase.perm {
			t.Errorf("File permission does not match expected. Got %v, expected %v", info.Mode().Perm(), testCase.perm)
		}
	})
}

// TestSaveFile 测试 SaveFile 函数的所有可能情况
func TestSaveFile(t *testing.T) {
	// 测试用例集合
	testCases := []struct {
		name     string
		content  string
		expected error
	}{
		// 测试用例1: 成功保存空文件
		{"save_empty_file", "", nil},
		// 测试用例2: 成功保存非空文件
		{"save_non_empty_file", "Hello, World!", nil},
		// 测试用例3: 保存文件到不存在的目录
		{"/save_to_non_existing_dir/file_name", "Hello, World!", os.ErrNotExist},
	}

	// 循环测试每个用例
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 创建模拟的io.ReadSeeker
			reader := bytes.NewReader([]byte(tc.content))

			// 创建临时文件路径
			targetPath := "temp_" + tc.name + ".txt"
			defer os.Remove(targetPath) // 测试结束后删除临时文件

			// 调用SaveFile函数
			err := SaveFile(reader, targetPath)
			if tc.expected != nil {
				if err == nil || !errors.Is(err, tc.expected) {
					t.Errorf("Expected error %v, but got %v", tc.expected, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				} else {
					// 验证文件内容是否正确
					content, err := os.ReadFile(targetPath)
					if err != nil {
						t.Errorf("Failed to read the saved file: %v", err)
					} else if string(content) != tc.content {
						t.Errorf("Saved file content does not match. Expected %s, got %s", tc.content, content)
					}
				}
			}
		})
	}
}

// 这里是测试文件和目录的辅助函数
func createTestFile(t *testing.T, filePath string, perm os.FileMode) {
	err := os.WriteFile(filePath, []byte("test"), perm)
	if err != nil {
		t.Fatal(err)
	}
}

func createTestDir(t *testing.T, dirPath string, perm os.FileMode) {
	err := os.MkdirAll(dirPath, perm)
	if err != nil {
		t.Fatal(err)
	}
}

func removeTestFileOrDir(t *testing.T, path string) {
	err := os.RemoveAll(path)
	if err != nil {
		t.Fatal(err)
	}
}

// 以下是单元测试
func TestEnsureFilePathWithPermission(t *testing.T) {
	// 创建一个临时目录用于测试

	// 测试文件不存在，创建文件并设置权限
	t.Run("CreateFileWithPermission", func(t *testing.T) {
		testDir := t.TempDir()
		defer removeTestFileOrDir(t, testDir)

		// 测试文件和测试目录的名称
		testFilePath := filepath.Join(testDir, "testfile.txt")
		// 确保测试文件不存在
		removeTestFileOrDir(t, testFilePath)

		// 尝试创建文件并设置权限
		err := EnsureFilePathWithPermission(testFilePath, 0644)
		if err != nil {
			t.Fatalf("Failed to create file with permission: %v", err)
		}

		// 检查文件是否存在
		_, err = os.Stat(testFilePath)
		if os.IsNotExist(err) {
			t.Fatal("File does not exist after creation")
		}

		// 检查文件权限是否正确
		info, err := os.Stat(testFilePath)
		if err != nil {
			t.Fatal(err)
		}
		if info.Mode().Perm()&0644 != 0644 {
			t.Fatalf("File permission is incorrect: got %v, want 0644", info.Mode().Perm())
		}
	})

	// 测试目录不存在，创建目录并设置权限
	t.Run("CreateDirWithPermission", func(t *testing.T) {
		testDir := t.TempDir()
		defer removeTestFileOrDir(t, testDir)

		// 测试文件和测试目录的名称
		testDirPath := filepath.Join(testDir, "testdir")
		// 确保测试目录不存在
		removeTestFileOrDir(t, testDirPath)

		// 尝试创建目录并设置权限
		err := EnsureFilePathWithPermission(testDirPath, 0755)
		if err != nil {
			t.Fatalf("Failed to create directory with permission: %v", err)
		}

		// 检查目录是否存在
		_, err = os.Stat(testDirPath)
		if os.IsNotExist(err) {
			t.Fatal("Directory does not exist after creation")
		}

		// 检查目录权限是否正确
		info, err := os.Stat(testDirPath)
		if err != nil {
			t.Fatal(err)
		}
		if info.Mode()&os.ModeDir != os.ModeDir || info.Mode().Perm()&0755 != 0755 {
			t.Fatalf("Directory permission is incorrect: got %v, want ModeDir|0755", info.Mode())
		}
	})

	// 测试文件已存在，检查并设置正确的权限
	t.Run("CheckExistingFilePermission", func(t *testing.T) {
		testDir := t.TempDir()
		defer removeTestFileOrDir(t, testDir)

		// 测试文件和测试目录的名称
		testFilePath := filepath.Join(testDir, "testfile.txt")
		createTestFile(t, testFilePath, 0755)

		// 尝试设置文件权限，预期权限已经是正确的
		err := EnsureFilePathWithPermission(testFilePath, 0644)
		if err != nil {
			t.Fatalf("Failed to ensure file permission: %v", err)
		}

		// 检查文件权限是否正确
		info, err := os.Stat(testFilePath)
		if err != nil {
			t.Fatal(err)
		}
		if info.Mode().Perm()&0644 != 0644 {
			t.Fatalf("File permission is incorrect: got %v, want 0644", info.Mode().Perm())
		}
	})

	// 测试目录已存在，检查并设置正确的权限
	t.Run("CheckExistingDirPermission", func(t *testing.T) {
		testDir := t.TempDir()
		defer removeTestFileOrDir(t, testDir)

		// 测试文件和测试目录的名称
		testDirPath := filepath.Join(testDir, "testdir")
		createTestDir(t, testDirPath, 0755)

		// 尝试设置目录权限，预期权限已经是正确的
		err := EnsureFilePathWithPermission(testDirPath, os.ModeDir|0755)
		if err != nil {
			t.Fatalf("Failed to ensure directory permission: %v", err)
		}

		// 检查目录权限是否正确
		info, err := os.Stat(testDirPath)
		if err != nil {
			t.Fatal(err)
		}
		if info.Mode()&os.ModeDir != os.ModeDir || info.Mode().Perm()&0755 != 0755 {
			t.Fatalf("Directory permission is incorrect: got %v, want ModeDir|0755", info.Mode())
		}
	})
}
