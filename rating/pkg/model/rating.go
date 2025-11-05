package model

type RecordID int64
type RecordType string

const (
	RecordTypeMovie = RecordType("movie")
)

type UserID string
type RatingValue int

type Rating struct {
	RecordID   RecordID    `json:"recordId"`
	RecordType RecordType  `json:"recordType"`
	UserID     UserID      `json:"userId"`
	Value      RatingValue `json:"value"`
}

type RatingEvent struct {
	Rating
	ProviderID string          `json:"providerId"`
	EventType  RatingEventType `json:"eventType"`
}

type RatingEventType string

const (
	RatingEventTypePut    = RatingEventType("put")
	RatingEventTypeDelete = RatingEventType("delete")
)
