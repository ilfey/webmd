package pages

import (
	"github.com/ilfey/webmd/internal/components"
	"github.com/ilfey/webmd/internal/fstree"
	"github.com/kyoto-framework/kyoto/v2"
)

type PDirState struct {
	*fstree.Dir
	DirMenu *kyoto.ComponentF[*components.CDirMenuState]
}

func PDir(dir *fstree.Dir) (kyoto.Component[*PDirState], error) {
	return func(ctx *kyoto.Context) *PDirState {
		kyoto.Template(ctx, "dir.html")

		return &PDirState{
			DirMenu: kyoto.Use(ctx, components.CDirMenu(&components.CDirMenuArgs{
				Dir: dir,
			})),
			Dir: dir,
		}
	}, nil
}
