#!/bin/bash

set -eux

private_key="$HOME/.ssh/id_rsa"
fingerprint=$(ssh-keygen -E md5 -lf "$private_key" | awk '{ print $2 }' | sed 's/MD5://')

(doctl compute ssh-key list | grep "$fingerprint") ||
  doctl compute ssh-key import knowhere --public-key-file "$private_key".pub

(doctl compute firewall list | grep knowhere) ||
  doctl compute firewall create knowhere \
    --inbound-rules "protocol:tcp,ports:22,address:0.0.0.0/0,address:::/0 protocol:tcp,ports:80,address:0.0.0.0/0,address:::/0" \
    --outbound-rules "protocol:icmp,address:0.0.0.0/0,address:::/0 protocol:tcp,ports:0,address:0.0.0.0/0,address:::/0 protocol:udp,ports:0,address:0.0.0.0/0,address:::/0" \
    --tag-names "knowhere"

doctl compute droplet create \
  knowhere-$(date +%s) \
  --region nyc3 \
  --image docker-20-04 \
  --size s-1vcpu-1gb \
  --ssh-keys "$fingerprint" \
  --enable-monitoring \
  --tag-name "knowhere" \
  --format "PublicIPv4" \
  --wait
