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