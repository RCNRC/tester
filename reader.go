package tester

import (
	"bufio"
	"strconv"
	"strings"
)

func ReadLine(in *bufio.Reader) string {
	line, _ := in.ReadString('\n')
	line = strings.ReplaceAll(line, "\r", "")
	line = strings.ReplaceAll(line, "\n", "")
	return line
}

func ReadInt(in *bufio.Reader) int {
	n, _ := strconv.Atoi(ReadLine(in))
	return n
}

func ReadLineNumbs(in *bufio.Reader) []string {
	return strings.Split(ReadLine(in), " ")
}

func ReadArrInt(in *bufio.Reader) []int {
	numbs := ReadLineNumbs(in)
	arr := make([]int, len(numbs))
	for i, n := range numbs {
		val, _ := strconv.Atoi(n)
		arr[i] = val
	}
	return arr
}

func ReadArrInt64(in *bufio.Reader) []int64 {
	numbs := ReadLineNumbs(in)
	arr := make([]int64, len(numbs))
	for i, n := range numbs {
		val, _ := strconv.ParseInt(n, 10, 64)
		arr[i] = val
	}
	return arr
}
