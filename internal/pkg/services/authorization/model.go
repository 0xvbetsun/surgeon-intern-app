package authorization

type (
	CasbinRequest struct {
		Subject string
		Domain  string
		Object  string
		Action  string
	}
)
