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
	"fmt"

	"github.com/pkg/errors"
	apps "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	"yunion.io/x/jsonutils"

	"yunion.io/x/onecloud/pkg/notify/options"

	"yunion.io/x/onecloud-operator/pkg/apis/constants"
	"yunion.io/x/onecloud-operator/pkg/apis/onecloud/v1alpha1"
	"yunion.io/x/onecloud-operator/pkg/controller"
	"yunion.io/x/onecloud-operator/pkg/manager"
)

const (
	NotifySocketFileDir       = "/etc/yunion/socket"
	NotifyPluginDingtalk      = "dingtalk"
	NotifyPluginConfig        = "plugin-config"
	NotifyPluginEmail         = "email"
	NotifyPluginSmsAliyun     = "smsaliyun"
	NotifyPluginWebsocket     = "websocket"
	NotifyPluginFeishu        = "feishu"
	NotifyPluginFeishuRobot   = "feishu-robot"
	NotifyPluginDingtalkRobot = "dingtalk-robot"
	NotifyPluginWorkwx        = "workwx"
	NotifyPluginWorkwxRobot   = "workwx-robot"
	NotifyPluginWebhook       = "webhook"
)

type notifyManager struct {
	*ComponentManager
}

func newNotifyManager(man *ComponentManager) manager.Manager {
	return &notifyManager{man}
}

func (m *notifyManager) getProductVersions() []v1alpha1.ProductVersion {
	return []v1alpha1.ProductVersion{
		v1alpha1.ProductVersionFullStack,
		v1alpha1.ProductVersionCMP,
		v1alpha1.ProductVersionEdge,
	}
}

func (m *notifyManager) getComponentType() v1alpha1.ComponentType {
	return v1alpha1.NotifyComponentType
}

func (m *notifyManager) Sync(oc *v1alpha1.OnecloudCluster) error {
	return syncComponent(m, oc, oc.Spec.Notify.Disable, "")
}

func (m *notifyManager) getDBConfig(cfg *v1alpha1.OnecloudClusterConfig) *v1alpha1.DBConfig {
	return &cfg.Notify.DB
}

func (m *notifyManager) getCloudUser(cfg *v1alpha1.OnecloudClusterConfig) *v1alpha1.CloudUser {
	return &cfg.Notify.CloudUser
}

func (m *notifyManager) getPhaseControl(man controller.ComponentManager, zone string) controller.PhaseControl {
	return controller.NewRegisterEndpointComponent(
		man, v1alpha1.NotifyComponentType,
		constants.ServiceNameNotify, constants.ServiceTypeNotify,
		man.GetCluster().Spec.Notify.Service.NodePort, "")
}

type NotifyPluginBaseConfig struct {
	SockFileDir string `default:"/etc/yunion/socket"`
	SenderNum   int    `default:"20"`
	LogLevel    string `default:"info"`
}

type NotifyPluginEmailConfig struct {
	NotifyPluginBaseConfig
	ChannelSize int `default:"100"`
}

type NotifyPluginWebsocketConfig struct {
	NotifyPluginBaseConfig
	Region string
}

func (m *notifyManager) getConfigMap(oc *v1alpha1.OnecloudCluster, cfg *v1alpha1.OnecloudClusterConfig, zone string) (*corev1.ConfigMap, bool, error) {
	opt := &options.Options
	if err := SetOptionsDefault(opt, constants.ServiceTypeNotify); err != nil {
		return nil, false, errors.Wrap(err, "set notify option")
	}
	config := cfg.Notify
	SetDBOptions(&opt.DBOptions, oc.Spec.Mysql, config.DB)
	SetOptionsServiceTLS(&opt.BaseOptions, false)
	SetServiceCommonOptions(&opt.CommonOptions, oc, config.ServiceCommonOptions)
	opt.SocketFileDir = NotifySocketFileDir
	opt.UpdateInterval = 30
	//opt.VerifyEmailUrlPath = fmt.Sprintf("/v2/email-verification/id/{0}/token/{1}?region=%s", oc.Spec.Region)
	opt.ReSendScope = 30
	opt.Port = config.Port

	cfgMap := m.newServiceConfigMap(v1alpha1.NotifyComponentType, "", oc, opt)
	pluginBaseOpt := &NotifyPluginBaseConfig{
		SockFileDir: NotifySocketFileDir,
		SenderNum:   5,
		LogLevel:    "info",
	}

	data := cfgMap.Data
	toStr := func(opt interface{}) string {
		return jsonutils.Marshal(opt).YAMLString()
	}
	// set plugins config
	// dingtalk and aliyunsms
	for _, pluginName := range []string{NotifyPluginDingtalk, NotifyPluginSmsAliyun, NotifyPluginFeishu,
		NotifyPluginFeishuRobot, NotifyPluginDingtalkRobot} {
		data[pluginName] = toStr(pluginBaseOpt)
	}
	// email
	data[NotifyPluginEmail] = toStr(NotifyPluginEmailConfig{
		NotifyPluginBaseConfig: *pluginBaseOpt,
		ChannelSize:            100,
	})
	// websocket
	data[NotifyPluginWebsocket] = toStr(NotifyPluginWebsocketConfig{
		NotifyPluginBaseConfig: *pluginBaseOpt,
		Region:                 oc.GetRegion(),
	})

	cfgMap.Data = data

	return cfgMap, false, nil
}

func (m *notifyManager) getService(oc *v1alpha1.OnecloudCluster, cfg *v1alpha1.OnecloudClusterConfig, zone string) []*corev1.Service {
	return []*corev1.Service{m.newSingleNodePortService(v1alpha1.NotifyComponentType, oc, int32(oc.Spec.Notify.Service.NodePort), int32(cfg.Notify.Port))}
}

func (m *notifyManager) getDeployment(oc *v1alpha1.OnecloudCluster, cfg *v1alpha1.OnecloudClusterConfig, zone string) (*apps.Deployment, error) {
	deploy, err := m.newCloudServiceSinglePortDeployment(v1alpha1.NotifyComponentType, "", oc, &oc.Spec.Notify.DeploymentSpec, int32(cfg.Notify.Port), true, false)
	if err != nil {
		return nil, err
	}
	socketVol := corev1.Volume{
		Name: "socket",
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		},
	}
	socketVolMount := corev1.VolumeMount{
		Name:      "socket",
		MountPath: NotifySocketFileDir,
	}
	cfgMapName := controller.ComponentConfigMapName(oc, v1alpha1.NotifyComponentType)
	pluginCfgVol := corev1.Volume{
		Name: NotifyPluginConfig,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: cfgMapName,
				},
				Items: []corev1.KeyToPath{
					{Key: NotifyPluginDingtalk, Path: fmt.Sprintf("%s.conf", NotifyPluginDingtalk)},
					{Key: NotifyPluginEmail, Path: fmt.Sprintf("%s.conf", NotifyPluginEmail)},
					{Key: NotifyPluginSmsAliyun, Path: fmt.Sprintf("%s.conf", NotifyPluginSmsAliyun)},
					{Key: NotifyPluginWebsocket, Path: fmt.Sprintf("%s.conf", NotifyPluginWebsocket)},
					{Key: NotifyPluginFeishu, Path: fmt.Sprintf("%s.conf", NotifyPluginFeishu)},
					{Key: NotifyPluginFeishuRobot, Path: fmt.Sprintf("%s.conf", NotifyPluginFeishuRobot)},
					{Key: NotifyPluginDingtalkRobot, Path: fmt.Sprintf("%s.conf", NotifyPluginDingtalkRobot)},
				},
			},
		},
	}
	newPluginC := func(name string) corev1.Container {
		return corev1.Container{
			Name:            name,
			Image:           oc.Spec.Notify.Plugins.Image,
			ImagePullPolicy: oc.Spec.Notify.Plugins.ImagePullPolicy,
			Command:         []string{fmt.Sprintf("/opt/yunion/bin/%s", name), "--config", fmt.Sprintf("/etc/yunion/%s.conf", name)},
			VolumeMounts: []corev1.VolumeMount{
				socketVolMount,
				{
					MountPath: constants.ConfigDir,
					Name:      NotifyPluginConfig,
				},
			},
		}
	}
	newPluginCWithoutConf := func(name string) corev1.Container {
		return corev1.Container{
			Name:            name,
			Image:           oc.Spec.Notify.Plugins.Image,
			ImagePullPolicy: oc.Spec.Notify.Plugins.ImagePullPolicy,
			Command:         []string{fmt.Sprintf("/opt/yunion/bin/%s", name), "--log-level", "info", "--sock-file-dir", NotifySocketFileDir},
			VolumeMounts: []corev1.VolumeMount{
				socketVolMount,
			},
		}
	}
	pluginCs := []corev1.Container{
		newPluginC(NotifyPluginDingtalk),
		newPluginC(NotifyPluginEmail),
		newPluginC(NotifyPluginSmsAliyun),
		newPluginC(NotifyPluginFeishu),
		newPluginC(NotifyPluginFeishuRobot),
		newPluginC(NotifyPluginDingtalkRobot),
		newPluginCWithoutConf(NotifyPluginWorkwx),
		newPluginCWithoutConf(NotifyPluginWorkwxRobot),
		newPluginCWithoutConf(NotifyPluginWebhook),
	}
	isEE := IsEnterpriseEdition(oc)
	if isEE {
		pluginCs = append(pluginCs, newPluginC(NotifyPluginWebsocket))
	}
	spec := &deploy.Spec.Template.Spec
	cs := spec.Containers
	notifyCs := &cs[0]
	notifyCs.VolumeMounts = append(notifyCs.VolumeMounts, socketVolMount)
	cs = append(cs, pluginCs...)
	spec.Containers = cs
	vols := spec.Volumes
	vols = append(vols, socketVol, pluginCfgVol)
	spec.Volumes = vols
	return deploy, nil
}

func (m *notifyManager) getDeploymentStatus(oc *v1alpha1.OnecloudCluster, zone string) *v1alpha1.DeploymentStatus {
	return &oc.Status.Notify
}
