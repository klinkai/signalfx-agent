package vsphere

import (
	"context"

	"github.com/signalfx/signalfx-agent/pkg/monitors/vsphere/model"
	log "github.com/sirupsen/logrus"
)

type runner struct {
	ctx                   context.Context
	monitor               *Monitor
	conf                  *model.Config
	vsm                   *vSphereMonitor
	vsphereReloadInterval int //seconds
}

func newRunner(ctx context.Context, conf *model.Config, monitor *Monitor) runner {
	vsphereReloadInterval := model.DefaultVSInfoReloadInterval
	if conf.InventoryRefreshIntervalSeconds > 0 {
		vsphereReloadInterval = conf.InventoryRefreshIntervalSeconds
	}
	return runner{
		ctx:                   ctx,
		monitor:               monitor,
		conf:                  conf,
		vsphereReloadInterval: vsphereReloadInterval,
		vsm:                   &vSphereMonitor{},
	}
}

// Called periodically. This is the entry point to the vSphere montor.
func (r *runner) run() {
	err := r.vsm.firstTimeSetup(r.ctx, r.conf)
	if err != nil {
		log.WithError(err).Error("firstTimeSetup failed")
		return
	}
	dps := r.vsm.retrieveDatapoints()
	r.monitor.Output.SendDatapoints(dps...)
	if r.vsm.isTimeForVSphereInfoReload(r.vsphereReloadInterval) {
		r.vsm.reloadVSphereInfo()
	}
}
