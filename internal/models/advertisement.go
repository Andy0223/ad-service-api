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
	AgeRange  AgeRange `json:"ageRange,omitempty" bson:"ageRange,omitempty"`
	Genders   []string `json:"genders,omitempty" bson:"genders,omitempty"`
	Countries []string `json:"countries,omitempty" bson:"countries,omitempty"`
	Platforms []string `json:"platforms,omitempty" bson:"platforms,omitempty"`
}

type AgeRange struct {
	AgeStart int `json:"ageStart" bson:"ageStart"`
	AgeEnd   int `json:"ageEnd" bson:"ageEnd"`
}

// type AdvertisementRepository interface {
// 	CreateAdvertisement(ctx context.Context, ad *Advertisement) error
// 	CountAdsCreatedToday(ctx context.Context, today string) (int, error)
// 	CountActiveAds(ctx context.Context, now time.Time) (int, error)
// 	ListAdvertisements(ctx context.Context, filter bson.M, limit, offset int) ([]*Advertisement, error)
// }

// type AdvertisementService interface {
// 	CreateAdvertisement(ctx context.Context, ad *Advertisement) error
// 	CountAdsCreatedToday(ctx context.Context, today string) (int, error)
// 	CountActiveAds(ctx context.Context, now time.Time) (int, error)
// 	ListAdvertisements(ctx context.Context, filter bson.M, limit, offset int) ([]*Advertisement, error)
// }
