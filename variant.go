// Package uiautomation 提供 Windows UI Automation 的 Go 语言封装
package uiautomation

// NewVariant 创建新的 VARIANT 结构
// 参数:
//   - vt: 变体类型
//   - val: 变体值
//
// 返回: 初始化的 VARIANT 结构
func NewVariant(vt TagVarenum, val int64) VARIANT {
	return VARIANT{
		VT:  vt,
		Val: val,
	}
}

// VariantFromString 从字符串创建 BSTR 类型的 VARIANT
// 参数: s - 字符串值
// 返回: VARIANT 结构和可能的错误
func VariantFromString(s string) (VARIANT, error) {
	bstr, err := string2Bstr(s)
	if err != nil {
		return VARIANT{}, err
	}
	return VARIANT{
		VT:  VT_BSTR,
		Val: int64(bstr),
	}, nil
}

// TagVarenum VARIANT 类型枚举
// 定义了 COM VARIANT 结构支持的数据类型
type TagVarenum uint16

// VARIANT 类型常量定义
// 参考: https://learn.microsoft.com/zh-cn/windows/win32/api/wtypes/ne-wtypes-varenum
const (
	VT_EMPTY            TagVarenum = 0      // 空类型
	VT_NULL             TagVarenum = 1      // NULL 值
	VT_I2               TagVarenum = 2      // 16位有符号整数
	VT_I4               TagVarenum = 3      // 32位有符号整数
	VT_R4               TagVarenum = 4      // 32位浮点数
	VT_R8               TagVarenum = 5      // 64位浮点数
	VT_CY               TagVarenum = 6      // 货币类型
	VT_DATE             TagVarenum = 7      // 日期类型
	VT_BSTR             TagVarenum = 8      // BSTR 字符串
	VT_DISPATCH         TagVarenum = 9      // IDispatch 接口
	VT_ERROR            TagVarenum = 10     // SCODE 错误码
	VT_BOOL             TagVarenum = 11     // 布尔值
	VT_VARIANT          TagVarenum = 12     // VARIANT 类型
	VT_UNKNOWN          TagVarenum = 13     // IUnknown 接口
	VT_DECIMAL          TagVarenum = 14     // 十进制类型
	VT_I1               TagVarenum = 16     // 8位有符号整数
	VT_UI1              TagVarenum = 17     // 8位无符号整数
	VT_UI2              TagVarenum = 18     // 16位无符号整数
	VT_UI4              TagVarenum = 19     // 32位无符号整数
	VT_I8               TagVarenum = 20     // 64位有符号整数
	VT_UI8              TagVarenum = 21     // 64位无符号整数
	VT_INT              TagVarenum = 22     // 有符号整数（平台相关）
	VT_UINT             TagVarenum = 23     // 无符号整数（平台相关）
	VT_VOID             TagVarenum = 24     // void 类型
	VT_HRESULT          TagVarenum = 25     // HRESULT 类型
	VT_PTR              TagVarenum = 26     // 指针类型
	VT_SAFEARRAY        TagVarenum = 27     // SAFEARRAY 类型
	VT_CARRAY           TagVarenum = 28     // C 风格数组
	VT_USERDEFINED      TagVarenum = 29     // 用户定义类型
	VT_LPSTR            TagVarenum = 30     // 以 null 结尾的字符串
	VT_LPWSTR           TagVarenum = 31     // 以 null 结尾的宽字符串
	VT_RECORD           TagVarenum = 36     // 记录类型
	VT_INT_PTR          TagVarenum = 37     // 整数指针
	VT_UINT_PTR         TagVarenum = 38     // 无符号整数指针
	VT_FILETIME         TagVarenum = 64     // FILETIME 类型
	VT_BLOB             TagVarenum = 65     // BLOB 类型
	VT_STREAM           TagVarenum = 66     // IStream 接口
	VT_STORAGE          TagVarenum = 67     // IStorage 接口
	VT_STREAMED_OBJECT  TagVarenum = 68     // 流式对象
	VT_STORED_OBJECT    TagVarenum = 69     // 存储对象
	VT_BLOB_OBJECT      TagVarenum = 70     // BLOB 对象
	VT_CF               TagVarenum = 71     // 剪贴板格式
	VT_CLSID            TagVarenum = 72     // CLSID 类型
	VT_VERSIONED_STREAM TagVarenum = 73     // 版本化流
	VT_BSTR_BLOB        TagVarenum = 0xfff  // BSTR BLOB 类型
	VT_VECTOR           TagVarenum = 0x1000 // 向量标志
	VT_ARRAY            TagVarenum = 0x2000 // 数组标志
	VT_BYREF            TagVarenum = 0x4000 // 引用标志
	VT_RESERVED         TagVarenum = 0x8000 // 保留标志
	VT_ILLEGAL          TagVarenum = 0xffff // 非法类型
	VT_ILLEGALMASKED    TagVarenum = 0xfff  // 非法类型掩码
	VT_TYPEMASK         TagVarenum = 0xfff  // 类型掩码
)
