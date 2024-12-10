package model

import (
	"testing"

	ovsdbops "github.com/ovn-org/ovn-kubernetes/go-controller/pkg/libovsdb/ops/ovsdb"
	"github.com/ovn-org/ovn-kubernetes/go-controller/pkg/nbdb"
)

var mapping = map[string]string{
	egressFirewallOwnerType:             ovsdbops.EgressFirewallOwnerType,
	adminNetworkPolicyOwnerType:         ovsdbops.AdminNetworkPolicyOwnerType,
	baselineAdminNetworkPolicyOwnerType: ovsdbops.BaselineAdminNetworkPolicyOwnerType,
	networkPolicyOwnerType:              ovsdbops.NetworkPolicyOwnerType,
	multicastNamespaceOwnerType:         ovsdbops.MulticastNamespaceOwnerType,
	multicastClusterOwnerType:           ovsdbops.MulticastClusterOwnerType,
	netpolNodeOwnerType:                 ovsdbops.NetpolNodeOwnerType,
	netpolNamespaceOwnerType:            ovsdbops.NetpolNamespaceOwnerType,
	udnIsolationOwnerType:               ovsdbops.UDNIsolationOwnerType,
	aclActionAllow:                      nbdb.ACLActionAllow,
	aclActionAllowRelated:               nbdb.ACLActionAllowRelated,
	aclActionAllowStateless:             nbdb.ACLActionAllowStateless,
	aclActionDrop:                       nbdb.ACLActionDrop,
	aclActionPass:                       nbdb.ACLActionPass,
	aclActionReject:                     nbdb.ACLActionReject,
}

// Protects from potential future renaming in ovn/ovs constants, since all constants are duplicated here
func TestConstantsMatch(t *testing.T) {
	for k, v := range mapping {
		if k != v {
			t.Fatalf("Constant %s does not match %s", k, v)
		}
	}
}
