package service

import (
	"context"

	"github.com/signalfx/signalfx-agent/pkg/monitors/vsphere/model"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

// A thin wrapper around the vmomi SDK so that callers don't have to use it directly.
type IGateway interface {
	retrievePerformanceManager() (*mo.PerformanceManager, error)
	retrieveTopLevelFolder() (*mo.Folder, error)
	retrieveRefProperties(mor types.ManagedObjectReference, dst interface{}) error
	queryAvailablePerfMetric(ref types.ManagedObjectReference) (*types.QueryAvailablePerfMetricResponse, error)
	queryPerfProviderSummary(mor types.ManagedObjectReference) (*types.QueryPerfProviderSummaryResponse, error)
	queryPerf(invObjs []*model.InventoryObject, maxSample int32) (*types.QueryPerfResponse, error)
}

type gateway struct {
	ctx    context.Context
	client *govmomi.Client
}

func NewGateway(ctx context.Context, client *govmomi.Client) *gateway {
	return &gateway{ctx, client}
}

func (g *gateway) retrievePerformanceManager() (*mo.PerformanceManager, error) {
	var pm mo.PerformanceManager
	err := mo.RetrieveProperties(
		g.ctx,
		g.client,
		g.client.ServiceContent.PropertyCollector,
		*g.client.Client.ServiceContent.PerfManager,
		&pm,
	)
	return &pm, err
}

func (g *gateway) retrieveTopLevelFolder() (*mo.Folder, error) {
	var folder mo.Folder
	err := mo.RetrieveProperties(
		g.ctx,
		g.client,
		g.client.ServiceContent.PropertyCollector,
		g.client.ServiceContent.RootFolder,
		&folder,
	)
	return &folder, err
}

func (g *gateway) retrieveRefProperties(mor types.ManagedObjectReference, dst interface{}) error {
	return mo.RetrieveProperties(
		g.ctx,
		g.client,
		g.client.ServiceContent.PropertyCollector,
		mor,
		dst,
	)
}

func (g *gateway) queryAvailablePerfMetric(ref types.ManagedObjectReference) (*types.QueryAvailablePerfMetricResponse, error) {
	req := types.QueryAvailablePerfMetric{
		This:       *g.client.Client.ServiceContent.PerfManager,
		Entity:     ref,
		IntervalId: model.VSMetricsInterval,
	}
	return methods.QueryAvailablePerfMetric(g.ctx, g.client.Client, &req)
}

func (g *gateway) queryPerfProviderSummary(mor types.ManagedObjectReference) (*types.QueryPerfProviderSummaryResponse, error) {
	req := types.QueryPerfProviderSummary{
		This:   *g.client.Client.ServiceContent.PerfManager,
		Entity: mor,
	}
	return methods.QueryPerfProviderSummary(g.ctx, g.client.Client, &req)
}

func (g *gateway) queryPerf(invObjs []*model.InventoryObject, maxSample int32) (*types.QueryPerfResponse, error) {
	specs := make([]types.PerfQuerySpec, 0, len(invObjs))
	for _, invObj := range invObjs {
		specs = append(specs, types.PerfQuerySpec{
			Entity:     invObj.Ref,
			MaxSample:  maxSample,
			IntervalId: model.VSMetricsInterval,
			MetricId:   invObj.MetricIds,
		})
	}
	queryPerf := types.QueryPerf{
		This:      *g.client.Client.ServiceContent.PerfManager,
		QuerySpec: specs,
	}
	return methods.QueryPerf(g.ctx, g.client.Client, &queryPerf)
}
