package stmgr

import (
	"context"

	"github.com/filecoin-project/go-lotus/chain/actors"
	"github.com/filecoin-project/go-lotus/chain/types"
	"github.com/filecoin-project/go-lotus/chain/vm"

	cid "github.com/ipfs/go-cid"
	"golang.org/x/xerrors"
)

func (sm *StateManager) CallRaw(ctx context.Context, msg *types.Message, bstate cid.Cid, bheight uint64) (*types.MessageReceipt, error) {
	vmi, err := vm.NewVM(bstate, bheight, actors.NetworkAddress, sm.cs)
	if err != nil {
		return nil, xerrors.Errorf("failed to set up vm: %w", err)
	}

	if msg.GasLimit == types.EmptyInt {
		msg.GasLimit = types.NewInt(10000000000)
	}
	if msg.GasPrice == types.EmptyInt {
		msg.GasPrice = types.NewInt(0)
	}
	if msg.Value == types.EmptyInt {
		msg.Value = types.NewInt(0)
	}

	fromActor, err := vmi.StateTree().GetActor(msg.From)
	if err != nil {
		return nil, err
	}

	msg.Nonce = fromActor.Nonce

	// TODO: maybe just use the invoker directly?
	ret, err := vmi.ApplyMessage(ctx, msg)
	if err != nil {
		return nil, xerrors.Errorf("apply message failed: %w", err)
	}

	if ret.ActorErr != nil {
		log.Warnf("chain call failed: %s", ret.ActorErr)
	}
	return &ret.MessageReceipt, nil

}

func (sm *StateManager) Call(ctx context.Context, msg *types.Message, ts *types.TipSet) (*types.MessageReceipt, error) {
	if ts == nil {
		ts = sm.cs.GetHeaviestTipSet()
	}

	state, err := sm.TipSetState(ts.Cids())
	if err != nil {
		return nil, err
	}

	return sm.CallRaw(ctx, msg, state, ts.Height())
}