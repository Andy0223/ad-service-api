package validators

import (
	"ad-service-api/internal/models"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/pariz/gountries"
	"go.mongodb.org/mongo-driver/bson"
)

func ValidateAgeRange(ageStart, ageEnd int) error {
	if ageStart < 1 || ageStart > 100 {
		return errors.New("ageStart out of range")
	}
	if ageEnd < 1 || ageEnd > 100 {
		return errors.New("ageEnd out of range")
	}
	if ageStart > ageEnd {
		return errors.New("ageStart must be less than or equal to ageEnd")
	}
	return nil
}

func ValidateAgeQueryParam(age int) error {
	if age < 1 || age > 100 {
		return errors.New("ageStart out of range")
	}
	return nil
}

func ValidateGenders(genders []string) error {
	validGenders := map[string]bool{"M": true, "F": true}
	for _, gender := range genders {
		if _, ok := validGenders[gender]; !ok {
			return fmt.Errorf("invalid gender: %v", gender)
		}
	}
	return nil
}

func ValidateCountries(countries []string) error {
	queryService := gountries.New()
	for _, country := range countries {
		_, err := queryService.FindCountryByAlpha(country)
		if err != nil {
			return err
		}
	}
	return nil
}

func ValidatePlatforms(platforms []string) error {
	validPlatforms := map[string]bool{"ios": true, "android": true, "web": true}
	for _, platform := range platforms {
		if _, ok := validPlatforms[platform]; !ok {
			return fmt.Errorf("invalid platform: %v", platform)
		}
	}
	return nil
}

func CreateAdValueValidation(ad models.Advertisement) error {
	// Validate startAt and endAt
	if ad.StartAt.After(ad.EndAt) {
		return errors.New("startAt must be before endAt")
	}

	// Validate age range
	if err := ValidateAgeRange(ad.Conditions.AgeStart, ad.Conditions.AgeEnd); err != nil {
		return err
	}

	// Validate genders
	if err := ValidateGenders(ad.Conditions.Genders); err != nil {
		return err
	}

	// Validate countries
	if err := ValidateCountries(ad.Conditions.Countries); err != nil {
		return err
	}

	// Validate platforms
	if err := ValidatePlatforms(ad.Conditions.Platforms); err != nil {
		return err
	}

	return nil
}

func ListAdParamsValidation(query url.Values) (bson.M, error) {
	filter := bson.M{
		"startAt": bson.M{"$lte": time.Now()},
		"endAt":   bson.M{"$gte": time.Now()},
	}

	// Age condition validation
	if ageStr := query.Get("age"); ageStr != "" {
		age, err := strconv.Atoi(ageStr) // String to int
		if err != nil {
			return nil, fmt.Errorf("invalid age: %v", err)
		}
		if err := ValidateAgeQueryParam(age); err != nil {
			return nil, err
		}

		filter["conditions.ageRange.ageStart"] = bson.M{"$lte": age}
		filter["conditions.ageRange.ageEnd"] = bson.M{"$gte": age}
	}

	// Gender condition validation
	if genders, ok := query["gender"]; ok && len(genders) > 0 {
		if err := ValidateGenders(genders); err != nil {
			return nil, err
		}
		filter["conditions.genders"] = bson.M{"$in": genders}
	}

	// Country condition validation
	if countries, ok := query["country"]; ok && len(countries) > 0 {
		if err := ValidateCountries(countries); err != nil {
			return nil, err
		}
		filter["conditions.countries"] = bson.M{"$in": countries}
	}

	// Platform condition validation
	if platforms, ok := query["platform"]; ok && len(platforms) > 0 {
		if err := ValidatePlatforms(platforms); err != nil {
			return nil, err
		}
		filter["conditions.platforms"] = bson.M{"$in": platforms}
	}

	return filter, nil
}
