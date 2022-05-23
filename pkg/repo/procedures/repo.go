//go:generate mockery --name IRepo
package procedures

import (
	"context"
	"database/sql"

	"github.com/vbetsun/surgeon-intern-app/pkg/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type (
	IRepo interface {
		Add(ctx context.Context, procedure *models.Procedure) (*models.Procedure, error)
		GetByID(ctx context.Context, procedureID string) (*models.Procedure, error)
		ListAll(ctx context.Context) ([]*models.Procedure, error)
	}
	Repo struct {
		db *sql.DB
	}
)

func NewRepo(db *sql.DB) IRepo {
	return &Repo{db: db}
}

func (r *Repo) Add(ctx context.Context, procedure *models.Procedure) (*models.Procedure, error) {
	err := procedure.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		return nil, err
	}

	return procedure, nil
}

func (r *Repo) GetByID(ctx context.Context, procedureID string) (*models.Procedure, error) {
	procedure, err := models.FindProcedure(ctx, r.db, procedureID)
	if err != nil {
		return nil, err
	}
	return procedure, nil
}

func (r *Repo) ListAll(ctx context.Context) ([]*models.Procedure, error) {
	activities, err := models.Procedures().All(ctx, r.db)
	if err != nil {
		return nil, err
	}
	return activities, nil
}
