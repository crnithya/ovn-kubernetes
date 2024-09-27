package ovn

import (
	libovsdbclient "github.com/ovn-org/libovsdb/client"

	libovsdbops "github.com/ovn-org/ovn-kubernetes/go-controller/pkg/libovsdb/ops/ovsdb"
	"github.com/ovn-org/ovn-kubernetes/go-controller/pkg/sbdb"
)

// GetNBGlobal looks up the SB Global entry from the cache
func GetSBGlobal(sbClient libovsdbclient.Client, sbGlobal *sbdb.SBGlobal) (*sbdb.SBGlobal, error) {
	found := []*sbdb.SBGlobal{}
	opModel := libovsdbops.OperationModel{
		Model:          sbGlobal,
		ModelPredicate: func(item *sbdb.SBGlobal) bool { return true },
		ExistingResult: &found,
		ErrNotFound:    true,
		BulkOp:         false,
	}

	m := libovsdbops.NewModelClient(sbClient)
	err := m.Lookup(opModel)
	if err != nil {
		return nil, err
	}

	return found[0], nil
}
