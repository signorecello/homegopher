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
	HA Connection
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
	s := SwitchAttributes{}

	// get initial state, assert type
	state := test.GetState()
	assert.IsType(t, SwitchState{}, state)

	// set on and check returning state to be on
	state = test.SetState("on", s)
	assert.Equal(t, "on", state.State)

	// same for off
	state = test.SetState("off", s)
	assert.Equal(t, "off", state.State)

	// changing it to on, the service won't trigger because it's a dummy entity
	// but should give us a status code
	status := test.Change("turn_on")
	assert.Equal(t, 200, status)
	log.Println(status)

	// listening to the state changed channel for that switch
	listen := test.Listen()

	// now we'll prepare the special payload with the "testing" flag
	go func() {
		s = SwitchAttributes{"Test": "testing"}
		test.SetState("on", s)
	}()

	// now we listen for the channel, see if we get mr. flag in the attributes
	func(listen chan SwitchStateChanged) {
		for l := range listen {
			if l.SwitchState.Attributes["Test"] == "testing" {
				assert.Equal(t, "on", l.SwitchState.State)

				s = SwitchAttributes{"Test": ""}
				test.SetState("off", s)
				return
			}
		}

	}(listen)
}

func TestLight(t *testing.T) {
	test := HA.NewLight("some_light")
	state := test.GetState()
	assert.IsType(t, LightState{}, state)

	s := LightAttributes{}
	state = test.SetState("on", s)
	assert.Equal(t, "on", state.State)

	state = test.SetState("off", s)
	assert.Equal(t, "off", state.State)

	state = test.Change("turn_on")
	assert.IsType(t, LightState{}, state)

	listen := test.Listen()
	go func() {
		s = LightAttributes{"Test": "testing"}
		test.SetState("on", s)
	}()

	func(listen chan LightStateChanged) {
		for l := range listen {
			if l.LightState.Attributes["Test"] == "testing" {
				assert.Equal(t, "on", l.LightState.State)

				s = LightAttributes{"Test": ""}
				test.SetState("off", s)
				return
			}
		}

	}(listen)

}

func TestSensor(t *testing.T) {
	test := HA.NewSensor("some_sensor")
	state := test.GetState()
	assert.IsType(t, SensorState{}, state)

	s := SensorAttributes{}
	state = test.SetState("on", s)
	assert.Equal(t, "on", state.State)

	state = test.SetState("off", s)
	assert.Equal(t, "off", state.State)

	listen := test.Listen()
	go func() {
		s = SensorAttributes{"Test": "testing"}
		test.SetState("on", s)
	}()

	func(listen chan SensorStateChanged) {
		for l := range listen {
			if l.SensorState.Attributes["Test"] == "testing" {
				assert.Equal(t, "on", l.SensorState.State)

				s = SensorAttributes{"Test": ""}
				test.SetState("off", s)
				return
			}
		}

	}(listen)

}

func TestBinarySensor(t *testing.T) {
	test := HA.NewBinarySensor("some_binary_sensor")
	state := test.GetState()
	assert.IsType(t, BinarySensorState{}, state)

	s := BinarySensorAttributes{}
	state = test.SetState("on", s)
	assert.Equal(t, "on", state.State)

	state = test.SetState("off", s)
	assert.Equal(t, "off", state.State)

	listen := test.Listen()
	go func() {
		s = BinarySensorAttributes{"Test": "testing"}
		test.SetState("on", s)
	}()

	func(listen chan BinarySensorStateChanged) {
		for l := range listen {
			if l.BinarySensorState.Attributes["Test"] == "testing" {
				assert.Equal(t, "on", l.BinarySensorState.State)

				s = BinarySensorAttributes{"Test": ""}
				test.SetState("off", s)
				return
			}
		}

	}(listen)
}
