## Pulsar Bridge
Apache Pulsar http bridge gives a simple api to produce messages to multiple topics with single transaction .

## Configurations
Following environment variables can be set to configure .

| VARIABLE                     | DESCRIPTION                                                                                                             |
|------------------------------|-------------------------------------------------------------------------------------------------------------------------|
| PORT                         | Server port to listen                                                                                                   |
| DISABLE_PULSAR_LOGS          | Disable logging for pulsar client if set to "yes"                                                                       |
| LOG_LEVEL                    | [zerolog](https://github.com/rs/zerolog/blob/7fa45a4dda359d9a657a2960078097415417ec73/log.go#L127 "zerolog") Log levels |
| PULSAR_CLIENT_TIMEOUT_SEC    | Timeout for pulsar client                                                                                               |
| PULSAR_OPERATION_TIMEOUT_SEC | pulsar publisher execution timeout                                                                                      |
| PULSAR_BROKERS               | List of pulsar brokers comma seperated e.g: ip:port[,ip:port]                                                           |
| TNX_TIMEOUT_SEC              | Transaction timeout                                                                                                     |
| MAX_PENDING_MSGS             | Max queue size for publisher                                                                                            |
| BATCHING_MAX_SIZE            | Max bytes permitted in a batch                                                                                          |

## Usage

    ~$ git clone git@github.com:arslanmughal99/pulsar-bridge.git PulsarBridge
	~$ cd PulsarBridge
    ~$ go run .
server will start listening on configured port exposing `/produce` endpoint .

##### Request Example
```bash
curl --location 'localhost:3000/produce' --header 'Content-Type: application/json' --data '{
    "data": [
        {
            "topic": "topic-1",
            "messages": [
                "Test data 1",
                "Test data 2",
                ...
            ]
        },
        {
            "topic": "topic-2",
            "messages": [
                "Test data 1",
                "Test data 2",
                ...
            ]
        }
    ]
}'
```
##### Success Response Example (Status 200 OK)

```json
{
	"total": 4
}
```