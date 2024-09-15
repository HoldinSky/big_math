package big_number

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
)

const (
	positive sign_t = 1
	negative sign_t = -1
)

var (
	base_power int   = 9
	base       int64 = int64(math.Pow10(int(base_power)))

	formatStr = "%0" + strconv.Itoa(base_power) + "d"

	ZERO BigNumber = new("0")
	ONE  BigNumber = new("1")
	NINE BigNumber = new("9")
)

type sign_t int8

type BigNumber struct {
	segments []int64
	sign     sign_t
}

// ================================================================================
// ================================================================================

func defaultBn() BigNumber {
	num := BigNumber{sign: positive}

	return num
}

func new(repr string) BigNumber {
	num, _ := NewBigNumber(repr)

	return num
}

func fromInt(val int) BigNumber {
	num, _ := NewBigNumber(strconv.Itoa(val))

	return num
}

func fromArray(values []int64, sign sign_t) BigNumber {
	num := defaultBn()

	leadingZeros := false
	var zerosCount int64 = 0
	for _, val := range values {
		leadingZeros = val == 0
		if leadingZeros {
			zerosCount++
			continue
		}

		num.appendEmptySegments(&zerosCount)
		num.appendSegment(val)

		zerosCount = 0
	}
	num.sign = sign

	return num
}

func RandomBigNumber(length int) BigNumber {
	var digit int
	numberStr := strings.Builder{}

	if rand.Intn(2) > 0 {
		numberStr.WriteString("-")
	}

	for ; length > 0; length-- {
		digit = 1 + rand.Intn(9)
		numberStr.WriteString(strconv.Itoa(digit))
	}

	return new(numberStr.String())
}

func NewBigNumber(repr string) (BigNumber, error) {
	result := defaultBn()
	if len(repr) > 0 && repr[0] == '-' {
		result.sign = negative
		repr = strings.TrimLeft(repr, "-")
	}

	var sliceToConvert string

	for i := len(repr); i > 0; i -= int(base_power) {
		if i < base_power {
			sliceToConvert = repr[0:i]

			val, err := strconv.Atoi(sliceToConvert)
			if err != nil {
				return result, err
			}

			result.appendSegment(int64(val))
		} else {
			sliceToConvert = repr[i-base_power : i]

			val, err := strconv.Atoi(sliceToConvert)
			if err != nil {
				return result, err
			}

			result.appendSegment(int64(val))
		}
	}

	return result, nil
}

func Copy(bn *BigNumber) BigNumber {
	res := defaultBn()

	len := len(bn.segments)
	for i := 0; i < len; i++ {
		res.appendSegment(bn.segments[i])
	}
	res.sign = bn.sign

	return res
}

func (bn BigNumber) String() string {
	i := len(bn.segments) - 1
	if i <= -1 {
		return "0"
	}

	builder := strings.Builder{}

	builder.WriteString(fmt.Sprintf("%d", int64(bn.sign)*bn.segments[i]))

	i--
	for ; i >= 0; i-- {
		builder.WriteString(fmt.Sprintf(formatStr, bn.segments[i]))
	}

	return builder.String()
}

func (bn *BigNumber) invertSign() {
	bn.sign = bn.oppositeSign()
}

func (bn *BigNumber) oppositeSign() sign_t {
	return bn.sign * -1
}

func (bn *BigNumber) Gt(other *BigNumber) bool {
	if bn.sign > other.sign {
		return true
	} else if bn.sign < other.sign {
		return false
	} else if bn.sign == positive {
		return bn.absGt(other)
	} else {
		return bn.absLt(other)
	}
}

func (bn *BigNumber) Lt(other *BigNumber) bool {
	if bn.sign < other.sign {
		return true
	} else if bn.sign > other.sign {
		return false
	} else if bn.sign == positive {
		return bn.absLt(other)
	} else {
		return bn.absGt(other)
	}
}

func (bn *BigNumber) Ge(other *BigNumber) bool {
	return !bn.Lt(other)
}

func (bn *BigNumber) Le(other *BigNumber) bool {
	return !bn.Gt(other)
}

func (bn *BigNumber) Eq(other *BigNumber) bool {
	return bn.sign == other.sign && bn.absEq(other)
}

func (bn *BigNumber) absGt(other *BigNumber) bool {
	thisLen, otherLen := len(bn.segments), len(other.segments)
	if thisLen > otherLen {
		return true
	}
	if thisLen < otherLen {
		return false
	}

	for i := thisLen - 1; i >= 0; i-- {
		if bn.segments[i] > other.segments[i] {
			return true
		}
		if bn.segments[i] < other.segments[i] {
			return false
		}
	}

	return false
}

func (bn *BigNumber) absLt(other *BigNumber) bool {
	thisLen, otherLen := len(bn.segments), len(other.segments)
	if thisLen < otherLen {
		return true
	}
	if thisLen > otherLen {
		return false
	}

	for i := thisLen - 1; i >= 0; i-- {
		if bn.segments[i] < other.segments[i] {
			return true
		}
		if bn.segments[i] > other.segments[i] {
			return false
		}
	}

	return false
}

func (bn *BigNumber) absGe(other *BigNumber) bool {
	return !bn.absLt(other)
}

func (bn *BigNumber) absLe(other *BigNumber) bool {
	return !bn.absGt(other)
}

func (bn *BigNumber) absEq(other *BigNumber) bool {
	thisLen, otherLen := len(bn.segments), len(other.segments)
	if thisLen != otherLen {
		return false
	}

	for i := thisLen - 1; i >= 0; i-- {
		if bn.segments[i] != other.segments[i] {
			return false
		}
	}

	return true
}

// ================================================================================
// ================================================================================

func (bn *BigNumber) appendSegment(val int64) {
	bn.segments = append(bn.segments, val)
}

func (bn *BigNumber) appendEmptySegments(count *int64) {
	for ; *count > 0; *count-- {
		bn.appendSegment(0)
	}
}
