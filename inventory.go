package main

import (
	"context"
	"fmt"
)

// InventoryDevice represents inventory data discovered from an external source
type InventoryDevice struct {
	Name         string
	IPAddress    string
	RackNumber   string
	Location     string
	Role         string
	Tenant       string
	Manufacturer string
	Model        string
}

// InventorySource is the read-only boundary for inventory providers
type InventorySource interface {
	ListPDUs(ctx context.Context) ([]InventoryDevice, error)
}

// SyntheticInventory provides development inventory without external systems.
type SyntheticInventory struct{}

func (s SyntheticInventory) ListPDUs(ctx context.Context) ([]InventoryDevice, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	pdus := samplePDUs()

	devices := make([]InventoryDevice, 0, len(pdus))

	for _, pdu := range pdus {
		devices = append(devices, InventoryDevice{
			Name:         pdu.Name,
			IPAddress:    pdu.IPAddress,
			RackNumber:   pdu.RackNumber,
			Location:     pdu.Location,
			Role:         pdu.Role,
			Tenant:       pdu.Tenant,
			Manufacturer: pdu.Manufacturer,
			Model:        pdu.Model,
		})
	}

	return devices, nil
}

func DiscoverPDUs(
	ctx context.Context,
	source InventorySource,
	config ModelConfig,
) ([]PDU, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid model configuration: %w", err)
	}

	devices, err := source.ListPDUs(ctx)
	if err != nil {
		return nil, fmt.Errorf("discover inventory: %w", err)
	}

	pdus := make([]PDU, 0, len(devices))

	for _, device := range devices {
		pdu := PDU{
			Name:         device.Name,
			IPAddress:    device.IPAddress,
			RackNumber:   device.RackNumber,
			Location:     device.Location,
			Role:         device.Role,
			Tenant:       device.Tenant,
			Manufacturer: device.Manufacturer,
			Model:        device.Model,
			Readings: Reading{
				Metrics: make(map[string]MetricReading),
			},
		}

		_, profile, matched := config.MatchModel(
			device.Manufacturer,
			device.Model,
		)

		if matched {
			pdu.Readings = InitializeReadings(profile)
		}

		pdus = append(pdus, pdu)
	}

	return pdus, nil
}
