/*
Copyright 2024.

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

// ActionType is a custom type for representing different calculation actions.
type ActionType string

const (
	// ActionTypeAdd represents the "add" action.
	ActionTypeAdd ActionType = ActionType("add")
	// ActionTypeSub represents the "subtract" action.
	ActionTypeSub ActionType = ActionType("sub")
	// ActionTypeMul represents the "multiply" action.
	ActionTypeMul ActionType = ActionType("mul")
	// ActionTypeDiv represents the "divide" action.
	ActionTypeDiv ActionType = ActionType("div")
)

// CalculateSpec defines the desired state of Calculate
type CalculateSpec struct {
	// The arithmetic action to be performed (add, sub, mul, or div).
	Action ActionType `json:"action,omitempty"`
	// The first operand in the calculation.
	First int `json:"first,omitempty"`
	// The second operand in the calculation.
	Second int `json:"second,omitempty"`
}

// CalculateStatus defines the observed state of Calculate
type CalculateStatus struct {
	// The result of the calculation.
	Result int `json:"result,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Action",type="string",JSONPath=".spec.action",description="The math action type"
// +kubebuilder:printcolumn:name="First",type="integer",JSONPath=".spec.first",description="Input first number"
// +kubebuilder:printcolumn:name="Second",type="integer",JSONPath=".spec.second",description="Input second number"
// +kubebuilder:printcolumn:name="Result",type="integer",JSONPath=".status.result",description="Calculate result"

// Calculate is the Schema for the calculates API
type Calculate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CalculateSpec   `json:"spec,omitempty"`
	Status CalculateStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// CalculateList contains a list of Calculate
type CalculateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Calculate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Calculate{}, &CalculateList{})
}
