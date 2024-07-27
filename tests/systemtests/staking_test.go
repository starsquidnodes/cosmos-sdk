//go:build system_test

package systemtests

import (
	"fmt"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestStakeUnstake(t *testing.T) {
	// Scenario:
	// delegate tokens to validator
	// undelegate some tokens

	sut.ResetChain(t)

	cli := NewCLIWrapper(t, sut, verbose)

	// add genesis account with some tokens
	account1Addr := cli.AddKey("account1")
	sut.ModifyGenesisCLI(t,
		[]string{"genesis", "add-genesis-account", account1Addr, "10000000stake"},
	)

	sut.StartChain(t)

	// query validator address to delegate tokens
	rsp := cli.CustomQuery("q", "staking", "validators")
	valAddr := gjson.Get(rsp, "validators.#.operator_address").Array()[0].String()

	// stake tokens
	rsp = cli.Run("tx", "staking", "delegate", valAddr, "10000stake", "--from="+account1Addr, "--fees=1stake")
	RequireTxSuccess(t, rsp)

	t.Log(cli.QueryBalance(account1Addr, "stake"))
	assert.Equal(t, int64(9989999), cli.QueryBalance(account1Addr, "stake"))

	rsp = cli.CustomQuery("q", "staking", "delegation", account1Addr, valAddr)
	assert.Equal(t, "10000", gjson.Get(rsp, "delegation_response.balance.amount").String(), rsp)
	assert.Equal(t, "stake", gjson.Get(rsp, "delegation_response.balance.denom").String(), rsp)

	// unstake tokens
	rsp = cli.Run("tx", "staking", "unbond", valAddr, "5000stake", "--from="+account1Addr, "--fees=1stake")
	RequireTxSuccess(t, rsp)

	rsp = cli.CustomQuery("q", "staking", "delegation", account1Addr, valAddr)
	assert.Equal(t, "5000", gjson.Get(rsp, "delegation_response.balance.amount").String(), rsp)
	assert.Equal(t, "stake", gjson.Get(rsp, "delegation_response.balance.denom").String(), rsp)

	rsp = cli.CustomQuery("q", "staking", "unbonding-delegation", account1Addr, valAddr)
	assert.Equal(t, "5000", gjson.Get(rsp, "unbond.entries.#.balance").Array()[0].String(), rsp)
}

func TestReduceValidatorSet(t *testing.T) {
	// Scenario:
	// Validator set starts with 4 validators
	// Validator set is reduced to 3 validators

	sut.ResetChain(t)

	cli := NewCLIWrapper(t, sut, verbose)

	// add genesis account with some tokens
	account1Addr := cli.AddKey("account1")
	sut.ModifyGenesisCLI(t,
		[]string{"genesis", "add-genesis-account", account1Addr, "10000000stake"},
	)

	votingPeriod := 5 * time.Second // enough time to vote

	sut.ModifyGenesisJSON(t,
		SetStakingParams(t, 4),
		SetGovVotingPeriod(t, votingPeriod),
	)

	sut.StartChain(t)
	govAddr := sdk.AccAddress(address.Module("gov")).String()

	// gov proposal to reduce validator set to 3
	proposal := fmt.Sprintf(`
{
 "messages": [
  {
   "@type": "/cosmos.staking.v1beta1.MsgUpdateParams",
   "authority": "%q",
   "params": {
    "unbonding_time": "504h0m0s",
    "max_validators": 3, 
    "max_entries": 7,
    "bond_denom": "stake",
    "min_commission_rate": "0.000000000000000000",
    "key_rotation_fee": {
     "denom": "stake",
     "amount": "1000000"
    }
   }
  }
 ],
}`, govAddr)
	proposalID := cli.SubmitAndVoteGovProposal(proposal)
	t.Logf("current_height: %d\n", sut.currentHeight)
	raw := cli.CustomQuery("q", "gov", "proposal", proposalID)
	t.Log(raw)

	t.Logf("current_height: %d\n", sut.currentHeight)
	raw = cli.CustomQuery("q", "gov", "proposal", proposalID)
	proposalStatus := gjson.Get(raw, "proposal.status").String()
	require.Equal(t, "PROPOSAL_STATUS_PASSED", proposalStatus, raw) // PROPOSAL_STATUS_PASSED

	rsp := cli.CustomQuery("q", "staking", "params")
	fmt.Println(rsp)

	rsp2 := cli.CustomQuery("q", "staking", "validators")
	fmt.Println(rsp2)

	t.Fail()
}
