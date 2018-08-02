package session

import (
	"sync"
	"time"
)

type SessionsMap map[string]*RobotSession

type SessionMgr struct {
	mgrlock  sync.RWMutex
	sessions SessionsMap
}

var gSessionMgr *SessionMgr

func init() {
	gSessionMgr = &SessionMgr{
		sessions: make(SessionsMap),
	}
}

func (this *SessionMgr) getSession(id string) *RobotSession {
	if s, ok := this.sessions[id]; ok {
		return s
	}

	return nil
}

func (this *SessionMgr) GetSession(id string) *RobotSession {
	this.mgrlock.RLock()
	defer this.mgrlock.RUnlock()
	return this.getSession(id)
}

func (this *SessionMgr) NewSession(id string) *RobotSession {
	this.mgrlock.Lock()
	defer this.mgrlock.Unlock()

	if s := this.getSession(id); s != nil {
		return s
	}

	s := &RobotSession{
		id:        id,
		ctime:     time.Now().Unix(),
		stopCtrl:  make(chan bool, 1024),
		stopRecv:  make(chan bool, 1024),
		remoteMap: make(map[string]string),
	}

	this.sessions[id] = s
	return s
}

func (this *SessionMgr) FreeSession(id string) {
	s := this.GetSession(id)
	if s == nil {
		return
	}

	s.Free()
	this.mgrlock.Lock()
	defer this.mgrlock.Unlock()
	delete(this.sessions, s.Id())
}

func GetSession(id string) *RobotSession {
	return gSessionMgr.GetSession(id)
}

func DestroySession(id string) {
	gSessionMgr.FreeSession(id)
}

func CreateSession(id string) *RobotSession {
	return gSessionMgr.NewSession(id)
}
