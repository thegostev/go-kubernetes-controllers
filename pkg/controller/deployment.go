package controller

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/client"
	crcontroller "sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// DeploymentReconciler logs reconcile requests for Deployments
// (no business logic, just logs for demonstration)
type DeploymentReconciler struct {
	client.Client
}

func (r *DeploymentReconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	log.FromContext(ctx).Info("Reconciling Deployment", "namespace", req.Namespace, "name", req.Name)
	return reconcile.Result{}, nil
}

// DeploymentEventHandler logs all Deployment events
var DeploymentEventHandler = handler.Funcs{
	CreateFunc: func(ctx context.Context, e event.CreateEvent, q workqueue.RateLimitingInterface) {
		log.FromContext(ctx).Info("Deployment created", "name", e.Object.GetName(), "namespace", e.Object.GetNamespace())
		q.Add(reconcile.Request{NamespacedName: client.ObjectKeyFromObject(e.Object)})
	},
	UpdateFunc: func(ctx context.Context, e event.UpdateEvent, q workqueue.RateLimitingInterface) {
		log.FromContext(ctx).Info("Deployment updated", "name", e.ObjectNew.GetName(), "namespace", e.ObjectNew.GetNamespace())
		q.Add(reconcile.Request{NamespacedName: client.ObjectKeyFromObject(e.ObjectNew)})
	},
	DeleteFunc: func(ctx context.Context, e event.DeleteEvent, q workqueue.RateLimitingInterface) {
		log.FromContext(ctx).Info("Deployment deleted", "name", e.Object.GetName(), "namespace", e.Object.GetNamespace())
		q.Add(reconcile.Request{NamespacedName: client.ObjectKeyFromObject(e.Object)})
	},
	GenericFunc: func(ctx context.Context, e event.GenericEvent, q workqueue.RateLimitingInterface) {
		log.FromContext(ctx).Info("Generic event", "name", e.Object.GetName(), "namespace", e.Object.GetNamespace())
		q.Add(reconcile.Request{NamespacedName: client.ObjectKeyFromObject(e.Object)})
	},
}

// SetupDeploymentController registers the controller-runtime controller for Deployments
func SetupDeploymentController(mgr manager.Manager) error {
	reconciler := &DeploymentReconciler{Client: mgr.GetClient()}
	c, err := crcontroller.New("deployment-logger", mgr, crcontroller.Options{
		Reconciler: reconciler,
	})
	if err != nil {
		return err
	}
	return c.Watch(
		source.Kind(mgr.GetCache(), &appsv1.Deployment{}),
		&DeploymentEventHandler,
	)
}
