package orthopedicSurgeries

type (
	Repo struct {
		ActivitiesRepo IActivitiesRepo
		ReviewsRepo    IReviewsRepo
	}
)

func NewRepo(activitiesRepo IActivitiesRepo, reviewsRepo IReviewsRepo) *Repo {
	return &Repo{ActivitiesRepo: activitiesRepo, ReviewsRepo: reviewsRepo}
}
