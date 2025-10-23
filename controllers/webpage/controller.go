package webpage

import (
	"context"
	"webpage/api/v1alpha1"

	"github.com/go-logr/logr"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// define reconciler
type WebReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

func (wr *WebReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	wr.Log.Info("reconciling webpage", "name", req.Name, "namespace", req.Namespace)

	// Step 1: Fetch the WebPage resource
	wp := v1alpha1.WebPage{}
	err := wr.Client.Get(ctx, req.NamespacedName, &wp)
	if err != nil && kerr.IsNotFound(err) {
		// Resource was deleted, nothing to do
		return ctrl.Result{}, nil
	} else if err != nil {
		// Error fetching resource, requeue with backoff
		return ctrl.Result{}, err
	}

	// Step 2: Update status to Pending if not set

	// Step 3: Create or Update ConfigMap with the website content
	// - ConfigMap name: wp.Name + "-content"
	// - Data: {"index.html": wp.Spec.Content}
	// - Set owner reference so ConfigMap is deleted when WebPage is deleted
	// - Use controllerutil.CreateOrUpdate() for idempotent operation

	// Step 4: Create or Update Deployment with nginx
	// - Deployment name: wp.Name
	// - Replicas: wp.Spec.Replicas (default to 1 if not set)
	// - Image: wp.Spec.Image (default to "nginx:latest" if not set)
	// - Volume: mount the ConfigMap as a volume
	// - VolumeMount: mount to /usr/share/nginx/html (nginx default serving directory)
	// - Set owner reference for garbage collection

	// Step 5: Create or Update Service to expose the Deployment
	// - Service name: wp.Name + "-service"
	// - Type: ClusterIP (will be exposed via Ingress)
	// - Port: 80
	// - Selector: match the Deployment's labels
	// - Set owner reference

	// Step 6: Create or Update Ingress to expose the Service
	// - Ingress name: wp.Name + "-ingress"
	// - Host: could be derived from wp.Name + ".example.com" or a custom domain field
	// - Path: / (serve everything)
	// - Backend service: the Service created above
	// - Optional: Add annotations for ingress controller (nginx, traefik, etc.)
	// - Optional: Add TLS configuration
	// - Set owner reference

	// Step 7: Check if Deployment is ready and update WebPage status
	// - Fetch the Deployment
	// - Check if Deployment.Status.ReadyReplicas == Deployment.Spec.Replicas
	// - If ready: update wp.Status.Phase to WebPhaseRunning
	// - If not ready: update wp.Status.Phase to WebPhasePending
	// - Update the status subresource

	// Step 8: Requeue after some time to periodically check status
	// - This ensures the status stays up-to-date even if events are missed
	// return ctrl.Result{RequeueAfter: 30 * time.Second}, nil

	wr.Log.Info("webpage reconciled successfully")
	return ctrl.Result{}, nil
}

// SetupWithManager should specify your resource explicitly
func (wr *WebReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.WebPage{}).
		// Uncomment these to watch owned resources for automatic reconciliation
		// when they change (e.g., if someone manually deletes the Deployment)
		// Owns(&appsv1.Deployment{}).
		// Owns(&corev1.Service{}).
		// Owns(&corev1.ConfigMap{}).
		// Owns(&networkingv1.Ingress{}).
		Complete(wr)
}
