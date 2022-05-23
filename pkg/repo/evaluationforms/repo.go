package evaluationforms

type (
	Repo struct {
		EvaluationForms IEvaluationFormsRepo
		Dops            IDopsRepo
		MiniCex         IMiniCexRepo
	}
)

func NewRepo(evaluationForms IEvaluationFormsRepo, dops IDopsRepo, miniCex IMiniCexRepo) Repo {
	return Repo{evaluationForms, dops, miniCex}
}
