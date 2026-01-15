package index

// SMA 计算给定天数的简单移动平均线
func SMA(prices []float64, days int) []float64 {
    n := len(prices)
    sma := make([]float64, n)

    // 初始化前days-1天的SMA为0
    for i := 0; i < days-1; i++ {
        sma[i] = 0
    }

    // 计算SMA
    for i := days - 1; i < n; i++ {
        var sum float64
        for j := 0; j < days; j++ {
            sum += prices[i-j]
        }
        sma[i] = sum / float64(days)
    }

    return sma
}