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

package v1

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var calculatelog = logf.Log.WithName("calculate-resource")

// SetupWebhookWithManager 用于设置 Controller Manager 以管理 Webhooks。
func (r *Calculate) SetupWebhookWithManager(mgr ctrl.Manager) error {
	// 通过调用 ctrl.NewWebhookManagedBy(mgr).For(r).Complete() 方法，为 Calculate 类型的对象创建了一个 Webhook。
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// 通过 +kubebuilder:webhook 注解，声明了一个 Webhook，指定其路径、是否进行变换、失败策略、以及操作等信息。
// `make manifests` 命令会根据该注解生成 MutatingWebhookConfiguration
// +kubebuilder:webhook:path=/mutate-math-superproj-com-v1-calculate,mutating=true,failurePolicy=fail,sideEffects=None,groups=math.superproj.com,resources=calculates,verbs=create;update,versions=v1,name=mcalculate.kb.io,admissionReviewVersions=v1

// 声明一个匿名变量，确保 Calculate 结构体实现了 webhook.Validator 接口。
var _ webhook.Defaulter = &Calculate{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Calculate) Default() {
	calculatelog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
}

// 通过 +kubebuilder:webhook 注解，声明了一个 Webhook，指定其路径、是否进行变换、失败策略、以及操作等信息。
// `make manifests` 命令会根据该注解生成 ValidatingWebhookConfiguration
// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: The 'path' attribute must follow a specific pattern and should not be modified directly here.
// Modifying the path for an invalid path can cause API server errors; failing to locate the webhook.
// +kubebuilder:webhook:path=/validate-math-superproj-com-v1-calculate,mutating=false,failurePolicy=fail,sideEffects=None,groups=math.superproj.com,resources=calculates,verbs=create;update,versions=v1,name=vcalculate.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Calculate{}

// ValidateCreate 实现了 webhook.Validator 接口中的 ValidateCreate() 方法，用于在创建对象时进行验证.
func (r *Calculate) ValidateCreate() (admission.Warnings, error) {
	calculatelog.Info("validate create", "name", r.Name)

	return r.validate()
}

// ValidateUpdate 实现了 webhook.Validator 接口中的 ValidateUpdate() 方法，用于在更新对象时进行验证.
func (r *Calculate) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	calculatelog.Info("validate update", "name", r.Name)

	return r.validate()
}

// ValidateDelete 实现了 webhook.Validator 接口中的 ValidateDelete() 方法，用于在删除对象时进行验证.
func (r *Calculate) ValidateDelete() (admission.Warnings, error) {
	calculatelog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil, nil
}

// validate 包含实际的验证逻辑.
func (r *Calculate) validate() (admission.Warnings, error) {
	// field.ErrorList 类型的变量，可以保存多个验证错误，我们可以进行所有参数的验证
	// 然后一次性返回验证结果，Kubernetes 基本都是一次性返回多个验证结果的
	allErrs := field.ErrorList{}

	specPath := field.NewPath("spec")
	// 当 Action 类型为 `div` 时，如果 `second` 字段值为 0，则拒绝，并返回错误信息
	if r.Spec.Action == ActionTypeDiv && r.Spec.Second == 0 {
		allErrs = append(allErrs, field.Invalid(specPath.Child("second"), r.Spec.Second, "the divisor cannot be zero whtn action is division"))
	}

	return nil, allErrs.ToAggregate()
}
