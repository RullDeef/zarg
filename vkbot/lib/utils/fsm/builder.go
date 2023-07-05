package fsm

import "log"

type FSMBuilder struct {
	states       map[int]any
	transitions  []transition
	initialState int
}

func New() *FSMBuilder {
	return &FSMBuilder{
		states:       make(map[int]any),
		transitions:  nil,
		initialState: -1,
	}
}

func (b *FSMBuilder) WithState(id int, state any) *FSMBuilder {
	if _, ok := b.states[id]; ok {
		log.Panicf("state with id=%d already exists", id)
	}

	b.states[id] = state
	return b
}

func (b *FSMBuilder) WithInitialState(id int, state any) *FSMBuilder {
	if b.initialState != -1 {
		log.Panic("initial state already set")
	}
	b.initialState = id
	return b.WithState(id, state)
}

func (b *FSMBuilder) WithTransition(fromID, toID int, fun TransitionFn) *FSMBuilder {
	for _, tr := range b.transitions {
		if tr.stateFrom == fromID && tr.stateTo == toID {
			log.Panicf("duplicate of transition #%d -> #%d", fromID, toID)
		}
	}
	tr := transition{
		stateFrom: fromID,
		stateTo:   toID,
		fun:       fun,
	}
	b.transitions = append(b.transitions, tr)
	return b
}

func (b *FSMBuilder) Build() *FSM {
	return &FSM{
		states:      b.states,
		transitions: b.transitions,
		currStateId: b.initialState,
	}
}
