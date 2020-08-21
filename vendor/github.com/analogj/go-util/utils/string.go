package utils

import (
	"strconv"
	"strings"
)

func StringToInt(input string) (int, error) {
	i, err := strconv.Atoi(input)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func SnakeCaseToCamelCase(inputUnderScoreStr string) (camelCase string) {
	//snake_case to camelCase

	isToUpper := false

	for k, v := range inputUnderScoreStr {
		if k == 0 {
			camelCase = strings.ToUpper(string(inputUnderScoreStr[0]))
		} else {
			if isToUpper {
				camelCase += strings.ToUpper(string(v))
				isToUpper = false
			} else {
				if v == '_' {
					isToUpper = true
				} else {
					camelCase += string(v)
				}
			}
		}
	}
	return
}

// https://github.com/DaddyOh/golang-samples/blob/master/pad.go
/*
 * leftPad and rightPad just repoeat the padStr the indicated
 * number of times
 *
 */

func LeftPad(s string, padStr string, pLen int) string {
	return strings.Repeat(padStr, pLen) + s
}
func RightPad(s string, padStr string, pLen int) string {
	return s + strings.Repeat(padStr, pLen)
}

/* the Pad2Len functions are generally assumed to be padded with short sequences of strings
 * in many cases with a single character sequence
 *
 * so we assume we can build the string out as if the char seq is 1 char and then
 * just substr the string if it is longer than needed
 *
 * this means we are wasting some cpu and memory work
 * but this always get us to want we want it to be
 *
 * in short not optimized to for massive string work
 *
 * If the overallLen is shorter than the original string length
 * the string will be shortened to this length (substr)
 *
 */

func RightPad2Len(s string, padStr string, overallLen int) string {
	padCountInt := 1 + ((overallLen - len(padStr)) / len(padStr))
	var retStr = s + strings.Repeat(padStr, padCountInt)
	return retStr[:overallLen]
}
func LeftPad2Len(s string, padStr string, overallLen int) string {
	padCountInt := 1 + ((overallLen - len(padStr)) / len(padStr))
	var retStr = strings.Repeat(padStr, padCountInt) + s
	return retStr[(len(retStr) - overallLen):]
}

func StripIndent(multilineStr string) string {
	return strings.Replace(multilineStr, "\t", "", -1)
}
