package ovn

import (
	libovsdbclient "github.com/ovn-org/libovsdb/client"
	"github.com/ovn-org/libovsdb/ovsdb"

	ovsdbops "github.com/ovn-org/ovn-kubernetes/go-controller/pkg/libovsdb/ops/ovsdb"
	"github.com/ovn-org/ovn-kubernetes/go-controller/pkg/nbdb"
)

type coppPredicate func(*nbdb.Copp) bool

// CreateOrUpdateCOPPsOps creates or updates the provided COPP returning the
// corresponding ops
func CreateOrUpdateCOPPsOps(nbClient libovsdbclient.Client, ops []ovsdb.Operation, copps ...*nbdb.Copp) ([]ovsdb.Operation, error) {
	opModels := make([]ovsdbops.OperationModel, 0, len(copps))
	for i := range copps {
		// can't use i in the predicate, for loop replaces it in-memory
		copp := copps[i]
		opModel := ovsdbops.OperationModel{
			Model:          copp,
			OnModelUpdates: ovsdbops.OnModelUpdatesAllNonDefault(),
			ErrNotFound:    false,
			BulkOp:         false,
		}
		opModels = append(opModels, opModel)
	}

	modelClient := ovsdbops.NewModelClient(nbClient)
	return modelClient.CreateOrUpdateOps(ops, opModels...)
}

// DeleteCOPPsOps deletes the provided COPPs found using the predicate, returning the
// corresponding ops
func DeleteCOPPsWithPredicateOps(nbClient libovsdbclient.Client, ops []ovsdb.Operation, p coppPredicate) ([]ovsdb.Operation, error) {
	copp := nbdb.Copp{}
	opModels := []ovsdbops.OperationModel{
		{
			Model:          &copp,
			ModelPredicate: p,
			ErrNotFound:    false,
			BulkOp:         true,
		},
	}

	modelClient := ovsdbops.NewModelClient(nbClient)
	return modelClient.DeleteOps(ops, opModels...)
}
