package k8s

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/merico-dev/stream/internal/pkg/log"
)

func (c *Client) CreateServiceAccount(name, namespace string) error {
	sa := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}

	_, err := c.CoreV1().ServiceAccounts(namespace).Create(context.TODO(), sa, metav1.CreateOptions{})
	if err != nil {
		log.Errorf("Failed to create ServiceAccount < %s >: %s.", sa.Name, err)
		return err
	}

	log.Debugf("The ServiceAccount < %s > has created.", sa.Name)
	return nil
}

type CROption struct {
	Name        string
	PolicyRules []rbacv1.PolicyRule
}

func (c *Client) CreateClusterRole(option *CROption) error {
	cr := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: option.Name,
		},
		Rules: option.PolicyRules,
	}

	_, err := c.RbacV1().ClusterRoles().Create(context.TODO(), cr, metav1.CreateOptions{})
	if err != nil {
		log.Errorf("Failed to create ClusterRole < %s >: %s.", cr.Name, err)
		return err
	}

	log.Debugf("The ClusterRole < %s > has created.", cr.Name)
	return nil
}

type CRBOption struct {
	Name    string
	SANames []string
	RName   string
}

func (c *Client) CreateClusterRoleBinding(option *CRBOption) error {
	subjects := make([]rbacv1.Subject, 0)
	for _, name := range option.SANames {
		subjects = append(subjects, rbacv1.Subject{
			Kind:     "Group",
			APIGroup: "rbac.authorization.k8s.io",
			Name:     fmt.Sprintf("system:serviceaccounts:%s", name),
		})
	}

	crb := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: option.Name,
		},
		Subjects: subjects,
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     option.RName,
		},
	}

	_, err := c.RbacV1().ClusterRoleBindings().Create(context.TODO(), crb, metav1.CreateOptions{})
	if err != nil {
		log.Errorf("Failed to create ClusterRoleBinding < %s >: %s.", crb.Name, err)
		return err
	}

	log.Debugf("The ClusterRoleBinding < %s > has created.", crb.Name)
	return nil
}
