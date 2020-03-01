package aftg

type Tag struct {
	Name string `json:"name"`
	TimestampBegin int64 `json:"timestampBegin"`
	TimestampEnd int64 `json:"timestampEnd"`
	ProductName string `json:"productName"`
	TagName string `json:"tagName"`
}

type NTP struct {
	SrvReceptionTime int64 `json:"srvReceptionTime"`
	ClientTransmissionTime int64 `json:"clientTransmissionTime"`
	SrvTransmissionTime int64 `json:"srvTransmissionTime"`
	ClientReceptionTime int64 `json:"clientReceptionTime"`
}

type RegisterApiKeyRequest struct {
	Key string `json:"key"`
}
