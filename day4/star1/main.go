package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	events, err := readEventsFromFile("input.txt")
	if err != nil {
		panic(err)
	}

	sortEventSlice(events)
	setEventGuardIDs(events)
	//fmt.Println(events)
	guards, err := eventListToGuardList(events)
	if err != nil {
		panic(err)
	}

	msind := mostSleepIndex(guards)
	minute := findMostCommonSleepMinute(guards[msind])
	fmt.Println("Minute", minute)
	fmt.Println("Guard #", msind)
	fmt.Println("Answer: ", minute*msind)
}

const (
	EVENT_SHIFT = iota
	EVENT_SLEEP
	EVENT_WAKE_UP
)

const (
	UNKNOWN_GUARD_ID = -1
)

type event struct {
	year      int
	month     int
	day       int
	hour      int
	minute    int
	eventType int
	guardID   int
}

func readEventsFromFile(path string) ([]event, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileReader := bufio.NewReader(file)

	var events []event = make([]event, 0)

	for line, err := fileReader.ReadString('\n'); err == nil; line, err = fileReader.ReadString('\n') {
		if len(line) < 18 {
			break
		}
		tempStringArr := strings.Split(line, "] ")

		if len(tempStringArr) != 2 {
			break
		}

		lTime := strings.Replace(tempStringArr[0], "[", "", -1)
		lEvent := tempStringArr[1]

		var (
			currentEvent event
		)
		if lEvent == "falls asleep\n" {
			currentEvent.eventType = EVENT_SLEEP
			currentEvent.guardID = UNKNOWN_GUARD_ID
		} else if lEvent == "wakes up\n" {
			currentEvent.eventType = EVENT_WAKE_UP
			currentEvent.guardID = UNKNOWN_GUARD_ID
		} else if strings.Contains(lEvent, "Guard #") {
			guardID := strings.Replace(strings.Replace(lEvent, "Guard #", "", -1), " begins shift\n", "", -1)
			currentEvent.guardID, err = strconv.Atoi(guardID)
			if err != nil {
				return nil, err
			}
			currentEvent.eventType = EVENT_SHIFT
		} else {
			return nil, errUnknownEvent{evt: lEvent}
		}

		tempStringArr = strings.Split(lTime, " ")
		if len(tempStringArr) != 2 {
			break
		}

		lDateSplit := strings.Split(tempStringArr[0], "-")
		lHourMinute := strings.Split(tempStringArr[1], ":")

		if len(lDateSplit) != 3 || len(lHourMinute) != 2 {
			break
		}

		lYear, err := strconv.Atoi(lDateSplit[0])
		lMonth, err := strconv.Atoi(lDateSplit[1])
		lDay, err := strconv.Atoi(lDateSplit[2])
		lHour, err := strconv.Atoi(lHourMinute[0])
		lMinute, err := strconv.Atoi(lHourMinute[1])
		if err != nil {
			return nil, err
		}
		currentEvent.year = lYear
		currentEvent.month = lMonth
		currentEvent.day = lDay
		currentEvent.hour = lHour
		currentEvent.minute = lMinute

		events = append(events, currentEvent)
	}

	return events, err
}

type errUnknownEvent struct {
	evt string
}

func (e errUnknownEvent) Error() string {
	return "Couldn't parse event: \"" + e.evt + "\""
}

func sortEventSlice(evts []event) {
	sort.Slice(evts, func(i, j int) bool {
		if evts[i].year < evts[j].year {
			return true
		} else if evts[i].year > evts[j].year {
			return false
		}
		if evts[i].month < evts[j].month {
			return true
		} else if evts[i].month > evts[j].month {
			return false
		}
		if evts[i].day < evts[j].day {
			return true
		} else if evts[i].day > evts[j].day {
			return false
		}
		if evts[i].hour < evts[j].hour {
			return true
		} else if evts[i].hour > evts[j].hour {
			return false
		}
		if evts[i].minute < evts[j].minute {
			return true
		} else if evts[i].minute > evts[j].minute {
			return false
		}

		return false
	})
}

func setEventGuardIDs(evts []event) error {
	if len(evts) == 0 {
		return emptySliceErr{}
	}

	var (
		currentID int = evts[0].guardID
		err       error
	)
	if currentID == UNKNOWN_GUARD_ID {
		err = unknownIDErr{}
	}

	for i := range evts {
		if evts[i].guardID == UNKNOWN_GUARD_ID {
			evts[i].guardID = currentID
		} else {
			currentID = evts[i].guardID
		}
	}
	return err
}

type emptySliceErr struct{}

func (e emptySliceErr) Error() string {
	return "The received slice is empty"
}

type unknownIDErr struct{}

func (e unknownIDErr) Error() string {
	return "The guard ID for the first event is unknown, trying to set the rest anyway"
}

type guard struct {
	sleepAmount [60]int
}

// Assumes the sleeping is always between 00:00 and 00:59
func eventListToGuardList(evts []event) (map[int]guard, error) {
	if len(evts) == 0 {
		return nil, emptySliceErr{}
	}

	guardList := make(map[int]guard, 0)
	currentGuard := evts[0].guardID
	var err error

	type dateAndTime struct {
		year   int
		month  int
		day    int
		minute int
	}
	var startSleepTime dateAndTime

	if currentGuard == UNKNOWN_GUARD_ID {
		err = unknownIDErr{}
	}

	_ = currentGuard

	for i := range evts {
		evt := evts[i]
		switch evt.eventType {
		case EVENT_WAKE_UP:
			if evt.year == startSleepTime.year &&
				evt.month == startSleepTime.month &&
				evt.day == startSleepTime.day {
				if evt.minute-startSleepTime.minute < 0 {
					return nil, notSortedErr{}
				}
				_, ok := guardList[evt.guardID]
				if !ok {
					guardList[evt.guardID] = guard{}
				}
				for minute := startSleepTime.minute; minute < evt.minute; minute++ {
					temp := guardList[evt.guardID].sleepAmount
					temp[minute]++
					guardList[evt.guardID] = guard{
						sleepAmount: temp,
					}
				}
			}
		case EVENT_SLEEP:
			startSleepTime = dateAndTime{
				year:   evt.year,
				month:  evt.month,
				day:    evt.day,
				minute: evt.minute,
			}
		}
	}

	return guardList, err
}

type notSortedErr struct{}

func (e notSortedErr) Error() string {
	return "The input probaly wasn't sorted"
}

func mostSleepIndex(guards map[int]guard) int {
	var maxIndex, maxValue int
	for i := range guards {
		var sum int
		for j := range guards[i].sleepAmount {
			sum += guards[i].sleepAmount[j]
		}
		if sum > maxValue {
			maxIndex = i
			maxValue = sum
		}
	}
	return maxIndex
}

func findMostCommonSleepMinute(g guard) int {
	var maxIndex, maxValue int
	for i := range g.sleepAmount {
		if g.sleepAmount[i] > maxValue {
			maxIndex = i
			maxValue = g.sleepAmount[i]
		}
	}
	return maxIndex
}
