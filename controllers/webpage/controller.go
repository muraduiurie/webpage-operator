package webpage

import (
	"context"
	"github.com/go-logr/logr"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"webpage/api/v1alpha1"
)

// define reconciler
type WebReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

func (wr *WebReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	wr.Log.Info("reconciling webpage", "name", req.Name, "namespace", req.Namespace)

	wp := v1alpha1.WebPage{}

	err := wr.Client.Get(ctx, req.NamespacedName, &wp)
	if err != nil && kerr.IsNotFound(err) {
		return ctrl.Result{}, nil
	} else if err != nil {
		return ctrl.Result{}, err
	}

	// reconciliation logic
	// For example, you might want to update the status of the WebPage resource, or create/update a Deployment based on the WebPage spec.

	wr.Log.Info("webpage reconciled")

	return ctrl.Result{}, nil
}

// SetupWithManager should specify your resource explicitly
func (wr *WebReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.WebPage{}).
		Complete(wr)
}
