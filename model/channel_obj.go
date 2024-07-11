package model

type ChannelObj struct {
	Key    string
	Result []byte
	Status bool
	Error  error
}

type Resp struct {
	Result string `json:"result"`
}
