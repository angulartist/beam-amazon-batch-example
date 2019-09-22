package ptrans

import (
	"github.com/apache/beam/sdks/go/pkg/beam"
	"github.com/apache/beam/sdks/go/pkg/beam/transforms/stats"
	"github.com/dariubs/percent"
	"projects/amazon_reviews_beam/core/custom"
	"projects/amazon_reviews_beam/core/models"
	"projects/amazon_reviews_beam/utils"
)

func ExtractRecommendationRatio(s beam.Scope, col beam.PCollection) beam.PCollection {
	s = s.Scope("Extract Recommendation Ratio")

	extractedKey := beam.ParDo(s.Scope("Extract DoRecommend Key"),
		func(r models.Review, emit func(string)) {
			emit(r.DoRecommend)
		}, col)

	pairedWithOne := custom.PairWithOne(s, extractedKey)

	summed := stats.SumPerKey(s, pairedWithOne)

	droppedKey := beam.DropKey(s, pairedWithOne)

	counted := stats.Sum(s, droppedKey)

	mapped := beam.ParDo(s,
		func(k string, v int,
			sideCounted int,
			emit func(ratio models.RecommendRatio)) {

			p := percent.PercentOf(v, sideCounted)

			emit(models.RecommendRatio{
				DoRecommend: k,
				NumVotes:    v,
				Percent:     p,
			})

		}, summed,
		beam.SideInput{Input: counted})

	return beam.ParDo(s, func(in models.RecommendRatio, emit func(models.RecommendOutput)) {
		emit(models.RecommendOutput{
			Data:      in,
			GraphType: utils.Pie,
		})
	}, mapped)
}
