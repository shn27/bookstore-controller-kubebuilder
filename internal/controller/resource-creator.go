package controller

import (
	appv1 "github.com/bookstore-controller-kubebuilder/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	//samplev1alpha1 "k8s.io/sample-controller/pkg/apis/samplecontroller/v1alpha1"
)

func getService(bookstore *appv1.BookStore) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      bookstore.Spec.DeploymentName,
			Namespace: bookstore.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				//*metav1.NewControllerRef(bookstore, samplev1alpha1.SchemeGroupVersion.WithKind("BookStore")), // TODO
				*metav1.NewControllerRef(bookstore, appv1.GroupVersion.WithKind("BookStore")),
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: bookstore.GetSelectorLabels(),
			Type:     corev1.ServiceTypeLoadBalancer,
			Ports: []corev1.ServicePort{
				{
					Port:       bookstore.Spec.ContainerPort,
					TargetPort: intstr.FromInt(int(bookstore.Spec.ContainerPort)),
					NodePort:   30000,
				},
			},
		},
	}
}

func getDeployment(bookstore *appv1.BookStore) *appsv1.Deployment {
	labels := bookstore.GetSelectorLabels()
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      bookstore.Spec.DeploymentName,
			Namespace: bookstore.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(bookstore, appv1.GroupVersion.WithKind("BookStore")),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: bookstore.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  bookstore.Spec.DeploymentName,
							Image: bookstore.Spec.DeploymentImageName + ":" + bookstore.Spec.DeploymentImageTag,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: bookstore.Spec.ContainerPort,
								},
							},
						},
					},
				},
			},
		},
	}
}
