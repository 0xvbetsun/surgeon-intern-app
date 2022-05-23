package surgeries

import (
	"context"
	"database/sql"

	"github.com/vbetsun/surgeon-intern-app/pkg/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	. "github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type (
	IRepo interface {
		GetByID(ctx context.Context, surgeryID string) (*models.Surgery, error)
		GetByDiagnoseAndMethod(ctx context.Context, diagnoseId string, methodId string, withRels bool) (*models.Surgery, error)
		GetByIDWithRels(ctx context.Context, surgeryID string) (*models.Surgery, error)
		ListWithRels(ctx context.Context) ([]*models.Surgery, error)
		ListByClinicWithRels(ctx context.Context, clinicId string) ([]*models.Surgery, error)
		ListByDiagnoseWithRels(ctx context.Context, diagnoseId string) ([]*models.Surgery, error)
		ListByMethodWithRels(ctx context.Context, methodId string) ([]*models.Surgery, error)
		GetByOrthopedicSurgeriesActivity(ctx context.Context, surgeryActivityID string) ([]*models.Surgery, error)
		Add(ctx context.Context, method *models.SurgeryMethod, fracture *models.SurgeryDiagnosis, specialtyId string) (*models.Surgery, error)
	}
	Repo struct {
		Db *sql.DB
	}
)

func NewRepo(db *sql.DB) IRepo {
	return &Repo{Db: db}
}

func (r *Repo) GetByDiagnoseAndMethod(ctx context.Context, diagnoseId string, methodId string, withRels bool) (*models.Surgery, error) {
	querymods := []QueryMod{models.SurgeryWhere.DiagnoseID.EQ(diagnoseId), models.SurgeryWhere.MethodID.EQ(methodId)}
	if withRels {
		querymods = append(querymods, Load(models.SurgeryRels.Method))
		querymods = append(querymods, Load(models.SurgeryRels.Diagnose))
	}
	surgery, err := models.Surgeries(
		querymods...).One(ctx, r.Db)
	if err != nil {
		return nil, err
	}
	return surgery, err

}

func (r *Repo) ListByDiagnoseWithRels(ctx context.Context, diagnoseId string) ([]*models.Surgery, error) {
	surgeries, err := models.Surgeries(models.SurgeryWhere.DiagnoseID.EQ(diagnoseId), Load(models.SurgeryRels.Method), Load(models.SurgeryRels.Diagnose)).All(ctx, r.Db)
	if err != nil {
		return nil, err
	}
	return surgeries, err
}

func (r *Repo) ListByMethodWithRels(ctx context.Context, methodId string) ([]*models.Surgery, error) {
	surgeries, err := models.Surgeries(models.SurgeryWhere.MethodID.EQ(methodId), Load(models.SurgeryRels.Method), Load(models.SurgeryRels.Diagnose)).All(ctx, r.Db)
	if err != nil {
		return nil, err
	}
	return surgeries, err
}

func (r *Repo) ListWithRels(ctx context.Context) ([]*models.Surgery, error) {
	surgeries, err := models.Surgeries(Load(models.SurgeryRels.Method), Load(models.SurgeryRels.Diagnose)).All(ctx, r.Db)
	if err != nil {
		return nil, err
	}
	return surgeries, err
}

func (r *Repo) ListByClinicWithRels(ctx context.Context, clinicId string) ([]*models.Surgery, error) {
	// TODO: We have no surgeru-clinic connetion right now, but maybe in future. This method returns all surgeries right now.
	surgeries, err := models.Surgeries(Load(models.SurgeryRels.Method), Load(models.SurgeryRels.Diagnose)).All(ctx, r.Db)
	if err != nil {
		return nil, err
	}
	return surgeries, err
}

func (r *Repo) GetByID(ctx context.Context, surgeryID string) (*models.Surgery, error) {
	surgery, err := models.Surgeries(Load(models.SurgeryRels.Method), Load(models.SurgeryRels.Diagnose), models.SurgeryWhere.ID.EQ(surgeryID)).One(ctx, r.Db)
	if err != nil {
		return nil, err
	}
	return surgery, err
}

func (r *Repo) GetByIDWithRels(ctx context.Context, surgeryID string) (*models.Surgery, error) {
	surgery, err := models.Surgeries(models.SurgeryWhere.ID.EQ(surgeryID), Load(models.SurgeryRels.Method), Load(models.SurgeryRels.Diagnose)).One(ctx, r.Db)
	if err != nil {
		return nil, err
	}
	return surgery, err
}

func (r *Repo) Add(ctx context.Context, method *models.SurgeryMethod, diagnose *models.SurgeryDiagnosis, specialtyId string) (*models.Surgery, error) {
	tx, err := r.Db.BeginTx(ctx, nil)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if _, err := models.FindSurgeryDiagnosis(ctx, tx, diagnose.ID); err != nil {
		if err = diagnose.Insert(ctx, tx, boil.Infer()); err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	if _, err := models.FindSurgeryMethod(ctx, tx, method.ID); err != nil {
		if err = method.Insert(ctx, tx, boil.Infer()); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	surgery := &models.Surgery{
		DiagnoseID:       diagnose.ID,
		MethodID:         method.ID,
		SurgerySpecialty: specialtyId,
	}
	if err = surgery.Insert(ctx, tx, boil.Infer()); err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return surgery, nil
}

func (r *Repo) GetByOrthopedicSurgeriesActivity(ctx context.Context, surgeryActivityID string) ([]*models.Surgery, error) {
	rActivitiesSurgeries, err := models.OrthopedicSurgeryActivitiesSurgeries(
		models.OrthopedicSurgeryActivitiesSurgeryWhere.OrthopedicSurgeryActivityID.EQ(surgeryActivityID)).All(ctx, r.Db)
	if err != nil {
		return nil, err
	}
	surgeries := make([]*models.Surgery, 0)
	for _, rActivitiesSurgery := range rActivitiesSurgeries {
		surgery, err := r.GetByIDWithRels(ctx, rActivitiesSurgery.SurgeryID)
		if err != nil {
			return nil, err
		}
		surgeries = append(surgeries, surgery)
	}

	return surgeries, nil
}
