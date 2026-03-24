// Package uiautomation 提供 Windows UI Automation 的 Go 语言封装
package uiautomation

// PropertyConditionFlags 属性条件标志
type PropertyConditionFlags int32

// 属性条件标志常量
const (
	PropertyConditionFlags_None           PropertyConditionFlags = 0   // 无标志
	PropertyConditionFlags_IgnoreCase     PropertyConditionFlags = 0x1 // 忽略大小写
	PropertyConditionFlags_MatchSubstring PropertyConditionFlags = 0x2 // 匹配子字符串
)

// TreeScope 树搜索范围类型
type TreeScope int32

// 树搜索范围常量
// 参考: https://learn.microsoft.com/zh-cn/windows/win32/api/uiautomationclient/ne-uiautomationclient-treescope
var (
	TreeScope_None        TreeScope = 0x0                                     // 无范围
	TreeScope_Element     TreeScope = 0x1                                     // 元素本身
	TreeScope_Children    TreeScope = 0x2                                     // 直接子元素
	TreeScope_Descendants TreeScope = 0x4                                     // 所有后代元素
	TreeScope_Parent      TreeScope = 0x8                                     // 父元素
	TreeScope_Ancestors   TreeScope = 0x10                                    // 所有祖先元素
	TreeScope_Subtree     TreeScope = TreeScope_Element | TreeScope_Children | TreeScope_Descendants // 元素及其所有后代
)

// SynchronizedInputType 同步输入类型
type SynchronizedInputType int32

// 同步输入类型常量
// 参考: https://learn.microsoft.com/zh-cn/windows/win32/api/uiautomationcore/ne-uiautomationcore-synchronizedinputtype
const (
	SynchronizedInputType_KeyUp          SynchronizedInputType = 0x1  // 键盘释放
	SynchronizedInputType_KeyDown        SynchronizedInputType = 0x2  // 键盘按下
	SynchronizedInputType_LeftMouseUp    SynchronizedInputType = 0x4  // 左键释放
	SynchronizedInputType_LeftMouseDown  SynchronizedInputType = 0x8  // 左键按下
	SynchronizedInputType_RightMouseUp   SynchronizedInputType = 0x10 // 右键释放
	SynchronizedInputType_RightMouseDown SynchronizedInputType = 0x20 // 右键按下
)

// TagDvAspect DVASPECT 枚举
type TagDvAspect int32

// DVASPECT 常量
const (
	DVASPECT_CONTENT   TagDvAspect = 1 // 内容视图
	DVASPECT_THUMBNAIL TagDvAspect = 2 // 缩略图视图
	DVASPECT_ICON      TagDvAspect = 4 // 图标视图
	DVASPECT_DOCPRINT  TagDvAspect = 8 // 文档打印视图
)

// CLSCTX COM 类上下文
type CLSCTX uint32

// COM 类上下文常量
// 参考: https://learn.microsoft.com/zh-cn/windows/win32/api/wtypesbase/ne-wtypesbase-clsctx
var (
	CLSCTX_INPROC_SERVER                   CLSCTX = 0x1       // 进程内服务器
	CLSCTX_INPROC_HANDLER                  CLSCTX = 0x2       // 进程内处理器
	CLSCTX_LOCAL_SERVER                    CLSCTX = 0x4       // 本地服务器
	CLSCTX_INPROC_SERVER16                 CLSCTX = 0x8       // 16位进程内服务器
	CLSCTX_REMOTE_SERVER                   CLSCTX = 0x10      // 远程服务器
	CLSCTX_INPROC_HANDLER16                CLSCTX = 0x20      // 16位进程内处理器
	CLSCTX_RESERVED1                       CLSCTX = 0x40      // 保留
	CLSCTX_RESERVED2                       CLSCTX = 0x80      // 保留
	CLSCTX_RESERVED3                       CLSCTX = 0x100     // 保留
	CLSCTX_RESERVED4                       CLSCTX = 0x200     // 保留
	CLSCTX_NO_CODE_DOWNLOAD                CLSCTX = 0x400     // 禁止代码下载
	CLSCTX_RESERVED5                       CLSCTX = 0x800     // 保留
	CLSCTX_NO_CUSTOM_MARSHAL               CLSCTX = 0x1000    // 禁止自定义封送
	CLSCTX_ENABLE_CODE_DOWNLOAD            CLSCTX = 0x2000    // 启用代码下载
	CLSCTX_NO_FAILURE_LOG                  CLSCTX = 0x4000    // 禁止失败日志
	CLSCTX_DISABLE_AAA                     CLSCTX = 0x8000    // 禁用 AAA
	CLSCTX_ENABLE_AAA                      CLSCTX = 0x10000   // 启用 AAA
	CLSCTX_FROM_DEFAULT_CONTEXT            CLSCTX = 0x20000   // 从默认上下文
	CLSCTX_ACTIVATE_X86_SERVER             CLSCTX = 0x40000   // 激活 x86 服务器
	_CLSCTX_ACTIVATE_32_BIT_SERVER         CLSCTX = 0         // 激活 32 位服务器
	CLSCTX_ACTIVATE_64_BIT_SERVER          CLSCTX = 0x80000   // 激活 64 位服务器
	CLSCTX_ENABLE_CLOAKING                 CLSCTX = 0x100000  // 启用伪装
	CLSCTX_APPCONTAINER                    CLSCTX = 0x400000  // 应用容器
	CLSCTX_ACTIVATE_AAA_AS_IU              CLSCTX = 0x800000  // 作为 IU 激活 AAA
	CLSCTX_RESERVED6                       CLSCTX = 0x1000000 // 保留
	CLSCTX_ACTIVATE_ARM32_SERVER           CLSCTX = 0x2000000 // 激活 ARM32 服务器
	_CLSCTX_ALLOW_LOWER_TRUST_REGISTRATION CLSCTX = 0         // 允许低信任注册
	CLSCTX_PS_DLL                          CLSCTX = 0x80000000 // PS DLL
	CLSCTX_ALL                             CLSCTX = CLSCTX_INPROC_SERVER | CLSCTX_INPROC_HANDLER | CLSCTX_LOCAL_SERVER | CLSCTX_REMOTE_SERVER // 所有上下文
)

// SelFlag 选择标志
type SelFlag int32

// 选择标志常量
// 参考: https://learn.microsoft.com/zh-cn/windows/win32/winauto/selflag
var (
	SELFLAG_NONE            SelFlag = 0    // 无标志
	SELFLAG_TAKEFOCUS       SelFlag = 0x1  // 获取焦点
	SELFLAG_TAKESELECTION   SelFlag = 0x2  // 获取选择
	SELFLAG_EXTENDSELECTION SelFlag = 0x4  // 扩展选择
	SELFLAG_ADDSELECTION    SelFlag = 0x8  // 添加到选择
	SELFLAG_REMOVESELECTION SelFlag = 0x10 // 从选择移除
)

// TagInvokeKind 调用类型
type TagInvokeKind int32

// 调用类型常量
const (
	INVOKE_FUNC           TagInvokeKind = 1 // 函数调用
	INVOKE_PROPERTYGET    TagInvokeKind = 2 // 属性获取
	INVOKE_PROPERTYPUT    TagInvokeKind = 4 // 属性设置
	INVOKE_PROPERTYPUTREF TagInvokeKind = 8 // 属性引用设置
)

// ParamFlag 参数标志
type ParamFlag uint16

// 参数标志常量
// 参考: https://learn.microsoft.com/zh-cn/previous-versions/windows/desktop/automat/paramflags
const (
	PARAMFLAG_NONE         ParamFlag = 0    // 无标志
	PARAMFLAG_FIN          ParamFlag = 0x1  // 输入参数
	PARAMFLAG_FOUT         ParamFlag = 0x2  // 输出参数
	PARAMFLAG_FLCID        ParamFlag = 0x4  // LCID 参数
	PARAMFLAG_FRETVAL      ParamFlag = 0x8  // 返回值参数
	PARAMFLAG_FOPT         ParamFlag = 0x10 // 可选参数
	PARAMFLAG_FHASDEFAULT  ParamFlag = 0x20 // 有默认值
	PARAMFLAG_FHASCUSTDATA ParamFlag = 0x40 // 有自定义数据
)
