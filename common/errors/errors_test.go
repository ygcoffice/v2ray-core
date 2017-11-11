package errors_test

import (
	"io"
	"testing"

	. "v2ray.com/core/common/errors"
	. "v2ray.com/ext/assert"
)

func TestError(t *testing.T) {
	assert := With(t)

	err := New("TestError")
	assert(GetSeverity(err), Equals, SeverityInfo)

	err = New("TestError2").Base(io.EOF)
	assert(GetSeverity(err), Equals, SeverityInfo)

	err = New("TestError3").Base(io.EOF).AtWarning()
	assert(GetSeverity(err), Equals, SeverityWarning)

	err = New("TestError4").Base(io.EOF).AtWarning()
	err = New("TestError5").Base(err)
	assert(GetSeverity(err), Equals, SeverityWarning)
	assert(err.Error(), HasSubstring, "EOF")
}

func TestErrorMessage(t *testing.T) {
	assert := With(t)

	data := []struct {
		err error
		msg string
	}{
		{
			err: New("a").Base(New("b")).Path("c", "d", "e"),
			msg: "c|d|e: a > b",
		},
		{
			err: New("a").Base(New("b").Path("c")).Path("d", "e"),
			msg: "d|e: a > c: b",
		},
	}

	for _, d := range data {
		assert(d.err.Error(), Equals, d.msg)
	}
}
