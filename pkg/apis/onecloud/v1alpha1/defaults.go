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

package v1alpha1

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/rand"

	"yunion.io/x/log"

	"yunion.io/x/onecloud-operator/pkg/apis/constants"
	"yunion.io/x/onecloud-operator/pkg/util/passwd"
)

const (
	DefaultVersion                 = "latest"
	DefaultOnecloudRegion          = "region0"
	DefaultOnecloudRegionDNSDomain = "cloud.onecloud.io"
	DefaultOnecloudZone            = "zone0"
	DefaultOnecloudWire            = "bcast0"
	DefaultImageRepository         = "registry.hub.docker.com/yunion"
	DefaultVPCId                   = "default"
	DefaultGlanceStorageSize       = "100G"
	DefaultMeterStorageSize        = "100G"
	DefaultInfluxdbStorageSize     = "20G"
	DefaultNotifyStorageSize       = "1G" // for plugin template
	DefaultBaremetalStorageSize    = "100G"
	DefaultEsxiAgentStorageSize    = "30G"
	// rancher local-path-provisioner: https://github.com/rancher/local-path-provisioner
	DefaultStorageClass = "local-path"

	DefaultOvnVersion   = "2.10.5"
	DefaultOvnImageName = "openvswitch"
	DefaultOvnImageTag  = DefaultOvnVersion + "-1"

	DefaultSdnAgentImageName = "sdnagent"

	DefaultWebOverviewImageName = "dashboard-overview"
	DefaultWebDocsImageName     = "docs-ee"

	DefaultNotifyPluginsImageName = "notify-plugins"

	DefaultHostImageName = "host-image"
	DefaultHostImageTag  = "v1.0.4"

	DefaultInfluxdbImageVersion = "1.7.7"

	DefaultTelegrafImageName     = "telegraf"
	DefaultTelegrafImageTag      = "release-1.19.2-0"
	DefaultTelegrafInitImageName = "telegraf-init"
	DefaultTelegrafInitImageTag  = "release-1.19.2-0"
	DefaultTelegrafRaidImageName = "telegraf-raid-plugin"
	DefaultTelegrafRaidImageTag  = "release-1.6.4"

	DefaultHostQemuVersion = "4.2.0"
)

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	return RegisterDefaults(scheme)
}

func SetDefaults_OnecloudCluster(obj *OnecloudCluster) {
	if _, ok := obj.GetLabels()[constants.InstanceLabelKey]; !ok {
		obj.SetLabels(map[string]string{constants.InstanceLabelKey: fmt.Sprintf("onecloud-cluster-%s", rand.String(4))})
	}

	SetDefaults_OnecloudClusterSpec(&obj.Spec, IsEnterpriseEdition(obj))
}

func GetEdition(oc *OnecloudCluster) string {
	edition := constants.OnecloudCommunityEdition
	if oc.Annotations == nil {
		return edition
	}
	curEdition := oc.Annotations[constants.OnecloudEditionAnnotationKey]
	if curEdition == constants.OnecloudEnterpriseEdition {
		return curEdition
	}
	return edition
}

func IsEnterpriseEdition(oc *OnecloudCluster) bool {
	return GetEdition(oc) == constants.OnecloudEnterpriseEdition
}

type hyperImagePair struct {
	*DeploymentSpec
	Supported bool
}

func newHyperImagePair(ds *DeploymentSpec, supported bool) *hyperImagePair {
	return &hyperImagePair{
		DeploymentSpec: ds,
		Supported:      supported,
	}
}

func SetDefaults_OnecloudClusterSpec(obj *OnecloudClusterSpec, isEE bool) {
	setDefaults_Mysql(&obj.Mysql)
	if obj.ProductVersion == "" {
		obj.ProductVersion = ProductVersionFullStack
	}
	if obj.Region == "" {
		obj.Region = DefaultOnecloudRegion
	}
	if obj.Zone == "" {
		obj.Zone = DefaultOnecloudZone
	}
	if obj.Version == "" {
		obj.Version = DefaultVersion
	}
	if obj.ImageRepository == "" {
		obj.ImageRepository = DefaultImageRepository
	}

	useHyperImage := obj.UseHyperImage

	SetDefaults_KeystoneSpec(&obj.Keystone, obj.ImageRepository, obj.Version, useHyperImage, isEE)
	SetDefaults_RegionSpec(&obj.RegionServer, obj.ImageRepository, obj.Version, useHyperImage, isEE)
	SetDefaults_RegionDNSSpec(&obj.RegionDNS, obj.ImageRepository, obj.Version)

	nHP := newHyperImagePair

	for cType, spec := range map[ComponentType]*hyperImagePair{
		WebconsoleComponentType:      nHP(&obj.Webconsole, useHyperImage),
		SchedulerComponentType:       nHP(&obj.Scheduler, useHyperImage),
		LoggerComponentType:          nHP(&obj.Logger, useHyperImage),
		YunionconfComponentType:      nHP(&obj.Yunionconf, useHyperImage),
		KubeServerComponentType:      nHP(&obj.KubeServer, false),
		AnsibleServerComponentType:   nHP(&obj.AnsibleServer, useHyperImage),
		CloudnetComponentType:        nHP(&obj.Cloudnet, false),
		CloudproxyComponentType:      nHP(&obj.Cloudproxy, false),
		CloudeventComponentType:      nHP(&obj.Cloudevent, useHyperImage),
		S3gatewayComponentType:       nHP(&obj.S3gateway, false),
		DevtoolComponentType:         nHP(&obj.Devtool, useHyperImage),
		AutoUpdateComponentType:      nHP(&obj.AutoUpdate, useHyperImage),
		OvnNorthComponentType:        nHP(&obj.OvnNorth, false),
		VpcAgentComponentType:        nHP(&obj.VpcAgent, false),
		MonitorComponentType:         nHP(&obj.Monitor, useHyperImage),
		ServiceOperatorComponentType: nHP(&obj.ServiceOperator, false),
		ItsmComponentType:            nHP(&obj.Itsm, false),
		CloudIdComponentType:         nHP(&obj.CloudId, useHyperImage),
		SuggestionComponentType:      nHP(&obj.Suggestion, useHyperImage),
		CloudmonComponentType:        nHP(&obj.Cloudmon.DeploymentSpec, useHyperImage),
		ScheduledtaskComponentType:   nHP(&obj.Scheduledtask, useHyperImage),
		ReportComponentType:          nHP(&obj.Report, useHyperImage),
		MspOperationComponentType:    nHP(&obj.MspOperation, useHyperImage),
	} {
		SetDefaults_DeploymentSpec(spec.DeploymentSpec, getImage(
			obj.ImageRepository, spec.Repository,
			cType, spec.ImageName,
			obj.Version, spec.Tag,
			spec.Supported, isEE,
		))
	}

	// CE or EE parts
	for cType, spec := range map[ComponentType]*hyperImagePair{
		APIGatewayComponentType: nHP(&obj.APIGateway.DeploymentSpec, useHyperImage),
		ClimcComponentType:      nHP(&obj.Climc, false),
		WebComponentType:        nHP(&obj.Web.DeploymentSpec, false),
	} {
		SetDefaults_DeploymentSpec(spec.DeploymentSpec,
			getEditionImage(
				obj.ImageRepository, spec.Repository,
				cType, spec.ImageName,
				obj.Version, spec.Tag,
				spec.Supported, isEE,
			))
	}

	// setting web overview image
	obj.Web.Overview.Image = getImage(
		obj.ImageRepository, obj.Web.Overview.Repository,
		DefaultWebOverviewImageName, obj.Web.Overview.ImageName,
		obj.Version, obj.Web.Overview.Tag,
		false, isEE,
	)
	obj.Web.Overview.ImagePullPolicy = corev1.PullIfNotPresent
	// setting web docs image
	obj.Web.Docs.Image = getImage(
		obj.ImageRepository, obj.Web.Docs.Repository,
		DefaultWebDocsImageName, obj.Web.Docs.ImageName,
		obj.Version, obj.Web.Docs.Tag,
		false, isEE,
	)
	obj.Web.Docs.ImagePullPolicy = corev1.PullIfNotPresent

	// setting host spec
	if obj.HostAgent.DefaultQemuVersion == "" {
		obj.HostAgent.DefaultQemuVersion = DefaultHostQemuVersion
	}

	for cType, spec := range map[ComponentType]*DaemonSetSpec{
		HostComponentType:         &obj.HostAgent.DaemonSetSpec,
		HostDeployerComponentType: &obj.HostDeployer,
		YunionagentComponentType:  &obj.Yunionagent,
		HostImageComponentType:    &obj.HostImage,
	} {
		useHI := false
		if cType == YunionagentComponentType {
			useHI = useHyperImage
		}
		SetDefaults_DaemonSetSpec(spec, getImage(
			obj.ImageRepository, spec.Repository,
			cType, spec.ImageName,
			obj.Version, spec.Tag,
			useHI, isEE,
		))
	}

	// setting sdnagent image
	obj.HostAgent.SdnAgent.Image = getImage(
		obj.ImageRepository, obj.HostAgent.SdnAgent.Repository,
		DefaultSdnAgentImageName, obj.HostAgent.SdnAgent.ImageName,
		obj.Version, obj.HostAgent.SdnAgent.Tag,
		false, isEE,
	)
	obj.HostAgent.SdnAgent.ImagePullPolicy = corev1.PullIfNotPresent

	// setting ovn image
	obj.HostAgent.OvnController.Image = getImage(
		obj.ImageRepository, obj.HostAgent.OvnController.Repository,
		DefaultOvnImageName, obj.HostAgent.OvnController.ImageName,
		DefaultOvnImageTag, obj.HostAgent.OvnController.Tag,
		false, isEE,
	)
	obj.HostAgent.OvnController.ImagePullPolicy = corev1.PullIfNotPresent

	obj.OvnNorth.Image = getImage(
		obj.ImageRepository, obj.OvnNorth.Repository,
		DefaultOvnImageName, obj.OvnNorth.ImageName,
		DefaultOvnImageTag, obj.OvnNorth.Tag,
		false, isEE,
	)
	obj.OvnNorth.ImagePullPolicy = corev1.PullIfNotPresent
	// host-image
	obj.HostImage.Image = getImage(
		obj.ImageRepository, obj.HostImage.Repository,
		DefaultHostImageName, obj.HostImage.ImageName,
		DefaultHostImageTag, obj.HostImage.Tag,
		false, isEE,
	)

	// telegraf spec
	obj.Telegraf.InitContainerImage = getImage(
		obj.ImageRepository, obj.Telegraf.Repository,
		DefaultTelegrafInitImageName, "",
		DefaultTelegrafInitImageTag, "",
		false, isEE,
	)
	obj.Telegraf.TelegrafRaidImage = getImage(
		obj.ImageRepository, obj.Telegraf.Repository,
		DefaultTelegrafRaidImageName, "",
		DefaultTelegrafRaidImageTag, "",
		false, isEE,
	)
	SetDefaults_DaemonSetSpec(
		&obj.Telegraf.DaemonSetSpec,
		getImage(obj.ImageRepository, obj.Telegraf.Repository,
			DefaultTelegrafImageName, obj.Telegraf.ImageName,
			DefaultTelegrafImageTag, obj.Telegraf.Tag,
			false, isEE,
		),
	)

	type stateDeploy struct {
		obj           *StatefulDeploymentSpec
		size          string
		version       string
		useHyperImage bool
	}
	for cType, spec := range map[ComponentType]*stateDeploy{
		GlanceComponentType:         {&obj.Glance.StatefulDeploymentSpec, DefaultGlanceStorageSize, obj.Version, useHyperImage},
		InfluxdbComponentType:       {&obj.Influxdb, DefaultInfluxdbStorageSize, DefaultInfluxdbImageVersion, false},
		NotifyComponentType:         {&obj.Notify.StatefulDeploymentSpec, DefaultNotifyStorageSize, obj.Version, useHyperImage},
		BaremetalAgentComponentType: {&obj.BaremetalAgent.StatefulDeploymentSpec, DefaultBaremetalStorageSize, obj.Version, false},
		MeterComponentType:          {&obj.Meter, DefaultMeterStorageSize, obj.Version, useHyperImage},
		EsxiAgentComponentType:      {&obj.EsxiAgent.StatefulDeploymentSpec, DefaultEsxiAgentStorageSize, obj.Version, useHyperImage},
	} {
		SetDefaults_StatefulDeploymentSpec(cType, spec.obj, spec.size, obj.ImageRepository, spec.version, spec.useHyperImage, isEE)
	}

	// setting web overview image
	obj.Notify.Plugins.Image = getImage(
		obj.ImageRepository, obj.Notify.Plugins.Repository,
		DefaultNotifyPluginsImageName, obj.Notify.Plugins.ImageName,
		obj.Version, obj.Notify.Plugins.Tag,
		useHyperImage, isEE,
	)
	obj.Notify.Plugins.ImagePullPolicy = corev1.PullIfNotPresent

	// cloudmon spec
	//SetDefaults_DeploymentSpec(&obj.Cloudmon.DeploymentSpec,
	//	getImage(obj.ImageRepository, obj.Cloudmon.Repository, APIGatewayComponentTypeEE,
	//		obj.Cloudmon.ImageName, obj.Version, obj.APIGateway.Tag))
	if obj.Cloudmon.CloudmonReportUsageDuration == 0 {
		obj.Cloudmon.CloudmonReportUsageDuration = 15
	}
	if obj.Cloudmon.CloudmonReportServerDuration == 0 {
		obj.Cloudmon.CloudmonReportServerDuration = 4
	}
	if obj.Cloudmon.CloudmonReportHostDuration == 0 {
		obj.Cloudmon.CloudmonReportHostDuration = 4
	}
	if obj.Cloudmon.CloudmonPingDuration == 0 {
		obj.Cloudmon.CloudmonPingDuration = 60
	}
	if obj.Cloudmon.CloudmonReportCloudAccountDuration == 0 {
		obj.Cloudmon.CloudmonReportCloudAccountDuration = 30
	}
	if obj.Cloudmon.CloudmonReportAlertRecordHistoryDuration == 0 {
		obj.Cloudmon.CloudmonReportAlertRecordHistoryDuration = 1
	}

	setDefaults_MonitorStackSpec(&obj.MonitorStack)
}

func setDefaults_Mysql(obj *Mysql) {
	if obj.Username == "" {
		obj.Username = "root"
	}
	if obj.Port == 0 {
		obj.Port = 3306
	}
}

func getHyperImageName(isEE bool) string {
	if isEE {
		return "cloudpods-ee"
	}
	return "cloudpods"
}

func getImage(
	globalRepo, specRepo string,
	componentType ComponentType, componentName string,
	globalVersion, tag string,
	useHyperImage bool, isEE bool) string {
	repo := specRepo
	if specRepo == "" {
		repo = globalRepo
	}
	version := tag
	if version == "" {
		version = globalVersion
	}
	if componentName == "" {
		if useHyperImage {
			componentName = getHyperImageName(isEE)
		} else {
			componentName = componentType.String()
		}
	}
	return fmt.Sprintf("%s/%s:%s", repo, componentName, version)
}

func getEditionImage(globalRepo, specRepo string, componentType ComponentType, componentName string, globalVersion, tag string, useHyperImage bool, isEE bool) string {
	if componentName == "" {
		componentName = componentType.String()
		if isEE {
			componentName = fmt.Sprintf("%s-ee", componentName)
		}
		if useHyperImage {
			componentName = getHyperImageName(isEE)
		}
	}
	return getImage(globalRepo, specRepo, componentType, componentName, globalVersion, tag, useHyperImage, isEE)
}

func SetDefaults_KeystoneSpec(
	obj *KeystoneSpec,
	imageRepo, version string,
	useHyperImage, isEE bool) {
	SetDefaults_DeploymentSpec(&obj.DeploymentSpec, getImage(
		imageRepo, obj.Repository,
		KeystoneComponentType, obj.ImageName,
		version, obj.Tag,
		useHyperImage, isEE,
	))
	if obj.BootstrapPassword == "" {
		obj.BootstrapPassword = passwd.GeneratePassword()
	}
}

func SetDefaults_RegionSpec(
	obj *RegionSpec,
	imageRepo, version string,
	useHyperImage, isEE bool,
) {
	SetDefaults_DeploymentSpec(&obj.DeploymentSpec, getImage(
		imageRepo, obj.Repository,
		RegionComponentType, obj.ImageName,
		version, obj.Tag,
		useHyperImage, isEE,
	))
	if obj.DNSDomain == "" {
		obj.DNSDomain = DefaultOnecloudRegionDNSDomain
	}
}

func SetDefaults_RegionDNSSpec(obj *RegionDNSSpec, imageRepo, version string) {
	SetDefaults_DaemonSetSpec(&obj.DaemonSetSpec, getImage(
		imageRepo, obj.Repository,
		RegionDNSComponentType, obj.ImageName,
		version, obj.Tag,
		false, false,
	))
}

func setPVCStoreage(obj *ContainerSpec, size string) {
	if obj.Requests == nil {
		obj.Requests = new(ResourceRequirement)
	}
	if obj.Requests.Storage == "" {
		obj.Requests.Storage = size
	}
}

func SetDefaults_StatefulDeploymentSpec(
	ctype ComponentType,
	obj *StatefulDeploymentSpec, defaultSize string,
	imageRepo, version string,
	useHyperImage bool, isEE bool) {
	SetDefaults_DeploymentSpec(&obj.DeploymentSpec, getImage(imageRepo, obj.Repository, ctype, obj.ImageName, version, obj.Tag, useHyperImage, isEE))
	if obj.StorageClassName == "" {
		obj.StorageClassName = DefaultStorageClass
	}
	setPVCStoreage(&obj.ContainerSpec, defaultSize)
}

func SetDefaults_DeploymentSpec(obj *DeploymentSpec, image string) {
	if obj.Replicas <= 0 {
		obj.Replicas = 1
	}
	if obj.Disable {
		obj.Replicas = 0
	}
	obj.Image = image
	if string(obj.ImagePullPolicy) == "" {
		obj.ImagePullPolicy = corev1.PullIfNotPresent
	}
	// add tolerations
	if len(obj.Tolerations) == 0 {
		obj.Tolerations = append(obj.Tolerations, []corev1.Toleration{
			{
				Key:    "node-role.kubernetes.io/master",
				Effect: corev1.TaintEffectNoSchedule,
			},
			{
				Key:    "node-role.kubernetes.io/controlplane",
				Effect: corev1.TaintEffectNoSchedule,
			},
		}...)
	}
}

func SetDefaults_DaemonSetSpec(obj *DaemonSetSpec, image string) {
	obj.Image = image
	if string(obj.ImagePullPolicy) == "" {
		obj.ImagePullPolicy = corev1.PullIfNotPresent
	}
	if len(obj.Tolerations) == 0 {
		obj.Tolerations = append(obj.Tolerations, []corev1.Toleration{
			{
				Key:    "node-role.kubernetes.io/master",
				Effect: corev1.TaintEffectNoSchedule,
			},
			{
				Key:    "node-role.kubernetes.io/controlplane",
				Effect: corev1.TaintEffectNoSchedule,
			},
		}...)
	}
}

func SetDefaults_CronJobSpec(obj *CronJobSpec, image string) {
	obj.Image = image
	if string(obj.ImagePullPolicy) == "" {
		obj.ImagePullPolicy = corev1.PullIfNotPresent
	}
	if len(obj.Tolerations) == 0 {
		obj.Tolerations = append(obj.Tolerations, []corev1.Toleration{
			{
				Key:    "node-role.kubernetes.io/master",
				Effect: corev1.TaintEffectNoSchedule,
			},
			{
				Key:    "node-role.kubernetes.io/controlplane",
				Effect: corev1.TaintEffectNoSchedule,
			},
		}...)
	}
}

func setDefaults_MonitorStackSpec(obj *MonitorStackSpec) {
	// set minio defaults
	minio := &obj.Minio
	if minio.AccessKey == "" {
		minio.AccessKey = "monitor-admin"
	}
	if minio.SecretKey == "" {
		minio.SecretKey = passwd.GeneratePassword()
	}

	// set grafana defaults
	grafana := &obj.Grafana
	if grafana.AdminUser == "" {
		grafana.AdminUser = "admin"
	}
	if grafana.AdminPassword == "" {
		grafana.AdminPassword = "admin@123"
	}
	if grafana.Subpath == "" {
		grafana.Subpath = "grafana"
	}
	oauth := &grafana.OAuth
	if oauth.Scopes == "" {
		oauth.Scopes = "user"
	}
	if oauth.RoleAttributePath == "" {
		oauth.RoleAttributePath = `projectName == 'system' && contains(roles, 'admin') && 'Admin' || 'Editor'`
	}

	// set loki defaults
	loki := &obj.Loki
	if loki.ObjectStoreConfig.Bucket == "" {
		loki.ObjectStoreConfig.Bucket = "loki"
	}

	// set prometheus defaults
	prom := &obj.Prometheus
	if prom.ThanosSidecarSpec.ObjectStoreConfig.Bucket == "" {
		prom.ThanosSidecarSpec.ObjectStoreConfig.Bucket = "thanos"
	}
	trueVar := true
	if prom.Disable == nil {
		prom.Disable = &trueVar
	}

	// set thanos defaults
	thanos := &obj.Thanos
	if thanos.ObjectStoreConfig.Bucket == "" {
		thanos.ObjectStoreConfig.Bucket = "thanos"
	}
}

func SetDefaults_OnecloudClusterConfig(obj *OnecloudClusterConfig) {
	setDefaults_KeystoneConfig(&obj.Keystone)

	type userPort struct {
		user string
		port int
	}

	for opt, userPort := range map[*ServiceCommonOptions]userPort{
		&obj.APIGateway:                          {constants.APIGatewayAdminUser, constants.APIGatewayPort},
		&obj.HostAgent.ServiceCommonOptions:      {constants.HostAdminUser, constants.HostPort},
		&obj.BaremetalAgent.ServiceCommonOptions: {constants.BaremetalAdminUser, constants.BaremetalPort},
		&obj.S3gateway:                           {constants.S3gatewayAdminUser, constants.S3gatewayPort},
		&obj.AutoUpdate:                          {constants.AutoUpdateAdminUser, constants.AutoUpdatePort},
		&obj.EsxiAgent.ServiceCommonOptions:      {constants.EsxiAgentAdminUser, constants.EsxiAgentPort},
		&obj.VpcAgent.ServiceCommonOptions:       {constants.VpcAgentAdminUser, 0},
		&obj.ServiceOperator:                     {constants.ServiceOperatorAdminUser, constants.ServiceOperatorPort},
	} {
		SetDefaults_ServiceCommonOptions(opt, userPort.user, userPort.port)
	}

	type userDBPort struct {
		user   string
		port   int
		db     string
		dbUser string
	}

	registryPorts := map[int]string{}

	for opt, tmp := range map[*ServiceDBCommonOptions]userDBPort{
		&obj.RegionServer.ServiceDBCommonOptions: {constants.RegionAdminUser, constants.RegionPort, constants.RegionDB, constants.RegionDBUser},
		&obj.Glance.ServiceDBCommonOptions:       {constants.GlanceAdminUser, constants.GlanceAPIPort, constants.GlanceDB, constants.GlanceDBUser},
		&obj.Logger:                              {constants.LoggerAdminUser, constants.LoggerPort, constants.LoggerDB, constants.LoggerDBUser},
		&obj.Yunionagent:                         {constants.YunionAgentAdminUser, constants.YunionAgentPort, constants.YunionAgentDB, constants.YunionAgentDBUser},
		&obj.Yunionconf:                          {constants.YunionConfAdminUser, constants.YunionConfPort, constants.YunionConfDB, constants.YunionConfDBUser},
		&obj.KubeServer:                          {constants.KubeServerAdminUser, constants.KubeServerPort, constants.KubeServerDB, constants.KubeServerDBUser},
		&obj.AnsibleServer:                       {constants.AnsibleServerAdminUser, constants.AnsibleServerPort, constants.AnsibleServerDB, constants.AnsibleServerDBUser},
		&obj.Cloudnet:                            {constants.CloudnetAdminUser, constants.CloudnetPort, constants.CloudnetDB, constants.CloudnetDBUser},
		&obj.Cloudproxy:                          {constants.CloudproxyAdminUser, constants.CloudproxyPort, constants.CloudproxyDB, constants.CloudproxyDBUser},
		&obj.Cloudevent:                          {constants.CloudeventAdminUser, constants.CloudeventPort, constants.CloudeventDB, constants.CloudeventDBUser},
		&obj.Notify:                              {constants.NotifyAdminUser, constants.NotifyPort, constants.NotifyDB, constants.NotifyDBUser},
		&obj.Devtool:                             {constants.DevtoolAdminUser, constants.DevtoolPort, constants.DevtoolDB, constants.DevtoolDBUser},
		&obj.Meter.ServiceDBCommonOptions:        {constants.MeterAdminUser, constants.MeterPort, constants.MeterDB, constants.MeterDBUser},
		&obj.Monitor:                             {constants.MonitorAdminUser, constants.MonitorPort, constants.MonitorDB, constants.MonitorDBUser},
		&obj.Itsm.ServiceDBCommonOptions:         {constants.ItsmAdminUser, constants.ItsmPort, constants.ItsmDB, constants.ItsmDBUser},
		&obj.CloudId:                             {constants.CloudIdAdminUser, constants.CloudIdPort, constants.CloudIdDB, constants.CloudIdDBUser},
		&obj.Webconsole:                          {constants.WebconsoleAdminUser, constants.WebconsolePort, constants.WebconsoleDB, constants.WebconsoleDBUser},
		&obj.Report:                              {constants.ReportAdminUser, constants.ReportPort, constants.ReportDB, constants.ReportDBUser},
		&obj.MspOperation:                        {constants.MspOperationAdminUser, constants.MspOperationPort, constants.MspOperationDB, constants.MspOperationDBUser},
	} {
		if user, ok := registryPorts[tmp.port]; ok {
			log.Fatalf("port %d has been registered by %s", tmp.port, user)
		}
		registryPorts[tmp.port] = tmp.user
		SetDefaults_ServiceDBCommonOptions(opt, tmp.db, tmp.dbUser, tmp.user, tmp.port)
	}
	setDefaults_ItsmConfig(&obj.Itsm)

	setDefaults_DBConfig(&obj.Grafana.DB, constants.MonitorStackGrafanaDB, constants.MonitorStackGrafanaDBUer)
}

func SetDefaults_ServiceBaseConfig(obj *ServiceBaseConfig, port int) {
	if obj.Port == 0 {
		obj.Port = port
	}
}

func setDefaults_KeystoneConfig(obj *KeystoneConfig) {
	SetDefaults_ServiceBaseConfig(&obj.ServiceBaseConfig, constants.KeystonePublicPort)
	setDefaults_DBConfig(&obj.DB, constants.KeystoneDB, constants.KeystoneDBUser)
}

func SetDefaults_ServiceCommonOptions(obj *ServiceCommonOptions, user string, port int) {
	SetDefaults_ServiceBaseConfig(&obj.ServiceBaseConfig, port)
	setDefaults_CloudUser(&obj.CloudUser, user)
}

func SetDefaults_ServiceDBCommonOptions(obj *ServiceDBCommonOptions, db, dbUser string, svcUser string, port int) {
	setDefaults_DBConfig(&obj.DB, db, dbUser)
	SetDefaults_ServiceCommonOptions(&obj.ServiceCommonOptions, svcUser, port)
}

func setDefaults_DBConfig(obj *DBConfig, database string, username string) {
	if obj.Database == "" {
		obj.Database = database
	}
	if obj.Username == "" {
		obj.Username = username
	}
	if obj.Password == "" {
		obj.Password = passwd.GeneratePassword()
	}
}

func setDefaults_CloudUser(obj *CloudUser, username string) {
	if obj.Username == "" {
		obj.Username = username
	}
	if obj.Password == "" {
		obj.Password = passwd.GeneratePassword()
	}
}

func setDefaults_ItsmConfig(obj *ItsmConfig) {
	obj.SecondDatabase = fmt.Sprintf("%s_engine", obj.DB.Database)
	obj.EncryptionKey = passwd.GeneratePassword()
}
