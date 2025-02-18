# Medicine logger

## Links

https://datatables.net/

## Configuration

| Setting | Environment variable | Default | Description |
| ------- | -------------------- | ------- | ----------- |
| DBString | FBF_DBPATH | medicine-logger.sqlite3 | The path to the database file, or "file::memory:?cache=shared" for an in-memory DB. |
| JWTSecret | FBF_JWTSECRET | very_very_secret | The secret used to sign JWTs. |
| Address | FBF_ADDRESS | :3019 | The address to listen on. |
| AutoTLS | FBF_AUTOTLS | false | Whether to automatically enable TLS. CURRENTLY NON-FUNCTIONAL |
| AccountCreation | FBF_ACCOUNTCREATION | true | Whether to allow account creation. |
| EmailHost | FBF_EMAILHOST | smtp.sendgrid.net | The host to use for sending emails. |
| EmailPort | FBF_EMAILPORT | 587 | The port to use for sending emails. |
| EmailUsername | FBF_EMAILUSERNAME | apikey | The username to use for sending emails. |
| EmailPassword | FBF_EMAILPASSWORD | | The password to use for sending emails. |
| EmailFrom | FBF_EMAILFROM | friend@5x5.com | The "From" address for sending emails. |

Non-default settings can be set by environment variables, or by a `config.toml` file:

```toml
DBString = 'medicine-logger.sqlite3'
JWTSecret = 'very_very_secret'
Address = ':3019'
AutoTLS = false
AccountCreation = true

```
