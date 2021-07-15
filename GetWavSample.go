package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/youpy/go-wav"
	"io"
	"math"
	"os"
	"strconv"
)

func main() {
	inFilePath := flag.String("i", "", "input wav file path")
	outFilePath := flag.String("o", "", "output pcm file path")
	bitDepth := flag.Uint("b", 10, "output pcm bit depth (Optional)")
	duration := flag.Uint("d", 0, "duration ms (Optional)")
	channel := flag.Uint("c", 0, "channel number(Optional)")
	positive := flag.Bool("p", false, "output is all positive or not (Optional true or false)")
	flag.Parse()
	//打印帮助
	if *inFilePath == "" || *outFilePath == "" {
		flag.Usage()
		return
	}
	//输入文件路径
	fileRead, _ := os.Open(*inFilePath)
	defer fileRead.Close()
	//输出文件路径
	fileWrite, _ := os.OpenFile(*outFilePath, os.O_CREATE|os.O_APPEND|os.O_TRUNC, 0600)
	defer fileWrite.Close()
	//输出十进制数据文件
	decWrite, _ := os.OpenFile(*outFilePath+".dec.txt", os.O_CREATE|os.O_APPEND|os.O_TRUNC, 0600)
	defer decWrite.Close()
	//创建wav结构
	reader := wav.NewReader(fileRead)
	//获取wav格式参数
	format, _ := reader.Format()
	count := 0
	for {
		samples, err := reader.ReadSamples()
		if err == io.EOF {
			break
		}
		for _, sample := range samples {
			//通过输入d时长参数计算采样数量判断
			if count >= int(*duration)*int(format.SampleRate)/1000 && *duration != 0 {
				break
			}
			downSample := sampleScaler(reader.IntValue(sample, *channel), uint(format.BitsPerSample), *bitDepth, *positive)
			s := decToBin(downSample, *bitDepth)
			//十进制数据写入文件
			decWrite.WriteString(strconv.Itoa(downSample) + ",\n")
			//二进制文件数据写入文件
			n, e := fileWrite.WriteString(s)
			if e == nil && n != len(s) {
				println(`错误代码：`, n)
				panic(err)
			}
			count++
		}

	}
}

//十进制转二进制
func decToBin(dec int, bit uint) string {
	//int转二进制并补全位数
	s := fmt.Sprintf("%0*b\n", bit, dec)
	//fmt.Printf("%s\n", insertNth(s))
	return insertNth(s)
}

//采样位深下变换
func sampleScaler(input int, bitPerSample uint, downBit uint, pos bool) int {
	bitScaler := bitPerSample - downBit
	if bitScaler < 0 {
		bitScaler = 0
	}
	//补正判断
	if pos {
		input += int(math.Pow(2, float64(bitPerSample-1)))
	}
	//采样位深下变换，采样除以2的bitScaler次方再加上2的bit-1次方
	return input / int(math.Pow(2, float64(bitScaler)))
}

//插入间隔符并判断转换-1
func insertNth(s string) string {
	var buffer bytes.Buffer
	str := "1"
	for i, r := range s {
		if string(r) == "-" {
			str = "-1"
			buffer.WriteString("0")
		}
		if string(r) == "1" {
			buffer.WriteString(str)
		}
		if string(r) == "0" {
			buffer.WriteString("0")
		}
		if i < len(s)-2 {
			buffer.WriteRune(',')
		}
		if i == len(s)-1 {
			buffer.WriteString("\n")
		}
	}
	return buffer.String()
}
