package main

import "context"

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
	ListPDUS(ctx context.Context) ([]InventoryDevice, error)
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
