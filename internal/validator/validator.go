package validator

import (
	"strings"
	"unicode/utf8"
)

type Validator struct {
	FieldErrors map[string]string
}

func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}


//Adding an error to the FieldsError map
func (v *Validator) AddFieldError(key,message string)  {
	// Create FieldErrors if it doesn't exist
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}
	if _,exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

func (v *Validator) CheckField(ok bool,key,message string){
	if !ok {
		v.AddFieldError(key,message)
	}
}

func NotBlank(s string) bool {
	return strings.TrimSpace(s) != ""
}

func MaxChars(s string, n int) bool {
	return utf8.RuneCountInString(s) <= n
}

func PermittedInt(value int, allowedVals ...int) bool {
	for i:= range allowedVals {
		if value == allowedVals[i] {
			return true
		}
	}
	return false
}