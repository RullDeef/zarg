package domain

import (
	"context"
)

// NPC - представление неигрового персонажа в целом (не в подземелье)
type NPC struct {
	Name      string     // имя персонажа
	Inventory *Inventory // инвернать персонажа
}

// HandleTradeRequest - обработчик запроса на куплю/продажу с игроком.
// Trader - игрок, запросивший обмен.
func (npc *NPC) HandleTradeRequest(_ context.Context, _ Trader, req *TradeRequest) error {
	// суть функции - проверить наличие запрошенных предметов в
	// инвентаре НИПа и адекватность предложенной цены

	reqItemsCost := 0 // реальная стоимость запрошенных вещей
	for _, item := range req.RequestItems {
		if !npc.Inventory.HasItem(item) {
			return ErrTradeUnfair
		}
		reqItemsCost += item.Cost
	}

	offItemsCost := 0 // реальная стоимость предлагаемых вещей
	for _, item := range req.OfferItems {
		offItemsCost += item.Cost
	}

	dealGain := reqItemsCost - offItemsCost // идеальный профит НИПа после сделки
	if req.OfferCost < dealGain {
		// предлагает недостаточно или просит слишком много монет.
		// Обновляем цену в условиях и сигнализируем о пересмотре условий
		req.OfferCost = dealGain
		return ErrTradeRequestChanged
	}

	// предложение в целом устраивает
	return nil
}
