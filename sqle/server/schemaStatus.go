package server

import (
	"actiontech.cloud/sqle/sqle/sqle/executor"
	"actiontech.cloud/sqle/sqle/sqle/log"
	"actiontech.cloud/sqle/sqle/sqle/model"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type InstanceStatus struct {
	DbType          string   `json:"db_type"`
	Name            string   `json:"name"`
	Host            string   `json:"host"`
	Port            string   `json:"port"`
	IsConnectFailed bool     `json:"is_connect_failed"`
	Schemas         []string `json:"schema_list"`
}

func (s *Sqled) statusLoop() {
	tick := time.Tick(1 * time.Hour)
	entry := log.NewEntry().WithField("type", "cron")
	s.UpdateAllInstanceStatus(entry)
	for {
		select {
		case <-s.exit:
			return
		case <-tick:
			s.UpdateAllInstanceStatus(entry)
		}
	}
}

func (s *Sqled) UpdateAllInstanceStatus(entry *logrus.Entry) error {
	st := model.GetStorage()
	instances, err := st.GetInstances()
	if err != nil {
		return err
	}
	instancesStatus := map[string]*InstanceStatus{}
	wait := sync.WaitGroup{}
	mutex := sync.Mutex{}
	for _, instance := range instances {
		wait.Add(1)
		currentInstance := instance
		go func() {
			status := &InstanceStatus{
				DbType:  currentInstance.DbType,
				Name:    currentInstance.Name,
				Host:    currentInstance.Host,
				Port:    currentInstance.Port,
				Schemas: []string{},
			}
			schemas, err := executor.ShowDatabases(entry, &currentInstance)
			if err != nil {
				status.IsConnectFailed = true
			} else {
				status.Schemas = schemas
			}
			mutex.Lock()
			instancesStatus[currentInstance.Name] = status
			mutex.Unlock()
			wait.Done()
		}()
	}
	wait.Wait()
	s.Lock()
	s.instancesStatus = instancesStatus
	s.Unlock()
	return nil
}

func (s *Sqled) UpdateAndGetInstanceStatus(entry *logrus.Entry, instance *model.Instance) (*InstanceStatus, error) {
	status := &InstanceStatus{
		DbType:  instance.DbType,
		Name:    instance.Name,
		Host:    instance.Host,
		Port:    instance.Port,
		Schemas: []string{},
	}
	schemas, err := executor.ShowDatabases(entry, instance)
	if err != nil {
		status.IsConnectFailed = true
	} else {
		status.Schemas = schemas
	}
	s.Lock()
	s.instancesStatus[instance.Name] = status
	s.Unlock()
	return status, err
}

func (s *Sqled) GetAllInstanceStatus() []InstanceStatus {
	statusList := make([]InstanceStatus, 0, len(s.instancesStatus))
	s.Lock()
	for _, status := range s.instancesStatus {
		statusList = append(statusList, *status)
	}
	s.Unlock()
	return statusList
}

func (s *Sqled) DeleteInstanceStatus(instance *model.Instance) {
	s.Lock()
	delete(s.instancesStatus, instance.Name)
	s.Unlock()
}
