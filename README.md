# Telegraf HTTP Harness
A HTTP based harness for [Telegraf](https://github.com/influxdata/telegraf) monitoring system. It has an embedded executable version of Telegraf. Now you can simply start and stop your Telegraf remotely and gathering your system metrics and send them to your Influxdb instant.
## Usage
In windows, download the  [latest binary package from release section](https://github.com/mohammadGh/telegraf-harness/releases/download/v0.2/telegraph-harness-v0_2.zip) (or clone the project and build it with `go build` ), unzip it and then run the executable `telegraf-harness.exe`. Now the harness listens to port `6663` for your http commands.

To test it works properly, use http request `/test`; this command run the Telegraf with `--test` arguments and responds the result which contains cpu, memory, disk and network metrics:
    
    curl http://your-telegraf-harness:6663/test  
To start the Telegraf, just send the following http request (using `curl` or just enter the request in your browser):
    
    curl http://your-telegraf-harness:6663/start?interval=5s&database=myDataBaseName&url=http://my-influxdb-server-address:8086  
To stop the `telegraf` just send the following request to our harness:
    
    curl http://your-telegraf-harness:6663/stop  
## Default Configuration File
Currently the default configuration file monitor cpu, memory, disk and network for windows machine.
   
   ###############################################################################
#                            INPUT PLUGINS                                    #
###############################################################################

[[inputs.win_perf_counters.object]]
    # Processor usage, alternative to native, reports on a per core.
    ObjectName = "Processor"
    Instances = ["*"]
    Counters = ["% Idle Time", "% Interrupt Time", "% Privileged Time", "% User Time", "% Processor Time"]
    Measurement = "win_cpu"
    IncludeTotal=true

  [[inputs.win_perf_counters.object]]
    # Disk times and queues
    ObjectName = "LogicalDisk"
    Instances = ["*"]
    Counters = ["% Idle Time", "% Disk Time","% Disk Read Time", "% Disk Write Time", "% User Time", "Current Disk Queue Length"]
    Measurement = "win_disk"
    #IncludeTotal=false #Set to true to include _Total instance when querying for all (*).

  [[inputs.win_perf_counters.object]]
    ObjectName = "System"
    Counters = ["Context Switches/sec","System Calls/sec", "Processor Queue Length"]
    Instances = ["------"]
    Measurement = "win_system"
    #IncludeTotal=false #Set to true to include _Total instance when querying for all (*).

  [[inputs.win_perf_counters.object]]
    # Example query where the Instance portion must be removed to get data back, such as from the Memory object.
    ObjectName = "Memory"
    Counters = ["Available Bytes","Cache Faults/sec","Demand Zero Faults/sec","Page Faults/sec","Pages/sec","Transition Faults/sec","Pool Nonpaged Bytes","Pool Paged Bytes"]
    Instances = ["------"] # Use 6 x - to remove the Instance bit from the query.
    Measurement = "win_mem"
    #IncludeTotal=false #Set to true to include _Total instance when querying for all (*).

  [[inputs.win_perf_counters.object]]
    # more counters for the Network Interface Object can be found at
    # https://msdn.microsoft.com/en-us/library/ms803962.aspx
    ObjectName = "Network Interface"
    Counters = ["Bytes Received/sec","Bytes Sent/sec","Packets Received/sec","Packets Sent/sec"]
    Instances = ["*"] # Use 6 x - to remove the Instance bit from the query.
    Measurement = "win_net"
    IncludeTotal=true #Set to true to include _Total instance when querying for all (*).
