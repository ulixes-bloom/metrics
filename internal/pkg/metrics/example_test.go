package metrics_test

import (
	"fmt"

	"github.com/ulixes-bloom/ya-metrics/internal/pkg/metrics"
)

func Example() {
	gm := metrics.NewGaugeMetric("gaugeID", 1)
	fmt.Println(gm.ID, gm.MType, gm.GetValue())

	cm := metrics.NewCounterMetric("counterID", 1)
	fmt.Println(cm.ID, cm.MType, cm.GetDelta())

	// Output:
	// gaugeID gauge 1
	// counterID counter 1
}
