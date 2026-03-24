# go-element API 文档

> 文档生成时间：2026-03-23  
> 项目版本：基于 go-element main 分支  
> 遵循规范：6A工作流

---

## 项目概述

**go-element** 是一个用于Windows窗体自动化的Go语言UI Automation库，通过Windows COM接口与应用程序中的UI元素进行高效交互和操作。

### 核心特性
- **极致性能**：深度集成 `IUIAutomationCacheRequest`，实现属性和模式的批量获取，大幅减少跨进程通信 (IPC) 开销
- **零反射设计**：彻底移除反射逻辑，改为显式属性填充，提升运行时效率并增强类型安全
- **完善的模式支持**：内置常用自动化模式 (Patterns) 的包装，包括：`Value`、`Invoke`、`Toggle`、`ExpandCollapse` 及 `SelectionItem`
- **资源安全**：严格的 COM 对象生命周期管理，通过显式的 `Release()` 调用确保无内存泄漏和句柄耗尽
- **便捷查询 API**：`Element` 结构体集成 `FindByName` 和 `FindByAutomationId` 等快捷方法

### 安装
```shell
go get -u github.com/auuunya/go-element
```

### 编译说明
- **64 位系统 (默认)**:
  ```shell
  GOOS=windows GOARCH=amd64 go build .
  ```
- **32 位系统 (高兼容性)**:
  ```shell
  GOOS=windows GOARCH=386 go build .
  ```

---

## 一、初始化与清理

### 1.1 CoInitialize
**功能描述**：初始化COM库，必须在使用任何UI Automation功能前调用。

**函数签名**：
```go
func CoInitialize() error
```

**返回值**：
- `error`：成功返回nil，失败返回错误信息

**使用示例**：
```go
err := uia.CoInitialize()
if err != nil {
    log.Fatal(err)
}
defer uia.CoUninitialize()
```

**注意事项**：
- 该函数会锁定OS线程 (`runtime.LockOSThread()`)
- 必须与 `CoUninitialize()` 配对使用

---

### 1.2 CoUninitialize
**功能描述**：清理COM库，释放相关资源。

**函数签名**：
```go
func CoUninitialize()
```

**注意事项**：
- 必须与 `CoInitialize()` 配对使用
- 建议使用 `defer` 确保执行
- 会解锁OS线程 (`runtime.UnlockOSThread()`)

---

## 二、窗口操作

### 2.1 GetWindowForString
**功能描述**：根据窗口类名或窗口标题获取窗口句柄。

**函数签名**：
```go
func GetWindowForString(classname, windowname string) (uintptr, error)
```

**参数说明**：
| 参数 | 类型 | 说明 | 示例 |
|------|------|------|------|
| classname | string | 窗口类名 | "Notepad"、"Chrome_WidgetWin_1" |
| windowname | string | 窗口标题 | "无标题 - 记事本" |

**返回值**：
| 类型 | 说明 |
|------|------|
| uintptr | 窗口句柄 |
| error | 未找到窗口返回 `ErrorNotFoundWindow` |

**使用示例**：
```go
// 通过类名查找记事本窗口
hwnd, err := uia.GetWindowForString("Notepad", "")
if err != nil {
    log.Fatal("窗口未找到")
}

// 通过窗口标题查找
hwnd, err := uia.GetWindowForString("", "无标题 - 记事本")
```

---

### 2.2 FindWindowW
**功能描述**：封装Windows API FindWindowW，查找顶层窗口。

**函数签名**：
```go
func FindWindowW(lpclass, lpwindow string) uintptr
```

**参数说明**：
- `lpclass`：窗口类名（可为空字符串）
- `lpwindow`：窗口标题（可为空字符串）

**返回值**：
- `uintptr`：窗口句柄，未找到返回0

---

### 2.3 FindWindowExW
**功能描述**：封装Windows API FindWindowExW，用于查找子窗口。

**函数签名**：
```go
func FindWindowExW(phwdn, chwdn uintptr, lpclass, lpwindow string) uintptr
```

**参数说明**：
| 参数 | 类型 | 说明 |
|------|------|------|
| phwdn | uintptr | 父窗口句柄 |
| chwdn | uintptr | 子窗口句柄（从该窗口之后开始查找） |
| lpclass | string | 窗口类名 |
| lpwindow | string | 窗口标题 |

---

## 三、IUIAutomation 核心接口

### 3.1 CreateInstance
**功能描述**：创建IUIAutomation实例，这是使用UI Automation的入口点。

**函数签名**：
```go
func CreateInstance(clsid, riid *syscall.GUID, clsctx CLSCTX) (unsafe.Pointer, error)
```

**参数说明**：
| 参数 | 类型 | 说明 | 常用值 |
|------|------|------|--------|
| clsid | *syscall.GUID | 类标识符 | `CLSID_CUIAutomation` |
| riid | *syscall.GUID | 接口标识符 | `IID_IUIAutomation` |
| clsctx | CLSCTX | 上下文类型 | `CLSCTX_INPROC_SERVER` / `CLSCTX_ALL` |

**返回值**：
- `unsafe.Pointer`：接口指针
- `error`：错误信息

**使用示例**：
```go
instance, err := uia.CreateInstance(
    uia.CLSID_CUIAutomation,
    uia.IID_IUIAutomation,
    uia.CLSCTX_INPROC_SERVER,
)
if err != nil {
    log.Fatal(err)
}
ppv := uia.NewIUIAutomation(uia.NewIUnKnown(instance))
defer ppv.Release()
```

---

### 3.2 ElementFromHandle
**功能描述**：从窗口句柄获取UI元素。

**函数签名**：
```go
func (v *IUIAutomation) ElementFromHandle(in uintptr) (*IUIAutomationElement, error)
```

**参数说明**：
- `in`：窗口句柄

**返回值**：
- `*IUIAutomationElement`：UI元素对象
- `error`：错误信息

**使用示例**：
```go
root, err := ppv.ElementFromHandle(hwnd)
if err != nil {
    log.Fatal(err)
}
defer root.Release()
```

---

### 3.3 GetRootElement
**功能描述**：获取桌面根元素，可用于遍历所有窗口。

**函数签名**：
```go
func (v *IUIAutomation) GetRootElement() (*IUIAutomationElement, error)
```

**使用示例**：
```go
desktop, err := ppv.GetRootElement()
if err != nil {
    log.Fatal(err)
}
defer desktop.Release()
```

---

### 3.4 ElementFromPoint
**功能描述**：根据屏幕坐标获取UI元素。

**函数签名**：
```go
func (v *IUIAutomation) ElementFromPoint(in *TagPoint) (*IUIAutomationElement, error)
```

**参数说明**：
- `in`：屏幕坐标点 `TagPoint{X: x, Y: y}`

**使用示例**：
```go
element, err := ppv.ElementFromPoint(&uia.TagPoint{X: 100, Y: 200})
```

---

### 3.5 GetFocusedElement
**功能描述**：获取当前焦点元素。

**函数签名**：
```go
func (v *IUIAutomation) GetFocusedElement() (*IUIAutomationElement, error)
```

**使用示例**：
```go
focused, err := ppv.GetFocusedElement()
if err != nil {
    log.Fatal(err)
}
defer focused.Release()
```

---

### 3.6 CreateCacheRequest
**功能描述**：创建缓存请求对象，用于高性能属性获取。

**函数签名**：
```go
func (v *IUIAutomation) CreateCacheRequest() (*IUIAutomationCacheRequest, error)
```

**使用示例**：
```go
cacheRequest, err := ppv.CreateCacheRequest()
if err != nil {
    log.Fatal(err)
}
defer cacheRequest.Release()

cacheRequest.AddProperty(uia.UIA_NamePropertyId)
cacheRequest.AddProperty(uia.UIA_ClassNamePropertyId)
cacheRequest.AddPattern(uia.UIA_ValuePatternId)
```

---

### 3.7 CreateTreeWalker
**功能描述**：创建树遍历器，用于手动遍历UI树。

**函数签名**：
```go
func (v *IUIAutomation) CreateTreeWalker(in *IUIAutomationCondition) (*IUIAutomationTreeWalker, error)
```

**参数说明**：
- `in`：遍历条件，传入 `nil` 或 `CreateTrueCondition()` 遍历所有元素

---

### 3.8 条件创建方法

#### CreateTrueCondition
**功能描述**：创建始终为真的条件，匹配所有元素。

```go
func (v *IUIAutomation) CreateTrueCondition() *IUIAutomationCondition
```

#### CreateFalseCondition
**功能描述**：创建始终为假的条件，不匹配任何元素。

```go
func (v *IUIAutomation) CreateFalseCondition() (*IUIAutomationCondition, error)
```

#### CreatePropertyCondition
**功能描述**：创建属性条件，根据属性值筛选元素。

```go
func (v *IUIAutomation) CreatePropertyCondition(id PropertyId, value VARIANT) (*IUIAutomationCondition, error)
```

**使用示例**：
```go
// 创建名称条件
nameVariant, _ := uia.VariantFromString("确定")
condition, err := ppv.CreatePropertyCondition(uia.UIA_NamePropertyId, nameVariant)
```

#### CreateAndCondition
**功能描述**：创建AND条件，同时满足多个条件。

```go
func (v *IUIAutomation) CreateAndCondition(in, in2 *IUIAutomationCondition) (*IUIAutomationCondition, error)
```

#### CreateOrCondition
**功能描述**：创建OR条件，满足任一条件即可。

```go
func (v *IUIAutomation) CreateOrCondition(in, in2 *IUIAutomationCondition) (*IUIAutomationCondition, error)
```

#### CreateNotCondition
**功能描述**：创建NOT条件，取反指定条件。

```go
func (v *IUIAutomation) CreateNotCondition(in *IUIAutomationCondition) (*IUIAutomationCondition, error)
```

---

### 3.9 其他常用方法

#### CompareElements
**功能描述**：比较两个元素是否相同。

```go
func (v *IUIAutomation) CompareElements(in, in2 *IUIAutomationElement) (int32, error)
```

#### Get_ControlViewWalker
**功能描述**：获取控件视图遍历器。

```go
func (v *IUIAutomation) Get_ControlViewWalker() *IUIAutomationTreeWalker
```

#### Get_ContentViewWalker
**功能描述**：获取内容视图遍历器。

```go
func (v *IUIAutomation) Get_ContentViewWalker() *IUIAutomationTreeWalker
```

---

## 四、Element 元素操作

### 4.1 Element 结构体
**功能描述**：封装IUIAutomationElement，提供便捷的属性访问和操作。

**结构定义**：
```go
type Element struct {
    UIAutoElement               *IUIAutomationElement
    CurrentAcceleratorKey       string
    CurrentAccessKey            string
    CurrentAutomationId         string
    CurrentBoundingRectangle    *TagRect
    CurrentClassName            string
    CurrentControlType          ControlTypeId
    CurrentName                 string
    CurrentProcessId            int32
    CurrentIsEnabled            int32
    CurrentLocalizedControlType string
    SupportedPatterns           []PatternId
    Child                       []*Element
}
```

---

### 4.2 TraverseUIElementTree
**功能描述**：遍历UI元素树，返回带缓存的Element结构。这是最常用的元素遍历方法。

**函数签名**：
```go
func TraverseUIElementTree(ppv *IUIAutomation, root *IUIAutomationElement) *Element
```

**参数说明**：
| 参数 | 类型 | 说明 |
|------|------|------|
| ppv | *IUIAutomation | IUIAutomation实例 |
| root | *IUIAutomationElement | 根元素 |

**返回值**：
- `*Element`：包含子元素的Element树结构

**使用示例**：
```go
root, _ := ppv.ElementFromHandle(hwnd)
tree := uia.TraverseUIElementTree(ppv, root)
uia.TreeString(tree, 0)
```

**缓存属性**：
该方法自动缓存以下属性和模式：
- 属性：Name、ClassName、ControlType、AutomationId、IsEnabled、ProcessId、LocalizedControlType
- 模式：Value、Invoke、Selection、ExpandCollapse、Toggle

---

### 4.3 TreeString
**功能描述**：以树形结构打印元素信息，用于调试。

**函数签名**：
```go
func TreeString(root *Element, level int)
```

**参数说明**：
- `root`：根元素
- `level`：起始缩进层级，通常传0

---

### 4.4 元素查找方法

#### FindByName
**功能描述**：根据名称查找元素（深度优先搜索）。

**函数签名**：
```go
func (e *Element) FindByName(name string) *Element
```

**使用示例**：
```go
editor := tree.FindByName("文本编辑器")
if editor != nil {
    // 操作元素
}
```

#### FindByAutomationId
**功能描述**：根据AutomationId查找元素（深度优先搜索）。

**函数签名**：
```go
func (e *Element) FindByAutomationId(id string) *Element
```

**使用示例**：
```go
button := tree.FindByAutomationId("btnOK")
```

---

### 4.5 模式获取方法

#### GetValuePattern
**功能描述**：获取Value模式，用于设置或获取文本值。

**函数签名**：
```go
func (e *Element) GetValuePattern() (*IUIAutomationValuePattern, error)
```

**使用示例**：
```go
vp, err := editor.GetValuePattern()
if err == nil {
    defer vp.Release()
    vp.SetValue("Hello World")
    value, _ := vp.Get_CurrentValue()
    fmt.Println("当前值:", value)
}
```

#### GetInvokePattern
**功能描述**：获取Invoke模式，用于执行点击操作。

**函数签名**：
```go
func (e *Element) GetInvokePattern() (*IUIAutomationInvokePattern, error)
```

**使用示例**：
```go
ip, err := element.GetInvokePattern()
if err == nil {
    defer ip.Release()
    ip.Invoke() // 执行点击
}
```

#### GetTogglePattern
**功能描述**：获取Toggle模式，用于切换状态。

**函数签名**：
```go
func (e *Element) GetTogglePattern() (*IUIAutomationTogglePattern, error)
```

#### GetExpandCollapsePattern
**功能描述**：获取展开/折叠模式。

**函数签名**：
```go
func (e *Element) GetExpandCollapsePattern() (*IUIAutomationExpandCollapsePattern, error)
```

#### GetSelectionItemPattern
**功能描述**：获取选择项模式。

**函数签名**：
```go
func (e *Element) GetSelectionItemPattern() (*IUIAutomationSelectionItemPattern, error)
```

---

### 4.6 属性填充方法

#### Populate
**功能描述**：填充Element的属性值。

**函数签名**：
```go
func (e *Element) Populate(cached bool)
```

**参数说明**：
- `cached`：是否使用缓存数据

---

## 五、自动化模式 (Patterns)

### 5.1 IUIAutomationValuePattern
**功能描述**：值模式，用于文本输入和获取。

**方法列表**：
| 方法 | 签名 | 说明 |
|------|------|------|
| SetValue | `func (v *IUIAutomationValuePattern) SetValue(val string) error` | 设置文本值 |
| Get_CurrentValue | `func (v *IUIAutomationValuePattern) Get_CurrentValue() (string, error)` | 获取当前值 |
| Release | `func (v *IUIAutomationValuePattern) Release() uint32` | 释放资源 |

**使用示例**：
```go
vp, _ := element.GetValuePattern()
defer vp.Release()
vp.SetValue("自动化输入文本")
value, _ := vp.Get_CurrentValue()
fmt.Println("当前值:", value)
```

---

### 5.2 IUIAutomationInvokePattern
**功能描述**：调用模式，用于执行按钮点击等操作。

**方法列表**：
| 方法 | 签名 | 说明 |
|------|------|------|
| Invoke | `func (v *IUIAutomationInvokePattern) Invoke() error` | 执行调用操作 |
| Release | `func (v *IUIAutomationInvokePattern) Release() uint32` | 释放资源 |

**使用示例**：
```go
ip, _ := element.GetInvokePattern()
defer ip.Release()
ip.Invoke() // 执行点击
```

---

### 5.3 IUIAutomationTogglePattern
**功能描述**：切换模式，用于复选框、单选按钮等。

**方法列表**：
| 方法 | 签名 | 说明 |
|------|------|------|
| Toggle | `func (v *IUIAutomationTogglePattern) Toggle() error` | 切换状态 |
| Get_CurrentToggleState | `func (v *IUIAutomationTogglePattern) Get_CurrentToggleState() (ToggleState, error)` | 获取当前状态 |
| Release | `func (v *IUIAutomationTogglePattern) Release() uint32` | 释放资源 |

**ToggleState 枚举**：
```go
const (
    ToggleState_Off           ToggleState = 0 // 关闭状态
    ToggleState_On            ToggleState = 1 // 开启状态
    ToggleState_Indeterminate ToggleState = 2 // 不确定状态
)
```

**使用示例**：
```go
tp, _ := element.GetTogglePattern()
defer tp.Release()
state, _ := tp.Get_CurrentToggleState()
if state == uia.ToggleState_Off {
    tp.Toggle() // 切换到开启状态
}
```

---

### 5.4 IUIAutomationExpandCollapsePattern
**功能描述**：展开/折叠模式，用于菜单、树节点、下拉框等。

**方法列表**：
| 方法 | 签名 | 说明 |
|------|------|------|
| Expand | `func (v *IUIAutomationExpandCollapsePattern) Expand() error` | 展开元素 |
| Collapse | `func (v *IUIAutomationExpandCollapsePattern) Collapse() error` | 折叠元素 |
| Release | `func (v *IUIAutomationExpandCollapsePattern) Release() uint32` | 释放资源 |

**使用示例**：
```go
ecp, _ := element.GetExpandCollapsePattern()
defer ecp.Release()
ecp.Expand()   // 展开
ecp.Collapse() // 折叠
```

---

### 5.5 IUIAutomationSelectionItemPattern
**功能描述**：选择项模式，用于列表项选择。

**方法列表**：
| 方法 | 签名 | 说明 |
|------|------|------|
| Select | `func (v *IUIAutomationSelectionItemPattern) Select() error` | 选中该项 |
| Release | `func (v *IUIAutomationSelectionItemPattern) Release() uint32` | 释放资源 |

**使用示例**：
```go
sip, _ := element.GetSelectionItemPattern()
defer sip.Release()
sip.Select() // 选中该项
```

---

## 六、条件查询

### 6.1 SearchElem
**功能描述**：使用自定义函数搜索单个元素（深度优先）。

**函数签名**：
```go
func SearchElem(elem *Element, searchFunc SearchFunc) *Element
```

**SearchFunc 类型**：
```go
type SearchFunc func(elem *Element) bool
```

**使用示例**：
```go
// 查找类名为Edit的第一个元素
result := uia.SearchElem(tree, func(elem *Element) bool {
    return elem.CurrentClassName == "Edit"
})

// 查找启用状态的按钮
result := uia.SearchElem(tree, func(elem *Element) bool {
    return elem.CurrentControlType == uia.UIA_ButtonControlTypeId && 
           elem.CurrentIsEnabled == 1
})
```

---

### 6.2 FindElems
**功能描述**：使用自定义函数查找所有匹配元素。

**函数签名**：
```go
func FindElems(elem *Element, searchFunc SearchFunc) []*Element
```

**使用示例**：
```go
// 查找所有按钮
buttons := uia.FindElems(tree, func(elem *Element) bool {
    return elem.CurrentControlType == uia.UIA_ButtonControlTypeId
})

for _, btn := range buttons {
    fmt.Println("按钮:", btn.CurrentName)
}
```

---

## 七、属性ID和模式ID常量

### 7.1 属性ID (PropertyId)

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
    UIA_HelpTextPropertyId             PropertyId = 30013
    UIA_ClickablePointPropertyId       PropertyId = 30014
    UIA_CulturePropertyId              PropertyId = 30015
    UIA_IsControlElementPropertyId     PropertyId = 30016
    UIA_IsContentElementPropertyId     PropertyId = 30017
    UIA_LabeledByPropertyId            PropertyId = 30018
    UIA_IsPasswordPropertyId           PropertyId = 30019
    UIA_NativeWindowHandlePropertyId   PropertyId = 30020
    UIA_ItemTypePropertyId             PropertyId = 30021
    UIA_IsOffscreenPropertyId          PropertyId = 30022
    UIA_OrientationPropertyId          PropertyId = 30023
    UIA_FrameworkIdPropertyId          PropertyId = 30024
    UIA_IsRequiredForFormPropertyId    PropertyId = 30025
    UIA_ItemStatusPropertyId           PropertyId = 30026
    // ... 更多属性ID见 id.go
)
```

---

### 7.2 模式ID (PatternId)

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
    UIA_MultipleViewPatternId      PatternId = 10008
    UIA_WindowPatternId            PatternId = 10009
    UIA_SelectionItemPatternId     PatternId = 10010
    UIA_DockPatternId              PatternId = 10011
    UIA_TablePatternId             PatternId = 10012
    UIA_TableItemPatternId         PatternId = 10013
    UIA_TextPatternId              PatternId = 10014
    UIA_TogglePatternId            PatternId = 10015
    UIA_TransformPatternId         PatternId = 10016
    UIA_ScrollItemPatternId        PatternId = 10017
    UIA_LegacyIAccessiblePatternId PatternId = 10018
    UIA_ItemContainerPatternId     PatternId = 10019
    UIA_VirtualizedItemPatternId   PatternId = 10020
    UIA_SynchronizedInputPatternId PatternId = 10021
    // ... 更多模式ID见 id.go
)
```

---

### 7.3 控件类型ID (ControlTypeId)

```go
const (
    UIA_ButtonControlTypeId       ControlTypeId = 50000
    UIA_CalendarControlTypeId     ControlTypeId = 50001
    UIA_CheckBoxControlTypeId     ControlTypeId = 50002
    UIA_ComboBoxControlTypeId     ControlTypeId = 50003
    UIA_EditControlTypeId         ControlTypeId = 50004
    UIA_HyperlinkControlTypeId    ControlTypeId = 50005
    UIA_ImageControlTypeId        ControlTypeId = 50006
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
    UIA_CustomControlTypeId       ControlTypeId = 50025
    UIA_GroupControlTypeId        ControlTypeId = 50026
    UIA_ThumbControlTypeId        ControlTypeId = 50027
    UIA_DataGridControlTypeId     ControlTypeId = 50028
    UIA_DataItemControlTypeId     ControlTypeId = 50029
    UIA_DocumentControlTypeId     ControlTypeId = 50030
    UIA_SplitButtonControlTypeId  ControlTypeId = 50031
    UIA_WindowControlTypeId       ControlTypeId = 50032
    UIA_PaneControlTypeId         ControlTypeId = 50033
    UIA_HeaderControlTypeId       ControlTypeId = 50034
    UIA_HeaderItemControlTypeId   ControlTypeId = 50035
    UIA_TableControlTypeId        ControlTypeId = 50036
    UIA_TitleBarControlTypeId     ControlTypeId = 50037
    UIA_SeparatorControlTypeId    ControlTypeId = 50038
    UIA_SemanticZoomControlTypeId ControlTypeId = 50039
    UIA_AppBarControlTypeId       ControlTypeId = 50040
)
```

---

## 八、数据类型

### 8.1 TagRect
**功能描述**：矩形区域结构。

```go
type TagRect struct {
    Left   int32
    Top    int32
    Right  int32
    Bottom int32
}
```

**使用示例**：
```go
rect := element.CurrentBoundingRectangle
width := rect.Right - rect.Left
height := rect.Bottom - rect.Top
```

---

### 8.2 TagPoint
**功能描述**：点坐标结构。

```go
type TagPoint struct {
    X int32
    Y int32
}
```

---

### 8.3 VARIANT
**功能描述**：COM变体类型，用于传递各种类型的值。

```go
type VARIANT struct {
    VT  TagVarenum
    Val int64
}
```

**构造函数**：
```go
// 创建通用变体
func NewVariant(vt TagVarenum, val int64) VARIANT

// 从字符串创建变体
func VariantFromString(s string) (VARIANT, error)
```

**TagVarenum 枚举**：
```go
const (
    VT_EMPTY   TagVarenum = 0
    VT_NULL    TagVarenum = 1
    VT_I2      TagVarenum = 2
    VT_I4      TagVarenum = 3
    VT_R4      TagVarenum = 4
    VT_R8      TagVarenum = 5
    VT_BSTR    TagVarenum = 8   // 字符串
    VT_BOOL    TagVarenum = 11
    VT_UNKNOWN TagVarenum = 13
    // ... 更多类型见 variant.go
)
```

---

### 8.4 TreeScope
**功能描述**：树遍历范围，用于指定搜索范围。

```go
var (
    TreeScope_None        TreeScope = 0x0
    TreeScope_Element     TreeScope = 0x1  // 元素本身
    TreeScope_Children    TreeScope = 0x2  // 直接子元素
    TreeScope_Descendants TreeScope = 0x4  // 所有后代元素
    TreeScope_Parent      TreeScope = 0x8  // 父元素
    TreeScope_Ancestors   TreeScope = 0x10 // 所有祖先元素
    TreeScope_Subtree     TreeScope = TreeScope_Element | TreeScope_Children | TreeScope_Descendants
)
```

---

### 8.5 CLSCTX
**功能描述**：COM类上下文，用于指定创建实例的方式。

```go
var (
    CLSCTX_INPROC_SERVER CLSCTX = 0x1
    CLSCTX_INPROC_HANDLER CLSCTX = 0x2
    CLSCTX_LOCAL_SERVER CLSCTX = 0x4
    CLSCTX_REMOTE_SERVER CLSCTX = 0x10
    CLSCTX_ALL CLSCTX = CLSCTX_INPROC_SERVER | CLSCTX_INPROC_HANDLER | CLSCTX_LOCAL_SERVER | CLSCTX_REMOTE_SERVER
)
```

---

## 九、事件处理

### 9.1 AddAutomationEventHandler
**功能描述**：添加自动化事件处理器。

**函数签名**：
```go
func (v *IUIAutomation) AddAutomationEventHandler(opt *EventHandler) error
```

**EventHandler 结构**：
```go
type EventHandler struct {
    EventId      UIA_EventId
    Element      *IUIAutomationElement
    Scope        TreeScope
    CacheRequest *IUIAutomationCacheRequest
    Handler      *IUIAutomationEventHandler
}
```

---

### 9.2 AddFocusChangedEventHandler
**功能描述**：添加焦点变化事件处理器。

**函数签名**：
```go
func (v *IUIAutomation) AddFocusChangedEventHandler(in *IUIAutomationCacheRequest, in2 *IUIAutomationFocusChangedEventHandler) error
```

---

### 9.3 AddPropertyChangedEventHandler
**功能描述**：添加属性变化事件处理器。

**函数签名**：
```go
func (v *IUIAutomation) AddPropertyChangedEventHandler(opt *ChangeEventHandler) error
```

**ChangeEventHandler 结构**：
```go
type ChangeEventHandler struct {
    Element       *IUIAutomationElement
    Scope         TreeScope
    CacheRequest  *IUIAutomationCacheRequest
    Handler       *IUIAutomationPropertyChangedEventHandler
    PropertyArray *TagSafeArray
}
```

---

### 9.4 RemoveAllEventHandlers
**功能描述**：移除所有事件处理器。

**函数签名**：
```go
func (v *IUIAutomation) RemoveAllEventHandlers() error
```

---

### 9.5 事件ID常量

```go
const (
    UIA_ToolTipOpenedEventId                       UIA_EventId = 20000
    UIA_ToolTipClosedEventId                       UIA_EventId = 20001
    UIA_StructureChangedEventId                    UIA_EventId = 20002
    UIA_MenuOpenedEventId                          UIA_EventId = 20003
    UIA_AutomationPropertyChangedEventId           UIA_EventId = 20004
    UIA_AutomationFocusChangedEventId              UIA_EventId = 20005
    UIA_AsyncContentLoadedEventId                  UIA_EventId = 20006
    UIA_MenuClosedEventId                          UIA_EventId = 20007
    UIA_LayoutInvalidatedEventId                   UIA_EventId = 20008
    UIA_Invoke_InvokedEventId                      UIA_EventId = 20009
    UIA_SelectionItem_ElementAddedToSelectionEventId UIA_EventId = 20010
    UIA_SelectionItem_ElementRemovedFromSelectionEventId UIA_EventId = 20011
    UIA_SelectionItem_ElementSelectedEventId       UIA_EventId = 20012
    UIA_Selection_InvalidatedEventId               UIA_EventId = 20013
    UIA_Text_TextSelectionChangedEventId           UIA_EventId = 20014
    UIA_Text_TextChangedEventId                    UIA_EventId = 20015
    UIA_Window_WindowOpenedEventId                 UIA_EventId = 20016
    UIA_Window_WindowClosedEventId                 UIA_EventId = 20017
    // ... 更多事件ID见 id.go
)
```

---

### 3.10 类型转换方法

#### RectToVariant
**功能描述**：将矩形结构转换为VARIANT类型。

**函数签名**：
```go
func (v *IUIAutomation) RectToVariant(in *TagRect) (*VARIANT, error)
```

**参数说明**：
- `in`：矩形结构 `TagRect{Left, Top, Right, Bottom}`

**返回值**：
- `*VARIANT`：转换后的VARIANT
- `error`：错误信息

---

#### VariantToRect
**功能描述**：将VARIANT转换为矩形结构。

**函数签名**：
```go
func (v *IUIAutomation) VariantToRect(in *VARIANT) (*TagRect, error)
```

**参数说明**：
- `in`：VARIANT类型

**返回值**：
- `*TagRect`：矩形结构
- `error`：错误信息

---

#### IntNativeArrayToSafeArray
**功能描述**：将整数原生数组转换为安全数组。

**函数签名**：
```go
func (v *IUIAutomation) IntNativeArrayToSafeArray(in, in2 int32) (*TagSafeArray, error)
```

**参数说明**：
- `in`：整数数组指针
- `in2`：数组长度

**返回值**：
- `*TagSafeArray`：安全数组
- `error`：错误信息

---

#### IntSafeArrayToNativeArray
**功能描述**：将安全数组转换为整数原生数组。

**函数签名**：
```go
func (v *IUIAutomation) IntSafeArrayToNativeArray(in *TagSafeArray) (int32, int32, error)
```

**参数说明**：
- `in`：安全数组

**返回值**：
- `int32`：数组指针
- `int32`：数组长度
- `error`：错误信息

---

#### PollForPotentialSupportedPatterns
**功能描述**：轮询元素潜在支持的模式。

**函数签名**：
```go
func (v *IUIAutomation) PollForPotentialSupportedPatterns(in *IUIAutomationElement) (*TagSafeArray, *TagSafeArray, error)
```

**参数说明**：
- `in`：UI元素

**返回值**：
- `*TagSafeArray`：模式ID数组
- `*TagSafeArray`：模式名称数组
- `error`：错误信息

---

#### PollForPotentialSupportedProperties
**功能描述**：轮询元素潜在支持的属性。

**函数签名**：
```go
func (v *IUIAutomation) PollForPotentialSupportedProperties(in *IUIAutomationElement) (*TagSafeArray, *TagSafeArray, error)
```

**参数说明**：
- `in`：UI元素

**返回值**：
- `*TagSafeArray`：属性ID数组
- `*TagSafeArray`：属性名称数组
- `error`：错误信息

---

## 十、完整使用示例

### 10.1 记事本自动化
```go
package main

import (
    "fmt"
    uia "github.com/auuunya/go-element"
)

func main() {
    // 初始化COM
    uia.CoInitialize()
    defer uia.CoUninitialize()

    // 查找记事本窗口
    hwnd, err := uia.GetWindowForString("Notepad", "")
    if err != nil {
        fmt.Println("未找到记事本窗口，请先打开记事本")
        return
    }

    // 创建IUIAutomation实例
    instance, err := uia.CreateInstance(
        uia.CLSID_CUIAutomation,
        uia.IID_IUIAutomation,
        uia.CLSCTX_ALL,
    )
    if err != nil {
        fmt.Println("创建实例失败:", err)
        return
    }
    ppv := uia.NewIUIAutomation(uia.NewIUnKnown(instance))
    defer ppv.Release()

    // 从窗口句柄获取元素
    root, err := ppv.ElementFromHandle(hwnd)
    if err != nil {
        fmt.Println("获取元素失败:", err)
        return
    }
    defer root.Release()

    // 遍历UI树
    tree := uia.TraverseUIElementTree(ppv, root)

    // 查找文本编辑器并输入文本
    if editor := tree.FindByName("文本编辑器"); editor != nil {
        if vp, err := editor.GetValuePattern(); err == nil {
            defer vp.Release()
            vp.SetValue("你好，来自 go-element 的自动输入！")
            fmt.Println("文本输入成功")
        }
    }
}
```

---

### 10.2 浏览器UI树遍历
```go
package main

import (
    uia "github.com/auuunya/go-element"
)

func main() {
    uia.CoInitialize()
    defer uia.CoUninitialize()

    // 查找Chrome浏览器窗口
    hwnd, err := uia.GetWindowForString("Chrome_WidgetWin_1", "")
    if err != nil {
        panic("未找到Chrome窗口")
    }

    instance, _ := uia.CreateInstance(
        uia.CLSID_CUIAutomation,
        uia.IID_IUIAutomation,
        uia.CLSCTX_INPROC_SERVER,
    )
    ppv := uia.NewIUIAutomation(uia.NewIUnKnown(instance))
    defer ppv.Release()

    root, _ := ppv.ElementFromHandle(hwnd)
    defer root.Release()

    // 遍历并打印UI树
    elems := uia.TraverseUIElementTree(ppv, root)
    uia.TreeString(elems, 0)
}
```

---

### 10.3 按钮点击操作
```go
package main

import (
    "fmt"
    uia "github.com/auuunya/go-element"
)

func main() {
    uia.CoInitialize()
    defer uia.CoUninitialize()

    hwnd, _ := uia.GetWindowForString("Notepad", "")
    instance, _ := uia.CreateInstance(
        uia.CLSID_CUIAutomation,
        uia.IID_IUIAutomation,
        uia.CLSCTX_ALL,
    )
    ppv := uia.NewIUIAutomation(uia.NewIUnKnown(instance))
    defer ppv.Release()

    root, _ := ppv.ElementFromHandle(hwnd)
    defer root.Release()

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

---

### 10.4 复选框切换
```go
package main

import (
    "fmt"
    uia "github.com/auuunya/go-element"
)

func main() {
    uia.CoInitialize()
    defer uia.CoUninitialize()

    // ... 初始化代码 ...

    tree := uia.TraverseUIElementTree(ppv, root)

    // 查找复选框并切换状态
    if checkbox := tree.FindByName("记住密码"); checkbox != nil {
        if tp, err := checkbox.GetTogglePattern(); err == nil {
            defer tp.Release()
            
            state, _ := tp.Get_CurrentToggleState()
            fmt.Printf("当前状态: %v\n", state)
            
            tp.Toggle() // 切换状态
        }
    }
}
```

---

### 10.5 自定义条件搜索
```go
package main

import (
    "fmt"
    uia "github.com/auuunya/go-element"
)

func main() {
    uia.CoInitialize()
    defer uia.CoUninitialize()

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

---

### 10.6 微信自动化示例
**功能描述**：自动选择联系人并发送消息。

```go
package main

import (
    "fmt"
    "log"
    "time"
    uia "github.com/auuunya/go-element"
)

func main() {
    // 初始化COM
    uia.CoInitialize()
    defer uia.CoUninitialize()

    // 查找微信窗口
    hwnd, err := uia.GetWindowForString("WeChatMainWndForPC", "")
    if err != nil {
        log.Fatal("未找到微信窗口，请先打开微信")
    }

    // 创建IUIAutomation实例
    instance, err := uia.CreateInstance(
        uia.CLSID_CUIAutomation,
        uia.IID_IUIAutomation,
        uia.CLSCTX_ALL,
    )
    if err != nil {
        log.Fatal("创建实例失败:", err)
    }
    ppv := uia.NewIUIAutomation(uia.NewIUnKnown(instance))
    defer ppv.Release()

    // 从窗口句柄获取元素
    root, err := ppv.ElementFromHandle(hwnd)
    if err != nil {
        log.Fatal("获取元素失败:", err)
    }
    defer root.Release()

    // 遍历UI树
    tree := uia.TraverseUIElementTree(ppv, root)

    // 1. 在搜索框输入联系人名称
    if searchBox := tree.FindByName("搜索"); searchBox != nil {
        if vp, err := searchBox.GetValuePattern(); err == nil {
            defer vp.Release()
            vp.SetValue("文件传输助手")
            time.Sleep(500 * time.Millisecond) // 等待搜索结果
        }
    }

    // 2. 重新遍历获取搜索结果
    tree = uia.TraverseUIElementTree(ppv, root)

    // 3. 查找并点击联系人（精确匹配）
    if contact := tree.FindByName("文件传输助手"); contact != nil {
        if ip, err := contact.GetInvokePattern(); err == nil {
            defer ip.Release()
            ip.Invoke() // 点击联系人
            time.Sleep(300 * time.Millisecond)
        }
    }

    // 4. 重新遍历获取聊天窗口
    tree = uia.TraverseUIElementTree(ppv, root)

    // 5. 在消息输入框输入文本
    if inputBox := tree.FindByAutomationId("chat_input_field"); inputBox != nil {
        if vp, err := inputBox.GetValuePattern(); err == nil {
            defer vp.Release()
            vp.SetValue("这是来自 go-element 的自动消息！")
            fmt.Println("消息已输入")
        }
    }
}
```

**注意事项**：
- 微信窗口类名为 `WeChatMainWndForPC`
- 搜索框和输入框需要使用 ValuePattern
- 联系人点击需要使用 InvokePattern
- 操作之间需要适当延时等待UI更新

---

## 十一、错误处理

### 11.1 常见错误类型

```go
var (
    // 窗口未找到错误
    ErrorNotFoundWindow = errors.New("not found window")
    
    // BSTR指针为空错误
    ErrorBstrPointerNil = errors.New("BSTR pointer is nil")
)
```

---

### 11.2 HResult错误
**功能描述**：将COM HRESULT转换为Go错误。

**函数签名**：
```go
func HResult(ret uintptr) error
```

**常见HRESULT值**：
| HRESULT | 值 | 说明 |
|---------|-----|------|
| S_OK | 0x00000000 | 成功 |
| E_FAIL | 0x80004005 | 一般性失败 |
| E_INVALIDARG | 0x80070057 | 参数无效 |
| E_NOINTERFACE | 0x80004002 | 不支持该接口 |
| E_OUTOFMEMORY | 0x8007000E | 内存不足 |

---

### 11.3 错误处理最佳实践

```go
// 推荐的错误处理方式
element, err := ppv.ElementFromHandle(hwnd)
if err != nil {
    if errors.Is(err, uia.ErrorNotFoundWindow) {
        fmt.Println("窗口未找到")
        return
    }
    fmt.Printf("其他错误: %v\n", err)
    return
}
defer element.Release()
```

---

## 十二、注意事项

### 12.1 COM初始化
- ✅ 必须在使用前调用 `CoInitialize()`
- ✅ 使用后必须调用 `CoUninitialize()`
- ✅ 建议使用 `defer` 确保清理

### 12.2 资源释放
- ✅ 所有COM对象使用完毕后必须调用 `Release()`
- ✅ 使用 `defer` 确保资源释放
- ⚠️ 避免循环引用导致内存泄漏

### 12.3 线程安全
- ⚠️ COM对象应在同一线程中使用
- ⚠️ `CoInitialize()` 会锁定OS线程
- ⚠️ 不要在goroutine间传递COM对象

### 12.4 架构兼容
- ⚠️ 编译时注意目标系统架构（32位/64位）
- ⚠️ 64位程序无法操作32位进程的UI
- ⚠️ 32位程序无法操作64位进程的UI

### 12.5 性能优化
- ✅ 使用 `TraverseUIElementTree` 自动缓存属性
- ✅ 批量获取属性减少IPC开销
- ✅ 避免频繁调用 `GetCurrentPattern`

---

## 十三、项目文件结构

```
go-element/
├── com.go              # COM初始化和窗口查找
├── uiautomation.go     # IUIAutomation核心接口
├── element.go          # Element结构和便捷方法
├── client.go           # Pattern客户端接口
├── condition.go        # 条件查询接口
├── constants.go        # 常量定义
├── id.go               # 属性ID、模式ID、控件类型ID
├── enum.go             # 枚举类型
├── typedef.go          # 类型定义
├── variant.go          # VARIANT类型
├── unknown.go          # IUnknown基础接口
├── dispatch.go         # IDispatch接口
├── accessible.go       # IAccessible接口
├── provider.go         # Provider接口
├── textserv.go         # 文本服务接口
├── drop.go             # 拖放接口
├── uia.go              # UIA相关定义
└── example/            # 示例代码
    ├── calc_demo/      # 计算器自动化示例
    ├── notepad_demo/   # 记事本自动化示例
    ├── find_demo/      # 元素查找示例
    ├── tree_demo/      # UI树遍历示例
    ├── json_demo/      # JSON导出示例
    └── software/       # 软件信息示例
```

---

## 十四、参考资源

- [Microsoft UI Automation 官方文档](https://learn.microsoft.com/zh-cn/windows/win32/winauto/entry-uiauto-win32)
- [UI Automation Property IDs](https://learn.microsoft.com/zh-cn/windows/win32/winauto/uiauto-entry-propids)
- [UI Automation Control Pattern IDs](https://learn.microsoft.com/zh-cn/windows/win32/winauto/uiauto-controlpattern-ids)
- [UI Automation Control Types](https://learn.microsoft.com/zh-cn/windows/win32/winauto/uiauto-controltypesoverview)

---

**文档版本**: v1.0  
**最后更新**: 2026-03-23  
**维护者**: go-element 项目组
