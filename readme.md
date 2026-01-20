# **SimpleDeploy**

> A lightweight, SSH-based deployment tool for static websites and
> Node/Express applications built from scratch in Go.

## Features of Simple Deploy

### Static Site Deployment

SimpleDeploy can:

- Build static sites locally
- Package artifacts
- Upload them via SSH to a remote server
- Serve sites using **nginx**
- Support **multiple static sites on one server**
- Route traffic using domain name

Example public URLs:

    http://timer.<public-ip>.nip.io
    http://portfolio.<public-ip>.nip.io

## Architecture Overview

    simpledeploy deploy
    Builds app locally
    Packages artifacts (.tar.gz)
    Uploads to EC2 via SSH
         
    Ubuntu Server
    ────────────────────────────────────
    /var/www/<app>/releases/<id>   (static)
    /home/ubuntu/simpledeploy/appName  (node)

    nginx  → routes by domain name
    systemd → process manager and keeps application alive

## Example `simpledeploy.yaml` (Static Site)

```yaml
app: Appname
type: static

target:
  host: "3.x.x.x"
  user: "ubuntu"
  port: 22
  keyPath: "~/.ssh/key.pem"

build:
  local:
    - "npm ci"
    - "npm run build"

static:
  webRoot: "/var/www/timer"
  distDir: "dist"

package:
  include:
    - "dist"

route:
  hostnames:
    - "timer.3.x.x.x.nip.io"
```

---

## Example `simpledeploy.yaml` (Node / Express)

```yaml
app: Appname
type: node

target:
  host: "3.x.x.x"
  user: "ubuntu"
  port: 22
  keyPath: "~/.ssh/key.pem"

build:
  local:
    - "npm ci"
    - "npm run build"

package:
  include:
    - "dist"
    - "package.json"
    - "package-lock.json"

node:
  port: 3001
  install: "npm ci"
  start: "node server.js" // "npm run start:prod"

route:
  hostnames:
    - "api.3.x.x.x.nip.io"
```
