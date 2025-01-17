// Copyright 2019 Yunion
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package component

import (
	apps "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	"yunion.io/x/pkg/errors"

	"yunion.io/x/onecloud/pkg/scheduler/options"

	"yunion.io/x/onecloud-operator/pkg/apis/constants"
	"yunion.io/x/onecloud-operator/pkg/apis/onecloud/v1alpha1"
	"yunion.io/x/onecloud-operator/pkg/controller"
	"yunion.io/x/onecloud-operator/pkg/manager"
)

type schedulerManager struct {
	*ComponentManager
}

func newSchedulerManager(man *ComponentManager) manager.Manager {
	return &schedulerManager{man}
}

func (m *schedulerManager) getProductVersions() []v1alpha1.ProductVersion {
	return []v1alpha1.ProductVersion{
		v1alpha1.ProductVersionFullStack,
		v1alpha1.ProductVersionCMP,
		v1alpha1.ProductVersionEdge,
	}
}

func (m *schedulerManager) getComponentType() v1alpha1.ComponentType {
	return v1alpha1.SchedulerComponentType
}

func (m *schedulerManager) Sync(oc *v1alpha1.OnecloudCluster) error {
	return syncComponent(m, oc, oc.Spec.Scheduler.Disable, "")
}

func (m *schedulerManager) getPhaseControl(man controller.ComponentManager, zone string) controller.PhaseControl {
	return controller.NewRegisterEndpointComponent(man, v1alpha1.SchedulerComponentType,
		constants.ServiceNameScheduler, constants.ServiceTypeScheduler, man.GetCluster().Spec.Scheduler.Service.NodePort, "")
}

func (m *schedulerManager) getConfigMap(oc *v1alpha1.OnecloudCluster, cfg *v1alpha1.OnecloudClusterConfig, zone string) (*corev1.ConfigMap, bool, error) {
	opt := &options.Options
	if err := SetOptionsDefault(opt, constants.ServiceTypeScheduler); err != nil {
		return nil, false, errors.Wrap(err, "scheduler: SetOptionsDefault")
	}
	// scheduler use region config directly
	config := cfg.RegionServer
	SetDBOptions(&opt.DBOptions, oc.Spec.Mysql, config.DB)
	SetOptionsServiceTLS(&opt.BaseOptions, false)
	SetServiceCommonOptions(&opt.CommonOptions, oc, config.ServiceDBCommonOptions.ServiceCommonOptions)

	opt.SchedulerPort = constants.SchedulerPort
	return m.newServiceConfigMap(v1alpha1.SchedulerComponentType, "", oc, opt), false, nil
}

func (m *schedulerManager) getService(oc *v1alpha1.OnecloudCluster, cfg *v1alpha1.OnecloudClusterConfig, zone string) []*corev1.Service {
	return []*corev1.Service{m.newSingleNodePortService(v1alpha1.SchedulerComponentType, oc, int32(oc.Spec.Scheduler.Service.NodePort), constants.SchedulerPort)}
}

func (m *schedulerManager) getDeployment(oc *v1alpha1.OnecloudCluster, cfg *v1alpha1.OnecloudClusterConfig, zone string) (*apps.Deployment, error) {
	return m.newCloudServiceSinglePortDeployment(v1alpha1.SchedulerComponentType, "", oc, &oc.Spec.Scheduler.DeploymentSpec, constants.SchedulerPort, false, false)
}

func (m *schedulerManager) getDeploymentStatus(oc *v1alpha1.OnecloudCluster, zone string) *v1alpha1.DeploymentStatus {
	return &oc.Status.Scheduler
}
