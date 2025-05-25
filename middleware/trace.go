// internal/pkg/interceptor/tracing.go
package interceptor

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TracingInterceptor(tracer opentracing.Tracer) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		span, ctx := opentracing.StartSpanFromContext(ctx, info.FullMethod)
		defer span.Finish()

		if userID, ok := ctx.Value("user_id").(string); ok {
			span.SetTag("user.id", userID)
		}

		return handler(ctx, req)
	}
}

// MetadataCarrier 实现 opentracing.TextMapWriter
type MetadataCarrier metadata.MD

func (mc MetadataCarrier) Set(key, val string) {
	mc[key] = append(mc[key], val)
}

func (mc MetadataCarrier) ForeachKey(handler func(key, val string) error) error {
	for k, vals := range mc {
		for _, v := range vals {
			if err := handler(k, v); err != nil {
				return err
			}
		}
	}
	return nil
}
