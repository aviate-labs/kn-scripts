package kn_scripts_test

import (
	"fmt"
	"testing"

	"github.com/aviate-labs/agent-go"
	"github.com/aviate-labs/agent-go/candid/idl"
	"github.com/aviate-labs/agent-go/ic"
	"github.com/aviate-labs/agent-go/ic/governance"
	"github.com/aviate-labs/agent-go/ic/registry"
	"github.com/aviate-labs/agent-go/principal"
)

const (
	Proposal134333 = 134333
)

func Test134333(t *testing.T) {
	govAgent, err := governance.NewAgent(ic.GOVERNANCE_PRINCIPAL, agent.DefaultConfig)
	if err != nil {
		t.Fatal(err)
	}
	info, err := govAgent.GetProposalInfo(Proposal134333)
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
	resp, err := regAgent.GetNodeOperatorsAndDcsOfNodeProvider(principal.MustDecode("mjnyf-lzqq6-s7fzb-62rqm-xzvge-5oa26-humwp-dvwxp-jxxkf-hoel7-fqe"))
	if err != nil {
		t.Fatal(err)
	}

	for _, resp := range *resp.Ok {
		fmt.Println(principal.Principal{Raw: resp.Field1.NodeOperatorPrincipalId})
	}
}
