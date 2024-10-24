package model

import (
	"fmt"
	"pmon3/protos"
	"pmon3/utils/conv"
)

func (g *Group) ToProtobuf() *protos.Group {
	return &protos.Group{
		Id:   g.ID,
		Name: g.Name,
	}
}

func (g *Group) GetIdStr() string {
	return conv.Uint32ToStr(g.ID)
}

func (g *Group) RenderTable() []string {

	return []string{
		g.GetIdStr(),
		g.Name,
	}
}

//non-receiver methods begin

func GroupsArrayToProtobuf(groups []*Group) []*protos.Group {

	pgs := make([]*protos.Group, len(groups))

	for i := range groups {
		pgs[i] = groups[i].ToProtobuf()
	}

	return pgs
}

// protobuf method begin

func GroupFromProtobuf(p *protos.Group) *Group {
	return &Group{
		ID:   p.GetId(),
		Name: p.GetName(),
	}
}

func GroupsArrayFromProtobuf(groups []*protos.Group) []*Group {

	pgs := make([]*Group, len(groups))

	for i := range groups {
		pgs[i] = GroupFromProtobuf(groups[i])
	}

	return pgs
}

func (g *Group) Hash() string {
	return fmt.Sprintf("%s:%d", g.Name, g.ID)
}
