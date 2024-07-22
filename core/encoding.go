package core

import "io"

// we pass Writer and Reader rather than bytes because we want to "stream it" -> we can pipe a connection into this encoder instead of buffering up the memory itself

type Encoder[T any] interface {
	Encode(w io.Writer, v T) error
}

type Decoder[T any] interface {
	Decode(r io.Reader, v T) error
}

