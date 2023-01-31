package reorder_referee

import (
	"log"
	"zarg/lib/model/player"
	"zarg/lib/model/user"
)

type ReorderReferee struct {
	players *player.PlayerSquad
	order   []*user.User
}

func New(players *player.PlayerSquad) *ReorderReferee {
	return &ReorderReferee{
		players: players,
		order:   nil,
	}
}

// returns true if state changed
func (r *ReorderReferee) VoteStarter(u *user.User) bool {
	if r.canVote(u) {
		r.order = nil
		r.order = append(r.order, u)
		return true
	}
	return false
}

// returns true if state changed
func (r *ReorderReferee) VoteNext(u *user.User) bool {
	if (r.order != nil && r.order[len(r.order)-1] == u) || !r.canVote(u) {
		return false
	}
	for i := 0; i < len(r.order); i += 1 {
		if r.order[i] == u {
			r.order = append(r.order[:i], r.order[i+1:]...)
			break
		}
	}
	r.order = append(r.order, u)
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
	for _, u := range r.order {
		res += u.FullName() + " -> "
	}
	return res + "..."
}

func (r *ReorderReferee) canVote(u *user.User) bool {
	p := r.players.GetByUser(u)
	return p != nil && p.Alive()
}
