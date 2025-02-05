package main

type ServerConf struct {
	Name      string `json:"name"`
	Scheme    string `json:"scheme"`
	Host      string `json:"host"`
	Port      int    `json:"port"`
	ProxyPort int    `json:"proxyPort"`
}

type ProxyConf []ServerConf

type CustomLogger func(mess string, val ...any)
type ConfProvider func() ServerConf
