package service

import (
	"github.com/signalfx/signalfx-agent/pkg/monitors/vsphere/model"
)

// Encapsulates services necessary to retrieve the inventory, available metrics, and build the metric index.
type VSphereInfoService struct {
	inventorySvc *InventorySvc
	metricsSvc   *MetricsSvc
}

func NewVSphereInfoService(inventorySvc *InventorySvc, metricsSvc *MetricsSvc) *VSphereInfoService {
	return &VSphereInfoService{inventorySvc: inventorySvc, metricsSvc: metricsSvc}
}

// Retrieves the inventory and available metrics and metric index.
func (l *VSphereInfoService) RetrieveVSphereInfo() (*model.VsphereInfo, error) {
	inv, err := l.retrievePopulatedInventory()
	if err != nil {
		return nil, err
	}
	idx, err := l.metricsSvc.RetrievePerfCounterIndex()
	if err != nil {
		return nil, err
	}
	return &model.VsphereInfo{Inv: inv, PerfCounterIndex: idx}, nil
}

// Retrieves the inventory and populates each inventory object with its available metrics.
func (l *VSphereInfoService) retrievePopulatedInventory() (*model.Inventory, error) {
	inv, err := l.inventorySvc.RetrieveInventory()
	if err != nil {
		return nil, err
	}
	l.metricsSvc.PopulateInvMetrics(inv)
	return inv, nil
}
