package common

import (
	"log/slog"
	"sync"
	"time"

	"github.com/goburrow/serial"
)

const (
	serialIdleTimeout = 60 * time.Second
)

type SerialClient struct {
	Config      *serial.Config
	IdleTimeout time.Duration

	mu           sync.Mutex
	port         serial.Port
	closeTimer   *time.Timer
	lastActivity time.Time
}

func NewSerialClient(config *serial.Config) SerialClient {
	return SerialClient{
		Config:      config,
		IdleTimeout: serialIdleTimeout,
	}
}

// Send 发送数据到串口，并获取响应数据
func (s *SerialClient) Send(requestData []byte, dataReader func(port serial.Port) error) (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err = s.connect(); err != nil {
		return
	}
	s.setCloseTimer()
	s.lastActivity = time.Now()
	// 发送数据
	if _, err = s.port.Write(requestData); err != nil {
		_ = s.close()
		return
	}
	err = dataReader(s.port)
	if err != nil {
		_ = s.close()
		s.drain()
		return
	}
	return
}

// Connect 封装给外部使用
func (s *SerialClient) Connect() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.connect()
}

// 如果不存在连接，创建连接
func (s *SerialClient) connect() error {
	if s.port == nil {
		port, err := serial.Open(s.Config)
		if err != nil {
			return err
		}
		s.port = port
	}
	return nil
}

// Close 封装给外部使用
func (s *SerialClient) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.close()
}

// 关闭连接
func (s *SerialClient) close() (err error) {
	if s.port != nil {
		err = s.port.Close()
		s.port = nil
	}
	return
}

// drain 清空当前连接中阻塞无效的数据
func (s *SerialClient) drain() {
	if s.port == nil {
		return
	}
	buf := make([]byte, 1024)
	_, _ = s.port.Read(buf)
	return
}

// 启动闲置连接检测
func (s *SerialClient) setCloseTimer() {
	if s.IdleTimeout <= 0 {
		return
	}
	if s.closeTimer == nil {
		s.closeTimer = time.AfterFunc(s.IdleTimeout, s.closeIdle)
	} else {
		s.closeTimer.Reset(s.IdleTimeout)
	}
}

// closeIdle 如果闲置时间超过 IdleTimeout，则关闭连接
func (s *SerialClient) closeIdle() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.IdleTimeout <= 0 {
		return
	}
	idle := time.Since(s.lastActivity)
	if idle >= s.IdleTimeout {
		slog.Info("serial client: closing connection due to idle timeout", "idle time", idle)
		_ = s.close()
	}
}
