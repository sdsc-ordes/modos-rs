# Test Services

## Start Services

```bash
just services-start
```

> [!NOTE]
>
> `Ctrl+D` to detach from the TUI and `just services-attach` to attach again.

## Keycloak

Keycloak has one realm `modos` where all clients and settings live. The
configuration file for it is located in
`tools/configs/keycloak/modos-realm.conf`.

### Test Users

Keycloak has the following test users:

- Email: `test@mail.ch`, Password: `test`.

### Account UI

Log in under
[http://localhost:8081/realms/modos/account](http://localhost:8081/realms/modos/account)

### Admin UI

Log in under
[http://localhost:8081/admin/master/console/#/modos](http://localhost:8081/admin/master/console/#/modos)
with name `admin` and password `admin`.

### Export Settings

Once you made changes in the UI you can stop keycloak and export the realm
`modos` with.
