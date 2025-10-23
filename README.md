# WebPage Operator - Kubernetes Operator Learning Project

A simple Kubernetes operator built with [controller-runtime](https://github.com/kubernetes-sigs/controller-runtime) to demonstrate the fundamental concepts and structure of Kubernetes operators. This operator manages a custom resource called `WebPage` that represents a simple web server configuration.

## Table of Contents

- [What is a Kubernetes Operator?](#what-is-a-kubernetes-operator)
- [Project Overview](#project-overview)
- [Architecture](#architecture)
- [Development Steps](#development-steps)
  - [Step 1: Define Custom Resource Types](#step-1-define-custom-resource-types)
  - [Step 2: Create Boilerplate Code](#step-2-create-boilerplate-code)
  - [Step 3: Generate CRD and DeepCopy Methods](#step-3-generate-crd-and-deepcopy-methods)
  - [Step 4: Set Up Manager and Controller](#step-4-set-up-manager-and-controller)
  - [Step 5: Implement Reconciler](#step-5-implement-reconciler)
  - [Step 6: Run the Operator](#step-6-run-the-operator)
- [Project Structure](#project-structure)
- [Building and Running](#building-and-running)
- [Understanding the Components](#understanding-the-components)
- [Next Steps](#next-steps)

## What is a Kubernetes Operator?

A Kubernetes operator is a method of packaging, deploying, and managing a Kubernetes application. It uses custom resources (CRs) to manage applications and their components, following Kubernetes principles, notably the control loop pattern.

The operator pattern consists of:
- **Custom Resource Definition (CRD)**: Extends Kubernetes API with custom resources
- **Controller**: Watches for changes to custom resources and reconciles the desired state with the actual state
- **Reconciliation Loop**: Continuously ensures the actual state matches the desired state

## Project Overview

This operator manages a `WebPage` custom resource that could be extended to deploy and manage simple web servers in Kubernetes. Currently, it demonstrates the core operator pattern with:

- A custom resource `WebPage` with spec fields for:
  - `content`: The static content to serve
  - `image`: Docker image to use (optional)
  - `replicas`: Number of replicas (optional)
- Status tracking with a `phase` field (Pending/Running)
- A reconciliation loop that watches for WebPage resource changes

## Architecture

```
┌─────────────────────────────────────────┐
│         Kubernetes API Server           │
│                                         │
│  ┌───────────────────────────────────┐  │
│  │   Custom Resource Definition      │  │
│  │      (WebPage CRD)                │  │
│  └───────────────────────────────────┘  │
└─────────────────────────────────────────┘
                    │
                    │ watches
                    ▼
┌─────────────────────────────────────────┐
│         WebPage Operator                │
│                                         │
│  ┌───────────────────────────────────┐  │
│  │  Manager                          │  │
│  │  ┌─────────────────────────────┐  │  │
│  │  │  Controller                 │  │  │
│  │  │  ┌───────────────────────┐  │  │  │
│  │  │  │  WebReconciler        │  │  │  │
│  │  │  │  - Reconcile()        │  │  │  │
│  │  │  └───────────────────────┘  │  │  │
│  │  └─────────────────────────────┘  │  │
│  └───────────────────────────────────┘  │
└─────────────────────────────────────────┘
                    │
                    │ manages
                    ▼
┌─────────────────────────────────────────┐
│    WebPage Custom Resources             │
└─────────────────────────────────────────┘
```

## Development Steps

### Step 1: Define Custom Resource Types

**File**: `api/v1alpha1/web_types.go`

The first step in building an operator is defining your Custom Resource. This involves:

1. **Define the API Group and Version**:
   ```go
   // +groupName=data-platform.qonto.co
   package v1alpha1
   ```
   - API Group: `data-platform.qonto.co`
   - Version: `v1alpha1` (alpha indicates it's still under development)

2. **Define the Spec** (desired state):
   ```go
   type WebPageSpec struct {
       Content  string `json:"content"`   // Required field
       Image    string `json:"image"`     // Optional field
       Replicas int    `json:"replicas"`  // Optional field
   }
   ```
   The spec represents what the user wants (desired state).

3. **Define the Status** (observed state):
   ```go
   type WebPageStatus struct {
       Phase WebPhase `json:"phase"` // Current state: Pending/Running
   }
   ```
   The status represents the current actual state of the resource.

4. **Define the Main Resource Type**:
   ```go
   type WebPage struct {
       metav1.TypeMeta   `json:",inline"`
       metav1.ObjectMeta `json:"metadata,omitempty"`
       Spec              WebPageSpec   `json:"spec"`
       Status            WebPageStatus `json:"status,omitempty"`
   }
   ```

5. **Add Kubebuilder Markers** (special comments):
   - `+kubebuilder:object:root=true`: Marks this as a root API object
   - `+kubebuilder:subresource:status`: Enables status subresource
   - `+kubebuilder:resource`: Defines resource properties (scope, path, shortName)
   - `+kubebuilder:printcolumn`: Adds columns when using `kubectl get`

6. **Register the Types**:
   ```go
   func init() {
       SchemeBuilder.Register(&WebPage{}, &WebPageList{})
   }
   ```

**Key File**: `api/v1alpha1/groupversion_info.go`

This file sets up the API group version information:
```go
var (
    GroupVersion = schema.GroupVersion{
        Group:   "data-platform.qonto.co", 
        Version: "v1alpha1"
    }
    SchemeBuilder = &scheme.Builder{GroupVersion: GroupVersion}
    AddToScheme   = SchemeBuilder.AddToScheme
)
```

### Step 2: Create Boilerplate Code

**File**: `hack/boilerplate.go.txt`

This file contains the license header that will be added to all generated files. The controller-gen tool uses this when generating code.

### Step 3: Generate CRD and DeepCopy Methods

After defining your types, you need to generate:

1. **DeepCopy Methods** - Required for all Kubernetes API types
2. **CRD Manifest** - The YAML file that extends Kubernetes API

**Generated Files**:
- `api/v1alpha1/zz_generated.deepcopy.go` - Auto-generated DeepCopy methods
- `config/crd/bases/data-platform.qonto.co_webpages.yaml` - CRD manifest

**Commands** (defined in Makefile):

```bash
# Generate DeepCopy methods
make generate
# This runs: controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./api/..."

# Generate CRD manifests
make manifests
# This runs: controller-gen crd paths="./api/..." output:crd:dir=./config/crd/bases
```

**Understanding `zz_generated.deepcopy.go`**:

Kubernetes requires all API types to implement the `runtime.Object` interface, which includes:
- `DeepCopy()`: Creates a complete copy of the object
- `DeepCopyInto()`: Copies the object into another
- `DeepCopyObject()`: Returns a deep copy as runtime.Object

These methods ensure that when Kubernetes manipulates your objects, it doesn't accidentally modify the original.

**Understanding the CRD**:

The generated CRD (`data-platform.qonto.co_webpages.yaml`) defines:
- The API group, version, and resource names
- The OpenAPI v3 schema for validation
- Additional printer columns for `kubectl get` output
- Status subresource configuration

### Step 4: Set Up Manager and Controller

**File**: `cmd/webpage/main.go`

The main.go file is the entry point of your operator. It sets up:

1. **Scheme** (type registry):
   ```go
   var scheme = runtime.NewScheme()
   
   func init() {
       utilruntime.Must(k8sscehme.AddToScheme(scheme))     // Add core K8s types
       utilruntime.Must(v1alpha1.AddToScheme(scheme))      // Add your custom types
   }
   ```
   The scheme is a registry that maps Go types to GroupVersionKinds (GVKs). It's essential for serialization/deserialization.

2. **Logger**:
   ```go
   logger := zap.New()
   ctrl.SetLogger(logger)
   ```
   Sets up structured logging using Zap.

3. **Manager**:
   ```go
   mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
       Scheme: scheme,
   })
   ```
   The Manager:
   - Connects to the Kubernetes API server
   - Provides shared dependencies (client, cache, scheme)
   - Manages the lifecycle of controllers
   - Handles leader election (for HA setups)

4. **Reconciler Initialization**:
   ```go
   wr := webpage.WebReconciler{
       Client: mgr.GetClient(),  // K8s client for API operations
       Scheme: mgr.GetScheme(),  // Type registry
       Log:    log.WithName("web-reconciler"),  // Scoped logger
   }
   ```

5. **Controller Setup**:
   ```go
   err = wr.SetupWithManager(mgr)
   ```
   Registers the reconciler with the manager and sets up watches.

6. **Start the Manager**:
   ```go
   ctx := ctrl.SetupSignalHandler()  // Handle SIGTERM/SIGINT gracefully
   if err = mgr.Start(ctx); err != nil {
       log.Error(err, "problem running manager")
       os.Exit(1)
   }
   ```

### Step 5: Implement Reconciler

**File**: `controllers/webpage/controller.go`

The reconciler is the heart of your operator. It implements the control loop pattern.

1. **Define the Reconciler Struct**:
   ```go
   type WebReconciler struct {
       client.Client           // For CRUD operations on K8s resources
       Scheme *runtime.Scheme  // Type registry
       Log    logr.Logger      // Logger
   }
   ```

2. **Implement the Reconcile Function**:
   ```go
   func (wr *WebReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error)
   ```
   
   The Reconcile function is called whenever:
   - A WebPage resource is created, updated, or deleted
   - A periodic resync occurs
   - The controller is restarted

   **Reconcile Logic**:
   ```go
   // 1. Fetch the WebPage resource
   wp := v1alpha1.WebPage{}
   err := wr.Client.Get(ctx, req.NamespacedName, &wp)
   
   // 2. Handle deletion (resource not found)
   if err != nil && kerr.IsNotFound(err) {
       return ctrl.Result{}, nil  // Resource deleted, nothing to do
   } else if err != nil {
       return ctrl.Result{}, err  // Requeue on error
   }
   
   // 3. Reconciliation logic goes here
   // - Create/update Deployments based on wp.Spec
   // - Update wp.Status with current state
   // - Handle any child resources (Services, Ingresses, etc.)
   
   return ctrl.Result{}, nil
   ```

   **Return Values**:
   - `ctrl.Result{}`: Reconciliation successful, don't requeue
   - `ctrl.Result{Requeue: true}`: Requeue immediately
   - `ctrl.Result{RequeueAfter: time.Minute}`: Requeue after a delay
   - `return err`: Requeue with exponential backoff (on error)

3. **Set Up Watches**:
   ```go
   func (wr *WebReconciler) SetupWithManager(mgr ctrl.Manager) error {
       return ctrl.NewControllerManagedBy(mgr).
           For(&v1alpha1.WebPage{}).  // Watch WebPage resources
           Complete(wr)
   }
   ```
   
   This sets up a watch on WebPage resources. You can also watch "owned" resources:
   ```go
   ctrl.NewControllerManagedBy(mgr).
       For(&v1alpha1.WebPage{}).
       Owns(&appsv1.Deployment{}).  // Watch Deployments owned by WebPage
       Complete(wr)
   ```

### Step 6: Run the Operator

To run your operator:

1. **Install the CRD**:
   ```bash
   kubectl apply -f config/crd/bases/data-platform.qonto.co_webpages.yaml
   ```

2. **Run the Operator**:
   ```bash
   # Option 1: Run locally (outside cluster)
   go run cmd/webpage/main.go
   
   # Option 2: Build and run in cluster
   docker build -t webpage-operator:latest .
   kubectl apply -f deploy/  # Deployment manifests
   ```

3. **Create a WebPage Resource**:
   ```yaml
   apiVersion: data-platform.qonto.co/v1alpha1
   kind: WebPage
   metadata:
     name: my-webpage
     namespace: default
   spec:
     content: "Hello, Kubernetes Operator!"
     image: nginx:latest
     replicas: 3
   ```

4. **Verify**:
   ```bash
   # Check the resource
   kubectl get webpages
   # or use the short name
   kubectl get web
   
   # Get details
   kubectl describe webpage my-webpage
   ```

## Project Structure

```
webpage-operator/
├── api/
│   └── v1alpha1/
│       ├── groupversion_info.go        # API group/version definitions
│       ├── web_types.go                # WebPage type definitions
│       └── zz_generated.deepcopy.go    # Generated DeepCopy methods
├── bin/
│   └── controller-gen                  # Code generation tool
├── cmd/
│   └── webpage/
│       └── main.go                     # Operator entry point
├── config/
│   └── crd/
│       └── bases/
│           └── data-platform.qonto.co_webpages.yaml  # Generated CRD
├── controllers/
│   └── webpage/
│       └── controller.go               # Reconciliation logic
├── hack/
│   └── boilerplate.go.txt             # License header template
├── go.mod                              # Go module dependencies
├── go.sum                              # Dependency checksums
└── Makefile                            # Build automation
```

## Building and Running

### Prerequisites

- Go 1.25+ installed
- Access to a Kubernetes cluster (minikube, kind, or production cluster)
- kubectl configured

### Development Workflow

```bash
# 1. Modify types in api/v1alpha1/web_types.go
vim api/v1alpha1/web_types.go

# 2. Regenerate code
make generate    # Generate DeepCopy methods
make manifests   # Generate/update CRD

# 3. Install/Update CRD in cluster
kubectl apply -f config/crd/bases/data-platform.qonto.co_webpages.yaml

# 4. Run operator locally
go run cmd/webpage/main.go

# 5. In another terminal, create a WebPage resource
kubectl apply -f examples/webpage-sample.yaml

# 6. Watch logs and verify reconciliation
```

### Example WebPage Resource

Create a file `examples/webpage-sample.yaml`:

```yaml
apiVersion: data-platform.qonto.co/v1alpha1
kind: WebPage
metadata:
  name: example-webpage
  namespace: default
spec:
  content: "Welcome to my website!"
  image: "nginx:1.21"
  replicas: 2
```

Apply it:
```bash
kubectl apply -f examples/webpage-sample.yaml
```

Check the results:
```bash
# List WebPages
kubectl get webpages

# Output:
# NAME              STATUS    AGE
# example-webpage             30s

# Get detailed information
kubectl describe webpage example-webpage

# Get YAML
kubectl get webpage example-webpage -o yaml
```

## Understanding the Components

### 1. Custom Resource Definition (CRD)

The CRD extends Kubernetes with a new resource type. When you apply the CRD YAML:
- Kubernetes API server learns about the `WebPage` type
- You can now use `kubectl` to manage WebPage resources
- The API server validates WebPage resources against the schema

### 2. Controller Manager

The Manager:
- Connects to the Kubernetes API server
- Maintains a cache of resources for efficient reads
- Runs multiple controllers
- Handles graceful shutdown
- Manages leader election for high availability

### 3. Reconciler

The Reconciler:
- Watches for changes to your custom resources
- Gets called for every create/update/delete event
- Implements idempotent logic (can be called multiple times safely)
- Updates resource status to reflect actual state
- Can create/update/delete other Kubernetes resources

### 4. Kubebuilder Markers

Special comments that control code generation:

| Marker | Purpose |
|--------|---------|
| `+kubebuilder:object:root=true` | Mark as top-level API type |
| `+kubebuilder:subresource:status` | Enable status subresource |
| `+kubebuilder:resource:scope=Namespaced` | Resource is namespace-scoped |
| `+kubebuilder:printcolumn` | Add column to `kubectl get` output |
| `+kubebuilder:validation:Required` | Mark field as required |
| `+kubebuilder:validation:Optional` | Mark field as optional |
| `+kubebuilder:validation:Minimum=1` | Set minimum value |

### 5. client.Client Interface

The Client provides CRUD operations:

```go
// Create
err := r.Client.Create(ctx, &deployment)

// Get
err := r.Client.Get(ctx, types.NamespacedName{Name: "name", Namespace: "ns"}, &obj)

// Update
err := r.Client.Update(ctx, &obj)

// Delete
err := r.Client.Delete(ctx, &obj)

// List
list := &v1alpha1.WebPageList{}
err := r.Client.List(ctx, list, client.InNamespace("default"))

// Patch
err := r.Client.Patch(ctx, &obj, patch)

// Update Status
err := r.Client.Status().Update(ctx, &obj)
```

### 6. Reconciliation Patterns

**Basic Pattern**:
```go
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    // 1. Fetch the resource
    // 2. Handle deletion (check for deletion timestamp)
    // 3. Create/update child resources
    // 4. Update status
    // 5. Return result
}
```

**With Finalizers** (for cleanup on deletion):
```go
const finalizerName = "webpage.finalizers.data-platform.qonto.co"

// If being deleted and has finalizer
if !obj.DeletionTimestamp.IsZero() && controllerutil.ContainsFinalizer(obj, finalizerName) {
    // Perform cleanup
    // ...
    // Remove finalizer
    controllerutil.RemoveFinalizer(obj, finalizerName)
    return ctrl.Result{}, r.Update(ctx, obj)
}

// Add finalizer if not present
if !controllerutil.ContainsFinalizer(obj, finalizerName) {
    controllerutil.AddFinalizer(obj, finalizerName)
    return ctrl.Result{}, r.Update(ctx, obj)
}
```

## Next Steps

To extend this operator into a fully functional web server manager:

1. **Implement Full Reconciliation Logic**:
   - Create a Deployment based on the WebPage spec
   - Create a Service to expose the Deployment
   - Create a ConfigMap with the content
   - Update WebPage status with phase and conditions

2. **Add More Fields**:
   ```go
   type WebPageSpec struct {
       Content   string            `json:"content"`
       Image     string            `json:"image"`
       Replicas  int               `json:"replicas"`
       Port      int               `json:"port"`             // New
       Domain    string            `json:"domain"`           // New
       TLS       bool              `json:"tls"`              // New
       Resources ResourceRequirements `json:"resources"`    // New
   }
   ```

3. **Add Status Conditions**:
   ```go
   type WebPageStatus struct {
       Phase      WebPhase           `json:"phase"`
       Conditions []metav1.Condition `json:"conditions,omitempty"`
       URL        string             `json:"url,omitempty"`
   }
   ```

4. **Watch Child Resources**:
   ```go
   ctrl.NewControllerManagedBy(mgr).
       For(&v1alpha1.WebPage{}).
       Owns(&appsv1.Deployment{}).
       Owns(&corev1.Service{}).
       Complete(r)
   ```

5. **Add Validation Webhooks**:
   - Validate spec fields before allowing creation
   - Prevent invalid updates

6. **Add Conversion Webhooks**:
   - Support multiple API versions (v1alpha1, v1beta1, v1)

7. **Add Metrics and Health Checks**:
   - Expose Prometheus metrics
   - Add health and ready endpoints

8. **Testing**:
   - Unit tests for reconciliation logic
   - Integration tests with envtest
   - E2E tests on real clusters

## Key Concepts Summary

- **Custom Resource**: Extends Kubernetes with new API types
- **Controller**: Watches resources and reconciles desired state with actual state
- **Reconciliation**: The process of making the actual state match the desired state
- **Manager**: Coordinates multiple controllers and shared dependencies
- **Scheme**: Registry mapping Go types to Kubernetes GroupVersionKinds
- **Client**: Interface for CRUD operations on Kubernetes resources
- **Watch**: Mechanism to get notified about resource changes
- **Operator Pattern**: CRD + Controller = Operator

## Resources

- [Kubebuilder Book](https://book.kubebuilder.io/)
- [controller-runtime Documentation](https://pkg.go.dev/sigs.k8s.io/controller-runtime)
- [Kubernetes API Conventions](https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md)
- [Operator Pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)

## License

Copyright 2025.

Licensed under the Apache License, Version 2.0.

