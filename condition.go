// Package uiautomation 提供 Windows UI Automation 的 Go 语言封装
package uiautomation

import (
	"syscall"
	"unsafe"
)

// IUIAutomationCondition UI Automation 条件接口基类
// 用于创建元素搜索条件，是所有条件类型的基类
type IUIAutomationCondition struct {
	vtbl *IUnKnown
}

// Release 释放条件对象
// 返回: 引用计数
func (v *IUIAutomationCondition) Release() uint32 {
	return (*IUnKnown)(unsafe.Pointer(v)).Release()
}

// IUIAutomationConditionVtbl 条件接口虚函数表
type IUIAutomationConditionVtbl struct {
	IUnKnownVtbl
}

// IUIAutomationAndCondition AND 条件接口
// 用于组合多个条件，所有条件都必须满足
type IUIAutomationAndCondition struct {
	vtbl *IUIAutomationCondition
}

// Release 释放 AND 条件对象
// 返回: 引用计数
func (v *IUIAutomationAndCondition) Release() uint32 {
	return (*IUnKnown)(unsafe.Pointer(v)).Release()
}

// IUIAutomationAndConditionVtbl AND 条件接口虚函数表
type IUIAutomationAndConditionVtbl struct {
	IUIAutomationConditionVtbl

	Get_ChildCount           uintptr // 获取子条件数量
	GetChildren              uintptr // 获取子条件数组（SAFEARRAY）
	GetChildrenAsNativeArray uintptr // 获取子条件数组（原生数组）
}

// Get_ChildCount 获取 AND 条件中的子条件数量
// 参数: v - AND 条件接口
// 返回: 子条件数量
func Get_ChildCount(v *IUIAutomationAndCondition) int32 {
	var retVal int32
	syscall.SyscallN(
		(*IUIAutomationAndConditionVtbl)(unsafe.Pointer(v.vtbl)).Get_ChildCount,
		uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(&retVal)),
	)
	return retVal
}

// GetChildren 获取 AND 条件中的子条件数组（SAFEARRAY 格式）
// 参数: v - AND 条件接口
// 返回: SAFEARRAY 指针和可能的错误
func GetChildren(v *IUIAutomationAndCondition) (*TagSafeArray, error) {
	var retVal *TagSafeArray
	ret, _, _ := syscall.SyscallN(
		(*IUIAutomationAndConditionVtbl)(unsafe.Pointer(v.vtbl)).GetChildren,
		uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(&retVal)),
	)
	if ret != 0 {
		return nil, HResult(ret)
	}
	return retVal, nil
}

// GetChildrenAsNativeArray 获取 AND 条件中的子条件数组（原生数组格式）
// 参数: v - AND 条件接口
// 返回: 条件数组、数组长度和可能的错误
func GetChildrenAsNativeArray(v *IUIAutomationAndCondition) (*IUIAutomationCondition, int32, error) {
	var retVal *IUIAutomationCondition
	var retVal2 int32
	ret, _, _ := syscall.SyscallN(
		(*IUIAutomationAndConditionVtbl)(unsafe.Pointer(v.vtbl)).GetChildrenAsNativeArray,
		uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(&retVal)),
		uintptr(unsafe.Pointer(&retVal2)),
	)
	if ret != 0 {
		return nil, -1, HResult(ret)
	}
	return retVal, retVal2, nil
}

// IUIAutomationBoolCondition 布尔条件接口
// 用于创建 True 或 False 条件
type IUIAutomationBoolCondition struct {
	vtbl *IUIAutomationCondition
}

// Release 释放布尔条件对象
// 返回: 引用计数
func (v *IUIAutomationBoolCondition) Release() uint32 {
	return (*IUnKnown)(unsafe.Pointer(v)).Release()
}

// IUIAutomationBoolConditionVtbl 布尔条件接口虚函数表
type IUIAutomationBoolConditionVtbl struct {
	IUIAutomationConditionVtbl
	Get_BooleanValue uintptr // 获取布尔值
}

// Get_BooleanValue 获取布尔条件的值
// 参数: v - 布尔条件接口
// 返回: 布尔值（0=false, 非0=true）
func Get_BooleanValue(v *IUIAutomationBoolCondition) int32 {
	var retVal int32
	syscall.SyscallN(
		(*IUIAutomationBoolConditionVtbl)(unsafe.Pointer(v.vtbl)).Get_BooleanValue,
		uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(&retVal)),
	)
	return retVal
}

// IUIAutomationNotCondition NOT 条件接口
// 用于对条件取反
type IUIAutomationNotCondition struct {
	vtbl *IUIAutomationCondition
}

// Release 释放 NOT 条件对象
// 返回: 引用计数
func (v *IUIAutomationNotCondition) Release() uint32 {
	return (*IUnKnown)(unsafe.Pointer(v)).Release()
}

// IUIAutomationNotConditionVtbl NOT 条件接口虚函数表
type IUIAutomationNotConditionVtbl struct {
	IUIAutomationConditionVtbl
	GetChild uintptr // 获取子条件
}

// GetChild 获取 NOT 条件的子条件
// 参数: v - NOT 条件接口
// 返回: 子条件接口和可能的错误
func GetChild(v *IUIAutomationNotCondition) (*IUIAutomationCondition, error) {
	var retVal *IUIAutomationCondition
	ret, _, _ := syscall.SyscallN(
		(*IUIAutomationNotConditionVtbl)(unsafe.Pointer(v.vtbl)).GetChild,
		uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(&retVal)),
	)
	if ret != 0 {
		return nil, HResult(ret)
	}
	return retVal, nil
}

// IUIAutomationPropertyCondition 属性条件接口
// 用于根据元素属性值创建搜索条件
type IUIAutomationPropertyCondition struct {
	vtbl *IUIAutomationCondition
}

// Release 释放属性条件对象
// 返回: 引用计数
func (v *IUIAutomationPropertyCondition) Release() uint32 {
	return (*IUnKnown)(unsafe.Pointer(v)).Release()
}

// IUIAutomationPropertyConditionVtbl 属性条件接口虚函数表
type IUIAutomationPropertyConditionVtbl struct {
	IUIAutomationConditionVtbl

	Get_PropertyConditionFlags uintptr // 获取条件标志
	Get_PropertyId             uintptr // 获取属性ID
	Get_PropertyValue          uintptr // 获取属性值
}

// Get_PropertyConditionFlags 获取属性条件的标志
// 参数: v - 属性条件接口
// 返回: 条件标志指针
func Get_PropertyConditionFlags(v *IUIAutomationPropertyCondition) *PropertyConditionFlags {
	var retVal *PropertyConditionFlags
	syscall.SyscallN(
		(*IUIAutomationPropertyConditionVtbl)(unsafe.Pointer(v.vtbl)).Get_PropertyConditionFlags,
		uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(&retVal)),
	)
	return retVal
}

// Get_PropertyId 获取属性条件的属性ID
// 参数: v - 属性条件接口
// 返回: 属性ID指针
func Get_PropertyId(v *IUIAutomationPropertyCondition) *PropertyId {
	var retVal *PropertyId
	syscall.SyscallN(
		(*IUIAutomationPropertyConditionVtbl)(unsafe.Pointer(v.vtbl)).Get_PropertyId,
		uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(&retVal)),
	)
	return retVal
}

// Get_PropertyValue 获取属性条件的属性值
// 参数: v - 属性条件接口
// 返回: VARIANT 值指针
func Get_PropertyValue(v *IUIAutomationPropertyCondition) *VARIANT {
	var retVal *VARIANT
	syscall.SyscallN(
		(*IUIAutomationPropertyConditionVtbl)(unsafe.Pointer(v.vtbl)).Get_PropertyValue,
		uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(&retVal)),
	)
	return retVal
}
