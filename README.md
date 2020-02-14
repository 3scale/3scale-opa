# 3scale OPA

[Open policy agent](https://www.openpolicyagent.org) built for 3scale.

Not much to see here yet. Just getting started.

## Instructions

- Compile:

```bash
go build
```

- Start the OPA server with the rules in `example.rego`:

```bash
./3scale-opa run --server --set=plugins.envoy_ext_authz_grpc.addr=:9191 --set=plugins.envoy_ext_authz_grpc.query=data.envoy.authz.allow --set=decision_logs.console=true --ignore=.* example.rego
```

- Start Envoy with the given config file that authorizes using the server above:

```bash
envoy -c envoy_config.yaml
```

- The `example.rego` policy defines a limit of 5 requests per second for the
"/abc" path, but each user identified by the "user_id" header can only make 3
requests per second. To test this, you'll need to start a server in 8080 or choose another one, but remember to change it in `envoy_config.yaml` as well. For testing purposes you can use: `python -m http.server 8080`. Now make requests to Envoy:
    - `curl -v http://localhost:8000/abc -H "user_id:a"`. Authorized: total
    (1/5), user_a (1/3), user_b (0/3).
    - `curl -v http://localhost:8000/abc -H "user_id:a"`. Authorized: total
    (2/5), user_a (2/3), user_b (0/3).
    - `curl -v http://localhost:8000/abc -H "user_id:a"`. Authorized: total
    (3/5), user_a (3/3), user_b (0/3).
    - `curl -v http://localhost:8000/abc -H "user_id:a"`. Limited: total (3/5),
    **user_a (4/3)**, user_b (0/3).
    - `curl -v http://localhost:8000/abc -H "user_id:b"`. Authorized: total
    (4/5), user_a (3/3), user_b (1/3).
    - `curl -v http://localhost:8000/abc -H "user_id:b"`. Authorized: total
    (5/5), user_a (3/3), user_b (2/3).
    - `curl -v http://localhost:8000/abc -H "user_id:b"`: Limited: **total
    (6/5)**, user_a (3/3), user_b (3/3).
