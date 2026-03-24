// Package uiautomation 提供 Windows UI Automation 的 Go 语言封装
package uiautomation

import (
	"errors"
	"runtime"
	"syscall"
	"unsafe"
)

// Windows API DLL 和函数声明
var (
	user32 = syscall.NewLazyDLL("user32.dll") // Windows 用户界面 API
	ole32  = syscall.NewLazyDLL("Ole32.dll")  // Windows COM/OLE API

	procCoCreateInstance = ole32.NewProc("CoCreateInstance") // 创建 COM 对象实例
	procCoInitialize     = ole32.NewProc("CoInitialize")     // 初始化 COM 库
	procCoUninitialize   = ole32.NewProc("CoUninitialize")   // 释放 COM 库
	procFindWindowW      = user32.NewProc("FindWindowW")     // 查找顶层窗口
	procFindWindowExW    = user32.NewProc("FindWindowExW")   // 查找子窗口

	// ErrorNotFoundWindow 窗口未找到错误
	ErrorNotFoundWindow = errors.New("not found window")
)

// GetWindowForString 通过类名或窗口名查找窗口
// 参数:
//   - classname: 窗口类名（可为空）
//   - windowname: 窗口标题名（可为空）
// 返回:
//   - uintptr: 窗口句柄
//   - error: 未找到窗口时返回错误
func GetWindowForString(classname, windowname string) (uintptr, error) {
	find := findWindowW(classname, windowname)
	if find == 0 {
		return 0, ErrorNotFoundWindow
	}
	return find, nil
}

// CoInitialize 初始化 COM 库
// 在使用任何 COM 对象之前必须调用此函数
// 返回: 初始化失败时返回错误
// 注意: 会锁定 OS 线程，需要配合 CoUninitialize 释放
func CoInitialize() error {
	runtime.LockOSThread()
	ret, _, _ := procCoInitialize.Call(
		uintptr(0),
	)
	// CoInitialize 返回 S_FALSE (1) 表示 COM 已在此线程初始化
	// 这仍然是一个成功的结果
	if ret != 0 && ret != 1 {
		runtime.UnlockOSThread()
		return HResult(ret)
	}
	return nil
}

// CoUninitialize 释放 COM 库
// 在程序结束时调用，与 CoInitialize 配对使用
// 注意: 会解锁 OS 线程
func CoUninitialize() {
	procCoUninitialize.Call()
	runtime.UnlockOSThread()
}

// CreateInstance 创建 COM 对象实例
// 参数:
//   - clsid: 类标识符（CLSID）
//   - riid: 接口标识符（IID）
//   - clsctx: 执行上下文
// 返回:
//   - unsafe.Pointer: COM 对象指针
//   - error: 创建失败时返回错误
func CreateInstance(clsid, riid *syscall.GUID, clsctx CLSCTX) (unsafe.Pointer, error) {
	return createInstance(clsid, riid, clsctx)
}

// FindWindowW 查找顶层窗口
// 参数:
//   - lpclass: 窗口类名
//   - lpwindow: 窗口标题名
// 返回: 窗口句柄，未找到返回 0
func FindWindowW(lpclass, lpwindow string) uintptr {
	return findWindowW(lpclass, lpwindow)
}

// FindWindowExW 查找子窗口
// 参数:
//   - phwdn: 父窗口句柄
//   - chwdn: 子窗口句柄（从该窗口之后开始查找）
//   - lpclass: 窗口类名
//   - lpwindow: 窗口标题名
// 返回: 窗口句柄，未找到返回 0
func FindWindowExW(phwdn, chwdn uintptr, lpclass, lpwindow string) uintptr {
	return findWindowExW(phwdn, chwdn, lpclass, lpwindow)
}

// createInstance 内部函数：调用 CoCreateInstance 创建 COM 对象
func createInstance(clsid, riid *syscall.GUID, clsctx CLSCTX) (unsafe.Pointer, error) {
	var retVal unsafe.Pointer
	ret, _, _ := procCoCreateInstance.Call(
		uintptr(unsafe.Pointer(clsid)),
		uintptr(0),
		uintptr(clsctx),
		uintptr(unsafe.Pointer(riid)),
		uintptr(unsafe.Pointer(&retVal)),
	)
	if ret != 0 {
		return nil, HResult(ret)
	}
	return retVal, nil
}

// findWindowW 内部函数：调用 FindWindowW API 查找窗口
func findWindowW(lpclass, lpwindow string) uintptr {
	var lpclassname, lpwindowname *uint16
	var err error
	if lpclass != "" {
		lpclassname, err = syscall.UTF16PtrFromString(lpclass)
		if err != nil {
			return 0
		}
	}
	if lpwindow != "" {
		lpwindowname, err = syscall.UTF16PtrFromString(lpwindow)
		if err != nil {
			return 0
		}
	}
	ret, _, _ := procFindWindowW.Call(
		uintptr(unsafe.Pointer(lpclassname)),
		uintptr(unsafe.Pointer(lpwindowname)),
	)
	return ret
}

// findWindowExW 内部函数：调用 FindWindowExW API 查找子窗口
func findWindowExW(phwdn, chwdn uintptr, lpclass, lpwindow string) uintptr {
	var lpclassname, lpwindowname *uint16
	var err error
	if lpclass != "" {
		lpclassname, err = syscall.UTF16PtrFromString(lpclass)
		if err != nil {
			return 0
		}
	}
	if lpwindow != "" {
		lpwindowname, err = syscall.UTF16PtrFromString(lpwindow)
		if err != nil {
			return 0
		}
	}
	ret, _, _ := procFindWindowExW.Call(
		phwdn,
		chwdn,
		uintptr(unsafe.Pointer(lpclassname)),
		uintptr(unsafe.Pointer(lpwindowname)),
	)
	return ret
}
