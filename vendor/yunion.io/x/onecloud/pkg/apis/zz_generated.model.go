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

// Code generated by model-api-gen. DO NOT EDIT.

package apis

import (
	time "time"
)

// SAdminSharableVirtualResourceBase is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudcommon/db.SAdminSharableVirtualResourceBase.
type SAdminSharableVirtualResourceBase struct {
	SSharableVirtualResourceBase
	Records string `json:"records"`
}

// SDomainLevelResourceBase is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudcommon/db.SDomainLevelResourceBase.
type SDomainLevelResourceBase struct {
	SStandaloneResourceBase
	SDomainizedResourceBase
	// 归属Domain信息的来源, local: 本地设置, cloud: 从云上同步过来
	// example: local
	DomainSrc string `json:"domain_src"`
}

// SDomainizedResourceBase is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudcommon/db.SDomainizedResourceBase.
type SDomainizedResourceBase struct {
	// 域Id
	DomainId string `json:"domain_id"`
}

// SEnabledResourceBase is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudcommon/db.SEnabledResourceBase.
type SEnabledResourceBase struct {
	// 资源是否启用
	Enabled *bool `json:"enabled,omitempty"`
}

// SEnabledStatusDomainLevelResourceBase is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudcommon/db.SEnabledStatusDomainLevelResourceBase.
type SEnabledStatusDomainLevelResourceBase struct {
	SStatusDomainLevelResourceBase
	SEnabledResourceBase
}

// SEnabledStatusInfrasResourceBase is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudcommon/db.SEnabledStatusInfrasResourceBase.
type SEnabledStatusInfrasResourceBase struct {
	SStatusInfrasResourceBase
	SEnabledResourceBase
}

// SEnabledStatusStandaloneResourceBase is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudcommon/db.SEnabledStatusStandaloneResourceBase.
type SEnabledStatusStandaloneResourceBase struct {
	SStatusStandaloneResourceBase
	SEnabledResourceBase
}

// SExternalizedResourceBase is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudcommon/db.SExternalizedResourceBase.
type SExternalizedResourceBase struct {
	// 外部Id, 对用公有云私有资源自身的Id
	ExternalId string `json:"external_id"`
}

// SInfrasResourceBase is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudcommon/db.SInfrasResourceBase.
type SInfrasResourceBase struct {
	SDomainLevelResourceBase
	SSharableBaseResource
}

// SJointResourceBase is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudcommon/db.SJointResourceBase.
type SJointResourceBase struct {
	SResourceBase
	RowId int64 `json:"row_id"`
}

// SKeystoneCacheObject is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudcommon/db.SKeystoneCacheObject.
type SKeystoneCacheObject struct {
	SStandaloneResourceBase
	DomainId  string    `json:"domain_id"`
	Domain    string    `json:"domain"`
	LastCheck time.Time `json:"last_check"`
}

// SProjectizedResourceBase is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudcommon/db.SProjectizedResourceBase.
type SProjectizedResourceBase struct {
	SDomainizedResourceBase
	// 项目Id
	ProjectId string `json:"tenant_id"`
}

// SResourceBase is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudcommon/db.SResourceBase.
type SResourceBase struct {
	// 资源创建时间
	CreatedAt time.Time `json:"created_at"`
	// 资源更新时间
	UpdatedAt time.Time `json:"updated_at"`
	// 资源被更新次数
	UpdateVersion int `json:"update_version"`
	// 资源删除时间
	DeletedAt time.Time `json:"deleted_at"`
	// 资源是否被删除
	Deleted bool `json:"deleted"`
}

// SScopedResourceBase is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudcommon/db.SScopedResourceBase.
type SScopedResourceBase struct {
	SProjectizedResourceBase
}

// SSharableBaseResource is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudcommon/db.SSharableBaseResource.
type SSharableBaseResource struct {
	// 是否共享
	IsPublic bool `json:"is_public"`
	// 默认共享范围
	PublicScope string `json:"public_scope"`
	// 共享设置的来源, local: 本地设置, cloud: 从云上同步过来
	// example: local
	PublicSrc string `json:"public_src"`
}

// SSharableVirtualResourceBase is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudcommon/db.SSharableVirtualResourceBase.
type SSharableVirtualResourceBase struct {
	SVirtualResourceBase
	SSharableBaseResource
}

// SSharedResource is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudcommon/db.SSharedResource.
type SSharedResource struct {
	SResourceBase
	Id           int64  `json:"id"`
	ResourceType string `json:"resource_type"`
	ResourceId   string `json:"resource_id"`
	// OwnerProjectId  string `width:"128" charset:"ascii" nullable:"false" index:"true" json:"owner_project_id"`
	TargetProjectId string `json:"target_project_id"`
	TargetType      string `json:"target_type"`
}

// SStandaloneResourceBase is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudcommon/db.SStandaloneResourceBase.
type SStandaloneResourceBase struct {
	SResourceBase
	// 资源UUID
	Id string `json:"id"`
	// 资源名称
	Name string `json:"name"`
	// 资源描述信息
	Description string `json:"description"`
	// 是否是模拟资源, 部分从公有云上同步的资源并不真实存在, 例如宿主机
	// list 接口默认不会返回这类资源，除非显示指定 is_emulate=true 过滤参数
	IsEmulated bool `json:"is_emulated"`
}

// SStatusDomainLevelResourceBase is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudcommon/db.SStatusDomainLevelResourceBase.
type SStatusDomainLevelResourceBase struct {
	SDomainLevelResourceBase
	SStatusResourceBase
}

// SStatusInfrasResourceBase is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudcommon/db.SStatusInfrasResourceBase.
type SStatusInfrasResourceBase struct {
	SInfrasResourceBase
	SStatusResourceBase
}

// SStatusResourceBase is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudcommon/db.SStatusResourceBase.
type SStatusResourceBase struct {
	// 资源状态
	Status string `json:"status"`
}

// SStatusStandaloneResourceBase is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudcommon/db.SStatusStandaloneResourceBase.
type SStatusStandaloneResourceBase struct {
	SStandaloneResourceBase
	SStatusResourceBase
}

// SVirtualJointResourceBase is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudcommon/db.SVirtualJointResourceBase.
type SVirtualJointResourceBase struct {
	SJointResourceBase
}

// SVirtualResourceBase is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudcommon/db.SVirtualResourceBase.
type SVirtualResourceBase struct {
	SStatusStandaloneResourceBase
	SProjectizedResourceBase
	// 云上同步资源是否在本地被更改过配置, local: 更改过, cloud: 未更改过
	// example: local
	ProjectSrc string `json:"project_src"`
	// 是否是系统资源
	IsSystem bool `json:"is_system"`
	// 资源放入回收站时间
	PendingDeletedAt time.Time `json:"pending_deleted_at"`
	// 资源是否处于回收站中
	PendingDeleted bool `json:"pending_deleted"`
}
