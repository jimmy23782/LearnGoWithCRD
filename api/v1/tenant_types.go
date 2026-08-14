/*
Copyright 2026.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// TenantSpec defines the desired state of a Tenant.
type TenantSpec struct {

	// Human-readable name of the tenant team.
	// +kubebuilder:validation:MinLength=1
	TeamName string `json:"teamName"`

	// Unique team identifier from the CMDB.
	// +kubebuilder:validation:MinLength=1
	CMDBTeamID string `json:"cmdbTeamId"`

	// Team owners responsible for this tenant.
	// +kubebuilder:validation:MinItems=1
	Owners []string `json:"owners"`

	// Target environment.
	// +kubebuilder:validation:Enum=dev;uat;prod
	Environment string `json:"environment"`
}

// TenantStatus defines the observed state of Tenant.
type TenantStatus struct {

	// Current reconciliation status.
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:subresource:status

// Tenant is the Schema for the tenants API.
type Tenant struct {
	metav1.TypeMeta `json:",inline"`

	// Standard Kubernetes metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Desired state.
	Spec TenantSpec `json:"spec"`

	// Current observed state.
	Status TenantStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TenantList contains a list of Tenant.
type TenantList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Tenant `json:"items"`
}

func init() {
	SchemeBuilder.Register(func(s *runtime.Scheme) error {
		s.AddKnownTypes(SchemeGroupVersion, &Tenant{}, &TenantList{})
		return nil
	})
}
