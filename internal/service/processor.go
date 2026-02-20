package service

import "intelligent-data-processing/internal/repository"

type Processor struct {
	repo          *repository.Storage
	maxDataPoints int
}

func NewProcessor(repo *repository.Storage) *Processor {
	return &Processor{repo: repo, maxDataPoints: 10}
}

func (p *Processor) ProcessFloat(serial, key string, val float64) float64 {
	data := p.repo.GetSensorData(serial, key)
	data.Mu.Lock()
	defer data.Mu.Unlock()

	if len(data.Values) >= p.maxDataPoints {
		data.Values = data.Values[1:]
	}
	data.Values = append(data.Values, val)

	var sum float64
	for _, v := range data.Values {
		sum += v
	}
	return sum / float64(len(data.Values))
}
