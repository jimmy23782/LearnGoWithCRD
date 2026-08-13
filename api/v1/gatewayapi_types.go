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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// GatewayAPISpec defines the desired state of GatewayAPI
type GatewayAPISpec struct {
	// Name of the tenant team requesting the gateway
	// +required
	TeamName string `json:"teamName"`

	// Name of the API being exposed
	// +required
	APIName string `json:"apiName"`

	// Backend configuration
	// +required
	Backend BackendSpec `json:"backend"`

	// Authentication configuration
	// +required
	Authentication AuthenticationSpec `json:"authentication"`

	// Optional list of routes
	// +optional
	Routes []RouteSpec `json:"routes,omitempty"`
}

type BackendSpec struct {
	// Backend URL
	// +required
	URL string `json:"url"`
}

type AuthenticationSpec struct {
	// Authentication type
	// +kubebuilder:validation:Enum=oidc;jwt;basic;none
	// +required
	Type string `json:"type"`
}

type RouteSpec struct {
	// Route path
	// +required
	Path string `json:"path"`

	// Optional HTTP methods
	// +optional
	Methods []string `json:"methods,omitempty"`
}

// GatewayAPIStatus defines the observed state of GatewayAPI.
type GatewayAPIStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// For Kubernetes API conventions, see:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	// conditions represent the current state of the GatewayAPI resource.
	// Each condition has a unique type and reflects the status of a specific aspect of the resource.
	//
	// Standard condition types include:
	// - "Available": the resource is fully functional
	// - "Progressing": the resource is being created or updated
	// - "Degraded": the resource failed to reach or maintain its desired state
	//
	// The status of each condition is one of True, False, or Unknown.
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// GatewayAPI is the Schema for the gatewayapis API
type GatewayAPI struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of GatewayAPI
	// +required
	Spec GatewayAPISpec `json:"spec"`

	// status defines the observed state of GatewayAPI
	// +optional
	Status GatewayAPIStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// GatewayAPIList contains a list of GatewayAPI
type GatewayAPIList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []GatewayAPI `json:"items"`
}

func init() {
	SchemeBuilder.Register(func(s *runtime.Scheme) error {
		s.AddKnownTypes(SchemeGroupVersion, &GatewayAPI{}, &GatewayAPIList{})
		return nil
	})
}
