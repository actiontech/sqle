package model

import "sqle/errors"

type Config struct {
	Name    string `json:"name"`
	Value   string `json:"value"`
	Default string `json:"default"`
	Desc    string `json:"desc"`
}

func (s *Storage) GetAllConfig() ([]Config, error) {
	configs := []Config{}
	err := s.db.Find(&configs).Error
	return configs, errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

func (s *Storage) GetConfigMap() (map[string]Config, error) {
	configs, err := s.GetAllConfig()
	if err != nil {
		return nil, err
	}
	configMap := make(map[string]Config, len(configs))
	for _, config := range configs {
		configMap[config.Name] = config
	}
	return configMap, nil
}

func (s *Storage) UpdateConfigValueByName(name, value string) error {
	err := s.db.Table("configs").Where("name = ?", name).
		Update(map[string]string{"value": value}).Error
	return errors.New(errors.CONNECT_STORAGE_ERROR, err)
}
