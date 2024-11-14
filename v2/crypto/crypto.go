package crypto

type Signer[T any] interface {
	Sign(msg []byte) T
}
