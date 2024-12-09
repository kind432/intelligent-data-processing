package mqtt

import (
	"strconv"
	"sync"
)

type SensorData struct {
	Values []float64
	mu     sync.Mutex
}

var sensorDataStore = map[string]map[string]*SensorData{}
var maxDataPoints = 10

var sensorDataStoreMu sync.Mutex

func getSensorData(serialNumber, dataKey string) *SensorData {
	sensorDataStoreMu.Lock()
	defer sensorDataStoreMu.Unlock()

	if data, exists := sensorDataStore[serialNumber]; exists {
		if sensor, exists := data[dataKey]; exists {
			return sensor
		}
	}

	if sensorDataStore[serialNumber] == nil {
		sensorDataStore[serialNumber] = map[string]*SensorData{}
	}
	data := &SensorData{}
	sensorDataStore[serialNumber][dataKey] = data
	return data
}

func processSensorData(serialNumber, dataKey string, data []byte) (map[string]interface{}, map[string]interface{}, error) {
	rawValue, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		return nil, nil, err
	}

	sensorData := getSensorData(serialNumber, dataKey)

	sensorData.mu.Lock()
	defer sensorData.mu.Unlock()

	if len(sensorData.Values) >= maxDataPoints {
		sensorData.Values = sensorData.Values[1:]
	}
	sensorData.Values = append(sensorData.Values, rawValue)

	var sum float64
	for _, value := range sensorData.Values {
		sum += value
	}
	average := sum / float64(len(sensorData.Values))

	rawData := map[string]interface{}{
		"sensorType": "default",
		dataKey:      rawValue,
	}

	processedData := map[string]interface{}{
		"sensorType":      "default",
		"proc_" + dataKey: average,
	}

	return rawData, processedData, nil
}
