package fabric

import (
	"fmt"
	"io/ioutil"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

type Client struct {
	// Profile is the file path to the JSON formatted connection profile as
	// exported from the IBM Blockchain Platform UI. When running in
	// production, it is recommended to store the profile contents in an
	// environment variable encoded in base64, and use an entrypoint script to
	// populate this file from the environment before startup. This strategy
	// combines the benefits of convenience in terms of directly using the file
	// exported from IBM's UI and the benefits of flexibility with regards to
	// the docker image not hardcoding any particular connection profile.
	Profile string

	// Channel is the channel ID that the client will connect to to invoke
	// smart contract functions. The channel ID can be found in the IBM UI's
	// channel tab.
	Channel string

	// User is the TODO: what is user vs Org?
	User string

	// Org is the TODO: what is Org vs user?
	Org string

	// Contract is the ID of the smart contract (aka chaincodeID) that the
	// client will invoke.
	Contract string

	channelClient *channel.Client
}

func (c *Client) Connect() error {
	b, err := ioutil.ReadFile(c.Profile)
	if err != nil {
		return err // TODO: wrap this error
	}
	sdk, err := fabsdk.New(config.FromRaw(b, "json"))
	if err != nil {
		return err // TODO: wrap this error
	}
	chClient, err := channel.New(sdk.ChannelContext(c.Channel, fabsdk.WithOrg(c.Org)))
	if err != nil {
		return err
	}
	c.channelClient = chClient
	return nil
}

// Invoke will connect to the and call the smart contract configured in the Client
// object with the function and args given. function is the name of an
// operation defined in the smart contract. args is a list of values to provide
// as arguments to the
// smart contract function.
// TODO: consider args helpers to avoid [][]byte
func (c *Client) Invoke(function string, args ...[]byte) ([]byte, error) {
	if c.channelClient == nil {
		c.Connect()
	}
	response, err := c.channelClient.Query(channel.Request{
		ChaincodeID: c.Contract,
		Fcn:         function,
		Args:        args,
	})
	if err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf("%#v", response)), nil
}
