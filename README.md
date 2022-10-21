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
$docker compose logs --follow
```

## Test WebSocker with [wscat](https://www.npmjs.com/package/wscat)
```
$wscat -c ws://localhost:8080/ws/user01
```