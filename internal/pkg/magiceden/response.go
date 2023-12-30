package magiceden

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

type FloorResponse struct {
	Symbol       string  `json:"symbol"`
	FloorPrice   float64 `json:"floorPrice"`
	ListedCount  int     `json:"listedCount"`
	AvgPrice24Hr float64 `json:"avgPrice24hr"`
	VolumeAll    float64 `json:"volumeAll"`
}
