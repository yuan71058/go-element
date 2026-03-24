// 微信自动化机器人示例程序
// 演示如何使用 wechat 包进行微信自动化操作
package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yuan71058/go-element/wechat"
)

func main() {
	// 创建微信机器人实例（带日志文件）
	bot, err := wechat.NewBotWithLog("wechat_demo.log")
	if err != nil {
		log.Fatalf("初始化失败: %v", err)
	}
	defer bot.Close()

	// 读取并打印会话列表
	fmt.Println("=== 读取会话列表 ===")
	sessions, err := bot.ReadSessionList()
	if err != nil {
		log.Fatalf("读取会话列表失败: %v", err)
	}

	for i, s := range sessions {
		fmt.Printf("[%d] 发送者: %s\n", i+1, s.Sender)
		fmt.Printf("    内容: %s\n", s.Content)
		fmt.Printf("    时间: %s\n", s.Time)
		fmt.Printf("    类型: %s\n", s.MsgType)
		fmt.Println("    ---")
	}
	fmt.Printf("共读取 %d 个会话\n\n", len(sessions))

	// 发送测试消息
	fmt.Println("=== 发送消息测试 ===")
	contactName := "文件传输助手"
	err = bot.SearchContact(contactName)
	if err != nil {
		log.Printf("搜索联系人失败: %v", err)
	} else {
		time.Sleep(2 * time.Second)
		message := fmt.Sprintf("测试消息 - %s", time.Now().Format("2006-01-02 15:04:05"))
		err = bot.SendMessage(message)
		if err != nil {
			log.Printf("发送消息失败: %v", err)
		} else {
			fmt.Printf("消息发送成功: %s -> %s\n", contactName, message)
		}
	}

	// 启动消息监控
	fmt.Println("\n=== 启动消息监控 ===")
	fmt.Println("按 Ctrl+C 退出程序...")

	go bot.StartMessageMonitor(3*time.Second, func(msg *wechat.NewMessage) {
		fmt.Printf("[新消息] %s: %s (%s)\n", msg.Sender, msg.Content, msg.MsgType)
	})

	// 等待退出信号
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	fmt.Println("\n程序退出")
}
