package kcp

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/xtaci/kcp-go/v5"
)

// Server represents a KCP server that manages incoming connections and sessions.
type Server struct {
	listener  *kcp.Listener
	cfg       *Config
	handler   Handler
	sessions  map[string]*Session
	sessionMu sync.RWMutex

	stopChan chan struct{}
	wg       sync.WaitGroup
}

// Handler defines the interface for handling KCP connections.
type Handler interface {
	HandleConn(conn *Conn, session *Session) error
}

// HandlerFunc is an adapter to allow the use of ordinary functions as KCP handlers.
type HandlerFunc func(conn *Conn, session *Session) error

// HandleConn implements the Handler interface for HandlerFunc.
func (f HandlerFunc) HandleConn(conn *Conn, session *Session) error {
	return f(conn, session)
}

// NewServer creates a new KCP server with the given configuration.
func NewServer(cfg *Config) *Server {
	if cfg == nil {
		cfg = DefaultConfig()
	}
	return &Server{
		cfg:      cfg,
		sessions: make(map[string]*Session),
		stopChan: make(chan struct{}),
	}
}

// Listen starts listening for incoming KCP connections on the specified address.
func (s *Server) Listen(addr string) error {
	listener, err := kcp.ListenWithOptions(addr, nil, 10, 3)
	if err != nil {
		return fmt.Errorf("failed liston on %s: %w", addr, err)
	}
	s.listener = listener
	return nil
}

func (s *Server) Serve(handler Handler) error {
	s.handler = handler

	s.wg.Add(1)
	go s.heartbeat()

	for {
		select {
		case <-s.stopChan:
			return nil
		default:
		}
		conn, err := s.listener.AcceptKCP()
		if err != nil {
			select {
			case <-s.stopChan:
				return nil
			default:
				continue
			}
		}
		s.wg.Add(1)
		go s.handleConn(conn)
	}
}

// ServeFunc is a convenience method to serve using a function as the handler.
func (s *Server) ServeFunc(fn func(conn *Conn, session *Session) error) error {
	return s.Serve(HandlerFunc(fn))
}

func (s *Server) handleConn(kcpConn *kcp.UDPSession) {
	defer s.wg.Done()

	conn := &Conn{
		conn: kcpConn,
		cfg:  s.cfg,
	}
	conn.applyConfig()
	defer conn.Close()

	session := &Session{
		LastActivity:  time.Now(),
		LastHeartbeat: time.Now(),
		Conn:          conn,
		IsAlive:       true,
		ID:            kcpConn.RemoteAddr().String(),
	}
	s.addSession(session)
	defer s.removeSession(session.ID)
	if s.handler != nil {
		if err := s.handler.HandleConn(conn, session); err != nil {
			// Handle error (log, etc.)
		}
	}
}

// addSession adds a new session to the server's session map.
func (s *Server) addSession(session *Session) {
	s.sessionMu.Lock()
	defer s.sessionMu.Unlock()
	s.sessions[session.ID] = session
}

// removeSession removes a session by its ID.
func (s *Server) removeSession(sessionID string) {
	s.sessionMu.Lock()
	defer s.sessionMu.Unlock()
	delete(s.sessions, sessionID)
}

// GetSession retrieves a session by its ID.
func (s *Server) GetSession(sessionID string) *Session {
	s.sessionMu.RLock()
	defer s.sessionMu.RUnlock()
	return s.sessions[sessionID]
}

// SessionCount returns the current number of active sessions.
func (s *Server) SessionCount() int {
	s.sessionMu.RLock()
	defer s.sessionMu.RUnlock()
	return len(s.sessions)
}

// heartbeat periodically checks sessions for heartbeats and removes inactive ones.
func (s *Server) heartbeat() {
	defer s.wg.Done()
	ticker := time.NewTicker(s.cfg.HeartbeatTimeout / 2)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopChan:
			return
		case <-ticker.C:
			s.checkHeartbeats()
		}
	}
}

// checkHeartbeats checks all sessions for heartbeat timeouts and marks inactive sessions.
func (s *Server) checkHeartbeats() {
	now := time.Now()
	timeout := s.cfg.HeartbeatTimeout

	s.sessionMu.Lock()
	var timeoutSessions []string
	for id, session := range s.sessions {
		session.mu.RLock()
		if session.IsAlive && now.Sub(session.LastHeartbeat) > timeout {
			timeoutSessions = append(timeoutSessions, id)
		}
		session.mu.RUnlock()
	}
	s.sessionMu.Unlock()
	for _, id := range timeoutSessions {
		s.sessionMu.Lock()
		if session, exists := s.sessions[id]; exists {
			session.mu.Lock()
			session.IsAlive = false
			session.mu.Unlock()
		}
		s.sessionMu.Unlock()
	}
}

// Stop gracefully stops the server, closing all connections and waiting for ongoing operations to finish.
func (s *Server) Stop(ctx context.Context) error {
	close(s.stopChan)

	if s.listener != nil {
		s.listener.Close()
	}

	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}

// ListenAndServe is a convenience method to create, listen, and serve a KCP server.
func ListenAndServe(addr string, cfg *Config, handler Handler) error {
	server := NewServer(cfg)
	if err := server.Listen(addr); err != nil {
		return err
	}
	return server.Serve(handler)
}

// ListenAndServeFunc is a convenience method to listen and serve using a function as the handler.
func ListenAndServeFunc(addr string, cfg *Config, fn func(conn *Conn, session *Session) error) error {
	return ListenAndServe(addr, cfg, HandlerFunc(fn))
}
