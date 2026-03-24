<div align="center">

# 🪟 Windows UI Automation Go 语言库

[![Go Reference](https://pkg.go.dev/badge/github.com/yuan71058/go-element.svg)](https://pkg.go.dev/github.com/yuan71058/go-element)
[![Go Report Card](https://goreportcard.com/badge/github.com/yuan71058/go-element)](https://goreportcard.com/report/github.com/yuan71058/go-element)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/dl/)

**一个高性能、易用的 Windows UI 自动化 Go 库**

基于 Windows UI Automation 框架，提供简洁的 Go API 进行桌面应用自动化

[快速开始](#-快速开始) • [功能特性](#-功能特性) • [示例](#-示例) • [API文档](#-api-文档)

</div>

---

## 📖 概述

**go-element** 是一个专为 Windows 平台设计的 UI 自动化 Go 语言库。它通过 Windows COM 接口与 UI Automation 框架交互，让您能够轻松地：

- 🔍 **查找和遍历** UI 元素树
- 🖱️ **模拟用户操作** 如点击、输入文本、选择等
- 📊 **读取 UI 属性** 如名称、类名、控件类型等
- ⚡ **高性能缓存** 批量获取属性，减少 IPC 开销

## ✨ 功能特性

| 特性 | 描述 |
|------|------|
| 🚀 **高性能** | 使用 `IUIAutomationCacheRequest` 批量获取属性，显著减少 IPC 开销 |
| 🎯 **零反射** | 移除反射逻辑，显式属性填充，提升运行时效率和类型安全 |
| 🔧 **模式封装** | 内置常用模式：`Value`、`Invoke`、`Toggle`、`ExpandCollapse`、`SelectionItem` |
| 🛡️ **资源安全** | 严格的 COM 生命周期管理，防止内存泄漏和句柄耗尽 |
| 🔍 **便捷搜索** | `FindByName`、`FindByAutomationId` 等快捷方法 |
| 📦 **JSON 支持** | UI 结构可序列化为 JSON，便于调试和分析 |

## 📦 安装

```shell
go get -u github.com/yuan71058/go-element
```

## 🔨 编译说明

为确保在各种 Windows 架构上正常运行：

| 架构 | 命令 | 说明 |
|------|------|------|
| **64 位** | `GOOS=windows GOARCH=amd64 go build .` | 默认，适用于现代系统 |
| **32 位** | `GOOS=windows GOARCH=386 go build .` | 高兼容性，适用于旧系统 |

> ⚠️ **注意**：如果 64 位程序出现 `is not valid win32 application` 错误，请使用 32 位架构重新编译。

## 🚀 快速开始

### 基础用法

```go
package main

import (
    "fmt"
    uia "github.com/yuan71058/go-element"
)

func main() {
    // 1. 初始化 COM
    uia.CoInitialize()
    defer uia.CoUninitialize()

    // 2. 查找目标窗口
    hwnd, err := uia.GetWindowForString("Notepad", "")
    if err != nil {
        panic("窗口未找到")
    }

    // 3. 创建 UI Automation 实例
    instance, _ := uia.CreateInstance(
        uia.CLSID_CUIAutomation,
        uia.IID_IUIAutomation,
        uia.CLSCTX_ALL,
    )
    ppv := uia.NewIUIAutomation(uia.NewIUnKnown(instance))
    defer ppv.Release()

    // 4. 获取窗口元素
    root, _ := ppv.ElementFromHandle(hwnd)
    defer root.Release()

    // 5. 遍历 UI 树
    tree := uia.TraverseUIElementTree(ppv, root)
    uia.TreeString(tree, 0)
}
```

## 📚 示例

### 1️⃣ 输出 UI 树结构

```go
// 遍历 Chrome 浏览器 UI 树
hwnd, _ := uia.GetWindowForString("Chrome_WidgetWin_1", "")
instance, _ := uia.CreateInstance(uia.CLSID_CUIAutomation, uia.IID_IUIAutomation, uia.CLSCTX_INPROC_SERVER)
ppv := uia.NewIUIAutomation(uia.NewIUnKnown(instance))
defer ppv.Release()

root, _ := ppv.ElementFromHandle(hwnd)
defer root.Release()

elems := uia.TraverseUIElementTree(ppv, root)
uia.TreeString(elems, 0)  // 打印树形结构
```

### 2️⃣ 自动化记事本

```go
// 查找记事本并输入文本
hwnd, _ := uia.GetWindowForString("Notepad", "")
// ... 初始化代码 ...

tree := uia.TraverseUIElementTree(ppv, root)

// 查找文本编辑器并设置值
if editor := tree.FindByName("文本编辑器"); editor != nil {
    if vp, err := editor.GetValuePattern(); err == nil {
        defer vp.Release()
        vp.SetValue("Hello from go-element!")
    }
}
```

### 3️⃣ 按钮点击

```go
// 查找并点击按钮
if button := tree.FindByName("确定"); button != nil {
    if ip, err := button.GetInvokePattern(); err == nil {
        defer ip.Release()
        ip.Invoke()  // 执行点击
    }
}
```

### 4️⃣ 复选框切换

```go
// 查找复选框并切换状态
if checkbox := tree.FindByName("记住密码"); checkbox != nil {
    if tp, err := checkbox.GetTogglePattern(); err == nil {
        defer tp.Release()
        state, _ := tp.Get_CurrentToggleState()
        fmt.Printf("当前状态: %v\n", state)
        tp.Toggle()  // 切换状态
    }
}
```

### 5️⃣ 自定义条件搜索

```go
// 查找所有启用的按钮
buttons := uia.FindElems(tree, func(elem *uia.Element) bool {
    return elem.CurrentControlType == uia.UIA_ButtonControlTypeId && 
           elem.CurrentIsEnabled == 1
})

for i, btn := range buttons {
    fmt.Printf("%d. %s\n", i+1, btn.CurrentName)
}
```

### 6️⃣ 微信自动化

```go
// 查找微信窗口
hwnd, _ := uia.GetWindowForString("WeChatMainWndForPC", "")
// ... 初始化代码 ...

tree := uia.TraverseUIElementTree(ppv, root)

// 在搜索框输入联系人
if searchBox := tree.FindByName("搜索"); searchBox != nil {
    if vp, err := searchBox.GetValuePattern(); err == nil {
        defer vp.Release()
        vp.SetValue("文件传输助手")
    }
}

// 查找并输入消息
if inputBox := tree.FindByAutomationId("chat_input_field"); inputBox != nil {
    if vp, err := inputBox.GetValuePattern(); err == nil {
        defer vp.Release()
        vp.SetValue("自动发送的消息")
    }
}
```

## 📁 项目结构

```
go-element/
├── 📄 com.go              # COM 初始化和窗口查找
├── 📄 uiautomation.go     # IUIAutomation 核心接口
├── 📄 element.go          # Element 结构和便捷方法
├── 📄 client.go           # Pattern 客户端接口
├── 📄 condition.go        # 条件查询接口
├── 📄 constants.go        # 常量定义
├── 📄 id.go               # 属性ID、模式ID、控件类型ID
├── 📄 enum.go             # 枚举类型
├── 📄 typedef.go          # 类型定义
├── 📄 variant.go          # VARIANT 类型
├── 📄 unknown.go          # IUnknown 基础接口
├── 📄 uia.go              # UIA 相关定义
├── 📂 example/            # 示例代码
│   ├── wechat_demo/       # 微信自动化示例
│   ├── notepad_demo/      # 记事本自动化示例
│   └── ...                # 更多示例
├── 📄 README.md           # 项目说明
└── 📄 API文档.md          # 详细 API 文档
```

## 📋 支持的模式 (Patterns)

| 模式 | 用途 | 主要方法 |
|------|------|----------|
| `ValuePattern` | 文本输入/获取 | `SetValue()`, `Get_CurrentValue()` |
| `InvokePattern` | 按钮点击 | `Invoke()` |
| `TogglePattern` | 复选框切换 | `Toggle()`, `Get_CurrentToggleState()` |
| `ExpandCollapsePattern` | 菜单展开/折叠 | `Expand()`, `Collapse()` |
| `SelectionItemPattern` | 列表项选择 | `Select()` |

## 📊 功能状态

- [x] 高性能属性缓存
- [x] UI 结构的 JSON 序列化
- [x] 常用模式封装（Invoke、Value、Toggle 等）
- [x] 便捷元素查找方法
- [x] 详细中文注释和文档
- [ ] UI 事件监听器支持
- [ ] 更多专门的 UI 操作

## ⚠️ 注意事项

### COM 初始化
- ✅ 必须在使用前调用 `CoInitialize()`
- ✅ 使用后必须调用 `CoUninitialize()`
- ✅ 建议使用 `defer` 确保清理

### 资源管理
- ✅ 所有 COM 对象使用完毕后必须调用 `Release()`
- ✅ 使用 `defer` 确保资源释放

### 架构兼容
- ⚠️ 64 位程序无法操作 32 位进程的 UI
- ⚠️ 32 位程序无法操作 64 位进程的 UI

## 📖 API 文档

详细的 API 文档请参阅 [API文档.md](./API文档.md)

## 🤝 贡献

欢迎贡献！如果您发现错误或想要增强库的功能，请：

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 提交 Pull Request

## 📄 许可证

本库采用 MIT 许可证分发。详见 [LICENSE](LICENSE) 文件。

## 🔗 参考资源

- [Microsoft UI Automation 官方文档](https://learn.microsoft.com/zh-cn/windows/win32/winauto/entry-uiauto-win32)
- [UI Automation Property IDs](https://learn.microsoft.com/zh-cn/windows/win32/winauto/uiauto-entry-propids)
- [UI Automation Control Pattern IDs](https://learn.microsoft.com/zh-cn/windows/win32/winauto/uiauto-controlpattern-ids)

---

<div align="center">

**如果这个项目对您有帮助，请给一个 ⭐️ Star！**

Made with ❤️ by go-element contributors

</div>
