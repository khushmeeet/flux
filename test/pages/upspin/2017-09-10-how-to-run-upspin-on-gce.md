---
template: post
title: How to run upspin on GCE
shortie: Tutorial on how to setup upspinserver on Google Compute Engine
date: 2017-09-10
categories: data
tags:
  Upspin
  File sharing
---

This tutorial is going to be a step by step guide on how to setup upspinserver-gcp on Google Compute
Engine. This is not a tutorial on upspin and its architecture. For more on upspin you can go to
[upspin.io](https://upspin.io/)

**Step 1 -** Download the upspin binaries from [here](https://upspin.io/dl/).
  Download the binaries according to your local machine.

**Step 2 -** Download Go programming language from [Golang](https://golang.org/).
  We need this for compiling our server binary. We need GOPATH for our packages.
  Add following line in your `.bash_profile` or `.zshrc` or any other shell you are using.
  ```bash
  export GOPATH=~/some_folder/go
  export PATH=$PATH:$GOPATH/bin
  ```

**Step 3 -** This is upspin signup process. We will use our email to register. Our email will
  act as our identity in upspin namespace. Do
  ```bash
  local$ upspin signup -server=upspin.your-domain.com you@gmail.com
  ```
  Use -`-server` since both dir and store is going to be hosted on GCE instance.
  `upspin.your-domain.com` will point to that GCE instance.
  It will spit config file, storing the dir and store reference among other things
  and seed key in case you lost your private key. Write down your seed key and keep it safe.

**Step 4 -** Now we need to setup our domain. Do
  ```bash
   local$ upspin setupdomain -domain=your-domain.com
   ```
  It will spit out the record that you need to put in your domain's DNS zone.
  Here's how you gonna do that.

  ![DNS Record]({{site.url}}/assets/dns.png){:style="max-width: -webkit-fill-available"}
  
  My registrat is cloudflare. UI could be different, but the record remains the same.
  ```bash
   local$ host -t TXT your-domain.com
  ```
  Use host utility to check if changes has propagated.

**Step 5 -** Compile upspinserver-gcp to deploy on GCE. This should be compiled for linux, since
  we are using ubuntu 16.04 LTS for GCE instance.
  ```bash
   local$ go get -d gcp.upspin.io/cmd/...
   local$ go install gcp.upspin.io/cmd/upspin-setupstorage-gcp
   local$ GOOS=linux GOARCH=amd64 go build gcp.upspin.io/cmd/upspinserver-gcp
  ```
  Create a google cloud project (standard stuff) and download sdk from [here](https://cloud.google.com/sdk/downloads).
  Now we need to enable certain APIs for our project.
  ```bash
   local$ gcloud components install beta
   local$ gcloud config set project example-com
   local$ gcloud auth login
   local$ gcloud beta service-management enable iam.googleapis.com
   local$ gcloud beta service-management enable storage_api
  ```
  Authenticate yourself.
  ```bash
   local$ gcloud auth application-default login
  ```
  Now create google storage bucket and its service account. This is going to used, when creating our VM instance.
  ```bash
   local$ upspin setupstorage-gcp -domain=your-domain.com -project=project-id bucket-name
  ```

**Step 6 -** Create a GCE VM. Use micro or small machine type. Select `Allow HTTPS traffic`.
  Make sure external IP is enabled, reserved static address (important).
  Generate a ssh key
  ```bash
   ssh-keygen -t rsa -f ~/.ssh/[KEY_FILENAME] -C [USERNAME]
   ```
  Add the public key of this key pair to your VM. This lets you access VM from your terminal using SSH.

**Step 7 -** Create a `A` DNS record. Your subdomain pointing to VM (external IP).
  ![A DNS]({{site.url}}/assets/adns.png){:style="max-width: -webkit-fill-available"}

**Step 8-** With domain in place. Login to your VM.
  Do `passwd` to set new UNIX password.
  Then add your public key that you generated in step 6 to this VM user.
  ```bash
   serve$ mkdir .ssh
   server$ chmod 0700 .ssh
   server$ touch .ssh/authorized_keys
   server$ chmod 0600 .ssh/authorized_keys
   server$ cat > .ssh/authorized_keys
   (Paste your SSH public key here and type Control-D and Enter)
  ```

**Step 9 -** Copy your compiled for linux `upspinserver-gcp` binary to your VM
  ```bash
   local$ scp -i ~/.ssh/ssh_key upspinserver-gcp user@upspin.your-domain.com:upspinserver
  ```
  `ssh_key` is your private key for which you uploaded its corresponding public key.

**Step 10 -** Create an upspin service to let it run on the background and also start itself on boot.
  Do this as root `sudo su`.
  Create the file `/etc/systemd/system/upspinserver.service` that contains the following service definition.
  ```bash
   [Unit]
   Description=Upspin server

   [Service]
   ExecStart=/home/upspin/upspinserver
   User=upspin
   Group=upspin
   Restart=on-failure

   [Install]
   WantedBy=multi-user.target
  ```

**Step 11 -** Bind the service to port `443`. Do as root.
  ```bash
   server% setcap cap_net_bind_service=+ep /home/user/upspinserver
  ```
  User is you.

**Step 12 -** Start the service. Do as root.
  ```bash
   server% systemctl enable --now /etc/systemd/system/upspinserver.service
  ```
  You may also use systemctl stop upspinserver and systemctl restart upspinserver to stop and restart the server, respectively.
  Check the log of the server.
  ```bash
   server% journalctl -f -u upspinserver
  ```
  Logs should output `Configuration file not found. Running in setup mode.`.
  And if you open `https://upspin.your-domain.com/` you should see
  `Unconfigured Upspin Server`. Which means its running.

**Step 13 -** On your local machine, run this command to send server keys to config to upspin server
  instance.
  ```bash
   local$ upspin setupserver -domain=your-domain.com -host=upspin.your-domain.com
  ```

**Step 14 -** Test your server by running
  ```bash
   local$ echo Hello, Upspin | upspin put you@gmail.com/hello
   local$ upspin get you@gmail.com/hello
   Hello, Upspin
  ```
  
  You can mount upspin as virtual file system by using osxfuse(FUSE filesystem) on macOS and similar program on linux. Not sure about windows. `upspinfs` lets you mount upspin namespace.
  ```bash
   local$ mkdir $HOME/up
   local$ upspinfs $HOME/up
   local$ ls $HOME/up/you@gmail.com
  ```

  Or you can use upspin [browser](https://github.com/jnglco/browser) built on electron. 