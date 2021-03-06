/*
 * Copyright (c) 2020 WSO2 Inc. (http:www.wso2.org) All Rights Reserved.
 *
 * WSO2 Inc. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http:www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package integration

import (
	wso2v1alpha1 "github.com/wso2/k8s-api-operator/api-operator/pkg/apis/wso2/v1alpha1"
	v1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"reflect"
	"strconv"
)

// labelsForIntegration returns the labels for selecting the resources
// belonging to the given integration CR name.
func labelsForIntegration(name string) map[string]string {
	return map[string]string{"app": "integration", "integration_cr": name}
}

// nameForDeployment gives the name for the deployment
func nameForDeployment(m *wso2v1alpha1.Integration) string {
	return m.Name + "-deployment"
}

// nameForService gives the name for the service
func nameForService(m *wso2v1alpha1.Integration) string {
	return m.Name + "-service"
}

// nameForInboundService gives the name for the inbound service
func nameForInboundService(m *wso2v1alpha1.Integration) string {
	return m.Name + "-inbound"
}

// nameForIngress gives the name for the ingress
func nameForIngress() string {
	return "ei-operator-ingress"
}

// nameForConfigMap gives the name for the config map
func nameForConfigMap() string {
	return "ei-operator-config"
}

// CheckIngressRulesExist checks the ingress rules are exist in current ingress
func CheckIngressRulesExist(m *wso2v1alpha1.Integration, eic *EIController, currentIngress *v1beta1.Ingress) ([]v1beta1.IngressRule, bool) {

	ingressPaths := GenerateIngressPaths(m)

	currentRules := currentIngress.Spec.Rules
	newRule := v1beta1.IngressRule{
		Host: eic.Host,
		IngressRuleValue: v1beta1.IngressRuleValue{
			HTTP: &v1beta1.HTTPIngressRuleValue{
				Paths: ingressPaths,
			},
		},
	}

	// check the rules are exists in the ingress, if not add the rules
	// checking because of reconsile is looping
	ruleExists := false
	for _, rule := range currentRules {
		if reflect.DeepEqual(rule, newRule) {
			ruleExists = true
		}
	}

	if !ruleExists {
		currentRules = append(currentRules, newRule)
	}
	return currentRules, ruleExists
}

// GenerateIngressPaths generates the ingress paths
func GenerateIngressPaths(m *wso2v1alpha1.Integration) []v1beta1.HTTPIngressPath {
	var ingressPaths []v1beta1.HTTPIngressPath

	//Set HTTP ingress path
	httpPath := "/" + nameForService(m) + "(/|$)(.*)"
	httpIngressPath := v1beta1.HTTPIngressPath{
		Path: httpPath,
		Backend: v1beta1.IngressBackend{
			ServiceName: nameForService(m),
			ServicePort: intstr.IntOrString{
				Type:   Int,
				IntVal: 8290,
			},
		},
	}
	ingressPaths = append(ingressPaths, httpIngressPath)

	// check inbound endpoint port is exist and update the ingress path
	for _, port := range m.Spec.InboundPorts {
		inboundPath := "/" + nameForInboundService(m) +
			"/" + strconv.Itoa(int(port)) + "(/|$)(.*)"
		inboundIngressPath := v1beta1.HTTPIngressPath{
			Path: inboundPath,
			Backend: v1beta1.IngressBackend{
				ServiceName: nameForService(m),
				ServicePort: intstr.IntOrString{
					Type:   Int,
					IntVal: port,
				},
			},
		}
		ingressPaths = append(ingressPaths, inboundIngressPath)
	}

	return ingressPaths
}
