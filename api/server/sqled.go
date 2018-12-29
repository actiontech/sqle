package server

import (
	"fmt"
	"math"
	"sqle/errors"
	"sqle/inspector"
	"sqle/log"
	"sqle/model"
	"sync"
)

var sqled *Sqled

func GetSqled() *Sqled {
	return sqled
}

type Sqled struct {
	sync.Mutex
	exit            chan struct{}
	currentTask     map[string]struct{}
	queue           chan *Action
	instancesStatus map[uint]*InstanceStatus
}

type Action struct {
	sync.Mutex
	Task  *model.Task
	Typ   int
	Error error
	Done  chan struct{}
}

func InitSqled(exit chan struct{}) {
	sqled = &Sqled{
		exit:        exit,
		currentTask: map[string]struct{}{},
		queue:       make(chan *Action, 1024),
	}
	sqled.Start()
}

func (s *Sqled) HasTask(taskId string) bool {
	s.Lock()
	_, ok := s.currentTask[taskId]
	s.Unlock()
	return ok
}

func (s *Sqled) addTask(taskId string, typ int) (*Action, error) {
	var err error
	action := &Action{
		Typ:  typ,
		Done: make(chan struct{}),
	}
	s.Lock()
	_, taskRunning := s.currentTask[taskId]
	if !taskRunning {
		s.currentTask[taskId] = struct{}{}
	}
	s.Unlock()
	if taskRunning {
		return action, errors.New(errors.TASK_RUNNING, fmt.Errorf("task is running"))
	}

	task, exist, err := model.GetStorage().GetTaskById(taskId)
	if err != nil {
		goto Error
	}
	if !exist {
		err = errors.New(errors.TASK_NOT_EXIST, fmt.Errorf("task not exist"))
		goto Error
	}

	err = task.ValidAction(typ)
	if err != nil {
		goto Error
	}

	if int(task.Action) < action.Typ {
		err := model.GetStorage().UpdateTask(task, map[string]interface{}{"action": action.Typ})
		if err != nil {
			goto Error
		}
	}
	action.Task = task
	s.queue <- action
	return action, nil

Error:
	s.Lock()
	delete(s.currentTask, taskId)
	s.Unlock()
	return action, err
}

func (s *Sqled) AddTask(taskId string, typ int) error {
	_, err := s.addTask(taskId, typ)
	return err
}

func (s *Sqled) AddTaskWaitResult(taskId string, typ int) (*model.Task, error) {
	action, err := s.addTask(taskId, typ)
	if err != nil {
		return nil, err
	}
	<-action.Done
	return action.Task, action.Error
}

func (s *Sqled) Start() {
	go s.taskLoop()
	go s.statusLoop()
}

func (s *Sqled) taskLoop() {
	for {
		select {
		case <-s.exit:
			return
		case action := <-s.queue:
			go s.do(action)
		}
	}
}

func (s *Sqled) do(action *Action) error {
	var err error
	switch action.Typ {
	case model.TASK_ACTION_INSPECT:
		err = s.inspect(action.Task)
	case model.TASK_ACTION_COMMIT:
		err = s.commit(action.Task)
	case model.TASK_ACTION_ROLLBACK:
		err = s.rollback(action.Task)
	}
	if err != nil {
		action.Error = err
	}
	s.Lock()
	taskId := fmt.Sprintf(fmt.Sprintf("%d", action.Task.ID))
	delete(s.currentTask, taskId)
	s.Unlock()

	select {
	case action.Done <- struct{}{}:
	default:
	}
	return err
}

func (s *Sqled) inspect(task *model.Task) error {
	entry := log.NewEntry().WithField("task_id", task.ID)

	st := model.GetStorage()

	rules, err := st.GetRulesByInstanceId(fmt.Sprintf("%v", task.InstanceId))
	if err != nil {
		entry.Logger.Errorf("get instance rule from storage failed, error: %v", err)
		return err
	}
	ruleMap := model.GetRuleMapFromAllArray(rules)
	ctx := inspector.NewContext(nil)
	i := inspector.NewInspector(entry, ctx, task, nil, ruleMap)
	err = i.Advise(rules)
	if err != nil {
		return err
	}

	// if sql type is DML and sql invalid, try to advise with other DDL.
	if i.SqlType() == model.SQL_TYPE_DML && i.SqlInvalid() {
		relateTasks, err := st.GetRelatedDDLTask(task)
		if err != nil {
			return err
		}
		i = inspector.NewInspector(entry, ctx, task, relateTasks, ruleMap)
		err = i.Advise(rules)
		if err != nil {
			return err
		}
	}

	var normalCount float64
	for _, sql := range task.CommitSqls {
		if sql.InspectLevel == model.RULE_LEVEL_NORMAL {
			normalCount += 1
		}
		if err := st.Save(&sql); err != nil {
			entry.Errorf("save commit sql to storage failed, error: %v", err)
			return err
		}
	}
	if len(task.CommitSqls) != 0 {
		task.NormalRate = round(normalCount/float64(len(task.CommitSqls)), 4)
	}

	err = st.UpdateTask(task, map[string]interface{}{
		"sql_type":    i.SqlType(),
		"normal_rate": task.NormalRate,
	})
	if err != nil {
		entry.Errorf("update task to storage failed, error: %v", err)
		return err
	}
	ctx = inspector.NewContext(nil)
	i = inspector.NewInspector(entry, ctx, task, nil, ruleMap)
	rollbackSqls, err := i.GenerateAllRollbackSql()
	if err != nil {
		return err
	}

	err = st.UpdateRollbackSql(task, rollbackSqls)
	if err != nil {
		entry.Errorf("save rollback sql to storage failed, error: %v", err)
		return err
	}
	return nil
}

func (s *Sqled) commit(task *model.Task) error {
	entry := log.NewEntry().WithField("task_id", task.ID)

	st := model.GetStorage()

	entry.Info("start commit")
	i := inspector.NewInspector(entry, inspector.NewContext(nil), task, nil, nil)
	for _, commitSql := range task.CommitSqls {
		currentSql := commitSql
		err := i.Add(&currentSql.Sql, func(sql *model.Sql) error {
			err := st.UpdateCommitSqlStatus(sql, model.TASK_ACTION_DOING, "")
			if err != nil {
				i.Logger().Errorf("update commit sql status to storage failed, error: %v", err)
				return err
			}
			i.Commit(sql)
			err = st.Save(currentSql)
			if err != nil {
				i.Logger().Errorf("save commit sql to storage failed, error: %v", err)
			}
			return err
		})
		if err != nil {
			entry.Error("add commit sql to task failed")
			return err
		}
	}
	err := i.Do()
	if err != nil {
		entry.Error("commit sql failed")
	} else {
		entry.Info("commit sql finish")
	}
	return err
}

func (s *Sqled) rollback(task *model.Task) error {
	entry := log.NewEntry().WithField("task_id", task.ID)
	entry.Info("start rollback sql")

	st := model.GetStorage()
	i := inspector.NewInspector(entry, inspector.NewContext(nil), task, nil, nil)

	for _, rollbackSql := range task.RollbackSqls {
		currentSql := rollbackSql
		if currentSql.Content == "" {
			continue
		}
		err := i.Add(&currentSql.Sql, func(sql *model.Sql) error {
			err := st.UpdateRollbackSqlStatus(sql, model.TASK_ACTION_DOING, "")
			if err != nil {
				i.Logger().Errorf("update rollback sql status to storage failed, error: %v", err)
				return err
			}
			i.Commit(sql)
			err = st.Save(currentSql)
			if err != nil {
				i.Logger().Errorf("save commit sql to storage failed, error: %v", err)
			}
			return err
		})
		if err != nil {
			entry.Error("add rollback sql to task failed")
			return err
		}
	}
	err := i.Do()
	if err != nil {
		entry.Error("rollback sql failed")
	} else {
		entry.Info("rollback sql finish")
	}
	return err
}

func round(f float64, n int) float64 {
	p := math.Pow10(n)
	return math.Trunc(f*p+0.5) / p
}
