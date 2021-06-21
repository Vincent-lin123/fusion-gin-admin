package hook

import (
	"fmt"
	"github.com/LyricTian/queue"
	"github.com/sirupsen/logrus"
	"os"
)

var defaultOptions = options{
	maxQueues:  512,
	maxWorkers: 1,
	levels: []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
		logrus.TraceLevel,
	},
}

type ExecCloser interface {
	Exec(entry *logrus.Entry) error
	Close() error
}

type FilterHandle func(*logrus.Entry) *logrus.Entry

type options struct {
	maxQueues  int
	maxWorkers int
	extra      map[string]interface{}
	filter     FilterHandle
	levels     []logrus.Level
}

type Option func(*options)

func SetMaxQueues(maxQueues int) Option {
	return func(o *options) {
		o.maxQueues = maxQueues
	}
}

func SetMaxWorkers(maxWorkers int) Option {
	return func(o *options) {
		o.maxWorkers = maxWorkers
	}
}

func SetExtra(extra map[string]interface{}) Option {
	return func(o *options) {
		o.extra = extra
	}
}

func SetFilter(filter FilterHandle) Option {
	return func(o *options) {
		o.filter = filter
	}
}

func SetLevels(levels ...logrus.Level) Option {
	return func(o *options) {
		if len(levels) == 0 {
			return
		}
		o.levels = levels
	}
}

type Hook struct {
	opts options
	q    *queue.Queue
	e    ExecCloser
}

func New(exec ExecCloser, opt ...Option) *Hook {
	opts := defaultOptions
	for _, o := range opt {
		o(&opts)
	}

	q := queue.NewQueue(opts.maxQueues, opts.maxWorkers)
	q.Run()

	return &Hook{
		opts: opts,
		q:    q,
		e:    exec,
	}
}

func (h *Hook) Levels() []logrus.Level {
	return h.opts.levels
}

func (h *Hook) Fire(entry *logrus.Entry) error {
	entry = h.copyEntry(entry)
	h.q.Push(queue.NewJob(entry, func(v interface{}) {
		h.exec(v.(*logrus.Entry))
	}))
	return nil
}

func (h *Hook) copyEntry(e *logrus.Entry) *logrus.Entry {
	entry := logrus.NewEntry(e.Logger)
	entry.Data = make(logrus.Fields)
	entry.Time = e.Time
	entry.Level = e.Level
	entry.Message = e.Message
	for k, v := range e.Data {
		entry.Data[k] = v
	}
	return entry
}

func (h *Hook) exec(entry *logrus.Entry) {
	for k, v := range h.opts.extra {
		if _, ok := entry.Data[k]; !ok {
			entry.Data[k] = v
		}
	}

	if filter := h.opts.filter; filter != nil {
		entry = filter(entry)
	}

	err := h.e.Exec(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[logrus-hook] execution error: %s", err.Error())
	}
}

func (h *Hook) Flush() {
	h.q.Terminate()
	h.e.Close()
}
