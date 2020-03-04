package env

import (
	"log"
	"os"
	"strconv"
	"time"
)

func CheckRunKey(key string) {
	if run, _ := strconv.ParseBool(os.Getenv(key)); !run {
		log.Printf("For run set environment %s to true.", key)
		os.Exit(1)
	}
}

func GetEnvOrPanic(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Panicf("%v is not defined", key)
	}
	return value
}

func GetEnvOrDefault(key, defval string) string {
	value := os.Getenv(key)
	if value == "" {
		return defval
	}
	return value
}

func GetEnvBoolOrPanic(key string) bool {
	value := GetEnvOrPanic(key)

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		log.Panicf("Error parse env '%v': %v", key, err)
	}
	return boolValue
}

func GetEnvBoolOrDefault(key string, defval bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defval
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		log.Panicf("Error parse env '%v': %v", key, err)
	}
	return boolValue
}

func GetEnvDurationOrPanic(key string) time.Duration {
	value := GetEnvOrPanic(key)

	durationValue, err := time.ParseDuration(value)
	if err != nil {
		log.Panicf("Error parse env '%v': %v", key, err)
	}
	return durationValue
}

func GetEnvDurationOrDefault(key string, defval time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defval
	}

	durationValue, err := time.ParseDuration(value)
	if err != nil {
		log.Panicf("Error parse env '%v': %v", key, err)
	}
	return durationValue
}

func GetEnvIntOrPanic(key string) int {
	value := GetEnvOrPanic(key)

	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Panicf("Error parse env '%v': %v", key, err)
	}
	return intValue
}

func GetEnvIntOrDefault(key string, defval int) int {
	value := os.Getenv(key)
	if value == "" {
		return defval
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Panicf("Error parse env '%v': %v", key, err)
	}
	return intValue
}

func GetUintOrDefault(key string, defval uint) uint {
	value := os.Getenv(key)
	if value == "" {
		return defval
	}
	v, err := strconv.ParseUint(value, 0, 0)
	if err == nil {
		return uint(v)
	}
	return defval
}

func GetFloatOrDefault(key string, defval float64) float64 {
	value := os.Getenv(key)
	if value == "" {
		return defval
	}
	v, err := strconv.ParseFloat(value, 0)
	if err == nil {
		return v
	}
	return defval
}
