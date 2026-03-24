// Package wechat 提供微信自动化操作功能
// 基于 Windows UI Automation 实现微信客户端的自动化控制
// 支持功能：搜索联系人、发送消息、读取会话列表、新消息监控等
package wechat

// Session 会话信息结构体
// 用于表示微信左侧会话列表中的单个会话项
type Session struct {
	// Sender 发送者名称（联系人或群组名称）
	Sender string
	// Content 消息内容预览（最新一条消息的内容摘要）
	Content string
	// Time 消息时间（如：星期五、昨天、12:30 等）
	Time string
	// AutomationId UI元素的自动化ID（格式：session_item_发送者名称）
	AutomationId string
	// IsSelf 是否为自己发送的消息
	IsSelf bool
	// MsgType 消息类型（文本、链接、语音、语音通话、视频等）
	MsgType string
}

// NewMessage 新消息通知结构体
// 当检测到新消息时，通过此结构体传递消息信息
type NewMessage struct {
	// Sender 发送者名称
	Sender string
	// Content 新消息内容
	Content string
	// Time 消息时间
	Time string
	// IsSelf 是否为自己发送的消息
	IsSelf bool
	// MsgType 消息类型
	MsgType string
}

// SessionManager 会话管理器
// 用于跟踪会话状态变化，检测新消息
type SessionManager struct {
	// sessions 当前会话列表（key为发送者名称）
	sessions map[string]*Session
	// lastContent 上次记录的消息内容（用于比较检测新消息）
	lastContent map[string]string
}

// NewSessionManager 创建会话管理器实例
// 返回初始化好的 SessionManager 指针
func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions:    make(map[string]*Session),
		lastContent: make(map[string]string),
	}
}

// GetSession 获取指定发送者的会话信息
// 参数: sender - 发送者名称
// 返回: 会话信息，不存在返回 nil
func (sm *SessionManager) GetSession(sender string) *Session {
	return sm.sessions[sender]
}

// GetAllSessions 获取所有会话信息
// 返回: 会话列表
func (sm *SessionManager) GetAllSessions() []*Session {
	sessions := make([]*Session, 0, len(sm.sessions))
	for _, s := range sm.sessions {
		sessions = append(sessions, s)
	}
	return sessions
}

// UpdateSession 更新会话信息
// 参数: session - 会话信息
func (sm *SessionManager) UpdateSession(session *Session) {
	if session == nil {
		return
	}
	sm.sessions[session.Sender] = session
}

// GetLastContent 获取指定发送者的上次消息内容
// 参数: sender - 发送者名称
// 返回: 消息内容，不存在返回空字符串
func (sm *SessionManager) GetLastContent(sender string) string {
	return sm.lastContent[sender]
}

// SetLastContent 设置指定发送者的上次消息内容
// 参数: sender - 发送者名称, content - 消息内容
func (sm *SessionManager) SetLastContent(sender, content string) {
	sm.lastContent[sender] = content
}

// HasSender 检查是否存在指定发送者的记录
// 参数: sender - 发送者名称
// 返回: 存在返回 true
func (sm *SessionManager) HasSender(sender string) bool {
	_, exists := sm.lastContent[sender]
	return exists
}

// Clear 清空所有会话记录
func (sm *SessionManager) Clear() {
	sm.sessions = make(map[string]*Session)
	sm.lastContent = make(map[string]string)
}

// MessageType 消息类型常量定义
const (
	MsgTypeText      = "文本"      // 文本消息
	MsgTypeLink      = "链接"      // 链接消息
	MsgTypeVoice     = "语音"      // 语音消息
	MsgTypeVoiceCall = "语音通话" // 语音通话
	MsgTypeVideo     = "视频"      // 视频消息
)

// ParseMessageType 根据消息内容判断消息类型
// 参数: content - 消息内容
// 返回: 消息类型字符串
func ParseMessageType(content string) string {
	switch {
	case hasPrefix(content, "[链接]"):
		return MsgTypeLink
	case hasPrefix(content, "[语音]"):
		return MsgTypeVoice
	case hasPrefix(content, "[语音通话]"):
		return MsgTypeVoiceCall
	case hasPrefix(content, "[视频]"):
		return MsgTypeVideo
	default:
		return MsgTypeText
	}
}

// hasPrefix 检查字符串是否有指定前缀
func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}
