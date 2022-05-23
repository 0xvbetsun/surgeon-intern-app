//go:generate mockery --name IGetter
package organisations

import (
	"context"
	"database/sql"

	"github.com/vbetsun/surgeon-intern-app/pkg/models"
)

type (
	IGetter interface {
		Get(organisationId string) (*models.OrganizationalUnit, error)
		All() ([]*models.OrganizationalUnit, error)
	}
)

type Getter struct {
	db  *sql.DB
	ctx context.Context
}

func NewGetter(db *sql.DB, ctx context.Context) IGetter {
	return &Getter{db: db, ctx: ctx}
}

func (g *Getter) Get(organisationId string) (*models.OrganizationalUnit, error) {
	organisation, err := models.FindOrganizationalUnit(g.ctx, g.db, organisationId)
	if err != nil {
		return nil, err
	}
	return organisation, nil
}

func (g *Getter) All() ([]*models.OrganizationalUnit, error) {
	rOrganisations, err := models.OrganizationalUnits().All(g.ctx, g.db)
	if err != nil {
		return nil, err
	}
	return rOrganisations, nil
}
