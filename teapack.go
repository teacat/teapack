package teapack

import (
	"errors"
)

var (
	// ErrUnknownType 表示封包格式錯誤。
	ErrUnknownType = errors.New("teapack: 未知的封包種類")
	// ErrNotLoaded 表示正在解析尚未編譯的封包。
	ErrNotLoaded = errors.New("teapack: `Unmarshal` 只能用在已經 `Load` 的封包")
	// ErrTooShort 表示封包的長度不正確而無法解析。
	ErrTooShort = errors.New("teapack: 不正確的封包長度。")
	// ErrIncorrectRange 表示封包裡要求的範圍遠超出封包本身的長度，可能資料不正確或被截斷了。
	ErrIncorrectRange = errors.New("teapack: 封包內的資料指示範圍已經超出資料長度。")
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

// concatCopyPreAllocate 會以更有效率的方式結合位元組切片。
// https://stackoverflow.com/questions/37884361/concat-multiple-slices-in-golang
func concatCopyPreAllocate(slices [][]byte) []byte {
	var totalLen int
	for _, s := range slices {
		totalLen += len(s)
	}
	tmp := make([]byte, totalLen)
	var i int
	for _, s := range slices {
		i += copy(tmp[i:], s)
	}
	return tmp
}

// Packet 定義了一個可以編譯的封包。
type Packet interface {
	marshal() (b []byte, err error)
	unmarshal(v interface{}) (err error)
	unmarshalContext(v interface{}) (err error)
	load(b []byte) (err error)
}

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

// Load 能夠解析封包，在映射或是讀取一個封包之前都必須先解析。
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

// Marshal 能夠將一個封包編譯成位元組陣列格式。
func Marshal(p Packet) (b []byte, err error) {
	return p.marshal()
}

// ID 能夠取得已經解析資料的封包工作編號。
func ID(p Packet) uint16 {
	if v, ok := p.(*PacketResponse); ok {
		return v.ID
	}
	if v, ok := p.(*PacketRequest); ok {
		return v.ID
	}
	return 0
}

// Method 能夠取得已經解析資料的封包目標函式代碼。
func Method(p Packet) uint8 {
	if v, ok := p.(*PacketRequest); ok {
		return v.Method
	}
	if v, ok := p.(*PacketEvent); ok {
		return v.Method
	}
	return 0
}

// Status 能夠取得已經解析資料的封包狀態代碼。
func Status(p Packet) (status StatusCode) {
	if v, ok := p.(*PacketResponse); ok {
		return v.StatusCode
	}
	return StatusCodeOK
}

// Unmarshal 能夠將已經解析資料的封包酬載映射到本地的資料。
func Unmarshal(p Packet, v interface{}) (err error) {
	return p.unmarshal(v)
}

// UnmarshalContext 能夠將已經解析資料的封包上下文映射到本地的資料。
func UnmarshalContext(p Packet, v interface{}) (err error) {
	return p.unmarshalContext(v)
}
