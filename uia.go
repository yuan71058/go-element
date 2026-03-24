// Package uiautomation 提供 Windows UI Automation 的 Go 语言封装
// 用于自动化 Windows 应用程序的 UI 元素操作
package uiautomation

// PropertyId UI Automation 属性标识符类型
// 用于标识 UI 元素的各种属性，如名称、类名、控件类型等
type PropertyId int32

// UI Automation 属性常量定义
// 参考: https://learn.microsoft.com/zh-cn/windows/win32/winauto/uiauto-entry-propids
const (
	UIA_AcceleratorKeyPropertyId           PropertyId = 30006  // 快捷键属性
	UIA_AccessKeyPropertyId                PropertyId = 30007  // 访问键属性（如 Alt+F）
	UIA_AnnotationObjectsPropertyId        PropertyId = 30156  // 注解对象属性
	UIA_AnnotationTypesPropertyId          PropertyId = 30155  // 注解类型属性
	UIA_AriaPropertiesPropertyId           PropertyId = 30102  // ARIA 属性
	UIA_AriaRolePropertyId                 PropertyId = 30101  // ARIA 角色属性
	UIA_AutomationIdPropertyId             PropertyId = 30011  // 自动化ID属性（唯一标识符）
	UIA_BoundingRectanglePropertyId        PropertyId = 30001  // 边界矩形属性（元素位置和大小）
	UIA_CenterPointPropertyId              PropertyId = 30165  // 中心点属性
	UIA_ClassNamePropertyId                PropertyId = 30012  // 类名属性
	UIA_ClickablePointPropertyId           PropertyId = 30014  // 可点击点属性
	UIA_ControllerForPropertyId            PropertyId = 30104  // 控制器属性
	UIA_ControlTypePropertyId              PropertyId = 30003  // 控件类型属性
	UIA_CulturePropertyId                  PropertyId = 30015  // 文化/区域属性
	UIA_DescribedByPropertyId              PropertyId = 30105  // 描述属性
	UIA_FillColorPropertyId                PropertyId = 30160  // 填充颜色属性
	UIA_FillTypePropertyId                 PropertyId = 30162  // 填充类型属性
	UIA_FlowsFromPropertyId                PropertyId = 30148  // 流向来源属性
	UIA_FlowsToPropertyId                  PropertyId = 30106  // 流向目标属性
	UIA_FrameworkIdPropertyId              PropertyId = 30024  // 框架ID属性（如 WPF, Win32）
	UIA_FullDescriptionPropertyId          PropertyId = 30159  // 完整描述属性
	UIA_HasKeyboardFocusPropertyId         PropertyId = 30008  // 是否有键盘焦点属性
	UIA_HeadingLevelPropertyId             PropertyId = 30173  // 标题级别属性
	UIA_HelpTextPropertyId                 PropertyId = 30013  // 帮助文本属性
	UIA_IsContentElementPropertyId         PropertyId = 30017  // 是否为内容元素属性
	UIA_IsControlElementPropertyId         PropertyId = 30016  // 是否为控件元素属性
	UIA_IsDataValidForFormPropertyId       PropertyId = 30103  // 表单数据是否有效属性
	UIA_IsDialogPropertyId                 PropertyId = 30174  // 是否为对话框属性
	UIA_IsEnabledPropertyId                PropertyId = 30010  // 是否启用属性
	UIA_IsKeyboardFocusablePropertyId      PropertyId = 30009  // 是否可获取键盘焦点属性
	UIA_IsOffscreenPropertyId              PropertyId = 30022  // 是否在屏幕外属性
	UIA_IsPasswordPropertyId               PropertyId = 30019  // 是否为密码字段属性
	UIA_IsPeripheralPropertyId             PropertyId = 30150  // 是否为外围设备属性
	UIA_IsRequiredForFormPropertyId        PropertyId = 30025  // 表单是否必填属性
	UIA_ItemStatusPropertyId               PropertyId = 30026  // 项目状态属性
	UIA_ItemTypePropertyId                 PropertyId = 300021 // 项目类型属性
	UIA_LabeledByPropertyId                PropertyId = 30018  // 标签来源属性
	UIA_LandmarkTypePropertyId             PropertyId = 30157  // 地标类型属性
	UIA_LevelPropertyId                    PropertyId = 30154  // 层级属性
	UIA_LiveSettingPropertyId              PropertyId = 30135  // 实时设置属性
	UIA_LocalizedControlTypePropertyId     PropertyId = 30004  // 本地化控件类型属性
	UIA_LocalizedLandmarkTypePropertyId    PropertyId = 30158  // 本地化地标类型属性
	UIA_NamePropertyId                     PropertyId = 30005  // 名称属性
	UIA_NativeWindowHandlePropertyId       PropertyId = 30020  // 原生窗口句柄属性
	UIA_OptimizeForVisualContentPropertyId PropertyId = 30111  // 视觉内容优化属性
	UIA_OrientationPropertyId              PropertyId = 300023 // 方向属性
	UIA_OutlineColorPropertyId             PropertyId = 30161  // 轮廓颜色属性
	UIA_OutlineThicknessPropertyId         PropertyId = 30164  // 轮廓厚度属性
	UIA_PositionInSetPropertyId            PropertyId = 30152  // 集合中位置属性
	UIA_ProcessIdPropertyId                PropertyId = 30002  // 进程ID属性
	UIA_ProviderDescriptionPropertyId      PropertyId = 30107  // 提供者描述属性
	UIA_RotationPropertyId                 PropertyId = 30166  // 旋转属性
	UIA_RuntimeIdPropertyId                PropertyId = 30000  // 运行时ID属性
	UIA_SizePropertyId                     PropertyId = 30167  // 尺寸属性
	UIA_SizeOfSetPropertyId                PropertyId = 30153  // 集合大小属性
	UIA_VisualEffectsPropertyId            PropertyId = 30163  // 视觉效果属性
)

// PatternId UI Automation 模式标识符类型
// 用于标识 UI 元素支持的各种操作模式，如点击、设置值、展开折叠等
type PatternId int32

// UI Automation 模式常量定义
// 参考: https://learn.microsoft.com/zh-cn/windows/win32/winauto/uiauto-controlpattern-ids
const (
	UIA_AnnotationPatternId        PatternId = 10023 // 注解模式
	UIA_CustomNavigationPatternId  PatternId = 10033 // 自定义导航模式
	UIA_DockPatternId              PatternId = 10011 // 停靠模式
	UIA_DragPatternId              PatternId = 10030 // 拖拽模式
	UIA_DropTargetPatternId        PatternId = 10031 // 放置目标模式
	UIA_ExpandCollapsePatternId    PatternId = 10005 // 展开/折叠模式
	UIA_GridItemPatternId          PatternId = 10007 // 网格项模式
	UIA_GridPatternId              PatternId = 10006 // 网格模式
	UIA_InvokePatternId            PatternId = 10000 // 调用模式（按钮点击）
	UIA_ItemContainerPatternId     PatternId = 10019 // 项目容器模式
	UIA_LegacyIAccessiblePatternId PatternId = 10018 // 旧版 IAccessible 模式
	UIA_MultipleViewPatternId      PatternId = 10008 // 多视图模式
	UIA_ObjectModelPatternId       PatternId = 10022 // 对象模型模式
	UIA_RangeValuePatternId        PatternId = 10003 // 范围值模式（滑块）
	UIA_ScrollItemPatternId        PatternId = 10017 // 滚动项模式
	UIA_ScrollPatternId            PatternId = 10004 // 滚动模式
	UIA_SelectionItemPatternId     PatternId = 10010 // 选择项模式
	UIA_SelectionPatternId         PatternId = 10001 // 选择模式
	UIA_SpreadsheetPatternId       PatternId = 10026 // 电子表格模式
	UIA_SpreadsheetItemPatternId   PatternId = 10027 // 电子表格项模式
	UIA_StylesPatternId            PatternId = 10025 // 样式模式
	UIA_SynchronizedInputPatternId PatternId = 10021 // 同步输入模式
	UIA_TableItemPatternId         PatternId = 10013 // 表格项模式
	UIA_TablePatternId             PatternId = 10012 // 表格模式
	UIA_TextChildPatternId         PatternId = 10029 // 文本子项模式
	UIA_TextEditPatternId          PatternId = 10032 // 文本编辑模式
	UIA_TextPatternId              PatternId = 10014 // 文本模式
	UIA_TextPattern2Id             PatternId = 10024 // 文本模式2
	UIA_TogglePatternId            PatternId = 10015 // 切换模式（复选框）
	UIA_TransformPatternId         PatternId = 10016 // 变换模式（移动、调整大小、旋转）
	UIA_TransformPattern2Id        PatternId = 10028 // 变换模式2
	UIA_ValuePatternId             PatternId = 10002 // 值模式（文本输入）
	UIA_VirtualizedItemPatternId   PatternId = 10020 // 虚拟化项模式
	UIA_WindowPatternId            PatternId = 10009 // 窗口模式
)

// ControlTypeId UI Automation 控件类型标识符类型
// 用于标识 UI 元素的控件类型，如按钮、文本框、列表等
type ControlTypeId int32

// UI Automation 控件类型常量定义
// 参考: https://learn.microsoft.com/zh-cn/windows/win32/winauto/uiauto-controlpatternmapping
const (
	UIA_AppBarControlTypeId       ControlTypeId = 50040 // 应用栏控件
	UIA_ButtonControlTypeId       ControlTypeId = 50000 // 按钮控件
	UIA_CalendarControlTypeId     ControlTypeId = 50001 // 日历控件
	UIA_CheckBoxControlTypeId     ControlTypeId = 50002 // 复选框控件
	UIA_ComboBoxControlTypeId     ControlTypeId = 50003 // 组合框控件
	UIA_CustomControlTypeId       ControlTypeId = 50025 // 自定义控件
	UIA_DataGridControlTypeId     ControlTypeId = 50028 // 数据网格控件
	UIA_DataItemControlTypeId     ControlTypeId = 50029 // 数据项控件
	UIA_DocumentControlTypeId     ControlTypeId = 50030 // 文档控件
	UIA_EditControlTypeId         ControlTypeId = 50004 // 编辑控件（文本框）
	UIA_GroupControlTypeId        ControlTypeId = 50026 // 分组控件
	UIA_HeaderControlTypeId       ControlTypeId = 50034 // 标题控件
	UIA_HeaderItemControlTypeId   ControlTypeId = 50035 // 标题项控件
	UIA_HyperlinkControlTypeId    ControlTypeId = 50005 // 超链接控件
	UIA_ImageControlTypeId        ControlTypeId = 50006 // 图像控件
	UIA_ListControlTypeId         ControlTypeId = 50008 // 列表控件
	UIA_ListItemControlTypeId     ControlTypeId = 50007 // 列表项控件
	UIA_MenuBarControlTypeId      ControlTypeId = 50010 // 菜单栏控件
	UIA_MenuControlTypeId         ControlTypeId = 50009 // 菜单控件
	UIA_MenuItemControlTypeId     ControlTypeId = 50011 // 菜单项控件
	UIA_PaneControlTypeId         ControlTypeId = 50033 // 面板控件
	UIA_ProgressBarControlTypeId  ControlTypeId = 50012 // 进度条控件
	UIA_RadioButtonControlTypeId  ControlTypeId = 50013 // 单选按钮控件
	UIA_ScrollBarControlTypeId    ControlTypeId = 50014 // 滚动条控件
	UIA_SemanticZoomControlTypeId ControlTypeId = 50039 // 语义缩放控件
	UIA_SeparatorControlTypeId    ControlTypeId = 50038 // 分隔符控件
	UIA_SliderControlTypeId       ControlTypeId = 50015 // 滑块控件
	UIA_SpinnerControlTypeId      ControlTypeId = 50016 // 微调控件
	UIA_SplitButtonControlTypeId  ControlTypeId = 50031 // 拆分按钮控件
	UIA_StatusBarControlTypeId    ControlTypeId = 50017 // 状态栏控件
	UIA_TabControlTypeId          ControlTypeId = 50018 // 选项卡控件
	UIA_TabItemControlTypeId      ControlTypeId = 50019 // 选项卡项控件
	UIA_TableControlTypeId        ControlTypeId = 50036 // 表格控件
	UIA_TextControlTypeId         ControlTypeId = 50020 // 文本控件
	UIA_ThumbControlTypeId        ControlTypeId = 50027 // 滑块把手控件
	UIA_TitleBarControlTypeId     ControlTypeId = 50037 // 标题栏控件
	UIA_ToolBarControlTypeId      ControlTypeId = 50021 // 工具栏控件
	UIA_ToolTipControlTypeId      ControlTypeId = 50022 // 工具提示控件
	UIA_TreeControlTypeId         ControlTypeId = 50023 // 树控件
	UIA_TreeItemControlTypeId     ControlTypeId = 50024 // 树项控件
	UIA_WindowControlTypeId       ControlTypeId = 50032 // 窗口控件
)

// UIA_EventId UI Automation 事件标识符类型
// 用于标识 UI 自动化事件，如焦点变化、属性变化、结构变化等
type UIA_EventId int32

// UI Automation 事件常量定义
const (
	UIA_ActiveTextPositionChangedEventId                 UIA_EventId = 20036 // 活动文本位置变化事件
	UIA_AsyncContentLoadedEventId                        UIA_EventId = 20006 // 异步内容加载事件
	UIA_AutomationFocusChangedEventId                    UIA_EventId = 20005 // 自动化焦点变化事件
	UIA_AutomationPropertyChangedEventId                 UIA_EventId = 20004 // 自动化属性变化事件
	UIA_ChangesEventId                                   UIA_EventId = 20034 // 变化事件
	UIA_Drag_DragCancelEventId                           UIA_EventId = 20027 // 拖拽取消事件
	UIA_Drag_DragCompleteEventId                         UIA_EventId = 20028 // 拖拽完成事件
	UIA_Drag_DragStartEventId                            UIA_EventId = 20026 // 拖拽开始事件
	UIA_DropTarget_DragEnterEventId                      UIA_EventId = 20029 // 拖拽进入目标事件
	UIA_DropTarget_DragLeaveEventId                      UIA_EventId = 20030 // 拖拽离开目标事件
	UIA_DropTarget_DroppedEventId                        UIA_EventId = 20031 // 拖拽放置事件
	UIA_HostedFragmentRootsInvalidatedEventId            UIA_EventId = 20025 // 托管片段根失效事件
	UIA_InputDiscardedEventId                            UIA_EventId = 20022 // 输入被丢弃事件
	UIA_InputReachedOtherElementEventId                  UIA_EventId = 20021 // 输入到达其他元素事件
	UIA_InputReachedTargetEventId                        UIA_EventId = 20020 // 输入到达目标事件
	UIA_Invoke_InvokedEventId                            UIA_EventId = 20009 // 调用事件（按钮点击）
	UIA_LayoutInvalidatedEventId                         UIA_EventId = 20008 // 布局失效事件
	UIA_LiveRegionChangedEventId                         UIA_EventId = 20024 // 实时区域变化事件
	UIA_MenuClosedEventId                                UIA_EventId = 20007 // 菜单关闭事件
	UIA_MenuModeEndEventId                               UIA_EventId = 20019 // 菜单模式结束事件
	UIA_MenuModeStartEventId                             UIA_EventId = 20018 // 菜单模式开始事件
	UIA_MenuOpenedEventId                                UIA_EventId = 20003 // 菜单打开事件
	UIA_NotificationEventId                              UIA_EventId = 20035 // 通知事件
	UIA_Selection_InvalidatedEventId                     UIA_EventId = 20013 // 选择失效事件
	UIA_SelectionItem_ElementAddedToSelectionEventId     UIA_EventId = 20010 // 元素添加到选择事件
	UIA_SelectionItem_ElementRemovedFromSelectionEventId UIA_EventId = 20011 // 元素从选择移除事件
	UIA_SelectionItem_ElementSelectedEventId             UIA_EventId = 20012 // 元素被选择事件
	UIA_StructureChangedEventId                          UIA_EventId = 20002 // 结构变化事件
	UIA_SystemAlertEventId                               UIA_EventId = 20023 // 系统警报事件
	UIA_Text_TextChangedEventId                          UIA_EventId = 20015 // 文本变化事件
	UIA_Text_TextSelectionChangedEventId                 UIA_EventId = 20014 // 文本选择变化事件
	UIA_TextEdit_ConversionTargetChangedEventId          UIA_EventId = 20033 // 文本编辑转换目标变化事件
	UIA_TextEdit_TextChangedEventId                      UIA_EventId = 20032 // 文本编辑变化事件
	UIA_ToolTipClosedEventId                             UIA_EventId = 20001 // 工具提示关闭事件
	UIA_ToolTipOpenedEventId                             UIA_EventId = 20000 // 工具提示打开事件
	UIA_Window_WindowClosedEventId                       UIA_EventId = 20017 // 窗口关闭事件
	UIA_Window_WindowOpenedEventId                       UIA_EventId = 20016 // 窗口打开事件
)
