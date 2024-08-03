package ops

import (
	"context"
	"fmt"
	libovsdbclient "github.com/ovn-org/libovsdb/client"
	"github.com/ovn-org/ovn-kubernetes/go-controller/pkg/types"
	"github.com/ovn-org/ovn-kubernetes/go-controller/pkg/vswitchd"
)

// ListOVSInterfaces looks up all ovs interfaces from the cache
func ListOVSInterfaces(ovsClient libovsdbclient.Client) ([]*vswitchd.Interface, error) {
	ctx, cancel := context.WithTimeout(context.Background(), types.OVSDBTimeout)
	defer cancel()
	searchedInterfaces := []*vswitchd.Interface{}
	err := ovsClient.List(ctx, &searchedInterfaces)
	return searchedInterfaces, err
}

// ListOVSInterfacesByPredicate returns all the ovs interfaces in the cache that matches the lookup function
func ListOVSInterfacesByPredicate(ovsClient libovsdbclient.Client, lookupFunction func(item *vswitchd.Interface) bool) ([]*vswitchd.Interface, error) {
	ctx, cancel := context.WithTimeout(context.Background(), types.OVSDBTimeout)
	defer cancel()
	searchedInterfaces := []*vswitchd.Interface{}

	err := ovsClient.WhereCache(lookupFunction).List(ctx, &searchedInterfaces)
	if err != nil {
		return nil, fmt.Errorf("failed listing interfaces : %v", err)
	}

	return searchedInterfaces, nil
}
