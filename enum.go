// Package uiautomation 提供 Windows UI Automation 的 Go 语言封装
package uiautomation

// TxtHitResult 文本命中测试结果
type TxtHitResult int32

// 文本命中测试结果常量
const (
	TXTHITRESULT_NOHIT        TxtHitResult = iota // 未命中
	TXTHITRESULT_TRANSPARENT                     // 透明区域
	TXTHITRESULT_CLOSE                           // 接近
	TXTHITRESULT_HIT                             // 命中
)

// StructureChangeType 结构变更类型
type StructureChangeType int32

// 结构变更类型常量
const (
	StructureChangeType_ChildAdded          StructureChangeType = iota // 子元素已添加
	StructureChangeType_ChildRemoved                                  // 子元素已移除
	StructureChangeType_ChildrenInvalidated                           // 子元素已失效
	StructureChangeType_ChildrenBulkAdded                             // 批量添加子元素
	StructureChangeType_ChildrenBulkRemoved                           // 批量移除子元素
	StructureChangeType_ChildrenReordered                             // 子元素已重新排序
)

// WindowVisualState 窗口可视状态
type WindowVisualState int32

// 窗口可视状态常量
// 参考: https://learn.microsoft.com/zh-cn/windows/win32/api/uiautomationcore/ne-uiautomationcore-windowvisualstate
const (
	WindowVisualState_Normal   WindowVisualState = iota // 正常状态
	WindowVisualState_Maximized                          // 最大化
	WindowVisualState_Minimized                          // 最小化
)

// ToggleState 切换状态
type ToggleState int32

// 切换状态常量
// 参考: https://learn.microsoft.com/zh-cn/windows/win32/api/uiautomationcore/ne-uiautomationcore-togglestate
const (
	ToggleState_Off          ToggleState = iota // 关闭
	ToggleState_On                             // 打开
	ToggleState_Indeterminate                  // 不确定（三态复选框）
)

// ZoomUnit 缩放单位
type ZoomUnit int32

// 缩放单位常量
// 参考: https://learn.microsoft.com/zh-cn/windows/win32/api/uiautomationcore/ne-uiautomationcore-zoomunit
const (
	ZoomUnit_NoAmount       ZoomUnit = iota // 无缩放
	ZoomUnit_LargeDecrement                 // 大幅度缩小
	ZoomUnit_SmallDecrement                 // 小幅度缩小
	ZoomUnit_LargeIncrement                 // 大幅度放大
	ZoomUnit_SmallIncrement                 // 小幅度放大
)

// WindowInteractionState 窗口交互状态
type WindowInteractionState int32

// 窗口交互状态常量
// 参考: https://learn.microsoft.com/zh-cn/windows/win32/api/uiautomationcore/ne-uiautomationcore-windowinteractionstate
const (
	WindowInteractionState_Running              WindowInteractionState = iota // 正在运行
	WindowInteractionState_Closing                                         // 正在关闭
	WindowInteractionState_ReadyForUserInteraction                         // 准备好用户交互
	WindowInteractionState_BlockedByModalWindow                            // 被模态窗口阻塞
	WindowInteractionState_NotResponding                                   // 无响应
)

// SupportedTextSelection 支持的文本选择类型
type SupportedTextSelection int32

// 支持的文本选择类型常量
// 参考: https://learn.microsoft.com/zh-cn/windows/win32/api/uiautomationcore/ne-uiautomationcore-supportedtextselection
const (
	SupportedTextSelection_None    SupportedTextSelection = iota // 不支持选择
	SupportedTextSelection_Single                               // 单选
	SupportedTextSelection_Multiple                             // 多选
)

// ScrollAmount 滚动量
type ScrollAmount int32

// 滚动量常量
// 参考: https://learn.microsoft.com/zh-cn/windows/win32/api/uiautomationcore/ne-uiautomationcore-scrollamount
const (
	ScrollAmount_LargeDecrement ScrollAmount = iota // 大幅度后退
	ScrollAmount_SmallDecrement                     // 小幅度后退
	ScrollAmount_NoAmount                           // 无滚动
	ScrollAmount_LargeIncrement                     // 大幅度前进
	ScrollAmount_SmallIncrement                     // 小幅度前进
)

// RowOrColumnMajor 行或列主导
type RowOrColumnMajor int32

// 行或列主导常量
// 参考: https://learn.microsoft.com/zh-cn/windows/win32/api/uiautomationcore/ne-uiautomationcore-roworcolumnmajor
const (
	RowOrColumnMajor_RowMajor     RowOrColumnMajor = iota // 行主导
	RowOrColumnMajor_ColumnMajor                          // 列主导
	RowOrColumnMajor_Indeterminate                        // 不确定
)

// TextPatternRangeEndpoint 文本范围端点
type TextPatternRangeEndpoint int32

// 文本范围端点常量
// 参考: https://learn.microsoft.com/zh-cn/windows/win32/api/uiautomationcore/ne-uiautomationcore-textpatternrangeendpoint
const (
	TextPatternRangeEndpoint_Start TextPatternRangeEndpoint = iota // 起点
	TextPatternRangeEndpoint_End                                   // 终点
)

// TextUnit 文本单位
type TextUnit int32

// 文本单位常量
// 参考: https://learn.microsoft.com/zh-cn/windows/win32/api/uiautomationcore/ne-uiautomationcore-textunit
const (
	TextUnit_Character TextUnit = iota // 字符
	TextUnit_Format                    // 格式
	TextUnit_Word                      // 单词
	TextUnit_Line                      // 行
	TextUnit_Paragraph                 // 段落
	TextUnit_Page                      // 页
	TextUnit_Document                  // 文档
)

// TextArrtibuteId 文本属性ID
type TextArrtibuteId int32

// 文本属性ID常量
// 参考: https://learn.microsoft.com/zh-cn/windows/win32/winauto/uiauto-textattribute-ids
const (
	// 文本属性ID定义
)

// OrientationType 方向类型
type OrientationType int32

// 方向类型常量
const (
	OrientationType_None       OrientationType = iota // 无方向
	OrientationType_Horizontal                        // 水平
	OrientationType_Vertical                          // 垂直
)

// TagFuncKind 函数类型
type TagFuncKind int32

// 函数类型常量
const (
	FUNC_VIRTUAL    TagFuncKind = iota // 虚函数
	FUNC_PUREVIRTUAL                   // 纯虚函数
	FUNC_NONVIRTUAL                    // 非虚函数
	FUNC_STATIC                        // 静态函数
	FUNC_DISPATCH                      // 分发函数
)

// TagCallConv 调用约定
type TagCallConv int32

// 调用约定常量
const (
	CC_FASTCALL   TagCallConv = iota // 快速调用
	CC_CDECL                         // C 声明
	CC_MSCPASCAL                     // MSC Pascal
	CC_PASCAL                        // Pascal
	CC_MACPASCAL                     // Mac Pascal
	CC_STDCALL                       // 标准调用
	CC_FPFASTCALL                    // FP 快速调用
	CC_SYSCALL                       // 系统调用
	CC_MPWCDECL                      // MPW C 声明
	CC_MPWPASCAL                     // MPW Pascal
	CC_MAX                           // 最大值
)
