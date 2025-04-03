package tester

import (
	"fmt"
	"math/rand"
	"reflect"
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
	Args            []interface{}
	PossibleResults []interface{}
}

// Check проверяет наличие ответа в возможных результатах
func (tc *TestCase) Check(answer interface{}) bool {
	for _, res := range tc.PossibleResults {
		if reflect.DeepEqual(answer, res) {
			return true
		}
	}
	return false
}

// GT интерфейс для генераторов тестовых данных
type GT interface {
	Generate()
	Val() interface{}
}

// GStatic обертка для статических значений
type GStatic struct {
	value interface{}
}

func (g *GStatic) Generate()        {}
func (g *GStatic) Val() interface{} { return g.value }

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

func (g *GInt) Val() interface{} { return g.val }

// GUInt генератор положительных целых чисел
type GUInt struct{ *GInt }

func NewGUInt(max GT) *GUInt {
	return &GUInt{NewGInt(&GStatic{1}, max)}
}

// GList генератор списков
type GList struct {
	elementType GT
	amount      GT
	listFunc    func([]interface{}) interface{}
	val         interface{}
}

func NewGList(elementType GT, amount GT, listFunc func([]interface{}) interface{}) *GList {
	return &GList{elementType, amount, listFunc, nil}
}

func (g *GList) Generate() {
	amount := g.amount.Val().(int)
	elements := make([]interface{}, amount)
	for i := 0; i < amount; i++ {
		g.elementType.Generate()
		elements[i] = g.elementType.Val()
	}
	rand.Shuffle(len(elements), func(i, j int) {
		elements[i], elements[j] = elements[j], elements[i]
	})
	g.val = g.listFunc(elements)
}

func (g *GList) Val() interface{} { return g.val }

// GTuple генератор кортежей
type GTuple struct {
	elements []GT
	val      []interface{}
}

func NewGTuple(elements []GT) *GTuple {
	return &GTuple{elements, nil}
}

func (g *GTuple) Generate() {
	g.val = make([]interface{}, len(g.elements))
	for i, e := range g.elements {
		e.Generate()
		g.val[i] = e.Val()
	}
}

func (g *GTuple) Val() interface{} { return g.val }

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

func (g *GChar) Val() interface{} { return string(g.val) }

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

func (g *GFrozStr) Val() interface{} { return g.val }

// GTester менеджер тестирования
type GTester struct {
	testFunc      func(TestCase) interface{}
	universalFunc func(TestCase) interface{}
	funcArgs      []GT
	allArgs       []GT
}

func NewGTester(
	testFunc func(TestCase) interface{},
	universalFunc func(TestCase) interface{},
	funcArgs []GT,
	allArgs []GT,
) *GTester {
	return &GTester{testFunc, universalFunc, funcArgs, allArgs}
}

func (gt *GTester) Generate1() TestCase {
	for _, arg := range gt.allArgs {
		arg.Generate()
	}
	args := make([]interface{}, len(gt.funcArgs))
	for i, arg := range gt.funcArgs {
		args[i] = arg.Val()
	}
	return TestCase{args, nil}
}

func (gt *GTester) Test(amount int, timeLimit float64, printRight int, failOn int, dontCountAnswers []interface{}) {
	failed := 0
	for i := 0; i < amount; i++ {
		rand.Seed(int64(i))
		tc := gt.Generate1()
		start := time.Now()
		res := gt.testFunc(tc)
		duration := time.Since(start).Seconds()

		expected := gt.universalFunc(tc)
		tc.PossibleResults = []interface{}{expected}

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

// Пример использования
func simpleTest() {
	n := NewGUInt(&GStatic{10})
	m := NewGInt(&GStatic{-10}, &GStatic{10})
	B := NewGTuple([]GT{n, m})
	A := NewGList(m, n, func(e []interface{}) interface{} { return e })

	f := func(tc TestCase) interface{} {
		// b := tc.Args[0].([]interface{})
		x := tc.Args[1].([]interface{})
		res := 0
		for _, v := range x {
			num := v.(int)
			if num != 2 { // Логическая ошибка из примера
				res += num
			}
		}
		return res
	}

	uf := func(tc TestCase) interface{} {
		x := tc.Args[1].([]interface{})
		sum := 0
		for _, v := range x {
			sum += v.(int)
		}
		return sum
	}

	tester := NewGTester(f, uf, []GT{B, A}, []GT{n, m, B, A})
	tester.Test(10, 0.5, 1, 3, nil)
}
