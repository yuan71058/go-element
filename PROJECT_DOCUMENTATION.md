# go-element 项目详细文档

> **文档版本**: v1.0  
> **生成时间**: 2026-03-27  
> **项目名称**: go-element  
> **项目类型**: Windows UI 自动化 Go 语言库  
> **许可证**: MIT

---

## 目录

1. [项目概述](#1-项目概述)
2. [项目架构](#2-项目架构)
3. [核心模块详解](#3-核心模块详解)
4. [API 参考](#4-api-参考)
5. [使用示例](#5-使用示例)
6. [微信自动化模块](#6-微信自动化模块)
7. [开发指南](#7-开发指南)
8. [最佳实践](#8-最佳实践)
9. [常见问题](#9-常见问题)
10. [附录](#10-附录)

---

## 1. 项目概述

### 1.1 项目简介

**go-element** 是一个高性能、易用的 Windows UI 自动化 Go 语言库，基于 Windows UI Automation 框架实现。它通过 Windows COM 接口与 UI Automation 框架交互，为开发者提供了简洁的 Go API 来进行桌面应用自动化操作。

### 1.2 核心特性

| 特性 | 描述 |
|------|------|
| **高性能缓存** | 使用 `IUIAutomationCacheRequest` 批量获取属性，显著减少 IPC 开销 |
| **零反射设计** | 移除反射逻辑，显式属性填充，提升运行时效率和类型安全 |
| **模式封装** | 内置常用模式：`Value`、`Invoke`、`Toggle`、`ExpandCollapse`、`SelectionItem` |
| **资源安全** | 严格的 COM 生命周期管理，防止内存泄漏和句柄耗尽 |
| **便捷搜索** | `FindByName`、`FindByAutomationId` 等快捷方法 |
| **JSON 支持** | UI 结构可序列化为 JSON，便于调试和分析 |

### 1.3 技术栈

- **编程语言**: Go 1.19+
- **系统依赖**: golang.org/x/sys v0.20.0
- **目标平台**: Windows (仅限 Windows)
- **架构支持**: 32位 (386) 和 64位 (amd64)

### 1.4 项目结构

```
go-element/
├── 核心库文件
│   ├── com.go              # COM 初始化和窗口查找
│   ├── uiautomation.go     # IUIAutomation 核心接口
│   ├── element.go          # Element 结构体和便捷方法
│   ├── client.go           # Pattern 客户端接口
│   ├── condition.go        # 条件查询接口
│   ├── constants.go        # 常量定义
│   ├── id.go               # 属性ID、模式ID、控件类型ID
│   ├── enum.go             # 枚举类型
│   ├── typedef.go          # 类型定义
│   ├── variant.go          # VARIANT 类型
│   ├── unknown.go          # IUnknown 基础接口
│   ├── dispatch.go         # IDispatch 接口
│   ├── accessible.go       # IAccessible 接口
│   ├── provider.go         # Provider 接口
│   ├── textserv.go         # 文本服务接口
│   ├── drop.go             # 拖放接口
│   ├── uia.go              # UIA 相关定义
│   ├── variant_amd64.go    # 64位变体实现
│   └── variant_386.go       # 32位变体实现
├── wechat/                 # 微信自动化模块
│   ├── wechat.go          # 微信机器人实现
│   └── types.go           # 微信相关类型定义
├── example/                # 示例代码
│   ├── notepad_demo/      # 记事本自动化示例
│   ├── wechat_demo/       # 微信自动化示例
│   ├── tree_demo/         # UI树遍历示例
│   ├── calc_demo/         # 计算器自动化示例
│   ├── find_demo/         # 元素查找示例
│   ├── json_demo/         # JSON导出示例
│   └── software/          # 软件信息示例
├── README.md              # 英文项目说明
├── README.zh.md           # 中文项目说明
├── API文档.md             # 详细 API 文档
└── go.mod                # Go 模块定义
```

---

## 2. 项目架构

### 2.1 架构层次图

```
┌─────────────────────────────────────────────────────────────┐
│                        用户代码层                              │
│           (业务逻辑、自动化脚本、测试代码等)                      │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                       便捷封装层                              │
│  ┌────────────┐ ┌────────────┐ ┌────────────┐ ┌───────────┐ │
│  │  Element   │ │ SearchFunc │ │Helper API  │ │ 微信Bot   │ │
│  │  高级封装  │ │  条件查询  │ │  工具函数  │ │  微信封装 │ │
│  └────────────┘ └────────────┘ └────────────┘ └───────────┘ │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                       模式层 (Patterns)                       │
│  ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐ ┌────────────┐ │
│  │ Value  │ │ Invoke │ │ Toggle │ │Expand  │ │Selection  │ │
│  │Pattern │ │Pattern │ │Pattern │ │Collapse│ │   Item    │ │
│  └────────┘ └────────┘ └────────┘ └────────┘ └────────────┘ │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    核心接口层 (IUIAutomation)                  │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │IUIAutomation │  │   Element    │  │  Condition   │     │
│  │    核心接口   │  │   元素操作   │  │   条件构建   │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │ TreeWalker  │  │CacheRequest  │  │   Variant    │     │
│  │   树遍历    │  │   缓存请求   │  │   变体类型   │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    基础层 (Foundation)                       │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │   COM 接口   │  │   Windows    │  │   内存管理   │     │
│  │   封装       │  │   API 调用   │  │   资源释放   │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│               Windows UI Automation API                      │
│                    (Windows 系统原生接口)                     │
└─────────────────────────────────────────────────────────────┘
```

### 2.2 模块依赖关系

```
com.go (COM 初始化)
    │
    └─► uiautomation.go (IUIAutomation 核心)
            │
            ├─► element.go (Element 封装)
            │       │
            │       └─► client.go (Pattern 客户端)
            │
            ├─► condition.go (条件查询)
            │
            ├─► provider.go (Provider 接口)
            │
            └─► accessible.go (IAccessible 接口)

unknown.go (IUnknown 基础)
    │
    ├─► dispatch.go (IDispatch)
    │
    └─► variant.go (VARIANT 类型)

id.go (常量定义)
    │
    └─► constants.go (枚举常量)
            │
            └─► enum.go (枚举类型)

typedef.go (类型定义)
    │
    └─► uia.go (UIA 相关)
```

---

## 3. 核心模块详解

### 3.1 COM 初始化模块 (com.go)

**文件路径**: `E:\SRC\go-element-main\go-element-main\com.go`

**功能描述**: 
- 负责 COM 库的初始化和清理
- 提供窗口查找功能

**核心函数**:

| 函数名 | 功能 | 参数 | 返回值 |
|--------|------|------|--------|
| `CoInitialize()` | 初始化 COM 库 | 无 | error |
| `CoUninitialize()` | 清理 COM 库 | 无 | 无 |
| `GetWindowForString()` | 根据类名/标题查找窗口 | classname, windowname | hwnd, error |
| `CreateInstance()` | 创建 COM 对象实例 | clsid, riid, clsctx | unsafe.Pointer, error |
| `FindWindowW()` | 查找顶层窗口 | lpclass, lpwindow | hwnd |
| `FindWindowExW()` | 查找子窗口 | phwdn, chwdn, lpclass, lpwindow | hwnd |

**使用示例**:
```go
// 初始化 COM
err := uia.CoInitialize()
if err != nil {
    log.Fatal(err)
}
defer uia.CoUninitialize()

// 查找窗口
hwnd, err := uia.GetWindowForString("Notepad", "")
if err != nil {
    log.Fatal("窗口未找到")
}
```

### 3.2 UI Automation 核心模块 (uiautomation.go)

**文件路径**: `E:\SRC\go-element-main\go-element-main\uiautomation.go`

**功能描述**:
- IUIAutomation 核心接口封装
- IUIAutomationElement 元素接口封装
- 事件处理器配置
- 条件创建方法

**核心结构**:

```go
// IUIAutomation 核心接口
type IUIAutomation struct {
    vtbl *IUnKnown
}

// IUIAutomationElement UI 元素接口
type IUIAutomationElement struct {
    vtbl *IUnKnown
}
```

**核心方法**:

| 方法 | 功能 |
|------|------|
| `ElementFromHandle()` | 从窗口句柄获取 UI 元素 |
| `GetRootElement()` | 获取桌面根元素 |
| `ElementFromPoint()` | 从屏幕坐标获取元素 |
| `GetFocusedElement()` | 获取当前焦点元素 |
| `CreateCacheRequest()` | 创建缓存请求对象 |
| `CreateTreeWalker()` | 创建树遍历器 |
| `CreateTrueCondition()` | 创建真条件 |
| `CreatePropertyCondition()` | 创建属性条件 |
| `CreateAndCondition()` | 创建 AND 条件 |
| `CreateOrCondition()` | 创建 OR 条件 |
| `CreateNotCondition()` | 创建 NOT 条件 |

### 3.3 Element 元素封装模块 (element.go)

**文件路径**: `E:\SRC\go-element-main\go-element-main\element.go`

**功能描述**:
- Element 结构体封装
- 属性自动填充
- 元素查找功能
- 模式获取方法

**Element 结构体**:

```go
type Element struct {
    UIAutoElement               *IUIAutomationElement  // 底层元素接口
    CurrentName                 string                 // 名称
    CurrentClassName            string                 // 类名
    CurrentControlType          ControlTypeId          // 控件类型
    CurrentAutomationId         string                 // 自动化ID
    CurrentIsEnabled            int32                  // 是否启用
    CurrentProcessId            int32                  // 进程ID
    CurrentBoundingRectangle    *TagRect               // 边界矩形
    CurrentLocalizedControlType string                 // 本地化控件类型
    SupportedPatterns           []PatternId            // 支持的模式
    Child                       []*Element             // 子元素
}
```

**核心方法**:

| 方法 | 功能 |
|------|------|
| `TraverseUIElementTree()` | 遍历 UI 元素树（带缓存） |
| `TreeString()` | 以树形结构打印元素 |
| `FindByName()` | 按名称查找元素 |
| `FindByAutomationId()` | 按 AutomationId 查找元素 |
| `GetValuePattern()` | 获取值模式 |
| `GetInvokePattern()` | 获取调用模式 |
| `GetTogglePattern()` | 获取切换模式 |
| `GetExpandCollapsePattern()` | 获取展开/折叠模式 |
| `GetSelectionItemPattern()` | 获取选择项模式 |
| `SearchElem()` | 自定义条件搜索（返回第一个匹配） |
| `FindElems()` | 自定义条件搜索（返回所有匹配） |
| `Populate()` | 填充元素属性 |

### 3.4 Pattern 客户端模块 (client.go)

**文件路径**: `E:\SRC\go-element-main\go-element-main\client.go`

**功能描述**:
- ValuePattern: 值模式，用于文本输入
- InvokePattern: 调用模式，用于按钮点击
- TogglePattern: 切换模式，用于复选框
- ExpandCollapsePattern: 展开/折叠模式，用于菜单
- SelectionItemPattern: 选择项模式，用于列表项
- LegacyIAccessiblePattern: 旧版可访问性模式

**支持的 Patterns**:

| Pattern | 功能 | 主要方法 |
|---------|------|----------|
| ValuePattern | 文本输入/获取 | `SetValue()`, `Get_CurrentValue()` |
| InvokePattern | 按钮点击 | `Invoke()`, `DoubleClick()` |
| TogglePattern | 复选框切换 | `Toggle()`, `Get_CurrentToggleState()` |
| ExpandCollapsePattern | 菜单展开/折叠 | `Expand()`, `Collapse()` |
| SelectionItemPattern | 列表项选择 | `Select()` |
| LegacyIAccessiblePattern | 默认操作 | `DoDefaultAction()` |

### 3.5 微信自动化模块 (wechat/)

**文件路径**: `E:\SRC\go-element-main\go-element-main\wechat\`

**功能描述**:
- 微信窗口自动化操作
- 联系人搜索和选择
- 消息读取和发送
- 会话列表管理
- 自动回复功能

**核心结构**:

```go
// Bot 微信机器人结构体
type Bot struct {
    ppv            *uia.IUIAutomation        // UI自动化接口
    hwnd           uintptr                   // 微信窗口句柄
    root           *uia.IUIAutomationElement // 微信窗口根元素
    logFile        *os.File                  // 日志文件
    sessionManager *SessionManager           // 会话管理器
    wechatId       string                    // 微信号
}
```

**核心方法**:

| 方法 | 功能 |
|------|------|
| `NewBot()` | 创建机器人实例 |
| `Close()` | 关闭机器人 |
| `FindElementByName()` | 按名称查找元素 |
| `FindElementByAutomationId()` | 按 AutomationId 查找 |
| `ClickElement()` | 点击元素（鼠标模拟） |
| `ClickElementWithInvoke()` | 点击元素（Invoke 模式） |
| `SetTextValue()` | 设置文本值 |
| `GetTextValue()` | 获取文本值 |
| `TypeText()` | 模拟键盘输入 |
| `PasteText()` | 剪贴板粘贴 |
| `SearchContact()` | 搜索并选择联系人 |
| `ReadChatMessages()` | 读取聊天消息 |
| `SendMessage()` | 发送消息 |
| `ReadSessionList()` | 读取会话列表 |
| `StartMessageMonitor()` | 启动消息监控 |
| `StartAutoReply()` | 启动自动回复 |

---

## 4. API 参考

### 4.1 常量定义

#### 属性 ID (PropertyId)

```go
const (
    UIA_RuntimeIdPropertyId            PropertyId = 30000
    UIA_BoundingRectanglePropertyId    PropertyId = 30001
    UIA_ProcessIdPropertyId            PropertyId = 30002
    UIA_ControlTypePropertyId          PropertyId = 30003
    UIA_LocalizedControlTypePropertyId PropertyId = 30004
    UIA_NamePropertyId                 PropertyId = 30005
    UIA_AcceleratorKeyPropertyId       PropertyId = 30006
    UIA_AccessKeyPropertyId            PropertyId = 30007
    UIA_HasKeyboardFocusPropertyId     PropertyId = 30008
    UIA_IsKeyboardFocusablePropertyId  PropertyId = 30009
    UIA_IsEnabledPropertyId            PropertyId = 30010
    UIA_AutomationIdPropertyId         PropertyId = 30011
    UIA_ClassNamePropertyId            PropertyId = 30012
    // ... 更多见 id.go
)
```

#### 模式 ID (PatternId)

```go
const (
    UIA_InvokePatternId            PatternId = 10000
    UIA_SelectionPatternId         PatternId = 10001
    UIA_ValuePatternId             PatternId = 10002
    UIA_RangeValuePatternId        PatternId = 10003
    UIA_ScrollPatternId            PatternId = 10004
    UIA_ExpandCollapsePatternId    PatternId = 10005
    UIA_GridPatternId              PatternId = 10006
    UIA_GridItemPatternId          PatternId = 10007
    UIA_WindowPatternId            PatternId = 10009
    UIA_SelectionItemPatternId     PatternId = 10010
    UIA_TogglePatternId            PatternId = 10015
    // ... 更多见 id.go
)
```

#### 控件类型 ID (ControlTypeId)

```go
const (
    UIA_ButtonControlTypeId       ControlTypeId = 50000
    UIA_CalendarControlTypeId     ControlTypeId = 50001
    UIA_CheckBoxControlTypeId     ControlTypeId = 50002
    UIA_ComboBoxControlTypeId     ControlTypeId = 50003
    UIA_EditControlTypeId         ControlTypeId = 50004
    UIA_HyperlinkControlTypeId    ControlTypeId = 50005
    UIA_ImageControlTypeId         ControlTypeId = 50006
    UIA_ListItemControlTypeId     ControlTypeId = 50007
    UIA_ListControlTypeId         ControlTypeId = 50008
    UIA_MenuControlTypeId         ControlTypeId = 50009
    UIA_MenuBarControlTypeId      ControlTypeId = 50010
    UIA_MenuItemControlTypeId     ControlTypeId = 50011
    UIA_ProgressBarControlTypeId  ControlTypeId = 50012
    UIA_RadioButtonControlTypeId  ControlTypeId = 50013
    UIA_ScrollBarControlTypeId    ControlTypeId = 50014
    UIA_SliderControlTypeId       ControlTypeId = 50015
    UIA_SpinnerControlTypeId      ControlTypeId = 50016
    UIA_StatusBarControlTypeId    ControlTypeId = 50017
    UIA_TabControlTypeId          ControlTypeId = 50018
    UIA_TabItemControlTypeId      ControlTypeId = 50019
    UIA_TextControlTypeId         ControlTypeId = 50020
    UIA_ToolBarControlTypeId      ControlTypeId = 50021
    UIA_ToolTipControlTypeId      ControlTypeId = 50022
    UIA_TreeControlTypeId         ControlTypeId = 50023
    UIA_TreeItemControlTypeId     ControlTypeId = 50024
    UIA_WindowControlTypeId       ControlTypeId = 50032
    // ... 更多见 id.go
)
```

### 4.2 枚举类型

#### ToggleState 切换状态

```go
type ToggleState int32

const (
    ToggleState_Off          ToggleState = 0  // 关闭
    ToggleState_On           ToggleState = 1  // 开启
    ToggleState_Indeterminate ToggleState = 2 // 不确定
)
```

#### OrientationType 方向类型

```go
type OrientationType int32

const (
    OrientationType_None       OrientationType = 0 // 无方向
    OrientationType_Horizontal OrientationType = 1 // 水平
    OrientationType_Vertical   OrientationType = 2 // 垂直
)
```

#### WindowVisualState 窗口可视状态

```go
type WindowVisualState int32

const (
    WindowVisualState_Normal   WindowVisualState = 0 // 正常
    WindowVisualState_Maximized WindowVisualState = 1 // 最大化
    WindowVisualState_Minimized WindowVisualState = 2 // 最小化
)
```

### 4.3 数据结构

#### TagRect 矩形区域

```go
type TagRect struct {
    Left   int32  // 左上角 X 坐标
    Top    int32  // 左上角 Y 坐标
    Right  int32  // 右下角 X 坐标
    Bottom int32  // 右下角 Y 坐标
}
```

#### TagPoint 点坐标

```go
type TagPoint struct {
    X int32  // X 坐标
    Y int32  // Y 坐标
}
```

#### VARIANT 变体类型

```go
type VARIANT struct {
    VT  TagVarenum  // 类型标记
    Val int64      // 值
}
```

#### TreeScope 树范围

```go
var (
    TreeScope_Element     TreeScope = 0x1  // 元素本身
    TreeScope_Children    TreeScope = 0x2  // 直接子元素
    TreeScope_Descendants TreeScope = 0x4  // 所有后代元素
    TreeScope_Parent      TreeScope = 0x8  // 父元素
    TreeScope_Ancestors   TreeScope = 0x10 // 所有祖先元素
    TreeScope_Subtree     TreeScope = TreeScope_Element | TreeScope_Children | TreeScope_Descendants
)
```

---

## 5. 使用示例

### 5.1 基础示例：输出 UI 树结构

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

    // 2. 查找窗口
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

### 5.2 自动化记事本

```go
package main

import (
    "fmt"
    uia "github.com/yuan71058/go-element"
)

func main() {
    uia.CoInitialize()
    defer uia.CoUninitialize()

    // 查找记事本
    hwnd, _ := uia.GetWindowForString("Notepad", "")
    instance, _ := uia.CreateInstance(
        uia.CLSID_CUIAutomation,
        uia.IID_IUIAutomation,
        uia.CLSCTX_ALL,
    )
    ppv := uia.NewIUIAutomation(uia.NewIUnKnown(instance))
    defer ppv.Release()

    root, _ := ppv.ElementFromHandle(hwnd)
    tree := uia.TraverseUIElementTree(ppv, root)

    // 查找文本编辑器并设置值
    if editor := tree.FindByName("文本编辑器"); editor != nil {
        if vp, err := editor.GetValuePattern(); err == nil {
            defer vp.Release()
            vp.SetValue("Hello from go-element!")
            value, _ := vp.Get_CurrentValue()
            fmt.Println("当前值:", value)
        }
    }
}
```

### 5.3 按钮点击

```go
package main

import (
    "fmt"
    uia "github.com/yuan71058/go-element"
)

func main() {
    // ... 初始化代码 ...

    tree := uia.TraverseUIElementTree(ppv, root)

    // 查找并点击按钮
    if button := tree.FindByName("确定"); button != nil {
        if ip, err := button.GetInvokePattern(); err == nil {
            defer ip.Release()
            ip.Invoke()
            fmt.Println("按钮已点击")
        }
    }
}
```

### 5.4 复选框切换

```go
package main

import (
    "fmt"
    uia "github.com/yuan71058/go-element"
)

func main() {
    // ... 初始化代码 ...

    tree := uia.TraverseUIElementTree(ppv, root)

    // 查找复选框并切换
    if checkbox := tree.FindByName("记住密码"); checkbox != nil {
        if tp, err := checkbox.GetTogglePattern(); err == nil {
            defer tp.Release()
            state, _ := tp.Get_CurrentToggleState()
            fmt.Printf("当前状态: %v\n", state)
            tp.Toggle()
        }
    }
}
```

### 5.5 自定义条件搜索

```go
package main

import (
    "fmt"
    uia "github.com/yuan71058/go-element"
)

func main() {
    // ... 初始化代码 ...

    tree := uia.TraverseUIElementTree(ppv, root)

    // 查找所有启用的按钮
    buttons := uia.FindElems(tree, func(elem *uia.Element) bool {
        return elem.CurrentControlType == uia.UIA_ButtonControlTypeId && 
               elem.CurrentIsEnabled == 1
    })

    fmt.Printf("找到 %d 个启用的按钮:\n", len(buttons))
    for i, btn := range buttons {
        fmt.Printf("%d. %s\n", i+1, btn.CurrentName)
    }
}
```

### 5.6 菜单展开/折叠

```go
package main

import (
    uia "github.com/yuan71058/go-element"
)

func main() {
    // ... 初始化代码 ...

    tree := uia.TraverseUIElementTree(ppv, root)

    // 查找菜单项并展开
    if menu := tree.FindByName("文件"); menu != nil {
        if ecp, err := menu.GetExpandCollapsePattern(); err == nil {
            defer ecp.Release()
            ecp.Expand()   // 展开
            // ... 操作 ...
            ecp.Collapse() // 折叠
        }
    }
}
```

---

## 6. 微信自动化模块

### 6.1 基本使用

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/yuan71058/go-element/wechat"
)

func main() {
    // 创建机器人
    bot, err := wechat.NewBot()
    if err != nil {
        log.Fatal(err)
    }
    defer bot.Close()

    // 搜索联系人
    err = bot.SearchContact("文件传输助手")
    if err != nil {
        log.Fatal(err)
    }

    // 发送消息
    err = bot.SendMessage("你好，来自 go-element！")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("消息发送成功")
}
```

### 6.2 读取聊天记录

```go
func readChatHistory() {
    bot, _ := wechat.NewBot()
    defer bot.Close()

    // 打开联系人
    bot.SearchContact("张三")

    // 读取消息
    messages, err := bot.ReadChatMessages()
    if err != nil {
        log.Fatal(err)
    }

    for _, msg := range messages {
        sender := "对方"
        if msg.IsSelf {
            sender = "自己"
        }
        fmt.Printf("[%s] %s: %s\n", msg.Time, sender, msg.Content)
    }
}
```

### 6.3 自动回复功能

```go
func autoReplyDemo() {
    bot, _ := wechat.NewBot()
    defer bot.Close()

    // 配置自动回复
    config := &wechat.AutoReplyConfig{
        WechatId: "your_wechat_id",
        Contacts: []string{"文件传输助手", "张三", "李四"},
        ReplyGenerator: func(content string) string {
            return "自动回复: " + content
        },
        OnMessage: func(sender, content string) {
            log.Printf("收到消息 - %s: %s", sender, content)
        },
        OnReply: func(sender, reply string) {
            log.Printf("发送回复 - %s: %s", sender, reply)
        },
        OnError: func(err error) {
            log.Printf("错误: %v", err)
        },
    }

    // 启动自动回复
    bot.StartAutoReply(config)
}
```

### 6.4 消息监控

```go
func monitorMessages() {
    bot, _ := wechat.NewBot()
    defer bot.Close()

    // 启动消息监控（每3秒检查一次）
    bot.StartMessageMonitor(3*time.Second, func(msg *wechat.NewMessage) {
        log.Printf("新消息 - %s: %s", msg.Sender, msg.Content)
    })
}
```

---

## 7. 开发指南

### 7.1 环境配置

#### 安装 Go

1. 下载并安装 Go 1.19+: https://go.dev/dl/
2. 验证安装: `go version`

#### 克隆项目

```bash
git clone https://github.com/yuan71058/go-element.git
cd go-element
```

#### 安装依赖

```bash
go mod download
```

### 7.2 编译说明

#### 64 位版本（默认）

```bash
GOOS=windows GOARCH=amd64 go build .
```

#### 32 位版本

```bash
GOOS=windows GOARCH=386 go build .
```

#### 运行示例

```bash
cd example/notepad_demo
go run main.go
```

### 7.3 资源管理

**重要**: 所有 COM 对象使用完毕后必须调用 `Release()`。

```go
// ✅ 正确：使用 defer 确保释放
element, _ := ppv.ElementFromHandle(hwnd)
defer element.Release()

// ✅ 正确：嵌套资源管理
vp, err := editor.GetValuePattern()
if err != nil {
    return err
}
defer vp.Release()

// ❌ 错误：忘记释放
element, _ := ppv.ElementFromHandle(hwnd)
// 忘记调用 element.Release()
```

### 7.4 错误处理

```go
// 检查错误
element, err := ppv.ElementFromHandle(hwnd)
if err != nil {
    return fmt.Errorf("获取元素失败: %w", err)
}
defer element.Release()

// 使用错误常量
if errors.Is(err, uia.ErrorNotFoundWindow) {
    fmt.Println("窗口未找到")
}
```

---

## 8. 最佳实践

### 8.1 性能优化

#### 使用缓存请求

```go
// ✅ 性能优化：使用缓存请求
cacheRequest, _ := ppv.CreateCacheRequest()
cacheRequest.AddProperty(UIA_NamePropertyId)
cacheRequest.AddProperty(UIA_ClassNamePropertyId)
cacheRequest.AddPattern(UIA_ValuePatternId)
elementArr, _ := root.FindAllBuildCache(TreeScope_Children, condition, cacheRequest)
```

#### 缩小搜索范围

```go
// ✅ 缩小范围提高性能
elementArr, _ := root.FindAll(TreeScope_Children, condition)

// ❌ 避免搜索整个子树
elementArr, _ := root.FindAll(TreeScope_Subtree, condition)
```

### 8.2 可靠性

#### 添加等待时间

```go
// UI 更新需要时间
time.Sleep(500 * time.Millisecond)

// 或使用重试机制
for i := 0; i < 3; i++ {
    if element := findElement(); element != nil {
        return element
    }
    time.Sleep(500 * time.Millisecond)
}
```

#### 检查元素状态

```go
func safeClick(element *uia.IUIAutomationElement) error {
    // 检查元素是否启用
    if element.Get_CurrentIsEnabled() == 0 {
        return fmt.Errorf("元素未启用")
    }
    
    // 检查元素是否在屏幕上
    if element.Get_CurrentIsOffscreen() != 0 {
        return fmt.Errorf("元素不在屏幕上")
    }
    
    // 执行点击
    // ...
}
```

### 8.3 调试技巧

#### 输出 UI 树结构

```go
tree := uia.TraverseUIElementTree(ppv, root)
uia.TreeString(tree, 0)
```

#### 格式化元素属性

```go
fmt.Println(editor.FormatString())
```

---

## 9. 常见问题

### Q1: "is not valid win32 application" 错误

**原因**: 架构不匹配（64位程序无法操作32位进程，或反之）

**解决**: 使用正确的编译目标
```bash
# 32 位系统
set GOARCH=386
go build .

# 64 位系统
set GOARCH=amd64
go build .
```

### Q2: 找不到窗口

**可能原因**:
1. 窗口类名不正确 - 使用 Spy++ 工具查看
2. 窗口未创建完成 - 添加等待时间
3. 权限不足 - 以管理员身份运行

**解决**:
```go
// 添加等待时间
time.Sleep(1 * time.Second)
hwnd, _ := uia.GetWindowForString("Notepad", "")
```

### Q3: 元素操作失败

**可能原因**:
1. 元素不支持该 Pattern
2. 元素未启用
3. 元素不在屏幕上

**解决**:
```go
// 检查元素状态
if elem.Get_CurrentIsEnabled() == 0 {
    return fmt.Errorf("元素未启用")
}

// 检查支持的 Pattern
patterns, _, _ := ppv.PollForPotentialSupportedPatterns(elem)
```

### Q4: 消息发送失败（微信）

**可能原因**:
1. 搜索框不支持 ValuePattern
2. 发送按钮位置变化

**解决**:
```go
// 使用多种方式查找输入框
msgEdit := bot.FindElementByAutomationId("chat_input_field")
if msgEdit == nil {
    msgEdit = bot.FindElementByClassName("mmui::ChatInputField")
}
```

---

## 10. 附录

### 10.1 参考资源

- [Microsoft UI Automation 官方文档](https://learn.microsoft.com/zh-cn/windows/win32/winauto/entry-uiauto-win32)
- [UI Automation Property IDs](https://learn.microsoft.com/zh-cn/windows/win32/winauto/uiauto-entry-propids)
- [UI Automation Control Pattern IDs](https://learn.microsoft.com/zh-cn/windows/win32/winauto/uiauto-controlpattern-ids)
- [UI Automation Control Types](https://learn.microsoft.com/zh-cn/windows/win32/winauto/uiauto-controltypesoverview)

### 10.2 调试工具

| 工具 | 说明 |
|------|------|
| [Inspect](https://learn.microsoft.com/zh-cn/windows/win32/winauto/inspect-objects) | UI Automation 检查工具 |
| [Accessibility Insights](https://accessibilityinsights.io/) | 可访问性测试工具 |
| Spy++ | 窗口查看工具（Visual Studio 内置） |

### 10.3 术语表

| 术语 | 英文 | 说明 |
|------|------|------|
| UI Automation | UI Automation | Windows UI 自动化框架 |
| COM | Component Object Model | 组件对象模型 |
| Pattern | Control Pattern | 控件模式 |
| Element | UI Automation Element | UI 自动化元素 |
| Cache Request | Cache Request | 缓存请求 |
| Tree Walker | Tree Walker | 树遍历器 |
| Condition | Condition | 查询条件 |
| VARIANT | VARIANT | 变体类型 |
| BSTR | Basic String | Basic 字符串类型 |
| HRESULT | HRESULT | COM 返回值类型 |

### 10.4 版本历史

| 版本 | 日期 | 主要变更 |
|------|------|----------|
| v1.0.0 | 2024-03-20 | 初始版本，核心功能完成 |

---

**文档结束**

*本文档由 go-element 项目组维护*
