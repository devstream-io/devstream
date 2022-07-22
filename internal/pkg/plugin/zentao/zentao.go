package zentao

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var (
	// Retry times for check zentao deployment status, currently this means 5 seconds * 120 times = 10 minutes
	retryTimes int = 120

	defaultPVPath map[string]string = map[string]string{
		"zentao-pv": "/www/zentaopms",
		"mysql-pv":  "/var/lib/mysql",
	}

	defaultZentaoPorts []corev1.ContainerPort = []corev1.ContainerPort{
		{
			Name:          "zentao",
			Protocol:      corev1.ProtocolTCP,
			ContainerPort: 80,
		},
		{
			Name:          "mysql",
			Protocol:      corev1.ProtocolTCP,
			ContainerPort: 3306,
		},
	}

	defaultSVCPort corev1.ServicePort = corev1.ServicePort{
		Name:       "zentao",
		Port:       80,
		TargetPort: intstr.IntOrString{IntVal: 80},
	}

	defaultZentaolabels map[string]string = map[string]string{
		"app": "zentao",
	}

	defaultName string = "zentao"
)
