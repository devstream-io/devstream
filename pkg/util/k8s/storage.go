package k8s

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/devstream-io/devstream/pkg/util/log"
)

type PVOption struct {
	Name             string
	StorageClassName string

	// ReadWriteOnce PersistentVolumeAccessMode = "ReadWriteOnce"
	// ReadOnlyMany PersistentVolumeAccessMode = "ReadOnlyMany"
	// ReadWriteMany PersistentVolumeAccessMode = "ReadWriteMany"
	// ReadWriteOncePod PersistentVolumeAccessMode = "ReadWriteOncePod"
	AccessMode []corev1.PersistentVolumeAccessMode

	// <binarySI>        ::= Ki | Mi | Gi | Ti | Pi | Ei
	//   (International System of units; See: http://physics.nist.gov/cuu/Units/binary.html)
	// <decimalSI>       ::= m | "" | k | M | G | T | P | E
	//   (Note that 1024 = 1Ki but 1000 = 1k; I didn't choose the capitalization.)
	// eg: 10Gi 200Mi
	Capacity string

	//PersistentVolumeReclaimRecycle PersistentVolumeReclaimPolicy = "Recycle"
	//PersistentVolumeReclaimDelete PersistentVolumeReclaimPolicy = "Delete"
	//PersistentVolumeReclaimRetain PersistentVolumeReclaimPolicy = "Retain"
	PersistentVolumeReclaimPolicy corev1.PersistentVolumeReclaimPolicy
	HostPath                      string
}

func (c *Client) CreatePersistentVolume(option *PVOption) error {
	quantity, err := resource.ParseQuantity(option.Capacity)
	if err != nil {
		log.Errorf("Failed to parse the Capacity string: %s.", err)
		return err
	}

	pv := &corev1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name: option.Name,
		},
		Spec: corev1.PersistentVolumeSpec{
			StorageClassName: option.StorageClassName,
			AccessModes:      option.AccessMode,
			Capacity: map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceStorage: quantity,
			},
			PersistentVolumeReclaimPolicy: option.PersistentVolumeReclaimPolicy,
			PersistentVolumeSource: corev1.PersistentVolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: option.HostPath,
				},
			},
		},
	}

	_, err = c.CoreV1().PersistentVolumes().Create(context.TODO(), pv, metav1.CreateOptions{})
	if err != nil {
		log.Errorf("Failed to create PersistentVolume < %s >: %s.", pv.Name, err)
		return err
	}

	log.Debugf("The PersistentVolume < %s > has created.", pv.Name)
	return nil
}
