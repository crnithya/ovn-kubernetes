package ovn

import (
	libovsdbclient "github.com/ovn-org/libovsdb/client"
	"github.com/ovn-org/libovsdb/ovsdb"

	ovsdbops "github.com/ovn-org/ovn-kubernetes/go-controller/pkg/libovsdb/ops/ovsdb"
	"github.com/ovn-org/ovn-kubernetes/go-controller/pkg/nbdb"
)

type DHCPOptionsPredicate func(*nbdb.DHCPOptions) bool

// CreateOrUpdateDhcpOptionsOps will configure logical switch port DHCPv4Options and DHCPv6Options fields with
// options at dhcpv4Options and dhcpv6Options arguments and create/update DHCPOptions objects that matches the
// pv4 and pv6 predicates. The missing DHCP options will default to nil in the LSP attributes.
func CreateOrUpdateDhcpOptionsOps(nbClient libovsdbclient.Client, ops []ovsdb.Operation, lsp *nbdb.LogicalSwitchPort, dhcpIPv4Options, dhcpIPv6Options *nbdb.DHCPOptions) ([]ovsdb.Operation, error) {
	opModels := []ovsdbops.OperationModel{}
	if dhcpIPv4Options != nil {
		opModel := ovsdbops.OperationModel{
			Model:          dhcpIPv4Options,
			OnModelUpdates: ovsdbops.OnModelUpdatesAllNonDefault(),
			DoAfter:        func() { lsp.Dhcpv4Options = &dhcpIPv4Options.UUID },
			ErrNotFound:    false,
			BulkOp:         false,
		}
		opModels = append(opModels, opModel)
	}
	if dhcpIPv6Options != nil {
		opModel := ovsdbops.OperationModel{
			Model:          dhcpIPv6Options,
			OnModelUpdates: ovsdbops.OnModelUpdatesAllNonDefault(),
			DoAfter:        func() { lsp.Dhcpv6Options = &dhcpIPv6Options.UUID },
			ErrNotFound:    false,
			BulkOp:         false,
		}
		opModels = append(opModels, opModel)
	}
	opModels = append(opModels, ovsdbops.OperationModel{
		Model: lsp,
		OnModelUpdates: []interface{}{
			&lsp.Dhcpv4Options,
			&lsp.Dhcpv6Options,
		},
		ErrNotFound: true,
		BulkOp:      false,
	})

	m := ovsdbops.NewModelClient(nbClient)
	return m.CreateOrUpdateOps(ops, opModels...)
}

func CreateOrUpdateDhcpOptions(nbClient libovsdbclient.Client, lsp *nbdb.LogicalSwitchPort, dhcpIPv4Options, dhcpIPv6Options *nbdb.DHCPOptions) error {
	ops, err := CreateOrUpdateDhcpOptionsOps(nbClient, nil, lsp, dhcpIPv4Options, dhcpIPv6Options)
	if err != nil {
		return err
	}
	_, err = ovsdbops.TransactAndCheck(nbClient, ops)
	return err
}

func DeleteDHCPOptions(nbClient libovsdbclient.Client, dhcpOptions *nbdb.DHCPOptions) error {
	opModels := []ovsdbops.OperationModel{}
	opModel := ovsdbops.OperationModel{
		Model:       dhcpOptions,
		ErrNotFound: false,
		BulkOp:      true,
	}
	opModels = append(opModels, opModel)
	m := ovsdbops.NewModelClient(nbClient)
	return m.Delete(opModels...)

}

func DeleteDHCPOptionsWithPredicate(nbClient libovsdbclient.Client, p DHCPOptionsPredicate) error {
	opModels := []ovsdbops.OperationModel{}
	opModel := ovsdbops.OperationModel{
		Model:          &nbdb.DHCPOptions{},
		ModelPredicate: p,
		ErrNotFound:    false,
		BulkOp:         true,
	}
	opModels = append(opModels, opModel)
	m := ovsdbops.NewModelClient(nbClient)
	return m.Delete(opModels...)

}
