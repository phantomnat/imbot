package tasks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
	"gocv.io/x/gocv"

	"github.com/pkg/errors"

	"github.com/phantomnat/imbot/pkg/domain"
)

type BaseTask struct {
	// Index          int
	Name           string
	StateChangedAt time.Time
	State          domain.TaskState
	StateTexts     map[domain.TaskState]string
	Log            *zap.SugaredLogger
	Manager        domain.Manager
	Im             domain.ImageManager
	Exiting        bool

	setting domain.TaskSetting
	status  domain.TaskStatus
}

var _ domain.Task = (*BaseTask)(nil)

func NewBaseTask(
	manager domain.Manager,
	setting domain.TaskSetting,
	status domain.TaskStatus,
	stateTexts map[domain.TaskState]string,
) BaseTask {
	name := strings.SplitN(fmt.Sprintf("%T", setting), ".", 2)[0]
	name = strings.TrimLeft(name, "*")
	b := BaseTask{
		Im: manager.GetImageManager(),
		// Index:   index,
		Name:    name,
		Manager: manager,
		Log:     zap.S().Named("task").Named(name),
		StateTexts: map[domain.TaskState]string{
			domain.TaskStateBegin: "begin",
			domain.TaskStateEnd:   "end",
		},
		setting: setting,
		status:  status,
	}
	for k, v := range stateTexts {
		b.StateTexts[k] = v
	}
	return b
}

func (t *BaseTask) GetName() string {
	return t.Name
}

func (t *BaseTask) Do(_ gocv.Mat) bool {
	return false
}

func (t *BaseTask) LoadStatus(in any) {
	err := t.ConvertTo(in, &t.status)
	if err != nil {
		t.Log.Warnf("reset status, cannot the current: %+v", err)
	} else {
		t.SetState(t.status.GetState())
	}
}

func (t *BaseTask) Reset() {
	t.SetState(domain.TaskStateBegin)
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
	t.StateChangedAt = time.Now()
	t.status.SetState(to)
	t.Log.Infof("state changed '%s' (%d) -> '%s' (%d)", t.StateTexts[from], from, t.StateTexts[to], to)
}

func (t *BaseTask) IsNeedMainScreen() bool {
	return true
}

func (t *BaseTask) CanInterrupt() bool {
	return true
}

func (t *BaseTask) IsReady() bool {
	if !t.setting.IsEnabled() {
		return false
	}
	return time.Now().After(t.status.GetNextExecution())
}

func (t *BaseTask) Exit() {
	t.SetState(domain.TaskStateBegin)
	t.Manager.ExitTask()
}

func (t *BaseTask) RequestExit() {
	t.Exiting = true
}

func (t *BaseTask) SaveStatus() {
	t.Manager.SaveStatus(t.Name, t.status)
}

func (t *BaseTask) ConvertTo(in, out any) error {
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

func (t *BaseTask) DurationSinceStateChanged() time.Duration {
	return time.Since(t.StateChangedAt)
}
