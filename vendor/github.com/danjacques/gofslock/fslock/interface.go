// Copyright 2017 by Dan Jacques. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package fslock is a cross-platform filesystem locking implementation.
//
// fslock aims to implement workable filesystem-based locking semantics. Each
// implementation offers its own nuances, and not all of those nuances are
// addressed by this package.
//
// Locks can either be exclusive (default) or shared. For a given file, exactly
// one of the following circumstances may be true at a given moment:
//	- No locks are held.
//	- A single exclusive lock may be held.
//	- Multiple shared locks may be held.
//
// Notably, an exclusive lock may not be held on a file that has shared locks,
// and a shared lock may not be held on a file that has an exclusive lock.
//
// fslock will work as long as you don't do anything particularly weird, such
// as:
//
//	- Manually forking your Go process.
//	- Circumventing the fslock package and opening and/or manipulating the lock
//	  file directly.
//
// An attempt to take a filesystem lock is non-blocking, and will return
// ErrLockHeld if the lock is already held elsewhere.
package fslock

// Lock acquires a filesystem lock for the given path.
//
// If the lock could not be acquired because it is held by another entity,
// ErrLockHeld will be returned. If an error is encountered while locking,
// that error will be returned.
//
// Lock is a convenience method for L's Lock.
func Lock(path string) (Handle, error) { return LockBlocking(path, nil) }

// LockBlocking acquires an exclusive filesystem lock for the given path. If the
// lock is already held, LockBlocking will repeatedly attempt to acquire it
// using the supplied Blocker in between attempts.
//
// If no Blocker is provided and the lock could not be acquired because it is
// held by another entity, ErrLockHeld will be returned. If an error is
// encountered while locking, or an error is returned by b, that error will be
// returned.
//
// LockBlocking is a convenience method for L's Lock.
func LockBlocking(path string, b Blocker) (Handle, error) {
	l := L{
		Path:  path,
		Block: b,
	}
	return l.Lock()
}

// With is a convenience function to create a lock, execute a function while
// holding that lock, and then release the lock on completion.
//
// See L's With method for details.
func With(path string, fn func() error) error { return WithBlocking(path, nil, fn) }

// WithBlocking is a convenience function to create a lock, execute a function
// while holding that lock, and then release the lock on completion. The
// supplied block function is used to retry (see L's Block field).
//
// See L's With method for details.
func WithBlocking(path string, b Blocker, fn func() error) error {
	l := L{
		Path:  path,
		Block: b,
	}
	return l.With(fn)
}

// Lock acquires a filesystem lock for the given path.
//
// If the lock could not be acquired because it is held by another entity,
// ErrLockHeld will be returned. If an error is encountered while locking,
// that error will be returned.
//
// Lock is a convenience method for L's Lock.
func LockShared(path string) (Handle, error) { return LockSharedBlocking(path, nil) }

// LockSharedBlocking acquires a shared filesystem lock for the given path. If
// the lock is already held, LockSharedBlocking will repeatedly attempt to
// acquire it using the supplied Blocker in between attempts.
//
// If no Blocker is provided and the lock could not be acquired because it is
// held by another entity, ErrLockHeld will be returned. If an error is
// encountered while locking, or an error is returned by b, that error will be
// returned.
//
// LockSharedBlocking is a convenience method for L's Lock.
func LockSharedBlocking(path string, b Blocker) (Handle, error) {
	l := L{
		Path:   path,
		Shared: true,
		Block:  b,
	}
	return l.Lock()
}

// WithShared is a convenience function to create a lock, execute a function while
// holding that lock, and then release the lock on completion.
//
// See L's With method for details.
func WithShared(path string, fn func() error) error { return WithSharedBlocking(path, nil, fn) }

// WithSharedBlocking is a convenience function to create a lock, execute a
// function while holding that lock, and then release the lock on completion.
// The supplied block function is used to retry (see L's Block field).
//
// See L's With method for details.
func WithSharedBlocking(path string, b Blocker, fn func() error) error {
	l := L{
		Path:   path,
		Shared: true,
		Block:  b,
	}
	return l.With(fn)
}
