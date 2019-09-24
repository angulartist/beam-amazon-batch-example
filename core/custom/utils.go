package custom

import (
	"github.com/apache/beam/sdks/go/pkg/beam"
	"math/rand"
	"strings"
)

/*
	PairWithOne takes a PCollection<T>, adds the value 1 to each element
	and returns back a PCollection<KV<T, 1>>.
*/
func PairWithOne(s beam.Scope, col beam.PCollection) beam.PCollection {
	return beam.ParDo(s, pairWithOneFn, col)
}

func pairWithOneFn(elm beam.T) (beam.T, int) {
	return elm, 1
}

/*
	Reshuffle takes a PCollection<T> and shuffles the data to help increase parallelism.
	Reshuffle adds a temporary random key to each element, performs a
  	GroupByKey, and finally removes the temporary key.
*/
func Reshuffle(s beam.Scope, col beam.PCollection) beam.PCollection {
	s = s.Scope("Reshuffle")

	col = beam.ParDo(s, func(x beam.X) (uint32, beam.X) {
		return rand.Uint32(), x
	}, col)

	col = beam.GroupByKey(s, col)

	return beam.ParDo(s, func(key uint32, values func(*beam.X) bool, emit func(beam.X)) {
		var x beam.X

		for values(&x) {
			emit(x)
		}
	}, col)
}

func Join(elem ...string) string {
	if len(elem) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(elem[0])
	for i, e := range elem[1:] {
		if !strings.HasSuffix(elem[i], "/") {
			sb.WriteRune('/')
		}
		sb.WriteString(e)
	}
	return sb.String()
}
