package thirdparty

import (
	"context"
	"github.com/spf13/viper"
	"gopkg.in/auth0.v5/management"
)

type (
	Auth0 struct {
		Management *management.Management
	}
)

func NewAuth0ManagementApi(viper *viper.Viper) (*Auth0, error) {
	m, err := management.New(viper.GetString("AUTH0_MANAGEMENT_DOMAIN"), management.WithContext(context.Background()), management.WithClientCredentials(viper.GetString("AUTH0_MANAGEMENT_CLIENT_ID"), viper.GetString("AUTH0_MANAGEMENT_CLIENT_SECRET")))

	if err != nil {
		return nil, err
	}
	return &Auth0{Management: m}, nil
}
