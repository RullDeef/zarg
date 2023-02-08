package fsm

import "log"

type FSM struct {
	states      map[int]any
	transitions []transition
	currStateId int
}

type TransitionFn func(fsm *FSM, fromID, toID int) bool

type transition struct {
	stateFrom int
	stateTo   int
	fun       TransitionFn
}

func (fsm *FSM) State(id int) any {
	return fsm.states[id]
}

func (fsm *FSM) TransitTo(stateID int) {
	for _, tr := range fsm.transitions {
		if tr.stateFrom == fsm.currStateId && tr.stateTo == stateID {
			prev := fsm.currStateId
			fsm.currStateId = stateID
			tr.fun(fsm, prev, stateID)
			return
		}
	}
	log.Panicf("transition from state #%d to state #%d not found", fsm.currStateId, stateID)
}
