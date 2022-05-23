package organisations

type (
	Service struct {
		IAdder
		IGetter
		IUpdater
		IDeleter
	}
)

func NewService(IAdder IAdder, IGetter IGetter, IUpdater IUpdater, IDeleter IDeleter) *Service {
	return &Service{
		IAdder:   IAdder,
		IGetter:  IGetter,
		IUpdater: IUpdater,
		IDeleter: IDeleter,
	}
}
