# Telegraf HTTP Harness
A HTTP based harness for [Telegraf](https://github.com/influxdata/telegraf) monitoring system. it has an embedded executable version of Telegraf. Now you can simply start and stop your Telegraf remotely and gathering your system metrics and send them to your Influxdb instant.
## Usage
In windows, download the  [latest binary package from release section](https://github.com/mohammadGh/telegraf-harness/releases/download/v0.2/telegraph-harness-v0_2.zip) (or clone the project and build it with `go build` ), unzip it and then run the executable `telegraf-harness.exe`
now the harness listens to port `6663` for your http commands.

To test it works properly, use http request `/test`; this command run the telgraf with `--test` arguments and responds the result which contains cpu, memory, disk and network metrics:
    
    curl http:://your-telegraf-harness:6663/test  
To start the Telegraf, just send the following http request (using `curl` or just enter the request in your browser):
    
    curl http:://your-telegraf-harness:6663/start?interval=5s&database=myDataBaseName&url=http://my-influxdb-server-address:8086  
To stop the `telegraf` just send the following request to our harness:
    
    curl http:://your-telegraf-harness:6663/stop  
