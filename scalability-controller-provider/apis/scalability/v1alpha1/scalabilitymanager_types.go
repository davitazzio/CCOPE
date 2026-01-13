/*
Copyright 2022 The Crossplane Authors.

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

package v1alpha1

import (
	"reflect"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
)

// ScalabilityManagerParameters are the configurable fields of a ScalabilityManager.
type ScalabilityManagerParameters struct {
	Trashold string `json:"trashold"`
}

// ScalabilityManagerObservation are the observable fields of a ScalabilityManager.
type ScalabilityManagerObservation struct {
	NumberOfProcesses string `json:"numberOfProcesses,omitempty"`
}

// A ScalabilityManagerSpec defines the desired state of a ScalabilityManager.
type ScalabilityManagerSpec struct {
	xpv1.ResourceSpec  `json:",inline"`
	NumConsumerDesired int                          `json:"numConsumerDesired,omitempty"`
	ForProvider        ScalabilityManagerParameters `json:"forProvider"`
}

// A ScalabilityManagerStatus represents the observed state of a ScalabilityManager.
type ScalabilityManagerStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	Consumers           *[]string                     `json:"Consumers,omitempty"`
	AtProvider          ScalabilityManagerObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A ScalabilityManager is an example API type.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,scalabilityprovider}
type ScalabilityManager struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ScalabilityManagerSpec   `json:"spec"`
	Status ScalabilityManagerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ScalabilityManagerList contains a list of ScalabilityManager
type ScalabilityManagerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ScalabilityManager `json:"items"`
}

// ScalabilityManager type metadata.
var (
	ScalabilityManagerKind             = reflect.TypeOf(ScalabilityManager{}).Name()
	ScalabilityManagerGroupKind        = schema.GroupKind{Group: Group, Kind: ScalabilityManagerKind}.String()
	ScalabilityManagerKindAPIVersion   = ScalabilityManagerKind + "." + SchemeGroupVersion.String()
	ScalabilityManagerGroupVersionKind = SchemeGroupVersion.WithKind(ScalabilityManagerKind)
)

func init() {
	SchemeBuilder.Register(&ScalabilityManager{}, &ScalabilityManagerList{})
}
