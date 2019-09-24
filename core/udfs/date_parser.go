package udfs

import (
	"github.com/apache/beam/sdks/go/pkg/beam"
	"github.com/apache/beam/sdks/go/pkg/beam/core/graph/mtime"
	"projects/amazon_reviews_beam/core/models"
	"projects/amazon_reviews_beam/utils"
	"time"
)

func ToTimestampFn(review models.Review, emit func(models.Review), errors func(string, string)) {
	timestamp, err := time.Parse(utils.DATE_LAYOUT, review.Date)

	/* Example of error handling with a ParDo :
	-> Errors are kicked to an other flow to be written in an error logs file
	or to a BigQuery Sink for later processing for instance... */
	if err != nil {
		/* In a real-world, we'd like to log the eventId
		or the event itself to correct it later */
		errors(review.Date, err.Error())

		return
	}

	review.Timestamp = timestamp.Unix()

	emit(review)
}

type AddTimestampFn struct{}

func (f *AddTimestampFn) ProcessElement(review models.Review) (beam.EventTime, models.Review) {
	timestamp := mtime.Normalize(mtime.Time(int64(time.Duration(review.Timestamp))))

	return timestamp, review
}
