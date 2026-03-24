//go:build amd64

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
		uintptr(unsafe.Pointer(&value)),
		uintptr(unsafe.Pointer(&retVal)),
	)
	if ret != 0 {
		return nil, HResult(ret)
	}
	return retVal, nil
}
