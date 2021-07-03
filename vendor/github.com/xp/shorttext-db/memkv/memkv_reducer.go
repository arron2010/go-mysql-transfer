package memkv

import (
	"github.com/xp/shorttext-db/easymr/artifacts/task"
	"sort"
)

type MemKVReducer struct {
}

func NewMemKVReducer() *MemKVReducer {
	m := &MemKVReducer{}
	return m
}

func (m *MemKVReducer) Reduce(sources map[int]*task.Task) (map[int]*task.Task, *task.TaskResult, error) {
	dbItems := make([]*DBItem, 0, 4)
	list := &DBItems{}
	for _, t := range sources {
		for _, r := range t.Result {
			items := r.(*DBItems)
			dbItems = append(dbItems, items.Items...)
		}
	}
	list.Items = dbItems
	if len(dbItems) > 1 {
		sort.Sort(list)
	}

	result := task.NewTaskResult(list)
	logger.Infof("数据库扫描汇总, 记录数:%d\n", len(list.Items))
	return sources, result, nil
}
