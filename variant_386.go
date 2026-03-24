//go:build 386

package uiautomation

type VARIANT struct {
	VT         TagVarenum
	wReserved1 uint16
	wReserved2 uint16
	wReserved3 uint16
	Val        int64
}
