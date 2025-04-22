package controller

import (
	"context"

	bookinfov1beta1 "com/bookinfocrd/api/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func pathTypePtr(pt networkingv1.PathType) *networkingv1.PathType {
	return &pt
}
func strPtr(s string) *string {
	return &s
}
func imagePullSecretsConvert(imagePullSecrets []string) []corev1.LocalObjectReference {
	res := make([]corev1.LocalObjectReference, len(imagePullSecrets))
	for index, v := range imagePullSecrets {
		res[index] = corev1.LocalObjectReference{
			Name: v,
		}
	}
	return res
}
func (r *BookinfoReconciler) deleteBookinfo(ctx context.Context, bookinfo bookinfov1beta1.Bookinfo) error {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      bookinfo.Name,
			Namespace: bookinfo.Namespace,
		},
	}
	// 删除该 Deployment
	return r.Client.Delete(ctx, deployment)
}
func (r *BookinfoReconciler) createService(ctx context.Context, deployment *appsv1.Deployment) error {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deployment.Name,
			Namespace: deployment.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(deployment, appsv1.SchemeGroupVersion.WithKind("Deployment")),
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": deployment.Name, // Service 会选择匹配相同标签的 Pod
			},
			Ports: []corev1.ServicePort{
				{
					Protocol:   corev1.ProtocolTCP,
					Port:       80,                   // Service 的端口
					TargetPort: intstr.FromInt(8080), // Pod 的端口
				},
			},
		},
	}
	return r.Client.Create(ctx, service)
}
func (r *BookinfoReconciler) createOrUpdateIngress(ctx context.Context, deployment *appsv1.Deployment, host string) error {
	// 定义 Ingress
	ingress := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deployment.Name,
			Namespace: deployment.Namespace,
			Annotations: map[string]string{
				"nginx.ingress.kubernetes.io/rewrite-target": "/",
			},
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(deployment, appsv1.SchemeGroupVersion.WithKind("Deployment")),
			},
		},
		Spec: networkingv1.IngressSpec{
			IngressClassName: strPtr("nginx"),
			Rules: []networkingv1.IngressRule{
				{
					Host: host,
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: pathTypePtr(networkingv1.PathTypePrefix),
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: deployment.Name, // 指向对应的 Service
											Port: networkingv1.ServiceBackendPort{
												Number: 80, // 对应 Service 的端口
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
	existingIngress := &networkingv1.Ingress{}
	err := r.Client.Get(ctx, client.ObjectKey{Name: deployment.Name, Namespace: deployment.Namespace}, existingIngress)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return r.Client.Create(ctx, ingress)
		} else {
			return err
		}
	} else {
		return r.Client.Update(ctx, ingress)
	}

}
func (r *BookinfoReconciler) createOrUpdateDeployment(ctx context.Context, deployment *appsv1.Deployment, host string) error {
	//如果已存在则更新
	existingDeployment := &appsv1.Deployment{}
	err := r.Client.Get(ctx, client.ObjectKey{Name: deployment.Name, Namespace: deployment.Namespace}, existingDeployment)
	if err != nil {
		if apierrors.IsNotFound(err) {
			if err := r.Client.Create(ctx, deployment); err != nil {
				return err
			}
			//创建 Service
			if err := r.createService(ctx, deployment); err != nil {
				return err
			}
			if err := r.createOrUpdateIngress(ctx, deployment, host); err != nil {
				return err
			}
			return nil
		}
		// 其他错误处理
		return err
	} else {
		// 已存在，更新操作
		if err := r.Client.Update(ctx, deployment); err != nil {
			return err
		}
		if err := r.createOrUpdateIngress(ctx, deployment, host); err != nil {
			return err
		}
		return nil
	}
}

func (r *BookinfoReconciler) createBookinfo(ctx context.Context, bookinfo bookinfov1beta1.Bookinfo) error {
	// 定义 Deployment
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      bookinfo.Name,
			Namespace: bookinfo.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &bookinfo.Spec.Replicas, // 副本数
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": bookinfo.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": bookinfo.Name,
					},
				},
				Spec: corev1.PodSpec{
					ImagePullSecrets: imagePullSecretsConvert(bookinfo.Spec.ImagePullSecrets),
					Containers: []corev1.Container{
						{
							Name:  "main",
							Image: bookinfo.Spec.Image,
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: 8080,
								},
							},
						},
					},
				},
			},
		},
	}
	return r.createOrUpdateDeployment(ctx, deployment, bookinfo.Spec.Host)
}
