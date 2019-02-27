# Go Application for test in Kubernetes Container Services.

Root directory is for health check and returns json "{status:up}".
"/dir" directory returns its IP addresses and os inforation with json format.

# USAGE
build the app with `docker build` command, then run app with port 8080.
If you would like to use docker, run `docker run -p $PORT:8080`.

When you run curl against container with directory "dir", it returns ip address, hostname and OS type.
```
~ ❯❯❯ curl localhost:30001/dir
{"ipaddress":["127.0.0.1,172.17.0.3"],"hostname":"da15c4a2f920","os":"linux"}%
```
