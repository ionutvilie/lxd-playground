# LXD Zentyal 5 Test

Zentyal test in LXD

 - http://www.zentyal.org/
 - https://help.ubuntu.com/lts/serverguide/zentyal.html
 - https://wiki.zentyal.org/wiki/Installation_Guide

## LocalHost

    sudo lxc launch ubuntu:16.04 zentyal-test
    sudo lxc exec zentyal-test bash

## Container

Execute inside the container 

    # https://wiki.zentyal.org/wiki/Installation_Guide
    # append Zentyal repository to the end of the sources list  and add apt key
    printf "\n#Zentyal\ndeb http://archive.zentyal.org/zentyal 5.0 main\n" >> /etc/apt/sources.list
    wget -q http://keys.zentyal.org/zentyal-5.0-archive.asc -O- | sudo apt-key add -


    sudo apt-get update
    sudo apt-get install zentyal
    # install additional modules
    ionut@zentyal-test:/$ apt install zentyal<hit TAB multiple times>
    zentyal             zentyal-ca          zentyal-dhcp        zentyal-groupware   zentyal-mailfilter  zentyal-objects     zentyal-samba       zentyal-sogo        
    zentyal-all         zentyal-common      zentyal-dns         zentyal-jabber      zentyal-network     zentyal-openchange  zentyal-services    zentyal-squid       
    zentyal-antivirus   zentyal-core        zentyal-firewall    zentyal-mail        zentyal-ntp         zentyal-openvpn     zentyal-software    



    groupadd ionut
    useradd --gid ionut --home /home/ionut ionut --shell /bin/bash
    passwd ionut
    sudo adduser ionut sudo


On localhost access your Zentyal test using last added user, in my case ionut, using the URL: `https://<lxd-container-ip>:8443/` # sudo lxc list | grep zentyal-test | awk '{print $6}'

- `sudo lxc stop zentyal-test` # stop container
- `sudo lxc delete zentyal-test` # destroy container
