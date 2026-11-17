package main

type PDU struct {
	Name         string  `json:"name"`
	IPAddress    string  `json:"ip_address"`
	RackNumber   string  `json:"rack_number"`
	Location     string  `json:"location"`
	Role         string  `json:"role"`
	Tenant       string  `json:"tenant"`
	Manufacturer string  `json:"manufacturer"`
	Model        string  `json:"model"`
	Readings     Reading `json:"readings"`
}

func samplePDUs() []PDU {
	pdu1 := NewPDU("pdu389")
	pdu1.IPAddress = "192.0.2.10"
	pdu1.RackNumber = "160"
	pdu1.Location = "Row 7"
	pdu1.Role = "Compute"
	pdu1.Tenant = "Tier1"
	pdu1.Manufacturer = "APC"
	pdu1.Model = "AP8953"

	pdu2 := NewPDU("pdu390")
	pdu2.IPAddress = "192.0.2.11"
	pdu2.RackNumber = "160"
	pdu2.Location = "Row 7"
	pdu2.Role = "Compute"
	pdu2.Tenant = "Tier1"
	pdu2.Manufacturer = "APC"
	pdu2.Model = "AP9953"

	pdu3 := NewPDU("pdu522")
	pdu3.IPAddress = "192.0.2.12"
	pdu3.RackNumber = "150"
	pdu3.Location = "Row 3"
	pdu3.Role = "Mixed"
	pdu3.Tenant = "Tier1"
	pdu3.Manufacturer = "Eaton"
	pdu3.Model = "Eaton-24-Port"

	return []PDU{pdu1, pdu2, pdu3}
}

func NewPDU(name string) PDU {
	return PDU{
		Name: name,
		Readings: Reading{
			Metrics: make(map[string]MetricReading),
		},
	}
}
