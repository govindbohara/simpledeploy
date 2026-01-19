# Simple Deploy

CLI to deploy applications to a server using a simple yaml config file.

## 📁 Example `simpledeploy.yaml` (Static Site)

``` yaml
app: Appname
type: static

target:
  host: "3.x.x.x"
  user: "ubuntu"
  port: 22
  keyPath: "~/.ssh/simpledeploy-key.pem"

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

------------------------------------------------------------------------

## 📁 Example `simpledeploy.yaml` (Node / Express)

``` yaml
app: api1
type: node

target:
  host: "3.x.x.x"
  user: "ubuntu"
  port: 22
  keyPath: "~/.ssh/simpledeploy-key.pem"

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
  start: "node server.js" //npm run start:prod

route:
  hostnames:
    - "api1.3.x.x.x.nip.io"
```