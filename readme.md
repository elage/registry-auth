# Docker Auth

Docker Registry with Authentication and Authorization

This project sets up a Docker registry with authentication and authorization using Casbin and htpasswd.

## Prerequisites

- Docker
- Docker Compose

## Setup

Create a file `model.conf` that contains the casbin model:

```
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = (p.sub == '*' || g(r.sub, p.sub)) && keyMatch(r.obj, p.obj) && regexMatch(r.act, p.act)
```

Create a file `policy.csv` that contains the casbin policy:

```
p, bob, /v2/foo/*, GET|HEAD
p, developers, /v2/*, POST|PATCH|PUT|HEAD
p, *, /v2/hello-world/*, GET|HEAD
g, bob, developers
```

Create a file `users.htpasswd` with a htpasswd formatted list of usernames and password hashes.

```
foo:$2y$05$BocMB//m0IAn4yJ8kaZLi.9wox5O4561iJo6KK4BI3W7.ccthVbQW
```

Create an `.env` file with your `HOST` set.

```
HOST=registry.example.com
```
