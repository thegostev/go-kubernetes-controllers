package k8s

import (
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/thegostev/go-kubernetes-controllers/api/v1alpha1"
)

func NewScheme() *runtime.Scheme {
	scheme := runtime.NewScheme()
	_ = appsv1.AddToScheme(scheme)
	_ = v1alpha1.AddToScheme(scheme) // Add FrontendPage types
	return scheme
}
