package utils

import (
	"io"
	"io/fs"
	"os"
)

// Grant file owners read and write execution permissions, group and other users read-only permissions
const OwnerPrivilegedAccessMode fs.FileMode = 0740

func EnsureFilePathWithPermission(filePath string, perm os.FileMode) error {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		err = MkdirWithPermission(filePath, perm)
		if err != nil {
			return err
		}
	} else {
		err = EnsureFilePermission(filePath, perm)
		if err != nil {
			return err
		}
	}
	return nil
}

func MkdirWithPermission(filePath string, perm os.FileMode) error {
	err := os.Mkdir(filePath, os.ModeDir)
	if err != nil {
		return err
	}
	err = os.Chmod(filePath, perm)
	if err != nil {
		return err
	}
	return nil
}

func EnsureFilePermission(filePath string, perm os.FileMode) error {
	info, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	currentPerm := info.Mode().Perm()
	// 使用按位与来检查每一个权限位是否满足要求
	if currentPerm&perm != perm {
		return os.Chmod(filePath, perm)
	}
	return nil
}

// SaveFile 从 io.ReadSeeker 读取内容并保存到指定路径的文件中。
func SaveFile(file io.ReadSeeker, targetPath string) (err error) {
	// 从文件头开始读取
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	// 创建目标文件
	target, err := os.Create(targetPath)
	if err != nil {
		return err
	}

	// 保存文件
	_, err = io.Copy(target, file)
	if err != nil {
		return err
	}
	return nil
}
