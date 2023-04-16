package summonerswar

import (
	"time"

	"github.com/phantomnat/imbot/pkg/domain"
)

type task struct {
}

const taskName = "repeatQuest"

func New() domain.Task {
	return &task{}
}

func (t *task) Do(index int, now time.Time) bool {
	return false
}

func (t *task) GetState() string {
	return ""
}
