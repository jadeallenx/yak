package yak

import (
	"fmt"
	"strings"
)

func aType2String(t Assettype) string {
	switch t {
	case Image:
		return "Image"
	case Link:
		return "External Link"
	case CSS:
		return "Stylesheet"
	case Script:
		return "Javascript"
	}
	return "Unknown"
}

func PrettyPrint(p *Page, generation int) {
	fmt.Println(strings.Repeat(" ", generation) + strings.Repeat(">", generation) + " " + p.Loc.String())
	fmt.Println(strings.Repeat("  ", generation) + "Assets:")
	if len(p.Assets) > 0 {
		for _, a := range p.Assets {
			if a.Loc != nil {
				fmt.Println(strings.Repeat("    ", generation) + aType2String(a.Atype) + " (" + a.Loc.String() + ")")
			} else {
				fmt.Println(strings.Repeat("    ", generation) + aType2String(a.Atype))
			}
		}
	} else {
		fmt.Println(strings.Repeat("    ", generation) + "None")
	}
	generation++
	for _, ch := range p.Children {
		PrettyPrint(ch, generation)
	}
}
