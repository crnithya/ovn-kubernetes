package util

import (
	"k8s.io/klog/v2"

	libovsdbclient "github.com/ovn-org/libovsdb/client"

	ovnops "github.com/ovn-org/ovn-kubernetes/go-controller/pkg/libovsdb/ops/ovn"
	ovsdbops "github.com/ovn-org/ovn-kubernetes/go-controller/pkg/libovsdb/ops/ovsdb"
	"github.com/ovn-org/ovn-kubernetes/go-controller/pkg/nbdb"
)

// GetACLCount returns the number of ACLs owned by idsType/controllerName
func GetACLCount(nbClient libovsdbclient.Client, idsType *ovsdbops.ObjectIDsType, controllerName string) int {
	predicateIDs := ovsdbops.NewDbObjectIDs(idsType, controllerName, nil)
	p := ovsdbops.GetPredicate[*nbdb.ACL](predicateIDs, nil)
	ACLs, err := ovnops.FindACLsWithPredicate(nbClient, p)
	if err != nil {
		klog.Warningf("Cannot find ACLs: %v; Resetting metrics...", err)
		return 0
	}
	return len(ACLs)
}

// GetAddressSetCount returns the number of AddressSets owned by idsType/controllerName
func GetAddressSetCount(nbClient libovsdbclient.Client, idsType *ovsdbops.ObjectIDsType, controllerName string) int {
	predicateIDs := ovsdbops.NewDbObjectIDs(idsType, controllerName, nil)
	p := ovsdbops.GetPredicate[*nbdb.AddressSet](predicateIDs, nil)
	ASes, err := ovnops.FindAddressSetsWithPredicate(nbClient, p)
	if err != nil {
		klog.Warningf("Cannot find AddressSets: %v; Resetting metrics...", err)
		return 0
	}
	return len(ASes)
}
