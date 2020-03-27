package pubsub

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/MartinNikolovMarinov/sb-infra/infra"
	"github.com/hashicorp/go-uuid"
)

// MetricRecord comment
type MetricRecord struct {
	ID           string        `json:"id"`
	URL          string        `json:"url"`
	HTTPMethod   string        `json:"http-method"`
	RespDuration time.Duration `json:"resp-duration"`
	RespFormated string        `json:"resp-formatted"`
	Timestamp    time.Time     `json:"timestamp"`
}

func writeMetrics(metrics *MetricsProducer, t time.Time, r *http.Request) {
	id, err := uuid.GenerateUUID()
	if err != nil {
		infra.Error(err.Error())
		return
	}

	record := &MetricRecord{
		ID:           id,
		URL:          r.RequestURI,
		RespDuration: time.Now().Sub(t),
		RespFormated: time.Now().Sub(t).String(),
		HTTPMethod:   r.Method,
	}

	recordJSON, err := json.Marshal(record)
	if err != nil {
		infra.Error(err.Error())
		return
	}

	err = metrics.Publish([]byte(id), recordJSON)
	if err != nil {
		infra.Error(err.Error())
		return
	}

	infra.Info(fmt.Sprintf("Metrics saved for %s", record.URL))
}

// MetricsMiddleware comment
func MetricsMiddleware(next http.Handler) http.Handler {
	metrics, err := NewMetrics()
	if err != nil {
		infra.Error(err.Error())
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if metrics != nil {
			// If metrics loaded ok use it :
			t := time.Now()
			defer func(t time.Time) {
				go writeMetrics(metrics, t, r)
			}(t)
		}

		next.ServeHTTP(w, r)
	})
}
