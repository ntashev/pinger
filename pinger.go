package main // import "github.com/johnpeterharvey/pinger"

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	log.Print("Starting up...")

	interval, settings, err := GetSettings()
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	client := &http.Client{}

	for {
		err := DoCall(client, settings)
		if err != nil {
			log.Print(err)
		}
		log.Printf("Sleeping for %d seconds", interval)
		time.Sleep(time.Duration(interval) * time.Second)
	}
}

// GetSettings Read required settings from environment variables
func GetSettings() (int, map[string]string, error) {
	target := os.Getenv("TARGET_URL")
	method := os.Getenv("METHOD")
	interval, err := strconv.Atoi(os.Getenv("INTERVAL"))

	if target == "" || method == "" || err != nil || interval < 0 {
		return -1, map[string]string{}, errors.New("Environment variables were not set, returning error")
	}

	return interval, map[string]string{
			"target": target,
			"method": method},
		nil
}

// DoCall Do an HTTP call out to the server, with the specified properties
func DoCall(client *http.Client, settings map[string]string) error {
	log.Printf("Trying HTTP %s to %s", settings["method"], settings["target"])

	req, err := http.NewRequest(settings["method"], settings["target"], strings.NewReader("{}"))
	if err != nil {
		return fmt.Errorf("Failed to create HTTP request! %s", err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return fmt.Errorf("Error! Received error while contacting target! %s", err)
	}
	log.Printf("Received response %d\n", resp.StatusCode)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("Error! Received unexpected status code from target! %d", resp.StatusCode)
	}
	return nil
}
