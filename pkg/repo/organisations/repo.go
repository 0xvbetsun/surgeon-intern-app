package organisations

type (
	Repo struct {
		IGetter
		IUpdater
		IAdder
		IDeleter
	}
)

func NewRepo(IGetter IGetter, IUpdater IUpdater, IAdder IAdder, IDeleter IDeleter) *Repo {
	return &Repo{IGetter: IGetter, IUpdater: IUpdater, IAdder: IAdder, IDeleter: IDeleter}
}
