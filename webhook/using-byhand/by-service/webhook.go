package main

import (
	"encoding/json" // 用于处理 JSON 格式的编码和解码  
	"fmt"
	"io/ioutil" // 用于读取输入输出流  
	"net/http"  // HTTP 客户端和服务器库  
	"strings"

	"github.com/golang/glog"                                      // 日志库
	admissionv1 "k8s.io/api/admission/v1"                         // Kubernetes Admission API
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1" // Admission Registration API
	appsv1 "k8s.io/api/apps/v1"                                   // Apps API
	corev1 "k8s.io/api/core/v1"                                   // Core API
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"                 // 元数据类型
	"k8s.io/apimachinery/pkg/runtime"                             // Kubernetes 运行时相关类型
	"k8s.io/apimachinery/pkg/runtime/serializer"                  // 编解码器相关
	"k8s.io/kubernetes/pkg/apis/core/v1"                          // Core API
)

var (
	// 创建一个运行时方案用于类型注册
	runtimeScheme = runtime.NewScheme()
	codecs        = serializer.NewCodecFactory(runtimeScheme) // 创建编解码器工厂
	deserializer  = codecs.UniversalDeserializer()            // 创建通用解码器

	// (https://github.com/kubernetes/kubernetes/issues/57982)
	defaulter = runtime.ObjectDefaulter(runtimeScheme) // 默认值处理器
)

var (
	// 要忽略的命名空间列表
	ignoredNamespaces = []string{
		metav1.NamespaceSystem, // 系统命名空间
		metav1.NamespacePublic, // 公共命名空间
	}

	// 需要的标签列表
	requiredLabels = []string{
		nameLabel,
		instanceLabel,
		versionLabel,
		componentLabel,
		partOfLabel,
		managedByLabel,
	}

	// 添加的标签映射
	addLabels = map[string]string{
		nameLabel:      NA,
		instanceLabel:  NA,
		versionLabel:   NA,
		componentLabel: NA,
		partOfLabel:    NA,
		managedByLabel: NA,
	}
)

const (
	// Admission Webhook 注解的键
	admissionWebhookAnnotationValidateKey = "admission-webhook-example.qikqiak.com/validate"
	admissionWebhookAnnotationMutateKey   = "admission-webhook-example.qikqiak.com/mutate"
	admissionWebhookAnnotationStatusKey   = "admission-webhook-example.qikqiak.com/status"

	// 常用标签
	nameLabel      = "app.kubernetes.io/name"
	instanceLabel  = "app.kubernetes.io/instance"
	versionLabel   = "app.kubernetes.io/version"
	componentLabel = "app.kubernetes.io/component"
	partOfLabel    = "app.kubernetes.io/part-of"
	managedByLabel = "app.kubernetes.io/managed-by"

	// 不可用状态值
	NA = "not_available"
)

// WebhookServer 结构体封装了 HTTP 服务器
type WebhookServer struct {
	server *http.Server // HTTP 服务器
}

// Webhook 服务器参数结构体
type WhSvrParameters struct {
	port     int    // webhook 服务器端口
	certFile string // HTTPS 的 x509 证书路径
	keyFile  string // 与 certFile 匹配的 x509 私钥路径
}

// 补丁操作结构体
type patchOperation struct {
	Op    string      `json:"op"`              // 操作类型（add、replace、remove）
	Path  string      `json:"path"`            // 操作的路径
	Value interface{} `json:"value,omitempty"` // 操作的值
}

// 初始化函数
func init() {
	_ = corev1.AddToScheme(runtimeScheme)                  // 将 Core API 注册到运行时方案
	_ = admissionregistrationv1.AddToScheme(runtimeScheme) // 将 Admission Registration API 注册到运行时方案
	_ = v1.AddToScheme(runtimeScheme)                      // 将 Apps API 注册到运行时方案
}

// 校验是否需要 Admission
func admissionRequired(ignoredList []string, admissionAnnotationKey string, metadata *metav1.ObjectMeta) bool {
	// 跳过特殊的 Kubernetes 系统命名空间
	for _, namespace := range ignoredList {
		if metadata.Namespace == namespace {
			glog.Infof("Skip validation for %v for it's in special namespace:%v", metadata.Name, metadata.Namespace)
			return false
		}
	}

	annotations := metadata.GetAnnotations() // 获取元数据中的注解
	if annotations == nil {
		annotations = map[string]string{}
	}

	var required bool
	// 根据注解判断是否需要执行 Admission
	switch strings.ToLower(annotations[admissionAnnotationKey]) {
	default:
		required = true // 默认需要
	case "n", "no", "false", "off":
		required = false // 如果注解值为这些，则不需要
	}
	return required
}

// 检查是否需要变更
func mutationRequired(ignoredList []string, metadata *metav1.ObjectMeta) bool {
	required := admissionRequired(ignoredList, admissionWebhookAnnotationMutateKey, metadata) // 检查是否需要变更
	annotations := metadata.GetAnnotations()
	if annotations == nil {
		annotations = map[string]string{}
	}
	// 检查状态注解
	status := annotations[admissionWebhookAnnotationStatusKey]

	if strings.ToLower(status) == "mutated" {
		required = false // 如果已变更，则不需要再变更
	}

	glog.Infof("Mutation policy for %v/%v: required:%v", metadata.Namespace, metadata.Name, required)
	return required
}

// 检查是否需要验证
func validationRequired(ignoredList []string, metadata *metav1.ObjectMeta) bool {
	required := admissionRequired(ignoredList, admissionWebhookAnnotationValidateKey, metadata) // 检查验证需求
	glog.Infof("Validation policy for %v/%v: required:%v", metadata.Namespace, metadata.Name, required)
	return required
}

// 更新注解并创建补丁
func updateAnnotation(target map[string]string, added map[string]string) (patch []patchOperation) {
	for key, value := range added {
		if target == nil || target[key] == "" {
			target = map[string]string{} // 如果目标为空，则初始化
			patch = append(patch, patchOperation{
				Op:   "add", // 添加操作
				Path: "/metadata/annotations",
				Value: map[string]string{
					key: value, // 添加的新注解
				},
			})
		} else {
			patch = append(patch, patchOperation{
				Op:    "replace", // 替换操作
				Path:  "/metadata/annotations/" + key,
				Value: value,
			})
		}
	}
	return patch
}

// 更新标签并创建补丁
func updateLabels(target map[string]string, added map[string]string) (patch []patchOperation) {
	values := make(map[string]string)
	for key, value := range added {
		if target == nil || target[key] == "" {
			values[key] = value // 仅添加未设置的标签
		}
	}
	patch = append(patch, patchOperation{
		Op:    "add", // 添加操作
		Path:  "/metadata/labels",
		Value: values,
	})
	return patch
}

// 创建补丁 JSON
func createPatch(availableAnnotations map[string]string, annotations map[string]string, availableLabels map[string]string, labels map[string]string) ([]byte, error) {
	var patch []patchOperation

	// 更新注解和标签
	patch = append(patch, updateAnnotation(availableAnnotations, annotations)...)
	patch = append(patch, updateLabels(availableLabels, labels)...)

	return json.Marshal(patch) // 返回补丁的 JSON 编码
}

// 验证 Deployment 和 Service
func (whsvr *WebhookServer) validate(ar *admissionv1.AdmissionReview) *admissionv1.AdmissionResponse {
	req := ar.Request
	var (
		availableLabels                 map[string]string  // 可用标签
		objectMeta                      *metav1.ObjectMeta // 对象元数据
		resourceNamespace, resourceName string             // 资源的命名空间和名称
	)

	glog.Infof("AdmissionReview for Kind=%v, Namespace=%v Name=%v (%v) UID=%v patchOperation=%v UserInfo=%v",
		req.Kind, req.Namespace, req.Name, resourceName, req.UID, req.Operation, req.UserInfo)

	// 根据 Kind 解码请求对象
	switch req.Kind.Kind {
	case "Deployment":
		var deployment appsv1.Deployment
		if err := json.Unmarshal(req.Object.Raw, &deployment); err != nil {
			glog.Errorf("Could not unmarshal raw object: %v", err)
			return &admissionv1.AdmissionResponse{
				Result: &metav1.Status{
					Message: err.Error(), // 返回错误信息
				},
			}
		}
		resourceName, resourceNamespace, objectMeta = deployment.Name, deployment.Namespace, &deployment.ObjectMeta
		availableLabels = deployment.Labels // 获取可用标签
	case "Service":
		var service corev1.Service
		if err := json.Unmarshal(req.Object.Raw, &service); err != nil {
			glog.Errorf("Could not unmarshal raw object: %v", err)
			return &admissionv1.AdmissionResponse{
				Result: &metav1.Status{
					Message: err.Error(), // 返回错误信息
				},
			}
		}
		resourceName, resourceNamespace, objectMeta = service.Name, service.Namespace, &service.ObjectMeta
		availableLabels = service.Labels // 获取可用标签
	}

	// 判断是否需要验证
	if !validationRequired(ignoredNamespaces, objectMeta) {
		glog.Infof("Skipping validation for %s/%s due to policy check", resourceNamespace, resourceName)
		return &admissionv1.AdmissionResponse{
			Allowed: true, // 允许请求通过
		}
	}

	allowed := true
	var result *metav1.Status
	glog.Info("available labels:", availableLabels) // 输出可用标签
	glog.Info("required labels", requiredLabels)    // 输出所需标签
	// 检查是否所有必需标签都存在
	for _, rl := range requiredLabels {
		if _, ok := availableLabels[rl]; !ok {
			allowed = false
			result = &metav1.Status{
				Reason: "required labels are not set", // 标签未设置
			}
			break
		}
	}

	return &admissionv1.AdmissionResponse{
		Allowed: allowed,
		Result:  result, // 验证结果
	}
}

// 主要变更处理过程
func (whsvr *WebhookServer) mutate(ar *admissionv1.AdmissionReview) *admissionv1.AdmissionResponse {
	req := ar.Request
	var (
		availableLabels, availableAnnotations map[string]string  // 可用标签和注解
		objectMeta                            *metav1.ObjectMeta // 对象元数据
		resourceNamespace, resourceName       string             // 资源的命名空间和名称
	)

	glog.Infof("AdmissionReview for Kind=%v, Namespace=%v Name=%v (%v) UID=%v patchOperation=%v UserInfo=%v",
		req.Kind, req.Namespace, req.Name, resourceName, req.UID, req.Operation, req.UserInfo)

	// 根据 Kind 解码请求对象
	switch req.Kind.Kind {
	case "Deployment":
		var deployment appsv1.Deployment
		if err := json.Unmarshal(req.Object.Raw, &deployment); err != nil {
			glog.Errorf("Could not unmarshal raw object: %v", err)
			return &admissionv1.AdmissionResponse{
				Result: &metav1.Status{
					Message: err.Error(), // 返回错误信息
				},
			}
		}
		resourceName, resourceNamespace, objectMeta = deployment.Name, deployment.Namespace, &deployment.ObjectMeta
		availableLabels = deployment.Labels // 获取可用标签
	case "Service":
		var service corev1.Service
		if err := json.Unmarshal(req.Object.Raw, &service); err != nil {
			glog.Errorf("Could not unmarshal raw object: %v", err)
			return &admissionv1.AdmissionResponse{
				Result: &metav1.Status{
					Message: err.Error(), // 返回错误信息
				},
			}
		}
		resourceName, resourceNamespace, objectMeta = service.Name, service.Namespace, &service.ObjectMeta
		availableLabels = service.Labels // 获取可用标签
	}

	// 判断是否需要变更
	if !mutationRequired(ignoredNamespaces, objectMeta) {
		glog.Infof("Skipping mutation for %s/%s due to policy check", resourceNamespace, resourceName)
		return &admissionv1.AdmissionResponse{
			Allowed: true, // 允许请求通过
		}
	}

	// 设置状态注解
	annotations := map[string]string{admissionWebhookAnnotationStatusKey: "mutated"}
	patchBytes, err := createPatch(availableAnnotations, annotations, availableLabels, addLabels) // 创建补丁
	if err != nil {
		return &admissionv1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(), // 返回错误信息
			},
		}
	}

	glog.Infof("AdmissionResponse: patch=%v\n", string(patchBytes))
	return &admissionv1.AdmissionResponse{
		Allowed: true,       // 允许请求通过
		Patch:   patchBytes, // 返回补丁
		PatchType: func() *admissionv1.PatchType {
			pt := admissionv1.PatchTypeJSONPatch
			return &pt
		}(),
	}
}

// Webhook 服务器的处理方法
func (whsvr *WebhookServer) serve(w http.ResponseWriter, r *http.Request) {
	var body []byte
	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			body = data // 读取请求体
		}
	}
	if len(body) == 0 {
		glog.Error("empty body")
		http.Error(w, "empty body", http.StatusBadRequest) // 返回错误
		return
	}

	// 验证 Content-Type 是否正确
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		glog.Errorf("Content-Type=%s, expect application/json", contentType)
		http.Error(w, "invalid Content-Type, expect `application/json`", http.StatusUnsupportedMediaType) // 返回错误
		return
	}

	var admissionResponse *admissionv1.AdmissionResponse
	ar := admissionv1.AdmissionReview{}
	// 解码请求体
	if _, _, err := deserializer.Decode(body, nil, &ar); err != nil {
		glog.Errorf("Can't decode body: %v", err)
		admissionResponse = &admissionv1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(), // 返回错误信息
			},
		}
	} else {
		fmt.Println(r.URL.Path)
		// 根据请求的路径调用相应的处理函数
		if r.URL.Path == "/mutate" {
			admissionResponse = whsvr.mutate(&ar) // 处理变更请求
		} else if r.URL.Path == "/validate" {
			admissionResponse = whsvr.validate(&ar) // 处理验证请求
		}
	}

	// 创建响应对象
	admissionReview := admissionv1.AdmissionReview{
		// 设置 API 元信息
		TypeMeta: metav1.TypeMeta{
			APIVersion: "admission.k8s.io/v1",
			Kind:       "AdmissionReview",
		},
	}
	if admissionResponse != nil {
		admissionReview.Response = admissionResponse
		if ar.Request != nil {
			admissionReview.Response.UID = ar.Request.UID // 复制请求的 UID
		}
	}

	// 编码响应
	resp, err := json.Marshal(admissionReview)
	if err != nil {
		glog.Errorf("Can't encode response: %v", err)
		http.Error(w, fmt.Sprintf("could not encode response: %v", err), http.StatusInternalServerError)
	}
	glog.Infof("Ready to write response ...")
	if _, err := w.Write(resp); err != nil {
		glog.Errorf("Can't write response: %v", err)
		http.Error(w, fmt.Sprintf("could not write response: %v", err), http.StatusInternalServerError)
	}
}
