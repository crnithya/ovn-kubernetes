package sampledecoder

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ovsdbops "github.com/ovn-org/ovn-kubernetes/go-controller/pkg/libovsdb/ops/ovsdb"
	libovsdbutil "github.com/ovn-org/ovn-kubernetes/go-controller/pkg/libovsdb/util"
	"github.com/ovn-org/ovn-kubernetes/go-controller/pkg/nbdb"
)

func TestCreateOrUpdateACL(t *testing.T) {
	event, err := newACLEvent(&nbdb.ACL{
		Action: nbdb.ACLActionAllow,
		ExternalIDs: map[string]string{
			ovsdbops.OwnerTypeKey.String():       ovsdbops.NetworkPolicyOwnerType,
			ovsdbops.ObjectNameKey.String():      "foo",
			ovsdbops.PolicyDirectionKey.String(): string(libovsdbutil.ACLIngress),
		},
	})
	require.ErrorContains(t, err, "expected format namespace:name for Object Name, but found: foo")
	assert.Nil(t, event)

	event, err = newACLEvent(&nbdb.ACL{
		Action: nbdb.ACLActionAllow,
		ExternalIDs: map[string]string{
			ovsdbops.OwnerTypeKey.String():       ovsdbops.NetworkPolicyOwnerType,
			ovsdbops.ObjectNameKey.String():      "bar:foo",
			ovsdbops.PolicyDirectionKey.String(): string(libovsdbutil.ACLIngress),
		},
	})
	require.NoError(t, err)
	assert.Equal(t, "Allowed by network policy foo in namespace bar, direction Ingress", event.String())

	event, err = newACLEvent(&nbdb.ACL{
		Action: nbdb.ACLActionAllow,
		ExternalIDs: map[string]string{
			ovsdbops.OwnerTypeKey.String():       ovsdbops.AdminNetworkPolicyOwnerType,
			ovsdbops.ObjectNameKey.String():      "foo",
			ovsdbops.PolicyDirectionKey.String(): string(libovsdbutil.ACLIngress),
		},
	})
	require.NoError(t, err)
	assert.Equal(t, "Allowed by admin network policy foo, direction Ingress", event.String())

	event, err = newACLEvent(&nbdb.ACL{
		Action: nbdb.ACLActionAllow,
		ExternalIDs: map[string]string{
			ovsdbops.OwnerTypeKey.String():  ovsdbops.EgressFirewallOwnerType,
			ovsdbops.ObjectNameKey.String(): "foo",
		},
	})
	require.NoError(t, err)
	assert.Equal(t, "Allowed by egress firewall in namespace foo", event.String())
	assert.Equal(t, "Egress", event.Direction)

	event, err = newACLEvent(&nbdb.ACL{
		Action: nbdb.ACLActionAllow,
		ExternalIDs: map[string]string{
			ovsdbops.OwnerTypeKey.String(): ovsdbops.NetpolNodeOwnerType,
		},
	})
	require.NoError(t, err)
	assert.Equal(t, "Allowed by default allow from local node policy, direction Ingress", event.String())
	assert.Equal(t, "Ingress", event.Direction)
}
