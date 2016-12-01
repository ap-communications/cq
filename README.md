# cq
[![release](https://img.shields.io/badge/release-1.0.0-blue.svg?style=flat-square)](https://github.com/ap-communications/cq/releases)
[![license: Apache](https://img.shields.io/badge/license-Apache-blue.svg?style=flat-square)](LICENSE)
* cq (cloud query) is a simple lightweight CLI tool for cloud environment control.


# License
Apache 2.0  
(read [./LICENSE](./LICENSE))


# Environment
* Windows
  * x64
  * 386
* Linux
  * x64
  * 386
  * arm
* darwin
  * x64
  * 386


# Start (Setup API key)

* Environment variables
```bash
[root@localhost ~]# export AWS_ACCESS_KEY_ID=xxxxxxxxxxxxxxxx
[root@localhost ~]# export AWS_SECRET_ACCESS_KEY=xxxxxxxxxxxxxxxx
```
or
* aws configure tool
```bash
[root@localhost ~]# aws configure
AWS Access Key ID [None]: xxxxxxxxxxxxxxxx
AWS Secret Access Key [None]: xxxxxxxxxxxxxxxx
Default region name [None]:       #There is no need to enter
Default output format [None]:       #There is no need to enter
```
or
* Attach IAM role


* Download a binary from [release](https://github.com/ap-communications/cq/releases) and execute :)


# Usage

## Command list
* rader
  * vm
    * list
      * --delimiter -d <string>
    * inspect
    * inspect [instance-ids]
    * start [instance-ids]
    * stop [instance-ids]
    * restart [instance-ids]
    * destroy [instance-ids]
    * easyup
      * --groupid
      * --imageid
      * --key
      * --region
      * --type
  * acl
    * list
      * --delimiter -d <string>
    * rule [securitygroup-ids]
    * add
      * --address [CIDR]
      * --groupid [securitygroup-id]
      * --protocol [tcp, udp, icmp, any]
      * --port [portnumber, any]
      * --way [ingress or egress]
    * delete
      * --address [CIDR]
      * --groupid [securitygroup-id]
      * --protocol [tcp, udp, icmp, any]
      * --port [portnumber, any]
      * --way [ingress or egress]
    * destroy [securitygroup-ids]


## VM list (default)
```
[root@localhost ~]# cq vm list
NAME-TAG                        INSTANCE-ID     STATE           GLOBAL          LOCAL           AZ                      PROVIDER
testserver01                    i-895ba307      stopped         NULL            172.20.1.68     ap-northeast-1c         AWS
testserver02                    i-d5bdf64a      stopped         NULL            172.20.0.220    ap-northeast-1a         AWS
testserver03                    i-e485796a      stopped         NULL            172.20.5.6      us-east-1d              AWS
```

## VM list  (delimiter)
```
[root@localhost ~]# cq vm list -d ,
NAME-TAG,INSTANCE-ID,STATE,GLOBAL,LOCAL,AZ,PROVIDER
testserver01,i-895ba307,stopped,NULL,172.20.1.68,ap-northeast-1c,AWS
testserver02,i-d5bdf64a,stopped,NULL,172.20.0.220,ap-northeast-1a,AWS
testserver03,i-e485796a,stopped,NULL,172.20.5.6,us-east-1d,AWS
```

## Start VM
```
[root@localhost ~]# cq vm start i-d5bdf64a i-e485796a
Success!  i-d5bdf64a   stopped  ===>  pending
Success!  i-e485796a   stopped  ===>  pending
```

## Stop VM
```
[root@localhost ~]# cq vm stop i-d5bdf64a i-e485796a
Success!  i-d5bdf64a   running  ===>  stopping
Success!  i-e485796a   running  ===>  stopping
```

## Reboot VM
```
[root@localhost ~]# cq vm reboot i-d5bdf64a i-e485796a
Success!  i-d5bdf64a   has started reboot sequence
Success!  i-e485796a   has started reboot sequence
```

## Destroy VM
```
[root@localhost ~]# cq vm destroy i-e485796a i-e485796a
Instance   i-e485796a i-e485796a   will be DESTROY, are you sure?  Y/N
Y
This is final warning. DESTROY instance   i-e485796a i-e485796a   ARE YOU SURE? (Check EBS data)  Y/N
Y
Success!   i-e485796a   destroyed
Success!   i-e485796a   destroyed
```

## Inspect VM
```
[root@localhost ~]# cq vm inspect i-d5bdf64a
{
  AmiLaunchIndex: 0,
  Architecture: "x86_64",
  BlockDeviceMappings: [{
      DeviceName: "/dev/xvda",
      Ebs: {
        AttachTime: 2016-10-31 06:36:55 +0000 UTC,
        DeleteOnTermination: true,
        Status: "attached",
        VolumeId: "vol-7787c1a9"
      }
    }],
  ClientToken: "HwPNQ1111111111111",
  EbsOptimized: enable,
  EnaSupport: true,
  Hypervisor: "xen",
  ImageId: "ami-cfr3ed3",
  InstanceId: "i-d5bdf64a",
  InstanceType: "p2.16xlarge",
--more--
```

## ACL list
```
[root@localhost ~]# cq acl list
GROUP-NAME              NAME-TAG        ID              DESCRIPTION                                             PROVIDER
default                 NULL            sg-7d44fc15     default VPC security group                              AWS
default                 NULL            sg-70937615     default VPC security group                              AWS
default                 NULL            sg-6b1a910e     default VPC security group                              AWS
default                 NULL            sg-6f52a50a     default VPC security group                              AWS
web-server              NULL            sg-76b2ac12     launch-wizard-2 created 2016-09-15T00:00:26.461+09:00   AWS
app-server              NULL            sg-414ff125     launch-wizard-2 created 2016-09-15T00:01:22.561+09:00   AWS
db-server               NULL            sg-cba3f4d8     launch-wizard-1 created 2016-07-16T14:43:33.313+09:00   AWS
```

## ACL rule
```
[root@localhost ~]# cq acl rule sg-76b2ac12
WAY             PROTOCOL        PORT                    ADDRESS
Ingress         tcp             80                      0.0.0.0/0
Ingress         tcp             443                     0.0.0.0/0
Ingress         tcp             22                      192.0.2.0/24
Egress          any             any                     0.0.0.0/0
```

## Add rule
```
[root@localhost ~]# cq acl add --groupid sg-76b2ac12 --way ingress --protocol tcp --port 8080 --address 192.0.2.0/24
Success!

[root@localhost ~]# cq acl rule sg-76b2ac12
WAY             PROTOCOL        PORT                    ADDRESS
Ingress         tcp             80                      0.0.0.0/0
Ingress         tcp             443                     0.0.0.0/0
Ingress         tcp             8080                    192.0.2.0/24
Ingress         tcp             22                      192.0.2.0/24
Egress          any             any                     0.0.0.0/0
```

## Delete rule
```
[root@localhost ~]# cq acl delete --groupid sg-76b2ac12 --way ingress --protocol tcp --port 8080 --address 192.0.2.0/24
Success!

[root@localhost ~]# cq acl rule sg-76b2ac12
WAY             PROTOCOL        PORT                    ADDRESS
Ingress         tcp             80                      0.0.0.0/0
Ingress         tcp             443                     0.0.0.0/0
Ingress         tcp             22                      192.0.2.0/24
Egress          any             any                     0.0.0.0/0
```

## Destroy ACL
```
[root@localhost ~]# cq acl list
GROUP-NAME              NAME-TAG        ID              DESCRIPTION                                             PROVIDER
default                 NULL            sg-7d44fc15     default VPC security group                              AWS
default                 NULL            sg-70937615     default VPC security group                              AWS
default                 NULL            sg-6b1a910e     default VPC security group                              AWS
default                 NULL            sg-6f52a50a     default VPC security group                              AWS
web-server              NULL            sg-76b2ac12     launch-wizard-2 created 2016-09-15T00:00:26.461+09:00   AWS
app-server              NULL            sg-414ff125     launch-wizard-2 created 2016-09-15T00:01:22.561+09:00   AWS
db-server               NULL            sg-cba3f4d8     launch-wizard-1 created 2016-07-16T14:43:33.313+09:00   AWS

[root@localhost ~]# cq acl destroy sg-76b2ac12
SecurityGroup   sg-76b2ac12   will be DESTROY, are you sure? (CAN NOT RESTORE) Y/N
Y
Success!

[root@localhost ~]# cq acl list
GROUP-NAME              NAME-TAG        ID              DESCRIPTION                                             PROVIDER
default                 NULL            sg-7d44fc15     default VPC security group                              AWS
default                 NULL            sg-70937615     default VPC security group                              AWS
default                 NULL            sg-6b1a910e     default VPC security group                              AWS
default                 NULL            sg-6f52a50a     default VPC security group                              AWS
app-server              NULL            sg-414ff125     launch-wizard-2 created 2016-09-15T00:01:22.561+09:00   AWS
db-server               NULL            sg-cba3f4d8     launch-wizard-1 created 2016-07-16T14:43:33.313+09:00   AWS
```

## Easy create VM
```
[root@localhost ~]# cq vm easyup
............
     Instance ID: i-4a1e6269
SecurityGroup ID: sg-a2ed64cc

  *** IMPORTANT: SSH (TCP22) is anyone can access!! ***

          Global: 52.198.207.112
         SSH Key:
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAyMBhQ0+QoXTSk6zrZmtslce9jgu0KJ7VVW6K3kopwkUoZmJnuK5JzVX7n7JL
dCf4stOSQkFPwFIHRoV4z0y0/rgiyKjZ6LF9fRAe3eNAO6CSKDlF/oyuytbJ5olAIMuxsSjSmX9W
1/7fwtbiYwqxhw1iXsx7ZiHyI22P9s/+eVJQYMCyqXeHvbhezXHiBtVQG831n9JgIrd4W2yD14e0
wvqw17omun5nkCb6ufI15CFQ4B+u4e1C6VGqMpXiBGNJs4wKBk//psaFIze+9OkzPn6Svx2iKeH0
me01Q4eC1uaAGDQ1XPVwie2rgqvOSSc83SDDE+MuPFNc/EgsQ6CtUwIDAQABAoIBAHwk/AD0Iyy1
YcORG8GqjOvTKZW+BxtXnfXG9nmgw1IwElu+XwYGQf2JPqHzUFX0ogd0bu4qFXeJQpaJ07veY89h
b6FHSfpsSH4eifgNoJs/ISNex7oypaUqTpESL2YYkTpNVG84ICxSoW2MFVPKOR6bWEnniigOtf7S
skSfO28qL2klExsjRodE8K7ucE7klYU7VhEYw+IjKFE8/PoUd86Gtfcw9JgVFUV3SVmxoRNbQnsk
R1LFukyIZptu17Act87upnNV92KSH6NrUTTlYS/zawyihzwiEOoPweA/lxbhVzjTiniwuv00Do0O
+JSmwh/mDH9L3trOmNFLakYdQAECgYEA7qoXkL4MopagFIw5S++qp5RIJ0REIq91cHqUMddo5mW6
vxCkvEK+potX3vJxjwC8WaF7VK2U9Vp3d/5ZA5TL1iWsf3VTM8LVLUwGegUhO+O3R2iv028NGsJq
wk1dSR2qGw3lBDLv9RODVMAqUcIQqypGWW85ObIHcXrtovlz+gECgYEA11VOTEPPWBlg5/6F6BJi
uNcOYDfjAL3ms8OAmzqG/kAm8+nOfoXWTvMk68RM1o9UqEWSpnb3TFS47PMd7EndvkT2wJ/EB96g
k9a4D+As5IXNsWnzMKvQn/ScwHqnMpur4eS+8zwKnQeNVaQZwdYGWxS8A7wozf3unPpNhtrAn1MC
gYB7M5TcpC7Tk5vpX4WIXJ0kytgAZS4jFZ4zRSPRItjE7sjbLtVLVc0jHhVrQo46eu+/+Ss4SC74
BQ5dBBChV6Nt1Z7ZYRBlTPNM4c0KGjHQv7lyBGmXZOFH6grAaL16K/UO6DpL05upULxf+J8f+2Ut
ZrJQIqlQzkUnYITht5IWAQKBgBf98zlgHMwOF9QBcd4HnmRF4d4qm+pLPlctIzkobj4J2801ccJ7
GwO055o1RvJpCd+t0McnXiptDiWoeN8I4+H9QRj5NBhEX9PZCq0KJzZXCjIIQgZcPmR0FrsDb3me
Cqqil/v15qDUZJT0McJ9HDwttT0dS9aXwmxPHzA1U1ztAoGBAMuDpx3BmUHeVJy53JT4yqj6yvrA
3alqPgo3mTpz+sB5sDp46MJy/kvJVD90HM1BUDfa7+jEn3LsgxGFH+Naf8lr/IP/Guzf2gAg38EX
xgskK8ZPmia3AVbmocNbX88lDavfpbMFGahoX07L7g1Sx59IWcRLSq94sWjpC8G1Yxom
-----END RSA PRIVATE KEY-----


[root@localhost ~]# cq vm easyup --groupid sg-a2ed64cb --imageid ami-0567c164 --key test-key --region ap-northeast-1 --type t2.nano
..
     Instance ID: i-4a1e6267
SecurityGroup ID: sg-a2ed64cb
          Global: 52.199.118.117
         SSH Key: test-key

```


# VS.
* AWS CLI
  * There is no need to be aware of JSON
  * There is no need to learn many subcommands and options
* Management Console
  * Lightweight
  * Simple UI


# Thanks
* [Cobra CLI framework](https://github.com/spf13/cobra)


# Author
Ryo Yamaoka
unite@ap-com.co.jp
