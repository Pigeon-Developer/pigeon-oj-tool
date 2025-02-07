#!/usr/bin/env bash


# copy from https://mirrors.tuna.tsinghua.edu.cn/help/docker-ce/
for pkg in docker.io docker-doc docker-compose podman-docker containerd runc; do apt-get remove $pkg; done

apt-get update
apt-get install ca-certificates curl gnupg

install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/debian/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://mirrors.tuna.tsinghua.edu.cn/docker-ce/linux/debian \
  "$(. /etc/os-release && echo "$VERSION_CODENAME")" stable" | \
  tee /etc/apt/sources.list.d/docker.list > /dev/null


apt-get update
apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin


echo '{
  "registry-mirrors": [
    "https://docker.1ms.run",
    "https://proxy.1panel.live",
    "https://docker.ketches.cn"
  ]
}' | tee /etc/docker/daemon.json

systemctl restart docker
