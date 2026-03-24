## Windows UI Automation Go 语言库
[![Go Reference](https://pkg.go.dev/badge/github.com/auuunya/go-element.svg)](https://pkg.go.dev/github.com/auuunya/go-element)

### 概述
这个 Go 库为 Windows 平台提供了强大且高性能的 UI 自动化接口，基于 UI Automation 框架。它允许您通过 Windows COM 接口与应用程序的 UI 元素进行交互和操作。

### 核心增强功能
- **高性能优化**：使用 `IUIAutomationCacheRequest` 批量获取属性和模式，显著减少 IPC（进程间通信）开销
- **无反射设计**：移除了基于反射的属性填充，提升运行时效率和类型安全性
- **改进的模式支持**：提供常用自动化模式的即用型封装：`Value`（值）、`Invoke`（调用）、`Toggle`（切换）、`ExpandCollapse`（展开折叠）和 `SelectionItem`（选择项）
- **资源安全**：严格的 COM 对象生命周期管理，通过显式 `Release()` 调用防止内存泄漏和句柄耗尽
- **更好的搜索 API**：在 `Element` 结构中集成了便捷的 `FindByName` 和 `FindByAutomationId` 方法

### 安装
```shell
go get -u github.com/auuunya/go-element
```

### 编译说明
为确保在各种 Windows 架构（如 XP/7/10/11 的 32 位或 64 位版本）上正常运行，请参考以下构建命令：

- **64 位系统（默认）**：
  ```shell
  GOOS=windows GOARCH=amd64 go build .
  ```
- **32 位系统（高兼容性）**：
  ```shell
  GOOS=windows GOARCH=386 go build .
  ```
> **注意**：如果在运行 64 位编译的程序时收到 `is not valid win32 application` 错误，请尝试使用 `386` 架构重新构建。

### 快速开始：

#### 1. 输出 UI 树（浏览器示例）
```go
import uia "github.com/auuunya/go-element"

func main() {
	uia.CoInitialize()
	defer uia.CoUninitialize()
	
	hwnd, _ := uia.GetWindowForString("Chrome_WidgetWin_1", "")
	instance, _ := uia.CreateInstance(uia.CLSID_CUIAutomation, uia.IID_IUIAutomation, uia.CLSCTX_INPROC_SERVER)
	ppv := uia.NewIUIAutomation(uia.NewIUnKnown(instance))
	defer ppv.Release()

	root, _ := uia.ElementFromHandle(ppv, hwnd)
	defer root.Release()

	// 使用缓存请求进行高性能遍历
	elems := uia.TraverseUIElementTree(ppv, root)
	uia.TreeString(elems, 0)
}
```

#### 2. 自动化记事本（模式使用示例）
```go
func main() {
	uia.CoInitialize()
	defer uia.CoUninitialize()
	
	hwnd, _ := uia.GetWindowForString("Notepad", "")
	instance, _ := uia.CreateInstance(uia.CLSID_CUIAutomation, uia.IID_IUIAutomation, uia.CLSCTX_ALL)
	ppv := uia.NewIUIAutomation(uia.NewIUnKnown(instance))
	defer ppv.Release()

	root, _ := uia.ElementFromHandle(ppv, hwnd)
	tree := uia.TraverseUIElementTree(ppv, root)

	// 通过名称查找元素并通过模式设置值
	if editor := tree.FindByName("文本编辑器"); editor != nil {
		if vp, err := editor.GetValuePattern(); err == nil {
			defer vp.Release()
			vp.SetValue("Hello from go-element!")
		}
	}
}
```

### 功能状态
- [x] 高性能属性缓存
- [x] UI 结构的 JSON 序列化
- [x] 常用模式封装（Invoke、Value、Toggle 等）
- [ ] UI 事件监听器支持
- [ ] 更多专门的 UI 操作

### 贡献
欢迎贡献！如果您发现错误或想要增强库的功能，请随时提出 issue 或提交 pull request。

### 许可证
本库采用 MIT 许可证分发