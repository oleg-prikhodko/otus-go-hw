package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func loggingInterceptor(logger common.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		resp, err := handler(ctx, req)

		latency := time.Since(start).Milliseconds()
		statusCode := status.Code(err)

		logger.Info(fmt.Sprintf(
			"[%s] %s %s %dms",
			start.Format(time.RFC3339),
			info.FullMethod,
			statusCode,
			latency,
		))

		return resp, err
	}
}
