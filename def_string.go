package gocc

type DefString struct {
	label string
	data  string
}

func NewDefString(label, data string) *DefString {
	return &DefString{
		label: label,
		data:  data,
	}
}
