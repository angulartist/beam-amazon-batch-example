package models

import "projects/amazon_reviews_beam/utils"

type ComputedRating struct {
	Mean float64 `json:"mean"`
	Min  int     `json:"min"`
	Max  int     `json:"max"`
}

type ComputedRatingOutput struct {
	ComputedRating `json:"data"`
	GraphType      utils.Graph `json:"graph_type"`
}
