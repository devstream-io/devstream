package configGetter

import (
	"fmt"
	"strings"

	"github.com/devstream-io/devstream/pkg/util/log"
)

// ItemGetter is used to get value from the specify key
type ItemGetter interface {
	Get() string
	DescribeWhereToSet() string
}

type itemGetterChain []ItemGetter

func NewItemGetterChain(getters ...ItemGetter) itemGetterChain {
	return getters
}

// Get value from the getters chain util non-empty value is found
// It will return an error to instruct user how to set the value if value not found
func (c itemGetterChain) Get() (string, error) {
	for _, getter := range c {
		if value := getter.Get(); value != "" {
			log.Debugf("get value <%s> from %s", value, getter.DescribeWhereToSet())
			return value, nil
		}
	}
	return "", NewErrItemNotFound(c...)
}

// SingleGet gets value from a single getter,
// it's just an alias of itemGetterChain.Get
func SingleGet(g ItemGetter) (string, error) {
	return NewItemGetterChain(g).Get()
}

// CheckItemExist checks if the item exists in the getter chain
// It will return an error to instruct user how to set the value if value not found
func CheckItemExist(getters ...ItemGetter) error {
	for _, getter := range getters {
		if value := getter.Get(); value != "" {
			return nil
		}
	}
	return NewErrItemNotFound(getters...)
}

type errItemNotFound struct {
	getters []ItemGetter
}

func NewErrItemNotFound(getters ...ItemGetter) *errItemNotFound {
	return &errItemNotFound{
		getters: getters,
	}
}

func (e errItemNotFound) Error() string {
	switch len(e.getters) {
	case 0:
		return ""
	default:
		var hints []string
		for _, getter := range e.getters {
			hints = append(hints, getter.DescribeWhereToSet())
		}
		return fmt.Sprintf("missing config settings, you could set it by: %s", strings.Join(hints, ", or "))
	}
}
