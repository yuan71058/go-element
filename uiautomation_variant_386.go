//go:build 386

// Package uiautomation 提供 Windows UI Automation 的 Go 语言封装
package uiautomation

import (
	"syscall"
	"unsafe"
)

// CreatePropertyCondition 创建属性条件（32位版本）
// 用于根据属性值创建查找条件
// 参数:
//   - id: 属性ID
//   - value: 属性值（VARIANT）
//
// 返回: 条件接口和可能的错误
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