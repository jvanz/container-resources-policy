package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kubewarden/container-resources-policy/resource"
	corev1 "github.com/kubewarden/k8s-objects/api/core/v1"
	api_resource "github.com/kubewarden/k8s-objects/apimachinery/pkg/api/resource"
	kubewarden "github.com/kubewarden/policy-sdk-go"
	kubewarden_protocol "github.com/kubewarden/policy-sdk-go/protocol"
)

func validateAndAdjustContainerResourceRequests(container *corev1.Container, settings *Settings) (bool, error) {
	mutated := false
	if !settings.MemoryRequestMaxQuantity.IsZero() {
		memoryStr, ok := container.Resources.Requests["memory"]
		if !ok || len(strings.TrimSpace(string(*memoryStr))) == 0 {
			newLimit := api_resource.Quantity(settings.MemoryRequestMaxQuantity.String())
			container.Resources.Requests["memory"] = &newLimit
			mutated = true
		} else {
			memoryLimit, err := resource.ParseQuantity(string(*memoryStr))
			if err != nil {
				return false, fmt.Errorf("invalid memory limit")
			}
			if memoryLimit.Cmp(settings.MemoryRequestMaxQuantity) > 0 {
				newLimit := api_resource.Quantity(settings.MemoryRequestMaxQuantity.String())
				container.Resources.Requests["memory"] = &newLimit
				mutated = true
			}
		}
	}

	if !settings.CpuRequestMaxQuantity.IsZero() {
		cpuStr, ok := container.Resources.Requests["cpu"]
		if !ok || len(strings.TrimSpace(string(*cpuStr))) == 0 {
			newLimit := api_resource.Quantity(settings.CpuRequestMaxQuantity.String())
			container.Resources.Requests["cpu"] = &newLimit
			mutated = true
		} else {
			cpuLimit, err := resource.ParseQuantity(string(*cpuStr))
			if err != nil {
				return false, fmt.Errorf("invalid cpu limit")
			}
			if cpuLimit.Cmp(settings.CpuRequestMaxQuantity) > 0 {
				newLimit := api_resource.Quantity(settings.CpuRequestMaxQuantity.String())
				container.Resources.Requests["cpu"] = &newLimit
				mutated = true
			}
		}
	}
	return mutated, nil
}

func validateAndAdjustContainerResourceLimits(container *corev1.Container, settings *Settings) (bool, error) {
	mutated := false
	if !settings.MemoryLimitMaxQuantity.IsZero() {
		memoryStr, ok := container.Resources.Limits["memory"]
		if !ok || len(strings.TrimSpace(string(*memoryStr))) == 0 {
			newLimit := api_resource.Quantity(settings.MemoryLimitMaxQuantity.String())
			container.Resources.Limits["memory"] = &newLimit
			mutated = true
		} else {
			memoryLimit, err := resource.ParseQuantity(string(*memoryStr))
			if err != nil {
				return false, fmt.Errorf("invalid memory limit")
			}
			if memoryLimit.Cmp(settings.MemoryLimitMaxQuantity) > 0 {
				newLimit := api_resource.Quantity(settings.MemoryLimitMaxQuantity.String())
				container.Resources.Limits["memory"] = &newLimit
				mutated = true
			}
		}
	}

	if !settings.CpuLimitMaxQuantity.IsZero() {
		cpuStr, ok := container.Resources.Limits["cpu"]
		if !ok || len(strings.TrimSpace(string(*cpuStr))) == 0 {
			newLimit := api_resource.Quantity(settings.CpuLimitMaxQuantity.String())
			container.Resources.Limits["cpu"] = &newLimit
			mutated = true
		} else {
			cpuLimit, err := resource.ParseQuantity(string(*cpuStr))
			if err != nil {
				return false, fmt.Errorf("invalid cpu limit")
			}
			if cpuLimit.Cmp(settings.CpuLimitMaxQuantity) > 0 {
				newLimit := api_resource.Quantity(settings.CpuLimitMaxQuantity.String())
				container.Resources.Limits["cpu"] = &newLimit
				mutated = true
			}
		}
	}
	return mutated, nil
}

func validateAndAdjustContainer(container *corev1.Container, settings *Settings) (bool, error) {
	if container.Resources == nil {
		container.Resources = &corev1.ResourceRequirements{
			Limits:   make(map[string]*api_resource.Quantity),
			Requests: make(map[string]*api_resource.Quantity),
		}
	}
	limitsMutation, err := validateAndAdjustContainerResourceLimits(container, settings)
	if err != nil {
		return false, err
	}
	requestsMutation, err := validateAndAdjustContainerResourceRequests(container, settings)
	if err != nil {
		return false, err
	}
	return limitsMutation || requestsMutation, nil

}

func shouldValidateContainer(image string, ignoreImages []string) bool {
	for _, ignoreImageName := range ignoreImages {
		if image == ignoreImageName {
			return false
		}
	}
	return true
}

func validatePodSpec(pod *corev1.PodSpec, settings *Settings) (bool, error) {
	mutated := false
	for _, container := range pod.Containers {
		if shouldValidateContainer(container.Image, settings.IgnoreImages) {
			containerMutated, err := validateAndAdjustContainer(container, settings)
			if err != nil {
				return false, err
			}
			mutated = mutated || containerMutated
		}
	}
	return mutated, nil

}

func validate(payload []byte) ([]byte, error) {
	// Create a ValidationRequest instance from the incoming payload
	validationRequest := kubewarden_protocol.ValidationRequest{}
	err := json.Unmarshal(payload, &validationRequest)
	if err != nil {
		return kubewarden.RejectRequest(
			kubewarden.Message(err.Error()),
			kubewarden.Code(400))
	}

	// Create a Settings instance from the ValidationRequest object
	settings, err := NewSettingsFromValidationReq(&validationRequest)
	if err != nil {
		return kubewarden.RejectRequest(
			kubewarden.Message(err.Error()),
			kubewarden.Code(400))
	}

	pod, err := kubewarden.ExtractPodSpecFromObject(validationRequest)
	if err == nil {
		mutatePod, err := validatePodSpec(&pod, &settings)
		if err != nil {
			return kubewarden.RejectRequest(
				kubewarden.Message(err.Error()),
				kubewarden.Code(400))
		}
		if mutatePod {
			return kubewarden.MutatePodSpecFromRequest(validationRequest, pod)
		}
	} else {
		return kubewarden.RejectRequest(kubewarden.Message(err.Error()), kubewarden.Code(400))
	}

	return kubewarden.AcceptRequest()
}
