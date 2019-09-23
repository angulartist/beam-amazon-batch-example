package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/apache/beam/sdks/go/pkg/beam"
	"github.com/apache/beam/sdks/go/pkg/beam/log"
	"github.com/apache/beam/sdks/go/pkg/beam/x/beamx"
	"projects/amazon_reviews_beam/core"
)

// Main entry point.
func main() {
	/* Parse the command line into the defined flags.
	flags may then be used directly. If you're using the flags themselves, they are all pointers;
	if you bind to variables, they're values. */
	flag.Parse()

	/* Init() is an initialization hook that must called on startup.
	On distributed runners, it is used to intercept control. */
	beam.Init()

	/* NewPipelineWithRoot creates a new empty pipeline and its root scope.

	Pipeline manages a directed acyclic graph of primitive PTransforms,
	and the PCollections that the PTransforms consume and produce.
	Each Pipeline is self-contained and isolated from any other Pipeline.
	The Pipeline owns the PCollections and PTransforms and they can by used by that Pipeline only.
	Pipelines can safely be executed concurrently.

	Scope is a hierarchical grouping for composite transforms.
	Scopes can be enclosed in other scopes and for a tree structure.
	For pipeline updates, the scope chain form a unique name.
	The scope chain can also be used for monitoring and visualization purposes.*/
	p, s := beam.NewPipelineWithRoot()

	/* Background returns a non-nil, empty Context.
	It is never canceled, has no values, and has no deadline.
	It is typically used by the main function, initialization, and tests,
	and as the top-level Context for incoming requests. */
	ctx := context.Background()

	/* Let's encapsulate the main pipeline logic in an outer fn.
	-> Takes a context and a scope parameters. */

	core.DescribePipeline(ctx, s)

	log.Info(ctx, "Pipeline is booting...")

	/* Run executes the pipeline using the selected registered runner.
	It is customary to define a "runner" with no default as a flag to
	let users control runner selection.
	Any [error] will be caught and stop the program execution. */
	if err := beamx.Run(ctx, p); err != nil {
		fmt.Println(err)
		log.Exitf(ctx, "Failed to execute job: on ctx=%v:", ctx)
	}
}
