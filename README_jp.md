# JSONを返す Go API Server

For English, see [English](/README.md).  
Docker image [IMAGE](https://hub.docker.com/r/sota0113/jpi-server)  
JSONを返すだけのAPIサーバーです。ローカル環境向けです。  
ルートディレクトリはヘルスチェック用で"{status:up}"を返します。  
デフォルトではディレクトリ"/dir"はランタイムが動作するホストのホスト名と全てのIPアドレスとos情報を返します。  

# 使い方
dockerコンテナーでビルドして動作させることを想定しています。  
`docker pull sota0113/jpi-server:go-1.11.5`を実行してイメージを取得するか`docker build`コマンドでイメージを生成し、コンテナーポート8080番と任意のホストポートを紐づけてイメージを起動します。  

上述の通り、コンテナーの"/dir"にアクセスすると、デフォルトでホスト名、ホストの全てのIPアドレス、OS情報を返します。  
コンテナーからのリターン内容は変更可能です。リターン内容を変更するには`CRUD操作`の章を参照してください。  

```
##　コンテナーがローカルホストで動作しており、コンテナーポート8080番がホストポート30001番と紐づいて起動していると仮定します。
~ ❯❯❯ curl -i localhost:30001/dir
HTTP/1.1 200 OK
Content-Type: Application/json
Date: Wed, 06 Mar 2019 08:49:14 GMT
Content-Length: 77

{"ipaddress":["127.0.0.1,172.17.0.3"],"hostname":"57d7191cc011","os":"linux"}
```

上記のcurlコマンドを実行すると、コンテナーの標準出力として以下のようなログが出力されます。  
```
[ debug ] 2019/02/27 04:51:51 return_OSInfo.go:140: Application is started.
[  info ] 2019/02/27 04:52:01 return_OSInfo.go:46: Received access.
[  info ] 2019/02/27 04:52:01 return_OSInfo.go:47: Protocol: HTTP/1.1
[  info ] 2019/02/27 04:52:01 return_OSInfo.go:48: Method: GET
[  info ] 2019/02/27 04:52:01 return_OSInfo.go:49: Header: map[User-Agent:[curl/7.54.0] Accept:[*/*]]
[  info ] 2019/02/27 04:52:01 return_OSInfo.go:50: Body: {}
[  info ] 2019/02/27 04:52:01 return_OSInfo.go:51: Host: localhost:30001
[  info ] 2019/02/27 04:52:01 return_OSInfo.go:52: Form: map[]
[  info ] 2019/02/27 04:52:01 return_OSInfo.go:53: PostForm: map[]
[  info ] 2019/02/27 04:52:01 return_OSInfo.go:54: RequestURI: /dir
[  info ] 2019/02/27 04:52:01 return_OSInfo.go:55: Response: <nil>
```


# CRUD操作
アプリケーションのリターン内容は変更可能です。  
現時点では、`text/plain` と `application/json` がContent-Typeとて利用可能です。

## Create
現時点では`Update`と同一操作です。`Update`を参照してください。  
現時点では、複数のディレクトリを作成して、ディレクトリ別に複数のJSONを返却させることはできません。

## Read
`使い方`を参照してください。コンテナーに対してGETリクエストを送付します。

## Update
コンテナーのリターン内容を変更するにはコンテナーにPUTリクエストを送付します。以下に例を示します。
```
~ ❯❯❯ curl -X PUT localhost:3001/dir  -H "Content-Type: x-www-form-urlencoded" -d '{"message": "Hello World."}'
CONTENT UPDATED. Contetnt type is Application/json.
```

更新が適用されたことを確認します。
```
~ ❯❯❯ curl -i http://localhost:30001/dir
HTTP/1.1 200 OK
Content-Type: Application/json
Date: Wed, 06 Mar 2019 08:49:27 GMT
Content-Length: 37

{"message": "Hello World."}
```

POST is not acceptable at this verstion.

## Delete

設定したUpdate/Createで設定した内容を削除するためには、Deleteリクエストを送付します。
```
~ ❯❯❯ curl -i -X DELETE http://localhost:30001/dir
HTTP/1.1 200 OK
Content-Type:
Date: Wed, 06 Mar 2019 08:49:34 GMT
Content-Length: 17

CONTENT DELETED.
```

更新が適用され、デフォルトのリターンが表示されることを確認します。
```
~ ❯❯❯ curl -i http://localhost:30001/dir
HTTP/1.1 200 OK
Content-Type: Application/json
Date: Wed, 06 Mar 2019 08:49:37 GMT
Content-Length: 77

{"ipaddress":["127.0.0.1,172.17.0.3"],"hostname":"57d7191cc011","os":"linux"}
```

デフォルトのリターン内容は削除されません。
