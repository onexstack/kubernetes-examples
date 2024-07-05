/*
Copyright 2024.

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
	"sigs.k8s.io/controller-runtime/pkg/log"

	calculatorv1 "github.com/superproj/k8sdemo/webhook/using-kubebuilder/api/v1"
)

// SumReconciler reconciles a Sum object
type SumReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=calculator.superproj.com,resources=sums,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=calculator.superproj.com,resources=sums/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=calculator.superproj.com,resources=sums/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Sum object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.18.2/pkg/reconcile
func (r *SumReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var sum calculatorv1.Sum
	if err := r.Get(ctx, req.NamespacedName, &sum); err != nil {
		klog.Error(err, "unable to fetch sum")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	klog.Infof("Found the sum object %v", sum)

	klog.Infof("Calculating the sum of %d and %d", sum.Spec.First, sum.Spec.Second)
	result := sum.Spec.First + sum.Spec.Second
	sum.Status.Result = result

	klog.Info("Updating the result of calculation")
	if err := r.Status().Update(ctx, &sum); err != nil {
		klog.Error(err, "Unable to update sum status")
		return ctrl.Result{}, err
	}

	klog.Infof("Successfully updated the sum status with result %d", result)
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SumReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&calculatorv1.Sum{}).
		Complete(r)
}
