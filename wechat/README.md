# WeChat 微信自动化包

基于 Windows UI Automation 实现的微信客户端自动化控制包。

## 功能特性

- ✅ 搜索联系人
- ✅ 发送消息
- ✅ 读取会话列表
- ✅ 新消息监控
- ✅ 消息类型识别（文本、链接、语音、语音通话、视频）

## 环境要求

- Windows 操作系统
- Go 1.16+
- 微信客户端已打开

## 安装

```bash
go get github.com/yuan71058/go-element/wechat
```

## 快速开始

### 1. 创建机器人实例

```go
package main

import (
    "log"
    "time"

    "github.com/yuan71058/go-element/wechat"
)

func main() {
    // 创建微信机器人实例
    bot, err := wechat.NewBot()
    if err != nil {
        log.Fatal(err)
    }
    // 确保在程序结束时释放资源
    defer bot.Close()

    // ... 其他操作
}
```

### 2. 发送消息

```go
// 搜索联系人
err = bot.SearchContact("文件传输助手")
if err != nil {
    log.Fatal(err)
}

// 等待进入聊天窗口
time.Sleep(2 * time.Second)

// 发送消息
err = bot.SendMessage("这是一条测试消息")
if err != nil {
    log.Fatal(err)
}
```

### 3. 读取会话列表

```go
// 读取会话列表
sessions, err := bot.ReadSessionList()
if err != nil {
    log.Fatal(err)
}

// 遍历会话
for _, session := range sessions {
    fmt.Printf("发送者: %s\n", session.Sender)
    fmt.Printf("内容: %s\n", session.Content)
    fmt.Printf("时间: %s\n", session.Time)
    fmt.Printf("类型: %s\n", session.MsgType)
    fmt.Println("---")
}
```

### 4. 监控新消息

```go
// 方式一：手动轮询
for {
    newMsgs, err := bot.CheckNewMessages()
    if err != nil {
        log.Println(err)
        time.Sleep(3 * time.Second)
        continue
    }

    for _, msg := range newMsgs {
        fmt.Printf("新消息 - %s: %s\n", msg.Sender, msg.Content)
    }
    time.Sleep(3 * time.Second)
}

// 方式二：使用监控器（推荐）
go bot.StartMessageMonitor(3*time.Second, func(msg *wechat.NewMessage) {
    fmt.Printf("收到新消息 - %s: %s\n", msg.Sender, msg.Content)
})
```

## API 文档

### 主要类型

#### Session - 会话信息

| 字段 | 类型 | 说明 |
|------|------|------|
| Sender | string | 发送者名称 |
| Content | string | 消息内容预览 |
| Time | string | 消息时间 |
| AutomationId | string | UI元素ID |
| IsSelf | bool | 是否自己发送 |
| MsgType | string | 消息类型 |

#### NewMessage - 新消息通知

| 字段 | 类型 | 说明 |
|------|------|------|
| Sender | string | 发送者名称 |
| Content | string | 新消息内容 |
| Time | string | 消息时间 |
| IsSelf | bool | 是否自己发送 |
| MsgType | string | 消息类型 |

### 消息类型常量

```go
const (
    MsgTypeText      = "文本"      // 文本消息
    MsgTypeLink      = "链接"      // 链接消息
    MsgTypeVoice     = "语音"      // 语音消息
    MsgTypeVoiceCall = "语音通话" // 语音通话
    MsgTypeVideo     = "视频"      // 视频消息
)
```

### 主要方法

#### Bot 创建与销毁

| 方法 | 说明 |
|------|------|
| `NewBot() (*Bot, error)` | 创建机器人实例 |
| `NewBotWithLog(path string) (*Bot, error)` | 创建机器人实例（带日志） |
| `bot.Close()` | 关闭机器人，释放资源 |

#### 联系人与消息

| 方法 | 说明 |
|------|------|
| `bot.SearchContact(name string) error` | 搜索并选择联系人 |
| `bot.SendMessage(msg string) error` | 发送消息 |

#### 会话列表

| 方法 | 说明 |
|------|------|
| `bot.ReadSessionList() ([]*Session, error)` | 读取会话列表 |
| `bot.CheckNewMessages() ([]*NewMessage, error)` | 检测新消息 |
| `bot.StartMessageMonitor(interval, callback)` | 启动消息监控 |

## 完整示例

```go
package main

import (
    "fmt"
    "log"
    "time"

    "github.com/yuan71058/go-element/wechat"
)

func main() {
    // 创建机器人
    bot, err := wechat.NewBotWithLog("wechat.log")
    if err != nil {
        log.Fatal(err)
    }
    defer bot.Close()

    // 读取并打印会话列表
    fmt.Println("=== 会话列表 ===")
    sessions, err := bot.ReadSessionList()
    if err != nil {
        log.Fatal(err)
    }
    for i, s := range sessions {
        fmt.Printf("[%d] %s: %s (%s)\n", i+1, s.Sender, s.Content, s.MsgType)
    }

    // 发送测试消息
    fmt.Println("\n=== 发送消息 ===")
    err = bot.SearchContact("文件传输助手")
    if err != nil {
        log.Fatal(err)
    }
    time.Sleep(2 * time.Second)

    msg := fmt.Sprintf("测试消息 - %s", time.Now().Format("2006-01-02 15:04:05"))
    err = bot.SendMessage(msg)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("消息发送成功！")

    // 启动消息监控
    fmt.Println("\n=== 启动消息监控 ===")
    fmt.Println("按 Ctrl+C 退出...")

    go bot.StartMessageMonitor(3*time.Second, func(msg *wechat.NewMessage) {
        fmt.Printf("[新消息] %s: %s (%s)\n", msg.Sender, msg.Content, msg.MsgType)
    })

    // 保持程序运行
    select {}
}
```

## 注意事项

1. **微信窗口必须打开**：程序运行前请确保微信客户端已打开并登录
2. **窗口会被激活**：执行操作时微信窗口会被置于前台
3. **中文输入**：中文消息通过剪贴板粘贴方式输入
4. **消息监控**：首次调用 `CheckNewMessages` 不会返回新消息，仅初始化状态
5. **资源释放**：使用完毕后务必调用 `bot.Close()` 释放资源

## 错误处理

```go
bot, err := wechat.NewBot()
if err != nil {
    // 可能的错误：
    // - "未找到微信窗口，请确保微信已打开"
    // - "COM 初始化失败"
    // - "无法创建 UIAutomation"
    log.Fatal(err)
}
```

## 依赖

- [github.com/yuan71058/go-element](https://github.com/yuan71058/go-element) - Windows UI Automation 库

## 许可证

MIT License
