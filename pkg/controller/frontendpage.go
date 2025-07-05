package controller

import (
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/client"
	crcontroller "sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/thegostev/go-kubernetes-controllers/api/v1alpha1"
)

// FrontendPageReconciler logs reconcile requests for FrontendPages
// (follows same pattern as DeploymentReconciler)
type FrontendPageReconciler struct {
	client.Client
}

func (r *FrontendPageReconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Reconciling FrontendPage", "namespace", req.Namespace, "name", req.Name)

	// Fetch the FrontendPage instance
	frontendPage := &v1alpha1.FrontendPage{}
	if err := r.Get(ctx, req.NamespacedName, frontendPage); err != nil {
		if client.IgnoreNotFound(err) != nil {
			logger.Error(err, "failed to get FrontendPage")
			return reconcile.Result{}, err
		}
		// Resource not found, likely deleted
		logger.Info("FrontendPage not found, likely deleted")
		return reconcile.Result{}, nil
	}

	// Simple reconciliation logic (following existing pattern)
	if err := r.reconcileFrontendPage(ctx, frontendPage); err != nil {
		logger.Error(err, "failed to reconcile FrontendPage")
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

// reconcileFrontendPage performs the actual reconciliation
func (r *FrontendPageReconciler) reconcileFrontendPage(ctx context.Context, frontendPage *v1alpha1.FrontendPage) error {
	logger := log.FromContext(ctx)

	// Update status to show reconciliation
	frontendPage.Status.Phase = "Ready"
	frontendPage.Status.Message = "Frontend page is ready"
	frontendPage.Status.ComponentCount = len(frontendPage.Spec.Components)
	frontendPage.Status.URL = fmt.Sprintf("http://localhost:8080/pages/%s/%s", frontendPage.Namespace, frontendPage.Name)
	frontendPage.Status.LastUpdated = &metav1.Time{Time: time.Now()}

	if err := r.Status().Update(ctx, frontendPage); err != nil {
		logger.Error(err, "failed to update FrontendPage status")
		return err
	}

	logger.Info("FrontendPage reconciled successfully",
		"namespace", frontendPage.Namespace,
		"name", frontendPage.Name,
		"components", len(frontendPage.Spec.Components))

	return nil
}

// FrontendPageEventHandler logs all FrontendPage events
// (follows exact same pattern as DeploymentEventHandler)
var FrontendPageEventHandler = handler.Funcs{
	CreateFunc: func(ctx context.Context, e event.CreateEvent, q workqueue.RateLimitingInterface) {
		log.FromContext(ctx).Info("FrontendPage created", "name", e.Object.GetName(), "namespace", e.Object.GetNamespace())
		q.Add(reconcile.Request{NamespacedName: client.ObjectKeyFromObject(e.Object)})
	},
	UpdateFunc: func(ctx context.Context, e event.UpdateEvent, q workqueue.RateLimitingInterface) {
		log.FromContext(ctx).Info("FrontendPage updated", "name", e.ObjectNew.GetName(), "namespace", e.ObjectNew.GetNamespace())
		q.Add(reconcile.Request{NamespacedName: client.ObjectKeyFromObject(e.ObjectNew)})
	},
	DeleteFunc: func(ctx context.Context, e event.DeleteEvent, q workqueue.RateLimitingInterface) {
		log.FromContext(ctx).Info("FrontendPage deleted", "name", e.Object.GetName(), "namespace", e.Object.GetNamespace())
		q.Add(reconcile.Request{NamespacedName: client.ObjectKeyFromObject(e.Object)})
	},
	GenericFunc: func(ctx context.Context, e event.GenericEvent, q workqueue.RateLimitingInterface) {
		log.FromContext(ctx).Info("Generic FrontendPage event", "name", e.Object.GetName(), "namespace", e.Object.GetNamespace())
		q.Add(reconcile.Request{NamespacedName: client.ObjectKeyFromObject(e.Object)})
	},
}

// SetupFrontendPageController registers the controller-runtime controller for FrontendPages
// (follows exact same pattern as SetupDeploymentController)
func SetupFrontendPageController(mgr manager.Manager) error {
	reconciler := &FrontendPageReconciler{Client: mgr.GetClient()}
	c, err := crcontroller.New("frontendpage-logger", mgr, crcontroller.Options{
		Reconciler: reconciler,
	})
	if err != nil {
		return err
	}
	return c.Watch(
		source.Kind(mgr.GetCache(), &v1alpha1.FrontendPage{}),
		&FrontendPageEventHandler,
	)
}
