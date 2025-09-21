/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package migrate

import (
	"context"

	ocpappsv1 "github.com/openshift/api/apps/v1"
	routev1 "github.com/openshift/api/route/v1"
	corev1 "k8s.io/api/core/v1"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Convert DeploymentConfigs → Deployments
func (r *MigrationPlanReconciler) convertDeploymentConfigs(ctx context.Context, ns string) ([]string, int, []string) {
	var dcs ocpappsv1.DeploymentConfigList
	var yamlDocs []string
	var notes []string
	if err := r.List(ctx, &dcs, client.InNamespace(ns)); err != nil {
		return nil, 0, []string{"error listing DeploymentConfigs: " + err.Error()}
	}
	for _, dc := range dcs.Items {
		yamlDocs = append(yamlDocs,
			"---",
			"apiVersion: apps/v1",
			"kind: Deployment",
			"metadata:",
			"  name: "+dc.Name,
			"  namespace: "+dc.Namespace,
			"  # TODO: copy labels/selector/template",
		)
		notes = append(notes, "Converted DeploymentConfig "+dc.Name+" → Deployment")
	}
	return yamlDocs, len(dcs.Items), notes
}

// Convert Routes → Ingresses
func (r *MigrationPlanReconciler) convertRoutes(ctx context.Context, ns string, ingressClass string) ([]string, int, []string) {
	var routes routev1.RouteList
	var yamlDocs []string
	var notes []string
	if err := r.List(ctx, &routes, client.InNamespace(ns)); err != nil {
		return nil, 0, []string{"error listing Routes: " + err.Error()}
	}
	for _, rt := range routes.Items {
		class := ingressClass
		if class == "" {
			class = "nginx"
		}
		yamlDocs = append(yamlDocs,
			"---",
			"apiVersion: networking.k8s.io/v1",
			"kind: Ingress",
			"metadata:",
			"  name: "+rt.Name,
			"  namespace: "+rt.Namespace,
			"  annotations:",
			"    kubernetes.io/ingress.class: "+class,
			"  # TODO: map host/path/service",
		)
		notes = append(notes, "Converted Route "+rt.Name+" → Ingress")
	}
	return yamlDocs, len(routes.Items), notes
}

// Convert Services
func (r *MigrationPlanReconciler) convertServices(ctx context.Context, ns string) ([]string, int, []string) {
	var svcs corev1.ServiceList
	var yamlDocs []string
	var notes []string
	if err := r.List(ctx, &svcs, client.InNamespace(ns)); err != nil {
		return nil, 0, []string{"error listing Services: " + err.Error()}
	}
	for _, svc := range svcs.Items {
		yamlDocs = append(yamlDocs,
			"---",
			"apiVersion: v1",
			"kind: Service",
			"metadata:",
			"  name: "+svc.Name,
			"  namespace: "+svc.Namespace,
			"  # TODO: copy spec.ports and selectors",
		)
		notes = append(notes, "Processed Service "+svc.Name)
	}
	return yamlDocs, len(svcs.Items), notes
}

// Convert PVCs
func (r *MigrationPlanReconciler) convertPVCs(ctx context.Context, ns string, targetCloud string) ([]string, int, []string) {
	var pvcs corev1.PersistentVolumeClaimList
	var yamlDocs []string
	var notes []string
	if err := r.List(ctx, &pvcs, client.InNamespace(ns)); err != nil {
		return nil, 0, []string{"error listing PVCs: " + err.Error()}
	}
	for _, pvc := range pvcs.Items {
		yamlDocs = append(yamlDocs,
			"---",
			"apiVersion: v1",
			"kind: PersistentVolumeClaim",
			"metadata:",
			"  name: "+pvc.Name,
			"  namespace: "+pvc.Namespace,
			"  # TODO: map storageClassName for "+targetCloud,
		)
		notes = append(notes, "Processed PVC "+pvc.Name+" (map to storage class for "+targetCloud+")")
	}
	return yamlDocs, len(pvcs.Items), notes
}
