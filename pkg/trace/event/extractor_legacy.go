package event

import (
	"github.com/DataDog/datadog-agent/pkg/trace/agent"
	"github.com/DataDog/datadog-agent/pkg/trace/sampler"
)

// legacyExtractor is an event extractor that decides whether to extract APM events from spans based on
// `serviceName => sampling rate` mappings.
type legacyExtractor struct {
	rateByService map[string]float64
}

// NewLegacyExtractor returns an APM event extractor that decides whether to extract APM events from spans following the
// specified extraction rates for a span's service.
func NewLegacyExtractor(rateByService map[string]float64) Extractor {
	return &legacyExtractor{
		rateByService: rateByService,
	}
}

// Extract decides to extract an apm event from the provided span if there's an extraction rate configured for that
// span's service. In this case the extracted event is returned along with the found extraction rate and a true value.
// If this rate doesn't exist or the provided span is not a top level one, then no extraction is done and false is
// returned as the third value, with the others being invalid.
func (e *legacyExtractor) Extract(s *agent.WeightedSpan, priority sampler.SamplingPriority) (float64, bool) {
	if !s.TopLevel {
		return 0, false
	}
	extractionRate, ok := e.rateByService[s.Service]
	if !ok {
		return 0, false
	}
	return extractionRate, true
}
