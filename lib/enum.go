package lib

type EmptyEnum struct {
}

func (e EmptyEnum) MarshalBCS() ([]byte, error) {
	return []byte{}, nil
}
