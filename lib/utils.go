package lib

import "zarg/lib/model"

func DupReplicaChannel(c chan model.Replica) (chan model.Replica, chan model.Replica) {
	ca := make(chan model.Replica)
	cb := make(chan model.Replica)

	go func(ca, cb chan model.Replica) {
		for {
			val, ok := <-c
			if !ok {
				break
			}
			ca <- val
			cb <- val
		}

		close(ca)
		close(cb)
	}(ca, cb)

	return ca, cb
}
