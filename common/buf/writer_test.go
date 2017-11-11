package buf_test

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"testing"

	"context"
	"io"

	"v2ray.com/core/common"
	. "v2ray.com/core/common/buf"
	"v2ray.com/core/transport/ray"
	. "v2ray.com/ext/assert"
)

func TestWriter(t *testing.T) {
	assert := With(t)

	lb := New()
	assert(lb.AppendSupplier(ReadFrom(rand.Reader)), IsNil)

	expectedBytes := append([]byte(nil), lb.Bytes()...)

	writeBuffer := bytes.NewBuffer(make([]byte, 0, 1024*1024))

	writer := NewBufferedWriter(NewWriter(writeBuffer))
	writer.SetBuffered(false)
	err := writer.WriteMultiBuffer(NewMultiBufferValue(lb))
	assert(err, IsNil)
	assert(writer.Flush(), IsNil)
	assert(expectedBytes, Equals, writeBuffer.Bytes())
}

func TestBytesWriterReadFrom(t *testing.T) {
	assert := With(t)

	cache := ray.NewStream(context.Background())
	reader := bufio.NewReader(io.LimitReader(rand.Reader, 8192))
	writer := NewBufferedWriter(cache)
	writer.SetBuffered(false)
	_, err := reader.WriteTo(writer)
	assert(err, IsNil)

	mb, err := cache.ReadMultiBuffer()
	assert(err, IsNil)
	assert(mb.Len(), Equals, 8192)
}

func TestDiscardBytes(t *testing.T) {
	assert := With(t)

	b := New()
	common.Must(b.Reset(ReadFullFrom(rand.Reader, Size)))

	nBytes, err := io.Copy(DiscardBytes, b)
	assert(nBytes, Equals, int64(Size))
	assert(err, IsNil)
}

func TestDiscardBytesMultiBuffer(t *testing.T) {
	assert := With(t)

	const size = 10240*1024 + 1
	buffer := bytes.NewBuffer(make([]byte, 0, size))
	common.Must2(buffer.ReadFrom(io.LimitReader(rand.Reader, size)))

	r := NewReader(buffer)
	nBytes, err := io.Copy(DiscardBytes, NewBufferedReader(r))
	assert(nBytes, Equals, int64(size))
	assert(err, IsNil)
}
