package teapack

import (
	"encoding/json"
	"testing"

	"github.com/vmihailenco/msgpack"
)

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
