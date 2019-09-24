package utils

import "flag"

var (
	Input   = flag.String("file", "./sources/amazon_reviews_sample*.csv", ".csv file input")
	Output  = flag.String("output", "./outputs/reporting.json", "output results")
	Logging = flag.String("logs", "./logs/", "output error logs")
)
