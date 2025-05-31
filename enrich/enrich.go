package enrich

import (
    "encoding/json"
    "fmt"
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
        CountryID string  `json:"country_id"`
        Probability float64 `json:"probability"`
    } `json:"country"`
}

func GetAge(name string) (int, error) {
    resp, err := http.Get(fmt.Sprintf("https://api.agify.io/?name=%s", name))
    if err != nil {
        return 0, err
    }
    defer resp.Body.Close()

    var result AgeResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return 0, err
    }
    return result.Age, nil
}

func GetGender(name string) (string, error) {
    resp, err := http.Get(fmt.Sprintf("https://api.genderize.io/?name=%s", name))
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    var result GenderResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return "", err
    }
    return result.Gender, nil
}

func GetNationality(name string) (string, error) {
    resp, err := http.Get(fmt.Sprintf("https://api.nationalize.io/?name=%s", name))
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    var result NationalityResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return "", err
    }
    if len(result.Country) > 0 {
        return result.Country[0].CountryID, nil
    }
    return "", nil
}
