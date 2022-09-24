package handler

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/uekiGityuto/go_todo_app/entity"
	"github.com/uekiGityuto/go_todo_app/store"
)

type ListTask struct {
	DB   *sqlx.DB
	Repo *store.Repository
}

type task struct {
	ID     entity.TaskID     `json:"id"`
	Title  string            `json:"title"`
	Status entity.TaskStatus `json:"status"`
}

func (lt *ListTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tasks, err := lt.Repo.ListTasks(ctx, lt.DB)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
	}
	// var rsp []task // これだと取得結果が存在しないときにnilを返してしまう
	rsp := []task{} // 空のsliceで初期化（取得結果が存在しないときに空のsliceを返せるように）
	for _, t := range tasks {
		rsp = append(rsp, task{
			ID:     t.ID,
			Title:  t.Title,
			Status: t.Status,
		})
	}
	RespondJSON(ctx, w, rsp, http.StatusOK)
}
