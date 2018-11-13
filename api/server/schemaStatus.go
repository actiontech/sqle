package server

import (
	"time"
	"sqle/model"
	"sync"
	"sqle/executor"
)

type InstanceStatus struct {
	ID              uint     `json:"id"`
	Name            string   `json:"name"`
	Host            string   `json:"host"`
	Port            string   `json:"port"`
	IsConnectFailed bool     `json:"is_connect_failed"`
	Schemas         []string `json:"schema_list"`
}

func (s *Sqled) statusLoop() {
	tick := time.Tick(1 * time.Hour)
	s.UpdateAllInstanceStatus()
	for {
		select {
		case <-s.exit:
			return
		case <-tick:
			s.UpdateAllInstanceStatus()
		}
	}
}

func (s *Sqled) UpdateAllInstanceStatus() error {
	st := model.GetStorage()
	instances, err := st.GetInstances()
	if err != nil {
		return err
	}
	instancesStatus := map[uint]*InstanceStatus{}
	wait := sync.WaitGroup{}
	mutex := sync.Mutex{}
	for _, instance := range instances {
		wait.Add(1)
		currentInstance := instance
		go func() {
			status := &InstanceStatus{
				ID:   currentInstance.ID,
				Name: currentInstance.Name,
				Host: currentInstance.Host,
				Port: currentInstance.Port,
			}
			schemas, err := executor.ShowDatabases(currentInstance)
			if err != nil {
				status.IsConnectFailed = true
			} else {
				status.Schemas = schemas
			}
			mutex.Lock()
			instancesStatus[currentInstance.ID] = status
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

func (s *Sqled) UpdateAndGetInstanceStatus(instance model.Instance) (*InstanceStatus, error) {
	status := &InstanceStatus{
		ID:   instance.ID,
		Name: instance.Name,
		Host: instance.Host,
		Port: instance.Port,
	}
	schemas, err := executor.ShowDatabases(instance)
	if err != nil {
		status.IsConnectFailed = true
	} else {
		status.Schemas = schemas
	}
	s.Lock()
	s.instancesStatus[instance.ID] = status
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

func (s *Sqled) DeleteInstanceStatus(instance model.Instance) {
	s.Lock()
	delete(s.instancesStatus, instance.ID)
	s.Unlock()
}
