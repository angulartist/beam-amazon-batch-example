package ptrans

import (
	"context"
	"github.com/apache/beam/sdks/go/pkg/beam"
	"github.com/apache/beam/sdks/go/pkg/beam/transforms/stats"
	"github.com/apache/beam/sdks/go/pkg/beam/transforms/top"
	"projects/amazon_reviews_beam/core/custom"
	"projects/amazon_reviews_beam/core/models"
	"projects/amazon_reviews_beam/utils"
)

type GraphTypeFn struct {
	Type utils.Graph `json:"graph_type"`
}

// The ProcessElement function will be called of each element of the input PCollection.
func (fn *GraphTypeFn) ProcessElement(
	ctx context.Context, in []models.UserReviews) models.UserReviewsOutput {

	return models.UserReviewsOutput{
		Data:      in,
		GraphType: fn.Type,
	}
}

func lessTopMostRatingsInterface(x, y models.UserReviews) bool {
	return x.NumRatings < y.NumRatings
}

func ExtractTopReviewers(s beam.Scope, col beam.PCollection, n int) beam.PCollection {
	s = s.Scope("Extract Top Reviewers")

	extractedKey := beam.ParDo(s.Scope("Extract Username Key"),
		func(r models.Review, emit func(string)) {
			emit(r.Username)
		}, col)

	pairedWithOne := custom.PairWithOne(s, extractedKey)

	summed := stats.SumPerKey(s, pairedWithOne)

	outerDoFnExample := func(k string, v int) models.UserReviews {
		return models.UserReviews{
			Username:   k,
			NumRatings: v,
		}
	}

	mapped := beam.ParDo(s, outerDoFnExample, summed)

	toTop := top.Largest(s, mapped, n, lessTopMostRatingsInterface)

	// Demo of how to use a stateful function:
	return beam.ParDo(s, &GraphTypeFn{Type: utils.Bar}, toTop)
}
