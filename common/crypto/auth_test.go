package crypto_test

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"testing"

	"v2ray.com/core/common/buf"
	. "v2ray.com/core/common/crypto"
	"v2ray.com/core/common/protocol"
	. "v2ray.com/ext/assert"
)

func TestAuthenticationReaderWriter(t *testing.T) {
	assert := With(t)

	key := make([]byte, 16)
	rand.Read(key)
	block, err := aes.NewCipher(key)
	assert(err, IsNil)

	aead, err := cipher.NewGCM(block)
	assert(err, IsNil)

	rawPayload := make([]byte, 8192*10)
	rand.Read(rawPayload)

	payload := buf.NewLocal(8192 * 10)
	payload.Append(rawPayload)

	cache := buf.NewLocal(160 * 1024)
	iv := make([]byte, 12)
	rand.Read(iv)

	writer := NewAuthenticationWriter(&AEADAuthenticator{
		AEAD: aead,
		NonceGenerator: &StaticBytesGenerator{
			Content: iv,
		},
		AdditionalDataGenerator: &NoOpBytesGenerator{},
	}, PlainChunkSizeParser{}, cache, protocol.TransferTypeStream)

	assert(writer.WriteMultiBuffer(buf.NewMultiBufferValue(payload)), IsNil)
	assert(cache.Len(), Equals, 83360)
	assert(writer.WriteMultiBuffer(buf.MultiBuffer{}), IsNil)
	assert(err, IsNil)

	reader := NewAuthenticationReader(&AEADAuthenticator{
		AEAD: aead,
		NonceGenerator: &StaticBytesGenerator{
			Content: iv,
		},
		AdditionalDataGenerator: &NoOpBytesGenerator{},
	}, PlainChunkSizeParser{}, cache, protocol.TransferTypeStream)

	var mb buf.MultiBuffer

	for mb.Len() < len(rawPayload) {
		mb2, err := reader.ReadMultiBuffer()
		assert(err, IsNil)

		mb.AppendMulti(mb2)
	}

	mbContent := make([]byte, 8192*10)
	mb.Read(mbContent)
	assert(mbContent, Equals, rawPayload)

	_, err = reader.ReadMultiBuffer()
	assert(err, Equals, io.EOF)
}

func TestAuthenticationReaderWriterPacket(t *testing.T) {
	assert := With(t)

	key := make([]byte, 16)
	rand.Read(key)
	block, err := aes.NewCipher(key)
	assert(err, IsNil)

	aead, err := cipher.NewGCM(block)
	assert(err, IsNil)

	cache := buf.NewLocal(1024)
	iv := make([]byte, 12)
	rand.Read(iv)

	writer := NewAuthenticationWriter(&AEADAuthenticator{
		AEAD: aead,
		NonceGenerator: &StaticBytesGenerator{
			Content: iv,
		},
		AdditionalDataGenerator: &NoOpBytesGenerator{},
	}, PlainChunkSizeParser{}, cache, protocol.TransferTypePacket)

	var payload buf.MultiBuffer
	pb1 := buf.New()
	pb1.Append([]byte("abcd"))
	payload.Append(pb1)

	pb2 := buf.New()
	pb2.Append([]byte("efgh"))
	payload.Append(pb2)

	assert(writer.WriteMultiBuffer(payload), IsNil)
	assert(cache.Len(), GreaterThan, 0)
	assert(writer.WriteMultiBuffer(buf.MultiBuffer{}), IsNil)
	assert(err, IsNil)

	reader := NewAuthenticationReader(&AEADAuthenticator{
		AEAD: aead,
		NonceGenerator: &StaticBytesGenerator{
			Content: iv,
		},
		AdditionalDataGenerator: &NoOpBytesGenerator{},
	}, PlainChunkSizeParser{}, cache, protocol.TransferTypePacket)

	mb, err := reader.ReadMultiBuffer()
	assert(err, IsNil)

	b1 := mb.SplitFirst()
	assert(b1.String(), Equals, "abcd")
	b2 := mb.SplitFirst()
	assert(b2.String(), Equals, "efgh")
	assert(mb.IsEmpty(), IsTrue)

	_, err = reader.ReadMultiBuffer()
	assert(err, Equals, io.EOF)
}
