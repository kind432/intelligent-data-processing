package repository

import "sync"

type SensorData struct {
	Values []float64
	Mu     sync.Mutex
}

type Storage struct {
	sensors map[string]map[string]*SensorData // serial -> key -> data
	devices map[string]bool                   // serial -> connected
	mu      sync.RWMutex
}

func NewStorage() *Storage {
	return &Storage{
		sensors: make(map[string]map[string]*SensorData),
		devices: make(map[string]bool),
	}
}

func (s *Storage) UpdateDeviceStatus(serial string, connected bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.devices[serial] = connected
}

func (s *Storage) IsConnected(serial string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.devices[serial]
}

func (s *Storage) GetSensorData(serial, key string) *SensorData {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.sensors[serial]; !ok {
		s.sensors[serial] = make(map[string]*SensorData)
	}
	if _, ok := s.sensors[serial][key]; !ok {
		s.sensors[serial][key] = &SensorData{Values: make([]float64, 0)}
	}
	return s.sensors[serial][key]
}
