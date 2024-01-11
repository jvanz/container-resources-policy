package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kubewarden/container-resources-policy/resource"
	corev1 "github.com/kubewarden/k8s-objects/api/core/v1"
	apimachinery_pkg_api_resource "github.com/kubewarden/k8s-objects/apimachinery/pkg/api/resource"
)

func TestContainerIsRequiredToHaveLimits(t *testing.T) {
	oneCore := resource.MustParse("1")
	oneGi := resource.MustParse("1Gi")
	oneCoreCpuQuantity := apimachinery_pkg_api_resource.Quantity("1")
	oneGiMemoryQuantity := apimachinery_pkg_api_resource.Quantity("1Gi")
	twoCore := resource.MustParse("1")
	twoGi := resource.MustParse("2Gi")
	twoCoreCpuQuantity := apimachinery_pkg_api_resource.Quantity("2")
	twoGiMemoryQuantity := apimachinery_pkg_api_resource.Quantity("2Gi")
	var tests = []struct {
		name                  string
		container             corev1.Container
		settings              Settings
		expectedResouceLimits *corev1.ResourceRequirements
		shouldMutate          bool
	}{
		{"no resources and limits defined", corev1.Container{}, Settings{CpuLimitMaxQuantity: oneCore, MemoryLimitMaxQuantity: oneGi, CpuRequestMaxQuantity: oneCore, MemoryRequestMaxQuantity: oneGi}, &corev1.ResourceRequirements{
			Limits: map[string]*apimachinery_pkg_api_resource.Quantity{
				"cpu":    &oneCoreCpuQuantity,
				"memory": &oneGiMemoryQuantity,
			},
			Requests: map[string]*apimachinery_pkg_api_resource.Quantity{
				"cpu":    &oneCoreCpuQuantity,
				"memory": &oneGiMemoryQuantity,
			},
		}, true},
		{"no memory limit", corev1.Container{
			Resources: &corev1.ResourceRequirements{
				Limits: map[string]*apimachinery_pkg_api_resource.Quantity{
					"cpu": &oneCoreCpuQuantity,
				},
				Requests: make(map[string]*apimachinery_pkg_api_resource.Quantity),
			},
		}, Settings{CpuLimitMaxQuantity: oneCore, MemoryLimitMaxQuantity: oneGi}, &corev1.ResourceRequirements{
			Limits: map[string]*apimachinery_pkg_api_resource.Quantity{
				"cpu":    &oneCoreCpuQuantity,
				"memory": &oneGiMemoryQuantity,
			},
			Requests: make(map[string]*apimachinery_pkg_api_resource.Quantity),
		}, true},
		{"no cpu limit", corev1.Container{
			Resources: &corev1.ResourceRequirements{
				Limits: map[string]*apimachinery_pkg_api_resource.Quantity{
					"memory": &oneGiMemoryQuantity,
				},
				Requests: make(map[string]*apimachinery_pkg_api_resource.Quantity),
			},
		}, Settings{CpuLimitMaxQuantity: oneCore, MemoryLimitMaxQuantity: oneGi}, &corev1.ResourceRequirements{
			Limits: map[string]*apimachinery_pkg_api_resource.Quantity{
				"cpu":    &oneCoreCpuQuantity,
				"memory": &oneGiMemoryQuantity,
			},
			Requests: make(map[string]*apimachinery_pkg_api_resource.Quantity),
		}, true},
		{"all limits within the expected range", corev1.Container{
			Resources: &corev1.ResourceRequirements{
				Limits: map[string]*apimachinery_pkg_api_resource.Quantity{
					"cpu":    &oneCoreCpuQuantity,
					"memory": &oneGiMemoryQuantity,
				},
				Requests: make(map[string]*apimachinery_pkg_api_resource.Quantity),
			},
		}, Settings{CpuLimitMaxQuantity: twoCore, MemoryLimitMaxQuantity: twoGi}, &corev1.ResourceRequirements{
			Limits: map[string]*apimachinery_pkg_api_resource.Quantity{
				"cpu":    &oneCoreCpuQuantity,
				"memory": &oneGiMemoryQuantity,
			},
			Requests: make(map[string]*apimachinery_pkg_api_resource.Quantity),
		}, false},
		{"all limits exceeding the expected range", corev1.Container{
			Resources: &corev1.ResourceRequirements{
				Limits: map[string]*apimachinery_pkg_api_resource.Quantity{
					"cpu":    &twoCoreCpuQuantity,
					"memory": &twoGiMemoryQuantity,
				},
				Requests: make(map[string]*apimachinery_pkg_api_resource.Quantity),
			},
		}, Settings{CpuLimitMaxQuantity: oneCore, MemoryLimitMaxQuantity: oneGi}, &corev1.ResourceRequirements{
			Limits: map[string]*apimachinery_pkg_api_resource.Quantity{
				"cpu":    &oneCoreCpuQuantity,
				"memory": &oneGiMemoryQuantity,
			},
			Requests: make(map[string]*apimachinery_pkg_api_resource.Quantity),
		}, true},

		{"no memory resource request", corev1.Container{
			Resources: &corev1.ResourceRequirements{
				Requests: map[string]*apimachinery_pkg_api_resource.Quantity{
					"cpu": &oneCoreCpuQuantity,
				},
				Limits: make(map[string]*apimachinery_pkg_api_resource.Quantity),
			},
		}, Settings{CpuRequestMaxQuantity: oneCore, MemoryRequestMaxQuantity: oneGi}, &corev1.ResourceRequirements{
			Requests: map[string]*apimachinery_pkg_api_resource.Quantity{
				"cpu":    &oneCoreCpuQuantity,
				"memory": &oneGiMemoryQuantity,
			},
			Limits: make(map[string]*apimachinery_pkg_api_resource.Quantity),
		}, true},
		{"no cpu resource request", corev1.Container{
			Resources: &corev1.ResourceRequirements{
				Requests: map[string]*apimachinery_pkg_api_resource.Quantity{
					"memory": &oneGiMemoryQuantity,
				},
				Limits: make(map[string]*apimachinery_pkg_api_resource.Quantity),
			},
		}, Settings{CpuRequestMaxQuantity: oneCore, MemoryRequestMaxQuantity: oneGi}, &corev1.ResourceRequirements{
			Requests: map[string]*apimachinery_pkg_api_resource.Quantity{
				"cpu":    &oneCoreCpuQuantity,
				"memory": &oneGiMemoryQuantity,
			},
			Limits: make(map[string]*apimachinery_pkg_api_resource.Quantity),
		}, true},
		{"all resource requests within the expected range", corev1.Container{
			Resources: &corev1.ResourceRequirements{
				Requests: map[string]*apimachinery_pkg_api_resource.Quantity{
					"cpu":    &oneCoreCpuQuantity,
					"memory": &oneGiMemoryQuantity,
				},
			},
		}, Settings{CpuRequestMaxQuantity: oneCore, MemoryRequestMaxQuantity: oneGi}, &corev1.ResourceRequirements{
			Requests: map[string]*apimachinery_pkg_api_resource.Quantity{
				"cpu":    &oneCoreCpuQuantity,
				"memory": &oneGiMemoryQuantity,
			},
		}, false},
		{"all resource requests exceeding the expected range", corev1.Container{
			Resources: &corev1.ResourceRequirements{
				Requests: map[string]*apimachinery_pkg_api_resource.Quantity{
					"cpu":    &twoCoreCpuQuantity,
					"memory": &twoGiMemoryQuantity,
				},
				Limits: make(map[string]*apimachinery_pkg_api_resource.Quantity),
			},
		}, Settings{CpuRequestMaxQuantity: oneCore, MemoryRequestMaxQuantity: oneGi}, &corev1.ResourceRequirements{
			Requests: map[string]*apimachinery_pkg_api_resource.Quantity{
				"cpu":    &oneCoreCpuQuantity,
				"memory": &oneGiMemoryQuantity,
			},
			Limits: make(map[string]*apimachinery_pkg_api_resource.Quantity),
		}, true},

		{"all resource configuration within the expected range", corev1.Container{
			Resources: &corev1.ResourceRequirements{
				Requests: map[string]*apimachinery_pkg_api_resource.Quantity{
					"cpu":    &oneCoreCpuQuantity,
					"memory": &oneGiMemoryQuantity,
				},
				Limits: map[string]*apimachinery_pkg_api_resource.Quantity{
					"cpu":    &oneCoreCpuQuantity,
					"memory": &oneGiMemoryQuantity,
				},
			},
		}, Settings{CpuLimitMaxQuantity: oneCore, MemoryLimitMaxQuantity: oneGi, CpuRequestMaxQuantity: oneCore, MemoryRequestMaxQuantity: oneGi}, &corev1.ResourceRequirements{
			Requests: map[string]*apimachinery_pkg_api_resource.Quantity{
				"cpu":    &oneCoreCpuQuantity,
				"memory": &oneGiMemoryQuantity,
			},
			Limits: map[string]*apimachinery_pkg_api_resource.Quantity{
				"cpu":    &oneCoreCpuQuantity,
				"memory": &oneGiMemoryQuantity,
			},
		}, false},

		{"all resource configuration exceeding the expected range", corev1.Container{
			Resources: &corev1.ResourceRequirements{
				Requests: map[string]*apimachinery_pkg_api_resource.Quantity{
					"cpu":    &twoCoreCpuQuantity,
					"memory": &twoGiMemoryQuantity,
				},
				Limits: map[string]*apimachinery_pkg_api_resource.Quantity{
					"cpu":    &twoCoreCpuQuantity,
					"memory": &twoGiMemoryQuantity,
				},
			},
		}, Settings{CpuLimitMaxQuantity: oneCore, MemoryLimitMaxQuantity: oneGi, CpuRequestMaxQuantity: oneCore, MemoryRequestMaxQuantity: oneGi}, &corev1.ResourceRequirements{
			Requests: map[string]*apimachinery_pkg_api_resource.Quantity{
				"cpu":    &oneCoreCpuQuantity,
				"memory": &oneGiMemoryQuantity,
			},
			Limits: map[string]*apimachinery_pkg_api_resource.Quantity{
				"cpu":    &oneCoreCpuQuantity,
				"memory": &oneGiMemoryQuantity,
			},
		}, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mutated, err := validateAndAdjustContainer(&test.container, &test.settings)
			if err != nil {
				t.Fatalf("unexpected error: %q", err)
			}
			if diff := cmp.Diff(test.container.Resources, test.expectedResouceLimits); diff != "" {
				t.Error(diff)
			}
			if mutated != test.shouldMutate {
				t.Errorf("validation function does not report mutation flag correctly. Got: %t, expected: %t", mutated, test.shouldMutate)
			}

		})
	}
}

func TestIgnoreImageSettings(t *testing.T) {
	oneCore := resource.MustParse("1")
	oneGi := resource.MustParse("1Gi")
	oneCoreCpuQuantity := apimachinery_pkg_api_resource.Quantity("1")
	oneGiMemoryQuantity := apimachinery_pkg_api_resource.Quantity("1Gi")
	twoCoreCpuQuantity := apimachinery_pkg_api_resource.Quantity("2")
	twoGiMemoryQuantity := apimachinery_pkg_api_resource.Quantity("2Gi")
	container1 := corev1.Container{
		Image: "image1:latest",
		Resources: &corev1.ResourceRequirements{
			Requests: map[string]*apimachinery_pkg_api_resource.Quantity{
				"cpu":    &twoCoreCpuQuantity,
				"memory": &twoGiMemoryQuantity,
			},
			Limits: map[string]*apimachinery_pkg_api_resource.Quantity{
				"cpu":    &twoCoreCpuQuantity,
				"memory": &twoGiMemoryQuantity,
			},
		},
	}
	container2 := corev1.Container{
		Image: "image2:latest",
		Resources: &corev1.ResourceRequirements{
			Requests: map[string]*apimachinery_pkg_api_resource.Quantity{
				"cpu":    &twoCoreCpuQuantity,
				"memory": &twoGiMemoryQuantity,
			},
			Limits: map[string]*apimachinery_pkg_api_resource.Quantity{
				"cpu":    &twoCoreCpuQuantity,
				"memory": &twoGiMemoryQuantity,
			},
		},
	}
	container3 := corev1.Container{
		Image: "image3:latest",
		Resources: &corev1.ResourceRequirements{
			Requests: map[string]*apimachinery_pkg_api_resource.Quantity{
				"cpu":    &twoCoreCpuQuantity,
				"memory": &twoGiMemoryQuantity,
			},
			Limits: map[string]*apimachinery_pkg_api_resource.Quantity{
				"cpu":    &twoCoreCpuQuantity,
				"memory": &twoGiMemoryQuantity,
			},
		},
	}
	settings := Settings{CpuLimitMaxQuantity: oneCore, MemoryLimitMaxQuantity: oneGi, CpuRequestMaxQuantity: oneCore, MemoryRequestMaxQuantity: oneGi, IgnoreImages: []string{"image1:latest"}}
	podSpec := &corev1.PodSpec{
		Containers: []*corev1.Container{&container1, &container2, &container3},
	}
	mutate, err := validatePodSpec(podSpec, &settings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !mutate {
		t.Error("pod should be mutated")
	}
	expectedPodSpec := &corev1.PodSpec{
		Containers: []*corev1.Container{
			{
				Image: "image1:latest",
				Resources: &corev1.ResourceRequirements{
					Requests: map[string]*apimachinery_pkg_api_resource.Quantity{
						"cpu":    &twoCoreCpuQuantity,
						"memory": &twoGiMemoryQuantity,
					},
					Limits: map[string]*apimachinery_pkg_api_resource.Quantity{
						"cpu":    &twoCoreCpuQuantity,
						"memory": &twoGiMemoryQuantity,
					},
				},
			},
			{
				Image: "image2:latest",
				Resources: &corev1.ResourceRequirements{
					Requests: map[string]*apimachinery_pkg_api_resource.Quantity{
						"cpu":    &oneCoreCpuQuantity,
						"memory": &oneGiMemoryQuantity,
					},
					Limits: map[string]*apimachinery_pkg_api_resource.Quantity{
						"cpu":    &oneCoreCpuQuantity,
						"memory": &oneGiMemoryQuantity,
					},
				},
			},
			{
				Image: "image3:latest",
				Resources: &corev1.ResourceRequirements{
					Requests: map[string]*apimachinery_pkg_api_resource.Quantity{
						"cpu":    &oneCoreCpuQuantity,
						"memory": &oneGiMemoryQuantity,
					},
					Limits: map[string]*apimachinery_pkg_api_resource.Quantity{
						"cpu":    &oneCoreCpuQuantity,
						"memory": &oneGiMemoryQuantity,
					},
				},
			},
		},
	}
	if diff := cmp.Diff(expectedPodSpec, podSpec); diff != "" {
		t.Errorf("invalid pod spec:\n %s", diff)
	}

}
