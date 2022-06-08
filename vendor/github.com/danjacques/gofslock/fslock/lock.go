// Copyright 2017 by Dan Jacques. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fslock

import (
	"errors"
	"os"
)

// ErrLockHeld is a sentinel error returned when the lock could not be
// acquired.
var ErrLockHeld = errors.New("fslock: lock is held")

// Handle is a reference to a held lock. It must be released via Unlock when
// finished. Multiple calls to Unlock will panic.
//
// Handle is NOT safe for concurrent use.
type Handle interface {
	// Unlock releases the held lock.
	//
	// This can error if the underlying filesystem operations fail. This should
	// not happen unless something has gone externally wrong, or the lock was
	// mishandled.
	Unlock() error

	// LockFile returns the underlying lock File. This is not generally useful,
	// and should be used circumstantially. Operating on the file occurs outside
	// of the scope of this package, and can result in unintended consequences.
	//
	// The file should NOT be directly closed or modified.
	//
	// The file will be valid for the duration of the Handle. Once the Handle is
	// closed with Unlock, the file's state is implementation-specific and
	// unspecified.
	LockFile() *os.File
}

// Blocker is used for the Delay field in a Lock.
type Blocker func() error

// L describes a filesystem lock.
//
// L's fields should not be modified concurrently, but L's methods are safe
// for concurrent use.
type L struct {
	// Path is the path of the file to lock.
	Path string

	// Shared, if true, indicates that this should be a shared lock rather than
	// an exclusive lock.
	//
	// See package documentation for details.
	Shared bool

	// Content, if populated, is the lock file content. Content is written to the
	// file when the lock call creates it, and only if the lock call actually
	// creates the file. Failure to write Content is non-fatal.
	//
	// Content should be used only as a convenience hint for users who want to
	// know what the lock file is, and not for actual programmatic management.
	// Several code paths can result in successful file locking and still fail to
	// write Content to that file.
	//
	// Content is not synchronized with the actual locking. Failure to write
	// Content to the lock file is considered non-fatal.
	Content []byte

	// Block is the configured blocking function.
	//
	// If not nil, an attempt to acquire the lock will loop indefinitely until an
	// error other than ErrLockHeld is encountered (fatal) or the lock is
	// acquired. Block will be called each time a lock attempt returns
	// ErrLockHeld, and should delay and/or cancel the acquisition by returning
	// nil or an error code respectively.
	//
	// If Block returns an error, it will be propagated as the error result of the
	// locking attempt.
	Block Blocker
}

// Lock attempts to acquire the configured lock.
func (l *L) Lock() (Handle, error) {
	// Loop repeatedly until the lock is held or an error is encountered.
	for {
		switch h, err := lockImpl(l); err {
		case nil:
			// Acquired the lock.
			return h, nil

		case ErrLockHeld:
			// If we have a Block function configured, invoke it, then try again.
			// Otherwise, propagate ErrLockHeld.
			if l.Block != nil {
				if err := l.Block(); err != nil {
					return nil, err
				}
				continue
			}
			fallthrough

		default:
			return nil, err
		}
	}
}

// With is a convenience method to acquire a lock via Lock, call fn, and release
// the lock on completion (via defer).
//
// If an error is encountered, it will be returned. Otherwise, the return value
// from fn will be returned.
func (l *L) With(fn func() error) (err error) {
	h, err := l.Lock()
	if err != nil {
		return
	}
	defer func() {
		uerr := h.Unlock()
		if uerr != nil && err == nil {
			err = uerr
		}
	}()
	return fn()
}
