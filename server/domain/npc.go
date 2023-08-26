package domain

import (
	"context"
)

const NPCBaseMaxHealth = 100

// NPC - представление неигрового персонажа в целом (не в подземелье)
type NPC struct {
	Name      string     // имя персонажа
	Inventory *Inventory // инвернать персонажа (используется для торговли)

	CanSell bool // возможна ли покупка у данного персонажа
	CanBuy  bool // возможна ли продажа данному персонажу

	// GuildScore - счет, добавляемый к рангу гильдии для сюжетного прогресса.
	// Обычно имеет значение при спасении обычных крестьян
	GuildScore int
}

// HandleTradeRequest - обработчик запроса на куплю/продажу с игроком.
// Trader - игрок, запросивший обмен
//
// Возвращает:
//   - ErrTradeRequestChanged - если запрос был изменен
//   - ErrTradeCanceled - если НИП не может торговать
//   - ErrTradeUnfair - если запрошены предметы, которых нет у НИПа
//   - nil - если запрос устраивает НИПа
func (npc *NPC) HandleTradeRequest(_ context.Context, _ Trader, req *TradeRequest) error {
	// суть функции - проверить наличие запрошенных предметов в
	// инвентаре НИПа и адекватность предложенной цены

	if !npc.CanSell && !npc.CanBuy {
		return ErrTradeCanceled
	}

	requestChanged := false

	// если НИП не может продавать - убрать предметы из списка запрошенных для выкупа
	if !npc.CanSell && len(req.RequestItems) != 0 {
		req.RequestItems = nil
		requestChanged = true
	}

	// если НИП не может покупать - убрать предметы из списка предлагаемых для продажи
	if !npc.CanBuy && len(req.OfferItems) != 0 {
		req.OfferItems = nil
		requestChanged = true
	}

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
		requestChanged = true
	}

	if requestChanged {
		return ErrTradeRequestChanged
	} else {
		return nil // предложение в целом устраивает
	}
}

// LostNPC - НИП, встреченный в подземелье.
// Может присоединиться к команде при обнаружении
type LostNPC struct {
	Entity
	*NPC
}

// NewLostNPC - делает из НИПа потерянного НИПа для помещения его в подземелье
func NewLostNPC(npc *NPC) *LostNPC {
	return &LostNPC{
		Entity: NewEntityBase(
			npc.Name,
			NPCBaseMaxHealth,
			npc.Inventory,
			NewEffectGroup(NoopEffectPolicy()),
		),
		NPC: npc,
	}
}
