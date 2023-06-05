# 请求打到http_server判断前缀有误
- 如果不加 `\`, 而`req.URL.Path`有斜杠开头, 会导致判断前缀不符合
```go
if !strings.HasPrefix(req.URL.Path, "/"+hs.basePath) {
    utils.Logger().Errorln("HttpServer serving unexpected path: " + req.URL.Path)
    http.Error(resp, "HttpServer serving unexpected path: "+req.URL.Path, http.StatusBadRequest)
    return
}
```