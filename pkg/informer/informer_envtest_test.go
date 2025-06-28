//go:build integration
// +build integration

package informer

import (
	"context"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/envtest"

	"github.com/yourusername/k8s-controller-tutorial/internal/types"
)

func TestInformerReceivesDeploymentEvents(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}
	g := NewWithT(t)

	testEnv := &envtest.Environment{
		CRDDirectoryPaths:     []string{},
		ErrorIfCRDPathMissing: false,
	}
	cfg, err := testEnv.Start()
	g.Expect(err).ToNot(HaveOccurred())
	defer func() {
		if err := testEnv.Stop(); err != nil {
			t.Logf("failed to stop test environment: %v", err)
		}
	}()

	clientset, err := kubernetes.NewForConfig(cfg)
	g.Expect(err).ToNot(HaveOccurred())

	// Wait for the API server to be ready
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test API server connectivity
	_, err = clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{Limit: 1})
	g.Expect(err).ToNot(HaveOccurred())

	informerConfig := &types.InformerConfig{
		Namespace:       "default",
		ResyncPeriod:    1 * time.Second,
		Workers:         1,
		MaxCacheSize:    100,
		MaxConnections:  5,
		EventBufferSize: 10,
	}

	inf, err := NewInformer(clientset, informerConfig)
	g.Expect(err).ToNot(HaveOccurred())

	// Create a new context for the informer
	informerCtx, informerCancel := context.WithCancel(context.Background())
	defer informerCancel()

	// Start informer in a goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- inf.Start(informerCtx)
	}()

	// Give the informer time to sync
	time.Sleep(3 * time.Second)

	// Create a deployment
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-deployment",
			Namespace: "default",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "test"},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"app": "test"}},
				Spec:       corev1.PodSpec{Containers: []corev1.Container{{Name: "nginx", Image: "nginx"}}},
			},
		},
	}
	_, err = clientset.AppsV1().Deployments("default").Create(ctx, dep, metav1.CreateOptions{})
	g.Expect(err).ToNot(HaveOccurred())

	// Wait for the event to be processed
	time.Sleep(2 * time.Second)

	// Check that the deployment is in the cache
	deployments, err := inf.ListDeployments()
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(deployments).ToNot(BeEmpty())
	found := false
	for _, d := range deployments {
		if d.Name == "test-deployment" {
			found = true
		}
	}
	g.Expect(found).To(BeTrue())

	// Stop the informer
	informerCancel()

	// Wait for informer to stop
	select {
	case err := <-errChan:
		if err != nil {
			t.Logf("informer stopped with error: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Log("informer stop timeout")
	}
}

func int32Ptr(i int32) *int32 { return &i }
