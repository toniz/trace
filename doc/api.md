### 函数说明：
1. 使用GRPC协议把日志传到远程的GPRC服务器上, 如阿里云的SLS.
```
SetGrpcExport(ctx context.Context, filename, serviceName, version string) error
参数说明:
  * ctx 传递上下文
  * filename GPRC服务的配置文件文件名，如：sls-demo/server/trace_config.json
  * serviceName 使用trace的应用的应用名. 如: OrderService, PaymentService
  * version 这个应用版本号. 如: v1.3.10
```
---  
2. 使用默认导出，既直接打印到stdout.
```
SetDefaultExport(ctx context.Context, serviceName, version string) error
参数说明:
  * serviceName 使用trace的应用的应用名. 具体可以参考: sls-demo/server
  * version 这个应用版本号.
```
---
3. 给http请求加上hook.
```
  NewHandler(handler http.Handler, name string) http.Handler
参数说明:
  * handler: http句柄
  * name: 这个http方法的名字,如：getuser
```
---
4. 创建一个span.
```
NewSpan(ctx context.Context, name string, kind int) (context.Context, error)
参数说明:
```
---

etc..

