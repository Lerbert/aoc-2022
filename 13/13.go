package main

import (
	util "aoc2022/aoc"
	inp "aoc2022/input"
	"encoding/json"
	"fmt"
	"log"
)

type packet []interface{}

func (p packet) rightOrder(other packet) int {
	pLen := len(p)
	oLen := len(other)
	for i := 0; ; i++ {
		if i >= pLen && i < oLen {
			return 1
		} else if i >= oLen && i < pLen {
			return -1
		} else if i >= pLen && i >= oLen {
			return 0
		}
		left := p[i]
		right := other[i]
		switch leftT := left.(type) {
		case float64:
			switch rightT := right.(type) {
			case float64:
				if leftT < rightT {
					return 1
				} else if leftT > rightT {
					return -1
				}
			case []interface{}:
				if res := (packet{leftT}.rightOrder(packet(rightT))); res != 0 {
					return res
				}
			default:
				panic("Unknown packet type")
			}
		case []interface{}:
			switch rightT := right.(type) {
			case float64:
				if res := packet(leftT).rightOrder(packet{rightT}); res != 0 {
					return res
				}
			case []interface{}:
				if res := packet(leftT).rightOrder(packet(rightT)); res != 0 {
					return res
				}
			default:
				panic("Unknown packet type")
			}
		default:
			panic("Unknown packet type")
		}
	}
}

func main() {
	lines := inp.ReadLines("input")
	lines = util.Filter(&lines, func(s string) bool { return s != "" })
	packets := make([]packet, len(lines))
	for i, l := range lines {
		if err := json.Unmarshal([]byte(l), &packets[i]); err != nil {
			log.Fatal("Could not parse JSON", err)
		}
	}

	indexSum := 0
	for i := 0; i < len(packets); i += 2 {
		if packets[i].rightOrder(packets[i+1]) > 0 {
			indexSum += i/2 + 1
		}
	}
	fmt.Printf("Part 1: %d\n", indexSum)
}
