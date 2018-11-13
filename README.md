# ipgo
ip2region go语言实现
依赖项目 <https://github.com/lionsoul2014/ip2region>

## 使用方法

1. 获取

```shell
go get github.com/WillVi/ipgo
```

2. 使用

```go
ipgo.GetFile("ip2region.db地址，自行进入上放github位置下载")
search, _ := ipgo.BtreeSearch("xxx.xxx.xxx.xxx")
fmt.Println(search)
```



