package models

type Base struct {
	Status int
	Error  string
}
type Rating struct {
	Base
	Stars int
	Color string
}
type Review struct {
	Rating   Rating
	Text     string
	Reviewer string
}
type ReviewsBody struct {
	Base
	Reviews     []Review
	Podname     string
	Clustername string
}
type ReviewsWithStatus struct {
	Base
	ReviewsBody
}
type DetailsWithStatus struct {
	Base
	Details
}
type Details struct {
	ISBN10    string
	Publisher string
	Pages     int
	Type      string
	Language  string
}
type Product struct {
	Title           string
	DescriptionHtml string
}
type ProductPage struct {
	User    string
	Product Product
	Details DetailsWithStatus
	Reviews ReviewsWithStatus
}
