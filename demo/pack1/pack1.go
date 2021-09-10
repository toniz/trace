
package pack1

import (
    "context"
    "errors"
    tracelog "github.com/toniz/otel"

    pack2 "github.com/toniz/otel/demo/pack2"
)

func CallLevel_1(ctx context.Context) string {
    ctxc, _ := tracelog.NewSpan(ctx, "pack1.CallLevel_1", OtelSpanKindProducer)
    result := "-> 1" + pack.CallLevel_2(ctxc)

    tracelog.SetSpanOK(ctxc, "Successssssss.")
    tracelog.AddSpanAttribute(ctxc, map[string]string{"user_id": "1000098"})
    tracelog.AddSpanAttribute(ctxc, map[string]string{"type": "10"})
    params := map[string]string{"a":"1", "b":"2", "c":"3"}
    tracelog.AddSpanEvent(ctxc, "UPDATE", params)

    tracelog.EndSpan(ctxc)
    return result
}

