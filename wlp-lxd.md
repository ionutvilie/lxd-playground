# WAS Liberty Test


## LocalHost

sudo lxc launch ubuntu:16.04 wlp-test
sudo lxc exec wlp-test bash


## LXD Container

IBM Documentation:
 - https://www.ibm.com/support/knowledgecenter/en/was_beta_liberty/com.ibm.websphere.wlp.nd.multiplatform.doc/ae/twlp_ui_setup.html#twlp_ui_setup__uiinstall


```bash
apt update && apt install -y unzip
cd /opt/
# wget http://public.dhe.ibm.com/ibmdl/export/pub/software/websphere/wasdev/downloads/wlp/17.0.0.2/wlp-javaee7-17.0.0.2.zip # does not contain java
wget http://public.dhe.ibm.com/ibmdl/export/pub/software/websphere/wasdev/downloads/wlp/17.0.0.2/wlp-webProfile7-java8-linux-x86_64-17.0.0.2.zip
unzip wlp-webProfile7-java8-linux-x86_64-17.0.0.2.zip
cd wlp/

# install
bin/installUtility install javaee-7.0
bin/installUtility install adminCenter-1.0

# Start WAS Server
bin/server start
bin/server create myServer

# Get a sample war file from tomcat
cd /opt/wlp/usr/servers/myServer/dropins
wget https://tomcat.apache.org/tomcat-6.0-doc/appdev/sample/sample.war

# configure myServer
cd /opt/wlp/usr/servers/myServer/
# replace server.xml
```

```xml
<?xml version="1.0" encoding="UTF-8"?>
<server description="new server">

    <!-- Enable features -->
    <featureManager>
        <feature>adminCenter-1.0</feature>
    </featureManager>

   <!-- Define the host name for use by the collective.
        If the host name needs to be changed, the server should be
        removed from the collective and re-joined. -->
   <variable name="defaultHostName" value="localhost" />

    <basicRegistry id="basic">
       <user name="admin" password="adminpwd" />
       <user name="nonadmin" password="nonadminpwd" />
    </basicRegistry>

   <!-- Assign 'admin' to Administrator -->
    <administrator-role>
       <user>admin</user>
    </administrator-role>

   <keyStore id="defaultKeyStore" password="Liberty" />

   <httpEndpoint id="defaultHttpEndpoint"
                 host="*"
                 httpPort="80"
                 httpsPort="443" />

    <!-- Automatically expand WAR files and EAR files -->
    <applicationManager autoExpand="true"/>
</server>

```
```bash
cd /opt/wlp/
# start myServer
nohup bin/server run myServer  >/dev/null 2>&1 &

```

## LocalHost

Limit container resources to use 2 CPU's and 2GB of RAM

```bash
sudo lxc config set wlp-test limits.cpu 2
sudo lxc config set wlp-test limits.memory 2GB
```
