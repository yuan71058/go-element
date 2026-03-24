// Package uiautomation 提供 Windows UI Automation 的 Go 语言封装
package uiautomation

import (
	"syscall"
	"unsafe"
)

// IUIAutomationInvokePattern 调用模式接口
// 用于调用按钮等可点击控件
type IUIAutomationInvokePattern struct {
	vtbl *IUIAutomationInvokePatternVtbl
}

// IUIAutomationInvokePatternVtbl 调用模式虚函数表
type IUIAutomationInvokePatternVtbl struct {
	IUnKnownVtbl
	Invoke uintptr // 调用控件
}

// Invoke 调用控件（如点击按钮）
// 返回: 错误信息
func (v *IUIAutomationInvokePattern) Invoke() error {
	ret, _, _ := syscall.SyscallN(
		v.vtbl.Invoke,
		uintptr(unsafe.Pointer(v)),
	)
	if ret != 0 {
		return HResult(ret)
	}
	return nil
}

// Release 释放调用模式对象
// 返回: 引用计数
func (v *IUIAutomationInvokePattern) Release() uint32 {
	return (*IUnKnown)(unsafe.Pointer(v)).Release()
}

// newIUIAutomationInvokePattern 内部函数：从 IUnKnown 创建调用模式
func newIUIAutomationInvokePattern(unk *IUnKnown) *IUIAutomationInvokePattern {
	return (*IUIAutomationInvokePattern)(unsafe.Pointer(unk))
}

// NewIUIAutomationInvokePattern 创建调用模式实例
// 参数: unk - IUnKnown 接口
// 返回: 调用模式实例
func NewIUIAutomationInvokePattern(unk *IUnKnown) *IUIAutomationInvokePattern {
	return newIUIAutomationInvokePattern(unk)
}

// IUIAutomationValuePattern 值模式接口
// 用于设置或获取文本输入控件的值
type IUIAutomationValuePattern struct {
	vtbl *IUIAutomationValuePatternVtbl
}

// IUIAutomationValuePatternVtbl 值模式虚函数表
type IUIAutomationValuePatternVtbl struct {
	IUnKnownVtbl
	SetValue              uintptr // 设置值
	Get_CurrentValue      uintptr // 获取当前值
	Get_CurrentIsReadOnly uintptr // 获取是否只读
	Get_CachedValue       uintptr // 获取缓存值
	Get_CachedIsReadOnly  uintptr // 获取缓存是否只读
}

// SetValue 设置控件的值
// 参数: val - 要设置的值
// 返回: 错误信息
func (v *IUIAutomationValuePattern) SetValue(val string) error {
	bstr, err := string2Bstr(val)
	if err != nil {
		return err
	}
	defer procSysFreeString.Call(bstr)
	ret, _, _ := syscall.SyscallN(
		v.vtbl.SetValue,
		uintptr(unsafe.Pointer(v)),
		uintptr(bstr),
	)
	if ret != 0 {
		return HResult(ret)
	}
	return nil
}

// Get_CurrentValue 获取控件的当前值
// 返回: 当前值和可能的错误
func (v *IUIAutomationValuePattern) Get_CurrentValue() (string, error) {
	var bstr uintptr
	ret, _, _ := syscall.SyscallN(
		v.vtbl.Get_CurrentValue,
		uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(&bstr)),
	)
	if ret != 0 {
		return "", HResult(ret)
	}
	defer procSysFreeString.Call(bstr)
	return bstr2str(bstr), nil
}

// newIUIAutomationValuePattern 内部函数：从 IUnKnown 创建值模式
func newIUIAutomationValuePattern(unk *IUnKnown) *IUIAutomationValuePattern {
	return (*IUIAutomationValuePattern)(unsafe.Pointer(unk))
}

// NewIUIAutomationValuePattern 创建值模式实例
// 参数: unk - IUnKnown 接口
// 返回: 值模式实例
func NewIUIAutomationValuePattern(unk *IUnKnown) *IUIAutomationValuePattern {
	return newIUIAutomationValuePattern(unk)
}

// Release 释放值模式对象
// 返回: 引用计数
func (v *IUIAutomationValuePattern) Release() uint32 {
	return (*IUnKnown)(unsafe.Pointer(v)).Release()
}

// IUIAutomationSelectionItemPattern 选择项模式接口
// 用于选择列表项等控件
type IUIAutomationSelectionItemPattern struct {
	vtbl *IUIAutomationSelectionItemPatternVtbl
}

// IUIAutomationSelectionItemPatternVtbl 选择项模式虚函数表
type IUIAutomationSelectionItemPatternVtbl struct {
	IUnKnownVtbl
	Select                        uintptr // 选择项
	AddToSelection                uintptr // 添加到选择
	RemoveFromSelection           uintptr // 从选择移除
	Get_CurrentIsSelected         uintptr // 获取是否已选中
	Get_CurrentSelectionContainer uintptr // 获取选择容器
	Get_CachedIsSelected          uintptr // 获取缓存是否已选中
	Get_CachedSelectionContainer  uintptr // 获取缓存选择容器
}

// Select 选择该项
// 返回: 错误信息
func (v *IUIAutomationSelectionItemPattern) Select() error {
	ret, _, _ := syscall.SyscallN(
		v.vtbl.Select,
		uintptr(unsafe.Pointer(v)),
	)
	if ret != 0 {
		return HResult(ret)
	}
	return nil
}

// Release 释放选择项模式对象
// 返回: 引用计数
func (v *IUIAutomationSelectionItemPattern) Release() uint32 {
	return (*IUnKnown)(unsafe.Pointer(v)).Release()
}

// IUIAutomationTogglePattern 切换模式接口
// 用于切换复选框等控件的状态
type IUIAutomationTogglePattern struct {
	vtbl *IUIAutomationTogglePatternVtbl
}

// IUIAutomationTogglePatternVtbl 切换模式虚函数表
type IUIAutomationTogglePatternVtbl struct {
	IUnKnownVtbl
	Toggle                 uintptr // 切换状态
	Get_CurrentToggleState uintptr // 获取当前切换状态
	Get_CachedToggleState  uintptr // 获取缓存切换状态
}

// Toggle 切换控件状态
// 返回: 错误信息
func (v *IUIAutomationTogglePattern) Toggle() error {
	ret, _, _ := syscall.SyscallN(
		v.vtbl.Toggle,
		uintptr(unsafe.Pointer(v)),
	)
	if ret != 0 {
		return HResult(ret)
	}
	return nil
}

// Get_CurrentToggleState 获取当前切换状态
// 返回: 切换状态和可能的错误
func (v *IUIAutomationTogglePattern) Get_CurrentToggleState() (ToggleState, error) {
	var state ToggleState
	ret, _, _ := syscall.SyscallN(
		v.vtbl.Get_CurrentToggleState,
		uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(&state)),
	)
	if ret != 0 {
		return 0, HResult(ret)
	}
	return state, nil
}

// newIUIAutomationTogglePattern 内部函数：从 IUnKnown 创建切换模式
func newIUIAutomationTogglePattern(unk *IUnKnown) *IUIAutomationTogglePattern {
	return (*IUIAutomationTogglePattern)(unsafe.Pointer(unk))
}

// NewIUIAutomationTogglePattern 创建切换模式实例
// 参数: unk - IUnKnown 接口
// 返回: 切换模式实例
func NewIUIAutomationTogglePattern(unk *IUnKnown) *IUIAutomationTogglePattern {
	return newIUIAutomationTogglePattern(unk)
}

// Release 释放切换模式对象
// 返回: 引用计数
func (v *IUIAutomationTogglePattern) Release() uint32 {
	return (*IUnKnown)(unsafe.Pointer(v)).Release()
}

// IUIAutomationExpandCollapsePattern 展开/折叠模式接口
// 用于展开或折叠树节点、菜单等控件
type IUIAutomationExpandCollapsePattern struct {
	vtbl *IUIAutomationExpandCollapsePatternVtbl
}

// IUIAutomationExpandCollapsePatternVtbl 展开/折叠模式虚函数表
type IUIAutomationExpandCollapsePatternVtbl struct {
	IUnKnownVtbl
	Expand                         uintptr // 展开
	Collapse                       uintptr // 折叠
	Get_CurrentExpandCollapseState uintptr // 获取当前展开/折叠状态
	Get_CachedExpandCollapseState  uintptr // 获取缓存展开/折叠状态
}

// Expand 展开控件
// 返回: 错误信息
func (v *IUIAutomationExpandCollapsePattern) Expand() error {
	ret, _, _ := syscall.SyscallN(
		v.vtbl.Expand,
		uintptr(unsafe.Pointer(v)),
	)
	if ret != 0 {
		return HResult(ret)
	}
	return nil
}

// Collapse 折叠控件
// 返回: 错误信息
func (v *IUIAutomationExpandCollapsePattern) Collapse() error {
	ret, _, _ := syscall.SyscallN(
		v.vtbl.Collapse,
		uintptr(unsafe.Pointer(v)),
	)
	if ret != 0 {
		return HResult(ret)
	}
	return nil
}

// Release 释放展开/折叠模式对象
// 返回: 引用计数
func (v *IUIAutomationExpandCollapsePattern) Release() uint32 {
	return (*IUnKnown)(unsafe.Pointer(v)).Release()
}
