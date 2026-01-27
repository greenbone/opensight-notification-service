package models

type ValidationErrors map[string]string

func (v ValidationErrors) Error() string {
	return "validation error"
}
