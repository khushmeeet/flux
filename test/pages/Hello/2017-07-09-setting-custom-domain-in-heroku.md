---
template: post
title: Setting custom domain in heroku
shortie: Few steps on how to setup domain to work with your heroku app
date: 2017-07-09
categories: web
tags:
  - heroku
  - godaddy
---

For adding a domain to your heroku app, process is pretty simple. For this tutorial, I will be using khushmeetsingh.com (my domain) as an example, which I bought from GoDaddy.

Here are few steps you need to follow, to get this going.

**Step 1 -** Add credit/debit card to your heroku app. No charges will be made. This is just for verification. You need to do this before you setup custom domain.

**Step 2 -** Go to your domain registrar and open **DNS Settings**. There is a list of records. Just create a new record of type **CNAME** with alias `www` and host name `proxy.heroku.com`.

**Step 3 -** To verify your previous step, do `ping www.domain.com`. This should give `proxy.heroku.com` for ping.

**Step 4 -** If that is successful, its time to add domain to heroku. Your do this either in the command line of the dashboard. Here is the command of that `heroku domains:add www.domain.com --app herokuapp`. This added domain for `www` but we also want our site to open when we enter domain without `www` like this `domain.com`.

For this enter the same command again, but this time do `heroku domains:add domain.com --app herokuapp`.

**Step 5 -** Visit your domain registrar (last time). Under domain forwarding(domain), add a domain. Add forward to `www.domain.com`, keep `http://` as is, unless you have SSL certificate.

Keep direct type as `301 (Permanent)`. Forward setting as `forward only`.
Check `update nameservers and DNS settings`.

That's it, now your domain is setup. Run domain.com and you will be greeted with heroku app.
