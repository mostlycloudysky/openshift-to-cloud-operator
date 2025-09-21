/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	demov1 "github.com/mostlycloudysky/openshift-to-cloud-operator/api/v1"
)

// HelloReconciler reconciles a Hello object
type HelloReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=demo.migrate.dev,resources=hellos,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=demo.migrate.dev,resources=hellos/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=demo.migrate.dev,resources=hellos/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Hello object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *HelloReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// 1. Fetch the Hello object
	var hello demov1.Hello
	if err := r.Get(ctx, req.NamespacedName, &hello); err != nil {
		// If the object no longer exists, ignore
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// 2. Compute the new Echo message
	newEcho := "You said: " + hello.Spec.Message

	// 3. Only update status if it has changed
	if hello.Status.Echo != newEcho {
		hello.Status.Echo = newEcho
		if err := r.Status().Update(ctx, &hello); err != nil {
			log.Error(err, "unable to update Hello status")
			return ctrl.Result{}, err
		}
		log.Info("Updated Hello status", "echo", newEcho)
	}

	// Done
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *HelloReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&demov1.Hello{}).
		Named("hello").
		Complete(r)
}
