package ptrans

import (
	"fmt"
	"github.com/apache/beam/sdks/go/pkg/beam"
	"github.com/apache/beam/sdks/go/pkg/beam/io/textio"
	"projects/amazon_reviews_beam/core/custom"
	"projects/amazon_reviews_beam/utils"
)

func WriteErrorLog(s beam.Scope, filename string, errors ...beam.PCollection) {
	s = s.Scope(fmt.Sprintf("Write %q", filename))

	flattenErrors := beam.Flatten(s, errors...)

	flattenErrors = beam.ParDo(s, func(key, value string) string {
		return fmt.Sprintf("%s,%s", key, value)
	}, flattenErrors)

	textio.Write(s, custom.Join(*utils.Logging, filename), flattenErrors)
}
