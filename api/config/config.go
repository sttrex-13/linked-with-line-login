package config

import (
	"os"
	"reflect"
)

type Configuration struct {
	// change to your url
	AppURL string `env:"APP_URL" default:"https://positive-muskrat-emerging.ngrok-free.app"`
	// line login credentials
	LINELoginClientID     string `env:"LINE_LOGIN_CLIENT_ID" default:"2005138101"`
	LINELoginChanelSecret string `env:"LINE_LOGIN_CHANEL_SECRET" default:"e6f25e7be49afa72300489c4a8241347"`
	// line messaging credentials
	LINEMessagingAccessToken string `env:"LINE_MESSAGING_ACCESS_TOKEN" default:"fEd+xwkjCLrKX0Y3kY6i3J1HXpPrv3MN825KF95ES6R1AvARfom5LY2FhgsdkXGpleUzR07LVM8rqMihUK4RVAexkD2hyp0aQ33fDhvTxM0qbQ/0Hw/qwtyO9EOeio0yUHKhKhNs8Ik9TQyeGbZ3JAdB04t89/1O/w1cDnyilFU="`
}

func New() Configuration {
	conf := Configuration{}
	v := reflect.ValueOf(&conf).Elem()
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		envKey := fieldType.Tag.Get("env")
		envValue, ok := os.LookupEnv(envKey)
		switch ok {
		case true:
			field.SetString(envValue)
		case false:
			field.SetString(fieldType.Tag.Get("default"))
		}
	}
	return conf
}
