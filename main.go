package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"stock/pkg/shares"
	"sync"
)

var (
	code     string
	detail   bool
	filename string

	currentDay = 4 // 1: 当前
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

	// 指定协程数量
	goroutineCount := 10

	// 创建任务 channel 和结果 channel
	fileChan := make(chan string, goroutineCount)
	var wg sync.WaitGroup
	var mu sync.Mutex // 保护 fmt.Println 输出

	// 启动 worker 协程
	for i := 0; i < goroutineCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file := range fileChan {
				gp, err := shares.ParseFile(file)
				if err != nil {
					continue
				}

				if gp == nil {
					continue
				}

				risk := IsMacdTopMountainRisk(gp)
				if risk {
					mu.Lock()
					fmt.Println(gp.Code, gp.Data[len(gp.Data)-currentDay].Date)
					mu.Unlock()
				}
			}
		}()
	}

	// 发送任务到 channel
	for _, file := range files {
		fileChan <- file
	}
	close(fileChan)

	// 等待所有协程完成
	wg.Wait()
}



func IsMacdTopMountainRisk(stock *shares.Stock) bool {
	data := stock.Data
	n := len(data)
	if n < 10 {
		return false
	}

	// 1. 找最近一个已完成的 MACD 顶部山包
	_, _, _, ok := findLastMacdTopMountain(data)
	if !ok {
		// fmt.Println(stock.Code, "false")
		return false
	}

	// i := n - currentDay // 当前交易日

	// // 2. 当前 MACD(Hist) > 前一日 MACD(Hist)
	// if data[i].Hist <= data[i-1].Hist {
	// 	return false
	// }
	
	// // ❗ 防止强趋势（当前动能不能超过前山包峰值）
	// if data[i].Hist >= data[peak].Hist {
	// 	return false
	// }

	// // 3. 当前最高价 > 前一山包内最高价
	// prevHigh := data[start].High
	// for j := start; j <= end; j++ {
	// 	if data[j].High > prevHigh {
	// 		prevHigh = data[j].High
	// 	}
	// }

	// if data[i].High <= prevHigh {
	// 	return false
	// }

	// // 4. 当前最高价 > 前一山包内最高价*1.05 假突破过滤（防强趋势）
	// if data[n-1].High > prevHigh*1.05 {
	// 	return false
	// }

	// // 5. DIF 顶背离增强
	// if data[n-1].Dif >= data[peak].Dif {
	// 	return false
	// }

	return true
}

func findLastMacdTopMountain(data []shares.Day) (peak, start, end int, ok bool) {
	n := len(data)
	// 当前MACD(Hist)必须0轴上方
	if data[n-currentDay].Hist <= 0 {
		return peak, start, end, false
	}

	// 从后往前找最近的山包
	for i := n - 2; i >= 2; i-- {
		if data[i].Hist <= 0 {
			return peak, start, end, false
		}
		// 通过(Hist)的递减找到山顶
		if data[i].Hist < data[i-1].Hist {
			peak = i
		} else {
			break
		}
	}

	if peak == 0 {
		return peak, start, end, false
	}

	// 从山顶向左找山包起点（递增开始）
	// start = peak - 1
	// start = peak
	for i := peak; i >= 2; i-- {
		if data[i].Hist > 0 && data[i-1].Hist < data[i].Hist {
			start = i
		} else {
			break
		}
	}
	if start == 0 {
		return peak, start, end, false
	}

	end = n - currentDay

	return peak, start, end, true
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