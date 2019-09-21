package udfs

import (
	"github.com/apache/beam/sdks/go/pkg/beam"
	"github.com/apache/beam/sdks/go/pkg/beam/core/graph/mtime"
	"projects/amazon_reviews_beam/core/models"
	"time"
)

type AddTimestampFn struct {}

func ToTimestampFn(r models.Review) models.Review {
	layout := "2006-01-02T03:04:05.000Z"

	if t, err := time.Parse(layout, r.Date); err != nil {
		r.Timestamp = 1475625600 // TODO: @Mo: Do a better date format handling.
	} else {
		r.Timestamp = t.Unix()
	}

	return r
}

func (f *AddTimestampFn) ProcessElement(r models.Review) (beam.EventTime, models.Review) {
	timestamp := mtime.Normalize(mtime.Time(int64(time.Duration(r.Timestamp))))

	return timestamp, r
}
