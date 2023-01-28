package controllers

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"zarg/lib/model"
)

func BeginSession(chat chan model.Replica, out chan string) {
	defer close(out)
	players := gatherPlayers(chat, out)

	if players == nil {
		return
	}

	time.Sleep(5 * time.Second)

	// StartJourney
	pickWeapons(chat, out, players)

	out <- makePlayersInfo(players)

	// make interruptive commands available
	// chat1, chat2 := lib.DupReplicaChannel(chat)

}

func gatherPlayers(chat chan model.Replica, out chan string) []*model.Player {
	var players []*model.Player

	out <- "Начинается сбор людей и нелюдей для похода в данж!\nЧтобы участвовать, напиши \"Я\"."

	tenSecs := time.After(10 * time.Second)
	fifteenSecs := time.After(15 * time.Second)
outer:
	for {
		select {
		case replica := <-chat:
			if strings.ToLower(strings.TrimSpace(replica.Message())) == "я" {
				newPlayer := model.NewPlayer(replica.PeerID(), replica.UserName())
				// if player already joined
				joined := false
				for _, player := range players {
					if player.UserID() == newPlayer.UserID() {
						joined = true
					}
				}
				if joined {
					out <- fmt.Sprintf("%s уже в списке!", newPlayer.Name())
				} else {
					players = append(players, newPlayer)
					out <- fmt.Sprintf("%s участвует в походе!", newPlayer.Name())
				}
			} else if strings.ToLower(strings.TrimSpace(replica.Message())) == "не я" {
				for _, player := range players {
					if player.UserID() == replica.PeerID() {
						//TODO: remove self from list
						out <- fmt.Sprintf("Извини, %s, твоя судьба говорит иначе.", player.Name())
					}
				}
			}
		case <-tenSecs:
			out <- "5 секунд до окончания сбора!"
		case <-fifteenSecs:
			break outer
		}
	}

	if players == nil {
		out <- "Cбор окончен! В поход не идёт никто."
	} else {
		res := "Сбор окончен! В поход собрались:\n"
		for i, player := range players {
			res += fmt.Sprintf("  %d. %s\n", i+1, player.Name())
		}
		out <- res
	}

	return players
}

func pickWeapons(chat chan model.Replica, out chan string, players []*model.Player) {
	// generate weapons on start
	const totalWeapons = 6
	weapons := model.RandomWeapons(totalWeapons, 0, 6, 2)
	chosen := make([]bool, totalWeapons)

	ask := "Приключения ждут Вас, господа. А пока подготовьтесь к ним как следует. Выберите ваше оружие среди представленных:\n"
	for i, weapon := range weapons {
		ask += fmt.Sprintf("  %d) %s\n", i+1, weapon.Summary())
	}

	ask += "И поторопитесь, через 15 секунд выдвигаемся!"
	out <- ask

	lat := time.After(15 * time.Second)
	nChosen := 0 // amount of players that already picked weapon
outer:
	for {
		select {
		case replica := <-chat:
			if opt, err := strconv.ParseInt(replica.Message(), 10, 32); err == nil {
				if opt < 1 || opt > totalWeapons {
					continue
				}

				// find player that chooses
				for _, player := range players {
					if player.UserID() == replica.PeerID() {
						if !chosen[opt-1] {
							if player.Weapon == nil {
								nChosen += 1
							}
							player.Weapon = weapons[opt-1]
							out <- fmt.Sprintf("%s выбирает %s.", player.Name(), weapons[opt-1].SummaryShort())
						} else {
							out <- fmt.Sprintf("%s уже выбрал другой игрок!", weapons[opt-1].SummaryShort())
						}
						break
					}
				}

				if nChosen == len(players) {
					out <- "Все выбрали по оружию, отправляемся в данж!"
					break outer
				}
			}
		case <-lat:
			out <- "Выдвигаемся! А кто не успел схватиться за оружие будет сражаться кулаками!"
			for _, player := range players {
				if player.Weapon == nil {
					player.Weapon = model.FistsWeapon(5, 1)
				}
			}
			break outer
		}
	}
}

func makePlayersInfo(players []*model.Player) string {
	inf := "Статы игроков:\n\n"

	for _, player := range players {
		inf += fmt.Sprintf("%s: %s\nHP: %d\n\n", player.Name(), player.Weapon.SummaryShort(), player.Health)
	}

	return inf
}
