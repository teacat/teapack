# 應用程式介面（API）規範

protoc -I ../../services/caris-events/authentication/pb --go_out=plugins=grpc:../../services/caris-events/authentication/pb --proto_path=. authentication.proto

protoc -I ../../services/caris-events/authentication/pb --go_out=plugins=grpc:../../services/caris-events/authentication/pb --proto_path=. pb.proto

所有連線（檔案上傳、資料請求）皆為 Message Pack 格式並經由 WebSocket 傳輸至伺服器。

## 內容格式

### 方法請求

當要送出一個請求到伺服器時，其 WebSocket 內容格式如下。

-   `tid`：工作編號，用以對應本地後續異步的處理函式。
-   `fun`：欲呼叫的遠端方法函式名稱。
-   `dat`：主要資料內容。
-   `met`：中繼資料（如：請求發起元件名稱、客戶端版本號碼）。

```json
{
    "tid": 1, // Task ID
    "fun": "CreateComment", // Function (Method)
    "dat": {
        // Data
        "content": "Hello, world!"
    },
    "met": {
        // Metadata
        "component": "CommentInput"
    }
}
```

### 回應內容

伺服器透過 WebSocket 所回傳的標準回應格式內容如下。

-   `tid`：工作編號，用以對應客戶端的處理函式。
-   `cod`：回應代號。
-   `dat`：主要成功資料內容（與 `err` 欄位不會同時存在）。
-   `err`：錯誤資料（與 `dat` 欄位不會同時存在）。
-   `met`：中繼資料（如：伺服器版本號碼、負載狀態）。

```json
{
    "tid": 1, // Task ID
    "cod": 1, // Status Code
    "dat": {
        // Data
        "id": 19
    },
    "err": null, // Error Data
    "met": {
        // Metadata
        "component": "CommentInput"
    }
}
```

### 通知內容

伺服器經由 WebSocket 主動發送至客戶端的通知事件內容格式如下。

-   `evt`：事件名稱。
-   `dat`：資料內容。
-   `met`：中繼資料（如：伺服器版本號碼、負載狀態）。

```json
{
    "evt": "NewChatMessage", // Event
    "dat": {
        // Data
        "content": "Foo, bar!"
    },
    "met": {
        // Metadata
        "version": "1.0.0+stable"
    }
}
```

### 檔案上傳

由客戶端經過 WebSocket 上傳檔案的時候，是切分為區塊並且以二進制資料進行傳遞。

```json
{
    "tid": 1, // ID
    "fun": "UploadFile", // Function (Method)
    "fil": {
        // Payload
        "key": "22523e48-7769-4369-b8f5-89c58a7804ed", // Key
        "bin": [128, 12, 83, 251, 90, 146] // Binary
    },
    "met": {
        // Metadata
        "component": "CommentInput"
    }
}
```

#### 區塊二進制格式

每個區塊被切分為 1024 KB（也就是 1 MB），在最起初的區塊位元資料中，前面有兩個區塊用來存放「檔案名稱」與「副檔名」。

```bash
+-----------+-----------+-----------+------------+
| File Name | Extension | File Size |   Binary   |
+-----------+-----------+-----------+------------+
| 256 Bytes | 16  Bytes | 4   Bytes | 748  Bytes | = 1024 Bytes
+-----------+-----------+-----------+------------+
```

然而在初次區塊傳遞完成，檔案已在伺服器產生時，接下來的區塊皆以檔案位元內容為主。

```bash
+------------------------------------------------+
|                      Binary                    |
+------------------------------------------------+
|                    1024 Bytes                  |
+------------------------------------------------+
```

## 狀態代號

成功

| 狀態碼 | 狀態字樣   | 說明                                                       |
| ------ | ---------- | ---------------------------------------------------------- |
| 0      | OK         | 完全沒有錯誤                                               |
| 1      | YES        | 是                                                         |
| 2      | NO         | 否                                                         |
| 3      | PROCESSING | 已接收到請求且正在處理中，不會馬上完成                     |
| 4      | NO_CHANGES | 提出的請求沒有改變任何事情，例如：請求刪除已經被刪除的資料 |

錯誤

| 狀態碼 | 狀態字樣           | 說明                                                       |
| ------ | ------------------ | ---------------------------------------------------------- |
| 50     | ERROR              | 內部不可預期的錯誤                                         |
| 51     | FULL               | 已滿而無法接受該請求，例如：正在加入已滿的聊天室、朋友清單 |
| 52     | EXISTS             | 某個東西已經存在，例如：使用者名稱、電子郵件地址           |
| 53     | INVALID            | 請求的格式不正確                                           |
| 54     | NOT_FOUND          | 找不到請求的資源                                           |
| 55     | NOT_AUTHORIZED     | 請求者必須登入後才能發送此請求                             |
| 56     | NO_PERMISSION      | 請求者在登入後沒有權限發送此請求而被拒絕                   |
| 57     | UNIMPLEMENTED      | 此功能尚未被實作完成                                       |
| 58     | TOO_MANY_REQUESTS  | 請求者在短時間內有太多的請求，需要暫緩一會才能重新發送請求 |
| 59     | RESOURCE_EXHAUSTED | 請求者可用的請求額度已被耗盡                               |
| 60     | BUSY               | 伺服器正在繁忙中而暫時無法處理此請求                       |
| 61     | DEAD               | 伺服器已經關閉，可能是正在維護或是遇到錯誤而長期停機       |
