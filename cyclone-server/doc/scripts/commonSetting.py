#!/usr/bin/env python
# -*- coding: utf-8 -*-
# =========================================================================
# FileName: commonSetting.py
# Creator: JosingCai
# Mail: caijiaoxing@idcos.com
# Created Time: 2018-08-28
# Description: 配置操作系统中需更新的项
# Usage: 1. Master存放路径：/opt/cloudboot/home/www/scripts
#           BootOS存放路径：/tmp
#        2. 手动调试：
#           （1）下载脚本，curl -o  /tmp/commonSetting.py http://osinstall/scripts/commonSetting.py
#           （2）执行脚本：python /tmp/commonSetting.py --network=Y --osuser=Y ...
#           （3）操作日志路径：/tmp/operator.log
#            参数在装机参数模板中，如:http://osinstall.idcos.com/api/cloudboot/v1/devices/{sn}/settings/os-users
# History: 2018-08-28/Create file
#          2019-03-22/Adapter config network in tlinux distribution
#          2019-03-25/Add function removeOldcfg() to remove old net config file
#          2019-03-25/Add function getDomainIP() to get the IP of Domain Name
#          2019-05-30/Optimized: Adapt to earlier version of dmidecode program
#          2019-07-24/Fit to SuSE python3.x Env; Fit to SueSE network config
#          2019-08-06/Add function to get MAC for eth port in SuSE
#          2019-08-09/Add funtion to restart network after config network
#          2019-10-15/Fit bondTxt： BONDING_MODULES_OPTS  to BONDING_MODULE_OPTS
#          2020-01-14/Add Eth4 and eth5 will be added to bond1 when there are 6 ethernet controller
#          2020-02-13/Add modifySSHDConfig set innner ip for ssh listen 
#          2020-03-16/Fix modifySSHDConfig set Port 36000 and replace listenaddress anyway
#          2020-04-14/Fix bondTxt，ethTxt: MN_CONTROLLED to NM_CONTROLLED=no
#          2020-05-20/Add mountData: sdb - sdo for /data1 - /dat14 if exist
#          2022-09-26/Add update network for openEuler
# Copyright (c) 2015-2028 iDCOS Tech. All Right Reserved.
# =========================================================================

import json
import os
import datetime
import sys
import socket
import time

if sys.version > '3':
    import urllib.request as urllib2
    import subprocess as commands
else:
    import urllib2
    import commands

reload(sys) 
sys.setdefaultencoding('utf8') 

DEBUG = False
POST_INTERVAL = 3  # 配置负载均衡后, POST两次请求过快，会导致主动拒绝客户端的连接，出现499

host = "osinstall.idcos.com"
pathHead = "http://%s/api/cloudboot/v1" %host

def getSN():
    cmd = "dmidecode -s system-serial-number | awk '/^[^#]/ { print $1 }'"
    status, output = commands.getstatusoutput(cmd)
    if status != 0:
        return False, output
    sn = output.strip()
    return True, sn

class SystemInfo(object):
    """docstring for SystemInfo"""
    def __init__(self, sn):
        self.sn = sn.replace(" ", '')

    def getDomainIP(self):
        domainName = host
        ret = socket.getaddrinfo(domainName, None)
        IP = ret[0][4][0]
        self.logRecord("%s: %s" %(host, IP))

    def uploadProgress(self, item, status=2):
        statusDict = {0: "Pass", 1: "Fail", 2: "Skip"}
        items = {"network":{"title":"Config Network", "progress":0.8, "log":"Q29uZmlnIE5ldHdvcms="},
                 "hostName":{"title":"Config HostName", "progress":0.82, "log":"Q29uZmlnIEhvc3ROYW1l"},
                 "OSUser":{"title":"Config OS User", "progress":0.85, "log":"Q29uZmlnIFVzZXI="},
                 "firmware":{"title":"Update Firmware", "progress":0.9, "log":"VXBkYXRlIEZpcm13YXJl"},
                 "service":{"title":"Config Service", "progress":0.95, "log":"Q29uZmlnIFN5c3RlbSBTZXJ2aWNl"},
                 "mountdata":{"title":"Mount Data", "progress":0.96, "log":"TW91bnQgRGF0YQ=="},
                 "complete":{"title":"Complete", "progress":1, "log":"SW5zdGFsbCBGaW5pc2hlZA=="}
                 }
        if item == "complete":
            cmd = r'curl -H "Content-Type: application/json" -X POST -d "{\"title\":\"%s\",\"progress\":%s,\"log\":\"%s\"}" %s/devices/%s/installations/progress'%(items[item]['title'], items[item]['progress'], items[item]['log'], pathHead, self.sn)
        else:
            cmd = r'curl -H "Content-Type: application/json" -X POST -d "{\"title\":\"%s %s\",\"progress\":%s,\"log\":\"%s\"}" %s/devices/%s/installations/progress'%(items[item]['title'], statusDict[status], items[item]['progress'], items[item]['log'], pathHead, self.sn)
        status, output = commands.getstatusoutput(cmd)
        self.logRecord(cmd)
        self.logRecord(output)
        if "error" in output:
            time.sleep(POST_INTERVAL)
            status, output = commands.getstatusoutput(cmd)
            self.logRecord(output)
        if status:
            return False
        time.sleep(POST_INTERVAL)
        return True
    
    def logRecord(self, string, Type=0):
        info = {0: "Info", 1: "Error"}
        curTime=datetime.datetime.now().strftime('%Y-%m-%d_%H:%M:%S')
        String = "[ %s ] %s: %s\n"%(curTime, info[Type], string)
        if DEBUG is True:
            print String
        with open("/tmp/operator.log", 'a+') as f:
            f.write(String)   

    def updateOSUser(self):
        url = '%s/devices/%s/settings/os-users' %(pathHead,self.sn)
        self.logRecord(url, 0)
        try:
            req = urllib2.Request(url)
            result = urllib2.urlopen(req)
        except:
            self.logRecord("Connect %s Error, Please check it ~ " %host, 1)
            self.uploadProgress("OSUser", 1)
            return False

        if sys.version > '3':
            ret = json.loads(str(result.read(), 'utf-8'))
        else:
            ret = json.loads(result.read())
        self.logRecord(ret["content"])

        item = ret['content']
        user_name = item['username']
        passwd = item['password']
        cmd = "cat /etc/passwd | awk -F ':' '{print $1}' | grep %s"%user_name
        status, output = commands.getstatusoutput(cmd)
        if status != 0 and len(output) == 0:
            cmd = 'useradd -d /home/%s -m %s -g root' % (user_name, user_name)
            status, output = commands.getstatusoutput(cmd)
            if status != 0:
                self.logRecord(output, 1)
                self.uploadProgress("OSUser", 1)
                return False
            self.logRecord(output)
        cmd = 'echo %s | passwd %s --stdin' % (passwd, user_name)
        status, output = commands.getstatusoutput(cmd)
        if status != 0:
            self.logRecord(output, 1)
            self.uploadProgress("OSUser", 1)
            return False
        self.logRecord(output)
        self.uploadProgress("OSUser", 0)
        return True

    def removeOldCfg(self, os_type="centos"):
        ethList = ["eth0", "eth1", "eth2", "eth3", "eth4", "eth5", "bond0", "bond1"]
        if os_type == "sles":
            base = "/etc/sysconfig/network/ifcfg"
        else:
            base = "/etc/sysconfig/network-scripts/ifcfg"
        for eth in ethList:
            netPath = "%s-%s" %(base, eth)
            if os.path.exists(netPath):
                os.remove(netPath)

    def remove70File(self):
        cmd = "ls /etc/udev/rules.d/70*"
        status, output = commands.getstatusoutput(cmd)
        self.logRecord(output)
        fileList = output.split("\n")
        for filePath in fileList:
            if os.path.exists(filePath):
                os.remove(filePath)
                return True
            else:
                return True

    def getNetMAC(self):
        portList = []
        cmd = "ip -4 -o l"
        status, output = commands.getstatusoutput(cmd)
        self.logRecord(output)
        if status != 0 :
            return False
        tmps = output.split("\n")
        portNum = len(tmps)
        mac_list = []
        for i in range(len(tmps)):
            line = tmps[i]
            parts = line.split(":")
            name = parts[1].strip()
            if name != "lo":
                parts = line.split("link/ether")
                mac = parts[1].split()[0].strip()
                mac_list.append(mac)
        self.logRecord(mac_list)
        return mac_list

    def modifySSHDConfig(self, intranet_ip):
        sshd_config_file = '/etc/ssh/sshd_config'
        sshd_listten_ip = intranet_ip

        cmd_check = 'grep "^Port" '+sshd_config_file+' | awk \'{print $2}\''
        status, output = commands.getstatusoutput(cmd_check)
        if output == "":
            self.logRecord("sshd listen port is not set.")
            cmd_modify = "sed -i '/^#Port 22/a\Port 36000' "+sshd_config_file
            self.logRecord(cmd_modify)
            commands.getstatusoutput(cmd_modify)

        if output != "" and output != "36000":
            self.logRecord("sshd listen port is not 36000.")
            cmd_modify = "sed -i 's/^Port.*/Port 36000/' "+sshd_config_file
            self.logRecord(cmd_modify)
            commands.getstatusoutput(cmd_modify)

        cmd_check = 'grep "^ListenAddress" '+sshd_config_file
        status, output = commands.getstatusoutput(cmd_check)
        if output != "":
            self.logRecord("ListenAddress already exist. Replace it.")
            cmd_modify = "sed -i 's/^ListenAddress.*/ListenAddress "+sshd_listten_ip+"/' "+sshd_config_file
            self.logRecord(cmd_modify)
            status, output = commands.getstatusoutput(cmd_modify)
        else:
            cmd_modify = "sed -i '/^Port 36000/a\ListenAddress "+sshd_listten_ip+"' "+sshd_config_file
            self.logRecord(cmd_modify)
            status, output = commands.getstatusoutput(cmd_modify)

    def setNetInfo(self, os_type, **info):
        if os_type == "sles":
            netPath = "/etc/sysconfig/network"
            if info["version"] == "ipv4":
                gatewayStr = "default %s\n"%info["gateway"]   #SuSE网关配置
                with open("%s/routes" %netPath, "w") as f:
                    f.write(gatewayStr)
        else:
            netPath = "/etc/sysconfig/network-scripts"

        cmd = "lspci | grep -i Ethernet | wc -l"
        status, output = commands.getstatusoutput(cmd)
        self.logRecord("There are %s ethernet controller." % output)
        countEth = 0
        if output == "6":
            countEth = 6 # 针对主板4个网口未禁用，使用外接两网口网卡场景
            self.logRecord("Eth4 and eth5 will be added to bond1 when there are 6 ethernet controller.")

        if info["bonding"] == "yes":
            if info["version"] == "ipv4":
                # 指定内外网对应的网卡设备
                if info["scope"] == "intranet":
                    if countEth == 6:
                        eths = ["eth1", "eth0", "eth4", "eth5"]
                    else:
                        eths = ["eth1", "eth0"]
                else:
                    # 简单地规避效的bond配置
                    if countEth >=4:
                        eths = ["eth2", "eth3"]
                    else:
                        eths = []
                # sles 网络配置兼容
                if os_type == "sles":
                    for eth in eths:
                        ethTxt = "# IP Config for %s:\nDEVICE='%s'\nBOOTPROTO='static'\nSTARTMODE='onboot'\n" %(eth, eth)
                        with open("%s/ifcfg-%s"%(netPath, eth), 'w') as f:
                            f.write(ethTxt)
                    bondTxt = "# IP Config for %s:\nDEVICE=%s\nBOOTPROTO=static\nSTARTMODE=onboot\nIPADDR=%s\nNETMASK=%s\nGATEWAY=%s\nBONDING_MASTER=yes\nBONDING_MODULE_OPTS='miimon=100 mode=4'\nBONDING_SLAVE0=%s\nBONDING_SLAVE1=%s\n"%(info["name"], info["name"], info['ip'], info['netmask'], info['gateway'],eths[0], eths[1])
                    with open("%s/ifcfg-%s"%(netPath, info['name']), 'w') as f:
                        f.write(bondTxt)
                # openEuler 网络配置兼容(该发行版使用NetworkManager服务管理网络)
                elif os_type == "openEuler":
                    for eth in eths:
                        ethTxt = "# IP Config for %s:\nDEVICE=%s\nTYPE=Ethernet\nBOOTPROTO=none\nONBOOT=yes\nSLAVE=yes\nMASTER=%s\n"%(eth, eth, info['name'])
                        with open("%s/ifcfg-%s"%(netPath, eth), 'w') as f:
                            f.write(ethTxt)
                    bondTxt = "# IP Config for %s:\nDEVICE=%s\nTYPE=Bond\nONBOOT=yes\nBOOTPROTO=none\nDELAY=0\nIPADDR=%s\nNETMASK=%s\nGATEWAY=%s\nBONDING_OPTS='mode=4 miimon=100 lacp_rate=fast xmit_hash_policy=layer3+4'\n"%(info["name"], info["name"], info['ip'], info['netmask'], info['gateway'])
                    with open("%s/ifcfg-%s"%(netPath, info['name']), 'a') as f:
                        f.write(bondTxt)
                    # 存在外网时，固化私网网段路由表，修改默认路由指向外网
                    if info["scope"] == "extranet":
                        routeTxt = "9.0.0.0/8 via %s dev %s\n10.0.0.0/8 via %s dev %s\n100.64.0.0/10 via %s dev %s\n"%(info['gateway'], info['name'], info['gateway'], info['name'], info['gateway'], info['name'])
                        with open("%s/route-%s"%(netPath, "bond1"), 'w') as f:
                            f.write(routeTxt)                
                        with open("%s/ifcfg-%s"%(netPath, "bond0"), 'a') as f:
                            f.write("DEFROUTE=yes")
                        with open("%s/ifcfg-%s"%(netPath, "bond1"), 'a') as f:
                            f.write("DEFROUTE=no")
                else:
                    for eth in eths:
                        ethTxt = "# IP Config for %s:\nDEVICE=%s\nTYPE=Ethernet\nBOOTPROTO=none\nONBOOT=yes\nNM_CONTROLLED=no\nSLAVE=yes\nMASTER=%s\n"%(eth, eth, info['name'])
                        with open("%s/ifcfg-%s"%(netPath, eth), 'w') as f:
                            f.write(ethTxt)
                    bondTxt = "# IP Config for %s:\nDEVICE=%s\nONBOOT=yes\nBOOTPROTO=static\nNM_CONTROLLED=no\nDELAY=0\nIPADDR=%s\nNETMASK=%s\nGATEWAY=%s\nBONDING_OPTS='mode=4 miimon=100 lacp_rate=fast xmit_hash_policy=layer3+4'\n"%(info["name"], info["name"], info['ip'], info['netmask'], info['gateway'])
                    with open("%s/ifcfg-%s"%(netPath, info['name']), 'w') as f:
                        f.write(bondTxt)

            if info["version"] == "ipv6":
                bondTxt = "IPV6INIT=yes\nIPV6_AUTOCONF=no\nIPV6_FAILURE_FATAL=no\nIPV6ADDR=%s\nIPV6_DEFAULTGW=%s\n"%((info['ip']+"/"+info['netmask']), info['gateway'])
                with open("%s/ifcfg-%s"%(netPath, info['name']), 'a') as f:
                    f.write(bondTxt)          
        else:
            if os_type == "sles":
                macs = self.getNetMAC()
                index = int(info['name'][-1])
                if index <= len(macs):
                    if info["version"] == "ipv4":
                        ethText = "# IP Config for %s:\nBOOTPROTO=static\nIPADDR=%s\nNETMASK=%s\nSTARTMODE=auto\nHWADDR=%s\n" %(info['name'], info['ip'], info['netmask'], macs[index])
                        with open("%s/ifcfg-%s"%(netPath, info['name']), 'w') as f:
                            f.write(ethText)
                else:
                    return False
            else:
                if info["version"] == "ipv4":
                    ethText = "# IP Config for %s:\nDEVICE=%s\nTYPE=Ethernet\nBOOTPROTO=static\nIPADDR=%s\nNETMASK=%s\nONBOOT=yes\nGATEWAY=%s\n"% (info['name'], info['name'], info['ip'], info['netmask'], info['gateway'])
                    with open("%s/ifcfg-%s"%(netPath, info['name']), 'w') as f:
                        f.write(ethText)
        return True

    def updateNetwork(self):
        url = '%s/devices/%s/settings/networks' %(pathHead, self.sn)
        self.logRecord(url, 0)
        try:
            req = urllib2.Request(url)
            result = urllib2.urlopen(req)
        except:
            self.logRecord("Connect %s Error, Please check it ~ " % host, 1)
            self.uploadProgress("network", 1)
            return False
        if sys.version > '3':
            ret = json.loads(str(result.read(), 'utf-8'))
        else:
            ret = json.loads(result.read())
        self.logRecord(ret["content"])

        if ret["content"]["ip_source"] == "dhcp":
            self.uploadProgress("network", 0)
            self.logRecord("DHCP Mode", 0)
            return True

        # 网络配置支持的OS-RELEASE
        supportList = ["centos", "redhat", "sles", "tlinux", "openEuler"]
        if os.path.exists("/etc/os-release"):
            cmdSys = "cat /etc/os-release | grep -i '^ID=' | awk -F '=' '{print $2}'"
            cmdVer = "cat /etc/os-release | grep -i '^VERSION=' | awk -F '=' '{print $2}'"
            status, output = commands.getstatusoutput(cmdSys)
            if status != 0:
                self.logRecord("Not Support Set Network in Current OS: %s" %output, 1)
                self.uploadProgress("network", 1)
                return False
            ID = output.replace("\"", "")
            status, output = commands.getstatusoutput(cmdVer)
            if status != 0:
                self.logRecord("Not Support Set Network in Current OS: %s" %output, 1)
                self.uploadProgress("network", 1)
                return False
            VERSION = output.replace("\"", "")
        elif os.path.exists("/etc/redhat-release"):
            ID = "centos"
            VERSION = "6"
        else:
            self.logRecord("Not Support Set Network in Current OS: %s" %output, 1)
            self.uploadProgress("network", 1)
            return False

        if ID not in supportList:
            self.logRecord("Not Support Set Network in Current OS: %s" %ID, 1)
            self.uploadProgress("network", 1)
            return False

        if ID == "sles":
            netPath = "/etc/sysconfig/network"
        else:
            netPath = "/etc/sysconfig/network-scripts"
        self.removeOldCfg(ID)

        netDict = {}
        netDict["bonding"] = ret['content']['bonding_required']
        netDict["ip_source"] = ret['content']["ip_source"]
        items = ret['content']['items']
        for info in items:
            netDict["scope"] = info["scope"]
            netDict["version"] = info["version"]
            if netDict["bonding"] == "no":
                if info["scope"] == "intranet":
                    netDict["name"] = "eth1"
                elif info["scope"] == "extranet":
                    netDict["name"] = "eth0"
                else:
                    self.logRecord("Not Support Current NetWork Mode [scope]: [%s]…… "%str(info["scope"]), 1)
                    self.uploadProgress("network", 1)
                    return False
            elif netDict["bonding"] == "yes":
                if info["scope"] == "intranet":
                    netDict["name"] = "bond1"
                elif info["scope"] == "extranet":
                    netDict["name"] = "bond0"
                else:
                    self.logRecord("Not Support Current NetWork Mode [scope]: [%s]…… "%str(info["scope"]), 1)
                    self.uploadProgress("network", 1)
                    return False
            else:
                self.logRecord("Not Support Current NetWork Mode [bonding]:[%s] …… "%str(netDict["bonding"]), 1)
                self.uploadProgress("network", 1)
                return False
            netDict["ip"] = info['ip']
            netDict["netmask"] = info['netmask']
            # there maybe more than one gw ip split by comma
            gateway_list = info['gateway'].split(",")
            if len(gateway_list) > 1:
                netDict["gateway"] = gateway_list[0]
            else:
                netDict["gateway"] = info['gateway']
            
            self.logRecord(str(netDict))
            # set intranet ipv4 for sshd listen
            if netDict["scope"] == "intranet" and netDict["version"] == "ipv4":
                self.modifySSHDConfig(netDict["ip"])
            # network scripts config
            if not self.setNetInfo(ID, **netDict):
                self.logRecord("Setting Network Failed, Please Check it ~", 1)
                self.uploadProgress("network", 1)
                return False

        self.uploadProgress("network", 0)
        self.remove70File()
        return True

    def updateService(self):
        ret = self.uploadProgress("service", 2)
        return ret

    def updateComplete(self):
        ret = self.uploadProgress("complete", 0)
        status, output = commands.getstatusoutput("systemctl restart network")
        if status != 0:
            self.logRecord(output,1)
            status, output = commands.getstatusoutput("service network restart")
            if status != 0:
                self.logRecord(output, 1)
            self.logRecord(output)
        self.logRecord(output)
        return ret

    def mountData(self):
        fail_count = 0
        # 查询设备的硬件配置参数(序列号、厂商、型号、设备类型、硬件备注)
        url = '%s/devices/%s/settings/hardwareinfo' %(pathHead, self.sn)
        self.logRecord(url, 0)
        try:
            req = urllib2.Request(url)
            result = urllib2.urlopen(req)
        except:
            self.logRecord("Connect %s Error, Please check it ~ " % host, 1)
            self.uploadProgress("mountdata", 2)
            return False
        if sys.version > '3':
            ret = json.loads(str(result.read(), 'utf-8'))
        else:
            ret = json.loads(result.read())
        self.logRecord(ret["content"])

        # 预定义块设备、挂载点、文件系统：sdb - sdak : /data1 - /data36
        mountPointList = [
            {
                "blockDevice":"/dev/sdb",
                "mountPoint":"/data1",
                "fsType":"ext4",
                "fsTab":"/dev/sdb        /data1        ext4        noatime,acl,user_xattr,nofail    0  0"
            },
            {
                "blockDevice":"/dev/sdc",
                "mountPoint":"/data2",
                "fsType":"ext4",
                "fsTab":"/dev/sdc        /data2        ext4        noatime,acl,user_xattr,nofail    0  0"
            },            
            {
                "blockDevice":"/dev/sdd",
                "mountPoint":"/data3",
                "fsType":"ext4",
                "fsTab":"/dev/sdd        /data3        ext4        noatime,acl,user_xattr,nofail    0  0"
            },
            {
                "blockDevice":"/dev/sde",
                "mountPoint":"/data4",
                "fsType":"ext4",
                "fsTab":"/dev/sde        /data4        ext4        noatime,acl,user_xattr,nofail    0  0"
            },
            {
                "blockDevice":"/dev/sdf",
                "mountPoint":"/data5",
                "fsType":"ext4",
                "fsTab":"/dev/sdf        /data5        ext4        noatime,acl,user_xattr,nofail    0  0"
            },
            {
                "blockDevice":"/dev/sdg",
                "mountPoint":"/data6",
                "fsType":"ext4",
                "fsTab":"/dev/sdg        /data6        ext4        noatime,acl,user_xattr,nofail    0  0"
            },
            {
                "blockDevice":"/dev/sdh",
                "mountPoint":"/data7",
                "fsType":"ext4",
                "fsTab":"/dev/sdh        /data7        ext4        noatime,acl,user_xattr,nofail    0  0"
            },
            {
                "blockDevice":"/dev/sdi",
                "mountPoint":"/data8",
                "fsType":"ext4",
                "fsTab":"/dev/sdi        /data8        ext4        noatime,acl,user_xattr,nofail    0  0"
            },                        
            {
                "blockDevice":"/dev/sdj",
                "mountPoint":"/data9",
                "fsType":"ext4",
                "fsTab":"/dev/sdj        /data9        ext4        noatime,acl,user_xattr,nofail    0  0"
            },
            {
                "blockDevice":"/dev/sdk",
                "mountPoint":"/data10",
                "fsType":"ext4",
                "fsTab":"/dev/sdk        /data10       ext4        noatime,acl,user_xattr,nofail    0  0"
            },                        
            {
                "blockDevice":"/dev/sdl",
                "mountPoint":"/data11",
                "fsType":"ext4",
                "fsTab":"/dev/sdl        /data11       ext4        noatime,acl,user_xattr,nofail    0  0"
            },            
            {
                "blockDevice":"/dev/sdm",
                "mountPoint":"/data12",
                "fsType":"ext4",
                "fsTab":"/dev/sdm        /data12       ext4        noatime,acl,user_xattr,nofail    0  0"
            },
            {
                "blockDevice":"/dev/sdn",
                "mountPoint":"/data13",
                "fsType":"ext4",
                "fsTab":"/dev/sdn        /data13       ext4        noatime,acl,user_xattr,nofail    0  0"
            },
            {
                "blockDevice":"/dev/sdo",
                "mountPoint":"/data14",
                "fsType":"ext4",
                "fsTab":"/dev/sdo        /data14       ext4        noatime,acl,user_xattr,nofail    0  0"
            },
            {
                "blockDevice":"/dev/sdp",
                "mountPoint":"/data15",
                "fsType":"ext4",
                "fsTab":"/dev/sdp        /data15       ext4        noatime,acl,user_xattr,nofail    0  0"
            },
            {
                "blockDevice":"/dev/sdq",
                "mountPoint":"/data16",
                "fsType":"ext4",
                "fsTab":"/dev/sdq        /data16       ext4        noatime,acl,user_xattr,nofail    0  0"
            },
            {
                "blockDevice":"/dev/sdr",
                "mountPoint":"/data17",
                "fsType":"ext4",
                "fsTab":"/dev/sdr        /data17       ext4        noatime,acl,user_xattr,nofail    0  0"
            },
            {
                "blockDevice":"/dev/sds",
                "mountPoint":"/data18",
                "fsType":"ext4",
                "fsTab":"/dev/sds        /data18       ext4        noatime,acl,user_xattr,nofail    0  0"
            },
            {
                "blockDevice":"/dev/sdt",
                "mountPoint":"/data19",
                "fsType":"ext4",
                "fsTab":"/dev/sdt        /data19       ext4        noatime,acl,user_xattr,nofail    0  0"
            },
            {
                "blockDevice":"/dev/sdu",
                "mountPoint":"/data20",
                "fsType":"ext4",
                "fsTab":"/dev/sdu        /data20       ext4        noatime,acl,user_xattr,nofail    0  0"
            },
            {
                "blockDevice":"/dev/sdv",
                "mountPoint":"/data21",
                "fsType":"ext4",
                "fsTab":"/dev/sdv        /data21       ext4        noatime,acl,user_xattr,nofail    0  0"
            },
            {
                "blockDevice":"/dev/sdw",
                "mountPoint":"/data22",
                "fsType":"ext4",
                "fsTab":"/dev/sdw        /data22       ext4        noatime,acl,user_xattr,nofail    0  0"
            },
            {
                "blockDevice":"/dev/sdx",
                "mountPoint":"/data23",
                "fsType":"ext4",
                "fsTab":"/dev/sdx        /data23       ext4        noatime,acl,user_xattr,nofail    0  0"
            },
            {
                "blockDevice":"/dev/sdy",
                "mountPoint":"/data24",
                "fsType":"ext4",
                "fsTab":"/dev/sdy        /data24       ext4        noatime,acl,user_xattr,nofail    0  0"
            },
            {
                "blockDevice":"/dev/sdz",
                "mountPoint":"/data25",
                "fsType":"ext4",
                "fsTab":"/dev/sdz        /data25       ext4        noatime,acl,user_xattr,nofail    0  0"
            },
            {
                "blockDevice":"/dev/sdaa",
                "mountPoint":"/data26",
                "fsType":"ext4",
                "fsTab":"/dev/sdaa        /data26       ext4        noatime,acl,user_xattr,nofail    0  0"
            },
            {
                "blockDevice":"/dev/sdab",
                "mountPoint":"/data27",
                "fsType":"ext4",
                "fsTab":"/dev/sdab        /data27       ext4        noatime,acl,user_xattr,nofail    0  0"
            },
            {
                "blockDevice":"/dev/sdac",
                "mountPoint":"/data28",
                "fsType":"ext4",
                "fsTab":"/dev/sdac        /data28       ext4        noatime,acl,user_xattr,nofail    0  0"
            },
            {
                "blockDevice":"/dev/sdad",
                "mountPoint":"/data29",
                "fsType":"ext4",
                "fsTab":"/dev/sdad        /data29       ext4        noatime,acl,user_xattr,nofail    0  0"
            },
            {
                "blockDevice":"/dev/sdae",
                "mountPoint":"/data30",
                "fsType":"ext4",
                "fsTab":"/dev/sdae        /data30       ext4        noatime,acl,user_xattr,nofail    0  0"
            },
            {
                "blockDevice":"/dev/sdaf",
                "mountPoint":"/data31",
                "fsType":"ext4",
                "fsTab":"/dev/sdaf        /data31       ext4        noatime,acl,user_xattr,nofail    0  0"
            },
            {
                "blockDevice":"/dev/sdag",
                "mountPoint":"/data32",
                "fsType":"ext4",
                "fsTab":"/dev/sdag        /data32       ext4        noatime,acl,user_xattr,nofail    0  0"
            },
            {
                "blockDevice":"/dev/sdah",
                "mountPoint":"/data33",
                "fsType":"ext4",
                "fsTab":"/dev/sdah        /data33       ext4        noatime,acl,user_xattr,nofail    0  0"
            },
            {
                "blockDevice":"/dev/sdai",
                "mountPoint":"/data34",
                "fsType":"ext4",
                "fsTab":"/dev/sdai        /data34       ext4        noatime,acl,user_xattr,nofail    0  0"
            },
            {
                "blockDevice":"/dev/sdaj",
                "mountPoint":"/data35",
                "fsType":"ext4",
                "fsTab":"/dev/sdaj        /data35       ext4        noatime,acl,user_xattr,nofail    0  0"
            },
            {
                "blockDevice":"/dev/sdak",
                "mountPoint":"/data36",
                "fsType":"ext4",
                "fsTab":"/dev/sdak        /data36       ext4        noatime,acl,user_xattr,nofail    0  0"
            }                                          
        ]
        # 每一个块设备进行检查是否存在，存在则格式化并将挂载关系写入fstab
        for each in mountPointList:
            cmd_check = "if [ -b %s ];then echo 'exist';fi" % each["blockDevice"]
            status, output = commands.getstatusoutput(cmd_check)
            if status != 0:
                self.logRecord(output, 1)
                continue
            elif output == "":
                msg = "%s block device not exist." % each["blockDevice"]
                self.logRecord(msg)
                continue
            elif output == "exist":
                cmd_mkfs = "mkfs -F -t %s %s" % (each["fsType"], each["blockDevice"])
                status, output = commands.getstatusoutput(cmd_mkfs)
                if status != 0:
                    self.logRecord(output, 1)
                    fail_count += 1
                    continue
                else:
                    msg = "CMD: %s " % cmd_mkfs
                    self.logRecord(msg)
                    self.logRecord(output)

                    cmd_check = "grep %s /etc/fstab|wc -l" % each["blockDevice"]
                    status, output = commands.getstatusoutput(cmd_check)
                    if status != 0:
                        self.logRecord(output, 1)
                        continue
                    elif output == "0":
                        msg = "%s not in /etc/fstab. Add one." % each["blockDevice"]
                        self.logRecord(msg)
                        cmd_mkdir_addfstab = "umask 022;mkdir -p "+each["mountPoint"]+";echo '"+each["fsTab"]+"' >> /etc/fstab"
                        status, output = commands.getstatusoutput(cmd_mkdir_addfstab)
                        if status != 0:
                            self.logRecord(output, 1)
                            continue
                    else:
                        msg = "%s already in /etc/fstab. Do nothing." % each["blockDevice"]
                        self.logRecord(msg)
        # 不存在格式化失败的硬盘则返回成功
        if fail_count != 0:
            ret = self.uploadProgress("mountdata", 1)
            return ret
        else:
            ret = self.uploadProgress("mountdata", 0)
            return ret                

if __name__ == '__main__':
    status, output = getSN()
    if not status:
        print "Get Device SN Error: %s" % output
        sys.exit(1)
    handler = SystemInfo(output)
    handler.getDomainIP()
    for info in sys.argv[1:]:
        if "network" in info:
            if info.split("=")[1].upper() == "Y":
                handler.logRecord("Start Update Network ... ")
                ret = handler.updateNetwork()
                if ret is not False:
                    handler.logRecord("End Update Network ... ")
                else:
                    handler.logRecord("Update Network ... ", 1)
        elif "complete" in info:
             if info.split("=")[1].upper() == "Y":
                handler.logRecord("Start Update Service ... ")
                ret = handler.updateComplete()
                if ret is not False:
                    handler.logRecord("End Update Complete ... ")
                else:
                    handler.logRecord("Update Complete ... ", 1)
        elif "osuser" in info:
            if info.split("=")[1].upper() == "Y":
                handler.logRecord("Start Update OSUser ... ")
                ret = handler.updateOSUser()
                if ret is not False:
                    handler.logRecord("End Update OSUser ... ")
                else:
                    handler.logRecord("Update OSUser ... ", 1)
        elif "mountdata" in info:
            if info.split("=")[1].upper() == "Y":
                handler.logRecord("Start mount data ... ")
                ret = handler.mountData()
                if ret is not False:
                    handler.logRecord("End mount data ... ")
                else:
                    handler.logRecord("mount data ... ", 1)                    