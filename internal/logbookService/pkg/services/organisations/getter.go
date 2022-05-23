//go:generate mockery --name IGetter
package organisations

import (
	authorization2 "github.com/vbetsun/surgeon-intern-app/internal/pkg/services/authorization"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/organisations"
)

type (
	IGetter interface {
		List(userID string) ([]*APIOrganisation, error)
		Get(organisationID int) (*APIOrganisation, error)
	}
)

type Getter struct {
	Repo                 organisations.IGetter
	authorizationService authorization2.IService
}

func NewGetter(repo organisations.IGetter, authorizationService authorization2.IService) IGetter {
	return &Getter{Repo: repo, authorizationService: authorizationService}
}

func (g Getter) List(userID string) ([]*APIOrganisation, error) {
	allOrganisations, err := g.Repo.All()
	if err != nil {
		return nil, err
	}
	organisations := make([]*APIOrganisation, 0)
	if g.authorizationService.UserIsSuperAdmin(userID) {
		for _, organisation := range allOrganisations {
			organisations = append(organisations, &APIOrganisation{
				ID:          organisation.ID,
				DisplayName: organisation.DisplayName,
			})
		}
	} else {
		for _, organisation := range allOrganisations {
			authorized, err := g.authorizationService.Authorize([]authorization2.CasbinRequest{{
				Subject: userID,
				Domain:  string(organisation.ID),
				Object:  "organisations",
				Action:  "read",
			}}, false)
			if err != nil {
				continue
			}
			if authorized {
				organisations = append(organisations, &APIOrganisation{
					ID:          organisation.ID,
					DisplayName: organisation.DisplayName,
				})
			}
		}
	}

	return organisations, nil
}

func (g Getter) Get(organisationID int) (*APIOrganisation, error) {
	panic("implement me")
}
