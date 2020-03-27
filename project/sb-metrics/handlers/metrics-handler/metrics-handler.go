package metricshandler

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/MartinNikolovMarinov/sb-infra/infra"
	"github.com/MartinNikolovMarinov/sb-infra/pubsub"
	"github.com/go-playground/validator"
)

// All commnet
func All(metrics pubsub.Cache, v *validator.Validate) infra.RootHandler {
	ret := func(w http.ResponseWriter, r *http.Request) error {
		all := metrics.All()
		sort.Slice(all, func(i, j int) bool {
			return all[i].Timestamp.After(all[j].Timestamp)
		})

		n := 20
		if n >= len(all) {
			n = len(all)
		}
		lastN := all[:n]
		result := make([]pubsub.MetricRecord, 0, n+1)

		for i := 0; i < len(lastN); i++ {
			var r pubsub.MetricRecord
			err := json.Unmarshal(all[i].Value, &r)
			if err != nil {
				infra.Error(err.Error())
			} else {
				r.Timestamp = all[i].Timestamp
				result = append(result, r)
			}
		}

		err := infra.WriteResponseJSON(w, result, 200)
		if err != nil {
			return err
		}

		return nil
	}

	return ret
}
