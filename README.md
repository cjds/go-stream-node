# GO Stream Node (WIP)
A ROS streams node built in GO that streams data to streams server. Data can be published to various channels and the stream node publishes it to the streams server which is eventually stored in Cassandra database. This a work in progress project and more updates will be added.

## Getting Started
### Installing Dependencies
Go stream node runs directly on the host machine which in this case is the robot. Assuming that `go-stream-node` is now in your host machine's `$GOPATH/src` directory and you have a stream server running on your local machine

Run the following command to install Go dependencies
```
cd $GOPATH/src/go-stream-node
dep ensure
```
This should install all the dependencies required to run go stream node.

In addition to above dependencies, we need `gengo` to automatically generate GO files for message types. Use the following commmand to get gengo.
```
go get github.com/akio/rosgo/gengo
```
Go to your `$GOPATH/src/github.com/akio/rosgo/gengo` and run `go install`. This should add the `gengo` executable to your `$GOPATH/bin`. Now run the following command to automatically generate GO files for message types from the message definitions
```
go generate
```
### Starting Stream Node
Now, run the following command to start the stream node server. Be sure to change streams server host address in `application.toml` to wherever your streams server is running if you run across `connection refused` error.
```
go install && go-stream-node
```
To add new listeners, following format can be used in `main.go` which will add new stream type to `streams` map. This will add additional listeners to the stream node.
```
streams := map[string]ros.MessageType{
                "<TopicName>":<MessageType>}
```
For example we can add new listener in the following way
```
streams := map[string]ros.MessageType{
		"string":  std_msgs.MsgString,
		"battery": power_msgs.MsgBatteryState}
```
To send messages from your host machine to go stream node use following format. Make sure that `roscore` is up and running before running `rostopic` to send messages.
```
rostopic pub /<TopicName>  <MessageType>  <Message>
```
For example we can send message of type `std_msgs/String` using the following command
```
rostopic pub /string std_msgs/String "data:'Test'"
```
Do note that `std_msgs/String` type data is not supported on streams server. Hence, you will probably see an error on your streams server.

Tests can be run using the following command. Make sure roscore is up and running before running the test.
```
go test -test.v=true
```

### Configurations
Changing auth0 credentials and streams server's properties can be done by editing the `application.toml` file in `conf/development` directory.
