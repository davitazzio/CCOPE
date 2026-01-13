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

// EmqxProcessParameters are the configurable fields of a EmqxProcess.
type EmqxProcessParameters struct {
	Host         string `json:"host"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	BrokerAPIKey string `json:"brokerAPIKey"`
}

// EmqxProcessObservation are the observable fields of a EmqxProcess.
type EmqxProcessObservation struct {
	ProcessPID int64 `json:"process_pid,omitempty"`
	Alive      bool  `json:"alive,omitempty"`
}

// A EmqxProcessSpec defines the desired state of a EmqxProcess.
type EmqxProcessSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       EmqxProcessParameters `json:"forProvider"`
}

// A EmqxProcessStatus represents the observed state of a EmqxProcess.
type EmqxProcessStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          EmqxProcessObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A EmqxProcess is an example API type.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,baremetalprovider}
type EmqxProcess struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EmqxProcessSpec   `json:"spec"`
	Status EmqxProcessStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// EmqxProcessList contains a list of EmqxProcess
type EmqxProcessList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EmqxProcess `json:"items"`
}

// EmqxProcess type metadata.
var (
	EmqxProcessKind             = reflect.TypeOf(EmqxProcess{}).Name()
	EmqxProcessGroupKind        = schema.GroupKind{Group: Group, Kind: EmqxProcessKind}.String()
	EmqxProcessKindAPIVersion   = EmqxProcessKind + "." + SchemeGroupVersion.String()
	EmqxProcessGroupVersionKind = SchemeGroupVersion.WithKind(EmqxProcessKind)
)

func init() {
	SchemeBuilder.Register(&EmqxProcess{}, &EmqxProcessList{})
}
