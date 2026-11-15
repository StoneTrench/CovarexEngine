package game

import (
	"github.com/StoneTrench/CovarexEngine/src/engine/resources"
	"github.com/StoneTrench/CovarexEngine/src/game/world"
)

var RegistryBlocks = new(resources.Registry[world.BlockType, world.BlockId]).Init("blocks")
