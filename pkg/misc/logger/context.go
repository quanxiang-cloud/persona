package logger

import (
	"context"
	"net/http"

	"git.internal.yunify.com/qxp/persona/pkg/misc/id2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	requestID     = "Request-Id"
	requestIDName = "requestID"

	roleName      = "Role"
	_userID       = "User-Id"
	_userName     = "User-Name"
	_departmentID = "Department-Id"
)

type CopyHeader map[string]string

// GINRequestID get request id from gin context
func GINRequestID(ctx *gin.Context) zap.Field {
	if ctx == nil {
		return zap.String(requestIDName, "")
	}
	id := ctx.Request.Header.Get(requestID)
	return zap.String(requestIDName, id)
}

// STDRequestID get request id from std context
func STDRequestID(ctx context.Context) zap.Field {
	if ctx == nil {
		return zap.String(requestIDName, "")
	}
	id, ok := ctx.Value(requestID).(string)
	if ok {
		return zap.String(requestIDName, id)
	}
	return zap.String(requestIDName, "")
}

// STDHeader get  User-ID ,Role,User-Name,Department-Id from std context
func STDHeader(ctx context.Context) CopyHeader {
	headersCTX := make(CopyHeader, 0)
	if ctx == nil {
		return headersCTX
	}
	var keys = []string{requestID, _userID, _userName, _departmentID, roleName} //copy header
	CopyCTX(ctx, headersCTX, keys...)
	return headersCTX
}

func CopyCTX(ctx context.Context, headers CopyHeader, keys ...string) {
	for _, headerKey := range keys {
		headerValue, ok := ctx.Value(headerKey).(string)
		if ok {
			headers[headerKey] = headerValue
		}
	}
}

// CTXTransfer transfer requestID from gin.context
// to context.Context
func CTXTransfer(ctx *gin.Context) context.Context {
	var id string
	var name interface{} = requestID
	id = ctx.Request.Header.Get(requestID)
	c := context.Background()
	c = context.WithValue(c, name, id)
	c = context.WithValue(c, _departmentID, ctx.Request.Header.Get(_departmentID))
	c = context.WithValue(c, _userName, ctx.Request.Header.Get(_userName))
	c = context.WithValue(c, _userID, ctx.Request.Header.Get(_userID))
	c = context.WithValue(c, roleName, ctx.Request.Header.Get(roleName))
	return c
}

// GenRequestID gen requestID
func GenRequestID(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.TODO()
	}
	var name interface{} = requestID
	return context.WithValue(ctx, name, id2.GenID())
}

// ReentryRequestID reentry requestID
func ReentryRequestID(ctx context.Context, id string) context.Context {
	if ctx == nil {
		ctx = context.TODO()
	}
	var name interface{} = requestID
	return context.WithValue(ctx, name, id)
}

// HeadAdd add requestID to header
func HeadAdd(header *http.Header, id string) {
	header.Add("Request-Id", id)
}
