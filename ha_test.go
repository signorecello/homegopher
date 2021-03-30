package haclient

import (
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

var (
	HA Connection
	ss = make(chan SensorState)
	ls = make(chan LightState)
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

	stateChanges := StateChanges{Sensor: ss, Light: ls}
	go NewWS(
		time.Second*10,
		stateChanges,
		os.Getenv("WSURL"),
		os.Getenv("AUTHORIZATION"),
	)

	code := m.Run()
	os.Exit(code)
}

func TestSwitch(t *testing.T) {
	test := HA.NewSwitch("some_switch")
	state := test.GetState()
	assert.IsType(t, SwitchState{}, state)

	state = test.SetState("on")
	assert.Equal(t, "on", state.State)

	state = test.SetState("off")
	assert.Equal(t, "off", state.State)

	state = test.Change("turn_on")
	assert.IsType(t, SwitchState{}, state)
}

func TestLight(t *testing.T) {
	test := HA.NewLight("some_light")
	state := test.GetState()
	assert.IsType(t, LightState{}, state)

	state = test.SetState("on")
	assert.Equal(t, "on", state.State)

	state = test.SetState("off")
	assert.Equal(t, "off", state.State)

	state = test.Change("turn_on")
	assert.IsType(t, LightState{}, state)
}

func TestSensor(t *testing.T) {
	test := HA.NewSensor("some_sensor")
	state := test.GetState()
	assert.IsType(t, SensorState{}, state)

	state = test.SetState("on")
	assert.Equal(t, "on", state.State)

	state = test.SetState("off")
	assert.Equal(t, "off", state.State)
}

func TestBinarySensor(t *testing.T) {
	test := HA.NewBinarySensor("some_binary_sensor")
	state := test.GetState()
	assert.IsType(t, BinarySensorState{}, state)

	state = test.SetState("on")
	assert.Equal(t, "on", state.State)

	state = test.SetState("off")
	assert.Equal(t, "off", state.State)
}

func TestWS(t *testing.T) {
	var s SensorState
	go func() {
		s = <-ss
	}()
	defer close(ss)

	test := HA.NewSensor("some_sensor")
	test.SetState("on")

	assert.IsType(t, SensorState{}, s)
}
