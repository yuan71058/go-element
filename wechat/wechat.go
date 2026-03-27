// Package wechat 提供微信自动化操作功能
package wechat

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"
	"time"
	"unsafe"

	uia "github.com/yuan71058/go-element"
)

// 微信窗口相关常量
const (
	// WeChatWindowTitle 微信窗口标题
	WeChatWindowTitle = "微信"
	// WeChatClassName 微信窗口类名
	WeChatClassName = "Qt51514QWindowIcon"
)

// Windows API 常量
const (
	WM_KEYDOWN = 0x0100 // 按键按下消息
	WM_KEYUP   = 0x0101 // 按键抬起消息
	VK_CONTROL = 0x11   // Ctrl 键虚拟键码
	VK_V       = 0x56   // V 键虚拟键码
)

// Windows API DLL
var (
	user32   = syscall.NewLazyDLL("user32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
)

// Bot 微信机器人结构体
// 封装了微信自动化操作的所有功能
type Bot struct {
	ppv            *uia.IUIAutomation        // UI自动化接口
	hwnd           uintptr                   // 微信窗口句柄
	root           *uia.IUIAutomationElement // 微信窗口根元素
	logFile        *os.File                  // 日志文件句柄
	sessionManager *SessionManager           // 会话管理器
	wechatId       string                    // 微信号
}

// NewBot 创建微信机器人实例
// 这是使用本包的入口函数，必须首先调用
//
// 返回值:
//   - *Bot: 初始化好的机器人实例
//   - error: 错误信息，如微信未打开则返回错误
//
// 使用示例:
//
//	bot, err := wechat.NewBot()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer bot.Close()
func NewBot() (*Bot, error) {
	return NewBotWithLog("")
}

// NewBotWithLog 创建微信机器人实例（带日志文件配置）
// 参数:
//   - logPath: 日志文件路径，为空则不记录日志
//
// 返回值:
//   - *Bot: 初始化好的机器人实例
//   - error: 错误信息
func NewBotWithLog(logPath string) (*Bot, error) {
	var f *os.File
	var err error

	if logPath != "" {
		f, err = os.OpenFile(logPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("无法创建日志文件: %v", err)
		}
		log.SetOutput(f)
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	}
	log.Println("--- 微信自动化启动 ---")

	// 初始化 COM 组件
	err = uia.CoInitialize()
	if err != nil {
		return nil, fmt.Errorf("COM 初始化失败: %v", err)
	}

	// 查找微信窗口
	log.Println("正在查找微信窗口...")
	hwnd, err := uia.GetWindowForString("", WeChatWindowTitle)
	if err != nil || hwnd == 0 {
		hwnd, _ = uia.GetWindowForString(WeChatClassName, "")
	}
	if hwnd == 0 {
		uia.CoUninitialize()
		return nil, fmt.Errorf("未找到微信窗口，请确保微信已打开")
	}
	log.Printf("找到微信窗口，句柄: 0x%X", hwnd)

	// 激活微信窗口（已禁用）
	// log.Println("激活微信窗口...")
	// user32.NewProc("SetForegroundWindow").Call(hwnd)
	// time.Sleep(500 * time.Millisecond)

	// 创建 UIAutomation 实例
	instance, err := uia.CreateInstance(uia.CLSID_CUIAutomation, uia.IID_IUIAutomation, uia.CLSCTX_ALL)
	if err != nil {
		uia.CoUninitialize()
		return nil, fmt.Errorf("无法创建 UIAutomation: %v", err)
	}
	ppv := uia.NewIUIAutomation(uia.NewIUnKnown(instance))

	// 从窗口句柄获取根元素
	root, err := uia.ElementFromHandle(ppv, hwnd)
	if err != nil {
		ppv.Release()
		uia.CoUninitialize()
		return nil, fmt.Errorf("无法获取窗口元素: %v", err)
	}

	return &Bot{
		ppv:            ppv,
		hwnd:           hwnd,
		root:           root,
		logFile:        f,
		sessionManager: NewSessionManager(),
	}, nil
}

// Close 关闭机器人，释放所有资源
// 使用完毕后必须调用此函数释放资源
//
// 使用示例:
//
//	defer bot.Close()
func (bot *Bot) Close() {
	if bot.root != nil {
		bot.root.Release()
	}
	if bot.ppv != nil {
		bot.ppv.Release()
	}
	uia.CoUninitialize()
	if bot.logFile != nil {
		bot.logFile.Close()
	}
}

// ==================== 元素查找方法 ====================

// FindElementByName 通过名称查找UI元素
// 参数:
//   - name: 元素名称
//
// 返回: 找到的元素，未找到返回 nil
func (bot *Bot) FindElementByName(name string) *uia.IUIAutomationElement {
	variant, _ := uia.VariantFromString(name)
	cond, _ := bot.ppv.CreatePropertyCondition(uia.UIA_NamePropertyId, variant)
	elem, _ := bot.root.FindFirst(uia.TreeScope_Descendants, cond)
	return elem
}

// FindElementByAutomationId 通过 AutomationId 查找UI元素
// 参数:
//   - id: 元素的 AutomationId
//
// 返回: 找到的元素，未找到返回 nil
func (bot *Bot) FindElementByAutomationId(id string) *uia.IUIAutomationElement {
	variant, _ := uia.VariantFromString(id)
	cond, _ := bot.ppv.CreatePropertyCondition(uia.UIA_AutomationIdPropertyId, variant)
	elem, _ := bot.root.FindFirst(uia.TreeScope_Descendants, cond)
	return elem
}

// FindElementByClassName 通过类名查找UI元素
// 参数:
//   - className: 元素的类名
//
// 返回: 找到的元素，未找到返回 nil
func (bot *Bot) FindElementByClassName(className string) *uia.IUIAutomationElement {
	variant, _ := uia.VariantFromString(className)
	cond, _ := bot.ppv.CreatePropertyCondition(uia.UIA_ClassNamePropertyId, variant)
	elem, _ := bot.root.FindFirst(uia.TreeScope_Descendants, cond)
	return elem
}

// ==================== 基础操作方法 ====================

// ClickElement 使用 BoundingRectangle 点击指定的UI元素（模拟鼠标点击）
// 参数:
//   - elem: 要点击的元素
//
// 返回: 错误信息，成功返回 nil
func (bot *Bot) ClickElement(elem *uia.IUIAutomationElement) error {
	if elem == nil {
		return fmt.Errorf("元素为空")
	}

	log.Println("使用BoundingRectangle点击")
	boundingRect := elem.Get_CurrentBoundingRectangle()
	if boundingRect == nil {
		return fmt.Errorf("无法获取元素位置")
	}

	clickX := (boundingRect.Left + boundingRect.Right) / 2
	clickY := (boundingRect.Top + boundingRect.Bottom) / 2
	log.Printf("点击位置: X=%d, Y=%d", clickX, clickY)

	user32.NewProc("SetCursorPos").Call(uintptr(clickX), uintptr(clickY))
	time.Sleep(300 * time.Millisecond)
	user32.NewProc("mouse_event").Call(0x0002, 0, 0, 0, 0)
	time.Sleep(300 * time.Millisecond)
	user32.NewProc("mouse_event").Call(0x0004, 0, 0, 0, 0)

	return nil
}

// ClickElementWithInvoke 使用 InvokePattern 或 SelectionItemPattern 点击指定的UI元素
// 参数:
//   - elem: 要点击的元素
//
// 返回: 错误信息，成功返回 nil
func (bot *Bot) ClickElementWithInvoke(elem *uia.IUIAutomationElement) error {
	if elem == nil {
		return fmt.Errorf("元素为空")
	}

	isEnabled := elem.Get_CurrentIsEnabled()
	isOffscreen := elem.Get_CurrentIsOffscreen()
	if isEnabled == 0 {
		return fmt.Errorf("元素未启用")
	}
	if isOffscreen != 0 {
		return fmt.Errorf("元素不在屏幕上")
	}

	name, _ := elem.Get_CurrentName()
	className, _ := elem.Get_CurrentClassName()
	automationId, _ := elem.Get_CurrentAutomationId()
	controlType := elem.Get_CurrentControlType()
	boundingRect := elem.Get_CurrentBoundingRectangle()
	log.Printf("[点击前] 元素信息 - Name: %q, ClassName: %q, AutomationId: %q, ControlType: %d", name, className, automationId, controlType)
	if boundingRect != nil {
		log.Printf("[点击前] 元素坐标 - Left: %d, Top: %d, Right: %d, Bottom: %d", boundingRect.Left, boundingRect.Top, boundingRect.Right, boundingRect.Bottom)
	}

	if controlType == uia.UIA_ButtonControlTypeId {
		log.Println("按钮控件，先使用 InvokePattern 单击")
		unk, err := elem.GetCurrentPattern(uia.UIA_InvokePatternId)
		if err != nil || unk == nil {
			return fmt.Errorf("无法获取 InvokePattern: %v", err)
		}
		ip := uia.NewIUIAutomationInvokePattern(unk)
		defer ip.Release()

		log.Println("调用 Invoke() 方法...")
		err = ip.Invoke()
		if err != nil {
			return fmt.Errorf("调用 Invoke 失败: %v", err)
		}
		log.Println("Invoke调用成功")

		time.Sleep(100 * time.Millisecond)

		log.Println("再使用 LegacyIAccessiblePattern 单击")
		unk, err = elem.GetCurrentPattern(uia.UIA_LegacyIAccessiblePatternId)
		if err != nil || unk == nil {
			return fmt.Errorf("无法获取 LegacyIAccessiblePattern: %v", err)
		}
		liap := uia.NewIUIAutomationLegacyIAccessiblePattern(unk)
		defer liap.Release()

		defaultAction, _ := liap.Get_CurrentDefaultAction()
		log.Printf("默认操作: %q", defaultAction)

		log.Println("调用 DoDefaultAction() 方法...")
		err = liap.DoDefaultAction()
		if err != nil {
			return fmt.Errorf("调用 DoDefaultAction 失败: %v", err)
		}
		log.Println("DoDefaultAction调用成功")
	} else if controlType == uia.UIA_ListItemControlTypeId {
		log.Println("列表项控件，使用 LegacyIAccessiblePattern 双击")
		unk, err := elem.GetCurrentPattern(uia.UIA_LegacyIAccessiblePatternId)
		if err != nil || unk == nil {
			return fmt.Errorf("无法获取 LegacyIAccessiblePattern: %v", err)
		}
		liap := uia.NewIUIAutomationLegacyIAccessiblePattern(unk)
		defer liap.Release()

		defaultAction, _ := liap.Get_CurrentDefaultAction()
		log.Printf("默认操作: %q", defaultAction)

		log.Println("调用 DoDefaultAction() 方法...")
		err = liap.DoDefaultAction()
		if err != nil {
			return fmt.Errorf("调用 DoDefaultAction 失败: %v", err)
		}
		log.Println("DoDefaultAction调用成功")
	} else if controlType == uia.UIA_EditControlTypeId {
		log.Println("编辑控件，使用 InvokePattern 获得焦点")
		unk, err := elem.GetCurrentPattern(uia.UIA_InvokePatternId)
		if err != nil || unk == nil {
			return fmt.Errorf("无法获取 InvokePattern: %v", err)
		}
		ip := uia.NewIUIAutomationInvokePattern(unk)
		defer ip.Release()

		log.Println("调用 Invoke() 方法...")
		err = ip.Invoke()
		if err != nil {
			return fmt.Errorf("调用 Invoke 失败: %v", err)
		}
		log.Println("Invoke调用成功")
	} else {
		log.Println("其他控件，尝试 LegacyIAccessiblePattern")
		unk, err := elem.GetCurrentPattern(uia.UIA_LegacyIAccessiblePatternId)
		if err != nil || unk == nil {
			return fmt.Errorf("无法获取 LegacyIAccessiblePattern: %v", err)
		}
		liap := uia.NewIUIAutomationLegacyIAccessiblePattern(unk)
		defer liap.Release()

		defaultAction, _ := liap.Get_CurrentDefaultAction()
		log.Printf("默认操作: %q", defaultAction)

		log.Println("调用 DoDefaultAction() 方法...")
		err = liap.DoDefaultAction()
		if err != nil {
			return fmt.Errorf("调用 DoDefaultAction 失败: %v", err)
		}
		log.Println("DoDefaultAction调用成功")
	}

	time.Sleep(500 * time.Millisecond)

	nameAfter, _ := elem.Get_CurrentName()
	classNameAfter, _ := elem.Get_CurrentClassName()
	automationIdAfter, _ := elem.Get_CurrentAutomationId()
	boundingRectAfter := elem.Get_CurrentBoundingRectangle()
	log.Printf("[点击后] 元素信息 - Name: %q, ClassName: %q, AutomationId: %q", nameAfter, classNameAfter, automationIdAfter)
	if boundingRectAfter != nil {
		log.Printf("[点击后] 元素坐标 - Left: %d, Top: %d, Right: %d, Bottom: %d", boundingRectAfter.Left, boundingRectAfter.Top, boundingRectAfter.Right, boundingRectAfter.Bottom)
	}

	return nil
}

// ClickElementWithDoubleClick 使用 InvokePattern 双击指定的UI元素（不占用鼠标）
// 参数:
//   - elem: 要双击的元素
//
// 返回: 错误信息，成功返回 nil
func (bot *Bot) ClickElementWithDoubleClick(elem *uia.IUIAutomationElement) error {
	if elem == nil {
		return fmt.Errorf("元素为空")
	}

	log.Println("使用InvokePattern双击")

	name, _ := elem.Get_CurrentName()
	className, _ := elem.Get_CurrentClassName()
	automationId, _ := elem.Get_CurrentAutomationId()
	log.Printf("双击元素信息 - Name: %q, ClassName: %q, AutomationId: %q", name, className, automationId)

	unk, err := elem.GetCurrentPattern(uia.UIA_InvokePatternId)
	if err != nil {
		return fmt.Errorf("无法获取 InvokePattern: %v", err)
	}

	if unk == nil {
		return fmt.Errorf("元素不支持 InvokePattern")
	}

	ip := uia.NewIUIAutomationInvokePattern(unk)
	defer ip.Release()

	log.Println("第一次调用 Invoke()...")
	err = ip.Invoke()
	if err != nil {
		return fmt.Errorf("第一次调用 Invoke 失败: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	log.Println("第二次调用 Invoke()...")
	err = ip.Invoke()
	if err != nil {
		return fmt.Errorf("第二次调用 Invoke 失败: %v", err)
	}

	log.Println("双击Invoke调用成功")
	return nil
}

// SetTextValue 使用 ValuePattern 设置元素的文本值
// 参数:
//   - elem: 目标元素
//   - text: 要设置的文本
//
// 返回: 错误信息，成功返回 nil
func (bot *Bot) SetTextValue(elem *uia.IUIAutomationElement, text string) error {
	if elem == nil {
		return fmt.Errorf("元素为空")
	}

	unk, err := elem.GetCurrentPattern(uia.UIA_ValuePatternId)
	if err != nil {
		return fmt.Errorf("无法获取 ValuePattern: %v", err)
	}

	if unk != nil {
		vp := uia.NewIUIAutomationValuePattern(unk)
		defer vp.Release()
		log.Printf("使用 ValuePattern 设置文本: %q", text)
		err := vp.SetValue(text)
		if err != nil {
			return err
		}
		log.Println("文本设置成功")
		return nil
	}
	return fmt.Errorf("元素不支持 ValuePattern")
}

// GetTextValue 获取元素的文本值
// 参数:
//   - elem: 目标元素
//
// 返回: 元素的文本值
func (bot *Bot) GetTextValue(elem *uia.IUIAutomationElement) string {
	if elem == nil {
		return ""
	}
	defer elem.Release()

	unk, err := elem.GetCurrentPattern(uia.UIA_ValuePatternId)
	if err == nil && unk != nil {
		vp := uia.NewIUIAutomationValuePattern(unk)
		defer vp.Release()
		val, err := vp.Get_CurrentValue()
		if err == nil {
			return val
		}
	}

	name, _ := elem.Get_CurrentName()
	return name
}

// ==================== 文本输入方法 ====================

// TypeText 模拟键盘输入文本
// 参数:
//   - text: 要输入的文本
//
// 返回: 错误信息，成功返回 nil
// 注意: 只支持英文字符，中文字符会调用 PasteText
func (bot *Bot) TypeText(text string) error {
	if hasChinese(text) {
		log.Println("检测到中文，使用 PasteText 方法")
		return bot.PasteText(text)
	}

	procKeybdEvent := user32.NewProc("keybd_event")
	for _, ch := range text {
		var vk uint32
		switch {
		case ch >= '0' && ch <= '9':
			vk = 0x30 + uint32(ch-'0')
		case ch >= 'a' && ch <= 'z':
			vk = 0x41 + uint32(ch-'a')
		case ch >= 'A' && ch <= 'Z':
			vk = 0x41 + uint32(ch-'A')
		case ch == ' ':
			vk = 0x20
		case ch == '\n':
			vk = 0x0D
		case ch == '-', ch == ':':
			vk = 0x0BD
		default:
			log.Printf("跳过不支持的字符: %c (Unicode: 0x%x)", ch, ch)
			continue
		}

		procKeybdEvent.Call(uintptr(vk), 0, 0, 0)
		time.Sleep(50 * time.Millisecond)
		procKeybdEvent.Call(uintptr(vk), 0, 0x0002, 0)
		time.Sleep(50 * time.Millisecond)
	}
	return nil
}

// PasteText 通过剪贴板粘贴文本
// 参数:
//   - text: 要粘贴的文本
//
// 返回: 错误信息，成功返回 nil
// 原理: 将文本放入剪贴板，然后模拟 Ctrl+V 粘贴
func (bot *Bot) PasteText(text string) error {
	procOpenClipboard := user32.NewProc("OpenClipboard")
	procEmptyClipboard := user32.NewProc("EmptyClipboard")
	procSetClipboardData := user32.NewProc("SetClipboardData")
	procCloseClipboard := user32.NewProc("CloseClipboard")
	procGlobalAlloc := kernel32.NewProc("GlobalAlloc")
	procGlobalLock := kernel32.NewProc("GlobalLock")
	procGlobalUnlock := kernel32.NewProc("GlobalUnlock")

	procOpenClipboard.Call(0)
	defer procCloseClipboard.Call()

	procEmptyClipboard.Call()

	size := (len(text) + 1) * 2
	hMem, _, _ := procGlobalAlloc.Call(0x0042, uintptr(size))
	if hMem == 0 {
		return fmt.Errorf("分配内存失败")
	}

	ptr, _, _ := procGlobalLock.Call(hMem)
	if ptr == 0 {
		return fmt.Errorf("锁定内存失败")
	}

	for i, ch := range []rune(text) {
		*(*uint16)(unsafe.Pointer(ptr + uintptr(i*2))) = uint16(ch)
	}
	*(*uint16)(unsafe.Pointer(ptr + uintptr(len(text)*2))) = 0

	procGlobalUnlock.Call(hMem)
	procSetClipboardData.Call(13, hMem)

	time.Sleep(500 * time.Millisecond)

	log.Println("按 Ctrl+V 粘贴文本...")
	procPostMessage := user32.NewProc("PostMessageW")
	procPostMessage.Call(bot.hwnd, WM_KEYDOWN, VK_CONTROL, 0)
	time.Sleep(50 * time.Millisecond)
	procPostMessage.Call(bot.hwnd, WM_KEYDOWN, VK_V, 0)
	time.Sleep(50 * time.Millisecond)
	procPostMessage.Call(bot.hwnd, WM_KEYUP, VK_V, 0)
	time.Sleep(50 * time.Millisecond)
	procPostMessage.Call(bot.hwnd, WM_KEYUP, VK_CONTROL, 0)
	time.Sleep(500 * time.Millisecond)

	log.Println("粘贴完成")
	return nil
}

// hasChinese 检测文本中是否包含中文字符
func hasChinese(text string) bool {
	for _, r := range text {
		if r >= 0x4e00 && r <= 0x9fff {
			return true
		}
	}
	return false
}

// ==================== 联系人操作方法 ====================

// SearchContact 搜索并选择联系人
// 参数:
//   - contactName: 联系人名称
//
// 返回: 错误信息，成功返回 nil
// 流程: 1.在搜索框输入联系人名称 2.等待搜索结果 3.匹配并点击正确的联系人
//
// 使用示例:
//
//	err := bot.SearchContact("文件传输助手")
//	if err != nil {
//	    log.Fatal(err)
//	}
func (bot *Bot) SearchContact(contactName string) error {
	log.Printf("搜索联系人: %s", contactName)

	searchEdit := bot.FindElementByName("搜索")
	if searchEdit == nil {
		return fmt.Errorf("未找到搜索框")
	}

	name, _ := searchEdit.Get_CurrentName()
	className, _ := searchEdit.Get_CurrentClassName()
	automationId, _ := searchEdit.Get_CurrentAutomationId()
	log.Printf("找到搜索框元素 - Name: %q, ClassName: %q, AutomationId: %q", name, className, automationId)

	// 先点击搜索框
	log.Println("点击搜索框...")
	bot.ClickElementWithInvoke(searchEdit)
	time.Sleep(2000 * time.Millisecond)

	// 设置搜索框文本值
	log.Println("设置搜索框文本值...")
	vp, _ := searchEdit.GetCurrentPattern(uia.UIA_ValuePatternId)
	if vp != nil {
		valuePattern := uia.NewIUIAutomationValuePattern(vp)
		defer valuePattern.Release()
		valuePattern.SetValue(contactName)
		log.Println("使用 ValuePattern 设置成功")
		currentValue, _ := valuePattern.Get_CurrentValue()
		log.Printf("搜索框当前值: %q", currentValue)
	} else {
		return fmt.Errorf("搜索框不支持 ValuePattern")
	}

	// 等待搜索结果
	time.Sleep(2000 * time.Millisecond)

	log.Println("检查搜索结果...")
	// 尝试多种匹配方式
	var matchedElem *uia.IUIAutomationElement

	// 方式1: 优先匹配搜索结果（search_item_function_）
	expectedSearchAutomationId := fmt.Sprintf("search_item_function_%s", contactName)
	expectedSessionAutomationId := fmt.Sprintf("session_item_%s", contactName)
	log.Printf("优先搜索 AutomationId: %q, 备选 AutomationId: %q", expectedSearchAutomationId, expectedSessionAutomationId)

	allElements, err := bot.root.FindAll(uia.TreeScope_Descendants, bot.ppv.CreateTrueCondition())
	if err != nil {
		return fmt.Errorf("查找元素失败: %v", err)
	}
	defer allElements.Release()

	length := allElements.Get_Length()
	log.Printf("窗口中有 %d 个元素", length)

	for i := int32(0); i < length; i++ {
		elem, _ := allElements.GetElement(i)
		if elem != nil {
			name, _ := elem.Get_CurrentName()
			className, _ := elem.Get_CurrentClassName()
			automationId, _ := elem.Get_CurrentAutomationId()

			// 记录所有可能的匹配项
			if strings.Contains(automationId, contactName) || strings.Contains(name, contactName) {
				firstLine := getFirstLine(name)
				log.Printf("[候选元素] AutomationId: %q, Name: %q (第一行: %q), ClassName: %q", automationId, name, firstLine, className)
			}

			// 方式1: 优先匹配搜索结果中的 mmui::XTableCell 类型
			if matchedElem == nil && automationId == expectedSearchAutomationId && className == "mmui::XTableCell" {
				matchedElem = elem
				firstLine := getFirstLine(name)
				log.Printf("[匹配成功] 通过搜索结果匹配 (XTableCell) - AutomationId: %q, Name: %q (第一行: %q), ClassName: %q", automationId, name, firstLine, className)
			}
			// 方式2: 备选匹配搜索结果中的 mmui::SearchContentCellView 类型
			if matchedElem == nil && automationId == expectedSearchAutomationId && className == "mmui::SearchContentCellView" {
				matchedElem = elem
				firstLine := getFirstLine(name)
				log.Printf("[匹配成功] 通过搜索结果匹配 (SearchContentCellView) - AutomationId: %q, Name: %q (第一行: %q), ClassName: %q", automationId, name, firstLine, className)
			}
			elem.Release()
		}
	}

	// 如果没找到，尝试通过 Name 匹配
	if matchedElem == nil {
		log.Println("通过 AutomationId 未找到，尝试通过 Name 匹配...")
		for i := int32(0); i < length; i++ {
			elem, _ := allElements.GetElement(i)
			if elem != nil {
				name, _ := elem.Get_CurrentName()
				className, _ := elem.Get_CurrentClassName()
				automationId, _ := elem.Get_CurrentAutomationId()

				// 提取 Name 的第一行（联系人名）
				firstLine := getFirstLine(name)

				// 方式2: 通过 Name 匹配（排除包含其他文本的）
				if matchedElem == nil && firstLine == contactName &&
					(className == "mmui::ChatSessionCell" || className == "mmui::SearchContentCellView") {
					matchedElem = elem
					log.Printf("通过 Name 匹配 - AutomationId: %q, Name: %q, ClassName: %q", automationId, name, className)
				}
				elem.Release()
			}
		}
	}

	if matchedElem != nil {
		log.Printf("找到联系人: %s", contactName)
		err := bot.ClickElementWithInvoke(matchedElem)
		matchedElem.Release()
		if err != nil {
			return err
		}
		log.Println("点击联系人成功，等待3秒进入聊天窗口...")
		time.Sleep(3000 * time.Millisecond)
		return nil
	}

	log.Printf("未找到联系人: %s", contactName)
	return fmt.Errorf("未找到联系人: %s", contactName)
}

// ReadChatMessages 读取当前聊天窗口的消息列表
// 返回: 消息列表，每个消息包含内容和是否是自己发送的
func (bot *Bot) ReadChatMessages() ([]*ChatMessage, error) {
	log.Println("读取聊天消息列表...")

	// 查找消息列表元素
	msgList := bot.FindElementByAutomationId("chat_message_list")
	if msgList == nil {
		log.Println("未通过 AutomationId 找到消息列表，尝试其他方式...")
		// 尝试通过 ClassName 查找
		msgList = bot.FindElementByClassName("mmui::RecyclerListView")
		if msgList == nil {
			return nil, fmt.Errorf("未找到消息列表 (AutomationId: chat_message_list, ClassName: mmui::RecyclerListView)")
		}
	}
	defer msgList.Release()

	name, _ := msgList.Get_CurrentName()
	className, _ := msgList.Get_CurrentClassName()
	log.Printf("找到消息列表 - Name: %q, ClassName: %q", name, className)

	// 获取所有子元素
	children, err := msgList.FindAll(uia.TreeScope_Children, bot.ppv.CreateTrueCondition())
	if err != nil {
		return nil, fmt.Errorf("查找消息子元素失败: %v", err)
	}
	defer children.Release()

	var messages []*ChatMessage
	count := children.Get_Length()
	log.Printf("找到 %d 个消息元素", count)

	for i := int32(0); i < count; i++ {
		elem, _ := children.GetElement(i)
		if elem == nil {
			continue
		}

		elemName, _ := elem.Get_CurrentName()
		elemName = strings.TrimSpace(elemName)

		log.Printf("元素 %d - Name: %q", i+1, elemName)

		// 跳过空消息
		if elemName == "" {
			elem.Release()
			continue
		}

		// 移除"列表项目"后缀
		elemName = strings.TrimSuffix(elemName, " 列表项目")
		elemName = strings.TrimSpace(elemName)

		// 判断是否是时间
		msgTime := ""
		if isTimeString(elemName) {
			msgTime = elemName
			log.Printf("识别为时间: %q", msgTime)
			elem.Release()
			continue
		}

		// 创建消息对象，默认为对方发送的消息
		msg := &ChatMessage{
			Content: elemName,
			IsSelf:  false,
			Time:    msgTime,
			MsgType: ParseMessageType(elemName),
		}

		messages = append(messages, msg)

		// 记录为 JSON
		jsonData, _ := json.Marshal(msg)
		log.Printf("消息 %d: %s", len(messages), string(jsonData))

		elem.Release()
	}

	log.Printf("成功读取 %d 条消息", len(messages))
	return messages, nil
}

// SaveChatMessages 保存聊天记录到文件
// 参数:
//   - wechatId: 微信号
//   - contact: 联系人名称
//   - messages: 消息列表
//
// 返回: 错误信息，成功返回 nil
func (bot *Bot) SaveChatMessages(wechatId, contact string, messages []*ChatMessage) error {
	if len(messages) == 0 {
		return nil
	}

	filename := fmt.Sprintf("%s_%s.json", wechatId, contact)
	jsonData, err := json.MarshalIndent(messages, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化消息失败: %v", err)
	}

	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("保存文件失败: %v", err)
	}

	log.Printf("保存聊天记录到文件: %s (%d 条消息)", filename, len(messages))
	return nil
}

// LoadChatMessages 从文件加载聊天记录
// 参数:
//   - wechatId: 微信号
//   - contact: 联系人名称
//
// 返回: 消息列表和错误信息
func (bot *Bot) LoadChatMessages(wechatId, contact string) ([]*ChatMessage, error) {
	filename := fmt.Sprintf("%s_%s.json", wechatId, contact)

	jsonData, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("读取文件失败: %v", err)
	}

	var messages []*ChatMessage
	err = json.Unmarshal(jsonData, &messages)
	if err != nil {
		return nil, fmt.Errorf("解析文件失败: %v", err)
	}

	log.Printf("加载聊天记录从文件: %s (%d 条消息)", filename, len(messages))
	return messages, nil
}

// ==================== 消息发送方法 ====================

// SendMessage 发送消息到当前聊天窗口
// 参数:
//   - message: 要发送的消息内容
//
// 返回: 错误信息，成功返回 nil
// 流程: 1.查找并点击消息输入框 2.输入消息文本 3.查找并点击发送按钮
// 注意: 发送前需要先调用 SearchContact 打开目标聊天窗口
//
// 使用示例:
//
//	bot.SearchContact("文件传输助手")
//	err := bot.SendMessage("测试消息")
//	if err != nil {
//	    log.Fatal(err)
//	}
func (bot *Bot) SendMessage(message string) error {
	log.Printf("发送消息: %s", message)

	// 先激活微信窗口
	log.Println("激活微信窗口...")
	user32.NewProc("SetForegroundWindow").Call(bot.hwnd)
	user32.NewProc("BringWindowToTop").Call(bot.hwnd)
	time.Sleep(300 * time.Millisecond)

	log.Println("查找消息输入框...")
	msgEdit := bot.FindElementByAutomationId("chat_input_field")
	if msgEdit == nil {
		log.Println("未通过 AutomationId 找到，尝试通过 ClassName 查找")
		msgEdit = bot.FindElementByClassName("mmui::ChatInputField")
	}
	if msgEdit == nil {
		log.Println("未通过 ClassName 找到，尝试通过 Name 查找")
		msgEdit = bot.FindElementByName("文件传输助手")
	}
	if msgEdit == nil {
		log.Println("未找到消息输入框，尝试查找所有 XView 元素")
		msgEdit = bot.FindElementByClassName("mmui::XView")
	}
	if msgEdit == nil {
		return fmt.Errorf("未找到消息输入框")
	}
	defer msgEdit.Release()

	name, _ := msgEdit.Get_CurrentName()
	className, _ := msgEdit.Get_CurrentClassName()
	automationId, _ := msgEdit.Get_CurrentAutomationId()
	log.Printf("找到消息输入框 - Name: %q, ClassName: %q, AutomationId: %q", name, className, automationId)

	err := bot.ClickElementWithInvoke(msgEdit)
	if err != nil {
		return err
	}
	time.Sleep(1000 * time.Millisecond)

	log.Println("尝试使用 ValuePattern 设置文本...")
	err = bot.SetTextValue(msgEdit, message)
	if err != nil {
		log.Printf("ValuePattern 设置失败: %v，尝试使用 TypeText", err)
		err = bot.TypeText(message)
		if err != nil {
			log.Printf("TypeText 失败: %v", err)
			return err
		}
		log.Println("使用 TypeText 设置成功")
	} else {
		log.Println("使用 ValuePattern 设置成功")
	}
	time.Sleep(2000 * time.Millisecond)

	log.Println("查找发送按钮...")
	allElements, err := bot.root.FindAll(uia.TreeScope_Descendants, bot.ppv.CreateTrueCondition())
	if err != nil {
		return fmt.Errorf("查找元素失败: %v", err)
	}
	defer allElements.Release()

	length := allElements.Get_Length()
	log.Printf("窗口中有 %d 个元素", length)

	var sendButton *uia.IUIAutomationElement
	for i := int32(0); i < length; i++ {
		elem, _ := allElements.GetElement(i)
		if elem != nil {
			name, _ := elem.Get_CurrentName()
			className, _ := elem.Get_CurrentClassName()

			if strings.Contains(name, "发送") && strings.Contains(name, "S") {
				log.Printf("找到发送按钮 - Name: %q, ClassName: %q", name, className)
				sendButton = elem
				break
			}
			elem.Release()
		}
	}

	if sendButton == nil {
		return fmt.Errorf("未找到发送按钮")
	}
	defer sendButton.Release()

	name, _ = sendButton.Get_CurrentName()
	className, _ = sendButton.Get_CurrentClassName()
	log.Printf("发送按钮信息 - Name: %q, ClassName: %q", name, className)

	log.Println("使用 Tab + Enter 点击发送按钮")

	// 激活窗口
	user32.NewProc("SetForegroundWindow").Call(bot.hwnd)
	user32.NewProc("BringWindowToTop").Call(bot.hwnd)
	time.Sleep(200 * time.Millisecond)

	// 按 Tab 键切换焦点到发送按钮
	log.Println("按 Tab 键...")
	procKeybdEvent := user32.NewProc("keybd_event")
	procKeybdEvent.Call(0x09, 0, 0, 0) // Tab 按下
	time.Sleep(100 * time.Millisecond)
	procKeybdEvent.Call(0x09, 0, 0x0002, 0) // Tab 抬起
	time.Sleep(300 * time.Millisecond)

	// 按 Enter 键触发发送
	log.Println("按 Enter 键...")
	procKeybdEvent.Call(0x0D, 0, 0, 0) // Enter 按下
	time.Sleep(100 * time.Millisecond)
	procKeybdEvent.Call(0x0D, 0, 0x0002, 0) // Enter 抬起
	time.Sleep(1000 * time.Millisecond)

	log.Println("消息发送完成！")
	return nil
}

// ==================== 会话列表读取方法 ====================

// ReadSessionList 读取会话列表（左侧聊天列表）
// 返回: 会话列表和可能的错误
// 功能: 从会话列表中提取所有会话信息
//
// 使用示例:
//
//	sessions, err := bot.ReadSessionList()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, s := range sessions {
//	    fmt.Printf("发送者: %s, 内容: %s, 类型: %s\n", s.Sender, s.Content, s.MsgType)
//	}
func (bot *Bot) ReadSessionList() ([]*Session, error) {
	log.Println("读取会话列表...")

	var sessions []*Session

	list := bot.FindElementByAutomationId("session_list")
	if list == nil {
		list = bot.FindElementByClassName("mmui::XTableView")
	}
	if list == nil {
		list = bot.FindElementByName("会话")
	}
	if list == nil {
		return nil, fmt.Errorf("未找到会话列表")
	}
	defer list.Release()

	name, _ := list.Get_CurrentName()
	className, _ := list.Get_CurrentClassName()
	automationId, _ := list.Get_CurrentAutomationId()
	log.Printf("找到会话列表 - Name: %q, ClassName: %q, AutomationId: %q", name, className, automationId)

	variant, _ := uia.VariantFromString("mmui::ChatSessionCell")
	cond, _ := bot.ppv.CreatePropertyCondition(uia.UIA_ClassNamePropertyId, variant)

	children, err := list.FindAll(uia.TreeScope_Descendants, cond)
	if err != nil {
		log.Printf("通过ClassName查找失败，尝试遍历所有子元素: %v", err)
		children, err = list.FindAll(uia.TreeScope_Descendants, bot.ppv.CreateTrueCondition())
		if err != nil {
			return nil, fmt.Errorf("查找子元素失败: %v", err)
		}
	}
	if children != nil {
		defer children.Release()
	}

	arrLen := children.Get_Length()
	log.Printf("找到 %d 个会话元素", arrLen)

	for i := int32(0); i < arrLen; i++ {
		elem, _ := children.GetElement(i)
		if elem == nil {
			continue
		}

		elemClassName, _ := elem.Get_CurrentClassName()
		if elemClassName != "mmui::ChatSessionCell" {
			elem.Release()
			continue
		}

		elemName, _ := elem.Get_CurrentName()
		elemAutomationId, _ := elem.Get_CurrentAutomationId()

		session := bot.parseSessionElement(elemName, elemAutomationId)
		if session != nil {
			sessions = append(sessions, session)
			bot.sessionManager.UpdateSession(session)
		}

		elem.Release()
	}

	log.Printf("成功读取 %d 个会话", len(sessions))
	return sessions, nil
}

// parseSessionElement 解析会话元素信息
// 参数:
//   - name: 元素Name属性（包含发送者、内容、时间）
//   - automationId: 元素AutomationId
//
// 返回: 解析后的Session对象
func (bot *Bot) parseSessionElement(name, automationId string) *Session {
	if name == "" {
		return nil
	}

	lines := strings.Split(name, "\n")
	if len(lines) < 2 {
		return nil
	}

	sender := strings.TrimSpace(lines[0])
	content := ""
	timeStr := ""

	// 第二行是消息内容
	if len(lines) >= 2 {
		content = strings.TrimSpace(lines[1])
	}

	// 第三行及以后，查找时间
	for i := 2; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}
		if isTimeString(line) {
			timeStr = line
			break
		}
	}

	if sender == "" {
		return nil
	}

	msgType := ParseMessageType(content)

	return &Session{
		Sender:       sender,
		Content:      content,
		Time:         timeStr,
		AutomationId: automationId,
		IsSelf:       false,
		MsgType:      msgType,
	}
}

// isTimeString 判断字符串是否为时间格式
func isTimeString(s string) bool {
	timePatterns := []string{"星期一", "星期二", "星期三", "星期四", "星期五", "星期六", "星期日",
		"昨天", "前天", "刚刚", "分钟前", "小时前"}
	for _, pattern := range timePatterns {
		if strings.Contains(s, pattern) {
			return true
		}
	}
	if len(s) <= 5 && strings.Contains(s, ":") {
		return true
	}
	return false
}

// ==================== 新消息检测方法 ====================

// CheckNewMessages 检测新消息
// 参数:
//   - replyRecord: 回复记录器，用于检查是否是自己发送的回复
//
// 返回: 新消息列表和错误信息
// 功能: 比较当前会话列表与上次保存的状态，检测是否有新消息
// 注意: 首次调用会初始化状态，不会返回新消息
//
// 使用示例:
//
//	for {
//	    newMsgs, err := bot.CheckNewMessages(replyRecord)
//	    if err != nil {
//	        log.Println(err)
//	        continue
//	    }
//	    for _, msg := range newMsgs {
//	        fmt.Printf("新消息: %s - %s\n", msg.Sender, msg.Content)
//	    }
//	    time.Sleep(3 * time.Second)
//	}
func (bot *Bot) CheckNewMessages(replyRecord *ReplyRecord) ([]*NewMessage, error) {
	currentSessions, err := bot.ReadSessionList()
	if err != nil {
		return nil, err
	}

	var newMessages []*NewMessage

	for _, session := range currentSessions {
		if !bot.sessionManager.HasSender(session.Sender) {
			bot.sessionManager.SetLastContent(session.Sender, session.Content)
			log.Printf("[初始化] %s - 内容: %q", session.Sender, session.Content)
			continue
		}

		lastContent := bot.sessionManager.GetLastContent(session.Sender)
		if session.Content != lastContent && session.Content != "" {
			// 去掉截断标记...，得到完整内容
			currentContent := strings.TrimSuffix(session.Content, "…")

			// 检查是否是自己发送的回复
			isSelfReply := false
			if replyRecord != nil {
				isSelfReply = replyRecord.IsSelfReply(session.Sender, currentContent)

				// 如果通过内容匹配判断为自己回复，再检查聊天记录确认
				if isSelfReply {
					// 加载聊天记录，检查最新一条消息是否是自己发送的
					historyMessages, err := bot.LoadChatMessages(bot.wechatId, session.Sender)
					if err == nil && len(historyMessages) > 0 {
						lastMsg := historyMessages[len(historyMessages)-1]
						// 如果最新一条消息是自己发送的，且内容匹配，则确认是自己回复
						if lastMsg.IsSelf && lastMsg.Content == currentContent {
							log.Printf("[确认自己回复] 发送者: %s, 内容: %q, 时间: %q", session.Sender, currentContent, lastMsg.Time)
						} else {
							// 不是自己回复，继续处理
							isSelfReply = false
							log.Printf("[不是自己回复] 发送者: %s, 内容: %q, 最新消息IsSelf: %v", session.Sender, currentContent, lastMsg.IsSelf)
						}
					} else {
						log.Printf("[无法加载聊天记录] 发送者: %s, 错误: %v", session.Sender, err)
					}
				}

				log.Printf("[检查自己回复] 发送者: %s, 当前内容: %q, 是否自己回复: %v", session.Sender, currentContent, isSelfReply)
			}

			if isSelfReply {
				log.Printf("[跳过自己回复] 发送者: %s, 内容: %q", session.Sender, currentContent)
			} else {
				newMsg := &NewMessage{
					Sender:  session.Sender,
					Content: session.Content,
					Time:    session.Time,
					IsSelf:  session.IsSelf,
					MsgType: session.MsgType,
				}
				newMessages = append(newMessages, newMsg)
				log.Printf("[检测到新消息] 发送者: %s, 上次内容: %q, 当前内容: %q", session.Sender, lastContent, session.Content)
			}
		}

		bot.sessionManager.SetLastContent(session.Sender, session.Content)
	}

	return newMessages, nil
}

// UpdateSessionState 手动更新指定联系人的会话状态
// 发送消息后调用此方法，避免自己的消息被误认为新消息
// 参数:
//   - sender: 联系人名称
func (bot *Bot) UpdateSessionState(sender string) {
	// 重试3次，等待会话列表更新
	for i := 0; i < 3; i++ {
		if i > 0 {
			time.Sleep(500 * time.Millisecond)
		}

		currentSessions, err := bot.ReadSessionList()
		if err != nil {
			log.Printf("更新会话状态失败（第%d次重试）: %v", i+1, err)
			continue
		}

		for _, session := range currentSessions {
			if session.Sender == sender {
				bot.sessionManager.SetLastContent(sender, session.Content)
				log.Printf("已更新会话状态（第%d次重试） - 发送者: %s, 内容: %s", i+1, sender, session.Content)
				return
			}
		}
		log.Printf("第%d次重试：未找到联系人 %s 的会话", i+1, sender)
	}

	log.Printf("更新会话状态失败：3次重试后仍未找到联系人 %s", sender)
}

// ShouldFilterMessage 判断是否应该过滤该消息
// 参数:
//   - sender: 发送者名称
//   - content: 消息内容
//   - config: 过滤配置，如果为 nil 则使用默认配置
//
// 返回: true 表示应该过滤（不回复），false 表示不过滤
func ShouldFilterMessage(sender, content string, config *FilterConfig) bool {
	if config == nil {
		config = DefaultFilterConfig
	}

	// 检查前缀
	for _, prefix := range config.PublicAccountPrefixes {
		if strings.HasPrefix(sender, prefix) {
			return true
		}
	}

	// 检查关键词
	for _, keyword := range config.PublicAccountKeywords {
		if strings.Contains(sender, keyword) {
			return true
		}
	}

	// 检查折叠的聊天
	if config.FilterCollapsedChats {
		// 折叠的聊天通常包含括号和数字，如 "群聊 (3)"、"联系人 (5)"
		if strings.Contains(sender, "(") && strings.Contains(sender, ")") {
			// 检查括号内是否是纯数字
			start := strings.Index(sender, "(")
			end := strings.Index(sender, ")")
			if start != -1 && end != -1 && end > start+1 {
				inside := sender[start+1 : end]
				if isAllDigits(inside) {
					return true
				}
			}
		}
	}

	// 检查自定义过滤函数
	for _, filter := range config.CustomFilters {
		if filter != nil && filter(sender, content) {
			return true
		}
	}

	return false
}

// isAllDigits 检查字符串是否全是数字
func isAllDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// getFirstLine 提取字符串的第一行
func getFirstLine(s string) string {
	if s == "" {
		return ""
	}

	// 查找第一个换行符
	index := strings.Index(s, "\n")
	if index == -1 {
		return s
	}

	return s[:index]
}

// StartMessageMonitor 启动消息监控
// 参数:
//   - interval: 检查间隔时间
//   - callback: 新消息回调函数
//
// 功能: 定期检查新消息并调用回调函数
// 注意: 此方法会阻塞当前 goroutine，建议在单独的 goroutine 中调用
//
// 使用示例:
//
//	go bot.StartMessageMonitor(3*time.Second, func(msg *wechat.NewMessage) {
//	    fmt.Printf("收到新消息: %s - %s\n", msg.Sender, msg.Content)
//	})
func (bot *Bot) StartMessageMonitor(interval time.Duration, callback func(*NewMessage)) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	log.Printf("启动消息监控，检查间隔: %v", interval)

	for range ticker.C {
		newMessages, err := bot.CheckNewMessages(nil)
		if err != nil {
			log.Printf("检查新消息失败: %v", err)
			continue
		}

		for _, msg := range newMessages {
			if callback != nil {
				callback(msg)
			}
		}
	}
}

// StartAutoReply 启动自动回复功能
// 参数:
//   - config: 自动回复配置
//
// 功能: 监控指定联系人的聊天消息并自动回复
// 注意: 此方法会阻塞当前 goroutine，建议在单独的 goroutine 中调用
//
// 使用示例:
//
//	config := &wechat.AutoReplyConfig{
//	    Contacts: []string{"文件传输助手", "张三"},
//	    ReplyGenerator: func(content string) string {
//	        return "自动回复: " + content
//	    },
//	}
//	go bot.StartAutoReply(config)
func (bot *Bot) StartAutoReply(config *AutoReplyConfig) {
	if config == nil {
		log.Println("自动回复配置为空，无法启动")
		return
	}

	if len(config.Contacts) == 0 {
		log.Println("联系人列表为空，无法启动自动回复")
		return
	}

	log.Printf("启动自动回复，监控联系人: %v", config.Contacts)

	bot.wechatId = config.WechatId

	replyRecord := NewReplyRecord()

	// 等待3秒，给用户准备时间
	time.Sleep(3 * time.Second)

	// 初始化：记录所有会话的当前消息状态
	log.Println("初始化：记录当前所有会话的消息状态...")
	currentSessions, err := bot.ReadSessionList()
	if err != nil {
		log.Printf("读取会话列表失败: %v", err)
		return
	}

	for _, session := range currentSessions {
		bot.sessionManager.SetLastContent(session.Sender, session.Content)
		log.Printf("初始化记录 - 发送者: %s, 内容: %s", session.Sender, session.Content)
	}

	log.Println("初始化完成，开始监控新消息...")

	for {
		// 监控会话列表，检测新消息
		newMessages, err := bot.CheckNewMessages(replyRecord)
		if err != nil {
			log.Printf("检查新消息失败: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		// 处理新消息
		for _, msg := range newMessages {
			// 检查是否在监控的联系人列表中
			if !contains(config.Contacts, msg.Sender) {
				continue
			}

			// 检查是否应该过滤
			if config.FilterConfig != nil && ShouldFilterMessage(msg.Sender, msg.Content, config.FilterConfig) {
				continue
			}

			// 打开联系人对话框
			err := bot.SearchContact(msg.Sender)
			if err != nil {
				if config.OnError != nil {
					config.OnError(fmt.Errorf("打开联系人 %s 失败: %v", msg.Sender, err))
				}
				continue
			}

			// 等待对话框加载
			time.Sleep(1 * time.Second)

			// 读取聊天消息
			messages, err := bot.ReadChatMessages()
			if err != nil {
				if config.OnError != nil {
					config.OnError(fmt.Errorf("读取 %s 的消息失败: %v", msg.Sender, err))
				}
				continue
			}

			if len(messages) == 0 {
				continue
			}

			// 加载历史聊天记录
			history, _ := bot.LoadChatMessages(config.WechatId, msg.Sender)

			// 获取最后一条消息
			lastMessage := messages[len(messages)-1]

			// 检查是否是新消息（与历史记录最后一条不同）
			if len(history) > 0 {
				lastHistory := history[len(history)-1]
				if lastMessage.Content == lastHistory.Content {
					log.Printf("[跳过] %s - 消息与历史记录相同: %q", msg.Sender, lastMessage.Content)
					continue
				}
			}

			// 检查是否是自己发送的
			if replyRecord.IsSelfReply(msg.Sender, lastMessage.Content) {
				log.Printf("[跳过] %s - 这是自己刚才的回复: %q", msg.Sender, lastMessage.Content)
				continue
			}

			// 触发收到消息回调
			if config.OnMessage != nil {
				config.OnMessage(msg.Sender, lastMessage.Content)
			}

			// 生成回复
			reply := ""
			if config.ReplyGenerator != nil {
				reply = config.ReplyGenerator(lastMessage.Content)
			}

			if reply == "" {
				continue
			}

			// 触发发送回复回调
			if config.OnReply != nil {
				config.OnReply(msg.Sender, reply)
			}

			// 记录回复内容
			log.Printf("[记录回复] %s - 回复内容: %q", msg.Sender, reply)
			replyRecord.RecordReplyWithInfo(msg.Sender, reply, time.Now().Format("15:04"), ParseMessageType(reply))

			// 发送回复
			err = bot.SendMessage(reply)
			if err != nil {
				if config.OnError != nil {
					config.OnError(fmt.Errorf("回复失败: %v", err))
				}
				continue
			}

			// 将回复消息添加到聊天记录中
			replyMsg := &ChatMessage{
				Content: reply,
				IsSelf:  true,
				Time:    time.Now().Format("15:04"),
				MsgType: ParseMessageType(reply),
			}
			messages = append(messages, replyMsg)

			// 保存聊天记录
			err = bot.SaveChatMessages(config.WechatId, msg.Sender, messages)
			if err != nil {
				log.Printf("保存聊天记录失败: %v", err)
			}

			// 更新会话状态，避免回复被误判为新消息
			bot.sessionManager.SetLastContent(msg.Sender, reply)
		}

		time.Sleep(2 * time.Second)
	}
}

// contains 检查字符串是否在切片中
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// processAutoReply 处理单个联系人的自动回复
func (bot *Bot) processAutoReply(contact string, config *AutoReplyConfig, replyRecord *ReplyRecord) {
	// 检查是否应该过滤
	if config.FilterConfig != nil && ShouldFilterMessage(contact, "", config.FilterConfig) {
		return
	}

	// 打开联系人对话框
	err := bot.SearchContact(contact)
	if err != nil {
		if config.OnError != nil {
			config.OnError(fmt.Errorf("打开联系人 %s 失败: %v", contact, err))
		}
		return
	}

	// 等待对话框加载
	time.Sleep(1 * time.Second)

	// 读取聊天消息
	messages, err := bot.ReadChatMessages()
	if err != nil {
		if config.OnError != nil {
			config.OnError(fmt.Errorf("读取 %s 的消息失败: %v", contact, err))
		}
		return
	}

	if len(messages) == 0 {
		return
	}

	// 获取最后一条消息
	lastMessage := messages[len(messages)-1]

	// 检查是否是自己发送的
	if replyRecord.IsSelfReply(contact, lastMessage.Content) {
		return
	}

	// 触发收到消息回调
	if config.OnMessage != nil {
		config.OnMessage(contact, lastMessage.Content)
	}

	// 生成回复
	reply := ""
	if config.ReplyGenerator != nil {
		reply = config.ReplyGenerator(lastMessage.Content)
	}

	if reply == "" {
		return
	}

	// 触发发送回复回调
	if config.OnReply != nil {
		config.OnReply(contact, reply)
	}

	// 记录回复内容
	replyRecord.RecordReply(contact, reply)

	// 发送回复
	err = bot.SendMessage(reply)
	if err != nil {
		if config.OnError != nil {
			config.OnError(fmt.Errorf("回复失败: %v", err))
		}
	}
}

// GetSessionManager 获取会话管理器
// 返回: 当前使用的会话管理器实例
func (bot *Bot) GetSessionManager() *SessionManager {
	return bot.sessionManager
}
