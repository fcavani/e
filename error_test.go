// Copyright 2015 Felipe A. Cavani. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Start date:		2013-05-08
// Last modification:	2013-

package e

import (
	"errors"
	"strconv"
	"testing"
)

var DUMMYERROR = errors.New("dummy error")
var SILLYERROR = errors.New("silly error")
var ANOTHERERROR = errors.New("another error")
var STILLAERROR = errors.New("still a error")

const STRERROR = "string error"

func TestNew(t *testing.T) {
	dummy := New(DUMMYERROR)
	if dummy.(*Error).String() != DUMMYERROR.Error() {
		t.Fatal("Invalid error:", dummy.(*Error).String())
	}
	str := New(STRERROR)
	if str.(*Error).String() != STRERROR {
		t.Fatal("Invalid error:", str.(*Error).String())
	}
	if dummy.Error() != "github.com/fcavani/e.TestNew - e/error_test.go - 23: dummy error" {
		t.Fatal("Wrong debug info:", dummy.Error())
	}
}

const trace = `github.com/fcavani/e.TestPush - e/error_test.go - 52: still a error
github.com/fcavani/e.TestPush - e/error_test.go - 51: another error
github.com/fcavani/e.TestPush - e/error_test.go - 50: silly error
github.com/fcavani/e.TestPush - e/error_test.go - 49: string error
github.com/fcavani/e.TestPush - e/error_test.go - 48: dummy error
`

const tracep1 = `github.com/fcavani/e.TestPush - e/error_test.go - 60: string error
github.com/fcavani/e.TestPush - e/error_test.go - 48: dummy error
`

func TestPush(t *testing.T) {
	dummy := New(DUMMYERROR)
	str := dummy.(*Error).Push(STRERROR)
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

const trace2 = `github.com/fcavani/e.TestForward - e/error_test.go - 83: dummy error
github.com/fcavani/e.TestForward - e/error_test.go - 82: dummy error
github.com/fcavani/e.TestForward - e/error_test.go - 81: dummy error
`
const trace3 = `github.com/fcavani/e.TestForward - e/error_test.go - 97: another error
github.com/fcavani/e.TestForward - e/error_test.go - 96: another error
`

const trace4 = `github.com/fcavani/e.TestForward - e/error_test.go - 101: silly error
`

const trace5 = `github.com/fcavani/e.TestForward - e/error_test.go - 105: string error
`

func TestForward(t *testing.T) {
	dummy := New(DUMMYERROR)
	f1 := dummy.(*Error).Forward()
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
	str1 := New(STRERROR).(*Error)
	str2 := New(STRERROR).(*Error)
	dummy := New(DUMMYERROR).(*Error)

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
	dummy := New(DUMMYERROR).(*Error)
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

func TestContains(t *testing.T) {
	dummy := New(DUMMYERROR)
	if !Contains(dummy, "dummy") {
		t.Fatal("Contains failed")
	}
	if !Contains(DUMMYERROR, "dummy") {
		t.Fatal("Contains failed")
	}
	if !Contains("ipsilum iptisum putsun dummy", "dummy") {
		t.Fatal("Contains failed")
	}
}

func TestFindStr(t *testing.T) {
	err := New(DUMMYERROR).(*Error).Push(STRERROR).Push(SILLYERROR).Push(ANOTHERERROR).Push(STILLAERROR)
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
	// err = Forward(nil)
	// if err != nil {
	// 	t.Fatalf("not nil: %#v", err)
	// }
	// err = Push(err, nil)
	// if err != nil {
	// 	t.Fatalf("not nil: %#v", err)
	// }
}

func TestCopy1(t *testing.T) {
	err := New(DUMMYERROR).(*Error).Push(STRERROR).Push(SILLYERROR).Push(ANOTHERERROR).Push(STILLAERROR)
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
