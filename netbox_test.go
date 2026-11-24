package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNetBoxListPDUs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}

		if r.Header.Get("Authorization") != "Token test-token" {
			t.Errorf("incorrect authorization header")
		}

		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("expected JSON Accept header")
		}

		if r.URL.Path != "/api/dcim/devices/" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		if r.URL.Query().Get("role") != "pdu" {
			t.Errorf("expected role=pdu, got %q", r.URL.Query().Get("role"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte(`{
			"count": 1,
			"next": null,
			"previous": null,
			"results": [
				{
					"name": "pdu389",
					"device_type": {
						"manufacturer": {"name": "APC"},
						"model": "AP8953"
					},
					"role": {"name": "PDU", "slug": "pdu"},
					"tenant": {"name": "Tier1"},
					"location": {"name": "Row 7"},
					"rack": {"name": "160"},
					"primary_ip4": {"address": "192.0.2.10/24"}
				}
			]
		}`))
	}))
	defer server.Close()

	inventory, err := NewNetBoxInventory(NetBoxConfig{
		BaseURL:     server.URL,
		Token:       "test-token",
		PDURoleSlug: "pdu",
	})
	if err != nil {
		t.Fatalf("failed to create NetBox inventory: %v", err)
	}

	devices, err := inventory.ListPDUs(context.Background())
	if err != nil {
		t.Fatalf("failed to list PDUs: %v", err)
	}

	if len(devices) != 1 {
		t.Fatalf("expected 1 device, got %d", len(devices))
	}

	device := devices[0]

	if device.Name != "pdu389" {
		t.Errorf("expected name pdu389, got %q", device.Name)
	}

	if device.Manufacturer != "APC" {
		t.Errorf("expected manufacturer APC, got %q", device.Manufacturer)
	}

	if device.Model != "AP8953" {
		t.Errorf("expected model AP8953, got %q", device.Model)
	}

	if device.IPAddress != "192.0.2.10" {
		t.Errorf("expected IP 192.0.2.10, got %q", device.IPAddress)
	}

	if device.RackNumber != "160" {
		t.Errorf("expected rack 160, got %q", device.RackNumber)
	}

	if device.Location != "Row 7" {
		t.Errorf("expected location Row 7, got %q", device.Location)
	}

	if device.Tenant != "Tier1" {
		t.Errorf("expected tenant Tier1, got %q", device.Tenant)
	}
}

func TestNetBoxPagination(t *testing.T) {
	var server *httptest.Server
	pageRequests := 0

	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pageRequests++
		w.Header().Set("Content-Type", "application/json")

		if r.URL.Query().Get("page") == "2" {
			_, _ = w.Write([]byte(`{
				"count": 2,
				"next": null,
				"previous": null,
				"results": [{
					"name": "pdu002",
					"device_type": {
						"manufacturer": {"name": "Eaton"},
						"model": "Eaton-24-Port"
					},
					"role": {"name": "PDU", "slug": "pdu"}
				}]
			}`))
			return
		}

		_, _ = w.Write([]byte(`{
			"count": 2,
			"next": "` + server.URL + `/api/dcim/devices/?page=2",
			"previous": null,
			"results": [{
				"name": "pdu001",
				"device_type": {
					"manufacturer": {"name": "APC"},
					"model": "AP8953"
				},
				"role": {"name": "PDU", "slug": "pdu"}
			}]
		}`))
	}))
	defer server.Close()

	inventory, err := NewNetBoxInventory(NetBoxConfig{
		BaseURL: server.URL,
		Token:   "test-token",
	})
	if err != nil {
		t.Fatal(err)
	}

	devices, err := inventory.ListPDUs(context.Background())
	if err != nil {
		t.Fatalf("pagination failed: %v", err)
	}

	if len(devices) != 2 {
		t.Fatalf("expected 2 devices, got %d", len(devices))
	}

	if pageRequests != 2 {
		t.Errorf("expected 2 page requests, got %d", pageRequests)
	}
}
