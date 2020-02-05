package client

import (
	govclient "web3space/ethermint/components/cosmos-sdk/x/gov/client"
	"web3space/ethermint/components/cosmos-sdk/x/params/client/cli"
	"web3space/ethermint/components/cosmos-sdk/x/params/client/rest"
)

// param change proposal handler
var ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
