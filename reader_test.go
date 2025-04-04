package tester

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func restoring(stdin *os.File, readIO *os.File) func() {
	return func() {
		os.Stdin = stdin
		readIO.Close()
	}
}

func mockIO(t *testing.T) (*os.File, func()) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	restoreFunc := restoring(os.Stdin, r)
	os.Stdin = r

	return w, restoreFunc
}

func TestReadLine(t *testing.T) {
	var inpText string

	cases := []string{
		"",
		"asdsad",
		"dsad asdas d  sf sd fd g dgsfdgfdg dfg",
		"   as a   ",
		"   ",
	}

	inputStrings := append(cases, "")

	inpText = strings.Join(inputStrings, "\n")
	input := []byte(inpText)

	w, deferFunc := mockIO(t)
	defer deferFunc()

	_, err := w.Write(input)
	if err != nil {
		t.Error(err)
	}

	in := bufio.NewReader(os.Stdin)
	for _, curCase := range cases {
		n := ReadLine(in)
		assert.Equal(t, curCase, n)
	}
}

func TestReadLineNumbs(t *testing.T) {
	var inpText string

	cases := []string{
		"",
		"asdsad",
		"dsad asdas d  sf sd fd g dgsfdgfdg dfg",
		"   as a   ",
		"   ",
	}

	inputStrings := append(cases, "")

	inpText = strings.Join(inputStrings, "\n")
	input := []byte(inpText)

	w, deferFunc := mockIO(t)
	defer deferFunc()

	_, err := w.Write(input)
	if err != nil {
		t.Error(err)
	}

	in := bufio.NewReader(os.Stdin)
	for _, curCase := range cases {
		n := ReadLineNumbs(in)
		assert.Equal(t, strings.Split(curCase, " "), n)
	}
}

func TestReadInt(t *testing.T) {
	var inpText string

	cases := []int{-120, -5, -1, 0, 1, 5, 123, 1234567890}
	inputStrings := make([]string, len(cases)+1)
	for i, curCase := range cases {
		inputStrings[i] = strconv.Itoa(curCase)
	}
	inputStrings[len(inputStrings)-1] = ""

	inpText = strings.Join(inputStrings, "\n")
	input := []byte(inpText)

	w, deferFunc := mockIO(t)
	defer deferFunc()

	_, err := w.Write(input)
	if err != nil {
		t.Error(err)
	}

	in := bufio.NewReader(os.Stdin)
	for _, curCase := range cases {
		n := ReadInt(in)
		assert.Equal(t, curCase, n)
	}
}

func TestReadArrInt(t *testing.T) {
	var inpText string
	var err error

	cases := [][]int{
		{-120, -5, -1, 0, 1, 5, 123, 1234567890},
		{1},
	}
	inputStrings := make([]string, len(cases)+1)
	for i, curCase := range cases {
		inputString := make([]string, len(curCase))
		for j, curCurCase := range curCase {
			inputString[j] = strconv.Itoa(curCurCase)
		}
		inputStrings[i] = strings.Join(inputString, " ")
	}
	inputStrings[len(inputStrings)-1] = ""

	inpText = strings.Join(inputStrings, "\n")
	input := []byte(inpText)

	w, deferFunc := mockIO(t)
	defer deferFunc()

	_, err = w.Write(input)
	if err != nil {
		t.Error(err)
	}

	in := bufio.NewReader(os.Stdin)
	for _, curCase := range cases {
		n := ReadArrInt(in)
		assert.ElementsMatch(t, curCase, n)
	}
}

func TestReadArrInt64(t *testing.T) {
	var inpText string
	var err error

	cases := [][]int64{
		{-120, -5, -1, 0, 1, 5, 123, 1234567890, 1234567890123456789},
		{1},
	}
	inputStrings := make([]string, len(cases)+1)
	for i, curCase := range cases {
		inputString := make([]string, len(curCase))
		for j, curCurCase := range curCase {
			inputString[j] = strconv.FormatInt(curCurCase, 10)
		}
		inputStrings[i] = strings.Join(inputString, " ")
	}
	inputStrings[len(inputStrings)-1] = ""

	inpText = strings.Join(inputStrings, "\n")
	input := []byte(inpText)

	w, deferFunc := mockIO(t)
	defer deferFunc()

	_, err = w.Write(input)
	if err != nil {
		t.Error(err)
	}

	in := bufio.NewReader(os.Stdin)
	for _, curCase := range cases {
		n := ReadArrInt64(in)
		assert.ElementsMatch(t, curCase, n)
	}
}
