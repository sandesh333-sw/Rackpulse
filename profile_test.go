package main

import "testing"

func TestLoadModelConfig(t *testing.T) {
	config, err := LoadModelConfig("config/pdu_models.toml")

	if err := config.Validate(); err != nil {
		t.Fatalf("configuration validation failed: %v", err)
	}

	if err != nil {
		t.Fatalf("failed to load model config: %v", err)
	}

	if len(config.Models) != 3 {
		t.Fatalf("expected 3 models, got %d", len(config.Models))
	}

	apdu, exists := config.Models["APDU10150SM"]
	if !exists {
		t.Fatal("APDU10150SM profile was not loaded")
	}

	if apdu.MatchManufacturer != "APC" {
		t.Errorf("expected manufacturer APC, got %q", apdu.MatchManufacturer)
	}

	if len(apdu.Readings) != 6 {
		t.Errorf("expected 6 APDU readings, got %d", len(apdu.Readings))
	}

	if apdu.Readings["bank_1"].OID != "" {
		t.Error("expected bank_1 to have no OID")
	}

	if apdu.Readings["input_voltage"].Unit != "V" {
		t.Errorf("expected input_voltage unit V, got %q",
			apdu.Readings["input_voltage"].Unit)
	}
	if apdu.Readings["input_voltage"].OID == "" {
		t.Error("expected input_voltage to have an OID")
	}
}

func TestMatchModel(t *testing.T) {
	config, err := LoadModelConfig("config/pdu_models.toml")
	if err != nil {
		t.Fatalf("failed to load model config: %v", err)
	}

	tests := []struct {
		name         string
		manufacturer string
		model        string
		wantProfile  string
		wantMatch    bool
	}{
		{
			name:         "APC AP8953",
			manufacturer: "APC",
			model:        "AP8953",
			wantProfile:  "apc-standard",
			wantMatch:    true,
		},
		{
			name:         "APC AP9953",
			manufacturer: "APC",
			model:        "AP9953",
			wantProfile:  "apc-standard",
			wantMatch:    true,
		},
		{
			name:         "APC APDU10150SM",
			manufacturer: "APC",
			model:        "APDU10150SM",
			wantProfile:  "APDU10150SM",
			wantMatch:    true,
		},
		{
			name:         "Eaton model",
			manufacturer: "Eaton",
			model:        "Eaton-24-Port",
			wantProfile:  "Eaton-24-Port",
			wantMatch:    true,
		},
		{
			name:         "unknown APC model",
			manufacturer: "APC",
			model:        "Unknown-PDU",
			wantMatch:    false,
		},
		{
			name:         "unknown manufacturer",
			manufacturer: "Unknown",
			model:        "Unknown-PDU",
			wantMatch:    false,
		},
		{
			name:         "case insensitive",
			manufacturer: "apc",
			model:        "ap8953",
			wantProfile:  "apc-standard",
			wantMatch:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profileName, _, matched := config.MatchModel(
				tt.manufacturer,
				tt.model,
			)

			if matched != tt.wantMatch {
				t.Fatalf("expected match %v, got %v", tt.wantMatch, matched)
			}

			if matched && profileName != tt.wantProfile {
				t.Errorf(
					"expected profile %q, got %q",
					tt.wantProfile,
					profileName,
				)
			}
		})
	}
}

func TestInitializeReadings(t *testing.T) {
	config, err := LoadModelConfig("config/pdu_models.toml")
	if err != nil {
		t.Fatalf("failed to load model config: %v", err)
	}

	_, profile, matched := config.MatchModel("APC", "APDU10150SM")
	if !matched {
		t.Fatal("failed to match APDU10150SM profile")
	}

	readings := InitializeReadings(profile)

	// A supported metric should initially be unavailable.
	inputVoltage := readings.Metrics["input_voltage"]

	if inputVoltage.Status != MetricStatusUnavailable {
		t.Errorf(
			"expected input_voltage status unavailable, got %q",
			inputVoltage.Status,
		)
	}

	if inputVoltage.Unit != "V" {
		t.Errorf("expected input_voltage unit V, got %q", inputVoltage.Unit)
	}

	// A metric without an OID is not applicable.
	bank1 := readings.Metrics["bank_1"]

	if bank1.Status != MetricStatusNotApplicable {
		t.Errorf(
			"expected bank_1 status not_applicable, got %q",
			bank1.Status,
		)
	}

	if bank1.Value != nil {
		t.Error("expected bank_1 value to be nil")
	}
}
