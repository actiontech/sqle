[![Build Status](https://travis-ci.com/bodgit/windows.svg?branch=master)](https://travis-ci.com/bodgit/windows)
[![Coverage Status](https://coveralls.io/repos/github/bodgit/windows/badge.svg?branch=master)](https://coveralls.io/github/bodgit/windows?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/bodgit/windows)](https://goreportcard.com/report/github.com/bodgit/windows)
[![GoDoc](https://godoc.org/github.com/bodgit/windows?status.svg)](https://godoc.org/github.com/bodgit/windows)

windows
=======

A collection of types native to Windows but are useful on non-Windows platforms.

The `FILETIME`-comparable type is the sole export which is a 1:1 copy of the one from `golang.org/x/sys/windows`. However, that package isn't available for all platforms and this particular type gets used in other protocols and file types such as NTLMv2 and 7-Zip.
