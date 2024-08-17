package jerr

import "log"

type Wrapper struct {
	Value error
}

func Wrap(err error) Wrapper {
	return Wrapper{Value: err}
}

func (w Wrapper) Unwrap() error {
	return w.Value
}

func (w Wrapper) Error() string {
	return w.Value.Error()
}

func (w Wrapper) Catch() {
	if w.Value != nil {
		log.Fatal(w.Value)
	}
}
