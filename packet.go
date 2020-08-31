package teapack

import (
	"encoding/binary"

	"github.com/vmihailenco/msgpack"
)

// PacketRequest 是一個請求封包。
type PacketRequest struct {
	// Method 是欲呼叫的方法名稱。
	Method uint8
	// ID 是封包工作編號，回應時會以相同編號回傳。
	ID uint16
	// Context 是中繼資料。
	Context interface{}
	// Data 是主要資料。
	Data interface{}
}

// Marshal 能夠將封包編譯成位元組資料。
func (p *PacketRequest) marshal() (b []byte, err error) {
	ctx, err := msgpack.Marshal(p.Context)
	if err != nil {
		return []byte{}, err
	}
	data, err := msgpack.Marshal(p.Data)
	if err != nil {
		return []byte{}, err
	}

	typ := []byte{uint8(PacketTypeRequest)}
	id := make([]byte, 2)
	binary.LittleEndian.PutUint16(id, p.ID)
	ctxLen := make([]byte, 2)
	binary.LittleEndian.PutUint16(ctxLen, uint16(len(ctx)))
	method := []byte{p.Method}

	b = concatCopyPreAllocate([][]byte{
		typ, id, ctxLen, method, ctx, data,
	})

	return b, nil
}

// load 會從位元組資料中解析資料並且轉換成一個封包。
func (p *PacketRequest) load(b []byte) (err error) {
	id := b[1:3]
	ctxLen := b[3:5]
	method := b[5:6]
	ctx := b[6 : 6+binary.LittleEndian.Uint16(ctxLen)]
	data := b[6+binary.LittleEndian.Uint16(ctxLen):]

	p.ID = binary.LittleEndian.Uint16(id)
	p.Method = uint8(method[0])
	p.Context = ctx
	p.Data = data
	return nil
}

// unmarshal 會將此解析過的封包酬載透過 Msgpack 映射到本地的資料。
func (p *PacketRequest) unmarshal(v interface{}) (err error) {
	if data, ok := p.Data.([]byte); ok {
		return msgpack.Unmarshal(data, v)
	}
	return ErrNotLoaded
}

// unmarshal 會將此解析過的封包上下文透過 Msgpack 映射到本地的資料。
func (p *PacketRequest) unmarshalContext(v interface{}) (err error) {
	if data, ok := p.Context.([]byte); ok {
		return msgpack.Unmarshal(data, v)
	}
	return ErrNotLoaded
}

// PacketResponse 是一個回應封包。
type PacketResponse struct {
	// ID 是封包工作編號。
	ID uint16
	// StatusCode 是狀態代碼。
	StatusCode StatusCode
	// Context 是中繼資料。
	Context interface{}
	// Data 是主要資料。
	Data interface{}
}

// Marshal 能夠將封包編譯成位元組資料。
func (p *PacketResponse) marshal() (b []byte, err error) {
	ctx, err := msgpack.Marshal(p.Context)
	if err != nil {
		return []byte{}, err
	}
	data, err := msgpack.Marshal(p.Data)
	if err != nil {
		return []byte{}, err
	}

	typ := []byte{uint8(PacketTypeResponse)}
	id := make([]byte, 2)
	binary.LittleEndian.PutUint16(id, p.ID)
	ctxLen := make([]byte, 2)
	binary.LittleEndian.PutUint16(ctxLen, uint16(len(ctx)))
	status := []byte{uint8(p.StatusCode)}

	b = concatCopyPreAllocate([][]byte{
		typ, id, ctxLen, status, ctx, data,
	})

	return b, nil
}

// load 會從位元組資料中解析資料並且轉換成一個封包。
func (p *PacketResponse) load(b []byte) (err error) {
	id := b[1:3]
	ctxLen := b[3:5]
	status := b[5:6]
	ctx := b[6 : 6+binary.LittleEndian.Uint16(ctxLen)]
	data := b[6+binary.LittleEndian.Uint16(ctxLen):]

	p.ID = binary.LittleEndian.Uint16(id)
	p.StatusCode = StatusCode(status[0])
	p.Context = ctx
	p.Data = data
	return nil
}

// unmarshal 會將此解析過的封包酬載透過 Msgpack 映射到本地的資料。
func (p *PacketResponse) unmarshal(v interface{}) (err error) {
	if data, ok := p.Data.([]byte); ok {
		return msgpack.Unmarshal(data, v)
	}
	return ErrNotLoaded
}

// unmarshal 會將此解析過的封包上下文透過 Msgpack 映射到本地的資料。
func (p *PacketResponse) unmarshalContext(v interface{}) (err error) {
	if data, ok := p.Context.([]byte); ok {
		return msgpack.Unmarshal(data, v)
	}
	return ErrNotLoaded
}

// PacketEvent 是一個事件封包。
type PacketEvent struct {
	// Method 是欲呼叫的方法名稱。
	Method uint8
	// Context 是中繼資料。
	Context interface{}
	// Data 是主要資料。
	Data interface{}
}

// Marshal 能夠將封包編譯成位元組資料。
func (p *PacketEvent) marshal() (b []byte, err error) {
	ctx, err := msgpack.Marshal(p.Context)
	if err != nil {
		return []byte{}, err
	}
	data, err := msgpack.Marshal(p.Data)
	if err != nil {
		return []byte{}, err
	}

	typ := []byte{uint8(PacketTypeEvent)}
	method := []byte{p.Method}
	ctxLen := make([]byte, 2)
	binary.LittleEndian.PutUint16(ctxLen, uint16(len(ctx)))

	b = concatCopyPreAllocate([][]byte{
		typ, method, ctxLen, ctx, data,
	})

	return b, nil
}

// load 會從位元組資料中解析資料並且轉換成一個封包。
func (p *PacketEvent) load(b []byte) (err error) {
	method := b[1:2]
	ctxLen := b[2:4]
	ctx := b[4 : 4+binary.LittleEndian.Uint16(ctxLen)]
	data := b[4+binary.LittleEndian.Uint16(ctxLen):]

	//
	p.Method = uint8(method[0])
	p.Context = ctx
	p.Data = data
	return nil
}

// unmarshal 會將此解析過的封包酬載透過 Msgpack 映射到本地的資料。
func (p *PacketEvent) unmarshal(v interface{}) (err error) {
	if data, ok := p.Data.([]byte); ok {
		return msgpack.Unmarshal(data, v)
	}
	return ErrNotLoaded
}

// unmarshal 會將此解析過的封包上下文透過 Msgpack 映射到本地的資料。
func (p *PacketEvent) unmarshalContext(v interface{}) (err error) {
	if data, ok := p.Context.([]byte); ok {
		return msgpack.Unmarshal(data, v)
	}
	return ErrNotLoaded
}
