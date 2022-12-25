package main

import (
	util "aoc2022/aoc"
	inp "aoc2022/input"
	"fmt"
	"strconv"
)

type snafuDigit int

const (
	DMINUS snafuDigit = -2
	MINUS  snafuDigit = -1
	ZERO   snafuDigit = 0
	ONE    snafuDigit = 1
	TWO    snafuDigit = 2
)

const SNAFU_BASE = 5

type snafu []snafuDigit

func (s snafu) String() string {
	bytes := make([]byte, len(s))
	for i := range bytes {
		switch s[len(s)-1-i] {
		case DMINUS:
			bytes[i] = '='
		case MINUS:
			bytes[i] = '-'
		case ZERO:
			bytes[i] = '0'
		case ONE:
			bytes[i] = '1'
		case TWO:
			bytes[i] = '2'
		default:
			panic("Unknown SNAFU digit")
		}
	}
	return string(bytes)
}

func (s snafu) toDecimal() int {
	dec := 0
	for i, digit := range s {
		dec += int(digit) * util.Pow(SNAFU_BASE, i)
	}
	return dec
}

func snafuFromDecimal(dec int) snafu {
	digits := util.Map([]byte(strconv.FormatInt(int64(dec), SNAFU_BASE)), func(b byte) int { return int(b - '0') })
	sn := make(snafu, 0, len(digits)+1)
	carry := 0
	for i := len(digits) - 1; i >= 0; i-- {
		d := digits[i]
		if d < 0 || d >= SNAFU_BASE {
			panic("Bad number base conversion")
		}

		d += carry

		if util.Between(d, 0, 2) {
			carry = 0
			sn = append(sn, snafuDigit(d))
		} else {
			carry = 1
			sn = append(sn, snafuDigit(d-SNAFU_BASE))
		}
	}
	if carry > 0 {
		sn = append(sn, snafuDigit(carry))
	}
	return sn
}

func parseSnafu(s string) snafu {
	sn := make(snafu, len(s))
	for i, b := range []byte(s) {
		switch b {
		case '=':
			sn[len(s)-1-i] = DMINUS
		case '-':
			sn[len(s)-1-i] = MINUS
		case '0':
			sn[len(s)-1-i] = ZERO
		case '1':
			sn[len(s)-1-i] = ONE
		case '2':
			sn[len(s)-1-i] = TWO
		}
	}
	return sn
}

func main() {
	lines := inp.ReadLines("input")
	snafus := util.Map(lines, parseSnafu)

	decSum := util.Sum(util.Map(snafus, func(s snafu) int { return s.toDecimal() }))
	snafuSum := snafuFromDecimal(decSum)
	fmt.Printf("Part 1: %s\n", snafuSum)
}
