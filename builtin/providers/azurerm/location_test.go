package azurerm

import "testing"

func TestAzureRMNormalizeLocation_westUSFull(t *testing.T) {
	s := azureRMNormalizeLocation("West US")
	if s != "westus" {
		t.Fatalf("expected location to equal westus, actual %s", s)
	}
}

func TestAzureRMNormalizeLocation_westUSShort(t *testing.T) {
	s := azureRMNormalizeLocation("westus")
	if s != "westus" {
		t.Fatalf("expected location to equal westus, actual %s", s)
	}
}

func TestAzureRMNormalizeLocation_southEastAsiaFull(t *testing.T) {
	s := azureRMNormalizeLocation("South East Asia")
	if s != "southeastasia" {
		t.Fatalf("expected location to equal southeastasia, actual %s", s)
	}
}

func TestAzureRMNormalizeLocation_southEastAsiaShort(t *testing.T) {
	s := azureRMNormalizeLocation("southeastasia")
	if s != "southeastasia" {
		t.Fatalf("expected location to equal southeastasia, actual %s", s)
	}
}
