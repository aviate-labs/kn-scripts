package kn_scripts_test

import (
	"fmt"
	"github.com/aviate-labs/agent-go"
	"github.com/aviate-labs/agent-go/candid/idl"
	"github.com/aviate-labs/agent-go/ic"
	"github.com/aviate-labs/agent-go/ic/governance"
	"github.com/aviate-labs/agent-go/ic/registry"
	"github.com/aviate-labs/agent-go/principal"
	"testing"
)

const (
	Proposal133082 = 133082
)

// https://github.com/dfinity/ic/blob/master/rs/registry/admin/src/main.rs#L1946
type ProposeToUpdateNodeOperatorConfig struct {
	NodeOperatorID  *principal.Principal `ic:"node_operator_id,omitempty"`
	NodeAllowance   *uint64              `ic:"node_allowance,omitempty"`
	DCID            *string              `ic:"dc_id,omitempty"`
	RewardableNodes []struct {
		Field0 string `ic:"0,omitempty"`
		Field1 uint32 `ic:"1,omitempty"`
	} `ic:"rewardable_nodes,omitempty"`
	NodeProviderID *principal.Principal `ic:"node_provider_id,omitempty"`
	Ipv6           *string              `ic:"ipv6,omitempty"`
	SetIpv6ToNone  *bool                `ic:"set_ipv6_to_none,omitempty"`
}

func Test133082(t *testing.T) {
	govAgent, err := governance.NewAgent(ic.GOVERNANCE_PRINCIPAL, agent.DefaultConfig)
	if err != nil {
		t.Fatal(err)
	}
	info, err := govAgent.GetProposalInfo(Proposal133082)
	if err != nil {
		t.Fatal(err)
	}
	action := (*info).Proposal.Action.ExecuteNnsFunction
	if action == nil {
		t.Fatal("no action found")
	}
	if action.NnsFunction != 16 {
		t.Fatal("wrong function")
	}
	var a ProposeToUpdateNodeOperatorConfig
	if err := idl.Unmarshal(action.Payload, []any{&a}); err != nil {
		t.Fatal(err)
	}
	nodeOperator := *a.NodeOperatorID
	rewardableNodes := a.RewardableNodes[0]
	fmt.Println(nodeOperator, rewardableNodes)

	regAgent, err := registry.NewAgent(ic.REGISTRY_PRINCIPAL, agent.DefaultConfig)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := regAgent.GetNodeOperatorsAndDcsOfNodeProvider(principal.MustDecode("4jjya-hlyyc-s766p-fd6gr-d6tvv-vo3ah-j5ptx-i73gw-mwgyd-rw6w2-rae"))
	if err != nil {
		t.Fatal(err)
	}

	var total uint32
	totalRewardableNodes := make(map[string]uint32)
	for _, r := range *resp.Ok {
		operatorID := principal.Principal{Raw: r.Field1.NodeOperatorPrincipalId}.String()
		if _, ok := totalRewardableNodes[operatorID]; !ok {
			totalRewardableNodes[operatorID] = 0
		}
		for _, rn := range r.Field1.RewardableNodes {
			total += rn.Field1
			totalRewardableNodes[operatorID] += rn.Field1
		}
	}

	// Check if the node provider has only 10 rewardable nodes.
	if len(totalRewardableNodes) != 2 || total != 10 {
		t.Errorf("unexpected totalRewardableNodes: %v", totalRewardableNodes)
	}

	// Check if the proposal sets the rewardable nodes correctly.
	if rewardableNodes.Field1 != 5 {
		t.Errorf("unexpected rewardableNodes: %v", rewardableNodes)
	}
}
