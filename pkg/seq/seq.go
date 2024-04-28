package seq

import "context"

type contextKeyMsgSeq struct{}

func SetSeq(ctx context.Context, seq int) context.Context {
	return context.WithValue(ctx, contextKeyMsgSeq{}, &seq)
}

func Seq(ctx context.Context) int {
	v := ctx.Value(contextKeyMsgSeq{})
	if v == nil {
		return 0
	}
	if seq, ok := v.(*int); ok {
		*seq++
		return *seq
	}
	return 0
}
