# Order Management Service

This is a minimal service which accepts orders in a certain format through the REST API and provides listing methods for order items and summary.

## Testing

Unit tests can be run by issuing following command

```bash
make unit
```

## Build and Run

This service stores data in program memory and doesnt require any external data store.


Service uses a config file inside:

```bash
./config/config.yaml
```

with the content
```yaml
cache_expiry_duration: 24
cache_cleanup_interval: 120
api_default_timeout: 5
server_port: 8000
```
To build the service run

```bash
make start
```

Which in turn generate the binary of the service at the folder root as "order-service" and run it

## Open API documentation
API documentation at /docs/openapi.yaml location can be viewed with an openapi viewer

## License
None