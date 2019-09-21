package ptrans

import (
	"github.com/apache/beam/sdks/go/pkg/beam"
	"github.com/apache/beam/sdks/go/pkg/beam/transforms/stats"
	"projects/amazon_reviews_beam/core/models"
	"projects/amazon_reviews_beam/utils"
)

func ExtractRatingsOverview(s beam.Scope, col beam.PCollection) beam.PCollection {
	s = s.Scope("Extract Ratings Overview")

	allRatings := beam.ParDo(s.Scope("Extract Rating Key"),
		func(r models.Review) int {
			return r.Rating
		}, col)

	averageRating := stats.Mean(s, allRatings)

	minRating := stats.Min(s, allRatings)

	maxRating := stats.Max(s, allRatings)

	mapped := beam.ParDo(s,
		func(avg float64,
			sideMin int,
			sideMax int,
			emit func(rating models.ComputedRating)) {

			emit(models.ComputedRating{
				Mean: avg,
				Min:  sideMin,
				Max:  sideMax,
			})
		}, averageRating,
		beam.SideInput{Input: minRating},
		beam.SideInput{Input: maxRating})

	return beam.ParDo(s, func(in models.ComputedRating, emit func(models.ComputedRatingOutput)) {
		emit(models.ComputedRatingOutput{
			ComputedRating: models.ComputedRating{
				Mean: in.Mean,
				Min:  in.Min,
				Max:  in.Max,
			},
			GraphType: utils.None,
		})
	}, mapped)
}
