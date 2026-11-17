package main

type PDU struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Rack     string  `json:"rack"`
	VoltageV float64 `json:"voltage_v"`
	CurrentA float64 `json:"current_a"`
	PowerW   float64 `json:"power_w"`
}

func samplePDUs() []PDU {
	return []PDU{
		{
			ID:       "pdu-001",
			Name:     "Rack A - PDU 01",
			Rack:     "Rack A",
			VoltageV: 230.4,
			CurrentA: 4.2,
			PowerW:   850.0,
		},
		{
			ID:       "pdu-002",
			Name:     "Rack B - PDU 01",
			Rack:     "Rack B",
			VoltageV: 229.8,
			CurrentA: 5.1,
			PowerW:   1100.0,
		},
		{
			ID:       "pdu-003",
			Name:     "Rack C - PDU 01",
			Rack:     "Rack C",
			VoltageV: 231.2,
			CurrentA: 6.5,
			PowerW:   1450.0,
		},
	}
}
