# tracksrv
Simple web service that counts views.

`docker pull sergeyfast/tracksrv`

## Flags
```
  -addr string
        address listen to (default ":8090")
  -types string
        allowed types (default "item,news")
```
## REST API

`/<type>/<id>` – register view. `?data=1` returns new count.

`/pop` – get all data in JSON format. `?keep=1` won't clear data.

