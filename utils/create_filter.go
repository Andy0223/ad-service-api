package utils

import (
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func CreateFilter(queryParams map[string]string) bson.M {
	now := time.Now()

	filter := bson.M{
		"startAt": bson.M{"$lte": now},
		"endAt":   bson.M{"$gte": now},
	}

	if age, ok := queryParams["age"]; ok {
		age, _ := strconv.Atoi(age)
		filter["conditions.ageStart"] = bson.M{"$lte": age}
		filter["conditions.ageEnd"] = bson.M{"$gte": age}
	}

	if gender, ok := queryParams["gender"]; ok {
		filter["conditions.gender"] = bson.M{"$in": []string{gender}}
	}

	if country, ok := queryParams["country"]; ok {
		filter["conditions.country"] = bson.M{"$in": []string{country}}
	}

	if platform, ok := queryParams["platform"]; ok {
		filter["conditions.platform"] = bson.M{"$in": []string{platform}}
	}

	return filter
}
