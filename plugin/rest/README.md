# CoreDNS REST Forward Plugin

This plugin is used to forward DNS requests to a RESTful API. 
When a DNS request is made, the plugin will forward the request to the RESTful
API specified in the `url` directive which returns the response in JSON format.
The plugin then constructs a DNS response from the JSON response
and returns it.

It will further handle the DNSSEC signing of the response if the `dnssec` plugin
is enabled in the same server block.

## Configuration of Corefile

```Corefile
example.com {
  rest http://localhost:8080/api/v1/
}
```

This means that any DNS request to `example.com` will be forwarded to the
RESTful API at `http://localhost:8080/api/v1/<domain>` and the response will be used
to construct a DNS response.

## Example of RESTful API

The following is an example of a RESTful API that returns a JSON response.

```bash
curl http://localhost:8080/api/v1/example.com
```

```json
{
  "RCODE": 0,              // DNS RCODE (0 = NOERROR)
  "AD": true,               // DNSSEC Authenticated Data flag
  "Answer": [
    {
      "name": "example.com.",
      "type": "A",
      "TTL": 3600,
      "data": "192.168.1.1"
    }
  ],
  "Question": [             // Echo original query
    {
      "name": "example.com.",
      "type": "A"
    }
  ]
}
```

## Configuration of the plugin

The plugin can be configured with the following directives:

```Corefile
rest <url> {
  <key> <value>
}
```

The `<url>` is the URL of the RESTful API.
The `<key> <value>` pairs are optional parameters that can be used to configure the plugin.

<!-- The following keys are supported:
- `dnssec`: Enable DNSSEC signing of the response. Default is `false`. -->









