package model

type WsMessage struct {
	Type string      `json:"type" v:"required"`
	Data interface{} `json:"data"`
	Msg  interface{} `json:"msg"`
}

type RequestData struct {
	Namespace string            `json:"namespace"`
	NodeName  string            `json:"nodeName"`
	Selector  map[string]string `json:"selector"`
	PodName   string            `json:"podName"`
}
