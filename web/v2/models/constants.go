package models

import "regexp"

const (
	EmailFormat string = "email"
	PhoneFormat string = "phone"
)

var EmailRegex = regexp.MustCompile(`[^@ \t\r\n]+@[^@ \t\r\n]+\.[^@ \t\r\n]+`)
var PhoneRegex = regexp.MustCompile(`(?:(?:\+|00)?(55)\s?)?(?:\(?([1-9][0-9])\)?\s?)(?:((?:9\d|[2-9])\d{3})\-?(\d{4}))`)
var WhateverRegex = regexp.MustCompile(`.*`)
