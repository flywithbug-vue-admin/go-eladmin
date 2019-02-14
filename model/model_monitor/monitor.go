package model_monitor

import (
	"strings"
	"sync"
	"time"
)

const (
	TimeLayout = "2006-01-02 15:04:05"
)

var (
	visitApiPool *sync.Pool
	visitUIdPool *sync.Pool
)

func init() {
	visitApiPool = &sync.Pool{New: func() interface{} {
		return &VisitApi{}
	}}
	visitUIdPool = &sync.Pool{New: func() interface{} {
		return &VisitUId{}
	}}
}

func (l Log) AddMonitorInfo() {
	if !strings.HasPrefix(l.Path, "/api") {
		return
	}
	if strings.HasPrefix(l.Path, "/api/log") && (l.ResponseCode == 200 || l.StatusCode == 204) {
		return
	}
	visitUid := visitUIdPool.Get().(*VisitUId)
	vApi := visitApiPool.Get().(*VisitApi)
	l.Insert()
	timeF := time.Now().Format(TimeLayout)
	if l.Latency > 0 {
		if len(l.UUID) > 0 {
			visitUid.TimeDate = timeF[:10]
			visitUid.UUID = l.UUID
			visitUid.UserId = l.UserId
			visitUid.ClientIp = l.ClientIp
			visitUid.IncrementVisitUId()
		}
		if len(l.Path) > 0 {
			vApi.TimeDate = timeF[:10]
			vApi.Method = l.Method
			vApi.Path = l.Path
			vApi.Para = l.Para
			vApi.IncrementVisitApi()
		}
	}
	visitUIdPool.Put(visitUid)
	visitApiPool.Put(vApi)
}
