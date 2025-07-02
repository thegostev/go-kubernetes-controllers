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

type DeploymentReconciler struct {
	client.Client
}

func (r *DeploymentReconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	// No business logic needed for just logging events
	return reconcile.Result{}, nil
}

var DeploymentEventHandler = handler.Funcs{
	CreateFunc: func(ctx context.Context, e event.CreateEvent, q workqueue.RateLimitingInterface) {
		log.Log.Info("Deployment created", "name", e.Object.GetName())
		q.Add(reconcile.Request{NamespacedName: client.ObjectKeyFromObject(e.Object)})
	},
	UpdateFunc: func(ctx context.Context, e event.UpdateEvent, q workqueue.RateLimitingInterface) {
		log.Log.Info("Deployment updated", "name", e.ObjectNew.GetName())
		q.Add(reconcile.Request{NamespacedName: client.ObjectKeyFromObject(e.ObjectNew)})
	},
	DeleteFunc: func(ctx context.Context, e event.DeleteEvent, q workqueue.RateLimitingInterface) {
		log.Log.Info("Deployment deleted", "name", e.Object.GetName())
		q.Add(reconcile.Request{NamespacedName: client.ObjectKeyFromObject(e.Object)})
	},
	GenericFunc: func(ctx context.Context, e event.GenericEvent, q workqueue.RateLimitingInterface) {
		log.Log.Info("Generic event", "name", e.Object.GetName())
		q.Add(reconcile.Request{NamespacedName: client.ObjectKeyFromObject(e.Object)})
	},
}

// SetupDeploymentController sets up the controller-runtime controller for Deployments.
func SetupDeploymentController(mgr manager.Manager) error {
	reconciler := &DeploymentReconciler{Client: mgr.GetClient()}
	c, err := crcontroller.New("deployment-controller", mgr, crcontroller.Options{
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
