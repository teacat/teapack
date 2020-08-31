package teapack

import (
	"encoding/binary"
	"errors"

	"github.com/vmihailenco/msgpack"
)

const (
	// PacketTypeUnknown 是未知的種類。
	PacketTypeUnknown PacketType = iota
	// PacketTypeRequest 是請求封包。
	PacketTypeRequest
	// PacketTypeResponse 是回應封包。
	PacketTypeResponse
	// PacketTypeEvent 是事件封包。
	PacketTypeEvent
)

// PacketType 是封包種類。
type PacketType uint8

// StatusCode 是回應的狀態代碼。
type StatusCode uint8

const (
	// StatusCodeOK 表示完全沒有錯誤。
	StatusCodeOK StatusCode = 1
	// StatusCodeProcessing 表示已接收到請求且正在處理中，不會馬上完成。
	StatusCodeProcessing StatusCode = 2
	// StatusCodeNoChanges 表示提出的請求沒有改變任何事情，例如：請求刪除已經被刪除的資料。
	StatusCodeNoChanges StatusCode = 3

	// StatusCodeError 表示內部不可預期的錯誤。
	StatusCodeError StatusCode = 50
	// StatusCodeFull 表示已滿而無法接受該請求，例如：正在加入已滿的聊天室、朋友清單。
	StatusCodeFull StatusCode = 51
	// StatusCodeExists 表示某個東西已經存在，例如：使用者名稱、電子郵件地址。
	StatusCodeExists StatusCode = 52
	// StatusCodeInvalid 表示請求的格式不正確。
	StatusCodeInvalid StatusCode = 53
	// StatusCodeNotFound 表示找不到請求的資源。
	StatusCodeNotFound StatusCode = 54
	// StatusCodeNotAuthorized 表示請求者必須登入後才能發送此請求。
	StatusCodeNotAuthorized StatusCode = 55
	// StatusCodeNoPermission 表示請求者在登入後沒有權限發送此請求而被拒絕。
	StatusCodeNoPermission StatusCode = 56
	// StatusCodeUnimplemented 表示此功能尚未被實作完成。
	StatusCodeUnimplemented StatusCode = 57
	// StatusCodeTooManyRequests 表示請求者在短時間內有太多的請求，需要暫緩一會才能重新發送請求。
	StatusCodeTooManyRequests StatusCode = 58
	// StatusCodeResourceExhausted 表示請求者可用的請求額度已被耗盡。
	StatusCodeResourceExhausted StatusCode = 59
	// StatusCodeBusy 表示伺服器正在繁忙中而暫時無法處理此請求。
	StatusCodeBusy StatusCode = 60
	// StatusCodeDead 表示伺服器已經關閉，可能是正在維護或是遇到錯誤而長期停機。
	StatusCodeDead StatusCode = 61
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

	binary []byte
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

	//
	typ := []byte{uint8(PacketTypeRequest)}
	b = append(b, typ...)
	//
	id := make([]byte, 2)
	binary.LittleEndian.PutUint16(id, p.ID)
	b = append(b, id...)
	//
	ctxLen := make([]byte, 2)
	binary.LittleEndian.PutUint16(ctxLen, uint16(len(ctx)))
	b = append(b, ctxLen...)
	//
	method := []byte{p.Method}
	b = append(b, method...)
	//
	b = append(b, ctx...)
	//
	b = append(b, data...)

	return b, nil
}

func (p *PacketRequest) load(b []byte) (err error) {
	p.binary = b
	//
	//b = b[1:]
	//
	id := b[1:3]
	//
	ctxLen := b[3:5]
	//
	method := b[5:6]
	//
	ctx := b[6 : 6+binary.LittleEndian.Uint16(ctxLen)]
	//
	data := b[6+binary.LittleEndian.Uint16(ctxLen):]
	//
	p.ID = binary.LittleEndian.Uint16(id)
	p.Method = uint8(method[0])
	p.Context = ctx
	p.Data = data
	return nil
}

func (p *PacketRequest) unmarshal(v interface{}) (err error) {
	if data, ok := p.Data.([]byte); ok {
		return msgpack.Unmarshal(data, v)
	}
	return ErrNotLoaded
}

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

	binary []byte
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

	//
	typ := []byte{uint8(PacketTypeResponse)}
	b = append(b, typ...)
	//
	id := make([]byte, 2)
	binary.LittleEndian.PutUint16(id, p.ID)
	b = append(b, id...)
	//
	ctxLen := make([]byte, 2)
	binary.LittleEndian.PutUint16(ctxLen, uint16(len(ctx)))
	b = append(b, ctxLen...)
	//
	status := []byte{uint8(p.StatusCode)}
	b = append(b, status...)
	//
	b = append(b, ctx...)
	//
	b = append(b, data...)

	return b, nil
}

func (p *PacketResponse) load(b []byte) (err error) {
	p.binary = b
	//
	id := b[1:3]
	//
	ctxLen := b[3:5]
	//
	status := b[5:6]
	//
	ctx := b[6 : 6+binary.LittleEndian.Uint16(ctxLen)]
	//
	data := b[6+binary.LittleEndian.Uint16(ctxLen):]
	//
	p.ID = binary.LittleEndian.Uint16(id)
	p.StatusCode = StatusCode(status[0])
	p.Context = ctx
	p.Data = data
	return nil
}

func (p *PacketResponse) unmarshal(v interface{}) (err error) {
	if data, ok := p.Data.([]byte); ok {
		return msgpack.Unmarshal(data, v)
	}
	return ErrNotLoaded
}

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

	binary []byte
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

	//
	typ := []byte{uint8(PacketTypeEvent)}
	b = append(b, typ...)
	//
	method := []byte{p.Method}
	b = append(b, method...)
	//
	ctxLen := make([]byte, 2)
	binary.LittleEndian.PutUint16(ctxLen, uint16(len(ctx)))
	b = append(b, ctxLen...)
	//
	b = append(b, ctx...)
	//
	b = append(b, data...)

	return b, nil
}

func (p *PacketEvent) load(b []byte) (err error) {
	p.binary = b

	//
	method := b[1:2]
	//
	ctxLen := b[2:4]
	//
	ctx := b[4 : 4+binary.LittleEndian.Uint16(ctxLen)]
	//
	data := b[4+binary.LittleEndian.Uint16(ctxLen):]

	//
	p.Method = uint8(method[0])
	p.Context = ctx
	p.Data = data
	return nil
}

func (p *PacketEvent) unmarshal(v interface{}) (err error) {
	if data, ok := p.Data.([]byte); ok {
		return msgpack.Unmarshal(data, v)
	}
	return ErrNotLoaded
}

func (p *PacketEvent) unmarshalContext(v interface{}) (err error) {
	if data, ok := p.Context.([]byte); ok {
		return msgpack.Unmarshal(data, v)
	}
	return ErrNotLoaded
}

// Packet 定義了一個可以編譯的封包。
type Packet interface {
	marshal() (b []byte, err error)
	unmarshal(v interface{}) (err error)
	unmarshalContext(v interface{}) (err error)
	load(b []byte) (err error)
}

var (
	// ErrUnknownType 表示這個封包格式錯誤。
	ErrUnknownType = errors.New("teapack: 未知的封包種類")
	// ErrNotLoaded 表示正在解析一個尚未編譯的封包。
	ErrNotLoaded = errors.New("teapack: `Unmarshal` 只能用在已經 `Load` 的封包")
)

// Type 能夠在解析之前刺探封包的種類為何。
func Type(data []byte) PacketType {
	if len(data) <= 0 {
		return PacketTypeUnknown
	}
	switch PacketType(data[0]) {
	case PacketTypeRequest:
		return PacketTypeRequest
	case PacketTypeResponse:
		return PacketTypeResponse
	case PacketTypeEvent:
		return PacketTypeEvent
	default:
		return PacketTypeUnknown
	}
}

// Load 能夠解析封包。
func Load(data []byte) (p Packet, err error) {
	switch Type(data) {
	case PacketTypeRequest:
		p = &PacketRequest{}
	case PacketTypeResponse:
		p = &PacketResponse{}
	case PacketTypeEvent:
		p = &PacketEvent{}
	default:
		return nil, ErrUnknownType
	}
	if err := p.load(data); err != nil {
		return nil, err
	}
	return p, nil
}

//
func Marshal(p Packet) (b []byte, err error) {
	return p.marshal()
}

//
func ID(p Packet) uint16 {
	if v, ok := p.(*PacketResponse); ok {
		return v.ID
	}
	if v, ok := p.(*PacketRequest); ok {
		return v.ID
	}
	return 0
}

//
func Method(p Packet) uint8 {
	if v, ok := p.(*PacketRequest); ok {
		return v.Method
	}
	if v, ok := p.(*PacketEvent); ok {
		return v.Method
	}
	return 0
}

//
func Status(p Packet) (status StatusCode) {
	if v, ok := p.(*PacketResponse); ok {
		return v.StatusCode
	}
	return StatusCodeOK
}

//
func Unmarshal(p Packet, v interface{}) (err error) {
	return p.unmarshal(v)
}

//
func UnmarshalContext(p Packet, v interface{}) (err error) {
	return p.unmarshalContext(v)
}
