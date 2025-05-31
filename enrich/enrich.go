package enrich

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type AgeResponse struct {
	Age int `json:"age"`
}

type GenderResponse struct {
	Gender string `json:"gender"`
}

type NationalityResponse struct {
	Country []struct {
		CountryID   string  `json:"country_id"`
		Probability float64 `json:"probability"`
	} `json:"country"`
}

// @Summary Get estimated age by name
// @Description Calls agify.io API to get estimated age
// @Tags Enrich
// @Produce json
// @Param name query string true "Name"
// @Success 200 {integer} int
// @Failure 500 {string} string
func GetAge(name string) (int, error) {
	log.Printf("[INFO] Fetching age for name: %s", name)
	resp, err := http.Get(fmt.Sprintf("https://api.agify.io/?name=%s", name))
	if err != nil {
		log.Printf("[ERROR] Failed to fetch age: %v", err)
		return 0, err
	}
	defer resp.Body.Close()

	var result AgeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("[ERROR] Failed to decode age response: %v", err)
		return 0, err
	}

	log.Printf("[DEBUG] Age API result for %s: %d", name, result.Age)
	return result.Age, nil
}

// @Summary Get estimated gender by name
// @Description Calls genderize.io API to get estimated gender
// @Tags Enrich
// @Produce json
// @Param name query string true "Name"
// @Success 200 {string} string
// @Failure 500 {string} string
func GetGender(name string) (string, error) {
	log.Printf("[INFO] Fetching gender for name: %s", name)
	resp, err := http.Get(fmt.Sprintf("https://api.genderize.io/?name=%s", name))
	if err != nil {
		log.Printf("[ERROR] Failed to fetch gender: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	var result GenderResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("[ERROR] Failed to decode gender response: %v", err)
		return "", err
	}

	log.Printf("[DEBUG] Gender API result for %s: %s", name, result.Gender)
	return result.Gender, nil
}

// @Summary Get estimated nationality by name
// @Description Calls nationalize.io API to get estimated nationality
// @Tags Enrich
// @Produce json
// @Param name query string true "Name"
// @Success 200 {string} string
// @Failure 500 {string} string
func GetNationality(name string) (string, error) {
	log.Printf("[INFO] Fetching nationality for name: %s", name)
	resp, err := http.Get(fmt.Sprintf("https://api.nationalize.io/?name=%s", name))
	if err != nil {
		log.Printf("[ERROR] Failed to fetch nationality: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	var result NationalityResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("[ERROR] Failed to decode nationality response: %v", err)
		return "", err
	}

	if len(result.Country) > 0 {
		log.Printf("[DEBUG] Nationality API result for %s: %s", name, result.Country[0].CountryID)
		return result.Country[0].CountryID, nil
	}

	log.Printf("[WARN] No nationality found for %s", name)
	return "", nil
}
