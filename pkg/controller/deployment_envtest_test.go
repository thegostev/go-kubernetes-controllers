//go:build integration

package controller

import (
	"context"
	"testing"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

func TestDeploymentEventLogging(t *testing.T) {
	testEnv := &envtest.Environment{}
	cfg, err := testEnv.Start()
	if err != nil {
		t.Fatalf("failed to start envtest: %v", err)
	}
	defer func() { _ = testEnv.Stop() }()

	s := runtime.NewScheme()
	_ = scheme.AddToScheme(s)
	_ = appsv1.AddToScheme(s)
	_ = corev1.AddToScheme(s)

	mgr, err := ctrl.NewManager(cfg, ctrl.Options{Scheme: s})
	if err != nil {
		t.Fatalf("failed to create manager: %v", err)
	}

	eventCh := make(chan string, 10)
	handler := handler.Funcs{
		CreateFunc: func(ctx context.Context, e event.CreateEvent, q workqueue.RateLimitingInterface) {
			eventCh <- "created"
		},
		UpdateFunc: func(ctx context.Context, e event.UpdateEvent, q workqueue.RateLimitingInterface) {
			eventCh <- "updated"
		},
		DeleteFunc: func(ctx context.Context, e event.DeleteEvent, q workqueue.RateLimitingInterface) {
			eventCh <- "deleted"
		},
		GenericFunc: func(ctx context.Context, e event.GenericEvent, q workqueue.RateLimitingInterface) {
			eventCh <- "generic"
		},
	}

	c, err := controller.New("deployment-event-test", mgr, controller.Options{
		Reconciler: &DeploymentReconciler{Client: mgr.GetClient()},
	})
	if err != nil {
		t.Fatalf("failed to build controller: %v", err)
	}
	if err := c.Watch(source.Kind(mgr.GetCache(), &appsv1.Deployment{}), handler); err != nil {
		t.Fatalf("failed to set up watch: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() { _ = mgr.Start(ctx) }()

	k8sClient := mgr.GetClient()
	ns := "default"
	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: "test-deploy", Namespace: ns},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"app": "test"}},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"app": "test"}},
				Spec:       corev1.PodSpec{Containers: []corev1.Container{{Name: "nginx", Image: "nginx"}}},
			},
		},
	}
	if err := k8sClient.Create(ctx, deploy); err != nil {
		t.Fatalf("failed to create deployment: %v", err)
	}

	// Wait for create event
	if err := waitForEvent(eventCh, "created"); err != nil {
		t.Fatalf("create event not logged: %v", err)
	}

	// Update deployment
	deploy.Spec.Replicas = int32Ptr(2)
	if err := k8sClient.Update(ctx, deploy); err != nil {
		t.Fatalf("failed to update deployment: %v", err)
	}
	if err := waitForEvent(eventCh, "updated"); err != nil {
		t.Fatalf("update event not logged: %v", err)
	}

	// Delete deployment
	if err := k8sClient.Delete(ctx, deploy); err != nil {
		t.Fatalf("failed to delete deployment: %v", err)
	}
	if err := waitForEvent(eventCh, "deleted"); err != nil {
		t.Fatalf("delete event not logged: %v", err)
	}
}

func waitForEvent(ch <-chan string, want string) error {
	timeout := time.After(5 * time.Second)
	for {
		select {
		case got := <-ch:
			if got == want {
				return nil
			}
		case <-timeout:
			return context.DeadlineExceeded
		}
	}
}

func int32Ptr(i int32) *int32 { return &i }
