package config

import (
	"bytes"
 
	"fmt"
	"os"
    "log"
	"github.com/joho/godotenv"
)
func LoadEmbeddedEnv(envFile []byte) error {
    // Parse the embedded .env content
    envMap, err := godotenv.Parse(bytes.NewReader(envFile))
    if err != nil {
        return fmt.Errorf("error parsing embedded .env: %v", err)
    }

    // Set each key-value pair into the environment
    for key, value := range envMap {
        if err := os.Setenv(key, value); err != nil {
            log.Printf("Error setting environment variable %s: %v", key, err)
        }
    }
    log.Println("Embedded .env variables have been loaded into the environment.")
    return nil
}

func LoadEnv() {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }
}

func GetEnv(key string) string {
    value, exists := os.LookupEnv(key)
    if !exists {
        log.Fatalf("Environment variable %s not set", key)
    }
    return value
}
