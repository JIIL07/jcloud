package role

const (
	NOTUSER = iota
	USER
	ADMIN
)

func Set(f bool) int {
	if f {
		return USER
	}
	return NOTUSER
}
