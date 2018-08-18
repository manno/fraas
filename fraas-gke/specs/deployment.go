package specs

import (
	"fmt"

	appv1 "k8s.io/api/apps/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"manno.name/mm/fraas/fraas-gke/models"
	fh "manno.name/mm/fraas/fraas-helpers"
)

func NewDeploymentSpec(d *models.Deployment, config *fh.SiteConfig) *appv1.Deployment {
	podSpec := deploymentPodSpec(d, config)
	deployment := &appv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "extensions/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   d.WebName(),
			Labels: map[string]string{"app": d.DeploymentID()},
		},
		Spec: appv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"tier": "web",
						"app":  d.DeploymentID(),
					},
				},
				Spec: podSpec,
			},
		},
	}
	return deployment
}

func deploymentPodSpec(d *models.Deployment, config *fh.SiteConfig) apiv1.PodSpec {
	configName := d.ConfigName()
	secretName := d.SecretName()
	return apiv1.PodSpec{
		Volumes: []apiv1.Volume{
			apiv1.Volume{
				Name: "cloudsql-instance-credentials",
				VolumeSource: apiv1.VolumeSource{
					Secret: &apiv1.SecretVolumeSource{
						SecretName: "cloudsql-fraas-credentials",
					},
				},
			},
		},
		Containers: []apiv1.Container{
			apiv1.Container{
				Name:  d.RailsContainerName(),
				Image: config.DockerImage,
				Ports: []apiv1.ContainerPort{
					apiv1.ContainerPort{
						Name:          "http-port",
						ContainerPort: 3000,
					},
				},
				EnvFrom: []apiv1.EnvFromSource{
					apiv1.EnvFromSource{
						ConfigMapRef: &apiv1.ConfigMapEnvSource{
							LocalObjectReference: apiv1.LocalObjectReference{Name: configName},
						},
					},
				},
				Env: []apiv1.EnvVar{
					apiv1.EnvVar{
						Name: "SECRET_KEY_BASE",
						ValueFrom: &apiv1.EnvVarSource{
							SecretKeyRef: &apiv1.SecretKeySelector{
								LocalObjectReference: apiv1.LocalObjectReference{Name: secretName},
								Key:                  "SECRET_KEY_BASE",
							},
						},
					},
					apiv1.EnvVar{
						Name: "DATABASE_URL",
						ValueFrom: &apiv1.EnvVarSource{
							SecretKeyRef: &apiv1.SecretKeySelector{
								LocalObjectReference: apiv1.LocalObjectReference{Name: secretName},
								Key:                  "DATABASE_URL",
							},
						},
					},
					apiv1.EnvVar{
						Name: "SMTP_USER_NAME",
						ValueFrom: &apiv1.EnvVarSource{
							SecretKeyRef: &apiv1.SecretKeySelector{
								LocalObjectReference: apiv1.LocalObjectReference{Name: secretName},
								Key:                  "SMTP_USER_NAME",
							},
						},
					},
					apiv1.EnvVar{
						Name: "SMTP_PASSWORD",
						ValueFrom: &apiv1.EnvVarSource{
							SecretKeyRef: &apiv1.SecretKeySelector{
								LocalObjectReference: apiv1.LocalObjectReference{Name: secretName},
								Key:                  "SMTP_PASSWORD",
							},
						},
					},
				},
			},
			apiv1.Container{
				Name:  "cloudsql-proxy",
				Image: "gcr.io/cloudsql-docker/gce-proxy:1.11",
				Command: []string{
					"/cloud_sql_proxy",
					fmt.Sprintf("-instances=%s:%s:%s=tcp:5432", config.Google.ProjectID, config.Google.Region, config.Database.Instance),
					"-credential_file=/secrets/cloudsql/credentials.json",
				},
				VolumeMounts: []apiv1.VolumeMount{
					apiv1.VolumeMount{
						Name:      "cloudsql-instance-credentials",
						ReadOnly:  true,
						MountPath: "/secrets/cloudsql",
					},
				},
			},
		},
	}
}
