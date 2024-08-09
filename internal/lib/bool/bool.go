package boolean

type Wrapper struct {
	Value bool
}

func (b Wrapper) Int() int {
	if b.Value {
		return 1
	}
	return 0
}
