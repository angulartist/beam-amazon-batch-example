package utils

import "time"

type Graph int

const (
	None Graph = iota
	Pie
	Bar
)

const (
	DATE_LAYOUT         = "2006-01-02T15:04:05.000Z"
	HALF_DAY_IN_SECONDS = (24 / 2) * time.Hour
)
