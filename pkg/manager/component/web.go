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

	apps "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/intstr"

	"yunion.io/x/onecloud-operator/pkg/apis/constants"
	"yunion.io/x/onecloud-operator/pkg/apis/onecloud/v1alpha1"
	"yunion.io/x/onecloud-operator/pkg/controller"
	"yunion.io/x/onecloud-operator/pkg/manager"
)

const (
	CEConfig = `
    location / {
        root /usr/share/nginx/html/web;
        index index.html;
        add_header Cache-Control no-cache;
        expires 1s;
        if (!-e $request_filename) {
            rewrite ^/(.*) /index.html last;
            break;
        }
    }
`

	EEConfig = `
    location ~ ^/v[12]/ { rewrite ^/v.(.*)$ $1 redirect; }

    location / {
        root   /usr/share/nginx/html/dashboard;
        index index.html;
        try_files $uri $uri/ /index.html;
    }

    location /overview {
        proxy_pass http://localhost:8080;
        proxy_redirect   off;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto https;
   }

    location /docs {
        proxy_pass http://localhost:8081;
        proxy_redirect   off;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto https;
   }
`

	WebNginxConfigTemplate = `
map $http_upgrade $connection_upgrade {
    default upgrade;
    '' close;
}

server {
    listen 443 default_server ssl;
    server_name _;
    ssl_certificate /etc/yunion/pki/service.crt;
    ssl_certificate_key /etc/yunion/pki/service.key;

    gzip_static on;
    gzip on;
    gzip_proxied any;
    gzip_min_length  1k;
    gzip_buffers     4 16k;
    gzip_http_version 1.0;
    gzip_comp_level 5;
    gzip_types text/plain application/javascript application/css text/css application/xml text/javascript application/x-httpd-php image/jpeg image/gif image/png;
    gzip_vary on;
    chunked_transfer_encoding off;

    client_body_buffer_size 16k;
    client_header_buffer_size 16k;
    client_max_body_size 8m;
    large_client_header_buffers 2 16k;
    client_body_timeout 20s;
    client_header_timeout 20s;

{{.EditionConfig}}

    location /static/ {
        # Some basic cache-control for static files to be sent to the browser
        root /usr/share/nginx/html/web;
        expires max;
        add_header Pragma public;
        add_header Cache-Control "public, must-revalidate, proxy-revalidate";
    }

    location /itsm {
        alias /usr/share/nginx/html/itsm;
        index index.html;
        add_header Cache-Control no-cache;
        expires 1s;
        if (!-e $request_filename) {
            rewrite ^/(.*) /itsm/index.html last;
            break;
        }
    }

    location /static-itsm/ {
        root /usr/share/nginx/html/itsm;
        expires max;
        add_header Pragma public;
        add_header Cache-Control "public, must-revalidate, proxy-revalidate";
    }

    location ~* /static/images/favicon.* {
        root /usr/share/nginx/html/login;
        expires max;
        add_header Pragma public;
        add_header Cache-Control "public, must-revalidate, proxy-revalidate";
    }

    location /api {
        proxy_pass {{.APIGatewayURL}};
        proxy_redirect   off;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;

        proxy_buffer_size  16k;
        proxy_buffers   32 16k;
        proxy_busy_buffers_size 16k;
        proxy_temp_file_write_size 16k;
        proxy_read_timeout 600;
    }

    location /api/v1/imageutils/upload {
        proxy_pass {{.APIGatewayURL}};
        client_max_body_size 0;
        proxy_http_version 1.1;
        proxy_request_buffering off;
        proxy_buffering off;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $remote_addr;
    }

    location /api/v1/s3uploads {
        proxy_pass {{.APIGatewayURL}};
        client_max_body_size 0;
        proxy_http_version 1.1;
        proxy_request_buffering off;
        proxy_buffering off;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $remote_addr;
    }

    location /query {
        proxy_pass {{.APIGatewayURL}};
        proxy_redirect   off;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    location ~ ^/(vnc|spice|wmks|sol) {
        proxy_pass {{.WebconsoleURL}};
        proxy_redirect   off;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    location ~ ^/(websockify|wsproxy|connect) {
        proxy_pass {{.WebconsoleURL}};
        proxy_redirect   off;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;

        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";

        proxy_read_timeout 86400;
    }

    location /web-console {
        alias /usr/share/nginx/html/web-console;
        index index.html;
        if (!-e $request_filename) {
            rewrite ^/(.*) /web-console/index.html last;
            break;
        }
    }

    location /ws {
        proxy_pass {{.APIGatewayWsURL}};
        proxy_redirect   off;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;

        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";

        proxy_read_timeout 86400;
    }

    location /socket.io/ {
        proxy_pass {{.APIGatewayWsURL}};
        proxy_redirect   off;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;

        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";

        proxy_read_timeout 86400;
    }

    location /baremetal-prepare/ {
        # Some basic cache-control for static files to be sent to the browser
        root /opt/cloud/yunion/baremetal/;
        expires max;
        add_header Pragma public;
        add_header Cache-Control "public, must-revalidate, proxy-revalidate";
    }
}
`
)

type WebNginxConfig struct {
	EditionConfig   string
	WebconsoleURL   string
	APIGatewayWsURL string
	APIGatewayURL   string
}

func (c WebNginxConfig) GetContent() (string, error) {
	return CompileTemplateFromMap(WebNginxConfigTemplate, c)
}

type webManager struct {
	*ComponentManager
}

func newWebManager(man *ComponentManager) manager.Manager {
	m := &webManager{man}
	return m
}

func (m *webManager) getComponentType() v1alpha1.ComponentType {
	return v1alpha1.WebComponentType
}

func (m *webManager) getProductVersions() []v1alpha1.ProductVersion {
	return []v1alpha1.ProductVersion{
		v1alpha1.ProductVersionFullStack,
		v1alpha1.ProductVersionCMP,
		v1alpha1.ProductVersionEdge,
	}
}

func (m *webManager) Sync(oc *v1alpha1.OnecloudCluster) error {
	if IsEnterpriseEdition(oc) {
		oc.Spec.Web.ImageName = constants.WebEEImageName
	} else {
		oc.Spec.Web.ImageName = constants.WebCEImageName
	}
	return syncComponent(m, oc, oc.Spec.Web.Disable, "")
}

func (m *webManager) getService(oc *v1alpha1.OnecloudCluster, cfg *v1alpha1.OnecloudClusterConfig, zone string) []*corev1.Service {
	ports := []corev1.ServicePort{
		{
			Name:       "https",
			Protocol:   corev1.ProtocolTCP,
			Port:       443,
			TargetPort: intstr.FromInt(443),
		},
		{
			Name:       "overview",
			Protocol:   corev1.ProtocolTCP,
			Port:       8080,
			TargetPort: intstr.FromInt(8080),
		},
		{
			Name:       "docs",
			Protocol:   corev1.ProtocolTCP,
			Port:       8081,
			TargetPort: intstr.FromInt(8081),
		},
	}
	return []*corev1.Service{m.newService(v1alpha1.WebComponentType, oc, corev1.ServiceTypeClusterIP, ports)}
}

func (m *webManager) getIngress(oc *v1alpha1.OnecloudCluster, zone string) *unstructured.Unstructured {
	ocName := oc.GetName()
	svcName := controller.NewClusterComponentName(ocName, v1alpha1.WebComponentType)
	appLabel := m.getComponentLabel(oc, v1alpha1.WebComponentType)
	secretName := controller.ClustercertSecretName(oc)

	obj := new(unstructured.Unstructured)
	objMeta := m.getObjectMeta(oc, svcName, appLabel)
	obj.SetAnnotations(objMeta.Annotations)
	obj.SetName(objMeta.Name)
	obj.SetNamespace(objMeta.Namespace)
	obj.SetLabels(objMeta.Labels)
	obj.SetOwnerReferences(objMeta.OwnerReferences)

	// unstructured.SetNestedSlice(obj.Object, []interface{}{
	// 	map[string]interface{}{
	// 		"secretName": secretName,
	// 	},
	// }, "spec", "tls")
	// unstructured.SetNestedSlice(obj.Object, []interface{}{
	// 	map[string]interface{}{
	// 		"http": map[string]interface{}{
	// 			"paths": []interface{}{
	// 				map[string]interface{}{
	// 					"path": "/",
	// 					"backend": map[string]interface{}{
	// 						"serviceName": svcName,
	// 						"servicePort": int64(443),
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// }, "spec", "rules")

	unstructured.SetNestedMap(obj.Object, map[string]interface{}{
		"tls": []interface{}{
			map[string]interface{}{
				"secretName": secretName,
			},
		},
		"rules": []interface{}{
			map[string]interface{}{
				"http": map[string]interface{}{
					"paths": []interface{}{
						map[string]interface{}{
							"pathType": "Exact",
							"path":     "/",
							"backend": map[string]interface{}{
								"serviceName": svcName,
								"servicePort": int64(443),
								"service": map[string]interface{}{
									"name": svcName,
									"port": map[string]interface{}{
										"number": int64(443),
									},
								},
							},
						},
					},
				},
			},
		},
	}, "spec")

	// ing := &extensions.Ingress{
	// 	ObjectMeta: m.getObjectMeta(oc, svcName, appLabel),
	// 	Spec: extensions.IngressSpec{
	// 		TLS: []extensions.IngressTLS{
	// 			{
	// 				SecretName: secretName,
	// 			},
	// 		},
	// 		Rules: []extensions.IngressRule{
	// 			{
	// 				IngressRuleValue: extensions.IngressRuleValue{
	// 					HTTP: &extensions.HTTPIngressRuleValue{
	// 						Paths: []extensions.HTTPIngressPath{
	// 							{
	// 								Path: "/",
	// 								Backend: extensions.IngressBackend{
	// 									ServiceName: svcName,
	// 									ServicePort: intstr.FromInt(443),
	// 								},
	// 							},
	// 						},
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// }

	// for nginx ingress
	anno := obj.GetAnnotations()
	if len(anno) == 0 {
		anno = map[string]string{}
	}
	anno["nginx.ingress.kubernetes.io/backend-protocol"] = "HTTPS"
	obj.SetAnnotations(anno)

	return m.addIngressPaths(IsEnterpriseEdition(oc), svcName, obj)
}

func (m *webManager) addIngressPaths(isEE bool, svcName string, ing *unstructured.Unstructured) *unstructured.Unstructured {
	if !isEE {
		return ing
	}
	/*
	 * rule := &ing.Spec.Rules[0]
	 * if !IsPathIngressRule("/overview", rule.HTTP.Paths) {
	 * 	rule.HTTP.Paths = append(rule.HTTP.Paths,
	 * 		extensions.HTTPIngressPath{
	 * 			Path: "/overview",
	 * 			Backend: extensions.IngressBackend{
	 * 				ServiceName: svcName,
	 * 				ServicePort: intstr.FromInt(8080),
	 * 			},
	 * 		},
	 * 	)
	 * }
	 * if !IsPathIngressRule("/docs", rule.HTTP.Paths) {
	 * 	rule.HTTP.Paths = append(rule.HTTP.Paths,
	 * 		extensions.HTTPIngressPath{
	 * 			Path: "/docs",
	 * 			Backend: extensions.IngressBackend{
	 * 				ServiceName: svcName,
	 * 				ServicePort: intstr.FromInt(8081),
	 * 			},
	 * 		},
	 * 	)
	 * }
	 */
	return ing
}

func IsPathIngressRule(path string, paths []interface{}) bool {
	exist := false
	for _, obj := range paths {
		ip := obj.(map[string]interface{})
		if ip["path"].(string) == path {
			exist = true
			break
		}
	}
	return exist
}

func (m *webManager) updateIngress(oc *v1alpha1.OnecloudCluster, oldIng *unstructured.Unstructured) *unstructured.Unstructured {
	newIng := oldIng.DeepCopy()
	spec, _, _ := unstructured.NestedMap(newIng.Object, "spec")
	doUpdate := false
	rules, _, _ := unstructured.NestedSlice(spec, "rules")
	for _, rule := range rules {
		http, _, _ := unstructured.NestedMap(rule.(map[string]interface{}), "http")
		if http == nil {
			continue
		}
		paths, _, _ := unstructured.NestedSlice(http, "paths")
		for _, path := range []string{"/overview", "/docs"} {
			if !IsPathIngressRule(path, paths) {
				doUpdate = true
				break
			}
		}
	}
	if doUpdate {
		svcName := m.getService(oc, nil, "")[0].GetName()
		newIng = m.addIngressPaths(IsEnterpriseEdition(oc), svcName, newIng)
	}
	return newIng
}

func (m *webManager) getConfigMap(oc *v1alpha1.OnecloudCluster, cfg *v1alpha1.OnecloudClusterConfig, zone string) (*corev1.ConfigMap, bool, error) {
	urlF := func(ct v1alpha1.ComponentType, port int) string {
		return fmt.Sprintf("https://%s:%d", controller.NewClusterComponentName(oc.GetName(), ct), port)
	}
	conf := CEConfig
	isEE := IsEnterpriseEdition(oc)
	if isEE {
		conf = EEConfig
	}
	config := WebNginxConfig{
		EditionConfig:   conf,
		WebconsoleURL:   urlF(v1alpha1.WebconsoleComponentType, constants.WebconsolePort),
		APIGatewayWsURL: urlF(v1alpha1.APIGatewayComponentType, constants.APIWebsocketPort),
		APIGatewayURL:   urlF(v1alpha1.APIGatewayComponentType, constants.APIGatewayPort),
	}
	content, err := config.GetContent()
	if err != nil {
		return nil, false, err
	}
	return m.newConfigMap(v1alpha1.WebComponentType, "", oc, content), false, nil
}

func (m *webManager) getDeployment(oc *v1alpha1.OnecloudCluster, cfg *v1alpha1.OnecloudClusterConfig, zone string) (*apps.Deployment, error) {
	cf := func(volMounts []corev1.VolumeMount) []corev1.Container {
		confVol := volMounts[len(volMounts)-1]
		confVol.MountPath = "/etc/nginx/conf.d"
		volMounts[len(volMounts)-1] = confVol
		containers := []corev1.Container{
			{
				Name:            "web",
				Image:           oc.Spec.Web.Image,
				ImagePullPolicy: oc.Spec.Web.ImagePullPolicy,
				Command: []string{
					"nginx",
					"-g",
					"daemon off;",
				},
				Ports: []corev1.ContainerPort{
					{
						Name:          "web",
						ContainerPort: 80,
						Protocol:      corev1.ProtocolTCP,
					},
				},
				VolumeMounts: volMounts,
			},
		}
		if IsEnterpriseEdition(oc) {
			containers = append(containers,
				corev1.Container{
					Name:            "overview",
					Image:           oc.Spec.Web.Overview.Image,
					ImagePullPolicy: oc.Spec.Web.Overview.ImagePullPolicy,
					Ports: []corev1.ContainerPort{
						{
							Name:          "overview",
							ContainerPort: 8080,
							Protocol:      corev1.ProtocolTCP,
						},
					},
					VolumeMounts: volMounts,
				},
				corev1.Container{
					Name:            "docs",
					Image:           oc.Spec.Web.Docs.Image,
					ImagePullPolicy: oc.Spec.Web.Docs.ImagePullPolicy,
					Ports: []corev1.ContainerPort{
						{
							Name:          "docs",
							ContainerPort: 8081,
							Protocol:      corev1.ProtocolTCP,
						},
					},
					VolumeMounts: volMounts,
				},
			)
		}
		return containers
	}
	deploy, err := m.newDefaultDeploymentNoInit(v1alpha1.WebComponentType, "", oc, NewVolumeHelper(oc, controller.ComponentConfigMapName(oc, v1alpha1.WebComponentType), v1alpha1.WebComponentType), &oc.Spec.Web.DeploymentSpec, cf)
	if err != nil {
		return nil, err
	}
	podSpec := &deploy.Spec.Template.Spec
	config := podSpec.Volumes[len(podSpec.Volumes)-1]
	config.ConfigMap.Items[0].Path = "default.conf"
	podSpec.Volumes[len(podSpec.Volumes)-1] = config
	return deploy, nil
}

func (m *webManager) getDeploymentStatus(oc *v1alpha1.OnecloudCluster, zone string) *v1alpha1.DeploymentStatus {
	return &oc.Status.Web
}
