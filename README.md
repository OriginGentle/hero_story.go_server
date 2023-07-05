
README
====

## ab 命令进行压力测试

ab -n 请求次数 -c 并发数 网站地址

例如：`ab -n 10000 -c 8 http://baidu.com/` （注意：最后这个 / 不能没有，否则会报错）

----

## 安装 Google Protobuf Golang 插件

```
go get google.golang.org/protobuf
```

安装完成之后，会在 gopath/bin 目录里多出来一个 protoc-gen-go.exe 文件。
这个 protoc-gen-go.exe 咱们直接是用不到的，
但是间接的会被 protoc.exe 调用到……

通过 GameMsgProtocol.proto 生成 go 代码

```
protoc --go_out=. .\GameMsgProtocol.proto
```

