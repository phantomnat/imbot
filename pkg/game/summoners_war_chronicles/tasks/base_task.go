package tasks

import (
	"bytes"
	"encoding/json"
	"image"
	"strings"
	"time"

	"go.uber.org/zap"
	"gocv.io/x/gocv"

	"github.com/pkg/errors"

	"github.com/phantomnat/imbot/pkg/domain"
)

type BaseTask struct {
	Index          int
	Name           string
	StateChangedAt time.Time
	State          domain.TaskState
	StateTexts     map[domain.TaskState]string
	Log            *zap.SugaredLogger
	Manager        domain.Manager
	Im             domain.ImageManager
	Exiting        bool
}

var _ domain.Task = (*BaseTask)(nil)

func (t *BaseTask) GetName() string {
	return t.Name
}

func (t *BaseTask) Do(_ gocv.Mat) bool {
	return false
}
func (t *BaseTask) LoadStatus(in any) {
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

type SearchROIOption struct {
	needClick   bool
	path        string
	roi         image.Rectangle
	wait        *time.Duration
	nextState   *domain.TaskState
	clickOffset image.Point

	debugTemplateMatch bool
}

func (o *SearchROIOption) getOffsetPt(pt image.Point) image.Point {
	if o.clickOffset == (image.Point{}) {
		return pt
	}
	return image.Pt(pt.X+o.clickOffset.X, pt.Y+o.clickOffset.Y)
}

func (o *SearchROIOption) apply(opts ...SearchROIOptionFunc) {
	for _, fn := range opts {
		fn(o)
	}
	if o.wait == nil {
		v := time.Millisecond * 1000
		o.wait = &v
	}
}
func (o *SearchROIOption) isEmpty() bool {
	if o == nil {
		return true
	}
	emptyPath := o.path == ""
	return emptyPath
}

type SearchROIOptionFunc func(o *SearchROIOption)

func WithClick() SearchROIOptionFunc {
	return func(o *SearchROIOption) {
		o.needClick = true
	}
}
func WithPath(paths ...string) SearchROIOptionFunc {
	return func(o *SearchROIOption) {
		o.path = strings.Join(paths, ".")
	}
}
func WithROI(roi image.Rectangle) SearchROIOptionFunc {
	return func(o *SearchROIOption) {
		o.roi = roi
	}
}
func WithDebugMatch() SearchROIOptionFunc {
	return func(o *SearchROIOption) {
		o.debugTemplateMatch = true
	}
}
func WithWaitMs(ms int) SearchROIOptionFunc {
	return func(o *SearchROIOption) {
		v := time.Duration(ms) * time.Millisecond
		o.wait = &v
	}
}
func WithNoWait() SearchROIOptionFunc {
	return WithWaitMs(0)
}
func WithNextState(s domain.TaskState) SearchROIOptionFunc {
	return func(o *SearchROIOption) {
		o.nextState = &s
	}
}
func WithClickOffset(roi image.Rectangle) SearchROIOptionFunc {
	return func(o *SearchROIOption) {
		o.clickOffset = roi.Min
	}
}

func (t *BaseTask) WaitMs(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func (t *BaseTask) SearchROI(
	m gocv.Mat,
	opts ...SearchROIOptionFunc,
) (found bool) {
	opt := SearchROIOption{}
	opt.apply(opts...)

	if opt.isEmpty() {
		return false
	}

	var pt image.Point

	if opt.roi == (image.Rectangle{}) {
		found, pt = t.Im.MatchWithOption(m, domain.MatchOption{
			Path: t.Manager.GetImagePath(opt.path),
		})
	} else {
		found, pt = t.Manager.MatchInROI(m, opt.roi, domain.MatchOption{
			Path:     opt.path,
			PrintVal: opt.debugTemplateMatch,
		})
	}

	if !found {
		return
	}
	if opt.needClick {
		t.Manager.ClickPt(opt.getOffsetPt(pt))
	}
	if opt.wait != nil || *opt.wait > 0 {
		time.Sleep(*opt.wait)
	}
	if opt.nextState != nil {
		t.SetState(*opt.nextState)
	}
	return
}

func (t *BaseTask) GetPtWithROI(roi image.Rectangle, pt image.Point) image.Point {
	return image.Pt(pt.X+roi.Min.X, pt.Y+roi.Min.Y)
}
