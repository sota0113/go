# Go API Server Just Returns your JSON.

For Japanese, see [日本語](/README_jp.md).  
just returns json object as an API server.  
Root directory is for health check and returns json "{status:up}".  
As default, "/dir" directory returns its IP addresses and os inforation with json format.

# USAGE
build the app with `docker build` command and run exposing container port 8080 to any host port.  
Using docker, run like `docker run ${THIS IMAGE} -p $PORT:8080`.  

When you access to container with directory "/dir", it returns json as default including ip address, hostname and OS type of running container.  
The return object is changeable. To configure your retrun, see next chapter `CRUD Operation`.
```
## Let's say, the application is running on localhost and listen port 30001 of localhost.
~ ❯❯❯ curl -i localhost:30001/dir
HTTP/1.1 200 OK
Content-Type: Application/json
Date: Wed, 06 Mar 2019 08:49:14 GMT
Content-Length: 77

{"ipaddress":["127.0.0.1,172.17.0.3"],"hostname":"57d7191cc011","os":"linux"}
```

Accessing container with curl as above leads to output logs to `stdout` as below.
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

# CRUD Operation
Return object of this application is changeable.  
`text/plain` and `application/json` is the usable Content-Type  for now.  
The type of content is automatically distinguished.  

## Create
Same operation with `Upadte` for now. See chapter `Upadte` below.  
For now, creating multiple directory and making JSON return on each directory at the same time is not capable.

## Read
See first chapter `USAGE`. Execute GET request to the application.

## Update
To change return object, execute `PUT` operation. Here is an example.
```
~ ❯❯❯ curl -X PUT localhost:3001/dir  -H "Content-Type: x-www-form-urlencoded" -d '{"message": "Hello World."}'
CONTENT UPDATED. Contetnt type is Application/json.
```

Make sure your update is reflected correctly.
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
To Delete your content, execute Delete request.

```
~ ❯❯❯ curl -i -X DELETE http://localhost:30001/dir
HTTP/1.1 200 OK
Content-Type:
Date: Wed, 06 Mar 2019 08:49:34 GMT
Content-Length: 17

CONTENT DELETED.
```

Make sure your update is reflected correctly.
```
~ ❯❯❯ curl -i http://localhost:30001/dir
HTTP/1.1 200 OK
Content-Type: Application/json
Date: Wed, 06 Mar 2019 08:49:37 GMT
Content-Length: 77

{"ipaddress":["127.0.0.1,172.17.0.3"],"hostname":"57d7191cc011","os":"linux"}
```

default return object is not deleatable.
