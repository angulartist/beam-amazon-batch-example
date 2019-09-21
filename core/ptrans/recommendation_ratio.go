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

	//beam.ParDo0(s, func(w beam.Window, e beam.EventTime, review models.Review) {
	//	log.Debugf(ctx, "window=%v eventTime=%v for %s", w, e, review.Date)
	//}, col)

	summed := stats.SumPerKey(s, pairedWithOne)

	droppedKey := beam.DropKey(s, pairedWithOne)

	counted := stats.Sum(s, droppedKey)

	// Demo: How to test your PTransforms?
	// Asserts that we get the sum value of 5000 (the total number of records in our .csv chunks files)
	//passert.Sum(s, counted, "total rows", 1, 5000)

	mapped := beam.ParDo(s,
		func(k string, v int,
			sideCounted func(*int) bool,
			emit func(ratio models.RecommendRatio)) {
			var total int

			if sideCounted(&total) {
				p := percent.PercentOf(v, total)

				emit(models.RecommendRatio{
					DoRecommend: k,
					NumVotes:    v,
					Percent:     p,
				})
			}

		}, summed,
		beam.SideInput{Input: counted})

	return beam.ParDo(s, func(in models.RecommendRatio, emit func(models.RecommendOutput)) {
		emit(models.RecommendOutput{
			Data:      in,
			GraphType: utils.Pie,
		})
	}, mapped)
}
