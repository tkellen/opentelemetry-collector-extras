package metricsasattributesprocessor

import (
	"time"

	"github.com/tkellen/opentelemetry-collector-extras/metricsasattributesprocessor/internal/common"
)

type Config struct {
	// time before state metrics are removed from cache (default 5m)
	CacheTtl time.Duration `mapstructure:"cache_ttl"`

	// lists of different groups of metrics to add to the target signal resources
	MetricGroups []common.MetricGroup `mapstructure:"metrics_groups,required"`
}
