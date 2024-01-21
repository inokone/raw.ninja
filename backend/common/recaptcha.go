package common

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/rs/zerolog/log"
)

// RecaptchaValidator is a struct for validating captcha using Google's ReCapthca validator
type RecaptchaValidator struct {
	secretKey string
}

// NewRecaptchaValidator is a function creating a new `RecaptchaValidator` based on the API secret
func NewRecaptchaValidator(secret string) *RecaptchaValidator {
	return &RecaptchaValidator{secret}
}

// Verify is a method of `RecaptchaValidator` verifying the captch token from the frontend
func (v RecaptchaValidator) Verify(token string) (bool, error) {
	values := url.Values{}
	values.Add("secret", v.secretKey)
	values.Add("response", token)

	// Send a request to Google's reCAPTCHA API for verification
	response, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify", values)
	if err != nil {
		log.Err(err).Msg("Captcha call failed.")
		return false, err
	}
	defer response.Body.Close()

	// Parse the JSON response
	var recaptchaResponse RecaptchaResponse
	err = json.NewDecoder(response.Body).Decode(&recaptchaResponse)
	if err != nil {
		log.Err(err).Msg("Captcha call failed.")
		return false, err
	}

	// Check if the reCAPTCHA verification was successful
	return recaptchaResponse.Success, nil
}
