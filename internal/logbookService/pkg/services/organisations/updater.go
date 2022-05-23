//go:generate mockery --name IUpdater
package organisations

type (
	IUpdater interface {
	}
	Updater struct{}
)

func NewUpdater() IUpdater {
	return &Updater{}
}
