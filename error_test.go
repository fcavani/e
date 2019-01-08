// Copyright 2015 Felipe A. Cavani. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Start date:		2013-05-08
// Last modification:	2013-

package e

import (
	"errors"
	"io"
	"strconv"
	"testing"
)

// Silly errors
var (
	ErrDummy   = errors.New("dummy error")
	ErrSilly   = errors.New("silly error")
	ErrAnother = errors.New("another error")
	ErrStill   = errors.New("still a error")
)

// ErrStr is a simple string error.
const ErrStr = "string error"

func TestNew(t *testing.T) {
	dummy := New(ErrDummy)
	if dummy.(*Error).String() != ErrDummy.Error() {
		t.Fatal("Invalid error:", dummy.(*Error).String())
	}
	str := New(ErrStr)
	if str.(*Error).String() != ErrStr {
		t.Fatal("Invalid error:", str.(*Error).String())
	}
	// if dummy.Error() != "github.com/fcavani/e.TestNew - e/error_test.go - 24: dummy error" {
	// 	t.Fatal("Wrong debug info:", dummy.Error())
	// }
}

const trace = `github.com/fcavani/e.TestPush - e/error_test.go - 53: still a error
github.com/fcavani/e.TestPush - e/error_test.go - 52: another error
github.com/fcavani/e.TestPush - e/error_test.go - 51: silly error
github.com/fcavani/e.TestPush - e/error_test.go - 50: string error
github.com/fcavani/e.TestPush - e/error_test.go - 49: dummy error
`

const tracep1 = `github.com/fcavani/e.TestPush - e/error_test.go - 61: string error
github.com/fcavani/e.TestPush - e/error_test.go - 49: dummy error
`

func TestPush(t *testing.T) {
	dummy := New(ErrDummy)
	str := dummy.(*Error).Push(ErrStr)
	silly := str.Push(ErrSilly)
	another := silly.Push(ErrAnother)
	still := another.Push(ErrStill)
	n := still.Push(nil)
	if n != nil {
		t.Fatal("Not nil.")
	}
}

const trace2 = `github.com/fcavani/e.TestForward - e/error_test.go - 84: dummy error
github.com/fcavani/e.TestForward - e/error_test.go - 83: dummy error
github.com/fcavani/e.TestForward - e/error_test.go - 82: dummy error
`
const trace3 = `github.com/fcavani/e.TestForward - e/error_test.go - 98: another error
github.com/fcavani/e.TestForward - e/error_test.go - 97: another error
`

const trace4 = `github.com/fcavani/e.TestForward - e/error_test.go - 102: silly error
`

const trace5 = `github.com/fcavani/e.TestForward - e/error_test.go - 106: string error
`

func TestForward(t *testing.T) {
	f3 := Forward(nil)
	if f3 != nil {
		t.Fatal("Not nil.")
	}
	e := new(Error)
	f3 = Forward(e)
	if f3 != nil {
		t.Fatalf("Not nil. (2) %#v\n", f3)
	}
}

func TestEqual(t *testing.T) {
	goerror := errors.New(ErrStr)
	str1 := New(ErrStr).(*Error)
	str2 := New(ErrStr).(*Error)
	dummy := New(ErrDummy).(*Error)

	if !str1.Equal(goerror) {
		t.Fatal("Plain error failed.")
	}
	if !str1.Equal(ErrStr) {
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
	if dummy.Equal(ErrStr) {
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
	if !Equal(str1, ErrStr) {
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
	if Equal(dummy, ErrStr) {
		t.Fatal("Dummy equals to string error.")
	}
	if Equal(dummy, str1) {
		t.Fatal("Dummy equals to Error.")
	}
}

func TestFind(t *testing.T) {
	dummy := New(ErrDummy).(*Error)
	str := dummy.Push(ErrStr)
	silly := str.Push(ErrSilly)
	another := silly.Push(ErrAnother)
	still := another.Push(ErrStill)
	if still.Find(ErrSilly) != 2 {
		t.Fatal("Find failed.")
	}
	if still.Find(ErrStr) != 3 {
		t.Fatal("Find failed.")
	}
	if still.Find("not found") != -1 {
		t.Fatal("Find failed.")
	}
	if Find(still, ErrSilly) != 2 {
		t.Fatal("Find failed.")
	}
	if Find(still, ErrStr) != 3 {
		t.Fatal("Find failed.")
	}
	if Find(io.EOF, io.EOF) != 0 {
		t.Fatal("Find failed.")
	}
}

func TestPhrase(t *testing.T) {
	if Phrase(New("foo bar")) != "Foo bar." {
		t.Fatal("Phrase failed")
	}
}

func TestContains(t *testing.T) {
	dummy := New(ErrDummy)
	if !Contains(dummy, "dummy") {
		t.Fatal("Contains failed")
	}
	if !Contains(ErrDummy, "dummy") {
		t.Fatal("Contains failed")
	}
	if !Contains("ipsilum iptisum putsun dummy", "dummy") {
		t.Fatal("Contains failed")
	}
}

func TestFindStr(t *testing.T) {
	err := New(ErrDummy).(*Error).Push(ErrStr).Push(ErrSilly).Push(ErrAnother).Push(ErrStill)
	if deep := FindStr(err, "dummy"); deep != 4 {
		t.Fatal("FindStr failed:", deep)
	}
	if deep := FindStr(err, "string"); deep != 3 {
		t.Fatal("FindStr failed:", deep)
	}
	if deep := FindStr(err, "silly"); deep != 2 {
		t.Fatal("FindStr failed:", deep)
	}
	if deep := FindStr(err, "another"); deep != 1 {
		t.Fatal("FindStr failed:", deep)
	}
	if deep := FindStr(err, "still"); deep != 0 {
		t.Fatal("FindStr failed:", deep)
	}
	if deep := FindStr(err, "bláblá"); deep != -1 {
		t.Fatal("FindStr failed:", deep)
	}
}

func TestNil(t *testing.T) {
	var err error
	err = New(nil)
	if err != nil {
		t.Fatalf("not nil: %#v", err)
	}
}

func TestCopy1(t *testing.T) {
	err := New(ErrDummy).(*Error).Push(ErrStr).Push(ErrSilly).Push(ErrAnother).Push(ErrStill)
	cp := Copy(err)
	if deep := FindStr(cp, "dummy"); deep != 4 {
		t.Fatal("FindStr failed:", deep)
	}
	if deep := FindStr(cp, "string"); deep != 3 {
		t.Fatal("FindStr failed:", deep)
	}
	if deep := FindStr(cp, "silly"); deep != 2 {
		t.Fatal("FindStr failed:", deep)
	}
	if deep := FindStr(cp, "another"); deep != 1 {
		t.Fatal("FindStr failed:", deep)
	}
	if deep := FindStr(cp, "still"); deep != 0 {
		t.Fatal("FindStr failed:", deep)
	}
}

func TestCopy2(t *testing.T) {
	err := errors.New("blá")
	cp := Copy(err)
	if cp.Error() != "blá" {
		t.Fatal("copy failed", cp.Error())
	}
}

func TestMerge(t *testing.T) {
	err := Merge(New("1"), Merge("2", "3")).(*Error)
	count := 3
	for e := err; e != nil; e = e.next {
		i, err := strconv.Atoi(e.err.Error())
		if err != nil {
			t.Fatal(err)
		}
		if i != count {
			t.Fatal("errors aren't in order")
		}
		count--
	}
	err = Merge(New("1"), Merge("2", errors.New("3"))).(*Error)
	count = 3
	for e := err; e != nil; e = e.next {
		i, err := strconv.Atoi(e.err.Error())
		if err != nil {
			t.Fatal(err)
		}
		if i != count {
			t.Fatal("errors aren't in order")
		}
		count--
	}
	er := Merge(nil, nil)
	if er != nil {
		t.Fatal("not nil")
	}
	er = Merge(nil, New("1"))
	if er == nil {
		t.Fatal("nil")
	}
	if er.(*Error).err.Error() != "1" {
		t.Fatal("error value is wrong")
	}
	er = Merge(New("1"), nil)
	if er == nil {
		t.Fatal("nil")
	}
	if er.(*Error).err.Error() != "1" {
		t.Fatal("error value is wrong")
	}
}
