#### Custom Error using the `error` interface
----
We are trying to develop a custom error object graph that will be more aligned to Go Gin framework's http response and also digestable by the log. To keep it compatible with older functions beneath the handlers it also implements (extends) the error interface. This helps in keeping up with the older function interface while using the `Throw` method in sending out the new error object.

#### Error correspondence to HTTP status code
---
This makes the error objects convertible to http responses. 

#### Logging the error and piggybacking on GoGin context
----
BEfore the errors are dispatched in http.Responses the appropriate HTTP status codes are important. And such that error codes and http responses will have a direct correspondence.