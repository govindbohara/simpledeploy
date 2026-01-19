# **SimpleDeploy (v0)**

> A lightweight, SSH-based deployment tool for static websites and
> Node/Express applications --- built from scratch in Go.

SimpleDeploy allows developers to deploy applications to a single Ubuntu
server with **one command**, using a simple YAML configuration file. It
automates build, packaging, uploading, release management, and nginx
routing while keeping the architecture close to real production systems.

------------------------------------------------------------------------

## 🚀 What SimpleDeploy Can Do (v0)

### ✅ Static Site Deployment

SimpleDeploy can:

-   Build static sites locally (`npm run build`, etc.)
-   Package artifacts into a `.tar.gz`
-   Upload them via SSH to a remote server
-   Create **release folders + atomic symlink switching**
-   Serve sites using **nginx**
-   Support **multiple static sites on one server**
-   Route traffic using hostnames (via nip.io or real DNS)

Example public URLs:

    http://timer.<server-ip>.nip.io
    http://portfolio.<server-ip>.nip.io

Each site gets its **own nginx config**, so deployments do **not
overwrite each other**.

------------------------------------------------------------------------

### 🟡 Node / Express Deployment (v0 in progress → core working model)

SimpleDeploy supports Node apps using:

-   Release-based deployment under:

```{=html}
<!-- -->
```
    /home/ubuntu/simpledeploy/apps/<app>/

-   Server-side dependency install:

```{=html}
<!-- -->
```
    npm ci --omit=dev

-   Automatic **systemd service creation** for:
    -   auto-restart on crash\
    -   auto-start on reboot\
    -   proper logging\
-   nginx **reverse proxy routing** so many apps can share port 80:

```{=html}
<!-- -->
```
    api1.<ip>.nip.io → 127.0.0.1:3001
    api2.<ip>.nip.io → 127.0.0.1:3002

------------------------------------------------------------------------

## 🧠 Architecture Overview

    Your Laptop
    │
    │  simpledeploy deploy
    │
    ├─ Builds app locally
    ├─ Packages artifacts (.tar.gz)
    └─ Uploads to EC2 via SSH
          │
          ▼
    Ubuntu Server
    ────────────────────────────────────
    /var/www/<app>/releases/<id>   (static)
    /home/ubuntu/simpledeploy/...   (node)

    nginx  → routes by hostname
    systemd → runs Node apps

**Key ideas:** - One server, many apps\
- nginx handles public traffic\
- systemd keeps Node apps alive\
- releases + symlinks enable safe deployments

------------------------------------------------------------------------

## 📁 Example `simpledeploy.yaml` (Static Site)

``` yaml
app: timer
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
  install: "npm ci --omit=dev"
  start: "/usr/bin/node dist/server.js"

route:
  hostnames:
    - "api1.3.x.x.x.nip.io"
```

------------------------------------------------------------------------

## 🛠 One-Time Server Setup (Required)

On your Ubuntu EC2 instance:

``` bash
# Install Node (for Node apps)
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
sudo apt install -y nodejs

# Prepare nginx
sudo rm -f /etc/nginx/sites-enabled/default

# Allow SimpleDeploy limited sudo access
sudo tee /etc/sudoers.d/simpledeploy > /dev/null <<'EOF'
ubuntu ALL=(root) NOPASSWD: /bin/mkdir, /bin/chown, /bin/ln, /bin/mv, /usr/sbin/nginx, /bin/systemctl reload nginx
EOF

sudo chmod 440 /etc/sudoers.d/simpledeploy
sudo mkdir -p /var/www
sudo chown ubuntu:ubuntu /var/www
```

This setup is done **once per server**.

------------------------------------------------------------------------

## 🧪 How to Deploy

From your project folder:

``` bash
simpledeploy deploy
```

Output example:

    Uploading artifact...
    Extracting...
    Configuring systemd...
    Configuring nginx...
    ✅ Deployed successfully.

    Public URL: http://api1.3.x.x.x.nip.io/

------------------------------------------------------------------------

## 📌 What's Coming in v1

Planned improvements:

-   Rollback command\
-   Health checks before traffic switch\
-   Automatic port allocation\
-   HTTPS via Let's Encrypt\
-   Release cleanup (keep last N versions)\
-   `simpledeploy doctor` to detect IP mismatches

------------------------------------------------------------------------

## 🎯 Why This Project Matters

SimpleDeploy mirrors real production concepts used in platforms like:

-   Vercel\
-   Netlify\
-   Render\
-   AWS ECS\
-   Capistrano

Key learning outcomes:

-   SSH automation\
-   nginx routing\
-   systemd services\
-   atomic deployments\
-   infrastructure as code

------------------------------------------------------------------------

## 👨‍💻 Built With

-   Go\
-   SSH + SFTP\
-   nginx\
-   systemd\
-   Ubuntu\
-   Node.js

------------------------------------------------------------------------

## 📄 License

MIT
