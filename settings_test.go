package main

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/kubewarden/container-resources-policy/resource"
	kubewarden_protocol "github.com/kubewarden/policy-sdk-go/protocol"
)

func TestParsingSettingsWithNoValueProvided(t *testing.T) {
	rawSettings := []byte(`{}`)
	settings := &Settings{}
	if err := json.Unmarshal(rawSettings, settings); err != nil {
		t.Fatalf("Unexpected error %+v", err)
	}
	valid, _ := settings.Valid()
	if valid {
		t.Fatal("settings reported as valid. User must set some required fields")
	}
}

func TestParsingSettings(t *testing.T) {
	var tests = []struct {
		name         string
		rawSettings  []byte
		errorMessage string
	}{
		{"valid settings", []byte(`{"maxCpuLimit": "1m", "maxMemoryLimit": "1G", "maxCpuRequest": "1m", "maxMemoryRequest": "1G", "ignoreImages": ["image:latest"]}`), ""},
		{"no suffix", []byte(`{"maxCpuLimit": "1", "maxMemoryLimit": "2","maxCpuRequest": "1", "maxMemoryRequest": "2"}`), ""},
		{"invalid cpu limit suffix", []byte(`{"maxCpuLimit": "1x", "maxMemoryLimit": "2", "maxCpuRequest": "1m", "maxMemoryRequest": "1G"}`), "failed to parse cpu limit quantity"},
		{"invalid memory limit suffix", []byte(`{"maxCpuLimit": "1", "maxMemoryLimit": "2x", "maxCpuRequest": "1m", "maxMemoryRequest": "1G"}`), "failed to parse memory limit quantity"},
		{"invalid cpu request suffix", []byte(`{"maxCpuLimit": "1", "maxMemoryLimit": "2", "maxCpuRequest": "1x", "maxMemoryRequest": "2G"}`), "failed to parse cpu request quantity"},
		{"invalid memory request suffix", []byte(`{"maxCpuLimit": "1", "maxMemoryLimit": "2", "maxCpuRequest": "1m", "maxMemoryRequest": "2x"}`), "failed to parse memory request quantity"},
		{"invalid settings", []byte(`{}`), "no settings provided. At least one resource limit or request must be verified"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			settings := &Settings{}
			if err := json.Unmarshal(test.rawSettings, settings); err != nil {
				t.Fatalf("Unexpected error %+v", err)
			}
			valid, err := settings.Valid()
			if len(test.errorMessage) == 0 && err != nil {
				t.Fatalf("unexpected validation error: %+v", err)
			}
			if len(test.errorMessage) > 0 && !strings.Contains(err.Error(), test.errorMessage) {
				t.Errorf("invalid error message. Expected the string '%s' in the error. Got '%s'", test.errorMessage, err.Error())
			}
			if !valid && err == nil {
				t.Fatal("settings reported as invalid but not error is returned")
			}
			if !settings.CpuLimitMaxQuantity.IsZero() {
				t.Error("validation should not parse CPU limits")
			}
			if !settings.MemoryLimitMaxQuantity.IsZero() {
				t.Error("validation should not parse CPU limits")
			}
			if !settings.CpuRequestMaxQuantity.IsZero() {
				t.Error("validation should not parse CPU limits")
			}
			if !settings.MemoryRequestMaxQuantity.IsZero() {
				t.Error("validation should not parse CPU limits")
			}
		})
	}

}

func TestNewSettingsFromValidationReq(t *testing.T) {
	validationReq := &kubewarden_protocol.ValidationRequest{
		Settings: []byte(`{"maxCpuLimit": "1m", "maxMemoryLimit": "1G", "maxCpuRequest": "1m", "maxMemoryRequest": "1G"}`),
	}
	settings, err := NewSettingsFromValidationReq(validationReq)
	if err != nil {
		t.Fatalf("Unexpected error %+v", err)
	}
	expectedCpuValue := resource.MustParse("1m")
	if !settings.CpuLimitMaxQuantity.Equal(expectedCpuValue) {
		t.Errorf("invalid memory limit quantity parsed. Expected %+v, go %+v", expectedCpuValue, settings.CpuLimitMaxQuantity)
	}
	if !settings.CpuRequestMaxQuantity.Equal(expectedCpuValue) {
		t.Errorf("invalid cpu request quantity parsed. Expected %+v, go %+v", expectedCpuValue, settings.CpuLimitMaxQuantity)
	}

	expectedMemoryValue := resource.MustParse("1G")
	if !settings.MemoryLimitMaxQuantity.Equal(expectedMemoryValue) {
		t.Errorf("invalid memory limit quantity parsed. Expected %+v, go %+v", expectedMemoryValue, settings.MemoryLimitMaxQuantity)
	}
	if !settings.MemoryRequestMaxQuantity.Equal(expectedMemoryValue) {
		t.Errorf("invalid cpu request quantity parsed. Expected %+v, go %+v", expectedMemoryValue, settings.MemoryLimitMaxQuantity)
	}
}
