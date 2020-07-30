package acn

type err string

func (e err) Error() string {
	return string(e)
}

const (
	// ErrNotFound is returned when a node tree is not found
	ErrNotFound err = "not found"
)
