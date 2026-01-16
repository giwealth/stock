package index

func EMA(prices []float64, n int) []float64 {
	ema := make([]float64, len(prices))
	k := 2.0 / float64(n+1)
	ema[0] = prices[0]
	for i := 1; i < len(prices); i++ {
		ema[i] = prices[i]*k + ema[i-1]*(1-k)
	}
	return ema
}

// MACD 使用所有的收盘价进行计算  参数6, 30, 9或6, 30, 6  默认值为12, 26, 9
func MACD(prices []float64) (dif, dea, hist []float64) {
	ema12 := EMA(prices, 12)
	ema26 := EMA(prices, 26)
	dif = make([]float64, len(prices))
	dea = make([]float64, len(prices))
	hist = make([]float64, len(prices))
	for i := range prices {
		if i < 26 {
			dif[i] = 0
			dea[i] = 0
			hist[i] = 0
		} else {
			dif[i] = ema12[i] - ema26[i]
			if i < 35 {
				dea[i] = 0
				hist[i] = dif[i]
			} else {
				dea[i] = EMA(dif[:i+1], 9)[i]
				hist[i] = 2 * (dif[i] - dea[i])
			}
		}
	}
	return dif, dea, hist
}
