package models

import "projects/amazon_reviews_beam/utils"

type Review struct {
	Manufacturer string `bigquery:"manufacturer" csv:"manufacturer" json:"manufacturer"`
	Date         string `bigquery:"date" csv:"date" json:"date"`
	DoRecommend  string `bigquery:"doRecommend" csv:"doRecommend" json:"doRecommend"`
	NumHelpful   int    `bigquery:"numHelpful" csv:"numHelpful" json:"numHelpful"`
	Rating       int    `bigquery:"rating" csv:"rating" json:"rating"`
	Text         string `bigquery:"text" csv:"text" json:"text"`
	Title        string `bigquery:"title" csv:"title" json:"title"`
	Username     string `bigquery:"username" csv:"username" json:"username"`
	Timestamp    int64  `bigquery:"timestamp" csv:"timestamp" json:"timestamp"`
}

type ReviewOutput struct {
	Data      []Review    `json:"data"`
	GraphType utils.Graph `json:"graph_type"`
}
