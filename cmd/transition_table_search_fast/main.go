package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func main() {
	// populate searchList with ints from the "/home/raph/Documents/GitHub/ur-lut-visualizer/b41.txt" file
	searchList := make([]int, 0)
	file, err := os.Open("/home/raph/Documents/GitHub/ur-lut-visualizer/b41.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("reading")

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		num, err := strconv.Atoi(line)
		if err != nil {
			fmt.Println("Error converting line to integer:", err)
			return
		}
		searchList = append(searchList, num)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	fmt.Println("Read", len(searchList), "integers from file.")

	sort.Slice(searchList, func(i, j int) bool {
		return searchList[i] < searchList[j]
	})
	fmt.Println("searchList read and sorted")
	searchIndex := 0
	searchElementInt := uint32(searchList[searchIndex])
	searchElement := make([]byte, 4)
	binary.BigEndian.PutUint32(searchElement, searchElementInt)
	file.Close()
	file, err = os.Open("/home/raph/Documents/GitHub/ur-lut-visualizer/lut_transition_merge_sort/round_8/output_0.data")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	c := 0
	var bytestr []byte
	var sb strings.Builder
	// TODO: Write to file but in binary
	outFile, err := os.Create("foundList.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	startTime := time.Now()
	for {
		bytestr = make([]byte, 100000000)
		byteLen, err := file.Read(bytestr)
		c += byteLen
		if c%100000000 == 0 {
			fmt.Println(float32(searchIndex) / float32(len(searchList)))
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			return
		}
		for i := 0; i < byteLen; i += 10 {
			cmp := bytes.Compare(bytestr[i:i+4], searchElement)
			if cmp == 0 {
				sb.WriteString(strconv.FormatUint(uint64(searchElementInt), 10) + "," + strconv.FormatUint(uint64(binary.BigEndian.Uint32(bytestr[i+4:i+8])), 10) + "," + strconv.Itoa(int(bytestr[i+8])) + "," + strconv.Itoa(int(bytestr[i+9])) + "\n")

				if sb.Len() > 3000000000 {
					outFile.WriteString(sb.String())
					sb.Reset()
				}
			} else if cmp > 0 {
				searchIndex++
				if searchIndex >= len(searchList) {
					break
				}
				searchElementInt = uint32(searchList[searchIndex])
				binary.BigEndian.PutUint32(searchElement, searchElementInt)
				i -= 10
			}
		}
	}
	if sb.Len() > 0 {
		outFile.WriteString(sb.String())
	}
	outFile.Close()
	fmt.Println("Time elapsed:", time.Since(startTime))
}
