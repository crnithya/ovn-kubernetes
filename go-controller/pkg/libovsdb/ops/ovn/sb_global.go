package ovn

import (
	libovsdbclient "github.com/ovn-org/libovsdb/client"

	ovsdbops "github.com/ovn-org/ovn-kubernetes/go-controller/pkg/libovsdb/ops/ovsdb"
	"github.com/ovn-org/ovn-kubernetes/go-controller/pkg/sbdb"
)

// GetSBGlobal looks up the SB Global entry from the cache
func GetSBGlobal(sbClient libovsdbclient.Client, sbGlobal *sbdb.SBGlobal) (*sbdb.SBGlobal, error) {
	found := []*sbdb.SBGlobal{}
	opModel := ovsdbops.OperationModel{
		Model:          sbGlobal,
		ModelPredicate: func(_ *sbdb.SBGlobal) bool { return true },
		ExistingResult: &found,
		ErrNotFound:    true,
		BulkOp:         false,
	}

	m := ovsdbops.NewModelClient(sbClient)
	err := m.Lookup(opModel)
	if err != nil {
		return nil, err
	}

	return found[0], nil
}
