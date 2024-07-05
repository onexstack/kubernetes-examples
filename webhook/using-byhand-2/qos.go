package main

type QoS struct {
	Spec QoSpec `json:"spec,omitempty"`
}

type QoSpec struct {
	Cpu    int64 `json:"cpu,omitempty"`
	Memory int64 `json:"memory,omitempty"`
}

var (
	QoSInst = &QoSpec{}
)

func (qos *QoSpec) getQoSpec() *QoSpec {
	//判断空对象
	if (*QoSInst == QoSpec{}) {
		//默认值
		return &QoSpec{
			Cpu:    200,
			Memory: 400,
		}
	} else {
		return QoSInst
	}
}

func (qos *QoSpec) setQoSpec(qoSpec *QoSpec) {
	QoSInst = qoSpec
}
