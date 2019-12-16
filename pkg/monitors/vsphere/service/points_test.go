package service

import (
	"testing"

	"github.com/signalfx/golib/v3/datapoint"
	"github.com/stretchr/testify/require"
)

func TestPointsRetrievePoints(t *testing.T) {
	gateway := newFakeGateway()
	inventorySvc := NewInventorySvc(gateway)
	metricsSvc := NewMetricsService(gateway)
	infoSvc := NewVSphereInfoService(inventorySvc, metricsSvc)
	vsphereInfo, _ := infoSvc.RetrieveVSphereInfo()
	svc := NewPointsSvc(gateway)
	pts, _ := svc.RetrievePoints(vsphereInfo, 1)
	pt := pts[0]
	require.Equal(t, "vsphere.cpu_core_utilization_percent", pt.Metric)
	require.Equal(t, datapoint.Counter, pt.MetricType)
	require.EqualValues(t, 111, pt.Value)
}
