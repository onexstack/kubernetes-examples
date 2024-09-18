// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/onex.
//

// +k8s:deepcopy-gen=package
// +k8s:defaulter-gen=TypeMeta
// +k8s:defaulter-gen-input=github.com/superproj/k8sdemo/resourcedefinition/apps/v1beta1
// +k8s:conversion-gen=github.com/superproj/k8sdemo/resourcedefinition/apps
// +k8s:conversion-gen=k8s.io/kubernetes/pkg/apis/core
// +k8s:conversion-gen-external-types=github.com/superproj/k8sdemo/resourcedefinition/apps/v1beta1

// Package v1beta1 is the v1beta1 version of the API.
package v1beta1
