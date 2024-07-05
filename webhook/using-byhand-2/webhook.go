package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/ghodss/yaml"
	"github.com/golang/glog"
	"github.com/wI2L/jsondiff"
	"k8s.io/api/admission/v1beta1"
	admissionregistrationv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	v1 "k8s.io/kubernetes/pkg/apis/core/v1"
)

var (
	runtimeScheme = runtime.NewScheme()
	codecs        = serializer.NewCodecFactory(runtimeScheme)
	deserializer  = codecs.UniversalDeserializer()

	// (https://github.com/kubernetes/kubernetes/issues/57982)
	defaulter = runtime.ObjectDefaulter(runtimeScheme)

	logFinalYaml  = true
	logFinalPatch = true
)

var (
	ignoredNamespaces = []string{
		metav1.NamespaceSystem,
		metav1.NamespacePublic,
	}
)

const (
	admissionWebhookAnnotationMutateKey = "admission-webhook-example.qikqiak.com/mutate"
)

type WebhookServer struct {
	server *http.Server
}

// Webhook Server parameters
type WhSvrParameters struct {
	port           int    // webhook server port
	certFile       string // path to the x509 certificate for https
	keyFile        string // path to the x509 private key matching `CertFile`
	sidecarCfgFile string // path to sidecar injector configuration file
}

func init() {
	_ = corev1.AddToScheme(runtimeScheme)
	_ = admissionregistrationv1beta1.AddToScheme(runtimeScheme)
	// defaulting with webhooks:
	// https://github.com/kubernetes/kubernetes/issues/57982
	_ = v1.AddToScheme(runtimeScheme)
}

func mutationRequired(ignoredList []string, metadata *metav1.ObjectMeta, log *bytes.Buffer) (required bool) {

	// K8S 默认的命名空间不做修改
	for _, namespace := range ignoredList {
		if metadata.Namespace == namespace {
			log.WriteString(fmt.Sprintf("\nSkip validation for [%v] for it's in special namespace: [%v]", metadata.Name, metadata.Namespace))
			required = false
			return required
		}
	}

	//拿到annotations
	annotations := metadata.GetAnnotations()
	if annotations == nil {
		annotations = map[string]string{}
	}

	//如果是没有开启ccnp.cib.io/runtimeEnable为真的话不做修改
	log.WriteString(fmt.Sprintf("\nadmission-webhook-example.qikqiak.com/mutate: [%v]", annotations[admissionWebhookAnnotationMutateKey]))
	switch strings.ToLower(annotations[admissionWebhookAnnotationMutateKey]) {
	default:
		required = true
	case "n", "no", "false", "off", "": //注意如果是空(没有配置这个注解ccnp.cib.io/runtimeEnable，也不会执行修改)
		required = false
	}

	log.WriteString(fmt.Sprintf("\nMutation policy for [%v/%v]: required:%v", metadata.Namespace, metadata.Name, required))

	return required
}

//修改Deployment
func mutateDeploy(deploy *appsv1.Deployment, log *bytes.Buffer) *v1beta1.AdmissionResponse {

	var (
		objectMeta                      *metav1.ObjectMeta
		resourceNamespace, resourceName string
	)
	resourceName, resourceNamespace, objectMeta = deploy.Name, deploy.Namespace, &deploy.ObjectMeta

	log.WriteString("\n---------begin mumate check---------")
	//检查本次是否需要修改
	if !mutationRequired(ignoredNamespaces, objectMeta, log) {
		log.WriteString(fmt.Sprintf("\nSkipping mutation for [%s/%s] due to policy check", resourceNamespace, resourceName))
		return &v1beta1.AdmissionResponse{
			Allowed: true,
		}
	}
	log.WriteString("\n---------ended mumate check---------")

	/********************************************************* 进行修改操作 */
	log.WriteString("\n---------begin mumate---------")
	newDeploy := deploy.DeepCopy()
	newPodSpec := &newDeploy.Spec.Template.Spec

	//添加一个initContainer
	var initContainer *corev1.Container

	//如果没有initContainer则新加一个
	if len(newPodSpec.InitContainers) == 0 {
		newPodSpec.InitContainers = []corev1.Container{
			{
				Name:    "init",
				Image:   "busybox",
				Command: []string{"/bin/sh", "-c", " echo 'init' && sleep 100 "},
			},
		}
		initContainer = &newPodSpec.InitContainers[0]
		log.WriteString("\nmutate add initContainer sucess!")
		//资源限制
		initConRequest := make(map[corev1.ResourceName]resource.Quantity)
		initConRequest[corev1.ResourceCPU] = *resource.NewMilliQuantity(QoSInst.getQoSpec().Cpu, resource.DecimalSI)           //cpu资源限制 100m
		initConRequest[corev1.ResourceMemory] = *resource.NewQuantity(QoSInst.getQoSpec().Memory*1024*1024, resource.BinarySI) //内存资源限制 100Mi
		initContainer.Resources.Requests = initConRequest
	}

	/********************************************************* 结束修改操作 */
	log.WriteString("\n---------ended mumate---------")

	//打印修改后的yaml
	if logFinalYaml {
		log.WriteString("\n---------begin mumated yaml---------")
		bytes, err := json.Marshal(newPodSpec)
		if err == nil {
			yamlStr, err := yaml.JSONToYAML(bytes)
			if err == nil {
				log.WriteString("\n" + string(yamlStr))
			}
		}
		log.WriteString("\n---------ended mumated yaml---------")
	}

	// 比较新旧deploy的不同，返回不同的bytes
	patch, err := jsondiff.Compare(deploy, newDeploy)
	if err != nil {
		log.WriteString(fmt.Sprintf("\nPatch Compare process error: %v", err.Error()))
		return &v1beta1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	}

	//打patch
	patchBytes, err := json.MarshalIndent(patch, "", "    ")
	if err != nil {
		log.WriteString(fmt.Sprintf("\nPatch process error: %v", err.Error()))
		return &v1beta1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	}

	if logFinalPatch {
		log.WriteString(fmt.Sprintf("\nAdmissionResponse: patch=%v\n", string(patchBytes)))
	}
	return &v1beta1.AdmissionResponse{
		Allowed: true,
		Patch:   patchBytes,
		PatchType: func() *v1beta1.PatchType {
			pt := v1beta1.PatchTypeJSONPatch
			return &pt
		}(),
	}
}

//设置QoS
func mutateQoS(qos *QoS, operation v1beta1.Operation, log *bytes.Buffer) *v1beta1.AdmissionResponse {
	if operation == "DELETE" {
		//删除设置空对象
		QoSInst.setQoSpec(&QoSpec{})
		log.WriteString(fmt.Sprintf("\ndelete QoS from crd : [%v] ,reset to default [%v]", qos.Spec, qoSpecFromCRD.getQoSpec()))
	} else {
		if qos != nil && (qos.Spec != QoSpec{}) {
			QoSInst.setQoSpec(&qos.Spec)
			log.WriteString(fmt.Sprintf("\nget QoS value : [%v]", qos.Spec))
		} else {
			log.WriteString(fmt.Sprintf("\nget an empty QoS value : [%v]", qos.Spec))
		}
	}

	return &v1beta1.AdmissionResponse{
		Allowed: true,
	}
}

// main mutation process
func (whsvr *WebhookServer) mutate(ar *v1beta1.AdmissionReview, log *bytes.Buffer) *v1beta1.AdmissionResponse {
	req := ar.Request

	log.WriteString(fmt.Sprintf("\n======begin Admission for Namespace=[%v], Kind=[%v], Name=[%v]======", req.Namespace, req.Kind.Kind, req.Name))
	log.WriteString("\n>>>>>>" + req.Kind.Kind)
	switch req.Kind.Kind {
	case "Deployment":
		var deployment appsv1.Deployment
		if err := json.Unmarshal(req.Object.Raw, &deployment); err != nil {
			log.WriteString(fmt.Sprintf("\nCould not unmarshal raw object: %v", err))
			glog.Errorf(log.String())
			return &v1beta1.AdmissionResponse{
				Result: &metav1.Status{
					Message: err.Error(),
				},
			}
		}
		return mutateDeploy(&deployment, log)
	case "QoS":
		var qos QoS
		var raw []byte
		if req.Operation == "DELETE" {
			raw = req.OldObject.Raw
		} else {
			raw = req.Object.Raw
		}
		if err := json.Unmarshal(raw, &qos); err != nil {
			log.WriteString(fmt.Sprintf("\nCould not unmarshal raw object: %v", err))
			glog.Errorf(log.String())
			return &v1beta1.AdmissionResponse{
				Result: &metav1.Status{
					Message: err.Error(),
				},
			}
		}
		return mutateQoS(&qos, req.Operation, log)
	default:
		msg := fmt.Sprintf("\nNot support for this Kind of resource  %v", req.Kind.Kind)
		log.WriteString(msg)
		return &v1beta1.AdmissionResponse{
			Result: &metav1.Status{
				Message: msg,
			},
		}
	}
}

// Serve method for webhook server
func (whsvr *WebhookServer) serve(w http.ResponseWriter, r *http.Request) {
	//记录日志
	var log bytes.Buffer

	//读取从ApiServer过来的数据放到body
	var body []byte
	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			body = data
		}
	}
	if len(body) == 0 {
		log.WriteString("empty body")
		glog.Info(log.String())
		//返回状态码400
		//如果在Apiserver调用此Webhook返回是400，说明APIServer自己传过来的数据是空
		http.Error(w, log.String(), http.StatusBadRequest)
		return
	}

	// verify the content type is accurate
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		log.WriteString(fmt.Sprintf("Content-Type=%s, expect `application/json`", contentType))
		glog.Errorf(log.String())
		//如果在Apiserver调用此Webhook返回是415，说明APIServer自己传过来的数据不是json格式，处理不了
		http.Error(w, log.String(), http.StatusUnsupportedMediaType)
		return
	}

	var admissionResponse *v1beta1.AdmissionResponse
	ar := v1beta1.AdmissionReview{}
	if _, _, err := deserializer.Decode(body, nil, &ar); err != nil {
		//组装错误信息
		log.WriteString(fmt.Sprintf("\nCan't decode body,error info is :  %s", err.Error()))
		glog.Errorln(log.String())
		//返回错误信息，形式表现为资源创建会失败，
		admissionResponse = &v1beta1.AdmissionResponse{
			Result: &metav1.Status{
				Message: log.String(),
			},
		}
	} else {
		fmt.Println(r.URL.Path)
		if r.URL.Path == "/mutate" {
			admissionResponse = whsvr.mutate(&ar, &log)
		}
	}

	admissionReview := v1beta1.AdmissionReview{}
	if admissionResponse != nil {
		admissionReview.Response = admissionResponse
		if ar.Request != nil {
			admissionReview.Response.UID = ar.Request.UID
		}
	}

	resp, err := json.Marshal(admissionReview)
	if err != nil {
		log.WriteString(fmt.Sprintf("\nCan't encode response: %v", err))
		http.Error(w, log.String(), http.StatusInternalServerError)
	}
	glog.Infof("Ready to write reponse ...")
	if _, err := w.Write(resp); err != nil {
		log.WriteString(fmt.Sprintf("\nCan't write response: %v", err))
		http.Error(w, log.String(), http.StatusInternalServerError)
	}

	log.WriteString("\n======ended Admission already writed to reponse======")
	//东八区时间
	datetime := time.Now().In(time.FixedZone("GMT", 8*3600)).Format("2006-01-02 15:04:05")
	//最后打印日志
	glog.Infof(datetime + " " + log.String())
}
