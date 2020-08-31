package teapack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPacketRequest(t *testing.T) {
	var data map[string]string
	var ctx map[string]string
	a := assert.New(t)

	req := &PacketRequest{
		Method: 12,
		ID:     12345,
		Context: map[string]string{
			"foo": "bar",
		},
		Data: map[string]string{
			"hello": "world",
		},
	}
	b, err := Marshal(req)
	a.NoError(err)
	p, err := Load(b)
	a.NoError(err)
	err = Unmarshal(p, &data)
	a.NoError(err)
	err = UnmarshalContext(p, &ctx)
	a.NoError(err)
	a.Equal("world", data["hello"])
	a.Equal("bar", ctx["foo"])
	a.Equal(uint16(12345), ID(req))
	a.Equal(uint8(12), Method(p))
}

func TestPacketEvent(t *testing.T) {
	var data map[string]string
	var ctx map[string]string
	a := assert.New(t)

	req := &PacketEvent{
		Method: 12,
		Context: map[string]string{
			"foo": "bar",
		},
		Data: map[string]string{
			"hello": "world",
		},
	}
	b, err := Marshal(req)
	a.NoError(err)
	p, err := Load(b)
	a.NoError(err)
	err = Unmarshal(p, &data)
	a.NoError(err)
	err = UnmarshalContext(p, &ctx)
	a.NoError(err)
	a.Equal("world", data["hello"])
	a.Equal("bar", ctx["foo"])
	a.Equal(uint8(12), Method(p))
}

func TestPacketResponse(t *testing.T) {
	var data map[string]string
	var ctx map[string]string
	a := assert.New(t)

	req := &PacketResponse{
		StatusCode: StatusCodeOK,
		ID:         12345,
		Context: map[string]string{
			"foo": "bar",
		},
		Data: map[string]string{
			"hello": "world",
		},
	}
	b, err := Marshal(req)
	a.NoError(err)
	p, err := Load(b)
	a.NoError(err)
	err = Unmarshal(p, &data)
	a.NoError(err)
	err = UnmarshalContext(p, &ctx)
	a.NoError(err)
	a.Equal("world", data["hello"])
	a.Equal("bar", ctx["foo"])
	a.Equal(StatusCodeOK, Status(req))
	a.Equal(uint16(12345), ID(p))
}
