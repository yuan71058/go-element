// 微信自动化机器人主程序
package main

import (
	"fmt"     // 格式化输出
	"log"     // 日志记录
	"os"      // 操作系统接口
	"strings" // 字符串处理
	"syscall" // 系统调用
	"time"    // 时间处理
	"unsafe"  // 不安全操作（用于指针操作）

	uia "github.com/auuunya/go-element" // UI自动化库
)

// 常量定义
const (
	WeChatWindowTitle = "微信"                 // 微信窗口标题
	WeChatClassName   = "Qt51514QWindowIcon" // 微信窗口类名
)

// WeChatBot 微信机器人结构体
type WeChatBot struct {
	ppv     *uia.IUIAutomation        // UI自动化接口
	hwnd    uintptr                   // 微信窗口句柄
	root    *uia.IUIAutomationElement // 微信窗口根元素
	logFile *os.File                  // 日志文件句柄
}

// main 主函数 - 程序入口
func main() {
	// 创建微信机器人实例
	bot, err := NewWeChatBot()
	if err != nil {
		log.Fatalf("初始化失败: %v", err)
	}
	defer bot.Close() // 程序退出时清理资源

	log.Println("=== 测试发送消息功能 ===")

	// 设置测试联系人
	contactName := "文件传输助手"
	// 生成带时间戳的测试消息
	message := "测试消息 - " + time.Now().Format("2006-01-02 15:04:05")

	log.Printf("联系人: %s", contactName)
	log.Printf("消息: %s", message)

	// 搜索联系人
	err = bot.SearchContact(contactName)
	if err != nil {
		log.Fatalf("搜索联系人失败: %v", err)
	}

	log.Println("搜索成功，等待2秒...")
	time.Sleep(2 * time.Second)

	// 发送消息
	err = bot.SendMessage(message)
	if err != nil {
		log.Fatalf("发送消息失败: %v", err)
	}

	log.Println("消息发送成功！")
}

// NewWeChatBot 创建微信机器人实例
// 返回初始化好的 WeChatBot 对象和可能的错误
func NewWeChatBot() (*WeChatBot, error) {
	// 创建日志文件
	f, err := os.OpenFile("wechat_demo.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("无法创建日志文件: %v", err)
	}
	// 设置日志输出到文件
	log.SetOutput(f)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("--- 微信自动化启动 ---")

	// 初始化 COM 组件
	err = uia.CoInitialize()
	if err != nil {
		return nil, fmt.Errorf("COM 初始化失败: %v", err)
	}

	// 查找微信窗口
	log.Println("正在查找微信窗口...")
	hwnd, err := uia.GetWindowForString("", WeChatWindowTitle)
	// 如果通过窗口标题找不到，尝试通过类名查找
	if err != nil || hwnd == 0 {
		hwnd, _ = uia.GetWindowForString(WeChatClassName, "")
	}
	// 如果还是找不到，返回错误
	if hwnd == 0 {
		uia.CoUninitialize()
		return nil, fmt.Errorf("未找到微信窗口，请确保微信已打开")
	}
	log.Printf("找到微信窗口，句柄: 0x%X", hwnd)

	// 激活微信窗口，将其置于前台
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

	// 返回初始化好的机器人实例
	return &WeChatBot{
		ppv:     ppv,
		hwnd:    hwnd,
		root:    root,
		logFile: f,
	}, nil
}

// Close 关闭机器人，释放所有资源
func (bot *WeChatBot) Close() {
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

// FindElementByName 通过名称查找UI元素
// 参数: name - 元素名称
// 返回: 找到的元素，未找到返回 nil
func (bot *WeChatBot) FindElementByName(name string) *uia.IUIAutomationElement {
	variant, _ := uia.VariantFromString(name)
	cond, _ := bot.ppv.CreatePropertyCondition(uia.UIA_NamePropertyId, variant)
	elem, _ := bot.root.FindFirst(uia.TreeScope_Descendants, cond)
	return elem
}

// FindElementByAutomationId 通过 AutomationId 查找UI元素
// 参数: id - 元素的 AutomationId
// 返回: 找到的元素，未找到返回 nil
func (bot *WeChatBot) FindElementByAutomationId(id string) *uia.IUIAutomationElement {
	variant, _ := uia.VariantFromString(id)
	cond, _ := bot.ppv.CreatePropertyCondition(uia.UIA_AutomationIdPropertyId, variant)
	elem, _ := bot.root.FindFirst(uia.TreeScope_Descendants, cond)
	return elem
}

// FindElementByClassName 通过类名查找UI元素
// 参数: className - 元素的类名
// 返回: 找到的元素，未找到返回 nil
func (bot *WeChatBot) FindElementByClassName(className string) *uia.IUIAutomationElement {
	variant, _ := uia.VariantFromString(className)
	cond, _ := bot.ppv.CreatePropertyCondition(uia.UIA_ClassNamePropertyId, variant)
	elem, _ := bot.root.FindFirst(uia.TreeScope_Descendants, cond)
	return elem
}

// ClickElement 点击指定的UI元素
// 参数: elem - 要点击的元素（注意：不会释放元素，由调用者负责释放）
// 返回: 错误信息，成功返回 nil
func (bot *WeChatBot) ClickElement(elem *uia.IUIAutomationElement) error {
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
// 参数: elem - 目标元素, text - 要设置的文本
// 返回: 错误信息，成功返回 nil
func (bot *WeChatBot) SetTextValue(elem *uia.IUIAutomationElement, text string) error {
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
// 参数: elem - 目标元素
// 返回: 元素的文本值
func (bot *WeChatBot) GetTextValue(elem *uia.IUIAutomationElement) string {
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

// Windows API DLL 加载
var user32 = syscall.NewLazyDLL("user32.dll")      // Windows 用户界面 API
var kernel32 = syscall.NewLazyDLL("kernel32.dll")  // Windows 内核 API
var shell32 = syscall.NewLazyDLL("shell32.dll")    // Windows Shell API
var procKeybdEvent = user32.NewProc("keybd_event") // 键盘事件函数

// TypeText 模拟键盘输入文本
// 参数: text - 要输入的文本
// 返回: 错误信息，成功返回 nil
// 注意: 只支持英文字符，中文字符会调用 PasteText
func (bot *WeChatBot) TypeText(text string) error {
	if hasChinese(text) {
		log.Println("检测到中文，使用 PasteText 方法")
		return bot.PasteText(text)
	}

	procKeybdEvent := user32.NewProc("keybd_event")
	for _, ch := range text {
		var vk uint32
		switch ch {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			vk = 0x30 + uint32(ch-'0')
		case 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z':
			vk = 0x41 + uint32(ch-'a')
		case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
			vk = 0x41 + uint32(ch-'A')
		case ' ':
			vk = 0x20
		case '\n':
			vk = 0x0D
		case '-', ':':
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

// hasChinese 检测文本中是否包含中文字符
// 参数: text - 要检测的文本
// 返回: 包含中文返回 true，否则返回 false
func hasChinese(text string) bool {
	for _, r := range text {
		if r >= 0x4e00 && r <= 0x9fff {
			return true
		}
	}
	return false
}

// PasteText 通过剪贴板粘贴文本
// 参数: text - 要粘贴的文本
// 返回: 错误信息，成功返回 nil
// 原理: 将文本放入剪贴板，然后模拟 Ctrl+V 粘贴
func (bot *WeChatBot) PasteText(text string) error {
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

	const (
		WM_KEYDOWN = 0x0100 // 按键按下消息
		WM_KEYUP   = 0x0101 // 按键抬起消息
		VK_CONTROL = 0x11   // Ctrl 键虚拟键码
		VK_V       = 0x56   // V 键虚拟键码
	)

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

// FindSiblingElements 查找同级元素
// 参数: startElem - 起始元素, filterFunc - 过滤函数
// 返回: 符合条件的元素，未找到返回 nil
func (bot *WeChatBot) FindSiblingElements(startElem *uia.IUIAutomationElement, filterFunc func(*uia.IUIAutomationElement) bool) *uia.IUIAutomationElement {
	if startElem == nil {
		return nil
	}

	walker, err := bot.ppv.CreateTreeWalker(bot.ppv.CreateTrueCondition())
	if err != nil {
		return nil
	}
	defer walker.Release()

	sibling, _ := walker.GetNextSiblingElement(startElem)
	siblingCount := 0
	for sibling != nil {
		siblingCount++
		if siblingCount > 100 {
			log.Println("警告：检查了超过100个同级元素，停止搜索")
			sibling.Release()
			break
		}

		matched := false
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("检查元素时发生 panic: %v", r)
				}
			}()
			matched = filterFunc(sibling)
		}()

		if matched {
			return sibling
		}
		sibling.Release()
		sibling, _ = walker.GetNextSiblingElement(sibling)
	}

	return nil
}

// SearchContact 搜索并选择联系人
// 参数: contactName - 联系人名称
// 返回: 错误信息，成功返回 nil
// 流程: 1. 在搜索框输入联系人名称 2. 等待搜索结果 3. 匹配并点击正确的联系人
func (bot *WeChatBot) SearchContact(contactName string) error {
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

	log.Println("查找所有Name为'文件传输助手'的元素...")

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
		log.Printf("共找到 %d 个Name为'%s'的元素，但都不符合条件（需要 ClassName=mmui::XTableCell 且 AutomationId=%s）", matchedCount, contactName, expectedAutomationId)
	} else {
		log.Printf("未找到任何Name为'%s'的元素", contactName)
	}

	return fmt.Errorf("未找到联系人: %s", contactName)
}

// SendMessage 发送消息到当前聊天窗口
// 参数: message - 要发送的消息内容
// 返回: 错误信息，成功返回 nil
// 流程: 1. 查找并点击消息输入框 2. 使用 TextPattern 选中输入框 3. 输入消息文本 4. 查找并点击发送按钮
func (bot *WeChatBot) SendMessage(message string) error {
	log.Printf("发送消息: %s", message)

	// 步骤1: 查找消息输入框
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
	defer msgEdit.Release() // 确保函数结束时释放元素

	// 记录找到的输入框属性
	name, _ := msgEdit.Get_CurrentName()
	className, _ := msgEdit.Get_CurrentClassName()
	automationId, _ := msgEdit.Get_CurrentAutomationId()
	log.Printf("找到消息输入框 - Name: %q, ClassName: %q, AutomationId: %q", name, className, automationId)

	// 点击输入框以激活焦点
	err := bot.ClickElement(msgEdit)
	if err != nil {
		return err
	}
	time.Sleep(1000 * time.Millisecond)

	// 步骤2: 使用 ValuePattern 设置文本
	log.Println("尝试使用 ValuePattern 设置文本...")
	err = bot.SetTextValue(msgEdit, message)
	if err != nil {
		log.Printf("ValuePattern 设置失败: %v，尝试使用 TypeText", err)
		// 如果 ValuePattern 失败，尝试使用 TypeText
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

	// 步骤4: 查找发送按钮
	log.Println("查找发送按钮...")

	allElements, err := bot.root.FindAll(uia.TreeScope_Descendants, bot.ppv.CreateTrueCondition())
	if err != nil {
		return fmt.Errorf("查找元素失败: %v", err)
	}
	defer allElements.Release()

	length := allElements.Get_Length()
	log.Printf("窗口中有 %d 个元素", length)

	// 遍历所有元素，查找包含"发送"文字的按钮
	var sendButton *uia.IUIAutomationElement
	for i := int32(0); i < length; i++ {
		elem, _ := allElements.GetElement(i)
		if elem != nil {
			name, _ := elem.Get_CurrentName()
			className, _ := elem.Get_CurrentClassName()

			if strings.Contains(name, "发送") {
				log.Printf("找到 Name 包含'发送'的元素 - Name: %q, ClassName: %q", name, className)
				// 确保是正确的按钮类型
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

	// 记录找到的发送按钮属性
	name, _ = sendButton.Get_CurrentName()
	className, _ = sendButton.Get_CurrentClassName()
	log.Printf("找到发送按钮 - Name: %q, ClassName: %q", name, className)

	// 点击发送按钮
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

// ReadMessages 读取当前聊天窗口中的消息列表
// 返回: 消息内容数组和可能的错误
// 功能: 从消息列表中提取所有消息文本
func (bot *WeChatBot) ReadMessages() ([]string, error) {
	log.Println("读取消息...")

	var messages []string

	// 查找消息列表元素
	list := bot.FindElementByAutomationId("chat_message_list")
	if list == nil {
		list = bot.FindElementByClassName("mmui::RecyclerListView")
	}
	if list == nil {
		return nil, fmt.Errorf("未找到消息列表")
	}
	defer list.Release()

	// 创建条件：只查找 ListItem 类型的元素
	variant, _ := uia.VariantFromString("ListItem")
	cond, _ := bot.ppv.CreatePropertyCondition(uia.UIA_ControlTypePropertyId, variant)

	// 查找所有子元素（消息项）
	children, err := list.FindAll(uia.TreeScope_Children, cond)
	if err != nil {
		return nil, err
	}
	if children != nil {
		defer children.Release()
	}

	// 遍历所有消息项，提取消息内容
	arrLen := children.Get_Length()
	for i := 0; i < int(arrLen); i++ {
		elem, _ := children.GetElement(int32(i))
		if elem != nil {
			name, _ := elem.Get_CurrentName()
			messages = append(messages, name)
			elem.Release()
		}
	}

	return messages, nil
}

// DemoSendMessage 交互式演示发送消息功能
// 功能: 通过控制台输入联系人和消息内容，然后发送
func (bot *WeChatBot) DemoSendMessage() {
	// 从控制台读取联系人名称
	fmt.Print("请输入联系人名称: ")
	var contactName string
	fmt.Scanln(&contactName)

	// 从控制台读取消息内容
	fmt.Print("请输入消息内容: ")
	var message string
	fmt.Scanln(&message)

	// 搜索联系人
	err := bot.SearchContact(contactName)
	if err != nil {
		log.Printf("搜索联系人失败: %v", err)
		return
	}
	time.Sleep(500 * time.Millisecond)

	// 发送消息
	err = bot.SendMessage(message)
	if err != nil {
		log.Printf("发送消息失败: %v", err)
		return
	}

	log.Println("消息发送成功！")
}

// DemoReadMessages 交互式演示读取消息功能
// 功能: 读取当前聊天窗口的消息并显示在控制台
func (bot *WeChatBot) DemoReadMessages() {
	// 读取消息列表
	messages, err := bot.ReadMessages()
	if err != nil {
		log.Printf("读取消息失败: %v", err)
		return
	}

	// 显示所有消息
	log.Println("=== 聊天记录 ===")
	for i, msg := range messages {
		log.Printf("[%d] %s", i+1, msg)
	}
}

func (bot *WeChatBot) DemoAutoReply() {
	fmt.Print("请输入联系人名称: ")
	var contactName string
	fmt.Scanln(&contactName)

	fmt.Print("请输入自动回复内容: ")
	var replyMsg string
	fmt.Scanln(&replyMsg)

	fmt.Print("请输入监听时长(秒): ")
	var duration int
	fmt.Scanln(&duration)

	err := bot.SearchContact(contactName)
	if err != nil {
		log.Printf("搜索联系人失败: %v", err)
		return
	}

	log.Printf("开始监听新消息，时长: %d秒...", duration)
	log.Println("按 Ctrl+C 可提前退出")

	timeout := time.After(time.Duration(duration) * time.Second)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	lastMessages := ""

	for {
		select {
		case <-timeout:
			log.Println("监听结束")
			return
		case <-ticker.C:
			messages, err := bot.ReadMessages()
			if err != nil {
				continue
			}

			currentMsgs := strings.Join(messages, "\n")
			if currentMsgs != lastMessages && len(messages) > 0 {
				newMsg := messages[len(messages)-1]
				log.Printf("收到新消息: %s", newMsg)

				time.Sleep(500 * time.Millisecond)
				err = bot.SendMessage(replyMsg)
				if err != nil {
					log.Printf("自动回复失败: %v", err)
				} else {
					log.Printf("已自动回复: %s", replyMsg)
				}

				lastMessages = currentMsgs
			}
		}
	}
}

func (bot *WeChatBot) DemoMassSend() {
	fmt.Print("请输入联系人列表(用逗号分隔): ")
	var contactsInput string
	fmt.Scanln(&contactsInput)
	contacts := strings.Split(contactsInput, ",")

	fmt.Print("请输入消息内容: ")
	var message string
	fmt.Scanln(&message)

	for i, contact := range contacts {
		contact = strings.TrimSpace(contact)
		if contact == "" {
			continue
		}

		log.Printf("[%d/%d] 发送给: %s", i+1, len(contacts), contact)

		err := bot.SearchContact(contact)
		if err != nil {
			log.Printf("搜索联系人失败: %v", err)
			continue
		}
		time.Sleep(500 * time.Millisecond)

		err = bot.SendMessage(message)
		if err != nil {
			log.Printf("发送消息失败: %v", err)
			continue
		}

		log.Printf("发送成功: %s", contact)
		time.Sleep(1 * time.Second)
	}

	log.Println("群发完成！")
}
