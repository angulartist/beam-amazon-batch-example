package models

import "projects/amazon_reviews_beam/utils"

type RecommendRatio struct {
	DoRecommend string  `json:"doRecommend"`
	NumVotes    int     `json:"numVotes"`
	Percent     float64 `json:"Percent"`
}

type RecommendOutput struct {
	Data      RecommendRatio `json:"data"`
	GraphType utils.Graph    `json:"graph_type"`
}
