package hub

// received value from server
type RV struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// client data send to server
type CD struct {
	Data []byte
}
