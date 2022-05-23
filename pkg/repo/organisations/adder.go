//go:generate mockery --name IAdder
package organisations

import (
	"context"
	"database/sql"

	"github.com/vbetsun/surgeon-intern-app/pkg/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type (
	IAdder interface {
		Add(organisation *models.OrganizationalUnit) (*models.OrganizationalUnit, error)
	}
	Adder struct {
		db  *sql.DB
		ctx context.Context
	}
)

func NewAdder(db *sql.DB, ctx context.Context) IAdder {
	return &Adder{
		db:  db,
		ctx: ctx,
	}
}

func (a *Adder) Add(organisation *models.OrganizationalUnit) (*models.OrganizationalUnit, error) {
	err := organisation.Insert(a.ctx, a.db, boil.Infer())
	if err != nil {
		return nil, err
	}

	return organisation, nil
}
