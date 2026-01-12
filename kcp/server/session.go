package server

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/xtaci/kcp-go/v5"
)

// Session represents a KCP session with activity tracking.
type Session struct {
	Conn         *kcp.UDPSession
	LastActivity atomic.Int64
	ID           string
	Closed       atomic.Bool
}

// SessionManager manages multiple KCP sessions.
type SessionManager struct {
	sessions sync.Map
	timeout  time.Duration
}

// NewSessionManager creates a new SessionManager with the given session timeout.
func NewSessionManager(timeout time.Duration) *SessionManager {
	return &SessionManager{
		sessions: sync.Map{},
		timeout:  timeout,
	}
}

// Add creates a new session for the given connection and stores it.
func (sm *SessionManager) Add(conn *kcp.UDPSession) *Session {
	id := conn.RemoteAddr().String()
	s := &Session{
		Conn: conn,
		ID:   id,
	}
	s.LastActivity.Store(time.Now().Unix())
	sm.sessions.Store(id, s)
	return s
}

// Get retrieves the session with the given id.
func (sm *SessionManager) Get(id string) (*Session, bool) {
	if s, ok := sm.sessions.Load(id); ok {
		return s.(*Session), true
	}
	return nil, false
}

// Remove closes and removes the session with the given id.
func (sm *SessionManager) Remove(id string) {
	if s, ok := sm.sessions.LoadAndDelete(id); ok {
		if session, ok := s.(*Session); ok {
			session.Closed.CompareAndSwap(false, true)
			session.Conn.Close()
		}
	}
}




