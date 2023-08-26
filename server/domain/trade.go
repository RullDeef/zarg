package domain

import (
	"context"
	"errors"
)

var (
	// ErrTradeRequestChanged - запрос изменен и должен
	// быть пересмотрен и подтвержден другой стороной
	ErrTradeRequestChanged = errors.New("trade request changed")

	ErrTradeUnfair   = errors.New("trade unfair")   // невыгодная торговая сделка
	ErrTradeCanceled = errors.New("trade canceled") // обмен отменен одной из сторон
	ErrTradeExpired  = errors.New("trade expired")  // контекст обмена отменен
)

// Trader - абстракция сущности, осуществляющей торговлю (куплю или продажу).
//
// Торговля осуществляется как процесс последовательных вызовов метода
// HandleRequest у торгующих с передачей измененной структуры TradeRequest
// до тех пор, пока обе стороны не перестанут изменять запрос и договоряться.
//
// Опасно: можно первым же вызовом HandleTradeRequest передать запрос с двумя
// подтверждениями. Данный случай хоть и опасный, но не отслеживается
type Trader interface {
	// HandleTradeRequest - обработать запрос на обмен с другим Trader.
	// При изменении запроса необходимо вернуть ErrTradeRequestChanged,
	// чтобы дать возможность первой стороне подтвердить изменения
	HandleTradeRequest(context.Context, Trader, *TradeRequest) error

	// ApplyTradeResult - безусловно применить результат обмена
	ApplyTradeResult(TradeResult)
}

// TradeRequest - структура, описывающая предложение для обмена
type TradeRequest struct {
	OfferItems   []*PickableItem // предметы, предлагаемые для продажи
	RequestItems []*PickableItem // предметы, запрашиваемые для выкупа

	// OfferCost - сумма, которую запрашивающий готов отдать за совершение обмена.
	// Если сделка нацелена на продажу, данное поле может быть отрицательным, что
	// означает, что запрашивающий предлагает предметы в обмен на монеты
	OfferCost int
}

type TradeResult struct {
	GainItems []*PickableItem // предметы, полученные в результате обмена
	LostItems []*PickableItem // предметы, отданные в результате обмена

	// GainMoney - сумма, приобретенная в результате обмена.
	// Может быть отрицательной, что означает передачу монет
	GainMoney int
}

// Flip - меняет местами участников обмена
func (r *TradeRequest) Flip() {
	r.OfferItems, r.RequestItems = r.RequestItems, r.OfferItems
	r.OfferCost = -r.OfferCost
}

// Result - формирует результат обмена по запросу
func (r *TradeRequest) Result() TradeResult {
	return TradeResult{
		GainItems: r.RequestItems,
		LostItems: r.OfferItems,
		GainMoney: -r.OfferCost,
	}
}

// PerformTrade - функция, отвечающая за корректность последовательности
// вызовов методов обмена между двумя сторонами
func PerformTrade(ctx context.Context, requester, responder Trader, req TradeRequest) error {
	confirms := 0 // число подтверждений (для завершения обмена нужно 2)
	for {
		err := responder.HandleTradeRequest(ctx, requester, &req)
		switch err {
		case nil:
			confirms++
		case ErrTradeRequestChanged:
			confirms = 1 // сторона изменила условия, необходимо доп. подтверждение
		default:
			return err
		}

		if confirms == 2 {
			// применяем условия к обоим сторонам обмена
			responder.ApplyTradeResult(req.Result())
			req.Flip()
			requester.ApplyTradeResult(req.Result())
			return nil
		}

		// поменять стороны и условия обмена местами и продолжить
		requester, responder = responder, requester
		req.Flip()
	}
}
