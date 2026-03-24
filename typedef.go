// Package uiautomation 提供 Windows UI Automation 的 Go 语言封装
package uiautomation

// TagPoint 表示一个点的坐标（整数）
type TagPoint struct {
	X int32 // X 坐标
	Y int32 // Y 坐标
}

// TagRect 表示一个矩形区域
type TagRect struct {
	Left   int32 // 左边界
	Top    int32 // 上边界
	Right  int32 // 右边界
	Bottom int32 // 下边界
}

// TagSafeArrayBound SAFEARRAY 边界描述
// 参考: https://learn.microsoft.com/zh-cn/windows/win32/api/oaidl/ns-oaidl-safearraybound
type TagSafeArrayBound struct {
	CElements uint32 // 元素数量
	LLbound   int32  // 下界
}

// TagSafeArray SAFEARRAY 结构
// 用于 COM 中的数组传递
// 参考: https://learn.microsoft.com/zh-cn/windows/win32/api/oaidl/ns-oaidl-safearray
type TagSafeArray struct {
	CDims     uint16              // 维度数量
	FFeatures uint16              // 特征标志
	CbElement uint32              // 元素大小（字节）
	CLocks    uint32              // 锁定计数
	PvData    uintptr             // 数据指针
	Rgsabound [1]TagSafeArrayBound // 边界数组（COM 中通常是尾随数组）
}

// UiaPoint UI Automation 点结构（浮点数）
// 参考: https://learn.microsoft.com/zh-cn/windows/win32/api/uiautomationcore/ns-uiautomationcore-uiapoint
type UiaPoint struct {
	X float64 // X 坐标
	Y float64 // Y 坐标
}

// TagDvTargetDevice 目标设备描述
type TagDvTargetDevice struct {
	TdSize             uint32   // 结构大小
	TdDriverNameOffset uint16   // 驱动名称偏移
	TdDeviceNameOffset uint16   // 设备名称偏移
	TdPortNameOffset   uint16   // 端口名称偏移
	TdExtDevmodeOffset uint16   // 扩展设备模式偏移
	TdData             [1]byte  // 数据
}

// UiaRect UI Automation 矩形结构（浮点数）
type UiaRect struct {
	Left   float64 // 左边界
	Top    float64 // 上边界
	Width  float64 // 宽度
	Height float64 // 高度
}

// TagDispParams DISPPARAMS 结构
// 用于 IDispatch::Invoke 方法的参数传递
type TagDispParams struct {
	Rgvarg            uintptr // 指向 VARIANT 数组
	RgdispidNamedArgs uintptr // 指向 DISPID 数组
	Cargs             uint32  // 参数数量
	CNamedArgs        uint32  // 命名参数数量
}

// TagExcepInfo EXCEPINFO 结构
// 用于存储异常信息
type TagExcepInfo struct {
	WCode             uint16  // 错误代码
	WReserved         uint16  // 保留
	BstrSource        uintptr // 错误源
	BstrDescription   uintptr // 错误描述
	BstrHelpFile      uintptr // 帮助文件
	DwHelpContext     uint32  // 帮助上下文
	PvReserved        uintptr // 保留
	PFnDeferredFillIn uintptr // 延迟填充函数
	Scode             int32   // SCODE
}

// TagSize 表示尺寸（宽度和高度）
type TagSize struct {
	Cx int32 // 宽度
	Cy int32 // 高度
}

// TYPEDESC 类型描述（未完成）
type TYPEDESC struct {
}

// TagIdlDesc IDL 描述
type TagIdlDesc struct {
	DwReserved uintptr // ULONG_PTR
	WIdlFlags  uint16  // IDL 标志
}

// TagParamDesc 参数描述
type TagParamDesc struct {
	Pparamdescex uintptr // LPPARAMDESCEX
	WParamFlags  uint16  // PARAMFLAGS
}

// TagTypeDesc 类型描述
type TagTypeDesc struct {
}

// TagElemDesc 元素描述
type TagElemDesc struct {
	Tdesc     TagTypeDesc  // 类型描述
	Paramdesc TagParamDesc // 参数描述
}

// TagFuncDesc 函数描述
type TagFuncDesc struct {
	Memid             int32        // MEMBERID
	LPrgsCode         uintptr      // SCODE*
	LPrgelemdescParam uintptr      // ELEMDESC*
	Funckind          int32        // FUNCKIND
	Invkind           int32        // INVOKEKIND
	Callconv          int32        // CALLCONV
	CParams           int16        // 参数数量
	CParamsOpt        int16        // 可选参数数量
	OVft              int16        // 偏移
	CScodes           int16        // SCODE 数量
	ElemdescFunc      TagElemDesc  // 函数元素描述
	WFuncFlags        uint16       // 函数标志
}
