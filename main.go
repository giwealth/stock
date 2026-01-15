package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"stock/pkg/shares"
)

var (
	code     string
	detail   bool
	filename string
)

var (
	sh = "/mnt/d/new_tdx/vipdoc/sh/lday/"
	sz = "/mnt/d/new_tdx/vipdoc/sz/lday/"
)

func main() {
	shfiles := list(sh)
	szfiles := list(sz)

	var files []string
	for _, v := range shfiles {
		s := filepath.Base(v)
		if string(s[0:4]) == "sh60" {
			files = append(files, v)
		}
	}
	for _, v := range szfiles {
		s := filepath.Base(v)
		if string(s[0:4]) == "sz00" {
			files = append(files, v)
		}
	}

	for _, file := range files {
		gp, err := shares.ParseFile(file)
		if err != nil {
			continue
		}

		if gp == nil {
			continue
		}

		risk := IsMacdTopMountainRisk(gp)
		if !risk {
			fmt.Println(gp.Code)
		}

	}
}

func list(path string) []string {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatalln(err)
	}
	var fileList []string
	for _, file := range files {
		if file.IsDir() {
			fileList = append(fileList, list(filepath.Join(path, file.Name()))...)
		} else {
			fileList = append(fileList, filepath.Join(path, file.Name()))
		}
	}
	return fileList
}

func IsMacdTopMountainRisk(stock *shares.Stock) bool {
	data := stock.Data
	n := len(data)
	if n < 10 {
		return false
	}

	// 1. 找最近一个已完成的 MACD 顶部山包
	peak, start, end, ok := findLastMacdTopMountain(data)
	if !ok {
		return false
	}

	i := n - 1 // 当前交易日

	// 2. 当前 MACD(Hist) > 前一日 MACD(Hist)
	if data[i].Hist <= data[i-1].Hist {
		return false
	}

	// ❗ 防止强趋势（当前动能不能超过前山包峰值）
	if data[i].Hist >= data[peak].Hist {
		return false
	}

	// 3. 当前最高价 > 前一山包内最高价
	prevHigh := data[start].High
	for j := start + 1; j <= end; j++ {
		if data[j].High > prevHigh {
			prevHigh = data[j].High
		}
	}

	if data[i].High <= prevHigh {
		return false
	}

	return true
}

func findLastMacdTopMountain(data []shares.Day) (peak, start, end int, ok bool) {
	n := len(data)

	// 从后往前找最近的山包
	for i := n - 3; i >= 2; i-- {

		// 必须在 0 轴上方
		if data[i].Hist <= 0 {
			continue
		}

		// 局部最大值（山顶）
		if data[i-2].Hist < data[i-1].Hist &&
			data[i-1].Hist < data[i].Hist &&
			data[i].Hist > data[i+1].Hist &&
			data[i+1].Hist > data[i+2].Hist {

			peak = i

			// 向左找山包起点（递增开始）
			start = i - 1
			for start > 0 &&
				data[start-1].Hist > 0 &&
				data[start-1].Hist < data[start].Hist {
				start--
			}

			// 向右找山包终点（递减结束）
			end = i + 1
			for end < n-1 &&
				data[end].Hist > 0 &&
				data[end].Hist < data[end-1].Hist {
				end++
			}

			// 山包必须“完成”（明显衰减）
			if end < n-1 && data[end].Hist < data[peak].Hist*0.3 {
				ok = true
				return
			}
		}
	}

	return
}
