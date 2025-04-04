package tester

import (
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"time"
)

// BColors содержит ANSI коды для цветного вывода
var BColors = struct {
	Header    string
	OKBlue    string
	OKCyan    string
	OKGreen   string
	Warning   string
	Fail      string
	End       string
	Bold      string
	Underline string
}{
	Header:    "\033[95m",
	OKBlue:    "\033[94m",
	OKCyan:    "\033[96m",
	OKGreen:   "\033[92m",
	Warning:   "\033[93m",
	Fail:      "\033[91m",
	End:       "\033[0m",
	Bold:      "\033[1m",
	Underline: "\033[4m",
}

// TestCase представляет тестовый случай
type TestCase struct {
	Args            []any
	PossibleResults []any
}

// Check проверяет наличие ответа в возможных результатах
func (tc *TestCase) Check(answer any) bool {
	for _, res := range tc.PossibleResults {
		if reflect.DeepEqual(answer, res) {
			return true
		}
	}
	return false
}

// CheckFunc проверяет функцию на готовых примерах
func CheckFunc(tcs []TestCase, tFunc func(TestCase) any) {
	for i, tc := range tcs {
		res := tFunc(tc)
		if !tc.Check(res) {
			fmt.Printf("%s======== TEST #%d FAILED ========%s\nArgs: %v\nGot: %v\nExpected one of: %v\n",
				BColors.Fail, i+1, BColors.End, tc.Args, res, tc.PossibleResults)
		}
	}
}

// GT интерфейс для генераторов тестовых данных
type GT interface {
	Generate()
	Val() any
}

// GStatic обертка для статических значений
type GStatic struct {
	value any
}

func (g *GStatic) Generate() {}
func (g *GStatic) Val() any  { return g.value }

func NewGStatic(value any) *GStatic {
	return &GStatic{value: value}
}

// GInt генератор целых чисел
type GInt struct {
	min GT
	max GT
	val int
}

func NewGInt(min, max GT) *GInt {
	return &GInt{min: min, max: max}
}

func (g *GInt) Generate() {
	minVal := g.min.Val().(int)
	maxVal := g.max.Val().(int)
	g.val = rand.Intn(maxVal-minVal+1) + minVal
}

func (g *GInt) Val() any { return g.val }

// GUInt генератор положительных целых чисел
type GUInt struct{ *GInt }

func NewGUInt(max GT) *GUInt {
	return &GUInt{NewGInt(&GStatic{1}, max)}
}

// GList генератор списков
type GList struct {
	elementType GT
	amount      GT
	val         []any
}

func NewGList(elementType GT, amount GT) *GList {
	return &GList{elementType, amount, []any{}}
}

func (g *GList) Generate() {
	amount := g.amount.Val().(int)
	var elements []any = make([]any, amount)
	for i := 0; i < amount; i++ {
		g.elementType.Generate()
		elements[i] = g.elementType.Val()
	}
	rand.Shuffle(len(elements), func(i, j int) {
		elements[i], elements[j] = elements[j], elements[i]
	})
	g.val = elements
}

func (g *GList) Val() any { return g.val }

// GTuple генератор кортежей
type GTuple struct {
	elements []GT
	val      []any
}

func NewGTuple(elements []GT) *GTuple {
	return &GTuple{elements, nil}
}

func (g *GTuple) Generate() {
	g.val = make([]any, len(g.elements))
	for i, e := range g.elements {
		e.Generate()
		g.val[i] = e.Val()
	}
}

func (g *GTuple) Val() any { return g.val }

// GChar генератор символов
type GChar struct {
	chars string
	val   rune
}

func NewGChar(chars string) *GChar {
	return &GChar{chars, 0}
}

func (g *GChar) Generate() {
	g.val = rune(g.chars[rand.Intn(len(g.chars))])
}

func (g *GChar) Val() any { return string(g.val) }

// GFrozStr генератор перемешанных строк
type GFrozStr struct {
	base string
	val  string
}

func NewGFrozStr(base string) *GFrozStr {
	return &GFrozStr{base, ""}
}

func (g *GFrozStr) Generate() {
	runes := []rune(g.base)
	rand.Shuffle(len(runes), func(i, j int) {
		runes[i], runes[j] = runes[j], runes[i]
	})
	g.val = string(runes)
}

func (g *GFrozStr) Val() any { return g.val }

// GStr генератор строки
type GStr struct {
	gChar  GChar
	length GT
	delim  string
	val    string
}

func NewGStr(gChar GChar, length GT, delim string) *GStr {
	return &GStr{gChar, length, delim, ""}
}

func (g *GStr) Generate() {
	lenght := g.length.Val().(int)
	newChars := make([]string, lenght)
	for i := 0; i < lenght; i++ {
		g.gChar.Generate()
		newChars[i] = g.gChar.Val().(string)
	}
	g.val = strings.Join(newChars, g.delim)
}

func (g *GStr) Val() any { return g.val }

// GTester менеджер тестирования
type GTester struct {
	testFunc      func(TestCase) any
	universalFunc func(TestCase) any
	funcArgs      []GT
	allArgs       []GT
}

func NewGTester(
	testFunc func(TestCase) any,
	universalFunc func(TestCase) any,
	funcArgs []GT,
	allArgs []GT,
) *GTester {
	return &GTester{testFunc, universalFunc, funcArgs, allArgs}
}

func (gt *GTester) Generate1() TestCase {
	for _, arg := range gt.allArgs {
		arg.Generate()
	}
	args := make([]any, len(gt.funcArgs))
	for i, arg := range gt.funcArgs {
		args[i] = arg.Val()
	}
	return TestCase{args, nil}
}

func (gt *GTester) Test(amount int, timeLimit float64, printRight int, failOn int) {
	failed := 0
	for i := 0; i < amount; i++ {
		rand.New(rand.NewSource(int64(i)))
		tc := gt.Generate1()
		start := time.Now()
		res := gt.testFunc(tc)
		duration := time.Since(start).Seconds()

		expected := gt.universalFunc(tc)
		tc.PossibleResults = []any{expected}

		if !tc.Check(res) {
			fmt.Printf("%s======== TEST #%d FAILED ========%s\nArgs: %v\nGot: %v\nExpected: %v\nTime: %.3fs\n",
				BColors.Fail, i+1, BColors.End, tc.Args, res, expected, duration)
			failed++
		} else if duration > timeLimit {
			fmt.Printf("%s====== TEST #%d TIMEOUT ======%s\nArgs: %v\nTime: %.3fs (limit %.3fs)\n",
				BColors.Warning, i+1, BColors.End, tc.Args, duration, timeLimit)
		} else if printRight > 0 {
			fmt.Printf("%s======== TEST #%d PASSED ========%s\nTime: %.3fs\n",
				BColors.OKGreen, i+1, BColors.End, duration)
			if printRight > 1 {
				fmt.Printf("Args: %v\nResult: %v\n", tc.Args, res)
			}
		}

		if failOn > 0 && failed >= failOn {
			fmt.Printf("%sToo many failures (%d), aborting%s\n", BColors.Fail, failed, BColors.End)
			return
		}
	}
}

func (gt *GTester) TestProfile(amount int, timeLimit float64, printRight int) {
	for i := 0; i < amount; i++ {
		rand.New(rand.NewSource(int64(i)))
		tc := gt.Generate1()
		start := time.Now()
		gt.testFunc(tc)
		duration := time.Since(start).Seconds()

		if duration > timeLimit {
			fmt.Printf("%s====== TEST #%d TIMEOUT ======%s\nTime: %.3fs (limit %.3fs)\n",
				BColors.Warning, i+1, BColors.End, duration, timeLimit)
		} else if printRight > 0 {
			fmt.Printf("%s======== TEST #%d PASSED ========%s\nTime: %.3fs\n",
				BColors.OKGreen, i+1, BColors.End, duration)
		}
	}
}
