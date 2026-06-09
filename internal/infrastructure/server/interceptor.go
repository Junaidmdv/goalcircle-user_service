package server

import (
	"context"
	"runtime/debug"

	"github.com/Junaidmdv/goalcircle-user_service/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RecoveryInterceptor(logger logger.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {

		defer func() {
			if r := recover(); r != nil {
				logger.Error(
					"panic recovered",
					"method", info.FullMethod,
					"panic", r,
					"stack", string(debug.Stack()),
				)

				err = status.Error(
					codes.Internal,
					"internal server error",
				)
			}
		}()

		return handler(ctx, req)
	}
}

// func LoggingInterceptor(logger logger.Logger) grpc.UnaryServerInterceptor {
// 	return func(
// 		ctx context.Context,
// 		req interface{},
// 		info *grpc.UnaryServerInfo,
// 		handler grpc.UnaryHandler,
// 	) (interface{}, error) {


// 		start := time.Now()

// 		resp, err := handler(ctx, req)

// 		logger.Info(
// 			"method=%s duration=%s error=%v",
// 			info.FullMethod,
// 			time.Since(start),
// 			err,
// 		)

// 		return resp, err
// 	}
// }
