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
