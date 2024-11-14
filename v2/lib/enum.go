package lib

import "io"

type EmptyEnum struct {
}

func (e EmptyEnum) MarshalBCS() ([]byte, error) {
	return []byte{}, nil
}

func (e *EmptyEnum) UnmarshalBCS(io.Reader) (int, error) {
	return 0, nil
}
