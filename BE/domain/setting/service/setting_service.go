package service

import "errors"

func (s *settingService) GetAll() (map[string]string, error) {
	settings, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	result := make(map[string]string)
	for _, item := range settings {
		result[item.Key] = item.Value
	}
	return result, nil
}

func (s *settingService) GetByKey(key string) (string, error) {
	setting, err := s.repo.GetByKey(key)
	if err != nil {
		return "", err
	}
	if setting == nil {
		return "", errors.New("setting not found")
	}
	return setting.Value, nil
}

func (s *settingService) Save(data map[string]string) error {
	for key, value := range data {
		if err := s.repo.Upsert(key, value); err != nil {
			return err
		}
	}
	return nil
}

func (s *settingService) Reset() error {
	return s.repo.ResetAll()
}
