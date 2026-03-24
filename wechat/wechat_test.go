package wechat

import (
	"testing"
)

func TestHasChinese(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"empty string", "", false},
		{"english only", "hello world", false},
		{"numbers only", "123456", false},
		{"chinese only", "你好世界", true},
		{"mixed chinese and english", "hello你好world", true},
		{"mixed chinese and numbers", "测试123", true},
		{"special chars", "!@#$%^&*()", false},
		{"chinese punctuation", "你好，世界", true},
		{"single chinese char", "中", true},
		{"japanese chars", "こんにちは", false},
		{"korean chars", "안녕하세요", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasChinese(tt.input)
			if result != tt.expected {
				t.Errorf("hasChinese(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsTimeString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"weekday monday", "星期一", true},
		{"weekday tuesday", "星期二", true},
		{"weekday wednesday", "星期三", true},
		{"weekday thursday", "星期四", true},
		{"weekday friday", "星期五", true},
		{"weekday saturday", "星期六", true},
		{"weekday sunday", "星期日", true},
		{"yesterday", "昨天", true},
		{"day before yesterday", "前天", true},
		{"just now", "刚刚", true},
		{"minutes ago", "5分钟前", true},
		{"hours ago", "2小时前", true},
		{"normal time", "12:30", true},
		{"date format", "2024-01-15", false},
		{"empty string", "", false},
		{"random text", "你好", false},
		{"partial match", "今天是星期一", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isTimeString(tt.input)
			if result != tt.expected {
				t.Errorf("isTimeString(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseMessageType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"text message", "这是一条普通消息", MsgTypeText},
		{"link message", "[链接]点击查看详情", MsgTypeLink},
		{"voice message", "[语音] 5\"", MsgTypeVoice},
		{"voice call", "[语音通话] 通话时长 02:30", MsgTypeVoiceCall},
		{"video message", "[视频] 00:15", MsgTypeVideo},
		{"file message", "[文件] document.pdf", MsgTypeFile},
		{"empty string", "", MsgTypeText},
		{"partial link text", "这是一个[链接]的描述", MsgTypeText},
		{"text with brackets", "【重要通知】", MsgTypeText},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseMessageType(tt.input)
			if result != tt.expected {
				t.Errorf("ParseMessageType(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestHasPrefix(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		prefix   string
		expected bool
	}{
		{"exact match", "hello", "hello", true},
		{"prefix match", "hello world", "hello", true},
		{"no match", "hello", "world", false},
		{"empty prefix", "hello", "", true},
		{"empty string", "", "hello", false},
		{"both empty", "", "", true},
		{"prefix longer than string", "hi", "hello", false},
		{"chinese prefix", "你好世界", "你好", true},
		{"case sensitive", "Hello", "hello", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasPrefix(tt.s, tt.prefix)
			if result != tt.expected {
				t.Errorf("hasPrefix(%q, %q) = %v, want %v", tt.s, tt.prefix, result, tt.expected)
			}
		})
	}
}

func TestNewSessionManager(t *testing.T) {
	sm := NewSessionManager()
	if sm == nil {
		t.Fatal("NewSessionManager() returned nil")
	}
	if sm.sessions == nil {
		t.Error("sessions map is nil")
	}
	if sm.lastContent == nil {
		t.Error("lastContent map is nil")
	}
}

func TestSessionManager_GetSession(t *testing.T) {
	sm := NewSessionManager()

	session := &Session{
		Sender:       "测试用户",
		Content:      "测试消息",
		Time:         "12:30",
		AutomationId: "session_item_测试用户",
		MsgType:      MsgTypeText,
	}
	sm.UpdateSession(session)

	tests := []struct {
		name       string
		sender     string
		wantNil    bool
		wantSender string
	}{
		{"existing session", "测试用户", false, "测试用户"},
		{"non-existing session", "不存在的用户", true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sm.GetSession(tt.sender)
			if tt.wantNil {
				if result != nil {
					t.Errorf("GetSession(%q) should return nil", tt.sender)
				}
			} else {
				if result == nil {
					t.Errorf("GetSession(%q) returned nil", tt.sender)
					return
				}
				if result.Sender != tt.wantSender {
					t.Errorf("GetSession(%q).Sender = %q, want %q", tt.sender, result.Sender, tt.wantSender)
				}
			}
		})
	}
}

func TestSessionManager_GetAllSessions(t *testing.T) {
	sm := NewSessionManager()

	sessions := sm.GetAllSessions()
	if len(sessions) != 0 {
		t.Errorf("new SessionManager should have 0 sessions, got %d", len(sessions))
	}

	sm.UpdateSession(&Session{Sender: "用户1", Content: "消息1"})
	sm.UpdateSession(&Session{Sender: "用户2", Content: "消息2"})
	sm.UpdateSession(&Session{Sender: "用户3", Content: "消息3"})

	sessions = sm.GetAllSessions()
	if len(sessions) != 3 {
		t.Errorf("GetAllSessions() returned %d sessions, want 3", len(sessions))
	}
}

func TestSessionManager_UpdateSession(t *testing.T) {
	sm := NewSessionManager()

	session1 := &Session{
		Sender:  "测试用户",
		Content: "第一条消息",
	}
	sm.UpdateSession(session1)

	result := sm.GetSession("测试用户")
	if result == nil || result.Content != "第一条消息" {
		t.Errorf("UpdateSession failed to add session")
	}

	session2 := &Session{
		Sender:  "测试用户",
		Content: "第二条消息",
	}
	sm.UpdateSession(session2)

	result = sm.GetSession("测试用户")
	if result == nil || result.Content != "第二条消息" {
		t.Errorf("UpdateSession failed to update session")
	}
}

func TestSessionManager_LastContent(t *testing.T) {
	sm := NewSessionManager()

	if sm.HasSender("测试用户") {
		t.Error("HasSender should return false for new sender")
	}

	sm.SetLastContent("测试用户", "旧消息")

	if !sm.HasSender("测试用户") {
		t.Error("HasSender should return true after SetLastContent")
	}

	content := sm.GetLastContent("测试用户")
	if content != "旧消息" {
		t.Errorf("GetLastContent() = %q, want %q", content, "旧消息")
	}

	content = sm.GetLastContent("不存在的用户")
	if content != "" {
		t.Errorf("GetLastContent() for non-existing user should return empty string, got %q", content)
	}
}

func TestSessionManager_Clear(t *testing.T) {
	sm := NewSessionManager()

	sm.UpdateSession(&Session{Sender: "用户1"})
	sm.UpdateSession(&Session{Sender: "用户2"})
	sm.SetLastContent("用户1", "消息1")

	sm.Clear()

	if len(sm.GetAllSessions()) != 0 {
		t.Error("Clear() failed to clear sessions")
	}
	if len(sm.lastContent) != 0 {
		t.Error("Clear() failed to clear lastContent")
	}
}

func TestSessionManager_NilSession(t *testing.T) {
	sm := NewSessionManager()

	sm.UpdateSession(nil)

	if len(sm.GetAllSessions()) != 0 {
		t.Error("UpdateSession(nil) should not add any session")
	}
}

func TestSession_Fields(t *testing.T) {
	session := &Session{
		Sender:       "测试发送者",
		Content:      "测试内容",
		Time:         "12:30",
		AutomationId: "test_id",
		IsSelf:       true,
		MsgType:      MsgTypeText,
	}

	if session.Sender != "测试发送者" {
		t.Errorf("Sender = %q, want %q", session.Sender, "测试发送者")
	}
	if session.Content != "测试内容" {
		t.Errorf("Content = %q, want %q", session.Content, "测试内容")
	}
	if session.Time != "12:30" {
		t.Errorf("Time = %q, want %q", session.Time, "12:30")
	}
	if session.AutomationId != "test_id" {
		t.Errorf("AutomationId = %q, want %q", session.AutomationId, "test_id")
	}
	if !session.IsSelf {
		t.Errorf("IsSelf = %v, want %v", session.IsSelf, true)
	}
	if session.MsgType != MsgTypeText {
		t.Errorf("MsgType = %q, want %q", session.MsgType, MsgTypeText)
	}
}

func TestNewMessage_Fields(t *testing.T) {
	msg := &NewMessage{
		Sender:  "发送者",
		Content: "消息内容",
		Time:    "刚刚",
		IsSelf:  false,
		MsgType: MsgTypeLink,
	}

	if msg.Sender != "发送者" {
		t.Errorf("Sender = %q, want %q", msg.Sender, "发送者")
	}
	if msg.Content != "消息内容" {
		t.Errorf("Content = %q, want %q", msg.Content, "消息内容")
	}
	if msg.Time != "刚刚" {
		t.Errorf("Time = %q, want %q", msg.Time, "刚刚")
	}
	if msg.IsSelf {
		t.Errorf("IsSelf = %v, want %v", msg.IsSelf, false)
	}
	if msg.MsgType != MsgTypeLink {
		t.Errorf("MsgType = %q, want %q", msg.MsgType, MsgTypeLink)
	}
}

func TestMessageTypeConstants(t *testing.T) {
	if MsgTypeText != "文本" {
		t.Errorf("MsgTypeText = %q, want %q", MsgTypeText, "文本")
	}
	if MsgTypeLink != "链接" {
		t.Errorf("MsgTypeLink = %q, want %q", MsgTypeLink, "链接")
	}
	if MsgTypeVoice != "语音" {
		t.Errorf("MsgTypeVoice = %q, want %q", MsgTypeVoice, "语音")
	}
	if MsgTypeVoiceCall != "语音通话" {
		t.Errorf("MsgTypeVoiceCall = %q, want %q", MsgTypeVoiceCall, "语音通话")
	}
	if MsgTypeVideo != "视频" {
		t.Errorf("MsgTypeVideo = %q, want %q", MsgTypeVideo, "视频")
	}
	if MsgTypeFile != "文件" {
		t.Errorf("MsgTypeFile = %q, want %q", MsgTypeFile, "文件")
	}
}
