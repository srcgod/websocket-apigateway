package grpcerr

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GrpcErrorToHTTP(err error) (int, gin.H) {
	st, ok := status.FromError(err)
	if !ok {
		return http.StatusInternalServerError, gin.H{"error": "internal server error"}
	}

	switch st.Code() {
	case codes.AlreadyExists:
		return http.StatusConflict, gin.H{"error": st.Message()}
	case codes.InvalidArgument:
		return http.StatusBadRequest, gin.H{"error": st.Message()}
	case codes.NotFound:
		return http.StatusNotFound, gin.H{"error": st.Message()}
	case codes.Unauthenticated:
		return http.StatusUnauthorized, gin.H{"error": st.Message()}
	case codes.PermissionDenied:
		return http.StatusForbidden, gin.H{"error": st.Message()}
	case codes.Unavailable:
		return http.StatusServiceUnavailable, gin.H{"error": "service unavailable, try later"}
	default:
		return http.StatusInternalServerError, gin.H{"error": "internal server error"}
	}
}
