// Package uiautomation 提供 Windows UI Automation 的 Go 语言封装
package uiautomation

import "syscall"

// UI Automation COM 类和接口的 GUID
// 参考: https://learn.microsoft.com/zh-cn/previous-versions/windows/desktop/legacy/ff384838(v=vs.85)
var (
	// CLSID_CUIAutomation UI Automation 组件类标识符
	// 用于创建 IUIAutomation 实例
	CLSID_CUIAutomation = &syscall.GUID{0xff48dba4, 0x60ef, 0x4201, [8]byte{0xaa, 0x87, 0x54, 0x10, 0x3e, 0xef, 0x59, 0x4e}}

	// IID_IUIAutomation IUIAutomation 接口标识符
	// 注册表路径: HKEY_LOCAL_MACHINE\SOFTWARE\Classes\Interface\{IID}
	IID_IUIAutomation = &syscall.GUID{0x30cbe57d, 0xd9d0, 0x452a, [8]byte{0xab, 0x13, 0x7a, 0xc5, 0xac, 0x48, 0x25, 0xee}}
)
