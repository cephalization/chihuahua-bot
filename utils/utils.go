package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

// GetEnv parses an environment variable and returns it if available
func GetEnv(name string) (string, error) {
	v := os.Getenv(name)
	if v == "" {
		return "", fmt.Errorf("missing required environment variable %s", name)
	}

	return v, nil
}

// AdjectivesFile - path to file with list of adjectives
const AdjectivesFile = "data/adjectives.txt"

// PopulateAdjectives populates a global list from a text file containing
// a list of adjectives.
func PopulateAdjectivesFromFile() []string {
	var adjectives []string

	file, err := os.Open(AdjectivesFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		adjectives = append(adjectives, scanner.Text())
	}

	return adjectives
}