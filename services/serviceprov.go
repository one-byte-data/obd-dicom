package services

// SCP - Interface to scp
type SCP interface {

}

type scp struct {

}

// NewSCP - Creates an interface to scu
func NewSCP() SCP {
	return &scp{}
}
