package main

type Config struct {
	Url          string `json:"url"`
	Authurl      string `json:"authurl"`
	Credentials  Creds
	MqttSettings []MqttSetting
}

type Lynustypeone struct {
	N string  `json:"n"`
	V float64 `json:"v"`
}

type Lynustype []struct {
	N string  `json:"n"`
	V float64 `json:"v"`
}

type Settings struct {
	Credentionals Creds
	AuthUrl       string `json:"authurl"`
	MqttSettings  []MqttSetting
}

type MqttSetting struct {
	ConnectionString string `json:"connectionstring"`
	Port             string `json:"port"`
	Username         string `json:"username"`
	Password         string `json:"password"`
	Topic            string `json:"topic"`
}

type Creds struct {
	Username string `json:"email"`
	Password string `json:"password"`
}

type Bakler struct {
	ID       string `json:"_id"`
	Token    string `json:"token"`
	ExpireAt int    `json:"expireAt"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Relays   []struct {
		Name         string `json:"name"`
		SerialNumber string `json:"serialNumber"`
		ID           string `json:"_id"`
		Index        int    `json:"index"`
	} `json:"relays"`
	IsAdmin bool `json:"isAdmin"`
}

type BaklerVoltage struct {
	Data []struct {
		Voltage struct {
			Min    float64 `json:"min"`
			Max    float64 `json:"max"`
			Actual float64 `json:"actual"`
		} `json:"voltage"`
		Timestamp int64 `json:"timestamp"`
	} `json:"data"`
	TypeRelay string `json:"typeRelay"`
	Limits    struct {
	} `json:"limits"`
}

type BaklerAmperage struct {
	Data []struct {
		Current struct {
			Min    float64 `json:"min"`
			Max    float64 `json:"max"`
			Actual float64 `json:"actual"`
		} `json:"current"`
		Timestamp int64 `json:"timestamp"`
	} `json:"data"`
	TypeRelay string `json:"typeRelay"`
	Limits    struct {
		Min int `json:"min"`
		Max int `json:"max"`
	} `json:"limits"`
}

type BaklerPower struct {
	Data []struct {
		Power struct {
			Min    float64 `json:"min"`
			Max    float64 `json:"max"`
			Actual float64 `json:"actual"`
		} `json:"power"`
		TotPosEnergyPA float64 `json:"tot_posEnergyPA"`
		TotNegEnergyPA int     `json:"tot_negEnergyPA"`
		Timestamp      int64   `json:"timestamp"`
	} `json:"data"`
	TypeRelay string `json:"typeRelay"`
	Limits    struct {
	} `json:"limits"`
}

type Voltage struct {
	Voltage struct {
		Min    float64 `json:"min"`
		Max    float64 `json:"max"`
		Actual float64 `json:"actual"`
	} `json:"voltage"`
	Timestamp int64 `json:"timestamp"`
}

type Amperage struct {
	Current struct {
		Min    float64 `json:"min"`
		Max    float64 `json:"max"`
		Actual float64 `json:"actual"`
	} `json:"current"`
	Timestamp int64 `json:"timestamp"`
}

type Power struct {
	Power struct {
		Min    float64 `json:"min"`
		Max    float64 `json:"max"`
		Actual float64 `json:"actual"`
	} `json:"power"`
	TotPosEnergyPA float64 `json:"tot_posEnergyPA"`
	TotNegEnergyPA int     `json:"tot_negEnergyPA"`
	Timestamp      int64   `json:"timestamp"`
}

type TooSend struct {
	Power    []Power
	Amperage []Amperage
	Voltage  []Voltage
}
