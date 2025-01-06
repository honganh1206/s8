// To trace parser function/method invocations
package parser

import (
	"fmt"
	"strings"
)

// Keep track of the current indentation level
var traceLevel int = 0

const traceIndentPlaceholder string = "\t" // Insert a horizontal tab space when used in a string

// Create indentation string based on current trace level
// Return a string with appropariate number of tabs
func indentLevel() string {
	return strings.Repeat(traceIndentPlaceholder, traceLevel)
}

func tracePrint(fs string) {
	fmt.Printf("%s%s\n", indentLevel(), fs)
}

func incIdent() { traceLevel = traceLevel + 1 }
func decIdent() { traceLevel = traceLevel - 1 }

// Mark the beginning of a parser operation
func trace(msg string) string {
	incIdent()
	tracePrint("BEGIN " + msg)
	return msg
}

// Mark the end of a parser operation
func untrace(msg string) {
	tracePrint("END " + msg)
	decIdent()
}
