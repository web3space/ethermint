package client

import (
	"web3space/ethermint/components/cosmos-sdk/x/distribution/client/cli"
	"web3space/ethermint/components/cosmos-sdk/x/distribution/client/rest"
	govclient "web3space/ethermint/components/cosmos-sdk/x/gov/client"
)

// param change proposal handler
var (
	ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
)
