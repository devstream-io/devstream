package k8s

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/devstream-io/devstream/pkg/util/log"
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
		log.Errorf("Failed to create the ServiceAccount < %s >: %s.", sa.Name, err)
		return err
	}

	log.Debugf("The ServiceAccount < %s > has been created.", sa.Name)
	return nil
}

func (c *Client) DeleteServiceAccount(name, namespace string) error {
	err := c.CoreV1().ServiceAccounts(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		log.Errorf("Failed to delete the ServiceAccount < %s >: %s.", name, err)
		return err
	}

	log.Debugf("The ServiceAccount < %s > has been deleted.", name)
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
		log.Errorf("Failed to create the ClusterRole < %s >: %s.", cr.Name, err)
		return err
	}

	log.Debugf("The ClusterRole < %s > has been created.", cr.Name)
	return nil
}

func (c *Client) DeleteClusterRole(name string) error {
	err := c.RbacV1().ClusterRoles().Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		log.Errorf("Failed to delete the ClusterRole < %s >: %s.", name, err)
		return err
	}

	log.Debugf("The ClusterRole < %s > has been deleted.", name)
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
		log.Errorf("Failed to create the ClusterRoleBinding < %s >: %s.", crb.Name, err)
		return err
	}

	log.Debugf("The ClusterRoleBinding < %s > has been created.", crb.Name)
	return nil
}

func (c *Client) DeleteClusterRoleBinding(name string) error {
	err := c.RbacV1().ClusterRoleBindings().Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		log.Errorf("Failed to delete the ClusterRoleBinding < %s >: %s.", name, err)
		return err
	}

	log.Debugf("The ClusterRoleBinding < %s > has been deleted.", name)
	return nil
}
