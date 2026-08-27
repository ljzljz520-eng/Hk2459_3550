package workflow

import (
	"context"
	"paymentconsole/internal/service"
	"paymentconsole/internal/store"
	"testing"
)

func TestWorkflowCreateReviewArchive(t *testing.T) {
	st, _ := store.Open(t.TempDir() + "/db")
	defer st.Close()
	e := New(service.New(st), st)
	if _, x := e.CreateReviewArchive(context.Background(), "main", []string{"create", "review", "confirm", "archive"}); x != nil {
		t.Fatal(x)
	}
}
func TestWorkflowSearchUpdatePublish(t *testing.T) {
	st, _ := store.Open(t.TempDir() + "/db")
	defer st.Close()
	e := New(service.New(st), st)
	if e.SearchUpdatePublish(context.Background(), []string{}, true) == nil {
		t.Fatal("expected map")
	}
}
