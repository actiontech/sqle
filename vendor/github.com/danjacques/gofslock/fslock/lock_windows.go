// Copyright 2017 by Dan Jacques. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fslock

import (
	"os"
	"syscall"
	"unsafe"
)

const errno_ERROR_SHARING_VIOLATION syscall.Errno = 32

func lockImpl(l *L) (Handle, error) {
	fd, created, err := getOrCreateFile(l.Path, l.Shared)
	if err != nil {
		return nil, err
	}

	// If we just created the file, write Contents to it. If this fails, it is
	// non-fatal.
	if created && len(l.Content) > 0 {
		_, _ = fd.Write(l.Content)
	}

	// We own the lock on virtue of having accessed the file exclusively.
	return &winLockHandle{fd}, nil
}

type winLockHandle struct {
	fd *os.File
}

func (h *winLockHandle) Unlock() error {
	if h.fd == nil {
		panic("lock is not held")
	}
	if err := h.fd.Close(); err != nil {
		return err
	}
	h.fd = nil
	return nil
}

func (h *winLockHandle) LockFile() *os.File { return h.fd }

func getOrCreateFile(path string, shared bool) (*os.File, bool, error) {
	mod := syscall.NewLazyDLL("kernel32.dll")
	proc := mod.NewProc("CreateFileW")

	pathp, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return nil, false, err
	}

	dwShareMode := uint32(0) // Exclusive (no sharing)
	if shared {
		dwShareMode = syscall.FILE_SHARE_READ | syscall.FILE_SHARE_WRITE
	}

	a, _, err := proc.Call(
		uintptr(unsafe.Pointer(pathp)),
		uintptr(syscall.GENERIC_READ|syscall.GENERIC_WRITE),
		uintptr(dwShareMode),
		uintptr(0), // No security attributes.
		uintptr(syscall.OPEN_ALWAYS),
		uintptr(syscall.FILE_ATTRIBUTE_NORMAL),
		0, // No template file.
	)
	fd := syscall.Handle(a)
	errno := err.(syscall.Errno)
	switch errno {
	case 0:
		return os.NewFile(uintptr(fd), path), true, nil
	case syscall.ERROR_ALREADY_EXISTS:
		// Opened the file, but did not create it.
		return os.NewFile(uintptr(fd), path), false, nil
	case errno_ERROR_SHARING_VIOLATION:
		// We could not open the file because someone else is holding it.
		return nil, false, ErrLockHeld
	default:
		syscall.CloseHandle(fd)
		return nil, false, err
	}
}
