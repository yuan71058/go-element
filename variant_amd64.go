//go:build amd64

package uiautomation

type VARIANT struct {
	VT         TagVarenum
	wReserved1 uint16
	wReserved2 uint16
	wReserved3 uint16
	_          [4]byte
	Val        int64
}
