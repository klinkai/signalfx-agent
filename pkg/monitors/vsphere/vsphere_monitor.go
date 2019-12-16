package vsphere

import (
	"context"
	"time"

	"github.com/signalfx/golib/v3/datapoint"
	"github.com/signalfx/signalfx-agent/pkg/monitors/vsphere/model"
	"github.com/signalfx/signalfx-agent/pkg/monitors/vsphere/service"
	"github.com/sirupsen/logrus"
	"github.com/vmware/govmomi"
)

// Encapsulates services and the current state of the vSphere monitor.
type vSphereMonitor struct {
	invSvc    *service.InventorySvc
	metricSvc *service.MetricsSvc
	vsInfoSvc *service.VSphereInfoService
	ptsSvc    *service.PointsSvc

	vSphereInfo              *model.VsphereInfo
	lastVsphereLoadTime      time.Time
	latestPointRetrievalTime time.Time
}

// Logs into vSphere, wires up service objects, and retrieves vSphereInfo (inventory, available metrics, and metric index).
func (vsm *vSphereMonitor) firstTimeSetup(ctx context.Context, conf *model.Config) error {
	if !vsm.lastVsphereLoadTime.IsZero() {
		return nil
	}
	client, err := service.LogIn(ctx, conf)
	if err != nil {
		return err
	}

	vsm.wireUpServices(ctx, client)

	vsm.vSphereInfo, err = vsm.vsInfoSvc.RetrieveVSphereInfo()
	if err != nil {
		return err
	}
	vsm.lastVsphereLoadTime = time.Now()
	return nil
}

// Creates the service objects and assigns them to the vSphereMonitor struct.
func (vsm *vSphereMonitor) wireUpServices(ctx context.Context, client *govmomi.Client) {
	gateway := service.NewGateway(ctx, client)
	vsm.ptsSvc = service.NewPointsSvc(gateway)
	vsm.invSvc = service.NewInventorySvc(gateway)
	vsm.metricSvc = service.NewMetricsService(gateway)
	vsm.vsInfoSvc = service.NewVSphereInfoService(vsm.invSvc, vsm.metricSvc)
}

// Retrieves datapoints for all the inventory for the number of 20-second intervals available since the last datapoint
// retrieval.
func (vsm *vSphereMonitor) retrieveDatapoints() []*datapoint.Datapoint {
	numSamples := vsm.getNumSamplesReqd()
	if numSamples == 0 {
		return nil
	}
	dps, latestRetrievalTime := vsm.ptsSvc.RetrievePoints(vsm.vSphereInfo, numSamples)
	if !latestRetrievalTime.IsZero() {
		vsm.latestPointRetrievalTime = latestRetrievalTime
	}
	return dps
}

func (vsm *vSphereMonitor) getNumSamplesReqd() int32 {
	return getNumSamplesReqd(vsm.latestPointRetrievalTime)
}

// Traverses the vSphere inventory and saves the result in vSphereInfo (hosts, VMs, available metrics, and metric index).
func (vsm *vSphereMonitor) reloadVSphereInfo() {
	var err error
	vsm.vSphereInfo, err = vsm.vsInfoSvc.RetrieveVSphereInfo()
	if err != nil {
		service.Log.WithError(err).Error("Failed to load vSphereInfo")
		return
	}
	vsm.lastVsphereLoadTime = time.Now()
}

// Compares the last vSphereInfo load time to the vSphere info reload interval, returning whether more time has elapsed
// than the configured duration.
func (vsm *vSphereMonitor) isTimeForVSphereInfoReload(vsphereReloadIntervalSeconds int) bool {
	secondsSinceLastVsReload := int(time.Since(vsm.lastVsphereLoadTime).Seconds())
	timeForReload := secondsSinceLastVsReload > vsphereReloadIntervalSeconds
	service.Log.WithFields(logrus.Fields{
		"secondsSinceLastVsReload": secondsSinceLastVsReload,
		"vsphereReloadInterval":    vsphereReloadIntervalSeconds,
	}).Debugf("Time for vs reload = %t", timeForReload)
	return timeForReload
}

// Returns the number of 20-second intervals available in vSphere since the passed-in Time. Assumes reasonably well-synced
// clocks between this monitor's host and the vCenter Server.
func getNumSamplesReqd(lastInterval time.Time) int32 {
	if lastInterval.IsZero() {
		return 1
	}
	fSecondsSinceLastInterval := time.Since(lastInterval).Seconds()
	intSecondsSinceLastInterval := int32(fSecondsSinceLastInterval)
	numSamples := intSecondsSinceLastInterval / model.VSMetricsInterval
	service.Log.WithFields(logrus.Fields{
		"now":                         time.Now(),
		"lastInterval":                lastInterval,
		"fSecondsSinceLastInterval":   fSecondsSinceLastInterval,
		"intSecondsSinceLastInterval": intSecondsSinceLastInterval,
	}).Debugf("numSamples = %d", numSamples)
	return numSamples
}
