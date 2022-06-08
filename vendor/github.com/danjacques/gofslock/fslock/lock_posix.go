// Copyright 2017 by Dan Jacques. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin linux freebsd netbsd openbsd android solaris

package fslock

import (
	"fmt"
	"golang.org/x/sys/unix"
	"os"
	"sync"
	"syscall"
)

var globalPosixLockState posixLockState

// lockImpl is an implementation of lock using POSIX locks via flock().
//
// flock() locks are released when *any* file handle to the locked file is
// closed. To address this, we hold actual file handles globally. Attempts to
// acquire a file lock first check to see if there is already a global entity
// holding the lock (fail), then attempt to acquire the lock at a filesystem
// level.
func lockImpl(l *L) (Handle, error) {
	return globalPosixLockState.lockImpl(l)
}

// posixLockEntry is an entry in a posixLockState represneting a single inode.
type posixLockEntry struct {
	// file is the underlying open file descriptor.
	file *os.File

	// isShared is true if this is a shared lock entry, false otherwise.
	shared bool
	// sharedCount, if "shared" is true, is the number of in-process open shared
	// handles.
	sharedCount uint64
}

// posixLockState maintains an internal state of filesystem locks.
//
// For runtime usage, this is maintained in the global variable,
// globalPosixLockState.
type posixLockState struct {
	sync.RWMutex
	held map[uint64]*posixLockEntry
}

func (pls *posixLockState) lockImpl(l *L) (Handle, error) {
	fd, err := getOrCreateLockFile(l.Path, l.Content)
	if err != nil {
		return nil, err
	}
	defer func() {
		// Close "fd". On success, we'll clear "fd", so this will become a no-op.
		if fd != nil {
			fd.Close()
		}
	}()

	st, err := fd.Stat()
	if err != nil {
		return nil, err
	}
	stat := st.Sys().(*syscall.Stat_t)

	// Do we already have a lock on this file?
	pls.RLock()
	ple := pls.held[stat.Ino]
	pls.RUnlock()

	if ple != nil {
		// If we are requesting an exclusive lock, or if "ple" is held exclusively,
		// then deny the request.
		if !(l.Shared && ple.shared) {
			return nil, ErrLockHeld
		}
	}

	// Attempt to register the lock.
	pls.Lock()
	defer pls.Unlock()

	// Check again, with write lock held.
	if ple := pls.held[stat.Ino]; ple != nil {
		if !(l.Shared && ple.shared) {
			return nil, ErrLockHeld
		}

		// We're requesting a shared lock, and "ple" is shared, so we can grant a
		// handle.
		ple.sharedCount++
		return &posixLockHandle{pls, ple, stat.Ino, true}, nil
	}

	// Use "flock()" to get a lock on the file.
	//
	lockcmd := unix.F_SETLK
	lockstr := unix.Flock_t{
		Type: unix.F_WRLCK, Start: 0, Len: 0, Whence: 1,
	}

	if l.Shared {
		lockstr.Type = unix.F_RDLCK
	}

	if err := unix.FcntlFlock(fd.Fd(), lockcmd, &lockstr); err != nil {
		if errno, ok := err.(syscall.Errno); ok {
			switch errno {
			case syscall.EWOULDBLOCK:
				// Someone else holds the lock on this file.
				return nil, ErrLockHeld
			default:
				return nil, err
			}
		}
		return nil, err
	}

	if pls.held == nil {
		pls.held = make(map[uint64]*posixLockEntry)
	}

	ple = &posixLockEntry{
		file:        fd,
		shared:      l.Shared,
		sharedCount: 1, // Ignored for exclusive.
	}
	pls.held[stat.Ino] = ple
	fd = nil // Don't Close in defer().
	return &posixLockHandle{pls, ple, stat.Ino, ple.shared}, nil
}

type posixLockHandle struct {
	pls    *posixLockState
	ple    *posixLockEntry
	ino    uint64
	shared bool
}

func (l *posixLockHandle) Unlock() error {
	if l.pls == nil {
		panic("lock is not held")
	}

	l.pls.Lock()
	defer l.pls.Unlock()

	ple := l.pls.held[l.ino]
	if ple == nil {
		panic(fmt.Errorf("lock for inode %d is not held", l.ino))
	}
	if l.shared {
		if !ple.shared {
			panic(fmt.Errorf("lock for inode %d is not shared, but handle is shared", l.ino))
		}
		ple.sharedCount--
	}

	if !ple.shared || ple.sharedCount == 0 {
		// Last holder of the lock. Clean it up and unregister.
		if err := ple.file.Close(); err != nil {
			return err
		}
		delete(l.pls.held, l.ino)
	}

	// Clear the lock's "pls" field so that future calls to Unlock will fail
	// immediately.
	l.pls = nil
	return nil
}

func (l *posixLockHandle) LockFile() *os.File { return l.ple.file }

func getOrCreateLockFile(path string, content []byte) (*os.File, error) {
	const mode = 0640 | os.ModeTemporary

	// Loop until we've either created or opened the file.
	for {
		// Attempt to open the file. This will succeed if the file already exists.
		fd, err := os.OpenFile(path, os.O_RDWR, mode)
		switch {
		case err == nil:
			// Successfully opened the file, return handle.
			return fd, nil

		case os.IsNotExist(err):
			// The file doesn't exist. Attempt to exclusively create it.
			//
			// If this fails, the file exists, so we will try opening it again.
			fd, err := os.OpenFile(path, (os.O_CREATE | os.O_EXCL | os.O_RDWR), mode)
			switch {
			case err == nil:
				// Successfully created the new file. If we have content to write, try
				// and write it.
				if len(content) > 0 {
					// Failure to write content is non-fatal.
					_, _ = fd.Write(content)
				}
				return fd, err

			case os.IsExist(err):
				// Loop, we will try to open the file.

			default:
				return nil, err
			}

		default:
			return nil, err
		}
	}
}
