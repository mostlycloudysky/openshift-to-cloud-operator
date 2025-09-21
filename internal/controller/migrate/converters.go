package migrate

import (
	"context"

	ocpappsv1 "github.com/openshift/api/apps/v1"
	routev1 "github.com/openshift/api/route/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"
)

func int32ptr(i int32) *int32 { return &i }

func toLabelSelector(m map[string]string) *metav1.LabelSelector {
	if m == nil {
		return &metav1.LabelSelector{}
	}
	return &metav1.LabelSelector{MatchLabels: m}
}

// Convert DeploymentConfigs → Deployments
func (r *MigrationPlanReconciler) convertDeploymentConfigs(ctx context.Context, ns string) ([]string, int, []string) {
	var dcs ocpappsv1.DeploymentConfigList
	var yamlDocs []string
	var notes []string

	if err := r.List(ctx, &dcs, client.InNamespace(ns)); err != nil {
		return nil, 0, []string{"error listing DeploymentConfigs: " + err.Error()}
	}

	for _, dc := range dcs.Items {
		deploy := appsv1.Deployment{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "apps/v1",
				Kind:       "Deployment",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      dc.Name,
				Namespace: dc.Namespace,
				Labels:    dc.Labels,
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: int32ptr(dc.Spec.Replicas),        // fix replicas pointer
				Selector: toLabelSelector(dc.Spec.Selector), // convert map[string]string → LabelSelector
				Template: *dc.Spec.Template,                 // this part is already compatible
			},
		}

		y, _ := yaml.Marshal(deploy)
		yamlDocs = append(yamlDocs, "---\n"+string(y))
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

	if ingressClass == "" {
		ingressClass = "nginx"
	}

	for _, rt := range routes.Items {
		ing := networkingv1.Ingress{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "networking.k8s.io/v1",
				Kind:       "Ingress",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      rt.Name,
				Namespace: rt.Namespace,
				Annotations: map[string]string{
					"kubernetes.io/ingress.class": ingressClass,
				},
			},
			// Minimal spec (host/path/service mapping could be improved later)
			Spec: networkingv1.IngressSpec{
				Rules: []networkingv1.IngressRule{
					{
						Host: rt.Spec.Host,
						IngressRuleValue: networkingv1.IngressRuleValue{
							HTTP: &networkingv1.HTTPIngressRuleValue{
								Paths: []networkingv1.HTTPIngressPath{
									{
										Path:     "/",
										PathType: func() *networkingv1.PathType { pt := networkingv1.PathTypePrefix; return &pt }(),
										Backend: networkingv1.IngressBackend{
											Service: &networkingv1.IngressServiceBackend{
												Name: rt.Spec.To.Name,
												Port: networkingv1.ServiceBackendPort{
													Number: rt.Spec.Port.TargetPort.IntVal,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}

		y, _ := yaml.Marshal(ing)
		yamlDocs = append(yamlDocs, "---\n"+string(y))
		notes = append(notes, "Converted Route "+rt.Name+" → Ingress")
	}
	return yamlDocs, len(routes.Items), notes
}

// Convert Services → Services (copy spec)
func (r *MigrationPlanReconciler) convertServices(ctx context.Context, ns string) ([]string, int, []string) {
	var svcs corev1.ServiceList
	var yamlDocs []string
	var notes []string

	if err := r.List(ctx, &svcs, client.InNamespace(ns)); err != nil {
		return nil, 0, []string{"error listing Services: " + err.Error()}
	}

	for _, svc := range svcs.Items {
		outSvc := corev1.Service{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "v1",
				Kind:       "Service",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      svc.Name,
				Namespace: svc.Namespace,
				Labels:    svc.Labels,
			},
			Spec: svc.Spec,
		}

		y, _ := yaml.Marshal(outSvc)
		yamlDocs = append(yamlDocs, "---\n"+string(y))
		notes = append(notes, "Processed Service "+svc.Name)
	}
	return yamlDocs, len(svcs.Items), notes
}

// Convert PVCs → PVCs (preserve spec + remap storageClassName)
func (r *MigrationPlanReconciler) convertPVCs(ctx context.Context, ns string, targetCloud string) ([]string, int, []string) {
	var pvcs corev1.PersistentVolumeClaimList
	var yamlDocs []string
	var notes []string

	if err := r.List(ctx, &pvcs, client.InNamespace(ns)); err != nil {
		return nil, 0, []string{"error listing PVCs: " + err.Error()}
	}

	for _, pvc := range pvcs.Items {
		outPVC := corev1.PersistentVolumeClaim{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "v1",
				Kind:       "PersistentVolumeClaim",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      pvc.Name,
				Namespace: pvc.Namespace,
				Labels:    pvc.Labels,
			},
			Spec: pvc.Spec,
		}

		// Example: adjust storageClassName for cloud mapping
		if targetCloud == "eks" {
			sc := "gp3"
			outPVC.Spec.StorageClassName = &sc
		}

		y, _ := yaml.Marshal(outPVC)
		yamlDocs = append(yamlDocs, "---\n"+string(y))
		notes = append(notes, "Processed PVC "+pvc.Name+" (mapped for "+targetCloud+")")
	}
	return yamlDocs, len(pvcs.Items), notes
}
