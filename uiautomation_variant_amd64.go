//go:build amd64

// Package uiautomation 提供 Windows UI Automation 的 Go 语言封装
package uiautomation

import (
	"syscall"
	"unsafe"
)

// CreatePropertyCondition 创建属性条件（64位版本）
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
		uintptr(unsafe.Pointer(&value)),
		uintptr(unsafe.Pointer(&retVal)),
	)
	if ret != 0 {
		return nil, HResult(ret)
	}
	return retVal, nil
}