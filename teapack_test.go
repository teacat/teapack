package teapack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vmihailenco/msgpack"
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
}

func BenchmarkTeaPackUnmarshal(b *testing.B) {
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
	bin, err := Marshal(req)
	if err != nil {
		panic(err)
	}

	for i := 0; i < b.N; i++ {
		var data map[string]string
		p, err := Load(bin)
		if err != nil {
			panic(err)
		}
		Method(p)
		if err != nil {
			panic(err)
		}
		ID(p)
		if err != nil {
			panic(err)
		}
		err = UnmarshalContext(p, &data)
		if err != nil {
			panic(err)
		}
		err = Unmarshal(p, &data)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkTeaPackUnmarshalData(b *testing.B) {
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
	bin, err := Marshal(req)
	if err != nil {
		panic(err)
	}

	for i := 0; i < b.N; i++ {
		var data map[string]string
		p, err := Load(bin)
		if err != nil {
			panic(err)
		}
		err = Unmarshal(p, &data)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkJSONUnmarshal(b *testing.B) {
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
	bin, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}

	for i := 0; i < b.N; i++ {
		var data PacketRequest
		err = json.Unmarshal(bin, &data)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkMsgpackUnmarshal(b *testing.B) {
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
	bin, err := msgpack.Marshal(req)
	if err != nil {
		panic(err)
	}

	for i := 0; i < b.N; i++ {
		var data PacketRequest
		err = msgpack.Unmarshal(bin, &data)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkTeaPackMarshal(b *testing.B) {
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
	for i := 0; i < b.N; i++ {
		_, err := Marshal(req)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkJSONMarshal(b *testing.B) {
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
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(req)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkMsgpackMarshal(b *testing.B) {
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
	for i := 0; i < b.N; i++ {
		_, err := msgpack.Marshal(req)
		if err != nil {
			panic(err)
		}
	}
}
