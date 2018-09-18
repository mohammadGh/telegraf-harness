# Telegraf Rest Harness
A HTTP based harness for Telegraf monitoring system. it has an embedded executable version of `Telegraf .`
## Usage
In windows, download the binary package (or clone the project and build it with `go build` ), unzip it and then run the executable `telegraf-harness.exe`
now the harness listens to port `6663` for your http commands.
for exmaple to start the follwoing http request (using `curl` or just enter the request in your browser):
    
    curl yourHarnessMachineAddress:6663/start?interval=5s&database=myDataBaseName&url=http://my-influxdb-server-address:8086  
To stop the `telegraf` just send the follwing request to our harness:
    
    curl yourHarnessMachineAddress:6663/stop  
