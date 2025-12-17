package err

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

var (
	httpToGrpcCode = map[int]codes.Code{
		http.StatusOK:                  codes.OK,
		http.StatusBadRequest:          codes.InvalidArgument,
		http.StatusNotFound:            codes.NotFound,
		http.StatusConflict:            codes.AlreadyExists,
		http.StatusForbidden:           codes.PermissionDenied,
		http.StatusUnauthorized:        codes.Unauthenticated,
		http.StatusServiceUnavailable:  codes.Unavailable,
		http.StatusGatewayTimeout:      codes.DeadlineExceeded,
		http.StatusInternalServerError: codes.Internal,
	}
	grpcToHttpCode = map[codes.Code]int{
		codes.OK:               http.StatusOK,
		codes.InvalidArgument:  http.StatusBadRequest,
		codes.NotFound:         http.StatusNotFound,
		codes.AlreadyExists:    http.StatusConflict,
		codes.PermissionDenied: http.StatusForbidden,
		codes.Unauthenticated:  http.StatusUnauthorized,
		codes.Unavailable:      http.StatusServiceUnavailable,
		codes.DeadlineExceeded: http.StatusGatewayTimeout,
		codes.Internal:         http.StatusInternalServerError,
		codes.Unknown:          http.StatusInternalServerError,
	}
)
