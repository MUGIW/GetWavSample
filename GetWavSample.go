package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/youpy/go-wav"
	"io"
	"math"
	"os"
)

func main() {
	inFilePath := flag.String("i", "", "input wav file path")
	outFilePath := flag.String("o", "", "output pcm file path")
	bit := flag.Int("b", 10, "output pcm bit depth (Optional)")
	duration := flag.Int("d", 0, "duration ms (Optional)")
	flag.Parse()
	if *inFilePath == "" || *outFilePath == "" {
		flag.Usage()
		return
	}

	fileRead, _ := os.Open(*inFilePath)
	reader := wav.NewReader(fileRead)
	format, _ := reader.Format()
	defer fileRead.Close()
	fileWrite, _ := os.OpenFile(*outFilePath, os.O_CREATE|os.O_APPEND|os.O_TRUNC, 0600)
	defer fileWrite.Close()
	count := 0
	for {
		samples, err := reader.ReadSamples()
		if err == io.EOF {
			break
		}
		for _, sample := range samples {
			//fmt.Printf("%d %d %d \n", count, format.SampleRate, *duration)
			if count >= *duration*int(format.SampleRate)/1000 && *duration != 0 {
				break
			}
			downSample := sampleScaler(reader.IntValue(sample, 0), int(format.BitsPerSample), *bit)
			s := decToBin(downSample, *bit)
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
func decToBin(dec int, bit int) string {
	//int转二进制并补全位数
	s := fmt.Sprintf("%0*b\n", bit, dec)
	//fmt.Printf("%s\n", insertNth(s))
	return insertNth(s)
}

//采样位深下变换
func sampleScaler(input int, bitPerSample int, downBit int) int {
	bitScaler := bitPerSample - downBit
	if bitScaler < 0 {
		bitScaler = 0
	}
	//采样位深下变换并补为正，采样除以2的bitScaler次方再加上2的bit-1次方
	return input/int(math.Pow(2, float64(bitScaler))) + int(math.Pow(2, float64(downBit-1)))
}

//插入间隔符
func insertNth(s string) string {
	var buffer bytes.Buffer
	for i, r := range s {
		buffer.WriteRune(r)
		if i < len(s)-2 {
			buffer.WriteRune(',')
		}
	}
	return buffer.String()
}
