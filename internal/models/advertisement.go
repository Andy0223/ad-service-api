package models

import (
	"time"
)

type Advertisement struct {
	Title      string     `json:"title" bson:"title"`
	StartAt    time.Time  `json:"startAt" bson:"startAt"`
	EndAt      time.Time  `json:"endAt" bson:"endAt"`
	Conditions Conditions `json:"conditions,omitempty" bson:"conditions,omitempty"`
}

type Conditions struct {
	AgeStart int      `json:"ageStart,omitempty" bson:"ageStart,omitempty"`
	AgeEnd   int      `json:"ageEnd,omitempty" bson:"ageEnd,omitempty"`
	Gender   []string `json:"gender,omitempty" bson:"gender,omitempty"`
	Country  []string `json:"country,omitempty" bson:"country,omitempty"`
	Platform []string `json:"platform,omitempty" bson:"platform,omitempty"`
}
