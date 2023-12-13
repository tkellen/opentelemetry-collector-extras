package cache

import (
	"github.com/honeycombio/metricsasattributesprocessor/internal/common"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
	"sync"
	"time"
)

type Metric struct {
	DataPoint *pmetric.NumberDataPoint
	updated   time.Time
}

type MatchedMetricsCache struct {
	// map key is the metric name
	Metrics map[string]*Metric
	Mutex   sync.RWMutex
}

type MetricGroup struct {
	// map key is the matched id
	matchedMetrics map[string]*MatchedMetricsCache
	Mutex          sync.RWMutex
}

func (mg *MetricGroup) SetMetric(id string, name string, data *pmetric.NumberDataPoint) {
	if _, ok := mg.matchedMetrics[id]; !ok {
		mg.Mutex.Lock()
		// we do an additional check here in case we were waiting for the lock
		if _, ok := mg.matchedMetrics[id]; !ok {
			mg.matchedMetrics[id] = &MatchedMetricsCache{
				Metrics: make(map[string]*Metric),
			}
		}
		mg.Mutex.Unlock()
	}

	mmc := mg.matchedMetrics[id]
	mmc.Mutex.Lock()
	mmc.Metrics[name] = &Metric{
		DataPoint: data,
		updated:   time.Now(),
	}
	mmc.Mutex.Unlock()
}

func (mg *MetricGroup) HasMatchedMetrics(id string) bool {
	_, ok := mg.matchedMetrics[id]
	return ok
}

func (mg *MetricGroup) GetMatchedMetricsCache(id string) *MatchedMetricsCache {
	if mmc, ok := mg.matchedMetrics[id]; ok {
		return mmc
	}
	return nil
}

const minimumTtl = time.Minute

type MetricCache struct {
	// map key is the metric group name
	MetricGroups map[string]*MetricGroup
	ttl          time.Duration
	logger       *zap.Logger
}

func (mc *MetricCache) cleanup() {
	cleaned := 0
	for _, mg := range mc.MetricGroups {
		for k, v := range mg.matchedMetrics {
			v.Mutex.Lock()
			for k2, v2 := range v.Metrics {
				if time.Since(v2.updated) > mc.ttl {
					delete(v.Metrics, k2)
					cleaned++
				}
			}
			v.Mutex.Unlock()
			if len(v.Metrics) == 0 {
				mg.Mutex.Lock()
				delete(mg.matchedMetrics, k)
				mg.Mutex.Unlock()
			}
		}
	}
	mc.logger.Debug("Cleaned up expired metrics", zap.Int("cleaned", cleaned))
}

func (mc *MetricCache) startCleanupTimer() {
	ticker := time.Tick(mc.ttl / 2)
	go (func() {
		for {
			select {
			case <-ticker:
				mc.cleanup()
			}
		}
	})()
}

var (
	caches = make(map[string]MetricCache)
	lock   sync.Mutex
)

func GetCache(id string, cacheTtl time.Duration, metricGroups []common.MetricGroup, logger *zap.Logger) MetricCache {
	lock.Lock()
	defer lock.Unlock()

	if c, ok := caches[id]; ok {
		return c
	}

	caches[id] = newCache(cacheTtl, metricGroups, logger)
	return caches[id]
}

func newCache(cacheTtl time.Duration, metricGroups []common.MetricGroup, logger *zap.Logger) MetricCache {
	var ttl time.Duration
	if cacheTtl < minimumTtl {
		ttl = minimumTtl
	} else {
		ttl = cacheTtl
	}

	mgs := make(map[string]*MetricGroup)
	for _, v := range metricGroups {
		mgs[v.Name] = &MetricGroup{
			matchedMetrics: make(map[string]*MatchedMetricsCache),
		}
	}

	mc := MetricCache{
		MetricGroups: mgs,
		ttl:          ttl,
		logger:       logger,
	}
	mc.startCleanupTimer()
	return mc
}
