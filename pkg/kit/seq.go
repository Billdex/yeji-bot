package kit

import (
	"context"
	"sync/atomic"
)

type contextKeyMsgSeq struct{}

func SetSeq(ctx context.Context, seq int64) context.Context {
	return context.WithValue(ctx, contextKeyMsgSeq{}, &seq)
}

func Seq(ctx context.Context) int64 {
	v := ctx.Value(contextKeyMsgSeq{})
	if v == nil {
		return 0
	}
	if seq, ok := v.(*int64); ok {
		atomic.AddInt64(seq, 1)
		return *seq
	}
	return 0
}
