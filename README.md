# gear-weight-api

Mainly been using it to pick the lighest trad rack possible for an objective.

Docker commands

```sh
docker compose build <service_name>
docker compose up <service_name> -d
```


## Seed data templates
### Cams (SLCD)
```json
{
  "type": "cam",
  "brand": <string>,
  "model": <string>,
  "sizes": [<array_of_strings>],
  "quantity": <int>
}
```

`quantity` is usually `1` unless you have a double rack of the same cams 

### Slings
```json
{
  "type": "sling",
  "brand": <string>,
  "model": <string>,
  "lengthInCentimeters": <int>,
  "quantity": <int>
}
```

### Carabiners
```json
{
  "type": "carabiner",
  "brand": <string>,
  "model": <string>,
  "quantity": <int>
}
```

## To-do
### Fix database backend for dev service
This line does not work anymore:
`go install github.com/mattn/go-sqlite3@1.14.16`

For some reason version `1.14.16` isn't recognized. I noticed version `1.14.24` in `go.sum`, but that does not work either. Installing sqlite is important because it takes a long time to build `server.go` otherwise.

Unless that git repo starts working again in the future, it would be good idea to switch to another database solution.

Thankfully, the prod image still builds.

Would probably not use sqlite with golang moving forward.