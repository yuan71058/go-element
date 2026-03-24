// Package uiautomation 提供 Windows UI Automation 的 Go 语言封装
package uiautomation

import (
	"fmt"
	"reflect"
	"strings"
	"syscall"
	"unsafe"
)

// Windows OLE Automation DLL 和函数声明
var (
	oleAut             = syscall.NewLazyDLL("OleAut32.dll") // OLE Automation DLL
	procSysFreeString  = oleAut.NewProc("SysFreeString")   // 释放 BSTR 字符串
	procSysAllocString = oleAut.NewProc("SysAllocString")  // 分配 BSTR 字符串
)

// Element UI 元素的高级封装结构体
// 提供了 UI 元素的属性缓存和便捷操作方法
type Element struct {
	UIAutoElement               *IUIAutomationElement        // 底层 UI Automation 元素接口
	CurrentAcceleratorKey       string                       // 快捷键
	CurrentAccessKey            string                       // 访问键（如 Alt+F）
	CurrentAriaProperties       string                       // ARIA 属性
	CurrentAriaRole             string                       // ARIA 角色
	CurrentAutomationId         string                       // 自动化ID（唯一标识符）
	CurrentBoundingRectangle    *TagRect                     // 边界矩形（位置和大小）
	CurrentClassName            string                       // 类名
	CurrentControllerFor        *IUIAutomationElementArray   // 控制器元素数组
	CurrentControlType          ControlTypeId                // 控件类型
	CurrentCulture              int32                        // 文化/区域设置
	CurrentDescribedBy          *IUIAutomationElementArray   // 描述元素数组
	CurrentFlowsTo              *IUIAutomationElementArray   // 流向元素数组
	CurrentFrameworkId          string                       // 框架ID（如 WPF, Win32）
	CurrentHasKeyboardFocus     int32                        // 是否有键盘焦点
	CurrentHelpText             string                       // 帮助文本
	CurrentIsContentElement     int32                        // 是否为内容元素
	CurrentIsControlElement     int32                        // 是否为控件元素
	CurrentIsDataValidForForm   int32                        // 表单数据是否有效
	CurrentIsEnabled            int32                        // 是否启用
	CurrentIsKeyboardFocusable  int32                        // 是否可获取键盘焦点
	CurrentIsOffscreen          int32                        // 是否在屏幕外
	CurrentIsPassword           int32                        // 是否为密码字段
	CurrentIsRequiredForForm    int32                        // 表单是否必填
	CurrentItemStatus           string                       // 项目状态
	CurrentItemType             string                       // 项目类型
	CurrentLabeledBy            *IUIAutomationElement        // 标签元素
	CurrentLocalizedControlType string                       // 本地化控件类型
	CurrentName                 string                       // 名称
	CurrentNativeWindowHandle   uintptr                      // 原生窗口句柄
	CurrentOrientation          OrientationType              // 方向
	CurrentProcessId            int32                        // 进程ID
	CurrentProviderDescription  string                       // 提供者描述

	SupportedPatterns []PatternId `json:"supported_patterns,omitempty"` // 支持的模式列表
	Child             []*Element  `json:"child,omitempty"`              // 子元素列表
}

// NewElement 创建新的 Element 实例
// 参数:
//   - uiaumation: 底层 UI Automation 元素接口
// 返回: 初始化的 Element 指针
func NewElement(uiaumation *IUIAutomationElement) *Element {
	return &Element{
		UIAutoElement: uiaumation,
	}
}

// FormatString 格式化输出元素的所有属性
// 返回: 格式化后的字符串
func (e *Element) FormatString() string {
	if e == nil {
		return ""
	}
	elemType := reflect.TypeOf(e)
	if elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
	}
	elemValue := reflect.ValueOf(e)
	if elemValue.Kind() == reflect.Ptr {
		elemValue = elemValue.Elem()
	}
	buf := ""
	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		fieldName := field.Name
		fieldValue := elemValue.Field(i)
		if fieldValue.Kind() == reflect.Ptr && !fieldValue.IsNil() {
			if fieldValue.CanAddr() {
				v := fieldValue.Addr().Interface()
				buf += fmt.Sprintf("[%s]:[%v],", fieldName, v)
			}
		} else {
			buf += fmt.Sprintf("[%s]:[%v],", fieldName, fieldValue.Interface())
		}
	}
	return buf
}

// SetUIAutomation 设置底层 UI Automation 元素接口
func (e *Element) SetUIAutomation(uiaumation *IUIAutomationElement) {
	e.UIAutoElement = uiaumation
}

// AcceleratorKey 获取快捷键属性
func (e *Element) AcceleratorKey() error {
	val, err := e.UIAutoElement.Get_CurrentAcceleratorKey()
	e.CurrentAcceleratorKey = val
	return err
}

// AccessKey 获取访问键属性
func (e *Element) AccessKey() error {
	val, err := e.UIAutoElement.Get_CurrentAccessKey()
	e.CurrentAccessKey = val
	return err
}

// AriaProperties 获取 ARIA 属性
func (e *Element) AriaProperties() error {
	val, err := e.UIAutoElement.Get_CurrentAriaProperties()
	e.CurrentAriaProperties = val
	return err
}

// AriaRole 获取 ARIA 角色属性
func (e *Element) AriaRole() error {
	val, err := e.UIAutoElement.Get_CurrentAriaRole()
	e.CurrentAriaRole = val
	return err
}

// AutomationId 获取自动化ID属性
func (e *Element) AutomationId() error {
	val, err := e.UIAutoElement.Get_CurrentAutomationId()
	e.CurrentAutomationId = val
	return err
}

// BoundingRectangle 获取边界矩形属性
func (e *Element) BoundingRectangle() {
	val := e.UIAutoElement.Get_CurrentBoundingRectangle()
	e.CurrentBoundingRectangle = val
}

// ClassName 获取类名属性
func (e *Element) ClassName() error {
	val, err := e.UIAutoElement.Get_CurrentClassName()
	e.CurrentClassName = val
	return err
}

// ControllerFor 获取控制器元素数组
func (e *Element) ControllerFor() {
	val := e.UIAutoElement.Get_CurrentControllerFor()
	e.CurrentControllerFor = val
}

// ControlType 获取控件类型
func (e *Element) ControlType() {
	val := e.UIAutoElement.Get_CurrentControlType()
	e.CurrentControlType = val
}

// Culture 获取文化/区域设置
func (e *Element) Culture() {
	val := e.UIAutoElement.Get_CurrentCulture()
	e.CurrentCulture = val
}

// DescribedBy 获取描述元素数组
func (e *Element) DescribedBy() {
	val := e.UIAutoElement.Get_CurrentDescribedBy()
	e.CurrentDescribedBy = val
}

// FlowsTo 获取流向元素数组
func (e *Element) FlowsTo() {
	val := e.UIAutoElement.Get_CurrentFlowsTo()
	e.CurrentFlowsTo = val
}

// FrameworkId 获取框架ID属性
func (e *Element) FrameworkId() error {
	val, err := e.UIAutoElement.Get_CurrentFrameworkId()
	e.CurrentFrameworkId = val
	return err
}

// HasKeyboardFocus 获取是否有键盘焦点
func (e *Element) HasKeyboardFocus() {
	val := e.UIAutoElement.Get_CurrentHasKeyboardFocus()
	e.CurrentHasKeyboardFocus = val
}

// HelpText 获取帮助文本属性
func (e *Element) HelpText() error {
	val, err := e.UIAutoElement.Get_CurrentHelpText()
	e.CurrentHelpText = val
	return err
}

// IsControlElement 获取是否为控件元素
func (e *Element) IsControlElement() {
	val := e.UIAutoElement.Get_CurrentIsControlElement()
	e.CurrentIsControlElement = val
}

// IsContentElement 获取是否为内容元素
func (e *Element) IsContentElement() {
	val := e.UIAutoElement.Get_CurrentIsContentElement()
	e.CurrentIsContentElement = val
}

// IsDataValidForForm 获取表单数据是否有效
func (e *Element) IsDataValidForForm() {
	val := e.UIAutoElement.Get_CurrentIsDataValidForForm()
	e.CurrentIsDataValidForForm = val
}

// IsEnabled 获取是否启用
func (e *Element) IsEnabled() {
	val := e.UIAutoElement.Get_CurrentIsEnabled()
	e.CurrentIsEnabled = val
}

// IsKeyboardFocusable 获取是否可获取键盘焦点
func (e *Element) IsKeyboardFocusable() {
	val := e.UIAutoElement.Get_CurrentIsKeyboardFocusable()
	e.CurrentIsKeyboardFocusable = val
}

// IsOffscreen 获取是否在屏幕外
func (e *Element) IsOffscreen() {
	val := e.UIAutoElement.Get_CurrentIsOffscreen()
	e.CurrentIsOffscreen = val
}

// IsPassword 获取是否为密码字段
func (e *Element) IsPassword() {
	val := e.UIAutoElement.Get_CurrentIsPassword()
	e.CurrentIsPassword = val
}

// IsRequiredForForm 获取表单是否必填
func (e *Element) IsRequiredForForm() {
	val := e.UIAutoElement.Get_CurrentIsRequiredForForm()
	e.CurrentIsRequiredForForm = val
}

// ItemStatus 获取项目状态属性
func (e *Element) ItemStatus() error {
	val, err := e.UIAutoElement.Get_CurrentItemStatus()
	e.CurrentItemStatus = val
	return err
}

// ItemType 获取项目类型属性
func (e *Element) ItemType() error {
	val, err := e.UIAutoElement.Get_CurrentItemType()
	e.CurrentItemType = val
	return err
}

// LabeledBy 获取标签元素
func (e *Element) LabeledBy() {
	val := e.UIAutoElement.Get_CurrentLabeledBy()
	e.CurrentLabeledBy = val
}

// LocalizedControlType 获取本地化控件类型
func (e *Element) LocalizedControlType() error {
	val, err := e.UIAutoElement.Get_CurrentLocalizedControlType()
	e.CurrentLocalizedControlType = val
	return err
}

// Name 获取名称属性
func (e *Element) Name() error {
	val, err := e.UIAutoElement.Get_CurrentName()
	e.CurrentName = val
	return err
}

// NativeWindowHandle 获取原生窗口句柄
func (e *Element) NativeWindowHandle() {
	val := e.UIAutoElement.Get_CurrentNativeWindowHandle()
	e.CurrentNativeWindowHandle = val
}

// Orientation 获取方向属性
func (e *Element) Orientation() {
	val := e.UIAutoElement.Get_CurrentOrientation()
	e.CurrentOrientation = val
}

// ProcessId 获取进程ID
func (e *Element) ProcessId() {
	val := e.UIAutoElement.Get_CurrentProcessId()
	e.CurrentProcessId = val
}

// ProviderDescription 获取提供者描述
func (e *Element) ProviderDescription() error {
	val, err := e.UIAutoElement.Get_CurrentProviderDescription()
	e.CurrentProviderDescription = val
	return err
}

// SearchFunc 元素搜索函数类型
// 参数: 要检查的元素
// 返回: 是否匹配
type SearchFunc func(elem *Element) bool

// SearchElem 在元素树中搜索匹配的元素（深度优先，返回第一个匹配）
// 参数:
//   - elem: 起始元素
//   - searchFunc: 搜索条件函数
// 返回: 匹配的元素，未找到返回 nil
func SearchElem(elem *Element, searchFunc SearchFunc) *Element {
	if elem == nil || searchFunc == nil {
		return nil
	}
	if searchFunc(elem) {
		return elem
	}
	for _, childElem := range elem.Child {
		if found := SearchElem(childElem, searchFunc); found != nil {
			return found
		}
	}
	return nil
}

// FindElems 在元素树中搜索所有匹配的元素
// 参数:
//   - elem: 起始元素
//   - searchFunc: 搜索条件函数
// 返回: 所有匹配的元素数组
func FindElems(elem *Element, searchFunc SearchFunc) (elems []*Element) {
	if searchFunc(elem) {
		elems = append(elems, elem)
	}
	for _, childElem := range elem.Child {
		if found := FindElems(childElem, searchFunc); found != nil {
			elems = append(elems, found...)
		}
	}
	return elems
}

// bstr2str 将 BSTR 字符串转换为 Go 字符串
// 参数: BSTR 指针
// 返回: Go 字符串
func bstr2str(bstr uintptr) string {
	if bstr == 0 {
		return ""
	}
	// 在 Windows 上 BSTR 是以 null 结尾的 UTF-16 字符串
	// 指针直接指向字符数组。我们手动计算长度以避免创建巨型切片。
	p := (*uint16)(unsafe.Pointer(bstr))
	n := 0
	for ptr := p; *ptr != 0; ptr = (*uint16)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + 2)) {
		n++
		if n > 0x100000 { // 限制最大 1MB 字符，防止无限循环
			break
		}
	}
	// 使用 syscall.UTF16ToString 将其转换为 Go 字符串
	var slice []uint16
	header := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	header.Data = bstr
	header.Len = n
	header.Cap = n
	return syscall.UTF16ToString(slice)
}

// string2Bstr 将 Go 字符串转换为 BSTR 字符串
// 参数: Go 字符串
// 返回: BSTR 指针和可能的错误
func string2Bstr(str string) (uintptr, error) {
	// 将Go字符串转换为UTF-16编码
	utf16Str, err := syscall.UTF16PtrFromString(str)
	if err != nil {
		return 0, err
	}
	bstrPtr, _, _ := procSysAllocString.Call(
		uintptr(unsafe.Pointer(utf16Str)),
	)
	if bstrPtr == 0 {
		return 0, ErrorBstrPointerNil
	}
	return bstrPtr, nil
}

// Populate 填充元素的属性
// 参数:
//   - cached: 是否使用缓存属性（true=缓存属性，false=当前属性）
func (e *Element) Populate(cached bool) {
	if cached {
		e.CurrentName, _ = e.UIAutoElement.Get_CachedName()
		e.CurrentClassName, _ = e.UIAutoElement.Get_CachedClassName()
		e.CurrentControlType = e.UIAutoElement.Get_CachedControlType()
		e.CurrentAutomationId, _ = e.UIAutoElement.Get_CachedAutomationId()
		e.CurrentIsEnabled = e.UIAutoElement.Get_CachedIsEnabled()
		e.CurrentProcessId = e.UIAutoElement.Get_CachedProcessId()
		e.CurrentLocalizedControlType, _ = e.UIAutoElement.Get_CachedLocalizedControlType()
	} else {
		e.CurrentName, _ = e.UIAutoElement.Get_CurrentName()
		e.CurrentClassName, _ = e.UIAutoElement.Get_CurrentClassName()
		e.CurrentControlType = e.UIAutoElement.Get_CurrentControlType()
		e.CurrentAutomationId, _ = e.UIAutoElement.Get_CurrentAutomationId()
		e.CurrentIsEnabled = e.UIAutoElement.Get_CurrentIsEnabled()
		e.CurrentProcessId = e.UIAutoElement.Get_CurrentProcessId()
		e.CurrentLocalizedControlType, _ = e.UIAutoElement.Get_CurrentLocalizedControlType()
	}

	// 检查支持的常用模式
	commonPatterns := []PatternId{
		UIA_ValuePatternId,
		UIA_InvokePatternId,
		UIA_SelectionPatternId,
		UIA_ExpandCollapsePatternId,
		UIA_TogglePatternId,
	}
	for _, pid := range commonPatterns {
		var pattern *IUnKnown
		var err error
		if cached {
			pattern, err = e.UIAutoElement.GetCachedPattern(pid)
		} else {
			pattern, err = e.UIAutoElement.GetCurrentPattern(pid)
		}
		if err == nil && pattern != nil {
			e.SupportedPatterns = append(e.SupportedPatterns, pid)
			pattern.Release()
		}
	}
}

// TraverseUIElementTree 遍历 UI 元素树并返回高级封装的 Element 结构
// 使用缓存请求提高性能，减少跨进程通信开销
// 参数:
//   - ppv: UI Automation 接口
//   - root: 根元素
// 返回: 封装后的元素树
func TraverseUIElementTree(ppv *IUIAutomation, root *IUIAutomationElement) *Element {
	cacheRequest, _ := ppv.CreateCacheRequest()
	if cacheRequest != nil {
		defer cacheRequest.Release()
		// 添加需要缓存的属性
		cacheRequest.AddProperty(UIA_NamePropertyId)
		cacheRequest.AddProperty(UIA_ClassNamePropertyId)
		cacheRequest.AddProperty(UIA_ControlTypePropertyId)
		cacheRequest.AddProperty(UIA_AutomationIdPropertyId)
		cacheRequest.AddProperty(UIA_IsEnabledPropertyId)
		cacheRequest.AddProperty(UIA_ProcessIdPropertyId)
		cacheRequest.AddProperty(UIA_LocalizedControlTypePropertyId)

		// 添加需要缓存的模式
		cacheRequest.AddPattern(UIA_ValuePatternId)
		cacheRequest.AddPattern(UIA_InvokePatternId)
		cacheRequest.AddPattern(UIA_SelectionPatternId)
		cacheRequest.AddPattern(UIA_ExpandCollapsePatternId)
		cacheRequest.AddPattern(UIA_TogglePatternId)

		cacheRequest.Put_TreeScope(TreeScope_Subtree)
	}

	return traverseUIElementTree(ppv, root, cacheRequest)
}

// traverseUIElementTree 内部递归函数：遍历 UI 元素树
func traverseUIElementTree(ppv *IUIAutomation, root *IUIAutomationElement, cacheRequest *IUIAutomationCacheRequest) *Element {
	var elementArr *IUIAutomationElementArray
	var err error

	if cacheRequest != nil {
		elementArr, err = root.FindAllBuildCache(TreeScope_Children, CreateTrueCondition(ppv), cacheRequest)
	} else {
		elementArr, err = root.FindAll(TreeScope_Children, CreateTrueCondition(ppv))
	}

	newElement := NewElement(root)
	if cacheRequest != nil {
		// 如果支持缓存，由于 FindAllBuildCache 只处理子节点，根节点需要手动更新缓存
		cachedRoot, _ := root.BuildUpdatedCache(cacheRequest)
		if cachedRoot != nil {
			newElement.UIAutoElement = cachedRoot
			newElement.Populate(true)
		} else {
			newElement.Populate(false)
		}
	} else {
		newElement.Populate(false)
	}

	if err == nil && elementArr != nil {
		defer elementArr.Release()
		arrLen := elementArr.Get_Length()
		for i := 0; i < int(arrLen); i++ {
			elem, _ := elementArr.GetElement(int32(i))
			if elem != nil {
				childElem := traverseUIElementTree(ppv, elem, cacheRequest)
				newElement.Child = append(newElement.Child, childElem)
			}
		}
	}

	return newElement
}

// TreeString 以树形结构打印元素
// 参数:
//   - root: 根元素
//   - level: 缩进层级
func TreeString(root *Element, level int) {
	if root == nil {
		return
	}
	fmt.Printf("%s- %v\n", getIndentation(level), root)
	for _, child := range root.Child {
		TreeString(child, level+1)
	}
}

// String 实现 Stringer 接口
func (e *Element) String() string {
	return fmt.Sprintf("Element{Name: %q, Class: %q, Type: %d}", e.CurrentName, e.CurrentClassName, e.CurrentControlType)
}

// getIndentation 获取缩进字符串
func getIndentation(level int) string {
	return strings.Repeat("  ", level)
}

// FindByName 通过名称查找子元素
// 参数: name - 元素名称
// 返回: 匹配的元素，未找到返回 nil
func (e *Element) FindByName(name string) *Element {
	if e.CurrentName == name {
		return e
	}
	for _, child := range e.Child {
		if found := child.FindByName(name); found != nil {
			return found
		}
	}
	return nil
}

// FindByAutomationId 通过 AutomationId 查找子元素
// 参数: id - 元素的 AutomationId
// 返回: 匹配的元素，未找到返回 nil
func (e *Element) FindByAutomationId(id string) *Element {
	if e.CurrentAutomationId == id {
		return e
	}
	for _, child := range e.Child {
		if found := child.FindByAutomationId(id); found != nil {
			return found
		}
	}
	return nil
}

// GetValuePattern 获取 ValuePattern 接口
// 用于设置或获取文本输入控件的值
// 返回: ValuePattern 接口和可能的错误
func (e *Element) GetValuePattern() (*IUIAutomationValuePattern, error) {
	unk, err := e.UIAutoElement.GetCurrentPattern(UIA_ValuePatternId)
	if err != nil {
		return nil, err
	}
	return (*IUIAutomationValuePattern)(unsafe.Pointer(unk)), nil
}

// GetInvokePattern 获取 InvokePattern 接口
// 用于调用按钮等可点击控件
// 返回: InvokePattern 接口和可能的错误
func (e *Element) GetInvokePattern() (*IUIAutomationInvokePattern, error) {
	unk, err := e.UIAutoElement.GetCurrentPattern(UIA_InvokePatternId)
	if err != nil {
		return nil, err
	}
	return (*IUIAutomationInvokePattern)(unsafe.Pointer(unk)), nil
}

// GetTogglePattern 获取 TogglePattern 接口
// 用于切换复选框等控件的状态
// 返回: TogglePattern 接口和可能的错误
func (e *Element) GetTogglePattern() (*IUIAutomationTogglePattern, error) {
	unk, err := e.UIAutoElement.GetCurrentPattern(UIA_TogglePatternId)
	if err != nil {
		return nil, err
	}
	return (*IUIAutomationTogglePattern)(unsafe.Pointer(unk)), nil
}

// GetExpandCollapsePattern 获取 ExpandCollapsePattern 接口
// 用于展开或折叠树节点、菜单等控件
// 返回: ExpandCollapsePattern 接口和可能的错误
func (e *Element) GetExpandCollapsePattern() (*IUIAutomationExpandCollapsePattern, error) {
	unk, err := e.UIAutoElement.GetCurrentPattern(UIA_ExpandCollapsePatternId)
	if err != nil {
		return nil, err
	}
	return (*IUIAutomationExpandCollapsePattern)(unsafe.Pointer(unk)), nil
}

// GetSelectionItemPattern 获取 SelectionItemPattern 接口
// 用于选择列表项等控件
// 返回: SelectionItemPattern 接口和可能的错误
func (e *Element) GetSelectionItemPattern() (*IUIAutomationSelectionItemPattern, error) {
	unk, err := e.UIAutoElement.GetCurrentPattern(UIA_SelectionItemPatternId)
	if err != nil {
		return nil, err
	}
	return (*IUIAutomationSelectionItemPattern)(unsafe.Pointer(unk)), nil
}
