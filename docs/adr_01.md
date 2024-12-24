## Number: 01
## Date: 2024-12-22
## Title: Architecture of Data Processing

In this context, Processors are defined as reading some input from a file, calling 
the [speedrun.com](https://www.speedrun.com) API, and writing some output to a file. 
This structure defines clear separation between where we read and write. This is not 
much of a problem, but becomes increasingly confusing when we write to multiple files. 

The architecture of processing different data can be seen below:
                                                                                         
```                                                                               
                                       ┌───────────────────┐    ┌───────────────────┐ 
                                       │                   │    │                   │ 
                                   ┌───►  Write Repository ┼────►    CSV Writer     │ 
                                   │   │                   │    │                   │ 
                                   │   └───────────────────┘    └───────────────────┘ 
                                   │                                                  
            ┌─────────────────┐    │   ┌───────────────────┐    ┌───────────────────┐ 
            │                 │    │   │                   │    │                   │ 
            │    Processor    ┼────┼───►  Read Repository  ┼────►    CSV Reader     │ 
            │                 │    │   │                   │    │                   │ 
            └────────▲────────┘    │   └───────────────────┘    └───────────────────┘ 
                     │             │                                                  
                     │             │   ┌───────────────────┐    ┌───────────────────┐ 
                     │             │   │                   │    │                   │ 
                     │             └───►   SRCOM Client    ┼────►    HTTP Client    │ 
 ┌─────────────────┐ │                 │                   │    │                   │ 
 │   Games  Data   ┼─┤                 └───────────────────┘    └───────────────────┘ 
 │    Processor    │ │                                                                
 └─────────────────┘ │                                                                
                     │                                                                
 ┌─────────────────┐ │                                                                
 │   Games  List   ┼─┤                                                                
 │    Processor    │ │                                                                
 └─────────────────┘ │                                                                
                     │                                                                
 ┌─────────────────┐ │                                                                
 │   Runs   Data   ┼─┘                                                                
 │    Processor    │                                                                  
 └─────────────────┘                                                                  
```
(This diagram was generated using [asciiflow.com](https://asciiflow.com/#/).)

### New Functionality

Want to do something new? Follow the list below on how to implement the new behaviour:

* To retrieve new data from the API— add function(s) to [speedrun.com Client](../internal/srcom_api/v1.go); 
* To write data to a new file— add the filename, comment, and header to [`constants.go`](../internal/repository/constants.go);
* To process new data— create a new executable in [`cmd`](../cmd/), and add a processor with a matching name in [`processor`](../internal/processor/);
* To describe new functionality— create a new ADR using [`adr_xx.md`](./adr_xx.md);