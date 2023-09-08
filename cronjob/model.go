package cronjob

import (
	"github.com/robfig/cron/v3"
	"github.com/yunduansing/gtools/logger"
	"go.uber.org/zap"
	"sync"
)

type TaskConfig struct {
	Code   string                       `gorm:"primarykey;" json:"code"` //英文编码
	Cmd    int                          `gorm:"type:int;" json:"cmd"`    //数字编码
	Name   string                       `gorm:"type:varchar(100)" json:"name"`
	Value  string                       `gorm:"type:varchar(100)" json:"value"`         //执行周期
	Remark string                       `gorm:"type:varchar(200)" json:"remark"`        //备注
	State  int                          `gorm:"type:tinyint(4);default:1" json:"state"` //0-禁用、1-启用
	Param  string                       `gorm:"type:text" json:"param"`                 //参数，Json
	Do     func() (cron.EntryID, error) `gorm:"-" json:"-"`                             //执行逻辑
}

type TaskEvent struct {
	Cmd       int           `json:"cmd"`
	Code      string        `json:"code"`
	Value     string        `json:"value"`     //当operation=4时，必须有值
	operation TaskEventType `json:"operation"` //1-启用、2-禁用、3-执行一次、4-修改执行周期
	Remark    string        `json:"remark"`
}

func (t TaskConfig) run() (cron.EntryID, error) {
	return t.Do()
}

type TaskEventType int

const (
	_           TaskEventType = iota
	EnableType                //启用
	DisableType               //禁用
	RunOnceType               //执行一次
	UpdateType                //修改任务
)

type Task struct {
	tasks     map[string]TaskConfig
	taskIdMap map[string]cron.EntryID
	sync.Mutex
	sync.Once
	serviceName string
}

func (t Task) run() {
	for _, v := range t.tasks {
		if v.Do == nil || v.State == 0 {
			continue
		}
		eid, err := v.Do()
		if err != nil {
			logger.Error("启动任务失败", zap.Any("task", v), zap.Error(err))
			continue
		}
		t.taskIdMap[v.Code] = eid
	}
}

func (t Task) addAndRun(task TaskConfig) error {
	t.Lock()
	defer t.Unlock()
	if _, ok := t.tasks[task.Code]; !ok {
		eid, err := task.Do()
		if err != nil {
			logger.Error("启动任务失败", zap.Any("task", task), zap.Error(err))
			return err
		}
		t.taskIdMap[task.Code] = eid
		t.tasks[task.Code] = task
	}

	return nil
}

func (t Task) stopAndRemove(task TaskConfig) error {
	t.Lock()
	defer t.Unlock()
	if _, ok := t.taskIdMap[task.Code]; ok {
		delete(t.taskIdMap, task.Code)
	}

	return nil
}
