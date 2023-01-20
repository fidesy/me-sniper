package models

type BuyResponse struct {
	Tx struct {
		Type string `json:"type"`
		Data []byte `json:"data"`
	} `json:"tx"`
	TxSigned struct {
		Type string `json:"type"`
		Data []byte `json:"data"`
	} `json:"txSigned"`
}
