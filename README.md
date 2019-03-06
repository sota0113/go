# Go Application for test in Kubernetes Container Services.

Root directory is for health check and returns json "{status:up}".
"/dir" directory returns its IP addresses and os inforation with json format.

# USAGE
build the app with `docker build` command, then run app with port 8080.
If you would like to use docker, run like `docker run -p $PORT:8080`.

When you run curl against container with directory "dir", it returns ip address, hostname and OS type.
```
~ ❯❯❯ curl localhost:30001/dir
{"ipaddress":["127.0.0.1,172.17.0.3"],"hostname":"da15c4a2f920","os":"linux"}%
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
