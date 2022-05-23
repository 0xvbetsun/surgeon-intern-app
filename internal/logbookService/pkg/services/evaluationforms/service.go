package evaluationforms

type (
	Service struct {
		EvaluationForms IEvaluationForms
		Dops            IDops
		MiniCex         IMiniCex
	}
)

func NewService(evaluationForms IEvaluationForms, dops IDops, miniCex IMiniCex) Service {
	return Service{
		EvaluationForms: evaluationForms,
		Dops:            dops,
		MiniCex:         miniCex,
	}
}
