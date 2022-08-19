package k8s

import "github.com/devstream-io/devstream/pkg/util/log"

func (c *Client) UpsertNameSpace(nameSpace string) error {
	// check if namespace exist
	exist, err := c.IsNamespaceExists(nameSpace)
	if err != nil {
		log.Debugf("Failed to check the namespace exist: %s.", nameSpace)
		return err
	}
	if !exist {
		return c.CreateNamespace(nameSpace)
	}
	log.Debugf("The namespace %s has been existed.", nameSpace)
	return nil
}
