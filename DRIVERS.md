# Drivers

Right now, I've only implemented a driver for Redis.
Fork it and feel free to develop other drivers (etcd, memcache, custom apis...).

## How it works

The driver shall implement the interface below. Simple like that.

```
type Driver interface {
  Start(name, url string) string
  Get() (map[string]string, error)
  Delete(currentName string)
}
```

### Methods and their returns

Here are the details of each method:

- Start(name, url string) string:
  - receives the service name (***"my-service-name"*** for example) and its URL (***https://...***).
  - responds an unique name for that service, note that it shall be formed with `<SERVICE_NAME>-<UNIQUE_ID>` for example: ***my-service-name-20160315122936.17256541***


- Get() (map[string]string, error)
  - responds with a key list of services unique names and their respective URLs, for example:

    ```
    provider-20160315122936.17256541  |  http://localhost:3333
    provider-20160315122713.15645582  |  http://localhost:3334
    consumer-20160315122543.15645582  |  http://localhost:3335
    ```

- Delete(currentName string)
  - receives a service unique name (like ***provider-20160315122936.17256541***) and removes it of it's database. The next time you run the command ***Get()*** that removed service data will not be retrieved.
