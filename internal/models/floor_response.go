package models

type FloorResponse struct {
	Symbol       string  `json:"symbol"`
	FloorPrice   float64 `json:"floorPrice"`
	ListedCount  int     `json:"listedCount"`
	AvgPrice24Hr float64 `json:"avgPrice24hr"`
	VolumeAll    float64 `json:"volumeAll"`
}
