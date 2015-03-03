// Copyright 2015 Felipe A. Cavani. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Start date:		2013-05-08
// Last modification:	2013-

// Helper functions to manipulate errors and trace then.
package e

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"unicode"

	"github.com/fcavani/util/types"
)

// Error expand go error type with debug information and error trace.
type Error struct {
	err  error
	args []interface{}
	pkg  string
	// File name
	file string
	// Line where the error occurred
	line      int
	debugInfo bool
	next      *Error
}

func init() {
	types.Insert(&Error{})
}

type GoError string

func (g GoError) Error() string {
	return string(g)
}

type messageType uint8

const (
	NextIsNill messageType = iota
	Next
	ErrorGo
	ErrorLocal
)

func (e *Error) GobEncode() ([]byte, error) {
	var err error
	buf := bytes.NewBuffer([]byte{})
	enc := gob.NewEncoder(buf)
	switch v := e.err.(type) {
	case *Error:
		err = enc.Encode(ErrorLocal)
		if err != nil {
			return nil, err
		}
		err = enc.Encode(v)
		if err != nil {
			return nil, err
		}
	case error:
		err = enc.Encode(ErrorGo)
		if err != nil {
			return nil, err
		}
		err = enc.Encode(GoError(v.Error()))
		if err != nil {
			return nil, err
		}
	default:
		panic("type not supported")
	}
	err = enc.Encode(e.args)
	if err != nil {
		return nil, err
	}
	err = enc.Encode(e.pkg)
	if err != nil {
		return nil, err
	}
	err = enc.Encode(e.file)
	if err != nil {
		return nil, err
	}
	err = enc.Encode(e.line)
	if err != nil {
		return nil, err
	}
	err = enc.Encode(e.debugInfo)
	if err != nil {
		return nil, err
	}
	if e.next == nil {
		err = enc.Encode(NextIsNill)
		if err != nil {
			return nil, err
		}
	} else {
		err = enc.Encode(Next)
		if err != nil {
			return nil, err
		}
		err = enc.Encode(e.next)
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func (e *Error) GobDecode(data []byte) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	var msg messageType
	err := dec.Decode(&msg)
	if err != nil {
		return err
	}
	switch msg {
	case ErrorLocal:
		var err_ *Error
		err := dec.Decode(&err_)
		if err != nil {
			return err
		}
		e.err = err_
	case ErrorGo:
		var err_ GoError
		err := dec.Decode(&err_)
		if err != nil {
			return err
		}
		e.err = err_
	default:
		return errors.New("protocol error")
	}
	err = dec.Decode(&e.args)
	if err != nil {
		return err
	}
	err = dec.Decode(&e.pkg)
	if err != nil {
		return err
	}
	err = dec.Decode(&e.file)
	if err != nil {
		return err
	}
	err = dec.Decode(&e.line)
	if err != nil {
		return err
	}
	err = dec.Decode(&e.debugInfo)
	if err != nil {
		return err
	}
	err = dec.Decode(&msg)
	if err != nil {
		return err
	}
	switch msg {
	case NextIsNill:
		return nil
	case Next:
		err = dec.Decode(&e.next)
		if err != nil {
			return err
		}
	default:
		return errors.New("protocol error")
	}
	return nil
}

func (e *Error) formatError() string {
	return fmt.Sprintf(e.err.Error(), e.args...)
}

// Error return the packed, file, line number and the error message.
func (e *Error) Error() string {
	if e.debugInfo {
		return fmt.Sprintf("%v - %v - %v: %v", e.pkg, e.file, strconv.Itoa(e.line), e.formatError())
	}
	return e.formatError()
}

// String return only the error message of the last error.
func (e *Error) String() string {
	return e.err.Error()
}

// GoString return the same as Error function, more verbose.
func (e *Error) GoString() string {
	if e.debugInfo {
		return fmt.Sprintf("package: %v - file: %v - line: %v - error: %v", e.pkg, e.file, strconv.Itoa(e.line), e.formatError())
	}
	return fmt.Sprintf("%#v", e.formatError())
}

func (e *Error) Arguments() []interface{} {
	return e.args
}

// Transform an error message in something readable.
func Phrase(i interface{}) string {
	msg := ""
	switch val := i.(type) {
	case fmt.Stringer:
		msg = val.String()
	case error:
		msg = val.Error()
	case string:
		msg = val
	default:
		panic("invalid type, must be Stringer or error")
	}
	return string(unicode.ToUpper(rune(msg[0]))) + msg[1:] + "."
}

// String return the string associated with the error
func String(i interface{}) string {
	msg := ""
	switch val := i.(type) {
	case fmt.Stringer:
		msg = val.String()
	case error:
		msg = val.Error()
	case string:
		msg = val
	default:
		panic("invalid type, must be Stringer or error")
	}
	return msg
}

func (e *Error) last() *Error {
	prev := e
	next := e.next
	for next != nil {
		prev = next
		next = next.next
	}
	return prev
}

func (e *Error) push(ie interface{}, n int) *Error {
	var err *Error
	if e2, ok := ie.(*Error); ok {
		err = e2.last()
	} else {
		err = newError(ie, n)
	}
	if err == nil {
		return nil
	}
	err.next = e
	return err
}

// Push one error on the top of the stack. ie must be *Error, error or string.
func (e *Error) Push(ie interface{}) *Error {
	return e.push(ie, 3)
}

// Push e2 error on the top of the stack (e1 error). Free function to use with other
// types of error beside the *Error. e1 must be *Error or error
// and e2 must be *Error, error or string.
func Push(e1, e2 interface{}) error {
	return PushN(e1, e2, 1)
}

func PushN(e1, e2 interface{}, n int) error {
	if e1 == nil {
		if e2b, ok := e2.(*Error); ok {
			return e2b.forward(3 + n)
		}
		return newError(e2, 2+n)
	}
	switch val := e1.(type) {
	case *Error:
		return val.push(e2, 3+n)
	case error:
		return newError(val, 2+n).push(e2, 3+n)
	case string:
		return newError(val, 2+n).push(e2, 3+n)
	default:
		panic("invalid type, e1 must be *Error")
	}
	return nil
}

func (e *Error) forward(n int) *Error {
	if e.err == nil {
		return nil
	}
	ne := newError(e, n)
	ne.next = e
	return ne
}

// Forward the error. Only stack the error menssage and the debug data.
func (e *Error) Forward() *Error {
	return e.forward(3)
}

// Forward the error. Only stack the error menssage and the debug data.
// Free function to use with other types of error besides the *Error.
// ie must be *Error and r must be *Error, error or string.
func Forward(ie interface{}) error {
	return ForwardN(ie, 1)
}

// ForwardN skip n levels from the statck when login the
// trace.
func ForwardN(ie interface{}, n int) error {
	if ie == nil {
		return nil
	}
	switch val := ie.(type) {
	case *Error:
		if val == nil {
			return nil
		}
		ret := val.forward(3 + n)
		if ret == nil {
			return nil
		}
		return ret
	case error:
		return newError(val, 2+n)
	case string:
		return newError(val, 2+n)
	default:
		panic("invalid type")
	}
	return nil
}

//Equal compare if the errors are the same. ie must be *Error and r must be
// *Error, error or string.
func (e *Error) Equal(ie interface{}) bool {
	if ie == nil || e == nil || e.err == nil {
		return false
	}
	switch val := ie.(type) {
	case *Error:
		if val == nil {
			return false
		}
		if val.err == nil {
			return false
		}
		return e.err.Error() == val.String()
	case error:
		if val == nil {
			return false
		}
		return e.err.Error() == val.Error()
	case string:
		return e.err.Error() == val
	default:
		panic("invalid type")
	}
	return false
}

// Equal compare if the errors are the same. l must be *Error and r must be
// *Error, error or string.
func Equal(l, r interface{}) bool {
	if l == nil || r == nil {
		return false
	}
	switch val := l.(type) {
	case *Error:
		return val.Equal(r)
	case error:
		return newError(val, 2).Equal(r)
	default:
		panic("invalid type, must be *Error")
	}
	return false
}

// Find an error in the chain. ie must be
// *Error, error or string.
func (e *Error) Find(ie interface{}) int {
	if ie == nil {
		return -1
	}
	deep := 0
	for err := e; err != nil; err = err.next {
		if err.Equal(ie) {
			return deep
		}
		deep = deep + 1
	}
	return -1
}

// Find an error in the chain. e must be *Error and ie must be
// *Error, error or string.
func Find(e, ie interface{}) int {
	if e == nil || ie == nil {
		return -1
	}
	switch err := e.(type) {
	case *Error:
		return err.Find(ie)
	default:
		panic("invalid type, must be *Error")
	}
	return -1
}

// Trace the error and return a string.
func (e *Error) Trace() (s string) {
	for err := e; err != nil; err = err.next {
		s = s + fmt.Sprintln(err)
	}
	return
}

// Trace the error and return a string. ie must be *Error.
func Trace(ie interface{}) string {
	if ie == nil {
		return "nil"
	}
	switch val := ie.(type) {
	case *Error:
		return val.Trace()
	default:
		panic("invalid type, must be *Error")
	}
	return ""
}

func newError(ie interface{}, level int, a ...interface{}) (err *Error) {
	if ie == nil {
		return nil
	}
	var e error
	switch val := ie.(type) {
	case *Error:
		if val == nil || val.err == nil {
			return nil
		}
		e = val.err
		a = val.args
	case error:
		if val == nil {
			return nil
		}
		e = val
	case string:
		e = errors.New(val)
	default:
		panic("invalid type")
	}
	pc, file, line, ok := runtime.Caller(level)
	if ok {
		s := strings.Split(file, "/")
		l := len(s)
		var file string
		if l >= 2 {
			file = strings.Join(s[l-2:l], "/")
		} else {
			file = s[0]
		}
		f := runtime.FuncForPC(pc)
		err = &Error{
			err:       e,
			args:      a,
			pkg:       f.Name(),
			file:      file,
			line:      line,
			debugInfo: true,
		}
		return
	}
	err = &Error{
		err:       e,
		args:      a,
		debugInfo: false,
	}
	return
}

// New initiates an error from a string, error or *Error. a is
// the verb in the error string that will be replaced when
// Error and GoString functions is called. The valids verbs are
// the same verbs in the fmt package.
func New(ie interface{}, a ...interface{}) *Error {
	return newError(ie, 2, a...)
}

func NewN(ie interface{}, n int, a ...interface{}) *Error {
	return newError(ie, 2+n, a...)
}

// Contains checks if the error message contains the sub string.
func (e *Error) Contains(sub string) bool {
	if e.err == nil {
		return false
	}
	return strings.Contains(e.formatError(), sub)
}

// Contains checks if the error message contains the sub string.
func Contains(ie interface{}, sub string) bool {
	if ie == nil || sub == "" {
		return false
	}
	switch val := ie.(type) {
	case *Error:
		if val == nil {
			return false
		}
		return val.Contains(sub)
	case error:
		if val == nil {
			return false
		}
		return strings.Contains(val.Error(), sub)
	case string:
		return strings.Contains(val, sub)
	default:
		panic("invalid type")
	}
	panic("don't get here")
}

// FindStr find a sub string int the chain of error and return
// the deep of the error.
func (e *Error) FindStr(sub string) int {
	deep := 0
	for err := e; err != nil; err = err.next {
		if err.Contains(sub) {
			return deep
		}
		deep = deep + 1
	}
	return -1
}

func FindStr(ie interface{}, sub string) int {
	if ie == nil || sub == "" {
		return -1
	}
	switch val := ie.(type) {
	case *Error:
		if val == nil {
			return -1
		}
		return val.FindStr(sub)
	default:
		panic("invalid type")
	}
	panic("don't get here")
}
