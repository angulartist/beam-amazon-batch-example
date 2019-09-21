package models

import "projects/amazon_reviews_beam/utils"

type UserReviews struct {
	Username   string `json:"username"`
	NumRatings int    `json:"numRatings"`
}

type UserReviewsOutput struct {
	Data      []UserReviews `json:"data"`
	GraphType utils.Graph   `json:"graph_type"`
}
