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
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	mathv1 "github.com/superproj/webhook-with-kubebuilder/api/v1"
)

// CalculateReconciler reconciles a Calculate object
type CalculateReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=math.superproj.com,resources=calculates,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=math.superproj.com,resources=calculates/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=math.superproj.com,resources=calculates/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Calculate object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.18.4/pkg/reconcile
func (r *CalculateReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var cal mathv1.Calculate
	if err := r.Get(ctx, req.NamespacedName, &cal); err != nil {
		klog.Error(err, "unable to fetch calculate")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	klog.Infof("Found the calculate object %v", cal)

	klog.Infof("Calculating the calculate of %d and %d with action %s", cal.Spec.First, cal.Spec.Second, cal.Spec.Action)
	switch cal.Spec.Action {
	case mathv1.ActionTypeAdd:
		cal.Status.Result = cal.Spec.First + cal.Spec.Second
	case mathv1.ActionTypeSub:
		cal.Status.Result = cal.Spec.First - cal.Spec.Second
	case mathv1.ActionTypeMul:
		cal.Status.Result = cal.Spec.First * cal.Spec.Second
	case mathv1.ActionTypeDiv:
		if cal.Spec.Second == 0 {
			return ctrl.Result{}, fmt.Errorf("the divisor cannot be zero whtn action is division")
		}
		cal.Status.Result = cal.Spec.First / cal.Spec.Second
	default:
		return ctrl.Result{}, fmt.Errorf("unknown action type")
	}

	klog.Info("Updating the result of calculation")
	if err := r.Status().Update(ctx, &cal); err != nil {
		klog.Error(err, "Unable to update calculate status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CalculateReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mathv1.Calculate{}).
		Complete(r)
}
