# Reverse Proxy Support

Scrutiny is designed so that it can be used with a reverse proxy, leveraging `domain`, `port` or `path` based matching to correctly route to the Scrutiny service.

For simple `domain` and/or `port` based routing, this is easy.

If your domain:port pair is similar to `http://scrutiny.example.com` or `http://localhost:54321`, just update your reverse proxy configuration
to route traffic to the Scrutiny backend, which is listening on `0.0.0.0:8080` by default.

```yaml
# default config
web:
  listen:
    port: 8080
    host: 0.0.0.0
```

However if you're using `path` based routing to differentiate your reverse proxy protected services, things become more complicated.

If you'd like to access Scrutiny using a path like: `http://example.com/scrutiny/`, then we need a way to configure Scrutiny so that it
understands `http://example.com/scrutiny/api/health` actually means `http://localhost:8080/api/health`.

Thankfully this can be done by changing **two** settings (both are required).

1. The webserver has a `web.listen.basepath` key
2. The collectors have a `api.endpoint` key.

## Webserver Configuration

When setting the `web.listen.basepath` key in the web config file, make sure the `basepath` key is prefixed with `/`.

```yaml
# customized webserver config
web:
  listen:
    port: 8080
    host: 0.0.0.0
    # if you're using a reverse proxy like apache/nginx, you can override this value to serve scrutiny on a subpath.
    # eg. http://example.com/custombasepath/* vs http://example.com:8080
    basepath: '/custombasepath'
```

## Collector Configuration 

Here's how you can update the collector `api.endpoint` key:

```yaml
# customized collector config
api:
  endpoint: 'http://localhost:8080/custombasepath'
```

# Environmental Variables.

You may also configure these values using the following environmental variables (both are required).

- `COLLECTOR_API_ENDPOINT=http://localhost:8080/custombasepath`
- `SCRUTINY_WEB_LISTEN_BASEPATH=/custombasepath`