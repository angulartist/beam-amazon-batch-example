package ptrans

import (
	"github.com/apache/beam/sdks/go/pkg/beam"
	"github.com/apache/beam/sdks/go/pkg/beam/transforms/top"
	"projects/amazon_reviews_beam/core/models"
	"projects/amazon_reviews_beam/utils"
)

func lessTopTenHelpfulInterface(x, y models.Review) bool {
	return x.NumHelpful < y.NumHelpful
}

func ExtractTopMostHelpful(s beam.Scope, col beam.PCollection, n int) beam.PCollection {
	s = s.Scope("Extract Most Helpful Reviews")

	toTop := top.Largest(s, col, n, lessTopTenHelpfulInterface)

	return beam.ParDo(s, func(r []models.Review, emit func(models.ReviewOutput)) {
		emit(models.ReviewOutput{
			Data:      r,
			GraphType: utils.Bar,
		})
	}, toTop)
}
