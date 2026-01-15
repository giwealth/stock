package shares

type (
	Stock struct {
		Code string `json:"code"`
		Data []Day  `json:"data"`
	}

	Day struct {
		Date   string  `json:"date"`
		Open   float64 `json:"open"`
		High   float64 `json:"high"`
		Low    float64 `json:"low"`
		Close  float64 `json:"close"`
		Amount float64 `json:"amount"`
		Vol    float64 `json:"vol"`
		Dif    float64 `json:"dif"`
		Dea    float64 `json:"dea"`
		Hist   float64 `json:"hist"`
		MA5    float64 `json:"ma5"`
		MA10   float64 `json:"ma10"`
		MA20   float64 `json:"ma20"`
		MA30   float64 `json:"ma30"`
		MA60   float64 `json:"ma60"`
		MA90   float64 `json:"ma90"`
		K float64 `json:"k"`
		D float64 `json:"d"`
		J float64 `json:"j"`
	}
)
