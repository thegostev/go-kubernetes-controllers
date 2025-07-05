package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// FrontendPageSpec defines the desired state of FrontendPage
type FrontendPageSpec struct {
	// Title is the display title of the frontend page
	Title string `json:"title"`

	// Template specifies the frontend template to use
	Template string `json:"template"`

	// Components defines the UI components to render on the page
	Components []Component `json:"components"`

	// Theme specifies the visual theme
	Theme string `json:"theme,omitempty"`
}

// Component defines a UI component on the page
type Component struct {
	// Name is the unique identifier for the component
	Name string `json:"name"`

	// Type specifies the component type
	Type string `json:"type"`

	// Config contains component-specific configuration
	Config map[string]interface{} `json:"config,omitempty"`
}

// FrontendPageStatus defines the observed state of FrontendPage
type FrontendPageStatus struct {
	// Phase represents the current phase of the frontend page
	Phase string `json:"phase"`

	// Message provides additional context about the status
	Message string `json:"message,omitempty"`

	// URL is the generated URL for accessing the page
	URL string `json:"url,omitempty"`

	// ComponentCount represents the number of components processed
	ComponentCount int `json:"componentCount,omitempty"`

	// LastUpdated tracks when the status was last updated
	LastUpdated *metav1.Time `json:"lastUpdated,omitempty"`
}

// FrontendPage is the Schema for the frontendpages API
type FrontendPage struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FrontendPageSpec   `json:"spec,omitempty"`
	Status FrontendPageStatus `json:"status,omitempty"`
}

// FrontendPageList contains a list of FrontendPage
type FrontendPageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FrontendPage `json:"items"`
}

// DeepCopy methods (minimal implementation)
func (in *FrontendPage) DeepCopy() *FrontendPage {
	if in == nil {
		return nil
	}
	out := new(FrontendPage)
	in.DeepCopyInto(out)
	return out
}

func (in *FrontendPage) DeepCopyInto(out *FrontendPage) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

func (in *FrontendPage) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *FrontendPageSpec) DeepCopy() *FrontendPageSpec {
	if in == nil {
		return nil
	}
	out := new(FrontendPageSpec)
	in.DeepCopyInto(out)
	return out
}

func (in *FrontendPageSpec) DeepCopyInto(out *FrontendPageSpec) {
	*out = *in
	if in.Components != nil {
		in, out := &in.Components, &out.Components
		*out = make([]Component, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

func (in *Component) DeepCopy() *Component {
	if in == nil {
		return nil
	}
	out := new(Component)
	in.DeepCopyInto(out)
	return out
}

func (in *Component) DeepCopyInto(out *Component) {
	*out = *in
	if in.Config != nil {
		in, out := &in.Config, &out.Config
		*out = make(map[string]interface{}, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

func (in *FrontendPageStatus) DeepCopy() *FrontendPageStatus {
	if in == nil {
		return nil
	}
	out := new(FrontendPageStatus)
	in.DeepCopyInto(out)
	return out
}

func (in *FrontendPageStatus) DeepCopyInto(out *FrontendPageStatus) {
	*out = *in
	if in.LastUpdated != nil {
		in, out := &in.LastUpdated, &out.LastUpdated
		*out = (*in).DeepCopy()
	}
}

func (in *FrontendPageList) DeepCopy() *FrontendPageList {
	if in == nil {
		return nil
	}
	out := new(FrontendPageList)
	in.DeepCopyInto(out)
	return out
}

func (in *FrontendPageList) DeepCopyInto(out *FrontendPageList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]FrontendPage, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

func (in *FrontendPageList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}
