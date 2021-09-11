# 封装多种trace协议
目的是让记录trace日志变得和使用glog一样简单高效。
* 支持OpenTelemetry协议。封装了OpenTelemetry-Go.
* 支持zipkin协议。

## 协议切换
* 默认使用：OpenTelemetry协议
* 如果需要切换协议，如切换到zipkin协议, 则在main文件头加上如下引用即可:
```
_ https://github.com/toniz/trace/zipkin
```

## 使用步骤(以sls-demo为例子):
1. **引入**：
```
improt tracelog "github.com/toniz/trace"
```
2. **初始化trace, 指定日志写到哪里**:
```
err := tracelog.SetGrpcExport(ctx, "trace_config.json", "OrderService", "v0.3.10")
```
3. **对http服务设置hook。这个http请求会被trace捕捉并上报**。
```
http.Handle("/hello", tracelog.NewHandler(http.HandlerFunc(helloHandler), "Hello"))
```
4. **在被调用的函数里面创建span,添加该函数关键属性或事件**:
```
func call_level_3(ctx context.Context) {
    // 继承函数调用的context, 创建子的span
    ctxc, _ := tracelog.NewSpan(ctx, "call_level_3", tracelog.OtelSpanKindProducer)

    // 添加关键属性:可选
    tracelog.AddSpanAttribute(ctxc, map[string]string{"user_id": "1000098"})

    // 添加关键事件:可选
    tracelog.AddSpanEvent(ctxc, "UPDATE", map[string]string{"a":"1", "b":"2", "c":"3"})

    // 设置span状态, 错误就调用 tracelog.SetSpanError(ctxc, err)
    tracelog.SetSpanOK(ctxc, "Successssssss.")
    // 标记该子span结束
    tracelog.EndSpan(ctxc)

    return
}
```

### 其它文档：
[api文档](doc/api.md)    
[阿里云SLS服务接入例子](https://github.com/toniz/SLS-Aliyun)    

 
etc..

