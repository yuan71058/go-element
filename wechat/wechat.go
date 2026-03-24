// Package wechat 提供微信自动化操作功能
package wechat

import (
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

	// 激活微信窗口
	log.Println("激活微信窗口...")
	user32.NewProc("SetForegroundWindow").Call(hwnd)
	time.Sleep(500 * time.Millisecond)

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

// ClickElement 点击指定的UI元素
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
	time.Sleep(2000 * time.Millisecond)

	log.Println("检查搜索结果...")
	expectedAutomationId := fmt.Sprintf("search_item_function_%s", contactName)
	log.Printf("期望的 AutomationId: %q", expectedAutomationId)

	allElements, err := bot.root.FindAll(uia.TreeScope_Descendants, bot.ppv.CreateTrueCondition())
	if err != nil {
		return fmt.Errorf("查找元素失败: %v", err)
	}
	defer allElements.Release()

	length := allElements.Get_Length()
	log.Printf("窗口中有 %d 个元素", length)

	matchedCount := 0
	for i := int32(0); i < length; i++ {
		elem, _ := allElements.GetElement(i)
		if elem != nil {
			name, _ := elem.Get_CurrentName()
			className, _ := elem.Get_CurrentClassName()
			automationId, _ := elem.Get_CurrentAutomationId()

			if name == contactName {
				matchedCount++
				log.Printf("匹配到第%d个Name为'%s'的元素 - ClassName: %q, AutomationId: %q", matchedCount, contactName, className, automationId)

				if className == "mmui::XTableCell" && automationId == expectedAutomationId {
					log.Printf("找到联系人: %s (AutomationId: %s)", contactName, automationId)
					err := bot.ClickElement(elem)
					elem.Release()
					if err != nil {
						return err
					}
					log.Println("点击联系人成功，等待2秒进入聊天窗口...")
					time.Sleep(2000 * time.Millisecond)
					return nil
				}
			}
			elem.Release()
		}
	}

	if matchedCount > 0 {
		log.Printf("共找到 %d 个Name为'%s'的元素，但都不符合条件", matchedCount, contactName)
	} else {
		log.Printf("未找到任何Name为'%s'的元素", contactName)
	}

	return fmt.Errorf("未找到联系人: %s", contactName)
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

	err := bot.ClickElement(msgEdit)
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
	time.Sleep(500 * time.Millisecond)

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

			if strings.Contains(name, "发送") {
				log.Printf("找到 Name 包含'发送'的元素 - Name: %q, ClassName: %q", name, className)
				if className == "mmui::XTextView" {
					sendButton = elem
					break
				}
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
	log.Printf("找到发送按钮 - Name: %q, ClassName: %q", name, className)

	log.Println("点击发送按钮")
	clickErr := bot.ClickElement(sendButton)
	if clickErr != nil {
		log.Printf("点击发送按钮失败: %v", clickErr)
		return clickErr
	}
	time.Sleep(500 * time.Millisecond)

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

	for i := 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}
		if isTimeString(line) {
			timeStr = line
		} else if content == "" {
			content = line
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
// 返回: 新消息列表和错误信息
// 功能: 比较当前会话列表与上次保存的状态，检测是否有新消息
// 注意: 首次调用会初始化状态，不会返回新消息
//
// 使用示例:
//
//	for {
//	    newMsgs, err := bot.CheckNewMessages()
//	    if err != nil {
//	        log.Println(err)
//	        continue
//	    }
//	    for _, msg := range newMsgs {
//	        fmt.Printf("新消息: %s - %s\n", msg.Sender, msg.Content)
//	    }
//	    time.Sleep(3 * time.Second)
//	}
func (bot *Bot) CheckNewMessages() ([]*NewMessage, error) {
	currentSessions, err := bot.ReadSessionList()
	if err != nil {
		return nil, err
	}

	var newMessages []*NewMessage

	for _, session := range currentSessions {
		if !bot.sessionManager.HasSender(session.Sender) {
			bot.sessionManager.SetLastContent(session.Sender, session.Content)
			continue
		}

		lastContent := bot.sessionManager.GetLastContent(session.Sender)
		if session.Content != lastContent && session.Content != "" {
			newMsg := &NewMessage{
				Sender:  session.Sender,
				Content: session.Content,
				Time:    session.Time,
				IsSelf:  session.IsSelf,
				MsgType: session.MsgType,
			}
			newMessages = append(newMessages, newMsg)
			log.Printf("检测到新消息 - 发送者: %s, 内容: %s", session.Sender, session.Content)
		}

		bot.sessionManager.SetLastContent(session.Sender, session.Content)
	}

	return newMessages, nil
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
		newMessages, err := bot.CheckNewMessages()
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

// GetSessionManager 获取会话管理器
// 返回: 当前使用的会话管理器实例
func (bot *Bot) GetSessionManager() *SessionManager {
	return bot.sessionManager
}
