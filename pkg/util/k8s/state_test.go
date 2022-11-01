package k8s

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/utils/pointer"
)

var _ = Describe("k8s state methods", func() {
	var (
		client                                                                                *Client
		labelFilter                                                                           map[string]string
		anFilter                                                                              map[string]string
		namespace, testFilterKey, testFilterVal, testAnKey, testAnVal, successName, faildName string
	)
	BeforeEach(func() {
		client = &Client{}
		namespace = "test"
		testFilterKey = "filter_key"
		testFilterVal = "filter_val"
		testAnKey = "an_key"
		testAnVal = "an_val"
		successName = "test_success"
		faildName = "test_fail"
		anFilter = map[string]string{
			testAnKey: testAnVal,
		}
		labelFilter = map[string]string{
			testFilterKey: testFilterVal,
		}
		podTemplate := corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: labelFilter,
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:  "test",
						Image: "test",
					},
				},
			},
		}

		depSpec := appsv1.DeploymentSpec{
			Replicas: pointer.Int32Ptr(int32(1)),
			Selector: &metav1.LabelSelector{MatchLabels: labelFilter},
			Template: podTemplate,
		}
		dsSpec := appsv1.DaemonSetSpec{
			Template: podTemplate,
		}
		stsSpec := appsv1.StatefulSetSpec{
			Replicas: pointer.Int32Ptr(int32(1)),
			Selector: &metav1.LabelSelector{MatchLabels: labelFilter},
			Template: podTemplate,
		}

		testResources := []runtime.Object{
			&appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:        successName,
					Labels:      labelFilter,
					Annotations: anFilter,
					Namespace:   namespace,
				},
				Spec: depSpec,
			},
			&appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      faildName,
					Labels:    map[string]string{testFilterKey: "not_val"},
					Namespace: namespace,
				},
				Spec: depSpec,
			},
			&appsv1.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:        successName,
					Labels:      labelFilter,
					Namespace:   namespace,
					Annotations: anFilter,
				},
				Spec: dsSpec,
			},
			&appsv1.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:        faildName,
					Labels:      labelFilter,
					Annotations: map[string]string{testAnKey: "not_exist"},
					Namespace:   namespace,
				},
				Spec: dsSpec,
			},
			&appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:        successName,
					Labels:      labelFilter,
					Annotations: anFilter,
					Namespace:   namespace,
				},
				Spec: stsSpec,
			},
			&appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:      faildName,
					Labels:    labelFilter,
					Namespace: namespace,
				},
				Spec: stsSpec,
			},
		}
		client.clientset = fake.NewSimpleClientset(testResources...)
	})
	Context("GetResourceStatus method", func() {
		It("should work normal", func() {
			allStatus, err := client.GetResourceStatus(namespace, anFilter, labelFilter)
			Expect(err).Error().ShouldNot(HaveOccurred())
			Expect(len(allStatus.DaemonSet)).Should(Equal(1))
			Expect(allStatus.DaemonSet[0].Name).Should(Equal(successName))
			Expect(len(allStatus.Deployment)).Should(Equal(1))
			Expect(allStatus.Deployment[0].Name).Should(Equal(successName))
			Expect(len(allStatus.StatefulSet)).Should(Equal(1))
			Expect(allStatus.StatefulSet[0].Name).Should(Equal(successName))
		})
	})
})

var _ = Describe("AllResourceStatus struct", func() {
	var (
		depName  string
		depReady bool
		s        *AllResourceStatus
	)
	Context("ToStringInterfaceMap method", func() {
		BeforeEach(func() {
			depName = "test_dep"
			depReady = true
			s = &AllResourceStatus{
				Deployment: []ResourceStatus{
					{
						Name:  "test_dep",
						Ready: depReady,
					},
				},
			}
		})
		It("should work", func() {
			m, err := s.ToStringInterfaceMap()
			Expect(err).ShouldNot(HaveOccurred())
			expectVal := fmt.Sprintf("deployment:\n  - name: %s\n    ready: %t\nstatefulset: []\ndaemonset: []\n", depName, depReady)
			Expect(m).Should(Equal(map[string]interface{}{
				"workflows": expectVal,
			}))
		})
	})
})
