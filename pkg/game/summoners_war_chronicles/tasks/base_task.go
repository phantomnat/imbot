package tasks

import (
	"bytes"
	"encoding/json"

	"go.uber.org/zap"
	"gocv.io/x/gocv"

	"github.com/pkg/errors"

	"github.com/phantomnat/imbot/pkg/domain"
)

type BaseTask struct {
	Index      int
	Name       string
	State      domain.TaskState
	StateTexts map[domain.TaskState]string
	Log        *zap.SugaredLogger
	Manager    domain.Manager
	Exiting    bool
}

var _ domain.Task = (*BaseTask)(nil)

func (t *BaseTask) GetName() string {
	return t.Name
}

func (t *BaseTask) Do(_ gocv.Mat) bool {
	return false
}

func (t *BaseTask) GetState() string {
	return t.StateTexts[t.State]
}

func (t *BaseTask) SetState(s domain.TaskState) {
	if _, exist := t.StateTexts[s]; !exist {
		t.Log.Panicf("please add state '%d' (%T)", s)
	}
	from := t.State
	t.State = s
	to := s
	t.Log.Infof("state changed '%s' (%d) -> '%s' (%d)", t.StateTexts[from], from, t.StateTexts[to], to)
}

func (t *BaseTask) IsNeedMainScreen() bool {
	return true
}

func (t *BaseTask) CanInterrupt() bool {
	return true
}

func (t *BaseTask) IsReady() bool {
	panic("implement me")
}

func (t *BaseTask) Exit() {
	t.SetState(domain.TaskStateBegin)
	t.Manager.ExitTask()
}

func (t *BaseTask) RequestExit() {
	t.Exiting = true
}

func (t *BaseTask) SaveStatus(v any) {
	t.Manager.SaveStatus(t.Index, t.Name, v)
}

func (t *BaseTask) LoadStatus(in any) {
	_ = convertTo(t.Manager.LoadStatus(t.Index, t.Name), in)
}

func convertTo(in any, out any) error {
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(in)
	if err != nil {
		return errors.Wrapf(err, "encode input (%T) to json", in)
	}
	err = json.NewDecoder(buf).Decode(&out)
	if err != nil {
		return errors.Wrapf(err, "decode input (%T) json to struct (%T)", in, out)
	}
	return nil
}
