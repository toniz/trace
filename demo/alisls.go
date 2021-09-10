
package main

import (
    "net/http"
    "fmt"
    "context"
    "errors"

    redis "github.com/go-redis/redis/v8"
    "github.com/go-redis/redis/extra/redisotel/v8"

    tracelog "github.com/toniz/otel"

    // test trace in package
    pack1 "github.com/toniz/otel/demo/pack1"
    pack2 "github.com/toniz/otel/demo/pack2"
)

var rc *redis.Client

func main() {
    ctx := context.TODO()
    err := tracelog.SetGrpcExport(ctx, "trace_ali_bbthis.json", "OrderService", "v0.3.10")
    if err != nil {
        panic(err)
    }

    helloHandler := func(w http.ResponseWriter, req *http.Request) {
        ctx := req.Context()
        tracelog.AddSpanAttribute(ctx, map[string]string{"user_id": "10086"})
        c3 := call_level_3(ctx)
        w.Write([]byte(c3))
    }

    hello2Handler := func(w http.ResponseWriter, req *http.Request) {
        ctx := req.Context()
        c1 := pack1.CallLevel_1(ctx)
        c2 := pack2.CallLevel_2(ctx)
        w.Write([]byte(c1+c2))
    }

    rc = redis.NewClient(&redis.Options{
        Addr:     "127.0.0.1:6379",
        Password: "",
        DB:       0,
    })

    rc.AddHook(redisotel.NewTracingHook())

    _, redisErr := rc.Ping(ctx).Result()
    if redisErr != nil {
        fmt.Println(redisErr)
    }

    fmt.Println(tracelog.IsWork())
    http.Handle("/hello", tracelog.NewHandler(http.HandlerFunc(helloHandler), "Hello"))
    http.Handle("/hello2", tracelog.NewHandler(http.HandlerFunc(hello2Handler), "Hello2"))
    fmt.Println("Now listen port 8080, you can visit 127.0.0.1:8080/hello .")
    err = http.ListenAndServe(":8080", nil)
    if err != nil {
        panic(err)
    }

}

func call_level_3(ctx context.Context) string {
    ctxc, _ := tracelog.NewSpan(ctx, "call_level_3", OtelSpanKindProducer)
    result := "-> 5" + call_level_3_1(ctxc)

    tracelog.SetSpanOK(ctxc, "Successssssss.")
    tracelog.AddSpanAttribute(ctxc, map[string]string{"user_id": "1000098"})
    tracelog.AddSpanAttribute(ctxc, map[string]string{"type": "10"})
    params := map[string]string{"a":"1", "b":"2", "c":"3"}
    tracelog.AddSpanEvent(ctxc, "UPDATE", params)

    tracelog.EndSpan(ctxc)
    return result
}


func call_level_3_1(ctx context.Context) string {
    ctxc, _ := tracelog.NewSpan(ctx, "call_level_3_1", OtelSpanKindProducer)

    err := errors.New("MyTestError")
    tracelog.SetSpanError(ctxc, err)
    params := map[string]string{"a":"1", "b":"2", "c":"3"}
    tracelog.AddSpanEvent(ctxc, "UPDATE", params)
    tracelog.AddSpanAttribute(ctxc, map[string]string{"user_id": "1000098"})

    res := call_level_redis(ctxc)
    tracelog.EndSpan(ctxc)
    return "-> 6" + res
}

func call_level_redis(ctx context.Context) string {
    result, redisErr := rc.Get(ctx, "test").Result()
    if redisErr != nil {
        fmt.Println(redisErr)
    }

    return result
}



