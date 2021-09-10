
package main

import (
    "fmt"
    "context"
    "time"

    //"go.opentelemetry.io/otel/trace"
    . "github.com/bbthis/frame"
    . "github.com/bbthis/codes"
    . "github.com/bbthis/gosuit/cipher"
    tracelog "github.com/bbthis/gosuit/trace"
    "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
    err := tracelog.SetGrpcExport(context.TODO(), "trace_ali_bbthis.json.enc", KeyConfuse(FRAME_ConfigKey, "secret", 32), "OrderService", "v0.3.20")
    if err != OK {
        panic(err)
    }

    res, err := otelhttp.Head(context.TODO(), "http://127.0.0.1:8080/hello")

    fmt.Println(res)
    fmt.Println(err)
    res.Body.Close()

    time.Sleep(30*time.Second)
}

