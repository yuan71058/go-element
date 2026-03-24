//go:build 386

package uiautomation

import (
	"syscall"
	"unsafe"
)

func (v *IUIAutomation) CreatePropertyCondition(id PropertyId, value VARIANT) (*IUIAutomationCondition, error) {
	var retVal *IUIAutomationCondition
	ret, _, _ := syscall.SyscallN(
		(*IUIAutomationVtbl)(unsafe.Pointer(v.vtbl)).CreatePropertyCondition,
		uintptr(unsafe.Pointer(v)),
		uintptr(id),
		uintptr(uint32(value.VT)|(uint32(value.wReserved1)<<16)),
		uintptr(uint32(value.wReserved2)|(uint32(value.wReserved3)<<16)),
		uintptr(uint32(value.Val&0xFFFFFFFF)),
		uintptr(uint32(value.Val>>32)),
		uintptr(unsafe.Pointer(&retVal)),
	)
	if ret != 0 {
		return nil, HResult(ret)
	}
	return retVal, nil
}
