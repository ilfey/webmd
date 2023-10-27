package components

import (
	"github.com/ilfey/webmd/internal/fstree"
	"github.com/kyoto-framework/kyoto/v2"
)

type CDirMenuArgs struct {
	*fstree.Dir
}

type CDirMenuState struct {
	Args *CDirMenuArgs
}

func CDirMenu(args *CDirMenuArgs) kyoto.Component[*CDirMenuState] {
	return func(ctx *kyoto.Context) *CDirMenuState {
		state := &CDirMenuState{
			Args: args,
		}

		kyoto.ActionPreload(ctx, state)

		return state
	}

}
