package models

import (
	"fmt"
	"regexp"
	"time"
	"unicode/utf8"

	"github.com/cayo-rodrigues/nff/web/utils"
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

type Fields []*struct {
	Name  string
	Value any
	Rules []*RuleSet
}
type ErrorMessages map[string]string

func Rules(ruleFuncs ...RuleFunc) []*RuleSet {
	ruleSets := make([]*RuleSet, len(ruleFuncs))
	for i := 0; i < len(ruleSets); i++ {
		ruleSets[i] = ruleFuncs[i]()
	}

	return ruleSets
}

func Validate(fields Fields) (ErrorMessages, bool) {
	var messages ErrorMessages

	for _, field := range fields {
		for _, rs := range field.Rules {
			if rs.ValidateFunc != nil {
				rs.FieldValue = field.Value
				isValid := rs.ValidateFunc(rs)
				if !isValid {
					if messages == nil {
						messages = make(ErrorMessages)
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

func HasValue(val any) bool {
	switch val := val.(type) {
	case string:
		return utf8.RuneCountInString(val) > 0
	case int:
		return val != 0
	case float64:
		return val != 0
	case time.Time:
		return !val.IsZero()
	}

	return false
}

func Required() *RuleSet {
	return &RuleSet{
		MessageFunc: func(rs *RuleSet) string {
			return utils.MandatoryFieldMsg
		},
		ValidateFunc: func(rs *RuleSet) bool {
			return HasValue(rs.FieldValue)
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

func Match(regexes ...*regexp.Regexp) RuleFunc {
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

				for _, regex := range regexes {
					match := regex.MatchString(str)
					if match {
						return true
					}
				}

				return false
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

func RequiredUnless(vals ...any) RuleFunc {
	return func() *RuleSet {
		return &RuleSet{
			MessageFunc: func(rs *RuleSet) string {
				return utils.MandatoryFieldMsg
			},
			ValidateFunc: func(rs *RuleSet) bool {
				if HasValue(rs) {
					return true
				}
				for _, val := range vals {
					if HasValue(val) {
						return true
					}
				}
				return false
			},
		}
	}
}

func NotAfter(dt time.Time) RuleFunc {
	return func() *RuleSet {
		return &RuleSet{
			MessageFunc: func(rs *RuleSet) string {
				return utils.IlogicalDatesMsg
			},
			ValidateFunc: func(rs *RuleSet) bool {
				switch val := rs.FieldValue.(type) {
				case time.Time:
					return !val.After(dt)
				}

				return false
			},
		}
	}
}

func MaxTimeRange(dt time.Time, days int) RuleFunc {
	return func() *RuleSet {
		return &RuleSet{
			MessageFunc: func(rs *RuleSet) string {
				return utils.TimeRangeTooLongMsg
			},
			ValidateFunc: func(rs *RuleSet) bool {
				switch val := rs.FieldValue.(type) {
				case time.Time:
					return int(dt.Sub(val).Hours()/24) > days
				}

				return false
			},
		}
	}
}
