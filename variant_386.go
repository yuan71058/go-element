//go:build 386

// Package uiautomation 提供 Windows UI Automation 的 Go 语言封装
package uiautomation

// VARIANT COM VARIANT 结构（32位版本）
// 用于存储各种类型的值，是 COM 中的通用数据类型
type VARIANT struct {
	VT         TagVarenum // 变体类型
	wReserved1 uint16     // 保留字段1
	wReserved2 uint16     // 保留字段2
	wReserved3 uint16     // 保留字段3
	Val        int64      // 值（根据 VT 类型解释）
}