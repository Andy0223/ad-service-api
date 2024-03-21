package validators

import (
	"errors"
	"fmt"

	"github.com/pariz/gountries"
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
