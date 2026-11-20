package main

import (
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"os"
	"strings"
)

type MetricProfile struct {
	OID         string  `toml:"oid"`
	ScaleFactor float64 `toml:"scale_factor"`
	Unit        string  `toml:"unit"`
}

type ModelProfile struct {
	MatchManufacturer  string                   `toml:"match_manufacturer"`
	MatchModelPrefixes []string                 `toml:"match_model_prefixes"`
	Readings           map[string]MetricProfile `toml:"readings"`
}

type ModelConfig struct {
	Models map[string]ModelProfile `toml:"models"`
}

func LoadModelConfig(path string) (ModelConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return ModelConfig{}, fmt.Errorf("read model config: %w", err)
	}

	var config ModelConfig

	if err := toml.Unmarshal(data, &config); err != nil {
		return ModelConfig{}, fmt.Errorf("parse model config: %w", err)
	}

	return config, nil
}

func (config ModelConfig) Validate() error {
	if len(config.Models) == 0 {
		return fmt.Errorf("no model profiles configured")
	}

	for modelName, model := range config.Models {
		if strings.TrimSpace(model.MatchManufacturer) == "" {
			return fmt.Errorf("model %q has no manufacturer", modelName)
		}

		if len(model.Readings) == 0 {
			return fmt.Errorf("model %q has no readings", modelName)
		}

		for readingName, reading := range model.Readings {
			if strings.TrimSpace(reading.Unit) == "" {
				return fmt.Errorf(
					"model %q reading %q has no unit",
					modelName,
					readingName,
				)
			}
		}
	}

	return nil
}

func (config ModelConfig) MatchModel(manufacturer, model string) (string, ModelProfile, bool) {
	manufacturer = strings.TrimSpace(manufacturer)
	model = strings.TrimSpace(model)

	for profileName, profile := range config.Models {
		if !strings.EqualFold(manufacturer, profile.MatchManufacturer) {
			continue
		}

		// An empty prefix list means this profile matches
		// every model from the specified manufacturer.
		if len(profile.MatchModelPrefixes) == 0 {
			return profileName, profile, true
		}

		for _, prefix := range profile.MatchModelPrefixes {
			if strings.HasPrefix(strings.ToUpper(model), strings.ToUpper(prefix)) {
				return profileName, profile, true
			}
		}
	}

	return "", ModelProfile{}, false
}

func InitializeReadings(profile ModelProfile) Reading {
	readings := Reading{
		Metrics: make(map[string]MetricReading),
	}

	for metricName, metricProfile := range profile.Readings {
		status := MetricStatusUnavailable

		if strings.TrimSpace(metricProfile.OID) == "" {
			status = MetricStatusNotApplicable
		}

		readings.Metrics[metricName] = MetricReading{
			Value:  nil,
			Unit:   metricProfile.Unit,
			Status: status,
		}
	}

	return readings
}
