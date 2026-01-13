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

// MqttTopicParameters are the configurable fields of a MqttTopic.
type MqttTopicParameters struct {
	Name             string `json:"name"`
	Host             string `json:"host"`
	Username         string `json:"username"`
	Password         string `json:"password"`
	DegradedTreshold string `json:"degradedtreshold"`
}

type MqttTopicObservation struct {
	Metrics    Metrics `json:"metrics,omitempty"`
	Topic      string  `json:"topic,omitempty"`
	CreateTime string  `json:"create_time,omitempty"`
}
type Metrics struct {
	MessagesQos2OutCount string `json:"messages.qos2.out.count,omitempty"`
	MessagesQos2InCount  string `json:"messages.qos2.in.count,omitempty"`
	MessagesQos1OutCount string `json:"messages.qos1.out.count,omitempty"`
	MessagesQos1InCount  string `json:"messages.qos1.in.count,omitempty"`
	MessagesQos0OutCount string `json:"messages.qos0.out.count,omitempty"`
	MessagesQos0InCount  string `json:"messages.qos0.in.count,omitempty"`
	MessagesOutCount     string `json:"messages.out.count,omitempty"`
	MessagesInCount      string `json:"messages.in.count,omitempty"`
	MessagesDroppedCount string `json:"messages.dropped.count,omitempty"`
	MessagesQos2OutRate  string `json:"messages.qos2.out.rate,omitempty"`
	MessagesQos2InRate   string `json:"messages.qos2.in.rate,omitempty"`
	MessagesQos1OutRate  string `json:"messages.qos1.out.rate,omitempty"`
	MessagesQos1InRate   string `json:"messages.qos1.in.rate,omitempty"`
	MessagesQos0OutRate  string `json:"messages.qos0.out.rate,omitempty"`
	MessagesQos0InRate   string `json:"messages.qos0.in.rate,omitempty"`
	MessagesOutRate      string `json:"messages.out.rate,omitempty"`
	MessagesInRate       string `json:"messages.in.rate,omitempty"`
	MessagesDroppedRate  string `json:"messages.dropped.rate,omitempty"`
}

type MqttTopicState struct {
	Exist       string `json:"exist"`
	Degraded    string `json:"degraded"`
	HostAddress string `json:"host_address"`
	QueueState  string `json:"queue_state"`
}

// A MqttTopicSpec defines the desired state of a MqttTopic.
type MqttTopicSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       MqttTopicParameters `json:"forProvider"`
}

// A MqttTopicStatus represents the observed state of a MqttTopic.
type MqttTopicStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          *MqttTopicObservation `json:"atProvider,omitempty"`
	AtResource          MqttTopicState        `json:"atResource,omitempty"`
}

// +kubebuilder:object:root=true

// A MqttTopic is an example API type.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,topicprovider}
type MqttTopic struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MqttTopicSpec   `json:"spec"`
	Status MqttTopicStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MqttTopicList contains a list of MqttTopic
type MqttTopicList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MqttTopic `json:"items"`
}

// MqttTopic type metadata.
var (
	MqttTopicKind             = reflect.TypeOf(MqttTopic{}).Name()
	MqttTopicGroupKind        = schema.GroupKind{Group: Group, Kind: MqttTopicKind}.String()
	MqttTopicKindAPIVersion   = MqttTopicKind + "." + SchemeGroupVersion.String()
	MqttTopicGroupVersionKind = SchemeGroupVersion.WithKind(MqttTopicKind)
)

func init() {
	SchemeBuilder.Register(&MqttTopic{}, &MqttTopicList{})
}
