# 从默认路径获取网络配置信息
$net_setting = @()
if (Test-Path "C:\firstboot\networkSetting.json") {
    $net_setting = Get-Content "C:\firstboot\networkSetting.json" | Out-String | ConvertFrom-Json
} else {
    exit 1
}

# 若配置文件有误则退出
if ($net_setting.status -ne "success") {
    exit 1
}

# get net adapter info {Name InterfaceDescription ifIndex Status MacAddress}
# eth0+eth1=bond1 -- 内网
# eth2+eth3=bond0 -- 外网
# 若只有2张网卡则只做bond1
$bond1_list = @()
$bond0_list = @()
$net_adapters = Get-NetAdapter | Sort-Object -Property ifIndex
if ($net_adapters.length -eq 2) {
    $bond1_list += $net_adapters[0].Name
    $bond1_list += $net_adapters[1].Name
    New-NetLbfoTeam -Name "Bond1" -TeamMembers ($bond1_list) -TeamingMode LACP -Confirm:$false
    foreach ($each_one in $net_setting.content) {
        if ($each_one.items.scope -eq "intranet") {
            netsh interface ip set address "Bond1" static $each_one.items.ip $each_one.items.netmask $each_one.items.gateway
        }
    }
}
# 若有4张网卡则做bond1 bond0，同时添加默认路由到外网卡
if ($net_adapters.length -eq 4) {
    $bond1_list += $net_adapters[0].Name
    $bond1_list += $net_adapters[1].Name
    $bond0_list += $net_adapters[2].Name
    $bond0_list += $net_adapters[3].Name    
    New-NetLbfoTeam -Name "Bond1" -TeamMembers ($bond1_list) -TeamingMode LACP -Confirm:$false
    New-NetLbfoTeam -Name "Bond0" -TeamMembers ($bond0_list) -TeamingMode LACP -Confirm:$false
    foreach ($each_one in $net_setting.content) {
        if ($each_one.items.scope -eq "intranet") {
            netsh interface ip set address "Bond1" static $each_one.items.ip $each_one.items.netmask $each_one.items.gateway
            route add -p 9.0.0.0 mask 255.0.0.0 $each_one.items.gateway
            route add -p 10.0.0.0 mask 255.0.0.0 $each_one.items.gateway
            route add -p 100.64.0.0 mask 255.192.0.0 $each_one.items.gateway
        }
        if ($each_one.items.scope -eq "extranet") {
            netsh interface ip set address "Bond0" static $each_one.items.ip $each_one.items.netmask $each_one.items.gateway
            route delete -p 0.0.0.0
            route add -p 0.0.0.0 mask 0.0.0.0 $each_one.items.gateway
        }
    }
}

# 添加防火墙规则，允许ping
$check_rule = Show-NetFirewallRule | findstr AllowPing
if ($check_rule.length -eq 0) {
    New-NetFirewallRule -Name "AllowPing" -DisplayName "AllowPing" -Direction Inbound -Action allow -Protocol ICMPv4
}