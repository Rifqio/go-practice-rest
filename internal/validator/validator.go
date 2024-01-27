package validator

type Validator struct {
	Errors map[string]string
}

// Valid returns true if there are no errors, otherwise it returns false.
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddFieldError(key, message string) {
	// Initialize the map if it nil.
	if v.Errors == nil {
		v.Errors = make(map[string]string)
	}

	if _, ok := v.Errors[key]; !ok {
		v.Errors[key] = message
	}
}

func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

func NotEmpty(value string) bool {
	return len(value) > 0
}

func MaxChars(value string, max int) bool {
	return len(value) <= max
}