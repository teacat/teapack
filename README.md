# TeaPack [![GoDoc](https://godoc.org/github.com/teacat/teapack?status.svg)](https://godoc.org/github.com/teacat/teapack) [![Coverage Status](https://coveralls.io/repos/github/teacat/teapack/badge.svg?branch=master)](https://coveralls.io/github/teacat/teapack?branch=master) [![Build Status](https://travis-ci.com/teacat/teapack.svg?branch=master)](https://travis-ci.com/teacat/teapack) [![Go Report Card](https://goreportcard.com/badge/github.com/teacat/teapack)](https://goreportcard.com/report/github.com/teacat/teapack)

基於 MsgPack 以達成比 JSON 還要更有效率、更節省頻寬的位元組資料傳遞方式。

# 這是什麼？

TeaPack 是基於位元組資料與 [vmihailenco/msgpack](https://github.com/vmihailenco/msgpack)（亦為 MsgPack）所實作的一個資料封包格式。更多相關的實作方式請參閱 [rfc/packet.md](./rfc/packet.md) 規範文件。

# 效能比較

因為 TeaPack 的部份資料是單純的位元組標記，因此會比起 JSON 或 MsgPack 在編譯時還要節省部份的字串與轉譯處理時間。

```
goos: windows
goarch: amd64
pkg: github.com/my/repo
BenchmarkTeaPackUnmarshal-12        	 1213075	      1002 ns/op	     920 B/op	      16 allocs/op
BenchmarkTeaPackUnmarshalData-12    	 1786981	       673 ns/op	     690 B/op	      11 allocs/op
BenchmarkJSONUnmarshal-12           	  630126	      2022 ns/op	     992 B/op	      17 allocs/op
BenchmarkMsgpackUnmarshal-12        	  704853	      1778 ns/op	    1016 B/op	      18 allocs/op

BenchmarkTeaPackMarshal-12          	 1423786	       848 ns/op	     448 B/op	       9 allocs/op
BenchmarkJSONMarshal-12             	  921212	      1363 ns/op	     560 B/op	      13 allocs/op
BenchmarkMsgpackMarshal-12          	  986062	      1198 ns/op	     224 B/op	       5 allocs/op
PASS
ok  	github.com/my/repo	11.363s
```

# 索引

* [安裝方式](#安裝方式)
* [使用方式](#使用方式)
    * [編譯資料](#編譯資料)
        * [請求](#請求)
        * [回應](#回應)
        * [事件](#事件)
    * [解析資料](#解析資料)
        * [取得附屬欄位](#取得附屬欄位)
        * [窺探封包種類](#窺探封包種類)
* [狀態代號](#狀態代號)

# 安裝方式

打開終端機並且透過 `go get` 安裝此套件即可。

```bash
$ go get github.com/teacat/teapack
```

# 使用方式

## 編譯資料

### 請求

當你希望發送請求給某個伺服端時，就可以使用請求封包（`PacketRequest`）。這個封包帶有幾個重要的欄位：

* `Method`：一個基於 `uint8` 的函式編號，遠端與本機都應該先定義好這些編號。遠端接收到編號時應找尋對應的函式呼叫。
* `ID`：請求的工作編號，回應封包也必須帶有相同的編號用以對應後續的處理事件。
* `Context`：上下文資料，可用來擺放版本、時間…等，非主要資料。
* `Data`：主要的資料酬載。

```go
var (
	createUser uint8 = iota
	updateUser uint8
	deleteUser uint8
)

func main() {
	b, err := teapack.Marshal(&PacketRequest{
		Method: createUser,
		ID:     1,
		Context: map[string]int64{
			"timestamp": time.Now().Unix(),
		},
		Data: map[string]string{
			"hello": "world",
		},
	})
	if err != nil {
		panic(err)
	}
}
```

### 回應

使用 `PacketResponse` 回應一個請求封包。這個封包帶有幾個重要的欄位：

* `StatusCode`：TeaPack 定義的狀態代碼。
* `ID`：回應的工作編號，這個編號應該來自請求封包。
* `Context`：上下文資料，可用來擺放版本、時間…等，非主要資料。
* `Data`：主要的資料酬載。

```go
func main() {
	b, err := teapack.Marshal(&PacketRequest{
		StatusCode: teapack.StatusCodeOK,
		ID:         1,
		Context: map[string]int64{
			"timestamp": time.Now().Unix(),
		},
		Data: map[string]string{
			"hello": "world",
		},
	})
	if err != nil {
		panic(err)
	}
}
```

### 事件

如果要在沒有請求的狀況下，主動發送訊息給客戶端就可以用上 `PacketEvent`。這個封包帶有幾個重要的欄位：

* `Method`：欲呼叫的函式編號或是事件種類編號。
* `Context`：上下文資料，可用來擺放版本、時間…等，非主要資料。
* `Data`：主要的資料酬載。

```go
var (
	newMessage uint8 = iota
	newEmail   uint8
)

func main() {
	b, err := teapack.Marshal(&PacketEvent{
		Method: newMessage,
		Context: map[string]int64{
			"timestamp": time.Now().Unix(),
		},
		Data: map[string]string{
			"hello": "world",
		},
	})
	if err != nil {
		panic(err)
	}
}
```

## 解析資料

透過 `Load` 從位元組陣列中載入封包資料，並且再透過 `Unmarshal` 等函式來將資料酬載映射到本機變數。

```go
func main() {
	// 編譯資料成為位元組陣列。
	b, err := teapack.Marshal(&PacketEvent{})
	if err != nil {
		panic(err)
	}
	// 從位元組陣列載入 TeaPack 資料。
	p, err := teapack.Load(b)
	if err != nil {
		panic(err)
	}
	// 將 TeaPack 封包裡的資料酬載映射到本機的變數。
	var data map[string]interface{}
	err = teapack.Unmarshal(p, &data)
	if err != nil {
		panic(err)
	}
}
```

### 取得附屬欄位

透過 `ID`、`Method` 或 `Status` 等函式可以取得請求或回應封包裡的對應編號。前提是：封包必須要先透過 `Load` 進行解析載入。

```go
func main() {
	// 編譯資料成為位元組陣列。
	b, err := teapack.Marshal(&PacketEvent{})
	if err != nil {
		panic(err)
	}
	// 從位元組陣列載入 TeaPack 資料。
	p, err := teapack.Load(b)
	if err != nil {
		panic(err)
	}
	// 取得此封包的目標函式編號。
	fmt.Println(teapack.Method(p))
}
```

### 窺探封包種類

透過 `Type` 來窺探位元組陣列為何種 TeaPack 封包種類。

```go
func main() {
	// 編譯資料成為位元組陣列。
	b, err := teapack.Marshal(&PacketEvent{})
	if err != nil {
		panic(err)
	}
	// 透過資料標頭檢視這個位元組陣列是什麼封包。
	fmt.Println(teapack.Type(b)) // 輸出：3（即為 teapack.PacketTypeEvent）
}
```

# 狀態代號

成功

| 狀態碼 | 狀態字樣   | 說明                                                          |
| ------ | ---------- | ---------------------------------------------------------- |
| 1      | OK         | 完全沒有錯誤                                                 |
| 2      | PROCESSING | 已接收到請求且正在處理中，不會馬上完成                            |
| 3      | NO_CHANGES | 提出的請求沒有改變任何事情，例如：請求刪除已經被刪除的資料           |

錯誤

| 狀態碼 | 狀態字樣           | 說明                                                          |
| ------ | ------------------ | ---------------------------------------------------------- |
| 50     | ERROR              | 內部不可預期的錯誤                                            |
| 51     | FULL               | 已滿而無法接受該請求，例如：正在加入已滿的聊天室、朋友清單           |
| 52     | EXISTS             | 某個東西已經存在，例如：使用者名稱、電子郵件地址                   |
| 53     | INVALID            | 請求的格式不正確                                             |
| 54     | NOT_FOUND          | 找不到請求的資源                                             |
| 55     | NOT_AUTHORIZED     | 請求者必須登入後才能發送此請求                                 |
| 56     | NO_PERMISSION      | 請求者在登入後沒有權限發送此請求而被拒絕                         |
| 57     | UNIMPLEMENTED      | 此功能尚未被實作完成                                         |
| 58     | TOO_MANY_REQUESTS  | 請求者在短時間內有太多的請求，需要暫緩一會才能重新發送請求          |
| 59     | RESOURCE_EXHAUSTED | 請求者可用的請求額度已被耗盡                                   |
| 60     | BUSY               | 伺服器正在繁忙中而暫時無法處理此請求                            |
| 61     | DEAD               | 伺服器已經關閉，可能是正在維護或是遇到錯誤而長期停機              |