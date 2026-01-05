package common

import (
	"errors"
	"log/slog"
	"net"
	"sync"
	"time"
)

const (
	tcpTimeout     = 10 * time.Second
	tcpIdleTimeout = 60 * time.Second
)

type TCPClient struct {
	Address     string
	Timeout     time.Duration
	IdleTimeout time.Duration

	mu           sync.Mutex
	conn         net.Conn
	closeTimer   *time.Timer
	lastActivity time.Time
}

func NewTCPClient(address string) TCPClient {
	return TCPClient{
		Address:     address,
		Timeout:     tcpTimeout,
		IdleTimeout: tcpIdleTimeout,
	}
}

// Send 发送数据到服务器，并获取响应数据
func (t *TCPClient) Send(requestData []byte, dataReader func(conn net.Conn) error) (err error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if err = t.connect(); err != nil {
		return
	}
	t.setCloseTimer()
	t.lastActivity = time.Now()
	// 设置读写超时时间
	var timeout time.Time
	if t.Timeout > 0 {
		timeout = t.lastActivity.Add(t.Timeout)
	}
	if err = t.conn.SetDeadline(timeout); err != nil {
		return
	}
	// 发送数据
	if _, err = t.conn.Write(requestData); err != nil {
		return
	}
	err = dataReader(t.conn)
	if err != nil {
		t.drain()
		return
	}
	return
}

// Connect 封装给外部使用
func (t *TCPClient) Connect() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.connect()
}

// 如果不存在连接，创建连接
func (t *TCPClient) connect() error {
	if t.conn == nil {
		dialer := net.Dialer{Timeout: t.Timeout}
		conn, err := dialer.Dial("tcp", t.Address)
		if err != nil {
			return err
		}
		t.conn = conn
	}
	return nil
}

// Close 封装给外部使用
func (t *TCPClient) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.close()
}

// 关闭连接
func (t *TCPClient) close() (err error) {
	if t.conn != nil {
		err = t.conn.Close()
		t.conn = nil
	}
	return
}

// flush 清空当前连接中阻塞无效的数据
func (t *TCPClient) drain() {
	if err := t.conn.SetReadDeadline(time.Now()); err != nil {
		return
	}
	buf := make([]byte, 1024)
	if _, err := t.conn.Read(buf); err != nil {
		// Ignore timeout error
		var netError net.Error
		if errors.As(err, &netError) && netError.Timeout() {
			err = nil
		}
	}
	return
}

// 启动闲置连接检测
func (t *TCPClient) setCloseTimer() {
	if t.IdleTimeout <= 0 {
		return
	}
	if t.closeTimer == nil {
		t.closeTimer = time.AfterFunc(t.IdleTimeout, t.closeIdle)
	} else {
		t.closeTimer.Reset(t.IdleTimeout)
	}
}

// closeIdle  如果闲置时间超过 IdleTimeout ，则关闭连接
func (t *TCPClient) closeIdle() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.IdleTimeout <= 0 {
		return
	}
	idle := time.Now().Sub(t.lastActivity)
	if idle >= t.IdleTimeout {
		slog.Info("tcp client: closing connection due to idle timeout: %v", idle)
		_ = t.close()
	}
}
