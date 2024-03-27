package models

import (
	"fmt"
	"regexp"
	"time"
	"unicode/utf8"

	"github.com/cayo-rodrigues/nff/web/utils"
	"github.com/google/uuid"
)

type RuleSet struct {
	FieldValue   any
	MessageFunc  func(*RuleSet) string
	ValidateFunc func(*RuleSet) bool
}

func (rs *RuleSet) WithMessage(msg string) RuleFunc {
	rs.MessageFunc = func(rs *RuleSet) string {
		return msg
	}

	return func() *RuleSet {
		return rs
	}
}

type RuleFunc func() *RuleSet

type Field struct {
	Name  string
	Value any
	Rules []*RuleSet
}
type Messages map[string]string

func Rules(ruleFuncs ...RuleFunc) []*RuleSet {
	ruleSets := make([]*RuleSet, len(ruleFuncs))
	for i := 0; i < len(ruleSets); i++ {
		ruleSets[i] = ruleFuncs[i]()
	}

	return ruleSets
}

func Validate(fields []*Field) (Messages, bool) {
	var messages Messages

	for _, field := range fields {
		for _, rs := range field.Rules {
			if rs.ValidateFunc != nil {
				rs.FieldValue = field.Value
				isValid := rs.ValidateFunc(rs)
				if !isValid {
					if messages == nil {
						messages = make(Messages)
					}
					msg := rs.MessageFunc(rs)
					messages[field.Name] = msg
					break // stop runing validate funcs after first fail
				}
			}
		}
	}

	return messages, len(messages) == 0
}

func Required() *RuleSet {
	return &RuleSet{
		MessageFunc: func(rs *RuleSet) string {
			return utils.MandatoryFieldMsg
		},
		ValidateFunc: func(rs *RuleSet) bool {
			switch val := rs.FieldValue.(type) {
			case string:
				return utf8.RuneCountInString(val) > 0
			case time.Time:
				return !val.IsZero()
			case uuid.UUID:
				return val != uuid.Nil
			}

			return false
		},
	}
}

func Email() *RuleSet {
	return &RuleSet{
		MessageFunc: func(rs *RuleSet) string {
			return utils.InvalidFormatMsg
		},
		ValidateFunc: func(rs *RuleSet) bool {
			str, ok := rs.FieldValue.(string)
			if !ok {
				return false
			}

			if str == "" {
				return true
			}

			return EmailRegex.MatchString(str)
		},
	}
}

func Phone() *RuleSet {
	return &RuleSet{
		MessageFunc: func(rs *RuleSet) string {
			return utils.InvalidFormatMsg
		},
		ValidateFunc: func(rs *RuleSet) bool {
			str, ok := rs.FieldValue.(string)
			if !ok {
				return false
			}

			if str == "" {
				return true
			}

			return PhoneRegex.MatchString(str)
		},
	}
}

func Format(format string) RuleFunc {
	return func() *RuleSet {
		return &RuleSet{
			MessageFunc: func(rs *RuleSet) string {
				return utils.InvalidFormatMsg
			},
			ValidateFunc: func(rs *RuleSet) bool {
				str, ok := rs.FieldValue.(string)
				if !ok {
					return false
				}

				if str == "" {
					return true
				}

				var regex *regexp.Regexp

				switch format {
				case EmailFormat:
					regex = EmailRegex
				case PhoneFormat:
					regex = PhoneRegex
				default:
					regex = WhateverRegex
				}

				return regex.MatchString(str)
			},
		}
	}
}

func Min(minValue int) RuleFunc {
	return func() *RuleSet {
		return &RuleSet{
			MessageFunc: func(rs *RuleSet) string {
				switch rs.FieldValue.(type) {
				case int, float32, float64:
					return fmt.Sprintf("Valor mínimo: %d", minValue)
				default:
					return fmt.Sprintf("Mínimo de %d caracteres", minValue)
				}
			},
			ValidateFunc: func(rs *RuleSet) bool {
				switch val := rs.FieldValue.(type) {
				case int:
					return val >= minValue
				case string:
					if val == "" {
						return true
					}
					return utf8.RuneCountInString(val) >= minValue
				}

				return false
			},
		}
	}
}

func Max(maxValue int) RuleFunc {
	return func() *RuleSet {
		return &RuleSet{
			MessageFunc: func(rs *RuleSet) string {
				switch rs.FieldValue.(type) {
				case int, float32, float64:
					return fmt.Sprintf("Valor máximo: %d", maxValue)
				default:
					return fmt.Sprintf("Máximo de %d caracteres", maxValue)
				}
			},
			ValidateFunc: func(rs *RuleSet) bool {
				switch val := rs.FieldValue.(type) {
				case int:
					return val <= maxValue
				case string:
					if val == "" {
						return true
					}
					return utf8.RuneCountInString(val) <= maxValue
				}

				return false
			},
		}
	}
}

func NotOneOf(vals ...any) RuleFunc {
	return func() *RuleSet {
		return &RuleSet{
			MessageFunc: func(rs *RuleSet) string {
				return utils.UnacceptableValueMsg
			},
			ValidateFunc: func(rs *RuleSet) bool {
				for _, val := range vals {
					if val == rs.FieldValue {
						return false
					}
				}
				return true
			},
		}
	}
}

func OneOf(vals ...any) RuleFunc {
	return func() *RuleSet {
		return &RuleSet{
			MessageFunc: func(rs *RuleSet) string {
				return utils.UnacceptableValueMsg
			},
			ValidateFunc: func(rs *RuleSet) bool {
				for _, val := range vals {
					if val == rs.FieldValue {
						return true
					}
				}
				return false
			},
		}
	}
}

func Datetime() *RuleSet {
	return &RuleSet{
		MessageFunc: func(rs *RuleSet) string {
			return utils.UnacceptableValueMsg
		},
		ValidateFunc: func(rs *RuleSet) bool {
			_, isDatetime := rs.FieldValue.(time.Time)
			return isDatetime
		},
	}
}
