package proxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const (
	proxyDefaultPort = 3000
)

type proxyTarget struct {
	Port    int          `json:"port"`
	Service proxyService `json:"service"`
}

type Settings struct {
	Routes  map[string]string      `json:"routes"`
	Targets map[string]proxyTarget `json:"targets"`
	Port    int                    `json:"port"`
	Default string                 `json:"default"`
}

func NewSettings(path string) (Settings, error) {
	s := Settings{}
	err := s.load(path)
	if err != nil {
		return s, err
	}
	return s, nil
}

func (s *Settings) load(path string) error {
	file, err := os.Open(path)
	if err != nil {
		log.Printf("Error opening file: %s", err)
		return err
	}
	defer file.Close()

	contents := make([]byte, 1024)
	_, err = file.Read(contents)
	if err != nil {
		log.Printf("Error reading file: %s", err)
		return err
	}
	contents = bytes.Trim(contents, "\x00")

	// Decode the JSON data into the Settings struct
	err = json.Unmarshal(contents, &s)
	if err != nil {
		log.Printf("Error decoding JSON: %s", err)
		return err
	}

	if len(s.Routes) == 0 {
		return fmt.Errorf("No routes found in settings file")
	}

	if len(s.Targets) == 0 {
		return fmt.Errorf("No targets found in settings file")
	}

	if s.Default != "" {
		if _, ok := s.Targets[s.Default]; !ok {
			return fmt.Errorf("Default target not found in settings file")
		}
	}

	return nil
}
