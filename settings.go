package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/kubewarden/container-resources-policy/resource"
	kubewarden "github.com/kubewarden/policy-sdk-go"
	kubewarden_protocol "github.com/kubewarden/policy-sdk-go/protocol"
)

type Settings struct {
	MaxCpuLimit            string
	MaxMemoryLimit         string
	CpuLimitMaxQuantity    resource.Quantity `json:"-"`
	MemoryLimitMaxQuantity resource.Quantity `json:"-"`

	MaxCpuRequest            string
	MaxMemoryRequest         string
	CpuRequestMaxQuantity    resource.Quantity `json:"-"`
	MemoryRequestMaxQuantity resource.Quantity `json:"-"`
	IgnoreImages             []string
}

// No special checks have to be done
func (s *Settings) Valid() (bool, error) {
	if len(strings.TrimSpace(s.MaxCpuLimit)) == 0 && len(strings.TrimSpace(s.MaxMemoryLimit)) == 0 && len(strings.TrimSpace(s.MaxCpuRequest)) == 0 && len(strings.TrimSpace(s.MaxMemoryRequest)) == 0 {
		return false, fmt.Errorf("no settings provided. At least one resource limit or request must be verified")
	}
	if _, err := resource.ParseQuantity(s.MaxCpuLimit); err != nil {
		return false, errors.Join(fmt.Errorf("failed to parse cpu limit quantity"), err)
	}
	if _, err := resource.ParseQuantity(s.MaxMemoryLimit); err != nil {
		return false, errors.Join(fmt.Errorf("failed to parse memory limit quantity"), err)
	}
	if _, err := resource.ParseQuantity(s.MaxCpuRequest); err != nil {
		return false, errors.Join(fmt.Errorf("failed to parse cpu request quantity"), err)
	}
	if _, err := resource.ParseQuantity(s.MaxMemoryRequest); err != nil {
		return false, errors.Join(fmt.Errorf("failed to parse memory request quantity"), err)
	}
	return true, nil
}

func (s *Settings) parseAllLimits() {
	if len(strings.TrimSpace(s.MaxCpuLimit)) != 0 {
		s.CpuLimitMaxQuantity = resource.MustParse(s.MaxCpuLimit)
	}
	if len(strings.TrimSpace(s.MaxMemoryLimit)) != 0 {
		s.MemoryLimitMaxQuantity = resource.MustParse(s.MaxMemoryLimit)
	}
	if len(strings.TrimSpace(s.MaxCpuRequest)) != 0 {
		s.CpuRequestMaxQuantity = resource.MustParse(s.MaxCpuRequest)
	}
	if len(strings.TrimSpace(s.MaxMemoryRequest)) != 0 {
		s.MemoryRequestMaxQuantity = resource.MustParse(s.MaxMemoryRequest)
	}
}

func NewSettingsFromValidationReq(validationReq *kubewarden_protocol.ValidationRequest) (Settings, error) {
	settings := Settings{}
	err := json.Unmarshal(validationReq.Settings, &settings)
	if err == nil {
		settings.parseAllLimits()
	}
	return settings, err
}

func validateSettings(payload []byte) ([]byte, error) {
	logger.Info("validating settings")
	settings := Settings{}
	err := json.Unmarshal(payload, &settings)
	if err != nil {
		return kubewarden.RejectSettings(kubewarden.Message(fmt.Sprintf("Provided settings are not valid: %v", err)))
	}

	valid, err := settings.Valid()
	if err != nil {
		return kubewarden.RejectSettings(kubewarden.Message(fmt.Sprintf("Provided settings are not valid: %v", err)))
	}
	if valid {
		return kubewarden.AcceptSettings()
	}

	logger.Warn("rejecting settings")
	return kubewarden.RejectSettings(kubewarden.Message("Provided settings are not valid"))
}
