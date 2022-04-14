package aurestbreakerprometheus

import (
	aurestbreaker "github.com/StephanHCB/go-autumn-restclient-circuitbreaker/implementation/breaker"
	aurestclientapi "github.com/StephanHCB/go-autumn-restclient/api"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sony/gobreaker"
)

var (
	cbStates            *prometheus.GaugeVec
	cbStateChangeCounts *prometheus.CounterVec
	cbCountsGauge       *prometheus.GaugeVec
)

func SetupCircuitBreakerMetrics() {
	SetupCommon()

	if cbStates == nil {
		cbStates = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "http_client_requests_circuitbreaker_states",
				Help: "Current state of the circuit breaker (-1 = open, -0.5 = half-open, 1 = closed) by circuit breaker.",
			},
			[]string{"circuitBreakerName"},
		)
		prometheus.MustRegister(cbStates)
	}

	if cbStateChangeCounts == nil {
		cbStateChangeCounts = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_client_requests_circuitbreaker_statechanges_count",
				Help: "Number of state changes by circuit breaker and target state.",
			},
			[]string{"circuitBreakerName", "targetState"},
		)
		prometheus.MustRegister(cbStateChangeCounts)
	}

	if cbCountsGauge == nil {
		cbCountsGauge = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "http_client_requests_circuitbreaker_details_count",
				Help: "Detailed counts by circuit breaker and count type.",
			},
			[]string{"circuitBreakerName", "counterType"},
		)
		prometheus.MustRegister(cbCountsGauge)
	}
}

//goland:noinspection GoUnusedExportedFunction
func InstrumentCircuitBreakerClient(client aurestclientapi.Client) {
	SetupCircuitBreakerMetrics()
	aurestbreaker.Instrument(client, StateChangeCallback, CountsCallback)
}

func StateChangeCallback(circuitBreakerName string, targetState string) {
	switch targetState {
	case "closed":
		cbStates.WithLabelValues(circuitBreakerName).Set(1.0)
	case "half-open":
		cbStates.WithLabelValues(circuitBreakerName).Set(-0.5)
	case "open":
		cbStates.WithLabelValues(circuitBreakerName).Set(-1.0)
	default:
		cbStates.WithLabelValues(circuitBreakerName).Set(0.0)
	}

	cbStateChangeCounts.WithLabelValues(circuitBreakerName, targetState).Inc()
}

func CountsCallback(circuitBreakerName string, counts gobreaker.Counts) {
	cbCountsGauge.WithLabelValues(circuitBreakerName, "Requests").Set(float64(counts.Requests))
	cbCountsGauge.WithLabelValues(circuitBreakerName, "TotalSuccesses").Set(float64(counts.TotalSuccesses))
	cbCountsGauge.WithLabelValues(circuitBreakerName, "TotalFailures").Set(float64(counts.TotalFailures))
	cbCountsGauge.WithLabelValues(circuitBreakerName, "ConsecutiveSuccesses").Set(float64(counts.ConsecutiveSuccesses))
	cbCountsGauge.WithLabelValues(circuitBreakerName, "ConsecutiveFailures").Set(float64(counts.ConsecutiveFailures))
}
