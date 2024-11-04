package config

import (
	"os"
	"strconv"
	"strings"
)

const (
	MIN_BAUD_RATE int = 9600
	MAX_BAUD_RATE int = 115200
)

type Config struct {
	DebugMode          bool
	SerialPort         string
	BaudRate           int
	AutodetectBaudRate bool
}

func New() *Config {
	return &Config{
		DebugMode:          getEnvAsBool("DEBUG_MODE", false),
		SerialPort:         getEnv("SERIAL_PORT", "/dev/ttyAMA0"),
		BaudRate:           getEnvAsInt("BAUD_RATE", MIN_BAUD_RATE),
		AutodetectBaudRate: getEnvAsBool("AUTODETECT_BAUD_RATE", false),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}

func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}

func getEnvAsSlice(name string, defaultVal []string, sep string) []string {
	valStr := getEnv(name, "")

	if valStr == "" {
		return defaultVal
	}

	val := strings.Split(valStr, sep)

	return val
}
