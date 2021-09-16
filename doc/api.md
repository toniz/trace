## 函数说明：

### 1. 设置日志导出方案
#### 1.1 SetGrpcExport  使用GRPC协议把日志传到远程的GPRC服务器上, 如阿里云的SLS.
```
SetGrpcExport(ctx context.Context, filename, serviceName, version string) error
参数说明:
  * ctx 传递上下文,如来自http请求, 可使用request.Context()取得。
  * filename GPRC服务的配置文件文件名，如：sls-demo/server/trace_config.json
  * serviceName 使用trace的应用的应用名. 如: OrderService, PaymentService
  * version 这个应用版本号. 如: v1.3.10
```
#### 1.2 SetHttpExport 使用HTTP协议导出到远程的HTTP服务器上, 如阿里云的链路追踪服务.
```
SetHttpExport(ctx context.Context, serviceName, version string) error
参数说明:
  * serviceName 使用trace的应用的应用名. 具体可以参考: tracing-demo/server
  * version 这个应用版本号.
```
#### 1.3 SetDefaultExport 使用默认导出，既直接打印到stdout.
```
SetDefaultExport(ctx context.Context, serviceName, version string) error
参数说明:
  * serviceName 使用trace的应用的应用名. 具体可以参考: sls-demo/server
  * version 这个应用版本号.
```
### 2. 设置Span
#### 2.1 NewSpan 创建一个span  
```
NewSpan(ctx context.Context, name string, kind int) (context.Context, error)
参数说明:
  * ctx: 用于传递上下文关联信息
  * name: 方法名
  * kind: span的类别
     0: SpanKindUnspecified
     1: SpanKindInternal
     2: SpanKindServer
     3: SpanKindClient
     4: SpanKindProducer
     5: SpanKindConsumer 
```
#### 2.2 EndSpan 结束当前span
```
EndSpan(ctx context.Context) error
参数说明：
  * ctx: NewSpan返回的context.
```
#### 2.3 AddSpanAttribute 添加属性
```
AddSpanAttribute(ctx context.Context, params map[string]string) error
参数说明:
  * ctx: NewSpan返回的context.
  * params: 参数kv
说明: 一次可以设置多个参数，也可以多次调用.
```
#### 2.4 AddSpanEvent 添加事件
```
AddSpanEvent(ctx context.Context, event string, params map[string]string) error
参数说明:
  * ctx: NewSpan返回的context.
  * event: 事件名称，可自定义，如：update insert等。
  * params: 事件想搞参数kv
说明: 一次可以设置多个参数，也可以多次调用.
```
#### 2.5 SetSpanOK 设置状态为成功
```
SetSpanOK(ctx context.Context, message string) error
参数说明:
  * ctx: NewSpan返回的context.
  * message: 随意文本
```
#### 2.6 SetSpanError 设置状态为失败
```
SetSpanError(ctx context.Context, err error) error
参数说明:
  * ctx: NewSpan返回的context.
  * err: 报错的error.
```

### 3 HTTP调用
#### 3.1 NewHandler 给http请求加上hook.
```
NewHandler(handler http.Handler, name string) http.Handler
参数说明:
  * handler: http句柄
  * name: 这个http方法的名字,如：getuser
```
#### 3.2 HttpDo 原始的HTTP调用
```
HttpDo(ctx context.Context, req *http.Request, name string, options ...interface{}) (*http.Response, error)
参数说明:
    * ctx: 上下文ctx
    * req: 系统包net/http的request结构
    * name: 这个HTTP请求的名称.
    * options: 传递参数(预留), 目前zipkin的包内有处理。
    * 返回结果：系统包net/http的response结构
```
#### 3.3 HttpGet 简化的封装的HTTP Get调用。
```
HttpGet(ctx context.Context, url, name string) (string, error)
参数说明:
  * ctx: 上下文ctx
  * url: http链接
  * name: 这个HTTP请求的名称.
  * 返回结果：string 是response body的内容
```
#### 3.4 HttpPost 简化的封装的HTTP Post调用。
```
HttpPost(ctx context.Context, url, contentType, body, name string) (string, error)
参数说明:
  * ctx: 上下文ctx
  * url: http链接
  * contentType: http头部的Content-Type。
  * body: HTTP请求body.
  * name: 这个HTTP请求的名称.
  * 返回结果：string 是response body的内容
```

### 其它函数
#### IsWork 检查trace是否可用
```
IsWork() error
前面所有函数内部都会调用这个方法。所以不必重复调用。
```

#### 关闭trace
```
Close(ctx context.Context) error
SetGrpcExport或者SetDefaultExport之后可以加上defer func(){Close(ctx)}();
```
etc..
