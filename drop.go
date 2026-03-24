// Package uiautomation 提供 Windows UI Automation 的 Go 语言封装
package uiautomation

import "unsafe"

// IDropTarget 拖放目标接口
// 用于支持拖放操作的目标对象
type IDropTarget struct {
	vtbl *IUnKnown
}

// IDropTargetVtbl 拖放目标接口虚函数表
// TODO:: IDropTarget method - 需要实现具体方法
type IDropTargetVtbl struct {
	IUnKnownVtbl
	DragEnter uintptr // 拖拽进入
	DragLeave uintptr // 拖拽离开
	DragOver  uintptr // 拖拽经过
	Drop      uintptr // 放置
}

// newIDropTarget 内部函数：从 IUnKnown 创建拖放目标接口
func newIDropTarget(unk *IUnKnown) *IDropTarget {
	return (*IDropTarget)(unsafe.Pointer(unk))
}

// NewIDropTarget 创建拖放目标接口实例
// 参数: unk - IUnKnown 接口
// 返回: 拖放目标接口实例
func NewIDropTarget(unk *IUnKnown) *IDropTarget {
	return newIDropTarget(unk)
}

// DragEnter 处理拖拽进入事件
// 返回: 错误信息（当前为空实现）
func (v *IDropTarget) DragEnter() error {
	return nil
}

// DragLeave 处理拖拽离开事件
// 返回: 错误信息（当前为空实现）
func (v *IDropTarget) DragLeave() error {
	return nil
}

// DragOver 处理拖拽经过事件
// 返回: 错误信息（当前为空实现）
func (v *IDropTarget) DragOver() error {
	return nil
}

// Drop 处理放置事件
// 返回: 错误信息（当前为空实现）
func (v *IDropTarget) Drop() error {
	return nil
}