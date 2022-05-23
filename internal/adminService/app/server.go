package adminService

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/friendsofgo/errors"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/vbetsun/surgeon-intern-app/internal/adminService/graph/generated"
	commonModel "github.com/vbetsun/surgeon-intern-app/internal/graph/model"
	"github.com/vbetsun/surgeon-intern-app/internal/pkg/commonDirectives"
	"github.com/vbetsun/surgeon-intern-app/internal/pkg/commonErrors"
	"github.com/vbetsun/surgeon-intern-app/internal/pkg/middleware"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/repoerrors"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"go.uber.org/zap"
)

type GraphQL struct {
	Server *gin.Engine
}

func NewQlConfig(rootResolver *Resolver) generated.Config {
	c := generated.Config{Resolvers: rootResolver}
	c.Directives.HasAtLeastRole = HasAtLeastRole
	c.Directives.HasOneOfRoles = HasOneOfRoles
	c.Directives.Binding = commonDirectives.Binding
	return c
}

func NewGraphQL(serverPrefix string, qlConfig generated.Config, restApi *RestApi) *GraphQL {
	// Setting up Gin
	s := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "https://portal.ogbook.se", "https://staging-portal.ogbook.se"}
	config.AllowCredentials = true
	config.AllowHeaders = []string{"*"}
	s.Use(cors.New(config))
	r := s.Group(serverPrefix)
	restApi.InitializeRouter(r.Group("api/v1/rest"))
	r.Use(middleware.GinContextToContextMiddleware())
	r.POST("/api/v1/query", middleware.CheckJWT(), graphqlHandler(qlConfig))
	r.GET("/api/v1", playgroundHandler())
	return &GraphQL{Server: s}
}

// Defining the Graphql handler
func graphqlHandler(config generated.Config) gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	h := handler.NewDefaultServer(generated.NewExecutableSchema(config))
	h.SetErrorPresenter(func(ctx context.Context, e error) *gqlerror.Error {
		err := graphql.DefaultErrorPresenter(ctx, e)
		zap.S().Error(err.Error())
		if errors.As(e, &repoerrors.NotFoundError{}) {
			err.Message = commonModel.ErrorItemNotFound.String()
			err.Extensions = map[string]interface{}{
				"error": commonModel.ErrorItemNotFound,
			}
		} else if errors.As(e, &repoerrors.QueryError{}) {
			err.Message = commonModel.ErrorUnknownServerError.String()
			err.Extensions = map[string]interface{}{
				"error": commonModel.ErrorUnknownServerError,
			}
		} else if errors.As(e, &commonErrors.InvalidInputError{}) {
			err.Message = commonModel.ErrorBadRequest.String()
			err.Extensions = map[string]interface{}{
				"error": commonModel.ErrorBadRequest,
			}
		} else if errors.As(e, &commonErrors.UnauthorizedError{}) {
			err.Message = commonModel.ErrorUnAuthorized.String()
			err.Extensions = map[string]interface{}{
				"error": commonModel.ErrorUnAuthorized,
			}
		}
		return err
	})
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/portal/api/v1/query")
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
