package sensor

import (
	"sync"
)

type topicsData struct {
	values []float64
	mu     sync.Mutex
}

var (
	sensorDataStore   = map[string]map[string]*topicsData{}
	sensorDataStoreMu sync.Mutex
	maxDataPoints     = 10
)

func getSensorData(serialNumber, dataKey string) *topicsData {
	sensorDataStoreMu.Lock()
	defer sensorDataStoreMu.Unlock()

	if data, exists := sensorDataStore[serialNumber]; exists {
		if sensor, exists := data[dataKey]; exists {
			return sensor
		}
	}

	if sensorDataStore[serialNumber] == nil {
		sensorDataStore[serialNumber] = map[string]*topicsData{}
	}
	data := &topicsData{}
	sensorDataStore[serialNumber][dataKey] = data
	return data
}

func ProcessSensorData(serialNumber, dataKey string, rawValue float64) (map[string]interface{}, map[string]interface{}, error) {
	sensorData := getSensorData(serialNumber, dataKey)

	sensorData.mu.Lock()
	defer sensorData.mu.Unlock()

	if len(sensorData.values) >= maxDataPoints {
		sensorData.values = sensorData.values[1:]
	}
	sensorData.values = append(sensorData.values, rawValue)

	var sum float64
	for _, value := range sensorData.values {
		sum += value
	}
	average := sum / float64(len(sensorData.values))

	rawData := map[string]interface{}{
		"sensorType": "default",
		dataKey:      rawValue,
	}
	procData := map[string]interface{}{
		"sensorType":      "default",
		"proc_" + dataKey: average,
	}
	return rawData, procData, nil
}
