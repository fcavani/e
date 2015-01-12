// Copyright 2015 Felipe A. Cavani. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Start date:		2013-05-08
// Last modification:	2013-

package e

import (
	"errors"
	"testing"
)

var DUMMYERROR = errors.New("dummy error")
var SILLYERROR = errors.New("silly error")
var ANOTHERERROR = errors.New("another error")
var STILLAERROR = errors.New("still a error")

const STRERROR = "string error"

func TestNew(t *testing.T) {
	dummy := New(DUMMYERROR)
	if dummy.String() != DUMMYERROR.Error() {
		t.Fatal("Invalid error:", dummy.String())
	}
	str := New(STRERROR)
	if str.String() != STRERROR {
		t.Fatal("Invalid error:", str.String())
	}
	if dummy.Error() != "github.com/fcavani/e.TestNew - e/error_test.go - 22: dummy error" {
		t.Fatal("Wrong debug info:", dummy.Error())
	}
}

const trace = `github.com/fcavani/e.TestPush - e/error_test.go - 51: still a error
github.com/fcavani/e.TestPush - e/error_test.go - 50: another error
github.com/fcavani/e.TestPush - e/error_test.go - 49: silly error
github.com/fcavani/e.TestPush - e/error_test.go - 48: string error
github.com/fcavani/e.TestPush - e/error_test.go - 47: dummy error
`

const tracep1 = `github.com/fcavani/e.TestPush - e/error_test.go - 59: string error
github.com/fcavani/e.TestPush - e/error_test.go - 47: dummy error
`

func TestPush(t *testing.T) {
	dummy := New(DUMMYERROR)
	str := dummy.Push(STRERROR)
	silly := str.Push(SILLYERROR)
	another := silly.Push(ANOTHERERROR)
	still := another.Push(STILLAERROR)
	if still.Trace() != trace {
		t.Fatal("Wrong trace:\n", still.Trace())
	}
	n := still.Push(nil)
	if n != nil {
		t.Fatal("Not nil.")
	}
	p1 := Push(dummy, STRERROR)
	if Trace(p1) != tracep1 {
		t.Fatal("Wrong trace:\n", Trace(p1))
	}
}

const trace2 = `github.com/fcavani/e.TestForward - e/error_test.go - 82: dummy error
github.com/fcavani/e.TestForward - e/error_test.go - 81: dummy error
github.com/fcavani/e.TestForward - e/error_test.go - 80: dummy error
`
const trace3 = `github.com/fcavani/e.TestForward - e/error_test.go - 96: another error
github.com/fcavani/e.TestForward - e/error_test.go - 95: another error
`

const trace4 = `github.com/fcavani/e.TestForward - e/error_test.go - 100: silly error
`

const trace5 = `github.com/fcavani/e.TestForward - e/error_test.go - 104: string error
`

func TestForward(t *testing.T) {
	dummy := New(DUMMYERROR)
	f1 := dummy.Forward()
	f2 := f1.Forward()
	if f2.Trace() != trace2 {
		t.Fatal("Wrong trace:\n", f2.Trace())
	}
	f3 := Forward(nil)
	if f3 != nil {
		t.Fatal("Not nil.")
	}
	e := new(Error)
	f3 = Forward(e)
	if f3 != nil {
		t.Fatalf("Not nil. (2) %#v\n", f3)
	}
	a := New(ANOTHERERROR)
	f4 := Forward(a)
	if f4.(*Error).Trace() != trace3 {
		t.Fatal("Wrong trace:\n", f4.(*Error).Trace())
	}
	f5 := Forward(SILLYERROR)
	if f5.(*Error).Trace() != trace4 {
		t.Fatal("Wrong trace:\n", f5.(*Error).Trace())
	}
	f6 := Forward(STRERROR)
	if f6.(*Error).Trace() != trace5 {
		t.Fatal("Wrong trace:\n", f6.(*Error).Trace())
	}
}

func TestEqual(t *testing.T) {
	goerror := errors.New(STRERROR)
	str1 := New(STRERROR)
	str2 := New(STRERROR)
	dummy := New(DUMMYERROR)

	if !str1.Equal(goerror) {
		t.Fatal("Plain error failed.")
	}
	if !str1.Equal(STRERROR) {
		t.Fatal("String failed.")
	}
	if !str1.Equal(str2) {
		t.Fatal("Error failed.")
	}
	if !str2.Equal(str1) {
		t.Fatal("Error failed (2).")
	}
	if str1.Equal(dummy) {
		t.Fatal("Different error failed.")
	}
	if dummy.Equal(goerror) {
		t.Fatal("Different error failed (2).")
	}
	if dummy.Equal(STRERROR) {
		t.Fatal("Different error failed (3).")
	}
	if dummy.Equal(str1) {
		t.Fatal("Different error failed (4).")
	}

	if Equal(nil, goerror) {
		t.Fatal("nil equal is false.")
	}
	if Equal(&Error{}, nil) {
		t.Fatal("nil equal is false (2).")
	}
	if Equal(&Error{}, goerror) {
		t.Fatal("nil Error is false.")
	}
	if Equal(str1, &Error{}) {
		t.Fatal("nil Error is false (2).")
	}

	if !Equal(str1, goerror) {
		t.Fatal("Not equal go error.")
	}
	if !Equal(str1, STRERROR) {
		t.Fatal("Not equal string error.")
	}
	if !Equal(str1, str2) {
		t.Fatal("Not equal string error (2).")
	}
	if !Equal(str2, str1) {
		t.Fatal("Not equal string error (3).")
	}
	if Equal(str1, dummy) {
		t.Fatal("Equal to dummy error.")
	}
	if Equal(dummy, goerror) {
		t.Fatal("Dummy equals to goerror.")
	}
	if Equal(dummy, STRERROR) {
		t.Fatal("Dummy equals to string error.")
	}
	if Equal(dummy, str1) {
		t.Fatal("Dummy equals to Error.")
	}
}

func TestFind(t *testing.T) {
	dummy := New(DUMMYERROR)
	str := dummy.Push(STRERROR)
	silly := str.Push(SILLYERROR)
	another := silly.Push(ANOTHERERROR)
	still := another.Push(STILLAERROR)
	if still.Find(SILLYERROR) != 2 {
		t.Fatal("Find failed.")
	}
	if still.Find(STRERROR) != 3 {
		t.Fatal("Find failed.")
	}
	if still.Find("not found") != -1 {
		t.Fatal("Find failed.")
	}
	if Find(still, SILLYERROR) != 2 {
		t.Fatal("Find failed.")
	}
	if Find(still, STRERROR) != 3 {
		t.Fatal("Find failed.")
	}
}

func TestPharse(t *testing.T) {
	if Phrase(New("foo bar")) != "Foo bar." {
		t.Fatal("Phrase failed")
	}
}
