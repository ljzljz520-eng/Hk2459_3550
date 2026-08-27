package workflow

import (
	"context"
	"errors"
	"paymentconsole/internal/model"
	"paymentconsole/internal/service"
	"paymentconsole/internal/store"
)

type Engine struct {
	s  *service.Service
	st *store.Store
}

func New(s *service.Service, st *store.Store) *Engine { return &Engine{s: s, st: st} }
func (e *Engine) CreateReviewArchive(ctx context.Context, name string, steps []string) (model.Workflow, error) {
	if len(steps) < 4 {
		return model.Workflow{}, errors.New("workflow needs four steps")
	}
	w := e.s.CreateWorkflow(name, steps)
	w.Status = "review"
	if err := e.st.SaveWorkflow(w); err != nil {
		return w, err
	}
	w.Status = "archived"
	if err := e.st.SaveWorkflow(w); err != nil {
		return w, err
	}
	return w, nil
}
func (e *Engine) SearchUpdatePublish(ctx context.Context, ids []string, enabled bool) map[string]error {
	return e.s.BatchUpdate(ctx, ids, enabled)
}
func (e *Engine) ImportReport(ctx context.Context, records []model.Record) error {
	for _, r := range records {
		if r.ID == "" {
			return errors.New("missing record id")
		}
		if err := e.st.SaveRecord(r); err != nil {
			return err
		}
	}
	return nil
}
func (e *Engine) Validate(w model.Workflow) error {
	if w.Name == "" || len(w.Steps) < 4 {
		return errors.New("invalid workflow")
	}
	return nil
}
