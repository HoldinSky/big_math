package main

import (
	"fmt"
	"os"
	"strconv"

	bn "github.com/HoldinSky/big_math/internal/big_number"
)

const usageStr = "USAGE: big_math.exe OPERATION [OPERAND1[ OPERAND2]|--random OPERAND1_LENGTH OPERAND2_LENGTH]\n\n" +
	"OPERATION:\n-a   addition\n-s   subtraction\n-m   multiplication\n-sq  square\n-d   division\n-mod modulus\n\n" +
	"Specify operands in sequence or use \"--random\" flag followed by the lengths of both operands to be created randomly"

const (
	addition       string = "-a"
	subtraction    string = "-s"
	multiplication string = "-m"
	square         string = "-sq"
	division       string = "-d"
	modulus        string = "-mod"
)

var operationsMap map[string]string = map[string]string{
	addition:       "+",
	subtraction:    "-",
	multiplication: "*",
	square:         "^2",
	division:       "/",
	modulus:        "%",
}

func printInfoAndUsageAndExit(info string) {
	fmt.Println(info)
	printUsageAndExit()
}

func printUsageAndExit() {
	fmt.Println(usageStr)
	os.Exit(-1)
}

func main() {
	args := os.Args

	if len(args) < 2 {
		printUsageAndExit()
	}

	if args[1] == "--help" {
		fmt.Println(usageStr)
		os.Exit(0)
	}

	opStr, exists := operationsMap[args[1]]

	if !exists {
		printUsageAndExit()
	}

	var operand1, operand2, result bn.BigNumber
	var err error

	if args[2] == "--random" {
		if len(args) < 5 {
			printUsageAndExit()
		}

		op1Len, err1 := strconv.Atoi(args[3])
		op2Len, err2 := strconv.Atoi(args[4])

		if err1 != nil || err2 != nil {
			printInfoAndUsageAndExit("Failed to parse operands' lengths")
		}

		operand1 = bn.RandomBigNumber(op1Len)
		operand2 = bn.RandomBigNumber(op2Len)
	} else {
		operand1, err = bn.NewBigNumber(args[2])
		if err != nil {
			printInfoAndUsageAndExit("Failed to parse OPERAND1")
		}

		if args[1] != square {
			operand2, err = bn.NewBigNumber(args[3])
			if err != nil {
				printInfoAndUsageAndExit("Failed to parse OPERAND2")
			}
		}
	}

	switch args[1] {
	case addition:
		result = operand1.Add(&operand2)
	case subtraction:
		result = operand1.Subtract(&operand2)
	case multiplication:
		result = operand1.Mul(&operand2)
	case square:
		result = operand1.Squared()
	case division:
		result = operand1.Div(&operand2)
	case modulus:
		result = operand1.Mod(&operand2)
	}

	if args[1] == square {
		fmt.Printf("%s %s\n =\n%s", operand1.String(), opStr, result.String())
	} else {
		fmt.Printf("%s\n %s\n%s\n =\n%s", operand1.String(), opStr, operand2.String(), result.String())
	}
}
