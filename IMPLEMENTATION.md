# Implementation Guide

This document provides a quick reference for implementing the complete WebPage operator functionality.

## Overview

The operator creates and manages these Kubernetes resources for each WebPage:
- **ConfigMap**: Stores website content
- **Deployment**: Runs nginx pods serving the content
- **Service**: Exposes the deployment internally
- **Ingress**: Routes external traffic to the service

## Quick Start

### 1. Uncomment the Reconciliation Code

Open `controllers/webpage/controller.go` and uncomment all the code blocks in the `Reconcile()` function (Steps 2-8).

### 2. Add Required Imports

Add these imports to `controllers/webpage/controller.go`:

```go
import (
    "context"
    "time"
    
    "github.com/go-logr/logr"
    appsv1 "k8s.io/api/apps/v1"
    corev1 "k8s.io/api/core/v1"
    networkingv1 "k8s.io/api/networking/v1"
    kerr "k8s.io/apimachinery/pkg/api/errors"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/runtime"
    "k8s.io/apimachinery/pkg/types"
    "k8s.io/apimachinery/pkg/util/intstr"
    ctrl "sigs.k8s.io/controller-runtime"
    "sigs.k8s.io/controller-runtime/pkg/client"
    "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
    
    "webpage/api/v1alpha1"
)
```

### 3. Update SetupWithManager

Uncomment the `.Owns()` calls in `SetupWithManager()`:

```go
func (wr *WebReconciler) SetupWithManager(mgr ctrl.Manager) error {
    return ctrl.NewControllerManagedBy(mgr).
        For(&v1alpha1.WebPage{}).
        Owns(&appsv1.Deployment{}).
        Owns(&corev1.Service{}).
        Owns(&corev1.ConfigMap{}).
        Owns(&networkingv1.Ingress{}).
        Complete(wr)
}
```

### 4. Verify Dependencies

Run:
```bash
go mod tidy
```

### 5. Test the Operator

```bash
# Install CRD
kubectl apply -f config/crd/bases/data-platform.qonto.co_webpages.yaml

# Run operator
go run cmd/webpage/main.go

# In another terminal, create a WebPage
kubectl apply -f examples/webpage-sample.yaml

# Check resources
kubectl get webpages
kubectl get pods
kubectl get services
kubectl get ingress
```

## Reconciliation Steps Explained

### Step 2: Initialize Status
Sets the initial phase to "Pending" when a new WebPage is created.

### Step 3: ConfigMap Creation
Creates a ConfigMap with the website content:
- Key: `index.html`
- Value: Content from `spec.content`
- Name: `{webpage-name}-content`

### Step 4: Deployment Creation
Creates an nginx Deployment:
- Uses `spec.image` (defaults to `nginx:latest`)
- Uses `spec.replicas` (defaults to 1)
- Mounts ConfigMap at `/usr/share/nginx/html`
- Labels: `app: {webpage-name}`

### Step 5: Service Creation
Creates a ClusterIP Service:
- Name: `{webpage-name}-service`
- Port: 80
- Selector: `app: {webpage-name}`

### Step 6: Ingress Creation
Creates an Ingress:
- Name: `{webpage-name}-ingress`
- Host: `{webpage-name}.example.com`
- Path: `/` → Service
- Optional TLS support (commented out)

### Step 7: Status Update
Monitors Deployment readiness:
- Checks `ReadyReplicas == Replicas`
- Updates phase to "Running" when ready
- Keeps "Pending" otherwise

### Step 8: Periodic Requeue
Optional requeue after 30 seconds to ensure status stays current.

## Testing the Implementation

### Create a Test WebPage

```yaml
apiVersion: data-platform.qonto.co/v1alpha1
kind: WebPage
metadata:
  name: test-site
  namespace: default
spec:
  content: |
    <!DOCTYPE html>
    <html>
      <head><title>Test Site</title></head>
      <body>
        <h1>Hello from Kubernetes Operator!</h1>
        <p>This page is managed by a custom operator.</p>
      </body>
    </html>
  image: "nginx:1.21"
  replicas: 2
```

### Verify Resources

```bash
# Check WebPage status
kubectl get webpage test-site -o yaml

# Check generated resources
kubectl get configmap test-site-content
kubectl get deployment test-site
kubectl get service test-site-service
kubectl get ingress test-site-ingress

# View ConfigMap content
kubectl get configmap test-site-content -o yaml

# Check pod logs
kubectl logs -l app=test-site

# Describe to see events
kubectl describe webpage test-site
```

### Test Updates

```bash
# Update content
kubectl edit webpage test-site
# Change spec.content and save

# Watch reconciliation
kubectl get pods -w

# Update replicas
kubectl patch webpage test-site -p '{"spec":{"replicas":3}}' --type=merge

# Watch deployment scaling
kubectl get deployment test-site -w
```

### Test Deletion

```bash
# Delete WebPage
kubectl delete webpage test-site

# Verify all resources are deleted (due to owner references)
kubectl get deployment,service,configmap,ingress -l app=test-site
# Should return "No resources found"
```

## Common Issues and Solutions

### Issue: Ingress not accessible
**Solution**: 
- Ensure you have an Ingress controller installed (nginx-ingress, traefik, etc.)
- Update `/etc/hosts` or DNS: `<ingress-ip> test-site.example.com`
- Add necessary annotations for your Ingress controller

### Issue: Pods not starting
**Solution**:
- Check pod logs: `kubectl logs -l app=test-site`
- Check events: `kubectl describe pod <pod-name>`
- Verify image exists and is accessible

### Issue: ConfigMap not mounting
**Solution**:
- Verify ConfigMap exists: `kubectl get configmap test-site-content`
- Check pod volume mounts: `kubectl describe pod <pod-name>`

### Issue: Status not updating
**Solution**:
- Check operator logs for errors
- Ensure Step 7 code is uncommented
- Verify RBAC permissions for status subresource

## RBAC Permissions Required

If deploying in-cluster, ensure the operator ServiceAccount has these permissions:

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: webpage-operator-role
rules:
- apiGroups:
  - data-platform.qonto.co
  resources:
  - webpages
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
  - watch
- apiGroups:
  - data-platform.qonto.co
  resources:
  - webpages/status
  verbs:
  - get
  - patch
  - update
- apiGroups:
  - ""
  resources:
  - configmaps
  - services
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
  - watch
- apiGroups:
  - apps
  resources:
  - deployments
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
  - watch
- apiGroups:
  - networking.k8s.io
  resources:
  - ingresses
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
  - watch
```

## Next Steps

1. **Add Error Handling**: Improve error messages and handling
2. **Add Events**: Emit Kubernetes events for important actions
3. **Add Metrics**: Instrument with Prometheus metrics
4. **Add Webhooks**: Validation and defaulting webhooks
5. **Add Tests**: Unit and integration tests
6. **Package for Production**: Create Deployment manifests, Helm charts

## Resources

- [Kubernetes API Reference](https://kubernetes.io/docs/reference/kubernetes-api/)
- [controller-runtime Documentation](https://pkg.go.dev/sigs.k8s.io/controller-runtime)
- [Kubebuilder Book](https://book.kubebuilder.io/)

