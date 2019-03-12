# JSONを返す Go API Server

For English, see [English](/README.md).  
Docker image [IMAGE](https://hub.docker.com/r/sota0113/jpi-server)  
JSONを返すだけのAPIサーバーです。ローカル環境用です。  
ルートディレクトリはヘルスチェック用で"{status:up}"を返します。  
デフォルトではディレクトリ"/list"はランタイムが動作するホストのホスト名と全てのIPアドレスとos情報を返します。  

# 使い方
dockerコンテナーでビルドして動作させることを想定しています。  
`docker pull sota0113/jpi-server:go-1.11.5`を実行してイメージを取得するか`docker build`コマンドでイメージを生成し、デフォルトのコンテナーポート8080番と任意のホストポートを紐づけてイメージを起動します。  

上述の通り、コンテナーの"/list"にGETリクエストを発行すると、ホスト名、ホストの全てのIPアドレス、OS情報を返します。  
コンテナーからのリターン内容はエンドポイントと共に作成/変更可能です。詳しくは`CRUD操作`の章を参照してください。  

```
##　コンテナーがローカルホストで動作しており、コンテナーポート8080番がホストポート3002番と紐づいて起動していると仮定します。
❯❯❯ curl -i -X GET http://localhost:3002/list
HTTP/1.1 200 OK
Content-Type: application/json
Date: Tue, 12 Mar 2019 12:22:48 GMT
Content-Length: 77

{"ipaddress":["127.0.0.1,172.17.0.2"],"hostname":"cb0a02e8984d","os":"linux"}
```
リクエストエンドポイントの末尾に"/"があると失敗します。


GETリクエストを発行すると、コンテナーの標準出力として以下のようなログが出力されます。  
```
[ debug ] 2019/02/27 04:51:51 return_OSInfo.go:140: Application is started.
[  info ] 2019/02/27 04:52:01 return_OSInfo.go:46: Received access.
[  info ] 2019/02/27 04:52:01 return_OSInfo.go:47: Protocol: HTTP/1.1
[  info ] 2019/02/27 04:52:01 return_OSInfo.go:48: Method: GET
[  info ] 2019/02/27 04:52:01 return_OSInfo.go:49: Header: map[User-Agent:[curl/7.54.0] Accept:[*/*]]
[  info ] 2019/02/27 04:52:01 return_OSInfo.go:50: Body: {}
[  info ] 2019/02/27 04:52:01 return_OSInfo.go:51: Host: localhost:3002
[  info ] 2019/02/27 04:52:01 return_OSInfo.go:52: Form: map[]
[  info ] 2019/02/27 04:52:01 return_OSInfo.go:53: PostForm: map[]
[  info ] 2019/02/27 04:52:01 return_OSInfo.go:54: RequestURI: /dir
[  info ] 2019/02/27 04:52:01 return_OSInfo.go:55: Response: <nil>
```


# CRUD操作
アプリケーションのリターン内容は変更可能です。  
現時点では、`text/plain` と `application/json` がContent-Typeとて利用可能です。  
その他のContent-Typeを指定すると409を返します。

## GET
`使い方`を参照してください。コンテナーに対してGETリクエストを送付します。  
存在しないエンドポイントに対してGETリクエストを発行すると404を返します。
```
❯❯❯ curl -i -X GET http://localhost:3002/list/api/v1/unExistContent
HTTP/1.1 404 Not Found
Date: Sun, 10 Mar 2019 13:35:36 GMT
Content-Length: 0
```

## POST
現時点ではPOSTをサポートしていません。リターン内容のカスタムについては`PUT`を参照してください。  
複数のエンドポイントを作成してエンドポイント別に異なるコンテンツを返却させるには`PUT`を参照してください。


## PUT
コンテナーのリターン内容を変更するにはコンテナーのパス`/list/api/v1/${CONTENTNAME}`にPUTリクエストを送付します。
`${CONTENTNAME}`には任意のコンテンツ名が入ります。
新規エンドポイント作成とコンテンツ保管が成功した場合、201を返します。
既存のエンドポイントのコンテンツの更新が成功した場合、204を返します。

以下に例を示します。
```
## create new content.
❯❯❯ curl -i -X PUT localhost:3002/list/api/v1/newContent -H "Content-Type: application/json" -d '{"message": "Hello World!"}'
HTTP/1.1 201 Created
Content-Type: application/json
Date: Tue, 12 Mar 2019 12:31:57 GMT
Content-Length: 52

CONTENT UPDATED. Contetnt type is application/json.

## update content.
❯❯❯ curl -i -X PUT localhost:3002/list/api/v1/newContent -H "Content-Type: application/json" -d '{"message": "May the force be with you."}'
HTTP/1.1 204 No Content
Content-Type: application/json
Date: Tue, 12 Mar 2019 12:36:10 GMT
```

パス`/list/api/v1/${CONTENTNAME}`の末尾に"/"が入っていないことを確認してください。
末尾に"/"が入っている場合409を返します。
```
❯❯❯ curl -i -X PUT localhost:3002/list/api/v1/newContent/  -H "Content-Type: application/json" -d '{"message": "Hello World!"}'
HTTP/1.1 409 Conflict
Date: Sun, 10 Mar 2019 15:53:41 GMT
Content-Length: 0
```

## Delete

設定したUpdate/Createで設定した内容を削除するためには、Deleteリクエストを送付します。
リクエストが成功すると204を返します。

```
❯❯❯ curl -i -X DELETE http://localhost:3002/list/api/v1/newContent
HTTP/1.1 204 No Content
Content-Type:
Date: Tue, 12 Mar 2019 12:39:08 GMT
```

削除したエンドポイントに再度リクエストを送付すると404が返ります。
```
❯❯❯ curl -i -X GET http://localhost:3002/list/api/v1/newContent
HTTP/1.1 404 Not Found
Content-Type:
Date: Tue, 12 Mar 2019 12:39:23 GMT
Content-Length: 0
```
