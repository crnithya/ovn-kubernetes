package ovn

import (
	"context"

	libovsdbclient "github.com/ovn-org/libovsdb/client"
	"github.com/ovn-org/libovsdb/ovsdb"

	"github.com/ovn-org/ovn-kubernetes/go-controller/pkg/config"
	ovsdbops "github.com/ovn-org/ovn-kubernetes/go-controller/pkg/libovsdb/ops/ovsdb"
	"github.com/ovn-org/ovn-kubernetes/go-controller/pkg/nbdb"
)

// CreateOrUpdateLoadBalancerGroupOps returns the ops to create or update the
// provided load balancer group
func CreateOrUpdateLoadBalancerGroupOps(nbClient libovsdbclient.Client, ops []ovsdb.Operation, group *nbdb.LoadBalancerGroup) ([]ovsdb.Operation, error) {
	// lb group has no fields other than name, safe to update just with non-default values
	opModel := ovsdbops.OperationModel{
		Model:          group,
		OnModelUpdates: ovsdbops.OnModelUpdatesAllNonDefault(),
		ErrNotFound:    false,
		BulkOp:         false,
	}

	m := ovsdbops.NewModelClient(nbClient)
	ops, err := m.CreateOrUpdateOps(ops, opModel)
	if err != nil {
		return nil, err
	}
	return ops, nil
}

// DeleteLoadBalancerGroupsOps DeleteLoadBalncerGroupOps creates the operations for deleting load balancer groups
func DeleteLoadBalancerGroupsOps(nbClient libovsdbclient.Client, ops []ovsdb.Operation, groups ...*nbdb.LoadBalancerGroup) ([]ovsdb.Operation, error) {
	opModels := make([]ovsdbops.OperationModel, 0, len(groups))
	for i := range groups {
		// can't use i in the predicate, for loop replaces it in-memory
		lb := groups[i]
		opModel := ovsdbops.OperationModel{
			Model:       lb,
			ErrNotFound: false,
			BulkOp:      false,
		}
		opModels = append(opModels, opModel)
	}

	modelClient := ovsdbops.NewModelClient(nbClient)
	return modelClient.DeleteOps(ops, opModels...)
}

// DeleteLoadBalancerGroups deletes the provided load balancer groups
func DeleteLoadBalancerGroups(nbClient libovsdbclient.Client, groups []*nbdb.LoadBalancerGroup) error {
	ops, err := DeleteLoadBalancerGroupsOps(nbClient, nil, groups...)
	if err != nil {
		return err
	}

	_, err = ovsdbops.TransactAndCheck(nbClient, ops)
	return err
}

// AddLoadBalancersToGroupOps adds the provided load balancers to the provided
// group and returns the corresponding ops
func AddLoadBalancersToGroupOps(nbClient libovsdbclient.Client, ops []ovsdb.Operation, group *nbdb.LoadBalancerGroup, lbs ...*nbdb.LoadBalancer) ([]ovsdb.Operation, error) {
	originalLBs := group.LoadBalancer
	group.LoadBalancer = make([]string, 0, len(lbs))
	for _, lb := range lbs {
		group.LoadBalancer = append(group.LoadBalancer, lb.UUID)
	}
	opModel := ovsdbops.OperationModel{
		Model:            group,
		ModelPredicate:   func(item *nbdb.LoadBalancerGroup) bool { return item.Name == group.Name },
		OnModelMutations: []interface{}{&group.LoadBalancer},
		ErrNotFound:      true,
		BulkOp:           false,
	}

	m := ovsdbops.NewModelClient(nbClient)
	ops, err := m.CreateOrUpdateOps(ops, opModel)
	group.LoadBalancer = originalLBs
	return ops, err
}

// RemoveLoadBalancersFromGroupOps removes the provided load balancers from the
// provided group and returns the corresponding ops
func RemoveLoadBalancersFromGroupOps(nbClient libovsdbclient.Client, ops []ovsdb.Operation, group *nbdb.LoadBalancerGroup, lbs ...*nbdb.LoadBalancer) ([]ovsdb.Operation, error) {
	originalLBs := group.LoadBalancer
	group.LoadBalancer = make([]string, 0, len(lbs))
	for _, lb := range lbs {
		group.LoadBalancer = append(group.LoadBalancer, lb.UUID)
	}
	opModel := ovsdbops.OperationModel{
		Model:            group,
		ModelPredicate:   func(item *nbdb.LoadBalancerGroup) bool { return item.Name == group.Name },
		OnModelMutations: []interface{}{&group.LoadBalancer},
		// if we want to delete loadbalancer from the port group that doesn't exist, that is noop
		ErrNotFound: false,
		BulkOp:      false,
	}

	m := ovsdbops.NewModelClient(nbClient)
	ops, err := m.DeleteOps(ops, opModel)
	group.LoadBalancer = originalLBs
	return ops, err
}

type loadBalancerGroupPredicate func(*nbdb.LoadBalancerGroup) bool

// FindLoadBalancerGroupsWithPredicate looks up load balancer groups from the
// cache based on a given predicate
func FindLoadBalancerGroupsWithPredicate(nbClient libovsdbclient.Client, p loadBalancerGroupPredicate) ([]*nbdb.LoadBalancerGroup, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Default.OVSDBTxnTimeout)
	defer cancel()
	groups := []*nbdb.LoadBalancerGroup{}
	err := nbClient.WhereCache(p).List(ctx, &groups)
	return groups, err
}
