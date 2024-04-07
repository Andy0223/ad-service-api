package validators

import (
	"ad-service-api/internal/models"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/pariz/gountries"
)

func ValidateAgeRange(ageStart, ageEnd int) error {
	if ageStart < 1 || ageStart > 100 {
		return errors.New("ageStart should be between 1 and 100")
	}
	if ageEnd < 1 || ageEnd > 100 {
		return errors.New("ageEnd should be between 1 and 100")
	}
	if ageStart > ageEnd {
		return errors.New("ageStart must be less than or equal to ageEnd")
	}
	return nil
}

func ValidateAgeQueryParam(age string) error {
	ageInt, err := strconv.Atoi(age)
	if err != nil {
		return fmt.Errorf("invalid age: %v", err)
	}
	if ageInt < 1 || ageInt > 100 {
		return errors.New("age should be between 1 and 100")
	}
	return nil
}

func ValidateGender(genders string) error {
	validGenders := map[string]bool{"M": true, "F": true}
	if _, ok := validGenders[genders]; !ok {
		return errors.New("invalid gender")
	}
	return nil
}

func ValidateCountry(country string) error {
	queryService := gountries.New()
	if _, err := queryService.FindCountryByAlpha(country); err != nil {
		return fmt.Errorf("invalid country: %v", country)
	}
	return nil
}

func ValidatePlatform(platforms string) error {
	validPlatforms := map[string]bool{"ios": true, "android": true, "web": true}
	if _, ok := validPlatforms[platforms]; !ok {
		return errors.New("invalid platform")
	}
	return nil
}

func ValidateLimit(limit string) error {
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return fmt.Errorf("invalid limit: %v", err)
	}
	if limitInt < 1 || limitInt > 100 {
		return errors.New("limit should be between 1 and 100")
	}
	return nil
}

func ValidateOffset(offset string) error {
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		return fmt.Errorf("invalid offset: %v", err)
	}
	if offsetInt < 0 {
		return errors.New("offset should be greater than or equal to 0")
	}
	return nil
}

func CreateAdValueValidation(ad models.Advertisement) error {
	// Validate startAt and endAt
	if ad.StartAt.After(ad.EndAt) {
		return errors.New("startAt must be before endAt")
	}

	if ad.EndAt.Before(time.Now()) {
		return errors.New("endAt must be after the current time")
	}

	// Validate age range
	if err := ValidateAgeRange(ad.Conditions.AgeStart, ad.Conditions.AgeEnd); err != nil {
		return err
	}

	// Validate genders
	for _, gender := range ad.Conditions.Gender {
		if err := ValidateGender(gender); err != nil {
			return err
		}
	}

	// Validate countries
	for _, country := range ad.Conditions.Country {
		if err := ValidateCountry(country); err != nil {
			return err
		}
	}

	// Validate platforms
	for _, platform := range ad.Conditions.Platform {
		if err := ValidatePlatform(platform); err != nil {
			return err
		}
	}

	return nil
}

func ListAdParamsValidation(query url.Values) (map[string]string, error) {
	validQueryParams := make(map[string]string)

	// AgeStart condition validation
	if age := query.Get("age"); age != "" {
		if err := ValidateAgeQueryParam(age); err != nil {
			return nil, fmt.Errorf("age validation failed: %w", err)
		}
		validQueryParams["age"] = age
	}

	// Gender condition validation
	if gender := query.Get("gender"); gender != "" {
		if err := ValidateGender(gender); err != nil {
			return nil, fmt.Errorf("gender validation failed: %w", err)
		}
		validQueryParams["gender"] = gender
	}

	// Country condition validation
	if country := query.Get("country"); country != "" {
		if err := ValidateCountry(country); err != nil {
			return nil, fmt.Errorf("country validation failed: %w", err)
		}
		validQueryParams["country"] = country
	}

	// Platform condition validation
	if platform := query.Get("platform"); platform != "" {
		if err := ValidatePlatform(platform); err != nil {
			return nil, fmt.Errorf("platform validation failed: %w", err)
		}
		validQueryParams["platform"] = platform
	}

	// Limit condition validation
	limit := query.Get("limit")
	if limit == "" {
		limit = "5" // Default value
	}
	if err := ValidateLimit(limit); err != nil {
		return nil, fmt.Errorf("limit validation failed: %w", err)
	}
	validQueryParams["limit"] = limit

	// Offset condition validation
	offset := query.Get("offset")
	if offset == "" {
		offset = "0" // Default value
	}
	if err := ValidateOffset(offset); err != nil {
		return nil, fmt.Errorf("offset validation failed: %w", err)
	}
	validQueryParams["offset"] = offset

	return validQueryParams, nil
}
