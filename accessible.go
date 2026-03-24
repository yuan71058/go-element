// Package uiautomation 提供 Windows UI Automation 的 Go 语言封装
package uiautomation

import (
	"syscall"
	"unsafe"
)

// IAccessible 可访问性接口
// 用于支持辅助功能，如屏幕阅读器等
// 参考: https://learn.microsoft.com/zh-cn/windows/win32/api/oleacc/nn-oleacc-iaccessible
type IAccessible struct {
	vtbl *IDispatch
}

// IAccessibleVtbl 可访问性接口虚函数表
type IAccessibleVtbl struct {
	// https://learn.microsoft.com/zh-cn/windows/win32/api/oleacc/nn-oleacc-iaccessible
	IDispatchVtbl
	AccDoDefaultAction      uintptr // 执行默认操作
	AccHitTest              uintptr // 命中测试
	AccLocation             uintptr // 获取位置
	AccNavigate             uintptr // 导航（不支持）
	AccSelect               uintptr // 选择
	Get_accChild            uintptr // 获取子元素
	Get_accChildCount       uintptr // 获取子元素数量
	Get_accDefaultAction    uintptr // 获取默认操作
	Get_accDescription      uintptr // 获取描述
	Get_accFocus            uintptr // 获取焦点
	Get_accHelp             uintptr // 获取帮助
	Get_accHelpTopic        uintptr // 获取帮助主题（不支持）
	Get_accKeyboardShortcut uintptr // 获取键盘快捷键
	Get_accName             uintptr // 获取名称
	Get_accParent           uintptr // 获取父元素
	Get_accRole             uintptr // 获取角色
	Get_accSelection        uintptr // 获取选择
	Get_accState            uintptr // 获取状态
	Get_accValue            uintptr // 获取值
	Put_accName             uintptr // 设置名称（不支持）
	Put_accValue            uintptr // 设置值
}

// Release 释放可访问性接口对象
// 返回: 引用计数
func (v *IAccessible) Release() uint32 {
	return (*IUnKnown)(unsafe.Pointer(v)).Release()
}

// newIAccessible 内部函数：从 IDispatch 创建可访问性接口
func newIAccessible(unk *IDispatch) *IAccessible {
	return (*IAccessible)(unsafe.Pointer(unk))
}

// NewIAccessible 创建可访问性接口实例
// 参数: unk - IDispatch 接口
// 返回: 可访问性接口实例
func NewIAccessible(unk *IDispatch) *IAccessible {
	return newIAccessible(unk)
}

// AccDoDefaultAction 执行默认操作
// 参数: in - VARIANT 参数
// 返回: 错误信息
func (v *IAccessible) AccDoDefaultAction(in *VARIANT) error {
	ret, _, _ := syscall.SyscallN(
		(*IAccessibleVtbl)(unsafe.Pointer(v.vtbl)).AccDoDefaultAction,
		uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(in)),
	)
	if ret != 0 {
		return HResult(ret)
	}
	return nil
}