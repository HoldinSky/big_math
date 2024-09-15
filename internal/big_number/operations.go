package big_number

import (
	"strconv"
	"strings"

	u "github.com/HoldinSky/big_math/internal/utils"
)

func sumInBase(val1 int64, val2 int64, carry *int64) int64 {
	var sum int64 = val1 + val2 + *carry
	*carry = 0

	for sum >= base {
		sum -= base
		*carry++
	}

	return sum
}

func diffInBase(minuend int64, subtrahend int64, carry *int64) int64 {
	var diff int64

	if minuend-*carry < subtrahend {
		diff = base + minuend - subtrahend - *carry
		*carry = 1
	} else {
		diff = minuend - subtrahend - *carry
		*carry = 0
	}

	return diff
}

func absSum(bn1 *BigNumber, bn2 *BigNumber) BigNumber {
	result := new("")

	var v1, v2, sum, carry int64

	i, len1, len2 := 0, len(bn1.segments), len(bn2.segments)
	maxLen := u.Max(len1, len2)

	for ; i < maxLen || carry > 0; i++ {
		if i < len1 {
			v1 = bn1.segments[i]
		}
		if i < len2 {
			v2 = bn2.segments[i]
		}
		sum = sumInBase(v1, v2, &carry)
		result.appendSegment(sum)

		v1, v2 = 0, 0
	}

	return result
}

func absDiff(minuend *BigNumber, subtrahend *BigNumber) BigNumber {
	result := new("")

	var lenM, lenS, mi, si int = len(minuend.segments), len(subtrahend.segments), 0, 0
	var diff, carry, zeroSegmentsCount int64

	iteration := func(a int64, b int64) {
		diff = diffInBase(a, b, &carry)

		if diff == 0 {
			zeroSegmentsCount++
		} else {
			result.appendEmptySegments(&zeroSegmentsCount)
			result.appendSegment(diff)
		}
	}

	for si < lenS {
		iteration(minuend.segments[mi], subtrahend.segments[si])

		mi++
		si++
	}

	for mi < lenM {
		iteration(minuend.segments[mi], 0)

		mi++
	}

	return result
}

// returns: quotient, remainder
func (bn *BigNumber) DivWithRem(divisor *BigNumber) (BigNumber, BigNumber) {
	if divisor.Eq(&ZERO) {
		panic("Division by zero!")
	}

	var quotient, remainder, numerator, windowProd BigNumber
	if bn.absLt(divisor) {
		quotient = new("0")
		if bn.sign != divisor.sign {
			remainder = divisor.Add(bn)
		} else {
			remainder = Copy(bn)
		}
		return quotient, remainder
	}

	if divisor.absEq(&ONE) {
		quotient = Copy(bn)
		remainder = defaultBn()

		if divisor.sign == negative {
			quotient.invertSign()
		}
		remainder.sign = bn.sign

		return quotient, remainder
	}

	numeratorStr := strings.TrimLeft(bn.String(), "-")
	ptr := len(strings.TrimLeft(divisor.String(), "-"))

	numeratorWindow := numeratorStr[:ptr]
	numerator = new(numeratorWindow)

	quotientStr := strings.Builder{}
	var carriedOnce bool

	carryIfNeeded := func() {
		carriedOnce = true
		for numerator.absLt(divisor) && ptr < len(numeratorStr) {
			if !carriedOnce {
				quotientStr.WriteString("0")
			}
			carriedOnce = false

			numeratorWindow = numerator.String() + numeratorStr[ptr:ptr+1]
			numerator = new(numeratorWindow)

			ptr++
		}
	}

	findIntermediateQuotient := func() int {
		windowProd = divisor.Mul(&NINE)
		i := 9

		for numerator.absLt(&windowProd) {
			windowProd = absDiff(&windowProd, divisor)
			i--
		}

		return i
	}

	updateResults := func(interQuot int) {
		numerator = absDiff(&numerator, &windowProd)
		quotientStr.WriteString(strconv.Itoa(interQuot))
	}

	for ptr < len(numeratorStr) {
		carryIfNeeded()
		if ptr >= len(numeratorStr) && numerator.absLt(divisor) {
			break
		}

		interQuotient := findIntermediateQuotient()
		updateResults(interQuotient)
	}

	quotient = new(quotientStr.String())
	remainder = Copy(&numerator)

	remainder.sign = bn.sign
	if bn.sign != divisor.sign {
		quotient.sign = negative
	}

	return quotient, remainder
}

// ================================================================================
// ================================================================================

func (bn *BigNumber) Add(other *BigNumber) BigNumber {
	if bn.Eq(&ZERO) {
		return Copy(other)
	}
	if other.Eq(&ZERO) {
		return Copy(bn)
	}

	var dummy BigNumber

	if bn.sign != other.sign {
		if bn.sign == positive {
			dummy = Copy(other)
			dummy.invertSign()

			return bn.Subtract(&dummy)
		}

		dummy = Copy(bn)
		dummy.invertSign()

		return other.Subtract(&dummy)
	}

	result := absSum(bn, other)
	result.sign = bn.sign

	return result
}

func (bn *BigNumber) Subtract(subtrahend *BigNumber) BigNumber {
	if subtrahend.Eq(&ZERO) {
		return Copy(bn)
	}

	var result BigNumber
	if bn.Eq(&ZERO) {
		result = Copy(subtrahend)
		result.invertSign()
		return result
	}

	if bn.sign != subtrahend.sign {
		result = absSum(bn, subtrahend)
		result.sign = bn.sign
		return result
	}

	if bn.absGt(subtrahend) {
		result = absDiff(bn, subtrahend)
		result.sign = bn.sign
	} else if bn.absLt(subtrahend) {
		result = absDiff(subtrahend, bn)
		result.sign = subtrahend.oppositeSign()
	} else {
		result = new("0")
	}

	return result
}

func (bn *BigNumber) Mul(factor *BigNumber) BigNumber {
	var result BigNumber
	if bn.Eq(&ZERO) || factor.Eq(&ZERO) {
		return Copy(&ZERO)
	}
	if bn.absEq(&ONE) {
		result = Copy(factor)
		if bn.sign == negative {
			result.invertSign()
		}
		return result
	}
	if factor.absEq(&ONE) {
		result = Copy(bn)
		if factor.sign == negative {
			result.invertSign()
		}
		return result
	}

	appendZeros := func(target *strings.Builder, count int) {
		for i := 0; i < count; i++ {
			target.WriteRune('0')
		}
	}

	prod := strings.Builder{}
	for i, seg1 := range bn.segments {
		for k, seg2 := range factor.segments {
			prod.Reset()
			prod.WriteString(strconv.Itoa(int(seg1 * seg2)))
			appendZeros(&prod, (i+k)*base_power)

			p := new(prod.String())
			result = result.Add(&p)
		}
	}

	if bn.sign != factor.sign {
		result.sign = negative
	}

	return result
}

func (bn *BigNumber) Squared() BigNumber {
	return bn.Mul(bn)
}

func (bn *BigNumber) Div(divisor *BigNumber) BigNumber {
	quot, _ := bn.DivWithRem(divisor)
	return quot
}

func (bn *BigNumber) Mod(divisor *BigNumber) BigNumber {
	_, remainder := bn.DivWithRem(divisor)
	return remainder
}

func parseInt64(str string) int64 {
	num, _ := strconv.Atoi(str)

	return int64(num)
}
