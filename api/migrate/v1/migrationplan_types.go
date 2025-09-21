/*
Copyright 2025.

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
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MigrationPlanSpec defines the desired state of MigrationPlan.
type MigrationPlanSpec struct {

	// Namespaces in the source OpenShift cluster to scan.
	// If empty, default to the CR's namespace.
	Namespaces []string `json:"namespaces,omitempty"`

	// Which resource kinds to include (case-insensitive):
	// e.g., ["deploymentconfigs","routes","services","pvcs"]
	Include []string `json:"include,omitempty"`

	// Target cloud to guide small deltas in output:
	// "eks" | "aks" | "gke" | "vanilla"
	TargetCloud string `json:"targetCloud,omitempty"`

	// Ingress class hint for Route -> Ingress conversion
	// e.g., "alb" (EKS), "azure/application-gateway" (AKS), "gce" (GKE), "nginx"
	IngressClass string `json:"ingressClass,omitempty"`

	// Name of the ConfigMap (in the same namespace) where the operator
	// will write the generated manifests as data["converted.yaml"].
	// If empty, defaults to "<CR name>-output".
	OutputConfigMap string `json:"outputConfigMap,omitempty"`
}

// MigrationPlanStatus defines the observed state of MigrationPlan.
type MigrationPlanStatus struct {

	// Phase of processing: "Scanning" | "Generated" | "Error"
	Phase string `json:"phase,omitempty"`

	// Counts of discovered resources by kind, e.g.:
	// { "deploymentconfigs": 1, "routes": 1, "services": 1, "persistentvolumeclaims": 1 }
	Found map[string]int `json:"found,omitempty"`

	// Notes / warnings for the user (e.g., "PVC RWO -> EBS (gp3)")
	Notes []string `json:"notes,omitempty"`

	// The actual ConfigMap name used for output (resolved from spec.OutputConfigMap)
	Output string `json:"output,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// MigrationPlan is the Schema for the migrationplans API.
type MigrationPlan struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MigrationPlanSpec   `json:"spec,omitempty"`
	Status MigrationPlanStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MigrationPlanList contains a list of MigrationPlan.
type MigrationPlanList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MigrationPlan `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MigrationPlan{}, &MigrationPlanList{})
}
