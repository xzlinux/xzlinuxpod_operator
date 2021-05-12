package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// XzlinuxPodSpec defines the desired state of XzlinuxPod
type XzlinuxPodSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Replicas int `json:"replicas"`
}

// XzlinuxPodStatus defines the observed state of XzlinuxPod
type XzlinuxPodStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Replicas int `json:"replicas"`
	PodNames []string `json:"padNames"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// XzlinuxPod is the Schema for the xzlinuxpods API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=xzlinuxpods,scope=Namespaced
type XzlinuxPod struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   XzlinuxPodSpec   `json:"spec,omitempty"`
	Status XzlinuxPodStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// XzlinuxPodList contains a list of XzlinuxPod
type XzlinuxPodList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []XzlinuxPod `json:"items"`
}

func init() {
	SchemeBuilder.Register(&XzlinuxPod{}, &XzlinuxPodList{})
}
