package orthopedicSurgeries

type (
	Service struct {
		ActivityService IActivities
		ReviewService   IReviews
		SurgeryService  ISurgeries
	}
)

func NewService(IActivities IActivities, IReviews IReviews, ISurgeries ISurgeries) *Service {
	return &Service{ActivityService: IActivities, ReviewService: IReviews, SurgeryService: ISurgeries}
}
