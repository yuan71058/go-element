// Package uiautomation 提供 Windows UI Automation 的 Go 语言封装
package uiautomation

import (
	"syscall"
	"time"
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

// DoubleClick 双击控件（调用两次 Invoke）
// 参数:
//   - interval: 两次点击之间的间隔时间（毫秒）
//
// 返回: 错误信息
func (v *IUIAutomationInvokePattern) DoubleClick(interval int) error {
	err := v.Invoke()
	if err != nil {
		return err
	}
	time.Sleep(time.Duration(interval) * time.Millisecond)
	err = v.Invoke()
	if err != nil {
		return err
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

// newIUIAutomationSelectionItemPattern 内部函数：从 IUnKnown 创建选择项模式
func newIUIAutomationSelectionItemPattern(unk *IUnKnown) *IUIAutomationSelectionItemPattern {
	return (*IUIAutomationSelectionItemPattern)(unsafe.Pointer(unk))
}

// NewIUIAutomationSelectionItemPattern 创建选择项模式实例
// 参数: unk - IUnKnown 接口
// 返回: 选择项模式实例
func NewIUIAutomationSelectionItemPattern(unk *IUnKnown) *IUIAutomationSelectionItemPattern {
	return newIUIAutomationSelectionItemPattern(unk)
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

// IUIAutomationLegacyIAccessiblePattern 旧版 IAccessible 模式接口
// 用于执行默认操作（如双击）
type IUIAutomationLegacyIAccessiblePattern struct {
	vtbl *IUIAutomationLegacyIAccessiblePatternVtbl
}

// IUIAutomationLegacyIAccessiblePatternVtbl 旧版 IAccessible 模式虚函数表
type IUIAutomationLegacyIAccessiblePatternVtbl struct {
	IUnKnownVtbl
	Get_CurrentChildId          uintptr // 获取子 ID
	Get_CurrentDefaultAction    uintptr // 获取默认操作
	Get_CurrentDescription      uintptr // 获取描述
	Get_CurrentHelp             uintptr // 获取帮助
	Get_CurrentKeyboardShortcut uintptr // 获取键盘快捷键
	Get_CurrentName             uintptr // 获取名称
	Get_CurrentRole             uintptr // 获取角色
	Get_CurrentState            uintptr // 获取状态
	Get_CurrentValue            uintptr // 获取值
	DoDefaultAction             uintptr // 执行默认操作
	SetCurrentValue             uintptr // 设置当前值
	Get_CurrentParent           uintptr // 获取父元素
	Get_Selection               uintptr // 获取选择
	Get_FocusedElement          uintptr // 获取焦点元素
}

// DoDefaultAction 执行默认操作（如双击）
// 返回: 错误信息
func (v *IUIAutomationLegacyIAccessiblePattern) DoDefaultAction() error {
	variant := NewVariant(VT_EMPTY, 0)
	ret, _, _ := syscall.SyscallN(
		v.vtbl.DoDefaultAction,
		uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(&variant)),
	)
	if ret != 0 {
		return HResult(ret)
	}
	return nil
}

// Get_CurrentDefaultAction 获取默认操作
// 返回: 默认操作和可能的错误
func (v *IUIAutomationLegacyIAccessiblePattern) Get_CurrentDefaultAction() (string, error) {
	var bstr uintptr
	ret, _, _ := syscall.SyscallN(
		v.vtbl.Get_CurrentDefaultAction,
		uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(&bstr)),
	)
	if ret != 0 {
		return "", HResult(ret)
	}
	defer procSysFreeString.Call(bstr)
	return bstr2str(bstr), nil
}

// Release 释放旧版 IAccessible 模式对象
// 返回: 引用计数
func (v *IUIAutomationLegacyIAccessiblePattern) Release() uint32 {
	return (*IUnKnown)(unsafe.Pointer(v)).Release()
}

// newIUIAutomationLegacyIAccessiblePattern 内部函数：从 IUnKnown 创建旧版 IAccessible 模式
func newIUIAutomationLegacyIAccessiblePattern(unk *IUnKnown) *IUIAutomationLegacyIAccessiblePattern {
	return (*IUIAutomationLegacyIAccessiblePattern)(unsafe.Pointer(unk))
}

// NewIUIAutomationLegacyIAccessiblePattern 创建旧版 IAccessible 模式实例
// 参数: unk - IUnKnown 接口
// 返回: 旧版 IAccessible 模式实例
func NewIUIAutomationLegacyIAccessiblePattern(unk *IUnKnown) *IUIAutomationLegacyIAccessiblePattern {
	return newIUIAutomationLegacyIAccessiblePattern(unk)
}
