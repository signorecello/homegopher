package homegopher

import (
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
	"time"
	"github.com/signorecello/homegopher/ha"
	"github.com/signorecello/homegopher/state"
	"github.com/signorecello/homegopher/entities"
	"github.com/signorecello/homegopher/events"
)

var (
	HA          ha.Connection
	TestTimeout time.Duration = time.Second * 2
)

func TestMain(m *testing.M) {
	_ = godotenv.Load()

	conn := ha.Connection{
		Prefix:        os.Getenv("PREFIX"),
		Host:          os.Getenv("HOST"),
		Path:          os.Getenv("HOST_PATH"),
		Port:          os.Getenv("PORT"),
		Authorization: os.Getenv("AUTHORIZATION"),
	}
	HA = ha.NewConnection(conn)
	go ha.NewWS(
		time.Second*10,
		os.Getenv("WSURL"),
		os.Getenv("AUTHORIZATION"),
	)

	code := m.Run()
	os.Exit(code)
}

func TestSwitch(t *testing.T) {
	// initialize switch
	test := entities.NewSwitch("some_switch", HA)

	// get initial state, assert type
	st := test.GetState()
	assert.IsType(t, state.State{}, st)

	// changing it
	status := test.TurnOff()
	assert.IsType(t, state.State{}, st)

	status = test.TurnOn(entities.SwitchOpts{})
	assert.IsType(t, state.State{}, st)

	log.Println(status)

}

func TestLight(t *testing.T) {
	test := entities.NewLight("some_light", HA)
	st := test.GetState()
	assert.IsType(t, state.State{}, st)

	st = test.TurnOff()
	assert.IsType(t, state.State{}, st)

	st = test.TurnOn(entities.LightOpts{Brightness: "255"})
	assert.IsType(t, state.State{}, st)

	assert.IsType(t, state.State{}, st)
}

func TestSensor(t *testing.T) {
	test := entities.NewSensor("some_sensor", HA)
	st := test.GetState()
	assert.IsType(t, state.State{}, st)

	s := state.Attributes{}
	st = test.SetState("on", s)
	assert.Equal(t, "on", st.State)

	st = test.SetState("off", s)
	assert.Equal(t, "off", st.State)

	listen := test.Listen()
	go func() {
		s = state.Attributes{"Test": "testing"}
		test.SetState("on", s)
		time.Sleep(TestTimeout)
		listen <- events.StateChangedEvent{Type: "fail"}
	}()

	func(listen chan events.StateChangedEvent) {
		for l := range listen {
			if l.Event.Data.NewState.Attributes["Test"] == "testing" {
				assert.Equal(t, "on", l.Event.Data.NewState.State)

				s = state.Attributes{"Test": ""}
				test.SetState("off", s)
				return
			} else if l.Type == "fail" {
				assert.Fail(t, "Timeout")
				return
			}
		}

	}(listen)

}

func TestBinarySensor(t *testing.T) {
	test := entities.NewBinarySensor("some_binary_sensor", HA)
	st := test.GetState()
	assert.IsType(t, state.State{}, st)

	s := state.Attributes{}
	st = test.SetState("on", s)
	assert.Equal(t, "on", st.State)

	st = test.SetState("off", s)
	assert.Equal(t, "off", st.State)

	listen := test.Listen()
	go func() {
		s = state.Attributes{"Test": "testing"}
		test.SetState("on", s)
		time.Sleep(TestTimeout)
		listen <- events.StateChangedEvent{Type: "fail"}
	}()

	func(listen chan events.StateChangedEvent) {
		for l := range listen {
			if l.Event.Data.NewState.Attributes["Test"] == "testing" {
				assert.Equal(t, "on", l.Event.Data.NewState.State)

				s = state.Attributes{"Test": ""}
				test.SetState("off", s)
				return
			} else if l.Type == "fail" {
				assert.Fail(t, "Timeout")
				return
			}
		}

	}(listen)
}
