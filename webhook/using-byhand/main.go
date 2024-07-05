package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	// 这里需要引入 Kubernetes 相关 Go 包
	// 引入 Kubernetes Admission API 的 v1 版本，用于处理 AdmissionReview 请求和响应
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	// 引入 Kubernetes 序列化器包，用于处理序列化和反序列化
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

func main() {
	// 注册 /mutate 路由，/mutate 路径的请求将由 mutate 函数处理
	http.HandleFunc("/mutate", mutate)
	// 启动 webhook 服务，指定证书路径，开启 TLS 认证，监听 9999 端口
	fmt.Println("Started mutating admission webhook server")
	panic(http.ListenAndServeTLS(":9999", "cert/server.crt", "cert/server.key", nil))
}

// JSON Patch 操作结构体
type patchOperation struct {
	// 操作类型（如 "add"、"replace"、"remove"）
	Op string `json:"op"`
	// 修改的路径
	Path string `json:"path"`
	// 设置的值，可以是任意类型
	Value interface{} `json:"value,omitempty"`
}

// 执行修改（Mutating）逻辑的函数
func mutate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received a request sent by the kube-apiserver")

	// 读取 HTTP 请求的 body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		// 如果读取 body 失败，返回 HTTP 500 错误
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 创建一个反序列化器，用于将 JSON 数据反序列化为 Go 结构体
	deserializer := serializer.NewCodecFactory(runtime.NewScheme()).UniversalDeserializer()
	ar := admissionv1.AdmissionReview{}
	// 将请求 body 解析为 AdmissionReview 结构体
	if _, _, err := deserializer.Decode(body, nil, &ar); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var pod corev1.Pod
	// 将 AdmissionReview 中的对象数据反序列化为 Pod 结构体
	if err := json.Unmarshal(ar.Request.Object.Raw, &pod); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 检查 Pod 是否有 Annotations，如果没有，初始化一个空的 Map
	if pod.ObjectMeta.Annotations == nil {
		pod.ObjectMeta.Annotations = make(map[string]string)
	}
	// 通过 Annotation 为 Pod 增加一个标识，表示矿机的机型
	pod.ObjectMeta.Annotations["apps.onex.io/miner-type"] = "S1.SMALL1"

	// 构造 Patch 对象，用于表示 JSON Patch 操作
	patch := []patchOperation{
		{
			Op:    "add",                      // 操作类型为 "add"
			Path:  "/metadata/annotations",    // 修改的路径为 Pod 的 Annotations
			Value: pod.ObjectMeta.Annotations, // 设置的值为增加的 Annotation
		},
	}
	// 将 Patch 对象序列化为 JSON 字节切片
	patchBytes, err := json.Marshal(patch)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// 构造 AdmissionReview 响应
	admissionReview := admissionv1.AdmissionReview{
		// 设置 API 元信息
		TypeMeta: metav1.TypeMeta{
			APIVersion: "admission.k8s.io/v1",
			Kind:       "AdmissionReview",
		},
		// 构造响应内容
		Response: &admissionv1.AdmissionResponse{
			UID:     ar.Request.UID, // 设置原请求的 UID
			Allowed: true,           // 允许该操作
			Patch:   patchBytes,     // 设置修改后的 Patch 内容
			PatchType: func() *admissionv1.PatchType { // 设置 Patch 操作的类型为 JSON Patch
				pt := admissionv1.PatchTypeJSONPatch
				return &pt
			}(),
		},
	}
	// 将 AdmissionReview 响应序列化为 JSON 字节切片
	resp, err := json.Marshal(admissionReview)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 将响应数据写入 HTTP 响应体
	if _, err := w.Write(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 打印日志，输出 Pod 的元数据
	fmt.Printf("%#vn\n", pod.ObjectMeta)
	// 打印日志，表示资源变更成功，包括命名空间和资源名称
	fmt.Printf("[%s/%s] Resource change succeeded\n", pod.ObjectMeta.Namespace, pod.ObjectMeta.Name)
}
