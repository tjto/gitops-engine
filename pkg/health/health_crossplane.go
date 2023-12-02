package health

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func getCrossplaneProviderHealth(obj *unstructured.Unstructured) (*HealthStatus, error) {
	status, ok := obj.Object["status"].(map[string]interface{})
	if !ok {
		return &HealthStatus{
			Status:  HealthStatusUnknown,
			Message: fmt.Sprintf("provider %s/%s doesn't have status yet", obj.GetKind(), obj.GetName()),
		}, nil
	}
	for _, condition := range status["conditions"].([]interface{}) {
		c, ok := condition.(map[string]interface{})
		if !ok {
			return &HealthStatus{
				Status:  HealthStatusUnknown,
				Message: fmt.Sprintf("cannot get conditions of provider %s/%s", obj.GetKind(), obj.GetName()),
			}, nil
		}
		if c["type"] == "Healthy" && c["status"] == "True" && c["reason"] == "HealthyPackageRevision" {
			return &HealthStatus{
				Status:  HealthStatusHealthy,
				Message: fmt.Sprintf("provider %s/%s is healthy", obj.GetKind(), obj.GetName()),
			}, nil
		}
	}
	return &HealthStatus{
		Status:  HealthStatusProgressing,
		Message: fmt.Sprintf("provider %s: %v", obj.GetName(), status),
	}, nil
}

func getCrossplaneResourceHealth(obj *unstructured.Unstructured) (*HealthStatus, error) {
	status, ok := obj.Object["status"].(map[string]interface{})
	if !ok {
		return &HealthStatus{
			Status:  HealthStatusUnknown,
			Message: fmt.Sprintf("resource %s/%s doesn't have status yet", obj.GetKind(), obj.GetName()),
		}, nil
	}
	for _, condition := range status["conditions"].([]interface{}) {
		c, ok := condition.(map[string]interface{})
		if !ok {
			return &HealthStatus{
				Status:  HealthStatusUnknown,
				Message: fmt.Sprintf("cannot get condition of resource %s/%s", obj.GetKind(), obj.GetName()),
			}, nil
		}
		if c["type"] == "Ready" && c["status"] == "True" {
			return &HealthStatus{
				Status:  HealthStatusHealthy,
				Message: fmt.Sprintf("resource %s/%s is healthy", obj.GetKind(), obj.GetName()),
			}, nil
		}
	}

	return &HealthStatus{
		Status:  HealthStatusProgressing,
		Message: fmt.Sprintf("resource %s: %v", obj.GetName(), status),
	}, nil
}
