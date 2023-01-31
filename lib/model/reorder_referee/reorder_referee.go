package reorder_referee

import (
	"log"
	"zarg/lib/model/player/squad"
)

type ReorderReferee struct {
	players *squad.PlayerSquad
	order   []int
}

func New(players *squad.PlayerSquad) *ReorderReferee {
	return &ReorderReferee{
		players: players,
		order:   nil,
	}
}

// returns true if state changed
func (r *ReorderReferee) VoteStarter(id int) bool {
	if r.canVote(id) {
		r.order = nil
		r.order = append(r.order, id)
		return true
	}
	return false
}

// returns true if state changed
func (r *ReorderReferee) VoteNext(id int) bool {
	if (r.order != nil && r.order[len(r.order)-1] == id) || !r.canVote(id) {
		return false
	}
	for i := 0; i < len(r.order); i += 1 {
		if r.order[i] == id {
			r.order = append(r.order[:i], r.order[i+1:]...)
			break
		}
	}
	r.order = append(r.order, id)
	return true
}

func (r *ReorderReferee) Completed() bool {
	return len(r.order) == r.players.LenAlive()
}

func (r *ReorderReferee) Apply() {
	if !r.Completed() {
		log.Panicf("reordering is not completed! ordering: %v", r.order)
	}

	r.players.SetOrdering(r.order)
}

func (r *ReorderReferee) OrderingInfo() string {
	res := ""
	for _, id := range r.order {
		p := r.players.GetByID(id)
		res += p.FullName() + " -> "
	}
	return res + "..."
}

func (r ReorderReferee) canVote(id int) bool {
	p := r.players.GetByID(id)
	return p != nil && p.Alive()
}
