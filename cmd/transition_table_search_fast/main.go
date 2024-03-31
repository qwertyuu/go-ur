package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func ReadStdin(in chan int) {
	defer close(in)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		num, err := strconv.Atoi(line)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading number:", err)
			return
		}
		in <- num
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading standard input:", err)
		return
	}
}

func main() {
	in := make(chan int)

	go ReadStdin(in)

	v, ok := <-in
	if !ok {
		return
	}
	searchElementInt := uint32(v)
	searchElement := make([]byte, 4)
	binary.BigEndian.PutUint32(searchElement, searchElementInt)
	file, err := os.Open("/home/raph/Documents/GitHub/ur-lut-visualizer/lut_transition_merge_sort/round_8/output_0.data")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error opening file:", err)
		return
	}
	defer file.Close()
	c := 0
	var bytestr []byte
	var sb strings.Builder
	hasErr := false
	// TODO: Write to file but in binary
	for {
		bytestr = make([]byte, 100000000)
		byteLen, err := file.Read(bytestr)
		c += byteLen
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Fprintln(os.Stderr, "Error reading file:", err)
			break
		}
		for i := 0; i < byteLen; i += 10 {
			cmp := bytes.Compare(bytestr[i:i+4], searchElement)
			if cmp == 0 {
				sb.WriteString(strconv.FormatUint(uint64(searchElementInt), 10) + "," + strconv.FormatUint(uint64(binary.BigEndian.Uint32(bytestr[i+4:i+8])), 10) + "," + strconv.Itoa(int(bytestr[i+8])) + "," + strconv.Itoa(int(bytestr[i+9])) + "\n")

				if sb.Len() > 3000000000 {
					os.Stdout.Write([]byte(sb.String()))
					sb.Reset()
				}
			} else if cmp > 0 {
				v, ok := <-in
				if !ok {
					hasErr = true
					break
				}
				searchElementInt = uint32(v)
				binary.BigEndian.PutUint32(searchElement, searchElementInt)
				i -= 10
			}
		}
		if hasErr {
			break
		}
	}
	if sb.Len() > 0 {
		os.Stdout.Write([]byte(sb.String()))
	}
}
