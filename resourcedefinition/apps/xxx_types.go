package apps

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +enum
type XXXPhase string

// These are the valid phases of a xxx.
const (
	// XXXRunning means the xxx is in the running state.
	XXXRunning XXXPhase = "Running"
	// XXXPending means the xxx is in the pending state.
	XXXPending XXXPhase = "Pending"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// XXX is an example definition of a Kubernetes resource object.
type XXX struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the behavior of the XXX.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
	// +optional
	Spec XXXSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`

	// Status describes the current status of a XXX.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
	// +optional
	Status XXXStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// XXXSpec describes the attributes on a XXX.
type XXXSpec struct {
	// DisplayName is the display name of the XXX resource.
	DisplayName string `json:"displayName" protobuf:"bytes,1,opt,name=displayName"`

	// Description provides a brief summary of the XXX's purpose or functionality.
	// It helps users understand what the XXX is intended for.
	// +optional
	Description string `json:"description,omitempty" protobuf:"bytes,2,opt,name=description"`

	// You can add more status fields as needed.
	// +optional
}

// XXXStatus is information about the current status of a XXX.
type XXXStatus struct {
	// Phase is the current lifecycle phase of the xxx.
	// +optional
	Phase XXXPhase `json:"phase,omitempty" protobuf:"bytes,1,opt,name=phase,casttype=NamespacePhase"`

	// ObservedGeneration reflects the generation of the most recently observed XXX.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty" protobuf:"varint,2,opt,name=observedGeneration"`

	// Represents the latest available observations of a xxx's current state.
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,3,rep,name=conditions"`

	// You can add more status fields as needed.
	// +optional
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// XXXList is a list of XXX objects.
type XXXList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Items is a list of schema objects.
	Items []XXX `json:"items" protobuf:"bytes,2,rep,name=items"`
}
