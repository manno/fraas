package specs

import (
	ext "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"manno.name/mm/fraas/fraas-gke/models"
)

func NewIngressSpec(d *models.Deployment) *ext.Ingress {
	return &ext.Ingress{
		TypeMeta: metav1.TypeMeta{
			Kind: "Ingress",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        d.WebName(),
			Labels:      map[string]string{"app": d.DeploymentID()},
			Annotations: map[string]string{"kubernetes.io/ingress.global-static-ip-name": d.IPName()},
		},
		Spec: ext.IngressSpec{
			Backend: &ext.IngressBackend{
				ServiceName: d.BackendName(),
				ServicePort: intstr.FromInt(3000),
			},
		},
	}
}

func NewIngressTLSSpec(d *models.Deployment) *ext.Ingress {
	return &ext.Ingress{
		TypeMeta: metav1.TypeMeta{
			Kind: "Ingress",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        d.WebName(),
			Labels:      map[string]string{"app": d.DeploymentID()},
			Annotations: map[string]string{"kubernetes.io/ingress.global-static-ip-name": d.IPName()},
		},
		Spec: ext.IngressSpec{
			Backend: &ext.IngressBackend{
				ServiceName: d.BackendName(),
				ServicePort: intstr.FromInt(3000),
			},
			TLS: []ext.IngressTLS{
				ext.IngressTLS{
					SecretName: d.TLSSecretName(),
					Hosts:      []string{d.Domain},
				},
			},
		},
	}
}
