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

package compute

import "yunion.io/x/onecloud/pkg/apis"

type VpcDetails struct {
	apis.EnabledStatusInfrasResourceBaseDetails
	ManagedResourceInfo
	CloudregionResourceInfo
	GlobalVpcResourceInfo

	SVpc

	// 二层网络数量
	// example: 1
	WireCount int `json:"wire_count"`
	// IP子网个数
	// example: 2
	NetworkCount int `json:"network_count"`
	// 路由表个数
	// example: 0
	RoutetableCount int `json:"routetable_count"`
	// NAT网关个数
	// example: 0
	NatgatewayCount int `json:"natgateway_count"`
}

type VpcResourceInfoBase struct {
	// Vpc Name
	Vpc string `json:"vpc"`

	// VPC外部Id
	VpcExtId string `json:"vpc_ext_id"`
}

type VpcResourceInfo struct {
	VpcResourceInfoBase

	// VPC归属区域ID
	CloudregionId string `json:"cloudregion_id"`

	CloudregionResourceInfo

	// VPC归属云订阅ID
	ManagerId string `json:"manager_id"`

	ManagedResourceInfo
}

type VpcSyncstatusInput struct {
}

type VpcCreateInput struct {
	apis.EnabledStatusInfrasResourceBaseCreateInput

	CloudregionResourceInput

	CloudproviderResourceInput

	// CIDR_BLOCK
	CidrBlock string `json:"cidr_block"`
}

type VpcResourceInput struct {
	// 关联VPC(ID或Name)
	Vpc string `json:"vpc"`
	// swagger:ignore
	// Deprecated
	// filter by vpc Id
	VpcId string `json:"vpc_id" deprecated-by:"vpc"`
}

type VpcFilterListInputBase struct {
	VpcResourceInput

	// 按VPC名称排序
	// pattern:asc|desc
	OrderByVpc string `json:"order_by_vpc"`
}

type VpcFilterListInput struct {
	VpcFilterListInputBase
	RegionalFilterListInput
	ManagedResourceListInput
}
