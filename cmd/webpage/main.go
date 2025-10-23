package main

import (
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	k8sscehme "k8s.io/client-go/kubernetes/scheme"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"webpage/api/v1alpha1"
	"webpage/controllers/webpage"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(k8sscehme.AddToScheme(scheme))
	utilruntime.Must(v1alpha1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func main() {
	// create main logger
	logger := zap.New()
	ctrl.SetLogger(logger)
	log := ctrl.Log.WithName("main")
	log.Info("set up manager")

	// create manager
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme: scheme,
	})

	// create reconciler
	wr := webpage.WebReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
		Log:    log.WithName("web-reconciler"),
	}

	// create new controller
	err = wr.SetupWithManager(mgr)
	if err != nil {
		log.Error(err, "unable to create controller")
		os.Exit(1)
	}

	// start manager
	ctx := ctrl.SetupSignalHandler()
	if err = mgr.Start(ctx); err != nil {
		log.Error(err, "problem running manager")
		os.Exit(1)
	}
}
