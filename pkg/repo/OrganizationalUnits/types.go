package OrganizationalUnits

import (
	"context"
	"database/sql"

	"github.com/friendsofgo/errors"
	"github.com/vbetsun/surgeon-intern-app/pkg/models"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/repoerrors"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type TypeName string

const CLINIC TypeName = "Clinic"
const HOSPITAL TypeName = "Hospital"
const CLINIC_DEPARTMENT TypeName = "ClinicDepartment"

func (n TypeName) String() string {
	return string(n)
}

func (n TypeName) Valid() bool {
	switch n {
	case HOSPITAL, CLINIC, CLINIC_DEPARTMENT:
		return true
	}
	return false
}

type (
	TypesRepo struct {
		db *sql.DB
	}
)

func NewTypesRepo(db *sql.DB) *TypesRepo {
	return &TypesRepo{db: db}
}

func (t *TypesRepo) Add(ctx context.Context, ouType *models.OrganizationalUnitType) (*models.OrganizationalUnitType, error) {
	if TypeName(ouType.Name).Valid() {
		err := ouType.Insert(ctx, t.db, boil.Infer())
		if err != nil {
			return nil, repoerrors.ErrorFromDbError(err)
		}
		return ouType, nil
	}
	return nil, errors.New("Wrong typename")
}

func (t *TypesRepo) GetByName(ctx context.Context, name TypeName) (*models.OrganizationalUnitType, error) {
	if !name.Valid() {
		return nil, repoerrors.NewNotFoundError(name)
	}
	rType, err := models.OrganizationalUnitTypes(models.OrganizationalUnitTypeWhere.Name.EQ(name.String())).One(ctx, t.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repoerrors.NewNotFoundError(name)
		}
		return nil, repoerrors.ErrorFromDbError(err)
	}
	return rType, nil
}
