package ops

import (
	"context"
	"fmt"
	libovsdbclient "github.com/ovn-org/libovsdb/client"
	"github.com/ovn-org/ovn-kubernetes/go-controller/pkg/types"
	"github.com/ovn-org/ovn-kubernetes/go-controller/pkg/vswitchd"
)

// ListOVSBridges looks up all bridges from the cache
func ListOVSBridges(ovsClient libovsdbclient.Client) ([]*vswitchd.Bridge, error) {
	ctx, cancel := context.WithTimeout(context.Background(), types.OVSDBTimeout)
	defer cancel()
	searchedBridges := []*vswitchd.Bridge{}
	err := ovsClient.List(ctx, &searchedBridges)
	return searchedBridges, err
}

// ListBridgeByPredicate returns all the bridges in the cache that matches the lookup function
func ListOVSBridgesByPredicate(ovsClient libovsdbclient.Client, lookupFunction func(item *vswitchd.Bridge) bool) ([]*vswitchd.Bridge, error) {
	ctx, cancel := context.WithTimeout(context.Background(), types.OVSDBTimeout)
	defer cancel()
	searchedBridges := []*vswitchd.Bridge{}

	err := ovsClient.WhereCache(lookupFunction).List(ctx, &searchedBridges)
	if err != nil {
		return nil, fmt.Errorf("failed listing bridges : %v", err)
	}

	return searchedBridges, nil
}
