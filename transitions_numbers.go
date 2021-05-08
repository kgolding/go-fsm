package fsm

import (
	"encoding/binary"
	"io"
)

// Uint read 1 or more bytes as an Integer
func Uint(v *uint, byteCount int) TransitionTest {
	return func(b []byte) (int, error) {
		if len(b) < byteCount {
			return 0, io.EOF
		}
		*v = 0
		for i := 0; i < byteCount; i++ {
			*v += uint(b[byteCount-i-1]) << (i * 8)
		}
		return byteCount, nil
	}
}

func Int8(v *int8) TransitionTest {
	return func(b []byte) (int, error) {
		if len(b) < 1 {
			return 0, io.EOF
		}
		*v = int8(b[0])
		return 1, nil
	}
}

func Uint8(v *uint8) TransitionTest {
	return func(b []byte) (int, error) {
		if len(b) < 1 {
			return 0, io.EOF
		}
		*v = uint8(b[0])
		return 1, nil
	}
}

func Int16(v *int16) TransitionTest {
	return func(b []byte) (int, error) {
		if len(b) < 2 {
			return 0, io.EOF
		}
		*v = int16(binary.BigEndian.Uint16(b))
		return 2, nil
	}
}

func Uint16(v *uint16) TransitionTest {
	return func(b []byte) (int, error) {
		if len(b) < 2 {
			return 0, io.EOF
		}
		*v = binary.BigEndian.Uint16(b)
		return 2, nil
	}
}

func Int16LE(v *int16) TransitionTest {
	return func(b []byte) (int, error) {
		if len(b) < 2 {
			return 0, io.EOF
		}
		*v = int16(binary.LittleEndian.Uint16(b))
		return 2, nil
	}
}

func Uint16LE(v *uint16) TransitionTest {
	return func(b []byte) (int, error) {
		if len(b) < 2 {
			return 0, io.EOF
		}
		*v = binary.LittleEndian.Uint16(b)
		return 2, nil
	}
}

func Int32(v *int32) TransitionTest {
	return func(b []byte) (int, error) {
		if len(b) < 4 {
			return 0, io.EOF
		}
		*v = int32(binary.BigEndian.Uint32(b))
		return 4, nil
	}
}

func Uint32(v *uint32) TransitionTest {
	return func(b []byte) (int, error) {
		if len(b) < 4 {
			return 0, io.EOF
		}
		*v = binary.BigEndian.Uint32(b)
		return 4, nil
	}
}

func Int32LE(v *int32) TransitionTest {
	return func(b []byte) (int, error) {
		if len(b) < 4 {
			return 0, io.EOF
		}
		*v = int32(binary.LittleEndian.Uint32(b))
		return 4, nil
	}
}

func Uint32LE(v *uint32) TransitionTest {
	return func(b []byte) (int, error) {
		if len(b) < 4 {
			return 0, io.EOF
		}
		*v = binary.LittleEndian.Uint32(b)
		return 4, nil
	}
}
