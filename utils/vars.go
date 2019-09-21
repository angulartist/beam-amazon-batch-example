package utils

import "time"

type Graph int

const (
	None Graph = iota
	Pie
	Bar
)

const HALF_DAY_IN_SECONDS = (24 / 2) * time.Hour
