package index

import (
	"math"
)

type KDJIndicator struct {
	K float64
	D float64
	J float64
}

// 计算 KDJ 指标
func CalculateKDJ(highs, lows, closes []float64, period int) ([]float64, []float64, []float64) {
	if len(highs) != len(lows) || len(lows) != len(closes) || len(closes) < period {
		panic("输入数据长度不一致或小于周期")
	}

	k := 50.0 // 初始值
	d := 50.0 // 初始值
	alpha := 1.0 / 3.0

	kValues := make([]float64, len(closes))
	dValues := make([]float64, len(closes))
	jValues := make([]float64, len(closes))

	for i := 0; i < len(closes); i++ {
		if i >= period-1 {
			// 计算最近 N 日的最高价和最低价
			highest := math.SmallestNonzeroFloat64
			lowest := math.MaxFloat64
			for j := i - period + 1; j <= i; j++ {
				if highs[j] > highest {
					highest = highs[j]
				}
				if lows[j] < lowest {
					lowest = lows[j]
				}
			}

			// 计算 RSV
			rsv := 0.0
			if highest != lowest {
				rsv = (closes[i] - lowest) / (highest - lowest) * 100
			}

			// 计算 K 值、D 值和 J 值
			k = (1-alpha)*k + alpha*rsv
			d = (1-alpha)*d + alpha*k
			j := 3*k - 2*d

			kValues[i] = k
			dValues[i] = d
			jValues[i] = j
		} else {
			// 初始阶段 K、D、J 为空值
			kValues[i] = math.NaN()
			dValues[i] = math.NaN()
			jValues[i] = math.NaN()
		}
	}

	return kValues, dValues, jValues
}