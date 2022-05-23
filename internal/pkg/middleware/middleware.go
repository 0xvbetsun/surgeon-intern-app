package middleware

import (
	"context"
	"fmt"

	"github.com/form3tech-oss/jwt-go"
	"github.com/friendsofgo/errors"
	"github.com/gin-gonic/gin"
	commonModel "github.com/vbetsun/surgeon-intern-app/internal/graph/model"
	gin3 "github.com/vbetsun/surgeon-intern-app/internal/pkg/gin"
	"github.com/vbetsun/surgeon-intern-app/internal/pkg/services/users"
)

const GinContextUserIDKey = "UID"
const AUTH0GinContextUserIDKey = "UID"
const AUTH0GinContextUserClaimsKey = "UserClaims"
const Auth0JWTMiddleWareUserKey = "user"
const GinContextUserClaimsKey = "UserClaims"
const RoleClaimNamespace = "http://schemas.microsoft.com/ws/2008/06/identity/claims/role"
const AUTH0UserUnauthorized = "AUTH0UserUnauthorized"

type requestBody struct {
	OperationName *string     `json:"operationName"`
	Query         *string     `json:"query"`
	Variables     interface{} `json:"variables"`
}

func GinContextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "GinContextKey", c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func UserExternalIdFromContext(ctx context.Context) (string, error) {
	gc, err := gin3.GinContextFromContext(ctx)
	if err != nil {
		return "", err
	}
	if uID, ok := gc.Get(AUTH0GinContextUserIDKey); ok != false {
		return uID.(string), nil
	}
	return "", errors.New("Unknown error")
}

func UserFromContext(ctx context.Context, service users.IService) (*commonModel.User, error) {
	gc, err := gin3.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if uID, ok := gc.Get(AUTH0GinContextUserIDKey); ok != false {
		user, err := service.GetByExternalID(ctx, fmt.Sprintf("%v", uID))
		if err != nil {
			return nil, err
		}
		return user, nil
	}
	return nil, errors.New("Unknown error")
}

func UserClaimsFromContext(ctx context.Context) (jwt.MapClaims, error) {
	gc, err := gin3.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if claims, ok := gc.Get(AUTH0GinContextUserClaimsKey); ok != false {
		return claims.(jwt.MapClaims), nil
	}
	return nil, errors.New("Unknown error")
}

func ClaimsViaContextContainsRole(ctx context.Context, role string) (bool, error) {
	claims, err := UserClaimsFromContext(ctx)
	if err != nil {
		return false, errors.Wrap(err, "unable to get claims from ctx")
	}
	if ClaimsContainsRole(claims, role) {
		return true, nil
	}
	return false, errors.New("Role not found in claims")
}

func ClaimsContainsRole(claims jwt.MapClaims, role string) bool {
	roles, ok := claims[RoleClaimNamespace].([]interface{})
	if ok {
		// this means aud is a slice of strings
		for _, v := range roles {
			return v == role
		}
	}
	return false
}

func CheckJWT() gin.HandlerFunc {
	return func(c *gin.Context) {

		jwtMid := *jwtMiddleware
		if err := jwtMid.CheckJWT(c.Writer, c.Request); err != nil {
			c.AbortWithStatus(401)
		} else {
			user := c.Request.Context().Value(Auth0JWTMiddleWareUserKey)
			c.Set(AUTH0GinContextUserClaimsKey, user.(*jwt.Token).Claims)
			c.Set(AUTH0GinContextUserIDKey, externalUserIDFromClaims(user.(*jwt.Token).Claims))
		}
		c.Next()
	}
}
