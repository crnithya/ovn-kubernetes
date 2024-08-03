package ops

import (
	"context"
	"fmt"
	libovsdbclient "github.com/ovn-org/libovsdb/client"
	"github.com/ovn-org/ovn-kubernetes/go-controller/pkg/types"
	"github.com/ovn-org/ovn-kubernetes/go-controller/pkg/vswitchd"
)

// ListOVSPorts looks up all ovs bridge ports from the cache
func ListOVSPorts(ovsClient libovsdbclient.Client) ([]*vswitchd.Port, error) {
	ctx, cancel := context.WithTimeout(context.Background(), types.OVSDBTimeout)
	defer cancel()
	searchedPorts := []*vswitchd.Port{}
	err := ovsClient.List(ctx, &searchedPorts)
	return searchedPorts, err
}

// ListOVSPortsByPredicate returns all the ovs ports in the cache that matches the lookup function
func ListOVSPortsByPredicate(ovsClient libovsdbclient.Client, lookupFunction func(item *vswitchd.Port) bool) ([]*vswitchd.Port, error) {
	ctx, cancel := context.WithTimeout(context.Background(), types.OVSDBTimeout)
	defer cancel()
	searchedPorts := []*vswitchd.Port{}

	err := ovsClient.WhereCache(lookupFunction).List(ctx, &searchedPorts)
	if err != nil {
		return nil, fmt.Errorf("failed listing ports : %v", err)
	}

	return searchedPorts, nil
}
