
package pack2

import (
    "context"
    "errors"
    tracelog "github.com/toniz/otel"
)

func CallLevel_2(ctx context.Context) string {
    ctxc, _ := tracelog.NewSpan(ctx, "call_level_3_1", OtelSpanKindProducer)

    err := errors.New("MyTestError")
    tracelog.SetSpanError(ctxc, err)
    params := map[string]string{"a":"1", "b":"2", "c":"3"}
    tracelog.AddSpanEvent(ctxc, "UPDATE", params)
    tracelog.AddSpanAttribute(ctxc, map[string]string{"user_id": "1000098"})

    res := call_level_redis(ctxc)
    tracelog.EndSpan(ctxc)
    return "-> 2" + res
}


