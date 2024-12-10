package ovn

import (
	"hash/fnv"

	"golang.org/x/net/context"

	libovsdbclient "github.com/ovn-org/libovsdb/client"
	"github.com/ovn-org/libovsdb/model"
	"github.com/ovn-org/libovsdb/ovsdb"

	"github.com/ovn-org/ovn-kubernetes/go-controller/pkg/config"
	ovsdbops "github.com/ovn-org/ovn-kubernetes/go-controller/pkg/libovsdb/ops/ovsdb"
	"github.com/ovn-org/ovn-kubernetes/go-controller/pkg/nbdb"
)

func CreateOrUpdateSampleCollector(nbClient libovsdbclient.Client, collector *nbdb.SampleCollector) error {
	opModel := ovsdbops.OperationModel{
		Model:          collector,
		OnModelUpdates: ovsdbops.OnModelUpdatesAllNonDefault(),
		ErrNotFound:    false,
		BulkOp:         false,
	}

	m := ovsdbops.NewModelClient(nbClient)
	_, err := m.CreateOrUpdate(opModel)
	return err
}

func UpdateSampleCollectorExternalIDs(nbClient libovsdbclient.Client, collector *nbdb.SampleCollector) error {
	opModel := ovsdbops.OperationModel{
		Model:          collector,
		OnModelUpdates: []interface{}{&collector.ExternalIDs},
		ErrNotFound:    true,
		BulkOp:         false,
	}

	m := ovsdbops.NewModelClient(nbClient)
	_, err := m.CreateOrUpdate(opModel)
	return err
}

func DeleteSampleCollector(nbClient libovsdbclient.Client, collector *nbdb.SampleCollector) error {
	opModel := ovsdbops.OperationModel{
		Model:       collector,
		ErrNotFound: false,
		BulkOp:      false,
	}
	m := ovsdbops.NewModelClient(nbClient)
	return m.Delete(opModel)
}

func DeleteSampleCollectorWithPredicate(nbClient libovsdbclient.Client, p func(collector *nbdb.SampleCollector) bool) error {
	opModel := ovsdbops.OperationModel{
		Model:          &nbdb.SampleCollector{},
		ModelPredicate: p,
		ErrNotFound:    false,
		BulkOp:         true,
	}
	m := ovsdbops.NewModelClient(nbClient)
	return m.Delete(opModel)
}

func FindSampleCollectorWithPredicate(nbClient libovsdbclient.Client, p func(*nbdb.SampleCollector) bool) ([]*nbdb.SampleCollector, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Default.OVSDBTxnTimeout)
	defer cancel()
	collectors := []*nbdb.SampleCollector{}
	err := nbClient.WhereCache(p).List(ctx, &collectors)
	return collectors, err
}

func ListSampleCollectors(nbClient libovsdbclient.Client) ([]*nbdb.SampleCollector, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Default.OVSDBTxnTimeout)
	defer cancel()
	collectors := []*nbdb.SampleCollector{}
	err := nbClient.List(ctx, &collectors)
	return collectors, err
}

func CreateOrUpdateSamplingAppsOps(nbClient libovsdbclient.Client, ops []ovsdb.Operation, samplingApps ...*nbdb.SamplingApp) ([]ovsdb.Operation, error) {
	opModels := make([]ovsdbops.OperationModel, 0, len(samplingApps))
	for i := range samplingApps {
		// can't use i in the predicate, for loop replaces it in-memory
		samplingApp := samplingApps[i]
		opModel := ovsdbops.OperationModel{
			Model:          samplingApp,
			OnModelUpdates: ovsdbops.OnModelUpdatesAllNonDefault(),
			ErrNotFound:    false,
			BulkOp:         false,
		}
		opModels = append(opModels, opModel)
	}

	modelClient := ovsdbops.NewModelClient(nbClient)
	return modelClient.CreateOrUpdateOps(ops, opModels...)
}

func DeleteSamplingAppsWithPredicate(nbClient libovsdbclient.Client, p func(collector *nbdb.SamplingApp) bool) error {
	opModel := ovsdbops.OperationModel{
		Model:          &nbdb.SamplingApp{},
		ModelPredicate: p,
		ErrNotFound:    false,
		BulkOp:         true,
	}
	m := ovsdbops.NewModelClient(nbClient)
	return m.Delete(opModel)
}

func FindSample(nbClient libovsdbclient.Client, sampleMetadata int) (*nbdb.Sample, error) {
	sample := &nbdb.Sample{
		Metadata: sampleMetadata,
	}
	return GetSample(nbClient, sample)
}

func GetSample(nbClient libovsdbclient.Client, sample *nbdb.Sample) (*nbdb.Sample, error) {
	found := []*nbdb.Sample{}
	opModel := ovsdbops.OperationModel{
		Model:          sample,
		ExistingResult: &found,
		ErrNotFound:    true,
		BulkOp:         false,
	}
	modelClient := ovsdbops.NewModelClient(nbClient)
	err := modelClient.Lookup(opModel)
	if err != nil {
		return nil, err
	}
	return found[0], err
}

type SampleFeature = string

const (
	EgressFirewallSample     SampleFeature = "EgressFirewall"
	NetworkPolicySample      SampleFeature = "NetworkPolicy"
	AdminNetworkPolicySample SampleFeature = "AdminNetworkPolicy"
	MulticastSample          SampleFeature = "Multicast"
	UDNIsolationSample       SampleFeature = "UDNIsolation"
)

// SamplingConfig is used to configure sampling for different db objects.
type SamplingConfig struct {
	featureCollectors map[SampleFeature][]string
}

func NewSamplingConfig(featureCollectors map[SampleFeature][]string) *SamplingConfig {
	return &SamplingConfig{
		featureCollectors: featureCollectors,
	}
}

func addSample(c *SamplingConfig, opModels []ovsdbops.OperationModel, model model.Model) []ovsdbops.OperationModel {
	switch t := model.(type) {
	case *nbdb.ACL:
		return createOrUpdateSampleForACL(opModels, c, t)
	}
	return opModels
}

// createOrUpdateSampleForACL should be called before acl ovsdbops.OperationModel is appended to opModels.
func createOrUpdateSampleForACL(opModels []ovsdbops.OperationModel, c *SamplingConfig, acl *nbdb.ACL) []ovsdbops.OperationModel {
	if c == nil {
		acl.SampleEst = nil
		acl.SampleNew = nil
		return opModels
	}
	collectors := c.featureCollectors[getACLSampleFeature(acl)]
	if len(collectors) == 0 {
		acl.SampleEst = nil
		acl.SampleNew = nil
		return opModels
	}
	aclID := GetACLSampleID(acl)
	sample := &nbdb.Sample{
		Collectors: collectors,
		// 32 bits
		Metadata: int(aclID),
	}
	opModel := ovsdbops.OperationModel{
		Model: sample,
		DoAfter: func() {
			acl.SampleEst = &sample.UUID
			acl.SampleNew = &sample.UUID
		},
		OnModelUpdates: []interface{}{&sample.Collectors},
		ErrNotFound:    false,
		BulkOp:         false,
	}
	opModels = append(opModels, opModel)
	return opModels
}

func GetACLSampleID(acl *nbdb.ACL) uint32 {
	// primaryID is unique for each ACL, but established connections will keep sampleID that is set on
	// connection creation. Here is the situation we want to avoid:
	// 1. ACL1 is created with sampleID=1 (e.g. based on ANP namespace+name+...+rule index with action Allow)
	// 2. connection A is established with sampleID=1, sample is decoded to say "Allowed by ANP namespace+name"
	// 3. ACL1 is updated with sampleID=1 (e.g. now same rule in ANP says Deny, but PrimaryIDKey is the same)
	// 4. connection A still generates samples with sampleID=1, but now it is "Denied by ANP namespace+name"
	// In reality, connection A is still allowed, as existing connections are not affected by ANP updates.
	// To avoid this, we encode Match and Action to the sampleID, to ensure a new sampleID is assigned on Match or action change.
	// In that case stale sampleIDs will just report messages like "sampling for this connection was updated or deleted".
	primaryID := acl.ExternalIDs[ovsdbops.PrimaryIDKey.String()] + acl.Match + acl.Action
	h := fnv.New32a()
	h.Write([]byte(primaryID))
	return h.Sum32()
}

func getACLSampleFeature(acl *nbdb.ACL) SampleFeature {
	switch acl.ExternalIDs[ovsdbops.OwnerTypeKey.String()] {
	case ovsdbops.AdminNetworkPolicyOwnerType, ovsdbops.BaselineAdminNetworkPolicyOwnerType:
		return AdminNetworkPolicySample
	case ovsdbops.MulticastNamespaceOwnerType, ovsdbops.MulticastClusterOwnerType:
		return MulticastSample
	case ovsdbops.NetpolNodeOwnerType, ovsdbops.NetworkPolicyOwnerType, ovsdbops.NetpolNamespaceOwnerType:
		return NetworkPolicySample
	case ovsdbops.EgressFirewallOwnerType:
		return EgressFirewallSample
	case ovsdbops.UDNIsolationOwnerType:
		return UDNIsolationSample
	}
	return ""
}
