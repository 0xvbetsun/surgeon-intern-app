package logbookService

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/friendsofgo/errors"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	commonModel "github.com/vbetsun/surgeon-intern-app/internal/graph/model"
	"github.com/vbetsun/surgeon-intern-app/internal/logbookService/graph/generated"
	"github.com/vbetsun/surgeon-intern-app/internal/pkg/commonDirectives"
	"github.com/vbetsun/surgeon-intern-app/internal/pkg/commonErrors"
	"github.com/vbetsun/surgeon-intern-app/internal/pkg/middleware"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/repoerrors"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type GraphQL struct {
	Server *gin.Engine
}

func NewQlConfig(rootResolver *Resolver, directives *AuthDirectives) generated.Config {
	c := generated.Config{Resolvers: rootResolver}
	c.Directives.HasAtLeastRole = directives.HasAtLeastRole
	c.Directives.HasOneOfRoles = directives.HasOneOfRoles
	c.Directives.Binding = commonDirectives.Binding
	return c
}

func NewGraphQL(serverPrefix string, qlConfig generated.Config) *GraphQL {
	// Setting up Gin
	s := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "https://ogbook-app-tfrat.ondigitalocean.app", "https://app.ogbook.se", "https://staging-app.ogbook.se"}
	config.AllowCredentials = true
	config.AllowHeaders = []string{"*"}
	s.Use(cors.New(config))
	r := s.Group(serverPrefix)
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

		if errors.As(e, &repoerrors.NotFoundError{}) {
			err.Extensions = map[string]interface{}{
				"error": commonModel.ErrorItemNotFound,
			}
		} else if errors.As(e, &repoerrors.QueryError{}) {
			err.Message = "An unknown server error occurred."
			err.Extensions = map[string]interface{}{
				"error": commonModel.ErrorUnknownServerError,
			}
		} else if errors.As(e, &commonErrors.InvalidInputError{}) {
			err.Extensions = map[string]interface{}{
				"error": commonModel.ErrorBadRequest,
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
	h := playground.Handler("GraphQL", "/app/api/v1/query")
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
