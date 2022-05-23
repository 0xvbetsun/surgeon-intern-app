//go:generate mockery --name IDeleter
package organisations

type (
	IDeleter interface{}
	Deleter  struct{}
)

func NewDeleter() IDeleter {
	return &Deleter{}
}
