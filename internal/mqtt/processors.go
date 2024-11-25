package mqtt

import "strconv"

func processTemperature1(data []byte) (map[string]interface{}, map[string]interface{}, error) {
	rawTemperature, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		return nil, nil, err
	}
	rawData := map[string]interface{}{
		"temperature_1": rawTemperature,
	}

	processedData := map[string]interface{}{
		"proc_temperature_1": rawTemperature + 1,
	}
	return rawData, processedData, nil
}

func processTemperature2(data []byte) (map[string]interface{}, map[string]interface{}, error) {
	rawTemperature, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		return nil, nil, err
	}
	rawData := map[string]interface{}{
		"temperature_2": rawTemperature,
	}

	processedData := map[string]interface{}{
		"proc_temperature_2": rawTemperature + 1,
	}
	return rawData, processedData, nil
}

func processMotorCurrent(data []byte) (map[string]interface{}, map[string]interface{}, error) {
	rawMotorCurrent, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		return nil, nil, err
	}
	rawData := map[string]interface{}{
		"motor_current": rawMotorCurrent,
	}

	processedData := map[string]interface{}{
		"proc_motor_current": rawMotorCurrent + 5.5,
	}
	return rawData, processedData, nil
}

func processGasSensor(data []byte) (map[string]interface{}, map[string]interface{}, error) {
	rawGasSensor, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		return nil, nil, err
	}
	rawData := map[string]interface{}{
		"gas_sensor": rawGasSensor,
	}

	processedData := map[string]interface{}{
		"proc_gas_sensor": rawGasSensor + 5.5,
	}
	return rawData, processedData, nil
}

func processMPU6050Temperature(data []byte) (map[string]interface{}, map[string]interface{}, error) {
	rawMPU6050Temperature, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		return nil, nil, err
	}
	rawData := map[string]interface{}{
		"mpu6050_temperature": rawMPU6050Temperature,
	}

	processedData := map[string]interface{}{
		"proc_mpu6050_temperature": rawMPU6050Temperature + 1.5,
	}
	return rawData, processedData, nil
}

func processMPU6050GyroZ(data []byte) (map[string]interface{}, map[string]interface{}, error) {
	rawMPU6050GyroZ, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		return nil, nil, err
	}
	rawData := map[string]interface{}{
		"mpu6050_gyro_z": rawMPU6050GyroZ,
	}

	processedData := map[string]interface{}{
		"proc_mpu6050_gyro_z": rawMPU6050GyroZ + 3.5,
	}
	return rawData, processedData, nil
}

func processMPU6050AccelZ(data []byte) (map[string]interface{}, map[string]interface{}, error) {
	rawMPU6050AccelZ, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		return nil, nil, err
	}
	rawData := map[string]interface{}{
		"mpu6050_accel_z": rawMPU6050AccelZ,
	}

	processedData := map[string]interface{}{
		"proc_mpu6050_accel_z": rawMPU6050AccelZ + 7.5,
	}
	return rawData, processedData, nil
}

func processMPU6050GyroY(data []byte) (map[string]interface{}, map[string]interface{}, error) {
	rawMPU6050GyroY, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		return nil, nil, err
	}
	rawData := map[string]interface{}{
		"mpu6050_gyro_y": rawMPU6050GyroY,
	}

	processedData := map[string]interface{}{
		"proc_mpu6050_gyro_y": rawMPU6050GyroY + 1.5,
	}
	return rawData, processedData, nil
}

func processMPU6050AccelY(data []byte) (map[string]interface{}, map[string]interface{}, error) {
	rawMPU6050AccelY, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		return nil, nil, err
	}
	rawData := map[string]interface{}{
		"mpu6050_accel_y": rawMPU6050AccelY,
	}

	processedData := map[string]interface{}{
		"proc_mpu6050_accel_y": rawMPU6050AccelY + 2.5,
	}
	return rawData, processedData, nil
}

func processMPU6050GyroX(data []byte) (map[string]interface{}, map[string]interface{}, error) {
	rawMPU6050GyroX, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		return nil, nil, err
	}
	rawData := map[string]interface{}{
		"mpu6050_gyro_x": rawMPU6050GyroX,
	}

	processedData := map[string]interface{}{
		"proc_mpu6050_gyro_x": rawMPU6050GyroX + 3.5,
	}
	return rawData, processedData, nil
}

func processMPU6050AccelX(data []byte) (map[string]interface{}, map[string]interface{}, error) {
	rawMPU6050AccelX, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		return nil, nil, err
	}
	rawData := map[string]interface{}{
		"mpu6050_accel_x": rawMPU6050AccelX,
	}

	processedData := map[string]interface{}{
		"proc_mpu6050_accel_x": rawMPU6050AccelX + 6.5,
	}
	return rawData, processedData, nil
}

func processINA226Power(data []byte) (map[string]interface{}, map[string]interface{}, error) {
	rawINA226Power, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		return nil, nil, err
	}
	rawData := map[string]interface{}{
		"ina226_power": rawINA226Power,
	}

	processedData := map[string]interface{}{
		"proc_ina226_power": rawINA226Power + 3.5,
	}
	return rawData, processedData, nil
}

func processINA226Current(data []byte) (map[string]interface{}, map[string]interface{}, error) {
	rawINA226Current, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		return nil, nil, err
	}
	rawData := map[string]interface{}{
		"ina226_current": rawINA226Current,
	}

	processedData := map[string]interface{}{
		"proc_ina226_current": rawINA226Current + 1.5,
	}
	return rawData, processedData, nil
}

func processINA226ShuntVoltage(data []byte) (map[string]interface{}, map[string]interface{}, error) {
	rawINA226ShuntVoltage, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		return nil, nil, err
	}
	rawData := map[string]interface{}{
		"ina226_shunt_voltage": rawINA226ShuntVoltage,
	}

	processedData := map[string]interface{}{
		"proc_ina226_shunt_voltage": rawINA226ShuntVoltage + 2.5,
	}
	return rawData, processedData, nil
}
