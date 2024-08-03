package ops

import (
	"context"
	"fmt"
	libovsdbclient "github.com/ovn-org/libovsdb/client"
	"github.com/ovn-org/ovn-kubernetes/go-controller/pkg/types"
	"github.com/ovn-org/ovn-kubernetes/go-controller/pkg/vswitchd"
)

// GetOVSOpenvSwitchTable from the cache
func GetOVSOpenvSwitchTable(ovsClient libovsdbclient.Client) (*vswitchd.OpenvSwitch, error) {
	ctx, cancel := context.WithTimeout(context.Background(), types.OVSDBTimeout)
	defer cancel()
	openvSwitchTableList := []*vswitchd.OpenvSwitch{}
	err := ovsClient.List(ctx, &openvSwitchTableList)
	if err != nil {
		return nil, err
	}
	if len(openvSwitchTableList) == 0 {
		return nil, fmt.Errorf("no openvSwitch table found")
	}

	return openvSwitchTableList[0], err
}
