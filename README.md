# Workshop
* Load balance with HAProxy
* WebSocket service with Go and Gorilla websocket
* Redis to keep data
  * Set
  * Pub/Sub

## Start all services
```
$docker compose build
$docker compose up -d
$docker compose ps
NAME                COMMAND                  SERVICE             STATUS              PORTS
demo-lb-1           "docker-entrypoint.s…"   lb                  running (healthy)   0.0.0.0:8080->8080/tcp
demo-redis-1        "docker-entrypoint.s…"   redis               running (healthy)   0.0.0.0:6379->6379/tcp
demo-ws1-1          "./app"                  ws1                 running (healthy)   8080/tcp
demo-ws2-1          "./app"                  ws2                 running (healthy)   8080/tcp

$docker compose logs --follow
```

## Test WebSocker with [wscat](https://www.npmjs.com/package/wscat)
```
$wscat -c ws://localhost:8080/ws/1
$wscat -c ws://localhost:8080/ws/2
```
