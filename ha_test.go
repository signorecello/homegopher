package homegopher

import (
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
	"time"
)

var (
	HA          Connection
	TestTimeout time.Duration = time.Second * 2
)

func TestMain(m *testing.M) {
	_ = godotenv.Load()

	conn := Connection{
		Prefix:        os.Getenv("PREFIX"),
		Host:          os.Getenv("HOST"),
		Path:          os.Getenv("HOST_PATH"),
		Port:          os.Getenv("PORT"),
		Authorization: os.Getenv("AUTHORIZATION"),
	}
	HA = NewConnection(conn)
	go NewWS(
		time.Second*10,
		os.Getenv("WSURL"),
		os.Getenv("AUTHORIZATION"),
	)

	code := m.Run()
	os.Exit(code)
}

func TestSwitch(t *testing.T) {
	// initialize switch
	test := HA.NewSwitch("some_switch")

	// initialize attributes so we can add a little flag for testing
	s := Attributes{}

	// get initial state, assert type
	state := test.GetState()
	assert.IsType(t, State{}, state)

	// set on and check returning state to be on
	state = test.SetState("on", s)
	assert.Equal(t, "on", state.State)

	// same for off
	state = test.SetState("off", s)
	assert.Equal(t, "off", state.State)

	// changing it to on, the service won't trigger because it's a dummy entity
	// but should give us a status code
	status := test.Change(ServiceCall{
		Service: "turn_on",
	})
	assert.Equal(t, 200, status)
	log.Println(status)

	// listening to the state changed channel for that switch
	listen := test.Listen()

	// now we'll prepare the special payload with the "testing" flag
	go func() {
		s = Attributes{"Test": "testing"}
		test.SetState("on", s)
		time.Sleep(TestTimeout)
		listen <- StateChangedEvent{Type: "fail"}
	}()

	// now we listen for the channel, see if we get mr. flag in the attributes
	func(listen chan StateChangedEvent) {
		for l := range listen {
			if l.Event.Data.NewState.Attributes["Test"] == "testing" {
				assert.Equal(t, "on", l.Event.Data.NewState.State)

				s = Attributes{"Test": ""}
				test.SetState("off", s)
				return
			} else if l.Type == "fail" {
				assert.Fail(t, "Timeout")
				return
			}
		}

	}(listen)
}

func TestLight(t *testing.T) {
	test := HA.NewLight("some_light")
	state := test.GetState()
	assert.IsType(t, State{}, state)

	s := Attributes{}
	state = test.SetState("on", s)
	assert.Equal(t, "on", state.State)

	state = test.SetState("off", s)
	assert.Equal(t, "off", state.State)

	state = test.Change(ServiceCall{
		Service: "turn_on",
	})
	assert.IsType(t, State{}, state)

	listen := test.Listen()
	go func() {
		s = Attributes{"Test": "testing"}
		test.SetState("on", s)
		time.Sleep(TestTimeout)
		listen <- StateChangedEvent{Type: "fail"}
	}()

	func(listen chan StateChangedEvent) {
		for l := range listen {
			if l.Event.Data.NewState.Attributes["Test"] == "testing" {
				assert.Equal(t, "on", l.Event.Data.NewState.State)

				s = Attributes{"Test": ""}
				test.SetState("off", s)
				return
			} else if l.Type == "fail" {
				assert.Fail(t, "Timeout")
				return
			}
		}

	}(listen)

}

func TestSensor(t *testing.T) {
	test := HA.NewSensor("some_sensor")
	state := test.GetState()
	assert.IsType(t, State{}, state)

	s := Attributes{}
	state = test.SetState("on", s)
	assert.Equal(t, "on", state.State)

	state = test.SetState("off", s)
	assert.Equal(t, "off", state.State)

	listen := test.Listen()
	go func() {
		s = Attributes{"Test": "testing"}
		test.SetState("on", s)
		time.Sleep(TestTimeout)
		listen <- StateChangedEvent{Type: "fail"}
	}()

	func(listen chan StateChangedEvent) {
		for l := range listen {
			if l.Event.Data.NewState.Attributes["Test"] == "testing" {
				assert.Equal(t, "on", l.Event.Data.NewState.State)

				s = Attributes{"Test": ""}
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
	test := HA.NewBinarySensor("some_binary_sensor")
	state := test.GetState()
	assert.IsType(t, State{}, state)

	s := Attributes{}
	state = test.SetState("on", s)
	assert.Equal(t, "on", state.State)

	state = test.SetState("off", s)
	assert.Equal(t, "off", state.State)

	listen := test.Listen()
	go func() {
		s = Attributes{"Test": "testing"}
		test.SetState("on", s)
		time.Sleep(TestTimeout)
		listen <- StateChangedEvent{Type: "fail"}
	}()

	func(listen chan StateChangedEvent) {
		for l := range listen {
			if l.Event.Data.NewState.Attributes["Test"] == "testing" {
				assert.Equal(t, "on", l.Event.Data.NewState.State)

				s = Attributes{"Test": ""}
				test.SetState("off", s)
				return
			} else if l.Type == "fail" {
				assert.Fail(t, "Timeout")
				return
			}
		}

	}(listen)
}
