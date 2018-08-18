package specs

import (
	"fmt"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/intstr"

	"manno.name/mm/fraas/fraas-gke/models"
	fh "manno.name/mm/fraas/fraas-helpers"
)

func int32Ptr(i int32) *int32 { return &i }

func NewRailsConfigMap(d *models.Deployment, c *fh.SiteConfig) *apiv1.ConfigMap {
	return &apiv1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: d.ConfigName(),
		},
		Data: map[string]string{
			"FRAB_CURRENCY_FORMAT":     "%n%u",
			"FRAB_CURRENCY_UNIT":       "€",
			"FRAB_HOST":                d.ExternalDomain,
			"FRAB_PROTOCOL":            "https",
			"FROM_EMAIL":               d.FromEmail,
			"SMTP_ADDRESS":             c.Mail.SMTPServer,
			"SMTP_NOTLS":               c.Mail.SMTPNOTLS,
			"SMTP_PORT":                c.Mail.SMTPServerPort,
			"RACK_ENV":                 "production",
			"RAILS_SERVE_STATIC_FILES": "true",
			"RAILS_LOG_TO_STDOUT":      "true",
			"EXCEPTION_EMAIL":          c.Mail.ExceptionEMail,
		},
	}
}

func NewRailsSecretsSpec(d *models.Deployment) *apiv1.Secret {
	return &apiv1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: d.SecretName(),
		},
		StringData: map[string]string{
			"SECRET_KEY_BASE": d.SecretKeyBase,
			"SMTP_USER_NAME":  fh.Config().Mail.SMTPUsername,
			"SMTP_PASSWORD":   fh.Config().Mail.SMTPPassword,
			"DATABASE_URL":    fmt.Sprintf("postgres://%s:%s@127.0.0.1:5432/%s?sslmode=disable", d.DatabaseID, d.DatabasePassword, d.DatabaseID),
		},
	}
}

func NewServiceSpec(d *models.Deployment) *apiv1.Service {
	return &apiv1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind: "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   d.BackendName(),
			Labels: map[string]string{"app": d.DeploymentID()},
		},
		Spec: apiv1.ServiceSpec{
			Type: "NodePort",
			Ports: []apiv1.ServicePort{
				apiv1.ServicePort{
					Port:       int32(3000),
					TargetPort: intstr.FromInt(3000),
				},
			},
			Selector: map[string]string{
				"tier": "web",
				"app":  d.DeploymentID(),
			},
		},
	}
}

func NewCertificate(d *models.Deployment) *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "certmanager.k8s.io/v1alpha1",
			"kind":       "Certificate",
			"metadata": map[string]string{
				"name":      d.TLSSecretName(),
				"namespace": "default",
			},
			"spec": map[string]interface{}{
				"secretName": d.TLSSecretName(),
				"issuerRef": map[string]string{
					"name": "letsencrypt-prod",
					"kind": "ClusterIssuer",
				},
				"commonName": d.Domain,
				"dnsNames":   []string{d.Domain},
				"acme": map[string]interface{}{
					"config": []interface{}{
						map[string]interface{}{
							"http01":  map[string]string{"ingress": d.WebName()},
							"domains": []string{d.Domain},
						},
					},
				},
			},
		},
	}
}
