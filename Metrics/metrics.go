package Metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ActionNotFoundErr = Error("action is not found")
)

type Metrics struct {
	Namespace string

	actions map[string][]ActionHandler

	registry *prometheus.Registry
}
type ActionHandler interface {
	Do(payload ...string) error
}

func (m *Metrics) Init() {
	m.registry = prometheus.NewRegistry()
	m.actions = make(map[string][]ActionHandler)
}

func (m *Metrics) AddCounter(c Counter) error {
	c.counter = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace:   m.Namespace,
		Name:        c.Name,
		Help:        c.Help,
		ConstLabels: c.Labels,
	})

	m.actions[c.Action] = append(m.actions[c.Action], &c)

	err := m.registry.Register(c.counter)
	if err != nil {
		return err
	}
	return nil
}
func (m *Metrics) Action(action string, payload ...string) error {
	val, ok := m.actions[action]

	if !ok {
		return ActionNotFoundErr
	}

	for _, e := range val {
		err := e.Do(payload...)
		if err != nil {
			return err
		}
	}
	return nil
}
func (m *Metrics) Actions(action ...string) error {
	for _, a := range action {
		err := m.Action(a)
		if err != nil {
			return err
		}
	}
	return nil
}

type Counter struct {
	Name   string
	Help   string
	Labels map[string]string

	Action string

	counter prometheus.Counter
	count   float64
}

func (c *Counter) Inc() {
	c.count++
	c.counter.Inc()
}
func (c *Counter) Add(n float64) {
	c.count += n
	c.counter.Add(n)
}
func (c *Counter) Do(_ ...string) error {
	c.Inc()
	return nil
}

func (m *Metrics) GetHandler() http.Handler {
	return promhttp.HandlerFor(
		m.registry,
		promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		})
}
