// Code generated by cmd/lexgen (see Makefile's lexgen); DO NOT EDIT.

package atproto

// schema: com.atproto.admin.updateAccountSigningKey

import (
	"context"

	"github.com/bluesky-social/indigo/xrpc"
)

// AdminUpdateAccountSigningKey_Input is the input argument to a com.atproto.admin.updateAccountSigningKey call.
type AdminUpdateAccountSigningKey_Input struct {
	Did string `json:"did" cborgen:"did"`
	// signingKey: Did-key formatted public key
	SigningKey string `json:"signingKey" cborgen:"signingKey"`
}

// AdminUpdateAccountSigningKey calls the XRPC method "com.atproto.admin.updateAccountSigningKey".
func AdminUpdateAccountSigningKey(ctx context.Context, c *xrpc.Client, input *AdminUpdateAccountSigningKey_Input) error {
	if err := c.Do(ctx, xrpc.Procedure, "application/json", "com.atproto.admin.updateAccountSigningKey", nil, input, nil); err != nil {
		return err
	}

	return nil
}
