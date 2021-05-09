# go-trainer-api

This project manages source code for a webserver used to help clients schedule a time with a physical trainer

## Getting started

Boot the webserver

```shell
# from project root

go run cmd/http/main.go
```

## Notes

- added numerical `user_id` value to JSON file "database" (was ommited in original JSON)
- switched to numerical UUID for `trainer_id` upon INSERT; recalculating the index incrementally based on the existing max seemed like an undesirably hacky way to handle this
- no tests :(

## Endpoints

|   endpoint    |   query parameters / request body    |   expected statuses | description |
| --------------| --------------------- | ------------------- | ----------- |
| `GET /availability/{trainerId}` | `?starts_at`, `?ends_at` (timestamps in RFC3339 format) | `200`: availability <br /> `400`: invalid start/end times <br /> `404`: no avails | List of available 30m slots, during business hours, for the queried trainer |
| `GET /appointment/{trainerId}` | | `200`: appointments | List of existing appointments for the queried trainer |
| `POST /appointment` | `trainer_id (int)`, `user_id (int)`, `starts_at (RFC3339 string)`, `ends_at (RFC3339 string)` | `204`: successful reservation <br /> `404`: invalid times | Make a reservation with a trainer at the provided times |
