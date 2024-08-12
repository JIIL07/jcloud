package params

type Params map[string]interface{}

func NewParams() Params {
	return make(Params)
}

func (p Params) Set(key string, value interface{}) {
	p[key] = value
}

func (p Params) Get(key string) interface{} {
	return p[key]
}

func (p Params) String(key string) string {
	if value, ok := p[key].(string); ok {
		return value
	}
	return ""
}

func (p Params) Int(key string) int {
	if value, ok := p[key].(int); ok {
		return value
	}
	return 0
}

func (p Params) Bool(key string) bool {
	if value, ok := p[key].(bool); ok {
		return value
	}
	return false
}
