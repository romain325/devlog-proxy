# Log Proxy

Proxy all your request from one server to another with way to much log so you don't miss information about your services requests 

**THIS TOOL IS ONLY MEANT FOR DEBUG PURPOSE AND SHOULD NOT BE USED IS PROD**

## How to use

Edit the file conf.json

```json
[
  {
    "name": "petstore",    // identifier that will be used in the logs
    "proxyPort": 3333,     // port that will be listening on this proxy
    "port": 1080,          // destination port
    "host": "localhost",   // destination host
    "scheme": "http"       
  },
  
  // you can run as much configuration as you want
  {
    "name": "petstore2",
    "proxyPort": 3335,
    "port": 1080,
    "host": "localhost",
    "scheme": "http"
  }
]
```

and then

```bash
go run main.go
```