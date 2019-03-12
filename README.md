# Go API Server with custom endpoints and objects.

For Japanese, see [日本語](/README_jp.md).  
Docker image [IMAGE](https://hub.docker.com/r/sota0113/jpi-server)  
just returns custom object as an API server.  

![image](https://github.com/sota0113/go/blob/images/image/goApiServer.png)  
Root directory is for health check and returns json "{status:up}".  
As default, "/dir" directory returns its IP addresses and os inforation with json format.

# USAGE
pull image `docker pull sota0113/jpi-server:go-1.11.5` or build the app with `docker build` command and run exposing container port 8080 to any host port.  
Using docker, run like `docker run ${THIS IMAGE} -p $PORT:8080`.  

When you access to container with directory "/dir", it returns json as default including ip address, hostname and OS type of running container.  
The return object is changeable. To configure your retrun, see next chapter `CRUD Operation`.
```
## Let's say, the application is running on localhost and listen port 3002 of localhost.
❯❯❯ curl -i -X GET http://localhost:3002/list
HTTP/1.1 200 OK
Content-Type: application/json
Date: Tue, 12 Mar 2019 12:22:48 GMT
Content-Length: 77

{"ipaddress":["127.0.0.1,172.17.0.2"],"hostname":"cb0a02e8984d","os":"linux"}
```
Note, request end with "/" does not work.
```
❯❯❯ curl -i -X GET http://localhost:3002/list/
HTTP/1.1 404 Not Found
Date: Tue, 12 Mar 2019 12:22:44 GMT
Content-Length: 0
```

Accessing container with curl as above leads to output logs to `stdout` as below.
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

# CRUD Operation
Return object of this application is changeable.  
`text/plain` and `application/json` is the usable Content-Type for now.  
The type of content is automatically distinguished.  

## GET
See first chapter `USAGE`. Execute GET request to the application.

If you operate GET request against unexist path, it would return 404.
```
❯❯❯ curl -i -X GET http://localhost:3002/list/api/v1/unExistContent
HTTP/1.1 404 Not Found
Date: Sun, 10 Mar 2019 13:35:36 GMT
Content-Length: 0
```

## POST
`POST` operation is not allowed. It would return 400.  
Same operation as `Upadte` for now. See chapter `Upadte` below. 

## PUT
To update return object, execute `PUT` request operation against a path `/list/api/v1/${CONTENTNAME}`. 
`${CONTENTNAME}` could be any string of content name. "/" is not allowed to include. It would return 409.  
If your content is new one and created successfully, it would return 201.  
If your content already exists, meaning update your content, it would return 204.  
Here is an example.

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

make sure ${CONTENTNAME} of your request path `/list/api/v1/${CONTENTNAME}` does not contain "/".  
If "/" is included in request path, it would return 409.
```
❯❯❯ curl -i -X PUT localhost:3002/list/api/v1/newContent/  -H "Content-Type: application/json" -d '{"message": "Hello World!"}'
HTTP/1.1 409 Conflict
Date: Sun, 10 Mar 2019 15:53:41 GMT
Content-Length: 0
```

## Delete
To Delete your content, execute Delete request.  
If your opeartion succeeded, it would return 204.

```
❯❯❯ curl -i -X DELETE http://localhost:3002/list/api/v1/newContent
HTTP/1.1 204 No Content
Content-Type:
Date: Tue, 12 Mar 2019 12:39:08 GMT
```

Once your content is deleted, the endpoint returns 404.
```
❯❯❯ curl -i -X GET http://localhost:3002/list/api/v1/newContent
HTTP/1.1 404 Not Found
Content-Type:
Date: Tue, 12 Mar 2019 12:39:23 GMT
Content-Length: 0
```

