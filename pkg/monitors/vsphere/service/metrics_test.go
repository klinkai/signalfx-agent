package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPopulateInvMetrics(t *testing.T) {
	gateway := newFakeGateway()
	metricsSvc := NewMetricsService(gateway)
	inventorySvc := NewInventorySvc(gateway)
	inv, _ := inventorySvc.RetrieveInventory()
	metricsSvc.PopulateInvMetrics(inv)
	invOjb := inv.Objects[0]
	perfMetricId := invOjb.MetricIds[0]
	require.EqualValues(t, "instance-0", perfMetricId.Instance)
}

func TestRetrievePerfCounterIndex(t *testing.T) {
	gateway := newFakeGateway()
	metricsSvc := NewMetricsService(gateway)
	idx, _ := metricsSvc.RetrievePerfCounterIndex()
	metric := idx[42]
	require.Equal(t, "vsphere.cpu_core_utilization_percent", metric.MetricName)
}
