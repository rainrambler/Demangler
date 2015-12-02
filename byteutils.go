package main

import (
)

func isNumberChar(c byte) bool {
	return (c >= '0') && (c <= '9')
}

func isUpperChar(c byte) bool {
	return (c >= 'A') && (c <= 'Z')
}

func isNonZeroNumberChar(c byte) bool {
	return (c >= '1') && (c <= '9')
}
