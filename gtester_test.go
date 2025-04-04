package tester

import "testing"

// Пример использования
func TestBasic(t *testing.T) {
	n := NewGUInt(NewGStatic(10))
	m := NewGInt(NewGStatic(-10), NewGStatic(10))
	B := NewGTuple([]GT{n, m})
	A := NewGList(m, n)

	f := func(tc TestCase) any {
		// b := tc.Args[0].([]any)
		x := tc.Args[1].([]any)
		res := 0
		for _, v := range x {
			num := v.(int)
			if num != 2 { // Логическая ошибка из примера
				res += num
			}
		}
		return res
	}

	uf := func(tc TestCase) any {
		x := tc.Args[1].([]any)
		sum := 0
		for _, v := range x {
			sum += v.(int)
		}
		return sum
	}

	tester := NewGTester(f, uf, []GT{B, A}, []GT{n, m, B, A})
	tester.Test(10, 0.5, 1, 3)
}
