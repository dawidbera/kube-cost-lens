package finops

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultPricingCard(t *testing.T) {
	card := DefaultPricingCard()

	if card.Currency != "USD" {
		t.Errorf("expected USD currency, got %s", card.Currency)
	}
	if card.CPURatePerHour != DefaultCPURatePerHour {
		t.Errorf("expected CPU rate %f, got %f", DefaultCPURatePerHour, card.CPURatePerHour)
	}
	if card.MemoryRatePerGBHour != DefaultMemoryRatePerGBHour {
		t.Errorf("expected Memory rate %f, got %f", DefaultMemoryRatePerGBHour, card.MemoryRatePerGBHour)
	}
	if card.HoursPerMonth != DefaultHoursPerMonth {
		t.Errorf("expected Hours per month %f, got %f", DefaultHoursPerMonth, card.HoursPerMonth)
	}
}

func TestCloudProviderPricingCards(t *testing.T) {
	aws := AWSPricingCard()
	if aws.Provider == "" || aws.CPURatePerHour <= 0 || aws.MemoryRatePerGBHour <= 0 {
		t.Errorf("invalid AWS pricing card: %+v", aws)
	}

	gcp := GCPPricingCard()
	if gcp.Provider == "" || gcp.CPURatePerHour <= 0 || gcp.MemoryRatePerGBHour <= 0 {
		t.Errorf("invalid GCP pricing card: %+v", gcp)
	}

	azure := AzurePricingCard()
	if azure.Provider == "" || azure.CPURatePerHour <= 0 || azure.MemoryRatePerGBHour <= 0 {
		t.Errorf("invalid Azure pricing card: %+v", azure)
	}
}

func TestLoadPricingCard(t *testing.T) {
	tmpDir := t.TempDir()
	validJSON := `{
		"provider": "CustomDC",
		"currency": "EUR",
		"cpuRatePerHour": 0.05,
		"memoryRatePerGbHour": 0.008,
		"hoursPerMonth": 720
	}`
	validFile := filepath.Join(tmpDir, "custom.json")
	if err := os.WriteFile(validFile, []byte(validJSON), 0600); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	card, err := LoadPricingCard(validFile)
	if err != nil {
		t.Fatalf("unexpected error loading pricing card: %v", err)
	}

	if card.Provider != "CustomDC" {
		t.Errorf("expected provider CustomDC, got %s", card.Provider)
	}
	if card.Currency != "EUR" {
		t.Errorf("expected currency EUR, got %s", card.Currency)
	}
	if card.CPURatePerHour != 0.05 {
		t.Errorf("expected cpu rate 0.05, got %f", card.CPURatePerHour)
	}
	if card.HoursPerMonth != 720 {
		t.Errorf("expected 720 hours, got %f", card.HoursPerMonth)
	}

	// Non-existent file
	_, err = LoadPricingCard(filepath.Join(tmpDir, "non-existent.json"))
	if err == nil {
		t.Error("expected error loading non-existent file, got nil")
	}

	// Malformed JSON
	badFile := filepath.Join(tmpDir, "bad.json")
	_ = os.WriteFile(badFile, []byte("not-json"), 0600)
	_, err = LoadPricingCard(badFile)
	if err == nil {
		t.Error("expected error loading malformed json, got nil")
	}
}
