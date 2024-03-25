package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
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

type MongodbAdRepository interface {
	Store(ctx context.Context, ad *Advertisement) error
	GetActiveAdCounts(ctx context.Context, now time.Time) (int, error)
	Fetch(ctx context.Context, filter primitive.M, limit, offset int) ([]*Advertisement, error)
}

type RedisAdRepository interface {
	IncrAdCountsByDate(ctx context.Context, key string) error
	GetAdCountsByDate(ctx context.Context, key string) (int, error)
}
