package cfg

import (
	"log"
	"os"
	"strconv"
)

const (
	defAutTokenKey              = "Authorization"
	defMinPassLen               = 8
	defAPIListen                = ":80"
	defaultSignUpTokenLifetime  = 604800
	defaultRecoverTokenLifetime = 86400
)

type Settings struct {
	APIListen            string
	APIHost              string
	APISchema            string
	AuthTokenLength      int
	AuthTokenLifetime    int
	AuthTokenKey         string
	MinPassLen           int
	SignUpTokenLifetime  int
	RecoverTokenLifetime int
	FrontURL             string
}

var App Settings //nolint

func Load() {
	const errorsCategory = "Settings"

	var errors []string

	val := os.Getenv("AUTH_TOKEN_LENGTH")
	key, err := strconv.Atoi(val)
	if err != nil {
		errors = append(errors, errorsCategory+": Invalid auth token length setting.")
	} else {
		App.AuthTokenLength = key
	}

	val = os.Getenv("API_LISTEN")
	if val == "" {
		App.AuthTokenKey = defAPIListen
	} else {
		App.AuthTokenKey = val
	}

	// API host
	val = os.Getenv("API_HOST")
	if val == "" {
		errors = append(errors, errorsCategory+": Undefined API_HOST.")
	} else {
		App.APIHost = val
	}
	// API host
	val = os.Getenv("API_SCHEMA")
	if val == "" {
		errors = append(errors, errorsCategory+": Undefined API_SCHEMA.")
	} else {
		App.APISchema = val
	}

	val = os.Getenv("AUTH_TOKEN_LIFETIME")
	key, err = strconv.Atoi(val)
	if err != nil {
		errors = append(errors, errorsCategory+": Invalid auth token lifetime setting.")
	} else {
		App.AuthTokenLifetime = key
	}

	val = os.Getenv("AUTH_TOKEN_KEY")
	if val == "" {
		App.AuthTokenKey = defAutTokenKey
	} else {
		App.AuthTokenKey = val
	}

	val = os.Getenv("MIN_PASS_LEN")
	key, err = strconv.Atoi(val)
	if err != nil {
		App.MinPassLen = defMinPassLen
	} else {
		App.MinPassLen = key
	}

	val = os.Getenv("SIGN_UP_TOKEN_LIFETIME")
	key, err = strconv.Atoi(val)
	if err != nil {
		App.SignUpTokenLifetime = defaultSignUpTokenLifetime
	} else {
		App.SignUpTokenLifetime = key
	}

	val = os.Getenv("RECOVER_TOKEN_LIFETIME")
	key, err = strconv.Atoi(val)
	if err != nil {
		App.RecoverTokenLifetime = defaultRecoverTokenLifetime
	} else {
		App.RecoverTokenLifetime = key
	}

	// front url
	val = os.Getenv("FRONT_URL")
	if val == "" {
		errors = append(errors, errorsCategory+": Undefined FRONT_URL.")
	} else {
		App.FrontURL = val
	}

	if len(errors) > 0 {
		log.Panicln(errors)
	}
}
