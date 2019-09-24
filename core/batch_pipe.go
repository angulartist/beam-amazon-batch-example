package core

import (
	"bramp.net/morebeam/csvio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/beam/sdks/go/pkg/beam"
	"github.com/apache/beam/sdks/go/pkg/beam/core/graph/window"
	"github.com/apache/beam/sdks/go/pkg/beam/io/textio"
	"github.com/apache/beam/sdks/go/pkg/beam/log"
	"projects/amazon_reviews_beam/core/models"
	"projects/amazon_reviews_beam/core/ptrans"
	"projects/amazon_reviews_beam/core/udfs"
	"projects/amazon_reviews_beam/utils"
	"reflect"
)

func DescribePipeline(ctx context.Context, s beam.Scope) {

	log.Infof(ctx, "Started pipeline on scope: %s", s)

	// [#START BATCH EXAMPLE]
	source := csvio.Read(s, *utils.Input, reflect.TypeOf(models.Review{}))

	/* Only (re)shuffle your data when it's needed. This could increase parallelism but will
	Introduce a huge performance hit because of too much network communication between nodes on the
	reduce phase (ie: GroupByKey) */
	//source = custom.Reshuffle(s, source)

	timestampFormatted, timestampAssignmentError := beam.ParDo2(s, udfs.ToTimestampFn, source)

	withTimestamp := beam.ParDo(s, &udfs.AddTimestampFn{}, timestampFormatted)

	// TODO: @Mo: Fixed window seems to cause troubles with side inputs.
	/* withFixedWindow := beam.WindowInto(s, window.NewFixedWindows(utils.HALF_DAY_IN_SECONDS), withTimestamp) */

	withGlobalWindow := beam.WindowInto(s, window.NewGlobalWindows(), withTimestamp)

	statsRatingsOverview := ptrans.ExtractRatingsOverview(s, withGlobalWindow)

	statsNHelpful := ptrans.ExtractTopMostHelpful(s, withGlobalWindow, 10)

	statsNMostReviews := ptrans.ExtractTopReviewers(s, withGlobalWindow, 5)

	statsRecommendRatio := ptrans.ExtractRecommendationRatio(s, withGlobalWindow)

	/* ParDo is the core element-wise PTransform in Apache Beam,
	invoking a user-specified function on each of the elements of the input PCollection to
	produce zero or more output elements, all of which are collected into the output PCollection.
	Use one of the ParDo variants for a different number of output PCollections.
	The PCollections do no need to have the same types.

	Elements are processed independently, and possibly in parallel across distributed cloud resources.
	The ParDo processing style is similar to what happens inside the "Mapper" or "Reducer" class of
	a MapReduce-style algorithm. */

	merged := beam.ParDo(s.Scope("Merge PCollections Together"),
		func(computedRating models.ComputedRatingOutput,
			sideRecommendRatio []models.RecommendOutput,
			sideNHelpful models.ReviewOutput,
			sideNMostReviews []models.UserReviewsOutput,
			emit func(models.OutputFormat)) {

			emit(models.OutputFormat{
				ReviewOutput:         sideNHelpful,
				UserReviewsOutput:    sideNMostReviews,
				ComputedRatingOutput: computedRating,
				RecommendOutput:      sideRecommendRatio,
			})

		}, statsRatingsOverview,
		beam.SideInput{Input: statsRecommendRatio},
		beam.SideInput{Input: statsNHelpful},
		beam.SideInput{Input: statsNMostReviews},
	)

	o := beam.ParDo(s.Scope("Stringify Output"),
		func(r models.OutputFormat) string {
			if data, err := json.MarshalIndent(r, "", ""); err != nil {
				return fmt.Sprintf("[Err]: %v", err)
			} else {
				return fmt.Sprintf("%s", data)
			}
		}, merged)

	/* I/O: Write writes a PCollection<string> to a file as separate lines.
	The writer add a newline after each element.

	PCollection<[]string> -> I/O Write. */
	textio.Write(s.Scope("Write to Text Sink"), *utils.Output, o)

	// TODO: @Mo: Add BiQuery output sink example.

	/* Write error log */
	ptrans.WriteErrorLog(s, "logs.log", timestampAssignmentError)

	// [#END BATCH EXAMPLE]
}

/*
	RegisterType inserts "external" types into a global type registry to
	bypass serialization and preserve full method information.
	It should be called in init() only.

	This is because in a distributed pipeline, the PTransforms could be processed on different workers,
	and the PCollection needs to be serialized to be passed between them...

	Note: /!\ Errors will only be triggered when using a runner different than the Direct runner.
*/
func init() {
	beam.RegisterType(reflect.TypeOf(models.Review{}))
	beam.RegisterType(reflect.TypeOf(models.UserReviews{}))
	beam.RegisterType(reflect.TypeOf(models.ComputedRating{}))
	beam.RegisterType(reflect.TypeOf(models.RecommendRatio{}))
	beam.RegisterType(reflect.TypeOf(models.RecommendOutput{}))
	beam.RegisterType(reflect.TypeOf(models.ComputedRatingOutput{}))
	beam.RegisterType(reflect.TypeOf(models.UserReviewsOutput{}))
	beam.RegisterType(reflect.TypeOf(models.ReviewOutput{}))
	beam.RegisterType(reflect.TypeOf(models.OutputFormat{}))
	beam.RegisterType(reflect.TypeOf((*udfs.AddTimestampFn)(nil)).Elem())
	beam.RegisterType(reflect.TypeOf((*ptrans.GraphTypeFn)(nil)).Elem())
}
