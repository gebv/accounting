package irecharge

import (
	"context"
	"log"
	"sync"

	"github.com/gebv/acca/engine"
	"github.com/gebv/acca/engine/strategies"
	"github.com/gebv/acca/ffsm"
	"github.com/pkg/errors"
)

const nameStrategy strategies.InvStrategyName = "invoice_recharge_strategy"

func init() {
	s := &Strategy{
		s: make(ffsm.Stack),
	}
	s.load()
	strategies.RegInvoiceStrategy(s)
}

type Strategy struct {
	s        ffsm.Stack
	syncOnce sync.Once
}

func (s *Strategy) Name() strategies.InvStrategyName {
	return nameStrategy
}

func (s *Strategy) Dispatch(ctx context.Context, state ffsm.State, payload ffsm.Payload) error {
	invID, ok := payload.(int64)
	if !ok {
		return errors.New("bad_payload")
	}
	tx := strategies.GetTXContext(ctx)
	if tx == nil {
		return errors.New("Not reform tx.")
	}
	inv := engine.Invoice{InvoiceID: invID}
	if err := tx.Reload(&inv); err != nil {
		return errors.Wrap(err, "Failed reload invoice by ID.")
	}
	st := ffsm.State(inv.Status)
	fsm := ffsm.MachineFrom(s.s, &st)
	err := fsm.Dispatch(context.Background(), state, payload)
	if err != nil {
		return err
	}
	return nil
}

var _ strategies.InvStrategy = (*Strategy)(nil)

func (s *Strategy) load() {
	s.syncOnce.Do(func() {
		s.s.Add(
			ffsm.State(engine.DRAFT_I),
			ffsm.State(engine.AUTH_I),
			func(ctx context.Context, payload ffsm.Payload) (context context.Context, e error) {
				// TODO load invoice from DB
				invID, ok := payload.(int64)
				if !ok {
					log.Println("Invoice bad Payload: ", payload)
					return
				}
				tx := strategies.GetTXContext(ctx)
				if tx == nil {
					return ctx, errors.New("Not reform tx in context.")
				}
				inv := engine.Invoice{InvoiceID: invID}
				if err := tx.Reload(&inv); err != nil {
					return ctx, errors.Wrap(err, "Failed reload invoice by ID.")
				}
				if inv.Strategy != nameStrategy.String() {
					return ctx, errors.New("Invoice strategy not nameStrategy.")
				}
				if inv.Status != engine.DRAFT_I {
					return ctx, errors.New("Invoice status not draft.")
				}
				// Установить статус куда происходит переход
				ns := engine.AUTH_I
				inv.NextStatus = &ns
				if err := tx.Save(&inv); err != nil {
					return ctx, errors.Wrap(err, "Failed save invoice by ID.")
				}
				// TODO load transactions from DB
				next := true
				list, err := tx.SelectAllFrom(engine.TransactionTable, "WHERE invoice_id = $1", invID)
				if err != nil {
					return ctx, errors.Wrap(err, "Failed list transaction by invoice ID.")
				}
				for _, v := range list {
					tr := v.(*engine.Transaction)
					if !tr.Status.Match(engine.AUTH_TX) {
						next = false
						break
					}
				}
				if !next {
					return ctx, nil
				}
				// Установить статус после проделанных операций
				inv.Status = engine.AUTH_I
				inv.NextStatus = nil
				if err := tx.Save(&inv); err != nil {
					return ctx, errors.Wrap(err, "Failed save invoice by ID.")
				}
				return ctx, nil
			},
			"draft>auth",
		)
		s.s.Add(
			ffsm.State(engine.DRAFT_I),
			ffsm.State(engine.REJECTED_I),
			func(ctx context.Context, payload ffsm.Payload) (context context.Context, e error) {
				// TODO load invoice from DB
				invID, ok := payload.(int64)
				if !ok {
					log.Println("Invoice bad Payload: ", payload)
					return
				}
				tx := strategies.GetTXContext(ctx)
				if tx == nil {
					return ctx, errors.New("Not reform tx in context.")
				}
				inv := engine.Invoice{InvoiceID: invID}
				if err := tx.Reload(&inv); err != nil {
					return ctx, errors.Wrap(err, "Failed reload invoice by ID.")
				}
				if inv.Strategy != nameStrategy.String() {
					return ctx, errors.New("Invoice strategy not nameStrategy.")
				}
				if inv.Status != engine.DRAFT_I {
					return ctx, errors.New("Invoice status not draft.")
				}
				// Установить статус куда происходит переход
				ns := engine.REJECTED_I
				inv.NextStatus = &ns
				if err := tx.Save(&inv); err != nil {
					return ctx, errors.Wrap(err, "Failed save invoice by ID.")
				}
				// TODO load transactions from DB
				next := true
				list, err := tx.SelectAllFrom(engine.TransactionTable, "WHERE invoice_id = $1", invID)
				if err != nil {
					return ctx, errors.Wrap(err, "Failed list transaction by invoice ID.")
				}
				for _, v := range list {
					tr := v.(*engine.Transaction)
					if !tr.Status.Match(engine.REJECTED_TX) {
						next = false
						break
					}
				}
				if !next {
					return ctx, nil
				}
				// Установить статус после проделанных операций
				inv.Status = engine.REJECTED_I
				inv.NextStatus = nil
				if err := tx.Save(&inv); err != nil {
					return ctx, errors.Wrap(err, "Failed save invoice by ID.")
				}
				return ctx, nil
			},
			"draft>rejected",
		)
		s.s.Add(
			ffsm.State(engine.AUTH_I),
			ffsm.State(engine.WAIT_I),
			func(ctx context.Context, payload ffsm.Payload) (context context.Context, e error) {
				// TODO load invoice from DB
				invID, ok := payload.(int64)
				if !ok {
					log.Println("Invoice bad Payload: ", payload)
					return
				}
				tx := strategies.GetTXContext(ctx)
				if tx == nil {
					return ctx, errors.New("Not reform tx in context.")
				}
				inv := engine.Invoice{InvoiceID: invID}
				if err := tx.Reload(&inv); err != nil {
					return ctx, errors.Wrap(err, "Failed reload invoice by ID.")
				}
				if inv.Strategy != nameStrategy.String() {
					return ctx, errors.New("Invoice strategy not nameStrategy.")
				}
				if inv.Status != engine.AUTH_I {
					return ctx, errors.New("Invoice status not auth.")
				}
				// Установить статус куда происходит переход
				ns := engine.WAIT_I
				inv.NextStatus = &ns
				if err := tx.Save(&inv); err != nil {
					return ctx, errors.Wrap(err, "Failed save invoice by ID.")
				}
				// TODO load transactions from DB
				next := true
				list, err := tx.SelectAllFrom(engine.TransactionTable, "WHERE invoice_id = $1", invID)
				if err != nil {
					return ctx, errors.Wrap(err, "Failed list transaction by invoice ID.")
				}
				for _, v := range list {
					tr := v.(*engine.Transaction)
					if !tr.Status.Match(engine.HOLD_TX) {
						next = false
						break
					}
				}
				if !next {
					return ctx, nil
				}
				// Установить статус после проделанных операций
				inv.Status = engine.WAIT_I
				inv.NextStatus = nil
				if err := tx.Save(&inv); err != nil {
					return ctx, errors.Wrap(err, "Failed save invoice by ID.")
				}
				return ctx, nil
			},
			"auth>wait",
		)
		s.s.Add(
			ffsm.State(engine.AUTH_I),
			ffsm.State(engine.ACCEPTED_I),
			func(ctx context.Context, payload ffsm.Payload) (context context.Context, e error) {
				// TODO load invoice from DB
				invID, ok := payload.(int64)
				if !ok {
					log.Println("Invoice bad Payload: ", payload)
					return
				}
				tx := strategies.GetTXContext(ctx)
				if tx == nil {
					return ctx, errors.New("Not reform tx in context.")
				}
				inv := engine.Invoice{InvoiceID: invID}
				if err := tx.Reload(&inv); err != nil {
					return ctx, errors.Wrap(err, "Failed reload invoice by ID.")
				}
				if inv.Strategy != nameStrategy.String() {
					return ctx, errors.New("Invoice strategy not nameStrategy.")
				}
				if inv.Status != engine.AUTH_I {
					return ctx, errors.New("Invoice status not auth.")
				}
				// Установить статус куда происходит переход
				ns := engine.ACCEPTED_I
				inv.NextStatus = &ns
				if err := tx.Save(&inv); err != nil {
					return ctx, errors.Wrap(err, "Failed save invoice by ID.")
				}
				// TODO load transactions from DB
				next := true
				list, err := tx.SelectAllFrom(engine.TransactionTable, "WHERE invoice_id = $1", invID)
				if err != nil {
					return ctx, errors.Wrap(err, "Failed list transaction by invoice ID.")
				}
				for _, v := range list {
					tr := v.(*engine.Transaction)
					if !tr.Status.Match(engine.ACCEPTED_TX) {
						next = false
						break
					}
				}
				if !next {
					return ctx, nil
				}
				// Установить статус после проделанных операций
				inv.Status = engine.ACCEPTED_I
				inv.NextStatus = nil
				if err := tx.Save(&inv); err != nil {
					return ctx, errors.Wrap(err, "Failed save invoice by ID.")
				}
				return ctx, nil
			},
			"auth>accepted",
		)
		s.s.Add(
			ffsm.State(engine.AUTH_I),
			ffsm.State(engine.REJECTED_I),
			func(ctx context.Context, payload ffsm.Payload) (context context.Context, e error) {
				// TODO load invoice from DB
				invID, ok := payload.(int64)
				if !ok {
					log.Println("Invoice bad Payload: ", payload)
					return
				}
				tx := strategies.GetTXContext(ctx)
				if tx == nil {
					return ctx, errors.New("Not reform tx in context.")
				}
				inv := engine.Invoice{InvoiceID: invID}
				if err := tx.Reload(&inv); err != nil {
					return ctx, errors.Wrap(err, "Failed reload invoice by ID.")
				}
				if inv.Strategy != nameStrategy.String() {
					return ctx, errors.New("Invoice strategy not nameStrategy.")
				}
				if inv.Status != engine.AUTH_I {
					return ctx, errors.New("Invoice status not auth.")
				}
				// Установить статус куда происходит переход
				ns := engine.REJECTED_I
				inv.NextStatus = &ns
				if err := tx.Save(&inv); err != nil {
					return ctx, errors.Wrap(err, "Failed save invoice by ID.")
				}
				// TODO load transactions from DB
				next := true
				list, err := tx.SelectAllFrom(engine.TransactionTable, "WHERE invoice_id = $1", invID)
				if err != nil {
					return ctx, errors.Wrap(err, "Failed list transaction by invoice ID.")
				}
				for _, v := range list {
					tr := v.(*engine.Transaction)
					if !tr.Status.Match(engine.REJECTED_TX) {
						next = false
						break
					}
				}
				if !next {
					return ctx, nil
				}
				// Установить статус после проделанных операций
				inv.Status = engine.REJECTED_I
				inv.NextStatus = nil
				if err := tx.Save(&inv); err != nil {
					return ctx, errors.Wrap(err, "Failed save invoice by ID.")
				}
				return ctx, nil
			},
			"auth>rejected",
		)
		s.s.Add(
			ffsm.State(engine.WAIT_I),
			ffsm.State(engine.ACCEPTED_I),
			func(ctx context.Context, payload ffsm.Payload) (context context.Context, e error) {
				// TODO load invoice from DB
				invID, ok := payload.(int64)
				if !ok {
					log.Println("Invoice bad Payload: ", payload)
					return
				}
				tx := strategies.GetTXContext(ctx)
				if tx == nil {
					return ctx, errors.New("Not reform tx in context.")
				}
				inv := engine.Invoice{InvoiceID: invID}
				if err := tx.Reload(&inv); err != nil {
					return ctx, errors.Wrap(err, "Failed reload invoice by ID.")
				}
				if inv.Strategy != nameStrategy.String() {
					return ctx, errors.New("Invoice strategy not nameStrategy.")
				}
				if inv.Status != engine.WAIT_I {
					return ctx, errors.New("Invoice status not wait.")
				}
				// Установить статус куда происходит переход
				ns := engine.ACCEPTED_I
				inv.NextStatus = &ns
				if err := tx.Save(&inv); err != nil {
					return ctx, errors.Wrap(err, "Failed save invoice by ID.")
				}
				// TODO load transactions from DB
				next := true
				list, err := tx.SelectAllFrom(engine.TransactionTable, "WHERE invoice_id = $1", invID)
				if err != nil {
					return ctx, errors.Wrap(err, "Failed list transaction by invoice ID.")
				}
				for _, v := range list {
					tr := v.(*engine.Transaction)
					if !tr.Status.Match(engine.ACCEPTED_TX) {
						next = false
						break
					}
				}
				if !next {
					return ctx, nil
				}
				// Установить статус после проделанных операций
				inv.Status = engine.ACCEPTED_I
				inv.NextStatus = nil
				if err := tx.Save(&inv); err != nil {
					return ctx, errors.Wrap(err, "Failed save invoice by ID.")
				}
				return ctx, nil
			},
			"wait>accepted",
		)
		s.s.Add(
			ffsm.State(engine.WAIT_I),
			ffsm.State(engine.REJECTED_I),
			func(ctx context.Context, payload ffsm.Payload) (context context.Context, e error) {
				// TODO load invoice from DB
				invID, ok := payload.(int64)
				if !ok {
					log.Println("Invoice bad Payload: ", payload)
					return
				}
				tx := strategies.GetTXContext(ctx)
				if tx == nil {
					return ctx, errors.New("Not reform tx in context.")
				}
				inv := engine.Invoice{InvoiceID: invID}
				if err := tx.Reload(&inv); err != nil {
					return ctx, errors.Wrap(err, "Failed reload invoice by ID.")
				}
				if inv.Strategy != nameStrategy.String() {
					return ctx, errors.New("Invoice strategy not nameStrategy.")
				}
				if inv.Status != engine.WAIT_I {
					return ctx, errors.New("Invoice status not wait.")
				}
				// Установить статус куда происходит переход
				ns := engine.REJECTED_I
				inv.NextStatus = &ns
				if err := tx.Save(&inv); err != nil {
					return ctx, errors.Wrap(err, "Failed save invoice by ID.")
				}
				// TODO load transactions from DB
				next := true
				list, err := tx.SelectAllFrom(engine.TransactionTable, "WHERE invoice_id = $1", invID)
				if err != nil {
					return ctx, errors.Wrap(err, "Failed list transaction by invoice ID.")
				}
				for _, v := range list {
					tr := v.(*engine.Transaction)
					if !tr.Status.Match(engine.REJECTED_TX) {
						next = false
						break
					}
				}
				if !next {
					return ctx, nil
				}
				// Установить статус после проделанных операций
				inv.Status = engine.REJECTED_I
				inv.NextStatus = nil
				if err := tx.Save(&inv); err != nil {
					return ctx, errors.Wrap(err, "Failed save invoice by ID.")
				}
				return ctx, nil
			},
			"wait>rejected",
		)
		s.s.Add(
			ffsm.State(engine.WAIT_I),
			ffsm.State(engine.DRAFT_I),
			func(ctx context.Context, payload ffsm.Payload) (context context.Context, e error) {
				// TODO load invoice from DB
				invID, ok := payload.(int64)
				if !ok {
					log.Println("Invoice bad Payload: ", payload)
					return
				}
				tx := strategies.GetTXContext(ctx)
				if tx == nil {
					return ctx, errors.New("Not reform tx in context.")
				}
				inv := engine.Invoice{InvoiceID: invID}
				if err := tx.Reload(&inv); err != nil {
					return ctx, errors.Wrap(err, "Failed reload invoice by ID.")
				}
				if inv.Strategy != nameStrategy.String() {
					return ctx, errors.New("Invoice strategy not nameStrategy.")
				}
				if inv.Status != engine.WAIT_I {
					return ctx, errors.New("Invoice status not wait.")
				}
				// Установить статус куда происходит переход
				ns := engine.DRAFT_I
				inv.NextStatus = &ns
				if err := tx.Save(&inv); err != nil {
					return ctx, errors.Wrap(err, "Failed save invoice by ID.")
				}
				// TODO load transactions from DB
				next := true
				list, err := tx.SelectAllFrom(engine.TransactionTable, "WHERE invoice_id = $1", invID)
				if err != nil {
					return ctx, errors.Wrap(err, "Failed list transaction by invoice ID.")
				}
				for _, v := range list {
					tr := v.(*engine.Transaction)
					if !tr.Status.Match(engine.DRAFT_TX) {
						next = false
						break
					}
				}
				if !next {
					return ctx, nil
				}
				// Установить статус после проделанных операций
				inv.Status = engine.DRAFT_I
				inv.NextStatus = nil
				if err := tx.Save(&inv); err != nil {
					return ctx, errors.Wrap(err, "Failed save invoice by ID.")
				}
				return ctx, nil
			},
			"wait>draft",
		)
		s.s.Add(
			ffsm.State(engine.AUTH_I),
			ffsm.State(engine.DRAFT_I),
			func(ctx context.Context, payload ffsm.Payload) (context context.Context, e error) {
				// TODO load invoice from DB
				invID, ok := payload.(int64)
				if !ok {
					log.Println("Invoice bad Payload: ", payload)
					return
				}
				tx := strategies.GetTXContext(ctx)
				if tx == nil {
					return ctx, errors.New("Not reform tx in context.")
				}
				inv := engine.Invoice{InvoiceID: invID}
				if err := tx.Reload(&inv); err != nil {
					return ctx, errors.Wrap(err, "Failed reload invoice by ID.")
				}
				if inv.Strategy != nameStrategy.String() {
					return ctx, errors.New("Invoice strategy not nameStrategy.")
				}
				if inv.Status != engine.AUTH_I {
					return ctx, errors.New("Invoice status not auth.")
				}
				// Установить статус куда происходит переход
				ns := engine.DRAFT_I
				inv.NextStatus = &ns
				if err := tx.Save(&inv); err != nil {
					return ctx, errors.Wrap(err, "Failed save invoice by ID.")
				}
				// TODO load transactions from DB
				next := true
				list, err := tx.SelectAllFrom(engine.TransactionTable, "WHERE invoice_id = $1", invID)
				if err != nil {
					return ctx, errors.Wrap(err, "Failed list transaction by invoice ID.")
				}
				for _, v := range list {
					tr := v.(*engine.Transaction)
					if !tr.Status.Match(engine.DRAFT_TX) {
						next = false
						break
					}
				}
				if !next {
					return ctx, nil
				}
				// Установить статус после проделанных операций
				inv.Status = engine.DRAFT_I
				inv.NextStatus = nil
				if err := tx.Save(&inv); err != nil {
					return ctx, errors.Wrap(err, "Failed save invoice by ID.")
				}
				return ctx, nil
			},
			"auth>draft",
		)
	})
}
