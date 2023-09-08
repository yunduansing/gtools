package cronjob

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/yunduansing/gtools/logger"
	"gorm.io/gorm"
)

var (
	taskContainer Task
)

func Init(db *gorm.DB, serviceName string, taskList []TaskConfig) {
	taskContainer.Do(func() {
		checkIsNeedSave(db, taskList)
		taskContainer.serviceName = serviceName
		taskContainer.tasks = make(map[string]TaskConfig)
		taskContainer.taskIdMap = make(map[string]cron.EntryID)
		for i, v := range taskList {
			taskContainer.tasks[v.Code] = taskList[i]
		}
	})

}
func InitAndRun(db *gorm.DB, serviceName string, taskList []TaskConfig) {
	taskContainer.Do(func() {
		checkIsNeedSave(db, taskList)
		taskContainer.serviceName = serviceName
		taskContainer.tasks = make(map[string]TaskConfig)
		taskContainer.taskIdMap = make(map[string]cron.EntryID)
		for i, v := range taskList {
			taskContainer.tasks[v.Code] = taskList[i]
		}
	})
	taskContainer.run()
}

func FindTaskConfigByCode(db *gorm.DB, code string) TaskConfig {

	var task TaskConfig
	err := db.First(&task, "code=?", code).Error
	if err != nil {
		logger.Error("查询task config err:", err)
	}

	return TaskConfig{}
}

func UpdateTask(db *gorm.DB, t TaskConfig) error {
	err := db.Updates(&t).Error
	if err != nil {
		return err
	}
	taskContainer.stopAndRemove(t)
	if t.State == 1 {
		taskContainer.addAndRun(t)
	}

	return nil
}

// Topic 当使用kafka来通知时
func topic() string {
	return fmt.Sprintf("%s:task:event", taskContainer.serviceName)
}

func checkIsNeedSave(db *gorm.DB, tasks []TaskConfig) {
	var exists, needSave []TaskConfig
	var existsMap = make(map[string]TaskConfig)
	db.Find(&exists)
	for i, v := range exists {
		existsMap[v.Code] = exists[i]
	}
	for i, v := range tasks {
		if _, ok := existsMap[v.Code]; !ok {
			needSave = append(needSave, tasks[i])
		}
	}
	if len(needSave) > 0 {
		db.Save(needSave)
	}
}

func AddTask(t TaskConfig) error {
	return taskContainer.addAndRun(t)
}

func RemoveTask(t TaskConfig) {
	taskContainer.stopAndRemove(t)
}
