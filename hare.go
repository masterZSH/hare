package hare

import (
	"bufio"
	"net"
)

const (
	connectionHost = "localhost"
	connectionType = "tcp"
)

// Listener tools for socket listening
// 包装net.Listener
// HasNewMessages 是否有新消息
// GetMessage 获取消息
// Stop 中断连接
type Listener struct {
	SocketListener net.Listener
	HasNewMessages func() bool
	GetMessage     func() string
	Stop           func()
}

// MessageManager manages message storage
type MessageManager struct {
	HasNewMessages bool
	Message        string
}

// 启动监听  用running控制停止
func listening(listener Listener, messageManager *MessageManager, running *bool) error {
	for *running {
		c, _ := listener.SocketListener.Accept()
		// \n分隔消息
		message, _ := bufio.NewReader(c).ReadString('\n')
		messageManager.Message = message
		messageManager.HasNewMessages = true
	}
	listener.SocketListener.Close()
	return nil
}

// Listen to socket port
func Listen(port string) (Listener, error) {
	var err error
	var listener Listener
	var messageManager MessageManager

	listener.SocketListener, err = net.Listen(connectionType, connectionHost+":"+port)
	if err != nil {
		return listener, err
	}

	// GetMessage return the last message
	listener.GetMessage = func() string {
		messageManager.HasNewMessages = false
		return messageManager.Message
	}

	// HasNewMessages return if there's new messages on socket
	listener.HasNewMessages = func() bool {
		return messageManager.HasNewMessages
	}

	running := true
	// Stop the listener
	listener.Stop = func() {
		running = false
	}

	go listening(listener, &messageManager, &running)

	return listener, nil
}

// Send message to socket port
func Send(port, message string) error {
	conn, err := net.Dial(connectionType, connectionHost+":"+port)
	if err != nil {
		return err
	}
	defer conn.Close()

	conn.Write([]byte(message))
	return nil
}

// 总体代码比较简单
