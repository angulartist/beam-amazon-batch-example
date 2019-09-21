package custom

import (
	"context"
	"fmt"
	"github.com/apache/beam/sdks/go/pkg/beam"
	"math/rand"
	"time"
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

/* Playground */

// TODO: @Mo: Add BatchSizeEstimator struct https://github.com/apache/beam/blob/511caa3d8025e6d51adc0caa2fa253222f325ec3/sdks/python/apache_beam/transforms/util.py#L231
type BatchSizeEstimator struct {
	minBatchSize            int
	maxBatchSize            int
	targetBatchOverhead     float64
	targetBatchDurationSecs int
	variance                float64
	clock                   time.Time
}

// TODO: @Mo: Implement GlobalWindowsBatching DoFunction https://github.com/apache/beam/blob/511caa3d8025e6d51adc0caa2fa253222f325ec3/sdks/python/apache_beam/transforms/util.py#L418
type GlobalWindowsBatchingDoFn struct {
	batchSizeEstimator BatchSizeEstimator
	//batch              []interface{}
	batchSize int
}

func BatchElements(s beam.Scope, col beam.PCollection) {
	beam.ParDo0(s, &GlobalWindowsBatchingDoFn{batchSizeEstimator: BatchSizeEstimator{
		minBatchSize:            0,
		maxBatchSize:            0,
		targetBatchOverhead:     0,
		targetBatchDurationSecs: 0,
		variance:                0,
		clock:                   time.Time{},
	}, batchSize: 2}, col)
}

func (fn *GlobalWindowsBatchingDoFn) Setup(ctx context.Context, estimator BatchSizeEstimator) {
	fn.batchSizeEstimator = estimator
}

//func (fn *GlobalWindowsBatchingDoFn) StartBundle(ctx context.Context, estimator BatchSizeEstimator) {
//
//}

func (fn GlobalWindowsBatchingDoFn) ProcessElement(ctx context.Context, col beam.X) {
	fmt.Printf("Yolo %v %v", fn, col)
}

func (fn GlobalWindowsBatchingDoFn) FinishBundle(ctx context.Context) {

}

func (fn GlobalWindowsBatchingDoFn) Teardown(ctx context.Context) {
	fmt.Println("Teardown @GlobalWindowsBatchingDoFn")
}
