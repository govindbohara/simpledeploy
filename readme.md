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

    Ubuntu server
    /var/www/<app>/releases/<id>   (static)
    /home/ubuntu/simpledeploy/appName  (node)

    nginx  → routes by domain name
    systemd → process manager and keeps application alive

## One time server setup

# Install Node (for Node apps)

curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
sudo apt install -y nodejs

# Prepare nginx

sudo rm -f /etc/nginx/sites-enabled/default

# Allow SimpleDeploy limited sudo access

sudo tee /etc/sudoers.d/simpledeploy > /dev/null <<'EOF'
ubuntu ALL=(root) NOPASSWD: \
/bin/mkdir, \
/bin/chown, \
/bin/ln, \
/bin/mv, \
/usr/sbin/nginx, \
/bin/systemctl reload nginx
EOF

sudo chmod 440 /etc/sudoers.d/simpledeploy
sudo mkdir -p /var/www
sudo chown ubuntu:ubuntu /var/www

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
    - "src"
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
