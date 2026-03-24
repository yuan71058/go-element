// Package uiautomation 提供 Windows UI Automation 的 Go 语言封装
package uiautomation

import (
	"errors"
	"syscall"
	"unsafe"
)

var (
	// ErrorBstrPointerNil BSTR 指针为空错误
	ErrorBstrPointerNil = errors.New("BSTR pointer is nil")
)

// IAnnotationProvider 注解提供者接口
// Windows 8 开始支持，用于 UIA_AnnotationPatternId
type IAnnotationProvider struct {
	vtbl *IUnKnown
}

// IAnnotationProviderVtbl 注解提供者虚函数表
type IAnnotationProviderVtbl struct {
	IUnKnownVtbl
	Get_AnnotationTypeId   uintptr // 获取注解类型ID
	Get_AnnotationTypeName uintptr // 获取注解类型名称
	Get_Author             uintptr // 获取作者
	Get_DateTime           uintptr // 获取日期时间
	Get_Target             uintptr // 获取目标
}

// IRawElementProviderSimple 简单原始元素提供者接口
type IRawElementProviderSimple struct {
	vtbl *IUnKnown
}

// IRawElementProviderSimpleVtbl 简单原始元素提供者虚函数表
type IRawElementProviderSimpleVtbl struct {
	IUnKnownVtbl
	Get_HostRawElementProvider uintptr // 获取宿主原始元素提供者
	Get_ProviderOptions        uintptr // 获取提供者选项
	GetPatternProvider         uintptr // 获取模式提供者
	GetPropertyValue           uintptr // 获取属性值
}

// NewIAnnotationProvider 创建注解提供者实例
// 参数: unk - IUnKnown 接口
// 返回: 注解提供者实例
func NewIAnnotationProvider(unk *IUnKnown) *IAnnotationProvider {
	return newIAnnotationProvider(unk)
}

// newIAnnotationProvider 内部函数：从 IUnKnown 创建注解提供者
func newIAnnotationProvider(unk *IUnKnown) *IAnnotationProvider {
	return (*IAnnotationProvider)(unsafe.Pointer(unk))
}

// Get_AnnotationTypeId 获取注解类型ID
// 返回: 注解类型ID
func (v *IAnnotationProvider) Get_AnnotationTypeId() int32 {
	var retVal int32
	syscall.SyscallN(
		(*IAnnotationProviderVtbl)(unsafe.Pointer(v.vtbl)).Get_AnnotationTypeId,
		uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(&retVal)),
	)
	return retVal
}