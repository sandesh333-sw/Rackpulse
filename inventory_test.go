package main

import (
	"context"
	"testing"
)

func TestSyntheticInventory(t *testing.T) {
	source := SyntheticInventory{}

	devices, err := source.ListPDUs(context.Background())
	if err != nil {
		t.Fatalf("failed to list synthetic PDUs: %v", err)
	}

	if len(devices) != 3 {
		t.Fatalf("expected 3 devices, got %d", len(devices))
	}

	if devices[0].Name == "" {
		t.Error("expected device name")
	}

	if devices[0].Manufacturer == "" {
		t.Error("expected manufacturer")
	}

	if devices[0].Model == "" {
		t.Error("expected model")
	}
}

func TestSyntheticInventoryCancelledContext(t *testing.T) {
	source := SyntheticInventory{}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := source.ListPDUs(ctx)
	if err == nil {
		t.Fatal("expected context cancellation error")
	}
}

func TestDiscoverPDUs(t *testing.T) {
	config, err := LoadModelConfig("config/pdu_models.toml")
	if err != nil {
		t.Fatal(err)
	}

	pdus, err := DiscoverPDUs(
		context.Background(),
		SyntheticInventory{},
		config,
	)
	if err != nil {
		t.Fatalf("discovery failed: %v", err)
	}

	if len(pdus) != 3 {
		t.Fatalf("expected 3 PDUs, got %d", len(pdus))
	}

	// APC AP8953 should receive the APC standard profile.
	if len(pdus[0].Readings.Metrics) != 8 {
		t.Errorf("expected 8 metrics for pdu389, got %d",
			len(pdus[0].Readings.Metrics))
	}

	// Eaton should receive its own profile.
	if len(pdus[2].Readings.Metrics) != 8 {
		t.Errorf("expected 8 metrics for pdu522, got %d",
			len(pdus[2].Readings.Metrics))
	}
}
