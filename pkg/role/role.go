package role

const (
	ADMIN = iota
	USER
)

func Set(f bool) int {
	if f {
		return ADMIN
	}
	return USER
}
