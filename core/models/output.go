package models

type OutputFormat struct {
	ReviewOutput         `bigquery:"top_helpful_reviews" json:"top_helpful_reviews"`
	UserReviewsOutput    []UserReviewsOutput `bigquery:"top_users_most_reviews" json:"top_users_most_reviews"`
	ComputedRatingOutput `bigquery:"ratings_overview" json:"ratings_overview"`
	RecommendOutput      []RecommendOutput `bigquery:"recommend_ratios" json:"recommend_ratios"`
}
