// Package uiautomation 提供 Windows UI Automation 的 Go 语言封装
package uiautomation

import (
	"syscall"
	"unsafe"
)

// IDispatch 调度接口
// COM 中的标准接口，用于支持动态调用和自动化
type IDispatch struct {
	vtbl *IUnKnown
}

// Release 释放调度接口对象
// 返回: 引用计数
func (v *IDispatch) Release() uint32 {
	return (*IUnKnown)(unsafe.Pointer(v)).Release()
}

// IDispatchVtbl 调度接口虚函数表
type IDispatchVtbl struct {
	IUnKnownVtbl
	GetIDsOfNames    uintptr // 获取名称ID
	GetTypeInfo      uintptr // 获取类型信息
	GetTypeInfoCount uintptr // 获取类型信息数量
	Invoke           uintptr // 调用方法
}

// GetIDsOfNames 获取方法或属性的名称ID
// 参数:
//   - v: IDispatch 接口
//   - in: GUID
//   - in2: 名称数量
//   - in3: 区域设置ID
//   - in4: 名称数组
//
// 返回: 名称ID和可能的错误
func GetIDsOfNames(v *IDispatch, in *syscall.GUID, in2 uint16, in3, in4 uint32) (int32, error) {
	var retVal int32
	ret, _, _ := syscall.SyscallN(
		(*IDispatchVtbl)(unsafe.Pointer(v.vtbl)).GetIDsOfNames,
		uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(in)),
		uintptr(in2),
		uintptr(in3),
		uintptr(in4),
		uintptr(unsafe.Pointer(&retVal)),
	)
	if ret != 0 {
		return -1, HResult(ret)
	}
	return retVal, nil
}

// GetTypeInfo 获取类型信息
// 参数:
//   - v: IDispatch 接口
//   - in: 类型信息索引
//   - in2: 区域设置ID
//
// 返回: 类型信息接口和可能的错误
// 注意: ITypeInfo 类型未定义，此函数暂时注释
/*
func GetTypeInfo(v *IDispatch, in, in2 uint32) (*ITypeInfo, error) {
	var retVal *ITypeInfo
	ret, _, _ := syscall.SyscallN(
		(*IDispatchVtbl)(unsafe.Pointer(v.vtbl)).GetTypeInfo,
		uintptr(unsafe.Pointer(v)),
		uintptr(in),
		uintptr(in2),
		uintptr(unsafe.Pointer(&retVal)),
	)
	if ret != 0 {
		return nil, HResult(ret)
	}
	return retVal, nil
}
*/