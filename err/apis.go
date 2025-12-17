package err

import (
	"errors"
	"net/http"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

func NewBadRequest(code, message string) Error {
	return &impl{
		httpCode:     http.StatusBadRequest,
		ErrorCode:    code,
		ErrorMessage: message,
	}
}

func NewNotAcceptable(message string) Error {
	return &impl{
		httpCode:     http.StatusNotAcceptable,
		ErrorCode:    "NOT_ACCEPTABLE",
		ErrorMessage: message,
	}
}

func NewUnsupportedMediaType(message string) Error {
	return &impl{
		httpCode:     http.StatusUnsupportedMediaType,
		ErrorCode:    "UNSUPPORTED_MEDIA_TYPE",
		ErrorMessage: message,
	}
}

func NewUnauthorized(message string) Error {
	return &impl{
		httpCode:     http.StatusUnauthorized,
		ErrorCode:    "UNAUTHORIZED",
		ErrorMessage: message,
	}
}

func NewForbidden(message string) Error {
	return &impl{
		httpCode:     http.StatusForbidden,
		ErrorCode:    "FORBIDDEN",
		ErrorMessage: message,
	}
}

func NewServerError(err error) Error {
	if st, ok := status.FromError(err); ok {
		errDetails := st.Details()
		if len(errDetails) > 0 {
			if errInfo, ok := errDetails[0].(*errdetails.ErrorInfo); ok {
				return &impl{
					err:          err,
					httpCode:     grpcToHttpCode[st.Code()],
					ErrorCode:    errInfo.Metadata["error_code"],
					ErrorMessage: errInfo.Metadata["error_message"],
				}
			}
			return &impl{
				err:          err,
				httpCode:     grpcToHttpCode[st.Code()],
				ErrorCode:    "INTERNAL_SERVER_ERROR",
				ErrorMessage: st.Message(),
			}
		}
		return &impl{
			err:          err,
			httpCode:     grpcToHttpCode[st.Code()],
			ErrorCode:    "INTERNAL_SERVER_ERROR",
			ErrorMessage: st.Message(),
		}
	}
	return &impl{
		err:          err,
		httpCode:     http.StatusInternalServerError,
		ErrorCode:    "INTERNAL_SERVER_ERROR",
		ErrorMessage: err.Error(),
	}
}

func NewUnauthorizedError(message string) Error {
	return &impl{
		httpCode:     http.StatusUnauthorized,
		ErrorCode:    "UNAUTHORIZED",
		ErrorMessage: message,
	}
}

func NewGRPCError(err error) error {
	var e Error
	if !errors.As(err, &e) {
		e = NewServerError(err)
	}
	st := status.Newf(httpToGrpcCode[e.HttpCode()], e.Message())
	errInfo := &errdetails.ErrorInfo{
		Reason: e.Error(),
		Metadata: map[string]string{
			"error_code":    e.Code(),
			"error_message": e.Message(),
		},
	}

	stErr, err := st.WithDetails(errInfo)
	if err != nil {
		return st.Err()
	}
	return stErr.Err()
}
