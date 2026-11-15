//go:build wasip1

package main

import (
	"fmt"

	eb "github.com/StoneTrench/go-mc-clone/bindings/out/go/guest"
)

func RegBlock(name string, r, g, b, a uint8, sld eb.BlockKind) {
	namespace := eb.RegAddBlock(eb.BlockType{
		Name:  name,
		Color: eb.ColorU8{R: r, G: g, B: b, A: a},
		Kind:  sld,
	})
	eb.LogInfo(fmt.Sprintf("Registered block: %s", namespace))
}

func init() {
	eb.LogInfo("Module initialized!")

	err := eb.AppendLang("en-us")
	if len(err) > 0 {
		eb.LogError(err)
	}
	RegBlock("air", 255, 0, 255, 255, eb.E_Air)
	RegBlock("stone", 128, 128, 128, 255, eb.E_Solid)
	RegBlock("dirt", 255, 128, 128, 255, eb.E_Solid)
	RegBlock("grass", 128, 255, 128, 255, eb.E_Solid)

	RegBlock("water", 51, 102, 204, 180, eb.E_Liquid)
	RegBlock("sand", 219, 205, 135, 255, eb.E_Solid)
	RegBlock("snow", 240, 245, 255, 255, eb.E_Solid)
	RegBlock("iron_ore", 168, 122, 90, 255, eb.E_Solid)
	RegBlock("log", 102, 71, 43, 255, eb.E_Solid)
	RegBlock("leaves", 58, 122, 47, 220, eb.E_Solid)
}

func main() {}
