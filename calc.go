package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

var (
	operatorsPriorities = map[rune]int{
		'(': 0,
		')': 1,
		'+': 2,
		'-': 2,
		'*': 3,
		'/': 3,
	}
)

func main() {
	var str string
	if len(os.Args) == 2 {
		str = os.Args[1]
	} else {
		log.Fatal(errors.New("incorrect program call"))
	}
	result, err := parse(str)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
}

type StackConstraint interface {
	rune | float64
}

type Stack[T StackConstraint] struct {
	data []T
}

func (o *Stack[T]) push(elem T) {
	o.data = append(o.data, elem)
}

func (o *Stack[T]) peek() T {
	return o.data[len(o.data)-1]
}

func (o *Stack[T]) top() T {
	n := len(o.data) - 1
	elem := o.data[n]
	o.data = o.data[:n]
	return elem
}

func (o *Stack[T]) isEmpty() bool {
	return len(o.data) == 0
}

func parse(expression string) (result float64, err error) {
	expression = strings.ReplaceAll(expression, " ", "")

	numbersStack := Stack[float64]{}
	operatorsStack := Stack[rune]{}

	expectNumOrLeftComa := true
	for utf8.RuneCountInString(expression) > 0 {
		if expectNumOrLeftComa {
			if expression[0] == '(' {
				operatorsStack.push(rune(expression[0]))
				expression = expression[1:]
				continue
			}
			number, numberLength, err := ParseNumber(expression)
			if err != nil {
				return result, err
			}
			expression = expression[numberLength:]
			numbersStack.push(number)
			expectNumOrLeftComa = false
		} else {
			operator := rune(expression[0])
			if !isOperator(operator) {
				return result, errors.New("incorrect operator " + string(operator))
			}
			expression = expression[1:]
			if err != nil {
				return result, err
			}
			if operatorsStack.isEmpty() {
				operatorsStack.push(operator)
				//expression = expression[1:]
				expectNumOrLeftComa = true
				continue
			} else if operatorsPriorities[operator] > operatorsPriorities[operatorsStack.peek()] {
				operatorsStack.push(operator)
				expectNumOrLeftComa = true
				continue
			} else { // Need to calc subexpression
				for !operatorsStack.isEmpty() && (operatorsPriorities[operator] <= operatorsPriorities[operatorsStack.peek()]) {
					lastOperatorInStack := operatorsStack.top()
					rightNumber := numbersStack.top()
					leftNumber := numbersStack.top()
					operationResult, err := ExecuteOperator(leftNumber, rightNumber, lastOperatorInStack)
					if err != nil {
						return result, err
					}

					numbersStack.push(operationResult)
				}
				if operator == ')' && operatorsStack.peek() == '(' {
					operatorsStack.top()
					continue
				} else {
					operatorsStack.push(operator)
					expectNumOrLeftComa = true
					continue
				}
			}
		}
	}

	for !operatorsStack.isEmpty() {
		lastOperatorInStack := operatorsStack.top()
		rightNumber := numbersStack.top()
		leftNumber := numbersStack.top()
		operationResult, err := ExecuteOperator(leftNumber, rightNumber, lastOperatorInStack)
		if err != nil {
			return result, err
		}

		numbersStack.push(operationResult)
	}

	return numbersStack.top(), err
}

func ExecuteOperator(leftNumber float64, rightNumber float64, operator rune) (result float64, err error) {
	if !isOperator(operator) {
		return result, errors.New("unsupported operator: " + string(operator))
	}
	switch operator {
	case '+':
		result = leftNumber + rightNumber
	case '-':
		result = leftNumber - rightNumber
	case '*':
		result = leftNumber * rightNumber
	case '/':
		result = leftNumber / rightNumber
	}

	return
}

func ParseNumber(expression string) (number float64, length int, err error) {
	numberStr := ""
	pointsCounter := 0
	for idx, v := range expression {
		if isDigit(v) {
			numberStr += string(v)
		} else if string(v) == "." && numberStr != "" && pointsCounter == 0 {
			pointsCounter++
			numberStr += string(v)
		} else {
			if !isOperator(v) {
				return number, length, errors.New("invalid symbol for number: " + string(v) + " at pos: " + strconv.Itoa(idx))
			}
			break
		}
	}
	number, err = strconv.ParseFloat(numberStr, 64)

	return number, utf8.RuneCountInString(numberStr), err
}

func isOperator(char rune) bool {
	return char == ')' || char == '+' || char == '-' || char == '*' || char == '/'
}

func isDigit(char rune) bool {
	return '0' <= char && char <= '9'
}
