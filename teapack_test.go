package teapack

import (
	"fmt"
	"testing"
)

func TestPacketRequest(t *testing.T) {
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
	if err != nil {
		panic(err)
	}
	fmt.Printf("Marshal: %+v\n", b)

	p, err := Load(b)
	if err != nil {
		panic(err)
	}

	var data map[string]string
	err = Unmarshal(p, &data)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Unmarshal: %+v\n", data)
}

func TestPacketEvent(t *testing.T) {

}
