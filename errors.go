package acn

type acnErr string

func (a acnErr) Error() string {
	return string(a)
}

const (
	// ErrNodeNotFound is returned whenever a node was not found
	ErrNodeNotFound acnErr = "could not find matching ast node"
)