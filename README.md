# homegopher

## What is dis?

This is a package for interacting with Homeassistant using Go.

## But... why?

I'm kind of a tinkerer and I wanted to learn Go. I also like Homeassistant and I wanted to create a telegram bot for it. Unfortunately I noticed there isn't an active maintained package for that, and those that exist lack some features.

So I decided to code my own using my basic Go skills.

## What does it do?

It currently implements three APIs:

- REST API
- Websocket API
- MQTT Api

Not all features are available for all APIs, but I believe it serves most of the purpose you may need. If it doesn't, feel free to submit a pull request, or just ping me and I'll think about it.

## How do I use it?

Fear not, here's a basic implementation for you:

```
conn := ha.Connection{
		Prefix:        http or https,
		Host:          <your host, example 192.168.1.1>
		Path:          /api,
		Port:          <your port, example 8123>,
		Authorization: <a long-lived authorization token>,
	}
	HA := ha.NewConnection(conn)

```
You need an authorization token, [follow these basic instructions](https://developers.home-assistant.io/docs/auth_api/#long-lived-access-token) on how to get one for your installation.

After this, you should be authenticated with your HA install. 

### Getting states
Want to get a sensor state?

```
sensor := entities.NewSensor("some_sensor", HA)
sensorState := sensor.GetState().State
```

### Setting states
You can set a sensor state using

```
sensor := entities.NewSensor("some_sensor", HA)
sensorState := sensor.SetState("whatever state your want")
```

### Switching things off
You can turn lights on and off, passing Opts if you want something on servicedata (such as Kelvin, Brightness, etc)

```
state = test.TurnOn(entities.LightOpts{})
```

## Examples please?

There are some examples in the self-explaining "example" folder. You can also dig into the tests, they should give you some more examples.

