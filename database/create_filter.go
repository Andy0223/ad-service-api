package database

import (
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func CreateFilter(validQueryParams map[string]string) bson.M {
	now := time.Now()

	filter := bson.M{
		"startAt": bson.M{"$lte": now},
		"endAt":   bson.M{"$gte": now},
	}

	if age, ok := validQueryParams["age"]; ok {
		age, _ := strconv.Atoi(age)
		filter["conditions.ageStart"] = bson.M{"$lte": age}
		filter["conditions.ageEnd"] = bson.M{"$gte": age}
	}

	if gender, ok := validQueryParams["gender"]; ok {
		filter["conditions.gender"] = bson.M{"$in": []string{gender}}
	}

	if country, ok := validQueryParams["country"]; ok {
		filter["conditions.country"] = bson.M{"$in": []string{country}}
	}

	if platform, ok := validQueryParams["platform"]; ok {
		filter["conditions.platform"] = bson.M{"$in": []string{platform}}
	}

	return filter
}
