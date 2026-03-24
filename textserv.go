// Package uiautomation 提供 Windows UI Automation 的 Go 语言封装
package uiautomation

import (
	"syscall"
	"unsafe"
)

var (
	// msftedit Microsoft Rich Edit 控件库
	msftedit = syscall.NewLazyDLL("Msftedit.dll")
	// procCreateTextServices 创建文本服务函数
	procCreateTextServices = msftedit.NewProc("CreateTextServices")
	// procShutdownTextServices 关闭文本服务函数
	procShutdownTextServices = msftedit.NewProc("ShutdownTextServices")
)

// CreateTextServices 创建文本服务
// 参数:
//   - unk: IUnKnown 接口
//   - thost: ITextHost 接口
//
// 返回: 文本服务接口和可能的错误
// 注意: ITextHost 类型未定义，此函数暂时注释
/*
func CreateTextServices(unk *IUnKnown, thost *ITextHost) (*IUnKnown, error) {
	return createTextServices(unk, thost)
}

// ShutdownTextServices 关闭文本服务
// 参数: unk - IUnKnown 接口
// 返回: 错误信息
func ShutdownTextServices(unk *IUnKnown) error {
	return shutdownTextServices(unk)
}

// createTextServices 内部函数：创建文本服务
// 参数:
//   - in: IUnKnown 接口
//   - in2: ITextHost 接口
//
// 返回: 文本服务接口和可能的错误
func createTextServices(in *IUnKnown, in2 *ITextHost) (*IUnKnown, error) {
	var retVal *IUnKnown
	ret, _, _ := procCreateTextServices.Call(
		uintptr(unsafe.Pointer(in)),
		uintptr(unsafe.Pointer(in2)),
		uintptr(unsafe.Pointer(&retVal)),
	)
	if ret != 0 {
		return nil, HResult(ret)
	}
	return retVal, nil
}
*/

// shutdownTextServices 内部函数：关闭文本服务
// 参数: in - IUnKnown 接口
// 返回: 错误信息
func shutdownTextServices(in *IUnKnown) error {
	ret, _, _ := procShutdownTextServices.Call(uintptr(unsafe.Pointer(in)))
	if ret != 0 {
		return HResult(ret)
	}
	return nil
}

// IRichEditUiaInformation Rich Edit UI Automation 信息接口
type IRichEditUiaInformation struct {
	vtbl *IUnKnown
}

// IRichEditUiaInformationVtbl Rich Edit UI Automation 信息接口虚函数表
type IRichEditUiaInformationVtbl struct {
	IUnKnownVtbl
	GetBoundaryRectangle uintptr // 获取边界矩形
	IsVisible            uintptr // 是否可见
}