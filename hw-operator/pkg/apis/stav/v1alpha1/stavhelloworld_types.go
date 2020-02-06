package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// StavHelloWorldSpec defines the desired state of StavHelloWorld
type StavHelloWorldSpec struct {
	Message string `json:"message"`
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// StavHelloWorldStatus defines the observed state of StavHelloWorld
type StavHelloWorldStatus struct {
	Message string `json:"message"`
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// StavHelloWorld is the Schema for the stavhelloworlds API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=stavhelloworlds,scope=Namespaced
type StavHelloWorld struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StavHelloWorldSpec   `json:"spec,omitempty"`
	Status StavHelloWorldStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// StavHelloWorldList contains a list of StavHelloWorld
type StavHelloWorldList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []StavHelloWorld `json:"items"`
}

func init() {
	SchemeBuilder.Register(&StavHelloWorld{}, &StavHelloWorldList{})
}
