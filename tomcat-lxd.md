# Tomcat Test

Tomcat in a Linux container

## LocalHost

sudo lxc launch ubuntu:16.04 tomcat-test
sudo lxc exec ubuntu:16.04 bash


## LXD Container

```bash
# https://www.digitalocean.com/community/tutorials/how-to-install-apache-tomcat-8-on-ubuntu-16-04

# Step 1: Install Java
sudo apt-get update
sudo apt-get install default-jdk

# Step 2: Create Tomcat User
sudo groupadd tomcat
sudo useradd -s /bin/false -g tomcat -d /opt/tomcat tomcat

# Step 3: Install Tomcat
cd /tmp
curl -O http://mirrors.m247.ro/apache/tomcat/tomcat-8/v8.5.20/bin/apache-tomcat-8.5.20.tar.gz
sudo mkdir /opt/tomcat
sudo tar xzvf apache-tomcat-8*tar.gz -C /opt/tomcat --strip-components=1

# Step 4: Update Permissions
cd /opt/tomcat
sudo chgrp -R tomcat /opt/tomcat
sudo chmod -R g+r conf
sudo chmod g+x conf
sudo chown -R tomcat webapps/ work/ temp/ logs/

# Step 5: Create a systemd Service File

sudo update-java-alternatives -l
export JAVA_HOME=/usr/lib/jvm/java-1.8.0-openjdk-amd64/jre
sudo nano /etc/systemd/system/tomcat.service
```

sudo vi /etc/systemd/system/tomcat.service
```
[Unit]
Description=Apache Tomcat Web Application Container
After=network.target

[Service]
Type=forking

Environment=JAVA_HOME=/usr/lib/jvm/java-1.8.0-openjdk-amd64/jre
Environment=CATALINA_PID=/opt/tomcat/temp/tomcat.pid
Environment=CATALINA_HOME=/opt/tomcat
Environment=CATALINA_BASE=/opt/tomcat
Environment='CATALINA_OPTS=-Xms512M -Xmx1024M -server -XX:+UseParallelGC'
Environment='JAVA_OPTS=-Djava.awt.headless=true -Djava.security.egd=file:/dev/./urandom'

ExecStart=/opt/tomcat/bin/startup.sh
ExecStop=/opt/tomcat/bin/shutdown.sh

User=tomcat
Group=tomcat
UMask=0007
RestartSec=10
Restart=always

[Install]
WantedBy=multi-user.target
```


    sudo systemctl daemon-reload
    sudo systemctl status tomcat

    # # Step 6: Adjust the Firewall and Test the Tomcat Server
    # sudo ufw allow 8080
    # http://server_domain_or_IP:8080
    # sudo systemctl enable tomcat


    # Step 7: Configure Tomcat Web Management Interface

sudo vi /opt/tomcat/conf/tomcat-users.xml

```xml
<!-- add admin:admin user:pass | demo server only -->
<tomcat-users xmlns="http://tomcat.apache.org/xml"
      xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
      xsi:schemaLocation="http://tomcat.apache.org/xml tomcat-users.xsd"
      version="1.0">
<user username="admin" password="admin" roles="manager-gui,admin-gui"/>
</tomcat-users>
```

    sudo vi /opt/tomcat/webapps/manager/META-INF/context.xml
    sudo vi /opt/tomcat/webapps/host-manager/META-INF/context.xml

```xml
<!-- Comment Valve to allow remote access -->
<Context antiResourceLocking="false" privileged="true" >
    <!--
    <Valve className="org.apache.catalina.valves.RemoteAddrValve"
 allow="127\.\d+\.\d+\.\d+|::1|0:0:0:0:0:0:0:1" />
    -->
</Context>

```
    sudo systemctl restart tomcat


## Sample .war Application

Exit the container shell and execute on localhost

Limit container resources to use 2 CPU's and 2GB of RAM

    sudo lxc config set tomcat-test limits.cpu 2
    sudo lxc config set tomcat-test limits.memory 2GB


Download the sample and deploy it in tomcat using the manager:
 - https://tomcat.apache.org/tomcat-6.0-doc/appdev/sample/
 - wget https://tomcat.apache.org/tomcat-6.0-doc/appdev/sample/sample.war

Download the request tool
- https://github.com/ionutvilie/learning-golang/blob/master/http/concurrency/concurrency_worker_pools.go

Edit `*go` file and change the url to point to the sample app endpoint
- http://10.171.180.219:8080/sample/hello


1000 requests using 5 concurrent workers resulted below:

```
995 - Worker:  2 Job: 997 http://10.171.180.219:8080/sample/hello 200 0.000312
996 - Worker:  4 Job: 994 http://10.171.180.219:8080/sample/hello 200 0.000689
997 - Worker:  1 Job: 999 http://10.171.180.219:8080/sample/hello 200 0.000195
998 - Worker:  3 Job: 985 http://10.171.180.219:8080/sample/hello 200 0.001846
999 - Worker:  5 Job: 998 http://10.171.180.219:8080/sample/hello 200 0.000437
1000 - Worker:  2 Job: 1000 http://10.171.180.219:8080/sample/hello 200 0.000386
Program executed in 0.246564 seconds
```
