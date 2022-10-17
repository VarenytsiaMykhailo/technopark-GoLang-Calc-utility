package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
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

func parse(expression string) (result float64, err error) {
	expression = strings.ReplaceAll(expression, " ", "")

	numbersStack := Stack[float64]{}
	operatorsStack := Stack[rune]{}

	expectNumOrLeftComa := true
	for utf8.RuneCountInString(expression) > 0 {
		if expectNumOrLeftComa {
			if isLeftBracket(rune(expression[0])) {
				operatorsStack.push(rune(expression[0]))
				expression = expression[1:]
				continue
			}
			number, numberLength, err := parseNumber(expression)
			if err != nil {
				return result, errors.New("error while parsing number: " + err.Error())
			}
			expression = expression[numberLength:]
			numbersStack.push(number)
			expectNumOrLeftComa = false
			continue
		}
		operator := rune(expression[0])
		if !isOperator(operator) {
			return result, errors.New("incorrect operator " + string(operator))
		}
		expression = expression[1:]
		if operatorsStack.isEmpty() {
			operatorsStack.push(operator)
			expectNumOrLeftComa = true
			continue
		}
		if operatorsPriorities[operator] > operatorsPriorities[operatorsStack.peek()] {
			operatorsStack.push(operator)
			expectNumOrLeftComa = true
			continue
		}

		// Need to calc subexpression
		err = calcSubExpressions(&numbersStack, &operatorsStack, operator)
		if err != nil {
			return result, errors.New("error while calc subexpression: " + err.Error())
		}

		if isRightBracket(operator) && isLeftBracket(operatorsStack.peek()) {
			operatorsStack.top()
			continue
		}
		operatorsStack.push(operator)
		expectNumOrLeftComa = true
	}

	for !operatorsStack.isEmpty() {
		lastOperatorInStack := operatorsStack.top()
		rightNumber := numbersStack.top()
		leftNumber := numbersStack.top()
		operationResult, err := executeOperator(leftNumber, rightNumber, lastOperatorInStack)
		if err != nil {
			return result, errors.New("error while execute operator: " + err.Error())
		}

		numbersStack.push(operationResult)
	}

	return numbersStack.top(), err
}

func calcSubExpressions(numbersStack *Stack[float64], operatorsStack *Stack[rune], operator rune) error {
	for !operatorsStack.isEmpty() && (operatorsPriorities[operator] <= operatorsPriorities[operatorsStack.peek()]) {
		lastOperatorInStack := operatorsStack.top()
		rightNumber := numbersStack.top()
		leftNumber := numbersStack.top()
		operationResult, err := executeOperator(leftNumber, rightNumber, lastOperatorInStack)
		if err != nil {
			return errors.New("error while execute operator: " + err.Error())
		}

		numbersStack.push(operationResult)
	}

	return nil
}

func executeOperator(leftNumber float64, rightNumber float64, operator rune) (result float64, err error) {
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

func parseNumber(expression string) (number float64, length int, err error) {
	numberStr := ""
	pointsCounter := 0
	for idx, v := range expression {
		if unicode.IsDigit(v) {
			numberStr += string(v)
			continue
		}
		if isDot(v) && numberStr != "" && pointsCounter == 0 {
			pointsCounter++
			numberStr += string(v)
			continue
		}
		if !isOperator(v) {
			return number, length, errors.New("invalid symbol for number: " + string(v) + " at pos: " + strconv.Itoa(idx))
		}
		break
	}
	number, err = strconv.ParseFloat(numberStr, 64)

	return number, utf8.RuneCountInString(numberStr), err
}

func isDot(v rune) bool {
	return v == '.'
}

func isLeftBracket(v rune) bool {
	return v == '('
}

func isRightBracket(v rune) bool {
	return v == ')'
}

func isOperator(char rune) bool {
	return char == ')' || char == '+' || char == '-' || char == '*' || char == '/'
}
