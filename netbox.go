package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type NetBoxConfig struct {
	BaseURL     string
	Token       string
	PDURoleSlug string
}

type NetBoxInventory struct {
	config NetBoxConfig
	client *http.Client
}

func NewnetBoxInventory(config NetBoxConfig) (*NetBoxInventory, error) {
	config.BaseURL = strings.TrimRight(config.BaseURL, "/")

	if config.BaseURL == "" {
		return nil, fmt.Errorf("NetBox base URL is required")
	}

	if config.Token == "" {
		return nil, fmt.Errorf("NetBox API token is required")
	}

	_, err := url.ParseRequestURI(config.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid NetBox base URL: %w", err)
	}

	return &NetBoxInventory{
		config: config,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}, nil
}

type netBoxDevice struct {
	Name       string `json:"name"`
	DeviceType struct {
		Manufacturer struct {
			Name string `json:"name"`
		} `json:"manufacturer"`
		Model string `json:"model"`
	} `json:"device_type"`
	Role struct {
		Name string `json:"name"`
		Slug string `json:"slug"`
	} `json:"role"`
	Tenant *struct {
		Name string `json:"name"`
	} `json:"tenant"`
	Location *struct {
		Name string `json:"name"`
	} `json:"location"`
	Rack *struct {
		Name string `json:"name"`
	} `json:"rack"`
	PrimaryIP4 *struct {
		Address string `json:"address"`
	} `json:"primary_ip4"`
}

type netBoxResponse struct {
	Count    int            `json:"count"`
	Next     *string        `json:"next"`
	Previous *string        `json:"previous"`
	Results  []netBoxDevice `json:"results"`
}

func (inventory *NetBoxInventory) get(ctx context.Context, endpoint string) (*netBoxResponse, error) {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		endpoint,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("create NetBox request: %w", err)
	}

	request.Header.Set("Authorization", "Token "+inventory.config.Token)
	request.Header.Set("Accept", "application/json")

	response, err := inventory.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("request NetBox API: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("NetBox returned HTTP %d", response.StatusCode)
	}

	var result netBoxResponse

	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode NetBox response: %w", err)
	}

	return &result, nil
}

func (inventory *NetBoxInventory) ListPDUs(ctx context.Context) ([]InventoryDevice, error) {
	endpoint := inventory.config.BaseURL + "/api/dcim/devices/?limit=100"

	if inventory.config.PDURoleSlug != "" {
		endpoint += "&role=" + url.QueryEscape(inventory.config.PDURoleSlug)
	}

	devices := make([]InventoryDevice, 0)

	for endpoint != "" {
		result, err := inventory.get(ctx, endpoint)
		if err != nil {
			return nil, err
		}

		for _, device := range result.Results {
			// If configured, only include devices with the PDU role.
			if inventory.config.PDURoleSlug != "" &&
				device.Role.Slug != inventory.config.PDURoleSlug {
				continue
			}

			ipAddress := ""

			if device.PrimaryIP4 != nil {
				ipAddress = strings.Split(
					device.PrimaryIP4.Address,
					"/",
				)[0]
			}

			rackNumber := ""

			if device.Rack != nil {
				rackNumber = device.Rack.Name
			}

			location := ""

			if device.Location != nil {
				location = device.Location.Name
			}

			tenant := ""

			if device.Tenant != nil {
				tenant = device.Tenant.Name
			}

			devices = append(devices, InventoryDevice{
				Name:         device.Name,
				IPAddress:    ipAddress,
				RackNumber:   rackNumber,
				Location:     location,
				Role:         device.Role.Name,
				Tenant:       tenant,
				Manufacturer: device.DeviceType.Manufacturer.Name,
				Model:        device.DeviceType.Model,
			})
		}

		if result.Next == nil {
			break
		}

		nextURL, err := url.Parse(*result.Next)
		if err != nil {
			return nil, fmt.Errorf("invalid NetBox pagination URL: %w", err)
		}

		baseURL, err := url.Parse(inventory.config.BaseURL)
		if err != nil {
			return nil, fmt.Errorf("invalid configured NetBox URL: %w", err)
		}

		// Never send the API token to a different host.
		if nextURL.Host != "" && nextURL.Host != baseURL.Host {
			return nil, fmt.Errorf("NetBox pagination URL changed host")
		}

		if nextURL.Host == "" {
			nextURL = baseURL.ResolveReference(nextURL)
		}

		endpoint = nextURL.String()
	}

	return devices, nil
}
