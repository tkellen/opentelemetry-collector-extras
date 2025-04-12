package metricsasattributesprocessor

import (
	"github.com/tkellen/opentelemetry-collector-extras/metricsasattributesprocessor/internal/common"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

func isSelectable(selectors []common.Selector, r pcommon.Resource, is pcommon.InstrumentationScope, attrs pcommon.Map) (string, bool) {
	id := ""
	for _, s := range selectors {
		var val pcommon.Value
		var ok bool
		switch s.AttributeType {
		case common.AttributeTypeResource:
			val, ok = r.Attributes().Get(s.Name)
		case common.AttributeTypeScope:
			val, ok = is.Attributes().Get(s.Name)
		default:
			val, ok = attrs.Get(s.Name)
		}
		if !ok {
			return "", false
		}
		id += val.AsString() + MatcherDelim
	}
	return id, true
}
