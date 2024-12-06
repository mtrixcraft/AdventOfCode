package main

import (
	fileparser "aoc/utils"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

var RULE_DIRECTIONS = [][2]int{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}
var symbolDIRECTIONS = []string{"|", "-", "|", "-"}

const symbolGuard = '^'
const symbolNewGuard = '^'
const symbolEmptyCell = '.'
const symbolObstruction = '#'

type Coordinates struct{ X, Y int }
type Guard struct {
	LocationStart     Coordinates
	LocationCurrent   Coordinates
	MovementDirection int
	MovementOffset    Coordinates
	LeftTheArea       bool
	StepSymbol        string
	StepsCount        int
}
type Map [][]string
type TheWatcher struct {
	GUARD             Guard
	MAP               Map
	SHOW_THE_TIMELINE bool
	SPEED             float64
	PART              bool
}

func (TheWatcher *TheWatcher) stopTimeline() {
	time.Sleep(time.Duration(TheWatcher.SPEED * float64(time.Second)))
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func (TheWatcher *TheWatcher) mapPrint() {
	TheWatcher.stopTimeline()
	fmt.Printf("GAME IS RUNNING...\n")
	for _, line := range TheWatcher.MAP {
		fmt.Printf("%v\n", line)
	}
}

func (TheWatcher *TheWatcher) mapUpdate(symbolUpdate string) {
	mapSymbol := &TheWatcher.MAP[TheWatcher.GUARD.LocationCurrent.Y][TheWatcher.GUARD.LocationCurrent.X]

	if TheWatcher.PART {
		*mapSymbol = "X"
		return
	}

	if *mapSymbol == string(symbolEmptyCell) || *mapSymbol == string(symbolNewGuard) {
		*mapSymbol = symbolUpdate
		return
	}
	if *mapSymbol == "|" || *mapSymbol == "-" {
		*mapSymbol = "+"
	}
}

func (TheWatcher *TheWatcher) initTheWatcher(INPUT *[]string) {
	var newGuard Guard
	newGuard.LeftTheArea = false
	newGuard.StepsCount = 0
	newGuard.MovementDirection = 0
	newGuard.MovementOffset.X = RULE_DIRECTIONS[newGuard.MovementDirection][0]
	newGuard.MovementOffset.Y = RULE_DIRECTIONS[newGuard.MovementDirection][1]
	for y, line := range *INPUT {
		x := strings.IndexRune(line, symbolGuard)
		if x != -1 {
			newGuard.LocationStart.X = x
			newGuard.LocationStart.Y = y
			break
		}
	}

	newGuard.LocationCurrent.X = newGuard.LocationStart.X
	newGuard.LocationCurrent.Y = newGuard.LocationStart.Y
	var newMap = make([][]string, len(*INPUT))
	for y, line := range *INPUT {
		newMap[y] = append(newMap[y], strings.Split(line, "")...)
	}

	TheWatcher.GUARD = newGuard
	TheWatcher.MAP = newMap
	TheWatcher.SPEED = 0
}

func (TheWatcher *TheWatcher) isValidNextMove() bool {
	guardOld := TheWatcher.GUARD
	var guardNew Coordinates
	guardNew.X = guardOld.LocationCurrent.X + guardOld.MovementOffset.X
	guardNew.Y = guardOld.LocationCurrent.Y + guardOld.MovementOffset.Y
	mapX := len(TheWatcher.MAP[0])
	mapY := len(TheWatcher.MAP)

	if guardNew.X < 0 || guardNew.X == mapX || guardNew.Y < 0 || guardNew.Y == mapY {
		TheWatcher.GUARD.LeftTheArea = true
		return false
	}

	if TheWatcher.MAP[guardNew.Y][guardNew.X] == string(symbolObstruction) {
		return false
	}

	TheWatcher.GUARD.StepSymbol = symbolDIRECTIONS[TheWatcher.GUARD.MovementDirection]
	return true
}

func (TheWatcher *TheWatcher) moveGuardForward() {
	guard := &TheWatcher.GUARD
	TheWatcher.mapUpdate(TheWatcher.GUARD.StepSymbol)
	guard.LocationCurrent.X += guard.MovementOffset.X
	guard.LocationCurrent.Y += guard.MovementOffset.Y
	if TheWatcher.MAP[guard.LocationCurrent.Y][guard.LocationCurrent.X] == string(symbolEmptyCell) {
		guard.StepsCount++
	}
	TheWatcher.mapUpdate(string(symbolNewGuard))
}

func (TheWatcher *TheWatcher) moveGuardRight() {
	guard := &TheWatcher.GUARD
	nextMovementDirection := (guard.MovementDirection + 1) % 4
	guard.MovementDirection = nextMovementDirection
	guard.MovementOffset.X = RULE_DIRECTIONS[nextMovementDirection][0]
	guard.MovementOffset.Y = RULE_DIRECTIONS[nextMovementDirection][1]
	TheWatcher.mapUpdate("+")
}

func partOne(TheWatcher *TheWatcher) int {
	TheWatcher.stopTimeline()
	for {
		if TheWatcher.SHOW_THE_TIMELINE {
			TheWatcher.mapPrint()
		}
		if !TheWatcher.isValidNextMove() {
			if TheWatcher.GUARD.LeftTheArea {
				break
			}
			TheWatcher.moveGuardRight()
			continue
		}
		TheWatcher.moveGuardForward()
	}

	return TheWatcher.GUARD.StepsCount + 1
}

func partTwo(TheWatcher *TheWatcher) int {
	return 0;
}


func main() {
	input := fileparser.ReadFileLines("inputTest", false)
	var TheWatcher TheWatcher
	TheWatcher.initTheWatcher(&input)
	TheWatcher.PART = true
	TheWatcher.SHOW_THE_TIMELINE = true
	TheWatcher.SPEED = 0.1
	fmt.Printf("--- Day 6: Guard Gallivant ---\n")
	fmt.Printf("PART[1] %v\n", partOne(&TheWatcher))
	var nextPart string;
	fmt.Println("Type Y to start Part 2: ")
	fmt.Scanln(&nextPart)
	fmt.Println("You entered:", nextPart)
	if nextPart != "Y" {
		os.Exit(0)
	}
	TheWatcher.initTheWatcher(&input)
	TheWatcher.PART = false
	TheWatcher.SHOW_THE_TIMELINE = true
	TheWatcher.SPEED = 0.001
	fmt.Printf("PART[2] %v\n", partTwo(&TheWatcher))
}
