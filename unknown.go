// Package uiautomation 提供 Windows UI Automation 的 Go 语言封装
package uiautomation

import (
	"fmt"
	"syscall"
	"unsafe"
)

// IUnKnown COM IUnknown 接口
// 所有 COM 接口的基类，提供引用计数和接口查询功能
type IUnKnown struct {
	Vtbl *IUnKnownVtbl
}

// IUnKnownVtbl IUnknown 接口虚函数表
type IUnKnownVtbl struct {
	QueryInterface uintptr // 查询接口
	AddRef         uintptr // 增加引用计数
	Release        uintptr // 释放引用计数
}

// NewIUnKnown 从指针创建 IUnKnown 实例
// 参数: v - COM 对象指针
// 返回: IUnKnown 实例
func NewIUnKnown(v unsafe.Pointer) *IUnKnown {
	return (*IUnKnown)(v)
}

// AddRef 增加对象引用计数
// 返回: 新的引用计数
func (v *IUnKnown) AddRef() uint32 {
	ret, _, _ := syscall.SyscallN(
		v.Vtbl.AddRef,
		uintptr(unsafe.Pointer(v)),
	)
	return uint32(ret)
}

// Release 释放对象引用计数
// 当引用计数为 0 时，对象会被销毁
// 返回: 新的引用计数
func (v *IUnKnown) Release() uint32 {
	if v == nil || v.Vtbl == nil {
		return 0
	}
	ret, _, _ := syscall.SyscallN(
		v.Vtbl.Release,
		uintptr(unsafe.Pointer(v)),
	)
	return uint32(ret)
}

// AddRef 增加对象引用计数（包级函数）
// 参数: v - IUnKnown 接口
// 返回: 新的引用计数
func AddRef(v *IUnKnown) uint32 {
	return v.AddRef()
}

// Release 释放对象引用计数（包级函数）
// 参数: v - IUnKnown 接口
// 返回: 新的引用计数
func Release(v *IUnKnown) uint32 {
	return v.Release()
}

// QueryInterface 查询对象是否支持指定接口
// 参数:
//   - v: IUnKnown 接口
//   - riid: 接口标识符（IID）
// 返回:
//   - unsafe.Pointer: 接口指针
//   - error: 查询失败时返回错误
func QueryInterface(v *IUnKnown, riid *syscall.GUID) (unsafe.Pointer, error) {
	var retVal unsafe.Pointer
	ret, _, _ := syscall.SyscallN(
		(*IUnKnownVtbl)(unsafe.Pointer(v.Vtbl)).QueryInterface,
		uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(riid)),
		uintptr(unsafe.Pointer(&retVal)),
	)
	if ret != 0 {
		return nil, HResult(ret)
	}
	return retVal, nil
}

// HResult 将 HRESULT 转换为错误
// 参数: ret - HRESULT 值
// 返回: 错误对象，成功返回 nil
func HResult(ret uintptr) error {
	if ret == 0 {
		return nil
	}
	// error info https://pkg.go.dev/golang.org/x/sys/windows
	return fmt.Errorf("COM Error: 0x%08X", uint32(ret))
}

// UnKnownToUintptr 将 IUnKnown 接口转换为 uintptr
// 用于 COM 方法调用
// 参数: obj - 接口对象
// 返回: 指针值
func UnKnownToUintptr(obj interface{}) uintptr {
	// 获取接口内部的指针
	type iface struct {
		typ uintptr
		ptr uintptr
	}
	return (*iface)(unsafe.Pointer(&obj)).ptr
}
