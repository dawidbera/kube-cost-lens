// Package finops provides cost calculation, price book models,
// and rightsizing recommendation logic for Kubernetes workloads.
package finops

import (
	"encoding/json"
	"fmt"
	"os"
)

// Standard hours in an average accounting month (365.25 days / 12 months * 24 hours).
const DefaultHoursPerMonth = 730.0

// Default hourly cloud rates based on blended average compute instances.
const (
	DefaultCPURatePerHour      = 0.0416 // ~$30.37 per vCPU-month
	DefaultMemoryRatePerGBHour = 0.0055 // ~$4.02 per GB-month
)

// PricingCard defines unit pricing parameters for compute resources.
type PricingCard struct {
	// Provider identifies the cloud or datacenter provider (e.g. AWS, GCP, Azure, Custom).
	Provider string `json:"provider"`
	// Currency denotes the 3-letter currency code (e.g. USD, EUR).
	Currency string `json:"currency"`
	// CPURatePerHour specifies the financial cost per 1 vCPU-hour.
	CPURatePerHour float64 `json:"cpuRatePerHour"`
	// MemoryRatePerGBHour specifies the financial cost per 1 Gigabyte of RAM per hour.
	MemoryRatePerGBHour float64 `json:"memoryRatePerGbHour"`
	// HoursPerMonth specifies the number of hours utilized for monthly projections (default: 730).
	HoursPerMonth float64 `json:"hoursPerMonth"`
}

// DefaultPricingCard returns standard blended cloud pricing (AWS/GCP/Azure average).
func DefaultPricingCard() PricingCard {
	return PricingCard{
		Provider:            "Blended Cloud Average",
		Currency:            "USD",
		CPURatePerHour:      DefaultCPURatePerHour,
		MemoryRatePerGBHour: DefaultMemoryRatePerGBHour,
		HoursPerMonth:       DefaultHoursPerMonth,
	}
}

// AWSPricingCard returns standard on-demand pricing rates for AWS EC2 general purpose instances.
func AWSPricingCard() PricingCard {
	return PricingCard{
		Provider:            "AWS (m6i/m7i on-demand blended)",
		Currency:            "USD",
		CPURatePerHour:      0.0404,
		MemoryRatePerGBHour: 0.0050,
		HoursPerMonth:       DefaultHoursPerMonth,
	}
}

// GCPPricingCard returns standard on-demand pricing rates for GCP E2/N2 standard instances.
func GCPPricingCard() PricingCard {
	return PricingCard{
		Provider:            "GCP (e2-standard on-demand blended)",
		Currency:            "USD",
		CPURatePerHour:      0.0382,
		MemoryRatePerGBHour: 0.0051,
		HoursPerMonth:       DefaultHoursPerMonth,
	}
}

// AzurePricingCard returns standard on-demand pricing rates for Azure D-series v5 instances.
func AzurePricingCard() PricingCard {
	return PricingCard{
		Provider:            "Azure (Dsv5 on-demand blended)",
		Currency:            "USD",
		CPURatePerHour:      0.0420,
		MemoryRatePerGBHour: 0.0053,
		HoursPerMonth:       DefaultHoursPerMonth,
	}
}

// LoadPricingCard loads a custom pricing configuration from a JSON file.
// If file reading or parsing fails, an error is returned.
func LoadPricingCard(filePath string) (PricingCard, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return PricingCard{}, fmt.Errorf("failed to read pricing file %q: %w", filePath, err)
	}

	card := DefaultPricingCard()
	if err := json.Unmarshal(data, &card); err != nil {
		return PricingCard{}, fmt.Errorf("failed to parse JSON pricing data from %q: %w", filePath, err)
	}

	if card.HoursPerMonth <= 0 {
		card.HoursPerMonth = DefaultHoursPerMonth
	}
	if card.Currency == "" {
		card.Currency = "USD"
	}

	return card, nil
}
