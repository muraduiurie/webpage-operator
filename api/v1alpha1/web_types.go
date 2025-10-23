// +kubebuilder:object:generate=true
// +groupName=data-platform.qonto.co
package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type WebPhase string

const (
	WebPhasePending WebPhase = "Pending"
	WebPhaseRunning WebPhase = "Running"
)

func init() {
	SchemeBuilder.Register(&WebPage{}, &WebPageList{})
}

// WebPage represents the web pod, pvc and ingress subresources
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,path=webpages,shortName=web,singular=webpage
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type WebPage struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              WebPageSpec   `json:"spec"`
	Status            WebPageStatus `json:"status,omitempty"`
}

// WebPageSpec defines the desired state of the website
type WebPageSpec struct {
	// Content is the static content exposed by our website
	// +equired
	Content string `json:"content"`
	// Image is the docker image used to expose the content of our website
	// +optional
	Image string `json:"image"`
	// Replicas defines the number of running pods of our website
	// +optional
	Replicas int `json:"replicas"`
}

// WebPageStatus represents the observed state of our website
type WebPageStatus struct {
	// Phase is the actual state of our website
	// +optional
	Phase WebPhase `json:"phase"`
}

// WebPageList is the list of multiple WebPage resources
// +kubebuilder:object:root=true
type WebPageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []WebPage `json:"items"`
}
