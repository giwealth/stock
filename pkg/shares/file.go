package shares

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"

	"stock/pkg/index"
)

func ParseFile(filename string) (*Stock, error) {
	buf, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	name := filepath.Base(filename)
	code := fmt.Sprint(name[2 : len(name)-4])

	stock := Stock{
		Code: code,
	}
	// 解包并写入每行数据
	prices := []float64{}
	highs := []float64{}
	lows := []float64{}
	closes := []float64{}
	const recSize = 32
	recCount := len(buf) / recSize
	for i := 0; i < recCount; i++ {
		begin := i * recSize
		end := begin + recSize

		a := unpack(buf[begin:end])
		date := fmt.Sprint(a[0].(uint32))
		open := float64(a[1].(uint32)) / 100
		high := float64(a[2].(uint32)) / 100
		low := float64(a[3].(uint32)) / 100
		close := float64(a[4].(uint32)) / 100
		amount := a[5].(float32)
		vol := a[6].(uint32) / 100
		prices = append(prices, close)
		highs = append(highs, high)
		lows = append(lows, low)
		closes = append(closes, close)
		stock.Data = append(stock.Data, Day{Date: date, Open: open, High: high, Low: low, Close: close, Amount: float64(amount), Vol: float64(vol)})
	}
	// 数据不足
	if len(prices) < 60 {
		return nil, nil
	}
	if len(prices) < 90 {
		return nil, nil
	}
	dif, dea, hist := index.MACD(prices)
	ks, ds, js := index.CalculateKDJ(highs, lows, closes, 9)
	ma5 := index.SMA(prices, 5)
	ma10 := index.SMA(prices, 10)
	ma20 := index.SMA(prices, 20)
	ma30 := index.SMA(prices, 30)
	ma60 := index.SMA(prices, 60)
	ma90 := index.SMA(prices, 90)

	for i := range stock.Data {
		stock.Data[i].Dif = dif[i]
		stock.Data[i].Dea = dea[i]
		stock.Data[i].Hist = hist[i]
		stock.Data[i].MA5 = ma5[i]
		stock.Data[i].MA10 = ma10[i]
		stock.Data[i].MA20 = ma20[i]
		stock.Data[i].MA30 = ma30[i]
		stock.Data[i].MA60 = ma60[i]
		stock.Data[i].MA90 = ma90[i]
		stock.Data[i].K = ks[i]
		stock.Data[i].D = ds[i]
		stock.Data[i].J = js[i]
 	}

	return &stock, nil
}

func unpack(buf []byte) []interface{} {
	// 解析的格式
	formats := []byte("IIIIIfI")

	// 解析结果
	values := make([]interface{}, len(formats))

	// 遍历格式字符串解析
	const recSize = 4
	for i, format := range formats {
		begin := i * recSize
		end := begin + recSize
		switch format {
		case 'I':
			var v uint32
			binary.Read(bytes.NewReader(buf[begin:end]), binary.LittleEndian, &v)
			values[i] = v
		case 'f':
			var v float32
			binary.Read(bytes.NewReader(buf[begin:end]), binary.LittleEndian, &v)
			values[i] = v
		}
	}

	return values
}
