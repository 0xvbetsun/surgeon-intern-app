package gin

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

const GinContextUserIDKey = "UID"
const GinContextKey = "GinContextKey"

func GinContextFromContext(ctx context.Context) (*gin.Context, error) {
	ginContext := ctx.Value(GinContextKey)
	if ginContext == nil {
		err := fmt.Errorf("could not retrieve gin.Context")
		return nil, err
	}

	gc, ok := ginContext.(*gin.Context)
	if !ok {
		err := fmt.Errorf("gin.Context has wrong type")
		return nil, err
	}
	return gc, nil
}

func FirebaseIdFromGinContext(ginContext *gin.Context) (string, error) {
	uid := ginContext.Value(GinContextUserIDKey)
	if uid == nil {
		return "", errors.New("no user firebase id found")
	}
	return uid.(string), nil
}
