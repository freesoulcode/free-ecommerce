package id

import (
	"fmt"

	"github.com/bwmarrin/snowflake"
)

type SnowflakeGenerator struct {
	node *snowflake.Node
}

func NewSnowflakeGenerator(nodeID int64) (*SnowflakeGenerator, error) {
	node, err := snowflake.NewNode(nodeID)
	if err != nil {
		return nil, fmt.Errorf("new snowflake node: %w", err)
	}

	return &SnowflakeGenerator{node: node}, nil
}

func (g *SnowflakeGenerator) NextID() (int64, error) {
	return g.node.Generate().Int64(), nil
}
