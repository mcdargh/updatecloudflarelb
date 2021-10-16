# Sync up your UNIFI wan1 and wan2 ips to your cloudflare load balancer

## cloudflare set-up

Set up load balancer with two pools named after your isp (or something else just be sure to add that name to the export below)
In each pool add your ip address

## Set up exports for your environment

Configure your unifi creds
Configure your cloudflare creds
Configure name of your pools

For example:

``` bash
export UNIFI_USER=joebob
export UNIFI_PWD=mysecurepwd
export UNIFI_URL=https://192.168.1.1/
export CLOUDFLARE_API_KEY='deadbeefdeadbeefdeadbeefdeadbeefdeadb'
export CLOUDFLARE_API_EMAIL='joebob@aol.com'
export WAN1_NAME=att
export WAN2_NAME=verizon
```

## Set up cron job or whatever and good to go