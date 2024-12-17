-- --------------------------------------------------------
-- 主机:                           x
-- 服务器版本:                        5.7.38-log - MySQL Community Server (GPL)
-- 服务器操作系统:                      Linux
-- HeidiSQL 版本:                  11.0.0.6063
-- --------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;


-- 导出 cloudboot_cyclone 的数据库结构
CREATE DATABASE IF NOT EXISTS `cloudboot_cyclone` /*!40100 DEFAULT CHARACTER SET utf8 */;
USE `cloudboot_cyclone`;

-- 导出  表 cloudboot_cyclone.api_log 结构
CREATE TABLE IF NOT EXISTS `api_log` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `method` varchar(64) NOT NULL COMMENT 'HTTP请求方法',
  `api` text NOT NULL COMMENT '接口方法',
  `description` varchar(128) DEFAULT NULL COMMENT '描述信息',
  `req_body` mediumtext COMMENT '请求参数',
  `status` varchar(16) DEFAULT NULL COMMENT '状态',
  `remote_addr` varchar(32) NOT NULL COMMENT '请求地址',
  `created_at` timestamp NULL DEFAULT NULL COMMENT '记录创建时间',
  `operator` varchar(64) DEFAULT NULL COMMENT '操作人',
  `msg` text COMMENT '返回信息',
  `result` text COMMENT '操作结果',
  `time` float(8,3) DEFAULT NULL COMMENT '操作耗时',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=37 DEFAULT CHARSET=utf8 COMMENT='API操作记录';

-- 正在导出表  cloudboot_cyclone.api_log 的数据：~36 rows (大约)
/*!40000 ALTER TABLE `api_log` DISABLE KEYS */;
REPLACE INTO `api_log` (`id`, `method`, `api`, `description`, `req_body`, `status`, `remote_addr`, `created_at`, `operator`, `msg`, `result`, `time`) VALUES
	(1, 'DELETE', '/api/cloudboot/v1/image-templates/2', '删除镜像模板', '', 'success', 'localhost:8083', '2023-12-06 21:33:03', 'admin', '操作成功', 'null', 0.002),
	(2, 'DELETE', '/api/cloudboot/v1/image-templates/1', '删除镜像模板', '', 'success', 'localhost:8083', '2023-12-06 21:33:04', 'admin', '操作成功', 'null', 0.001),
	(3, 'PUT', '/api/cloudboot/v1/device-setting-rules', '修改装机参数规则', '{"condition":"[{\\"value\\": [\\"x86\\"], \\"operator\\": \\"equal\\", \\"attribute\\": \\"arch\\"}]","action":"CentOS 7.9","rule_category":"os","id":9}', 'success', 'localhost:8083', '2023-12-06 21:35:46', 'admin', '操作成功', '{"id":9}', 0.001),
	(4, 'DELETE', '/api/cloudboot/v1/device-setting-rules', '删除装机参数规则', '{"ids":[8]}', 'success', 'localhost:8083', '2023-12-06 21:36:06', 'admin', '操作成功', '{"affected":1}', 0.001),
	(5, 'DELETE', '/api/cloudboot/v1/image-templates/8', '删除镜像模板', '', 'success', 'localhost:8083', '2023-12-06 21:36:59', 'admin', '操作成功', 'null', 0.001),
	(6, 'DELETE', '/api/cloudboot/v1/image-templates/1', '删除镜像模板', '', 'success', 'localhost:8083', '2023-12-06 21:37:01', 'admin', '操作成功', 'null', 0.001),
	(7, 'DELETE', '/api/cloudboot/v1/image-templates/3', '删除镜像模板', '', 'success', 'localhost:8083', '2023-12-06 21:37:04', 'admin', '操作成功', 'null', 0.001),
	(8, 'DELETE', '/api/cloudboot/v1/image-templates/4', '删除镜像模板', '', 'success', 'localhost:8083', '2023-12-06 21:37:06', 'admin', '操作成功', 'null', 0.001),
	(9, 'DELETE', '/api/cloudboot/v1/image-templates/5', '删除镜像模板', '', 'success', 'localhost:8083', '2023-12-06 21:37:09', 'admin', '操作成功', 'null', 0.001),
	(10, 'DELETE', '/api/cloudboot/v1/image-templates/7', '删除镜像模板', '', 'success', 'localhost:8083', '2023-12-06 21:37:12', 'admin', '操作成功', 'null', 0.001),
	(11, 'PUT', '/api/cloudboot/v1/image-templates/6', '修改镜像模板', '{"family":"CentOS","boot_mode":"uefi","arch":"x86_64","os_lifecycle":"active_default","name":"CentOS 7.9","url":"http://osinstall.idcos.com/images/centos7u9.tar.gz","username":"root","password":"Cyclone@1234","disks":[{"name":"/dev/sda","partitions":[{"size":"200","fstype":"vfat","mountpoint":"/boot/efi","name":"/dev/sda1"},{"size":"1024","fstype":"ext4","mountpoint":"/boot","name":"/dev/sda2"},{"size":"30720","fstype":"ext4","mountpoint":"/","name":"/dev/sda3"},{"size":"10240","fstype":"ext4","mountpoint":"/tmp","name":"/dev/sda4"},{"size":"5120","fstype":"ext4","mountpoint":"/home","name":"/dev/sda5"},{"size":"20480","fstype":"ext4","mountpoint":"/usr/local","name":"/dev/sda6"},{"size":"free","fstype":"ext4","mountpoint":"/data","name":"/dev/sda7"}]}],"pre_script":"","post_script":"# dowload commonSetting.py for setting items in OS\\ncurl -o /tmp/commonSetting.py \\"http://osinstall.idcos.com/scripts/commonSetting.py\\"\\n\\n# config network\\npython /tmp/commonSetting.py --network=Y \\n\\n# change root passwd\\npython /tmp/commonSetting.py --osuser=Y\\n\\n# complete\\npython /tmp/commonSetting.py --complete=Y\\n"}', 'success', 'localhost:8083', '2023-12-06 21:37:54', 'admin', '操作成功', '{"id":6}', 0.001),
	(12, 'PUT', '/api/cloudboot/v1/image-templates/6', '修改镜像模板', '{"family":"CentOS","boot_mode":"uefi","arch":"x86_64","os_lifecycle":"active","name":"CentOS 7.9","url":"http://osinstall.idcos.com/images/centos7u9.tar.gz","username":"root","password":"Cyclone@1234","disks":[{"name":"/dev/sda","partitions":[{"size":"200","fstype":"vfat","mountpoint":"/boot/efi","name":"/dev/sda1"},{"size":"1024","fstype":"ext4","mountpoint":"/boot","name":"/dev/sda2"},{"size":"30720","fstype":"ext4","mountpoint":"/","name":"/dev/sda3"},{"size":"10240","fstype":"ext4","mountpoint":"/tmp","name":"/dev/sda4"},{"size":"5120","fstype":"ext4","mountpoint":"/home","name":"/dev/sda5"},{"size":"20480","fstype":"ext4","mountpoint":"/usr/local","name":"/dev/sda6"},{"size":"free","fstype":"ext4","mountpoint":"/data","name":"/dev/sda7"}]}],"pre_script":"","post_script":"# dowload commonSetting.py for setting items in OS\\ncurl -o /tmp/commonSetting.py \\"http://osinstall.idcos.com/scripts/commonSetting.py\\"\\n\\n# config network\\npython /tmp/commonSetting.py --network=Y \\n\\n# change root passwd\\npython /tmp/commonSetting.py --osuser=Y\\n\\n# complete\\npython /tmp/commonSetting.py --complete=Y\\n"}', 'success', 'localhost:8083', '2023-12-06 21:38:04', 'admin', '操作成功', '{"id":6}', 0.001),
	(13, 'DELETE', '/api/cloudboot/v1/system-templates/27', '删除系统模板', '', 'success', 'localhost:8083', '2023-12-06 21:39:19', 'admin', '操作成功', 'null', 0.001),
	(14, 'DELETE', '/api/cloudboot/v1/system-templates/18', '删除系统模板', '', 'success', 'localhost:8083', '2023-12-06 21:39:21', 'admin', '操作成功', 'null', 0.001),
	(15, 'DELETE', '/api/cloudboot/v1/system-templates/19', '删除系统模板', '', 'success', 'localhost:8083', '2023-12-06 21:39:24', 'admin', '操作成功', 'null', 0.001),
	(16, 'DELETE', '/api/cloudboot/v1/system-templates/31', '删除系统模板', '', 'success', 'localhost:8083', '2023-12-06 21:39:30', 'admin', '操作成功', 'null', 0.001),
	(17, 'DELETE', '/api/cloudboot/v1/system-templates/14', '删除系统模板', '', 'success', 'localhost:8083', '2023-12-06 21:39:34', 'admin', '操作成功', 'null', 0.001),
	(18, 'DELETE', '/api/cloudboot/v1/system-templates/15', '删除系统模板', '', 'success', 'localhost:8083', '2023-12-06 21:39:37', 'admin', '操作成功', 'null', 0.001),
	(19, 'DELETE', '/api/cloudboot/v1/system-templates/3', '删除系统模板', '', 'success', 'localhost:8083', '2023-12-06 21:39:42', 'admin', '操作成功', 'null', 0.001),
	(20, 'DELETE', '/api/cloudboot/v1/system-templates/23', '删除系统模板', '', 'success', 'localhost:8083', '2023-12-06 21:39:45', 'admin', '操作成功', 'null', 0.001),
	(21, 'DELETE', '/api/cloudboot/v1/system-templates/30', '删除系统模板', '', 'success', 'localhost:8083', '2023-12-06 21:39:48', 'admin', '操作成功', 'null', 0.001),
	(22, 'DELETE', '/api/cloudboot/v1/system-templates/29', '删除系统模板', '', 'success', 'localhost:8083', '2023-12-06 21:39:53', 'admin', '操作成功', 'null', 0.001),
	(23, 'DELETE', '/api/cloudboot/v1/system-templates/37', '删除系统模板', '', 'success', 'localhost:8083', '2023-12-06 21:39:57', 'admin', '操作成功', 'null', 0.001),
	(24, 'DELETE', '/api/cloudboot/v1/system-templates/38', '删除系统模板', '', 'success', 'localhost:8083', '2023-12-06 21:40:01', 'admin', '操作成功', 'null', 0.001),
	(25, 'DELETE', '/api/cloudboot/v1/system-templates/13', '删除系统模板', '', 'success', 'localhost:8083', '2023-12-06 21:40:03', 'admin', '操作成功', 'null', 0.001),
	(26, 'DELETE', '/api/cloudboot/v1/system-templates/20', '删除系统模板', '', 'success', 'localhost:8083', '2023-12-06 21:40:29', 'admin', '操作成功', 'null', 0.001),
	(27, 'DELETE', '/api/cloudboot/v1/system-templates/26', '删除系统模板', '', 'success', 'localhost:8083', '2023-12-06 21:40:41', 'admin', '操作成功', 'null', 0.001),
	(28, 'DELETE', '/api/cloudboot/v1/system-templates/21', '删除系统模板', '', 'success', 'localhost:8083', '2023-12-06 21:41:28', 'admin', '操作成功', 'null', 0.001),
	(29, 'DELETE', '/api/cloudboot/v1/system-templates/16', '删除系统模板', '', 'success', 'localhost:8083', '2023-12-06 21:42:51', 'admin', '操作成功', 'null', 0.001),
	(30, 'PUT', '/api/cloudboot/v1/system-templates/2', '修改系统模板', '{"family":"BootOS","name":"bootos_arm64","boot_mode":"uefi","arch":"aarch64","os_lifecycle":"testing","username":"","password":"","pxe":"#!ipxe\\nkernel http://osinstall.idcos.com/bootos/aarch64/vmlinuz initrd=initrd.img console=tty0 selinux=0 biosdevname=0 SERVER_ADDR=http://osinstall.idcos.com LOOP_INTERVAL=5 DEVELOPER=1 BOOTIF=01-${netX/mac:hexhyp}\\ninitrd http://osinstall.idcos.com/bootos/aarch64/initrd.img\\nboot\\n","content":"#"}', 'success', 'localhost:8083', '2023-12-06 21:46:22', 'admin', '操作成功', '{"id":2}', 0.002),
	(31, 'DELETE', '/api/cloudboot/v1/device-setting-rules', '删除装机参数规则', '{"ids":[16]}', 'success', 'localhost:8083', '2023-12-06 21:56:12', 'admin', '操作成功', '{"affected":1}', 0.001),
	(32, 'PUT', '/api/cloudboot/v1/system-templates/17', '修改系统模板', '{"family":"CentOS","name":"CentOS_7.6_aarch64","boot_mode":"uefi","arch":"aarch64","os_lifecycle":"active_default","username":"root","password":"Cyclone@1234","pxe":"#!ipxe\\nkernel http://osinstall.idcos.com/centos/7.6/os/aarch64/images/pxeboot/vmlinuz initrd=initrd.img ksdevice=bootif ks=http://osinstall.idcos.com/api/cloudboot/v1/devices/{sn}/settings/system-template console=tty0 selinux=0 net.ifnames=0 biosdevname=0 BOOTIF=01-${netX/mac:hexhyp}\\ninitrd http://osinstall.idcos.com/centos/7.6/os/aarch64/images/pxeboot/initrd.img\\nboot","content":"install\\nurl --url=http://osinstall.idcos.com/centos/7.6/os/aarch64/\\nlang en_US.UTF-8\\nkeyboard us\\nnetwork --onboot yes --device bootif --bootproto dhcp --noipv6\\nrootpw  Cyclone@1234\\nfirewall --disabled\\nauthconfig --enableshadow --passalgo=sha512\\nselinux --disabled\\ntimezone Asia/Shanghai\\ntext\\nreboot\\nzerombr\\nbootloader --location=mbr\\nclearpart --all --initlabel\\npart /boot/efi --fstype=efi --size=200 --ondisk=sda \\npart /boot --fstype=xfs --size=1024 --ondisk=sda\\npart swap --size=8192 --ondisk=sda\\npart / --fstype=xfs --size=30720 --ondisk=sda\\npart /tmp --fstype=xfs --size=10240 --ondisk=sda\\npart /home --fstype=xfs --size=5120 --ondisk=sda\\npart /usr/local --fstype=xfs --size=20480 --ondisk=sda\\npart /data --fstype=xfs --size=1 --grow --ondisk=sda\\n\\n%packages --ignoremissing\\n@base\\n@core\\n@development\\ndmidecode\\n%end\\n\\n%pre\\n_sn=$(dmidecode -s system-serial-number 2>/dev/null | awk \'/^[^#]/ { print $1 }\')\\n\\ncurl -H \\"Content-Type: application/json\\" -X POST -d \\"{\\\\\\"title\\\\\\":\\\\\\"启动OS安装程序\\\\\\",\\\\\\"progress\\\\\\":0.6,\\\\\\"log\\\\\\":\\\\\\"SW5zdGFsbCBPUwo=\\\\\\"}\\" http://osinstall.idcos.com/api/cloudboot/v1/devices/$_sn/installations/progress\\ncurl -H \\"Content-Type: application/json\\" -X POST -d \\"{\\\\\\"title\\\\\\":\\\\\\"分区并安装软件包\\\\\\",\\\\\\"progress\\\\\\":0.7,\\\\\\"log\\\\\\":\\\\\\"SW5zdGFsbCBPUwo=\\\\\\"}\\" http://osinstall.idcos.com/api/cloudboot/v1/devices/$_sn/installations/progress\\n%end\\n\\n%post\\n# dowload commonSetting.py for setting items in OS\\ncurl -o /tmp/commonSetting.py \\"http://osinstall.idcos.com/scripts/commonSetting.py\\"\\n\\n# config network\\npython /tmp/commonSetting.py --network=Y \\n\\n# config osuser\\npython /tmp/commonSetting.py --osuser=Y\\n\\n# complete\\npython /tmp/commonSetting.py --complete=Y\\n%end"}', 'success', 'localhost:8083', '2023-12-07 09:17:50', 'admin', '操作成功', '{"id":17}', 0.003),
	(33, 'PUT', '/api/cloudboot/v1/system-templates/17', '修改系统模板', '{"family":"CentOS","name":"CentOS_7.6_aarch64","boot_mode":"uefi","arch":"aarch64","os_lifecycle":"active","username":"root","password":"Cyclone@1234","pxe":"#!ipxe\\nkernel http://osinstall.idcos.com/centos/7.6/os/aarch64/images/pxeboot/vmlinuz initrd=initrd.img ksdevice=bootif ks=http://osinstall.idcos.com/api/cloudboot/v1/devices/{sn}/settings/system-template console=tty0 selinux=0 net.ifnames=0 biosdevname=0 BOOTIF=01-${netX/mac:hexhyp}\\ninitrd http://osinstall.idcos.com/centos/7.6/os/aarch64/images/pxeboot/initrd.img\\nboot","content":"install\\nurl --url=http://osinstall.idcos.com/centos/7.6/os/aarch64/\\nlang en_US.UTF-8\\nkeyboard us\\nnetwork --onboot yes --device bootif --bootproto dhcp --noipv6\\nrootpw  Cyclone@1234\\nfirewall --disabled\\nauthconfig --enableshadow --passalgo=sha512\\nselinux --disabled\\ntimezone Asia/Shanghai\\ntext\\nreboot\\nzerombr\\nbootloader --location=mbr\\nclearpart --all --initlabel\\npart /boot/efi --fstype=efi --size=200 --ondisk=sda \\npart /boot --fstype=xfs --size=1024 --ondisk=sda\\npart swap --size=8192 --ondisk=sda\\npart / --fstype=xfs --size=30720 --ondisk=sda\\npart /tmp --fstype=xfs --size=10240 --ondisk=sda\\npart /home --fstype=xfs --size=5120 --ondisk=sda\\npart /usr/local --fstype=xfs --size=20480 --ondisk=sda\\npart /data --fstype=xfs --size=1 --grow --ondisk=sda\\n\\n%packages --ignoremissing\\n@base\\n@core\\n@development\\ndmidecode\\n%end\\n\\n%pre\\n_sn=$(dmidecode -s system-serial-number 2>/dev/null | awk \'/^[^#]/ { print $1 }\')\\n\\ncurl -H \\"Content-Type: application/json\\" -X POST -d \\"{\\\\\\"title\\\\\\":\\\\\\"启动OS安装程序\\\\\\",\\\\\\"progress\\\\\\":0.6,\\\\\\"log\\\\\\":\\\\\\"SW5zdGFsbCBPUwo=\\\\\\"}\\" http://osinstall.idcos.com/api/cloudboot/v1/devices/$_sn/installations/progress\\ncurl -H \\"Content-Type: application/json\\" -X POST -d \\"{\\\\\\"title\\\\\\":\\\\\\"分区并安装软件包\\\\\\",\\\\\\"progress\\\\\\":0.7,\\\\\\"log\\\\\\":\\\\\\"SW5zdGFsbCBPUwo=\\\\\\"}\\" http://osinstall.idcos.com/api/cloudboot/v1/devices/$_sn/installations/progress\\n%end\\n\\n%post\\n# dowload commonSetting.py for setting items in OS\\ncurl -o /tmp/commonSetting.py \\"http://osinstall.idcos.com/scripts/commonSetting.py\\"\\n\\n# config network\\npython /tmp/commonSetting.py --network=Y \\n\\n# config osuser\\npython /tmp/commonSetting.py --osuser=Y\\n\\n# complete\\npython /tmp/commonSetting.py --complete=Y\\n%end"}', 'success', 'localhost:8083', '2023-12-07 09:18:01', 'admin', '操作成功', '{"id":17}', 0.001),
	(34, 'PUT', '/api/cloudboot/v1/system-templates/32', '修改系统模板', '{"family":"EulerOS","name":"openEuler release 20.03 (LTS-SP3-x86_64)","boot_mode":"uefi","arch":"x86_64","os_lifecycle":"active_default","username":"root","password":"Cyclone@1234","pxe":"#!ipxe\\nkernel http://osinstall.idcos.com/openEuler/20.03_LTS_SP3/os/x86_64/images/pxeboot/vmlinuz initrd=initrd.img ksdevice=bootif ks=http://osinstall.idcos.com/api/cloudboot/v1/devices/{sn}/settings/system-template console=tty0 selinux=0 net.ifnames=0 biosdevname=0 BOOTIF=01-${netX/mac:hexhyp}\\ninitrd http://osinstall.idcos.com/openEuler/20.03_LTS_SP3/os/x86_64/images/pxeboot/initrd.img\\nboot","content":"#version=openEuler release 20.03 (LTS-SP3)\\n# Reboot after installation\\nreboot\\n# Use text mode install\\ntext\\n\\n\\n%pre --log=/tmp/kickstart_pre.log\\n\\necho \\"which dmidecode:\\"\\nwhich dmidecode\\n_sn=$(/usr/sbin/dmidecode -s system-serial-number 2>/dev/null | awk \'/^[^#]/ { print $1 }\')\\necho \\"system-serial-number：$_sn\\"\\necho \\"curl：http://osinstall.idcos.com/api/cloudboot/v1/devices/$_sn/installations/progress\\"\\ncurl -H \\"Content-Type: application/json\\" -X POST -d \\"{\\\\\\"title\\\\\\":\\\\\\"启动OS安装程序\\\\\\",\\\\\\"progress\\\\\\":0.6,\\\\\\"log\\\\\\":\\\\\\"SW5zdGFsbCBPUwo=\\\\\\"}\\" http://osinstall.idcos.com/api/cloudboot/v1/devices/$_sn/installations/progress\\ncurl -H \\"Content-Type: application/json\\" -X POST -d \\"{\\\\\\"title\\\\\\":\\\\\\"分区并安装软件包\\\\\\",\\\\\\"progress\\\\\\":0.7,\\\\\\"log\\\\\\":\\\\\\"SW5zdGFsbCBPUwo=\\\\\\"}\\" http://osinstall.idcos.com/api/cloudboot/v1/devices/$_sn/installations/progress\\n\\n%end\\n\\n\\n%post --nochroot --log=/mnt/sysimage/tmp/kickstart_post_nochroot.log\\n\\necho \\"Copying %pre stage log files\\"\\n/usr/bin/cp -rv /tmp/kickstart_pre.log /mnt/sysimage/tmp/\\n\\n%end\\n\\n\\n%post --log=/tmp/kickstart_post.log\\n#enable kdump\\n#sed  -i \\"s/ crashkernel=512M/ crashkernel=1024M,high /\\" /boot/efi/EFI/openEuler/grub.cfg\\n\\necho \\"dowload commonSetting.py for setting items in OS(chroot):\\"\\ncurl -o /tmp/commonSetting.py \\"http://osinstall.idcos.com/scripts/commonSetting.py\\"\\n\\ncurl -o /tmp/driver_upgrade.sh \\"http://osinstall.idcos.com/scripts/driver_upgrade.sh\\"\\nbash /tmp/driver_upgrade.sh\\n\\n# config network\\npython /tmp/commonSetting.py --network=Y\\n\\n# config osuser\\npython /tmp/commonSetting.py --osuser=Y\\n\\n# complete\\npython /tmp/commonSetting.py --complete=Y\\n%end\\n\\n%packages --ignoremissing\\n#@^minimal-environment\\n@^server-product-environment\\n@standard\\n@system-tools\\n#@development\\n#@performance\\n\\n%end\\n\\n# Keyboard layouts\\nkeyboard --vckeymap=us --xlayouts=\'us\'\\n# System language\\nlang en_US.UTF-8\\n\\n# Firewall configuration\\nfirewall --disabled\\n# Network information\\nnetwork  --bootproto=dhcp --device=bootif --ipv6=auto --activate\\nnetwork  --hostname=openeuler.webank\\n\\n# Use network installation\\nurl --url=\\"http://osinstall.idcos.com/openEuler/20.03_LTS_SP3/os/x86_64/\\"\\n\\n# System authorization information\\nauth --enableshadow --passalgo=sha512\\n# SELinux configuration\\nselinux --disabled\\n\\n# Do not configure the X Window System\\nskipx\\n# System services\\nservices --disabled=\\"chronyd\\"\\n\\nignoredisk --only-use=sda\\n# Partition clearing information\\nclearpart --all --initlabel --drives=sda\\n# Disk partitioning information\\npart /boot/efi --fstype=\\"efi\\" --ondisk=sda --size=200\\npart /boot --fstype=\\"ext4\\" --ondisk=sda --size=1024\\npart / --fstype=\\"ext4\\" --ondisk=sda --size=30720\\npart /tmp --fstype=\\"ext4\\" --ondisk=sda --size=10240\\npart /home --fstype=\\"ext4\\" --ondisk=sda --size=5120\\npart /usr/local --fstype=\\"ext4\\" --ondisk=sda --size=20480\\npart /data --fstype=\\"ext4\\" --size=1 --grow --ondisk=sda\\n\\n\\n# System timezone\\ntimezone Asia/Shanghai --utc\\n\\n# Root password\\nrootpw Cyclone@1234\\n\\n%addon com_redhat_kdump --disable --reserve-mb=\'128\'\\n\\n%end"}', 'success', 'localhost:8083', '2023-12-07 09:18:40', 'admin', '操作成功', '{"id":32}', 0.001),
	(35, 'PUT', '/api/cloudboot/v1/system-templates/28', '修改系统模板', '{"family":"EulerOS","name":"openEuler release 20.03 (LTS-SP3-aarch64)","boot_mode":"uefi","arch":"aarch64","os_lifecycle":"active","username":"root","password":"Cyclone@1234","pxe":"#!ipxe\\nkernel http://osinstall.idcos.com/openEuler/20.03_LTS_SP3/os/aarch64/images/pxeboot/vmlinuz initrd=initrd.img ksdevice=bootif ks=http://osinstall.idcos.com/api/cloudboot/v1/devices/{sn}/settings/system-template console=tty0 selinux=0 net.ifnames=0 biosdevname=0 BOOTIF=01-${netX/mac:hexhyp}\\ninitrd http://osinstall.idcos.com/openEuler/20.03_LTS_SP3/os/aarch64/images/pxeboot/initrd.img\\nboot","content":"#version=openEuler release 20.03 (LTS-SP3)\\n# Reboot after installation\\nreboot\\n# Use text mode install\\ntext\\n\\n\\n%pre --log=/tmp/kickstart_pre.log\\n\\necho \\"which dmidecode:\\"\\nwhich dmidecode\\n_sn=$(/usr/sbin/dmidecode -s system-serial-number 2>/dev/null | awk \'/^[^#]/ { print $1 }\')\\necho \\"system-serial-number：$_sn\\"\\necho \\"curl：http://osinstall.idcos.com/api/cloudboot/v1/devices/$_sn/installations/progress\\"\\ncurl -H \\"Content-Type: application/json\\" -X POST -d \\"{\\\\\\"title\\\\\\":\\\\\\"启动OS安装程序\\\\\\",\\\\\\"progress\\\\\\":0.6,\\\\\\"log\\\\\\":\\\\\\"SW5zdGFsbCBPUwo=\\\\\\"}\\" http://osinstall.idcos.com/api/cloudboot/v1/devices/$_sn/installations/progress\\ncurl -H \\"Content-Type: application/json\\" -X POST -d \\"{\\\\\\"title\\\\\\":\\\\\\"分区并安装软件包\\\\\\",\\\\\\"progress\\\\\\":0.7,\\\\\\"log\\\\\\":\\\\\\"SW5zdGFsbCBPUwo=\\\\\\"}\\" http://osinstall.idcos.com/api/cloudboot/v1/devices/$_sn/installations/progress\\n\\n%end\\n\\n\\n%post --nochroot --log=/mnt/sysimage/tmp/kickstart_post_nochroot.log\\n\\necho \\"Copying %pre stage log files\\"\\n/usr/bin/cp -rv /tmp/kickstart_pre.log /mnt/sysimage/tmp/\\n\\n%end\\n\\n\\n%post --log=/tmp/kickstart_post.log\\n#enable kdump\\n#sed  -i \\"s/ crashkernel=512M/ crashkernel=1024M,high /\\" /boot/efi/EFI/openEuler/grub.cfg\\n\\nvirsh net-destroy default\\n\\necho \\"dowload commonSetting.py for setting items in OS(chroot):\\"\\ncurl -o /tmp/commonSetting.py \\"http://osinstall.idcos.com/scripts/commonSetting.py\\"\\n#curl -o /tmp/driver_upgrade_aarch64.sh \\"http://osinstall.idcos.com/scripts/driver_upgrade_aarch64.sh\\"\\n#bash /tmp/driver_upgrade_aarch64.sh\\n# config network\\npython /tmp/commonSetting.py --network=Y\\n\\n# config osuser\\npython /tmp/commonSetting.py --osuser=Y\\n\\n# complete\\npython /tmp/commonSetting.py --complete=Y\\n%end\\n\\n%packages --ignoremissing\\n#@^minimal-environment\\n@^server-product-environment\\n@standard\\n@system-tools\\n#@development\\n#@performance\\n\\n%end\\n\\n# Keyboard layouts\\nkeyboard --vckeymap=us --xlayouts=\'us\'\\n# System language\\nlang en_US.UTF-8\\n\\n# Firewall configuration\\nfirewall --disabled\\n# Network information\\nnetwork  --bootproto=dhcp --device=bootif --ipv6=auto --activate\\nnetwork  --hostname=openeuler.webank\\n\\n# Use network installation\\nurl --url=\\"http://osinstall.idcos.com/openEuler/20.03_LTS_SP3/os/aarch64/\\"\\ndriverdisk --source=\\"http://osinstall.idcos.com/openEuler/20.03_LTS_SP3/drivers/aarch64/openEuler-20.03-LTS-SP3.iso\\"\\n# System authorization information\\nauth --enableshadow --passalgo=sha512\\n# SELinux configuration\\nselinux --disabled\\n\\n# Do not configure the X Window System\\nskipx\\n# System services\\nservices --disabled=\\"chronyd\\"\\n\\nignoredisk --only-use=sda\\n# Partition clearing information\\nclearpart --all --initlabel --drives=sda\\n# Disk partitioning information\\npart /boot/efi --fstype=\\"efi\\" --ondisk=sda --size=200\\npart /boot --fstype=\\"ext4\\" --ondisk=sda --size=1024\\npart / --fstype=\\"ext4\\" --ondisk=sda --size=30720\\npart /tmp --fstype=\\"ext4\\" --ondisk=sda --size=10240\\npart /home --fstype=\\"ext4\\" --ondisk=sda --size=5120\\npart /usr/local --fstype=\\"ext4\\" --ondisk=sda --size=20480\\npart /data --fstype=\\"ext4\\" --size=1 --grow --ondisk=sda\\n\\n\\n# System timezone\\ntimezone Asia/Shanghai --utc\\n\\n# Root password\\nrootpw Cyclone@1234\\n\\n%addon com_redhat_kdump --disable --reserve-mb=\'128\'\\n\\n%end"}', 'success', 'localhost:8083', '2023-12-07 09:18:58', 'admin', '操作成功', '{"id":28}', 0.001),
	(36, 'PUT', '/api/cloudboot/v1/system-templates/25', '修改系统模板', '{"family":"CentOS","name":"CentOS 7.9","boot_mode":"uefi","arch":"x86_64","os_lifecycle":"active","username":"root","password":"Cyclone@1234","pxe":"#!ipxe\\nkernel http://osinstall.idcos.com/centos/7.9/os/x86_64/images/pxeboot/vmlinuz initrd=initrd.img ksdevice=bootif ks=http://osinstall.idcos.com/api/cloudboot/v1/devices/${serial}/settings/system-template console=tty0 selinux=0 net.ifnames=0 biosdevname=0 BOOTIF=01-${netX/mac:hexhyp}\\ninitrd http://osinstall.idcos.com/centos/7.9/os/x86_64/images/pxeboot/initrd.img\\nboot\\n","content":"install\\nurl --url=http://osinstall.idcos.com/centos/7.9/os/x86_64/\\nlang en_US.UTF-8\\nkeyboard us\\nnetwork --onboot yes --device bootif --bootproto dhcp --noipv6\\nrootpw  Cyclone@1234\\nfirewall --disabled\\nauthconfig --enableshadow --passalgo=sha512\\nselinux --disabled\\ntimezone Asia/Shanghai\\ntext\\nreboot\\nzerombr\\nbootloader --location=mbr\\nclearpart --all --initlabel\\npart /boot/efi --fstype=efi --size=200 --ondisk=sda \\npart /boot --fstype=ext4 --size=1024 --ondisk=sda\\npart swap --size=8192 --ondisk=sda\\npart / --fstype=ext4 --size=30720 --ondisk=sda\\npart /tmp --fstype=ext4 --size=10240 --ondisk=sda\\npart /home --fstype=ext4 --size=5120 --ondisk=sda\\npart /usr/local --fstype=ext4 --size=20480 --ondisk=sda\\npart /data --fstype=ext4 --size=1 --grow --ondisk=sda\\n\\n%packages --ignoremissing\\n@base\\n@core\\n@development\\ndmidecode\\n%end\\n\\n%pre\\n_sn=$(dmidecode -s system-serial-number 2>/dev/null | awk \'/^[^#]/ { print $1 }\')\\n\\ncurl -H \\"Content-Type: application/json\\" -X POST -d \\"{\\\\\\"title\\\\\\":\\\\\\"启动OS安装程序\\\\\\",\\\\\\"progress\\\\\\":0.6,\\\\\\"log\\\\\\":\\\\\\"SW5zdGFsbCBPUwo=\\\\\\"}\\" http://osinstall.idcos.com/api/cloudboot/v1/devices/$_sn/installations/progress\\ncurl -H \\"Content-Type: application/json\\" -X POST -d \\"{\\\\\\"title\\\\\\":\\\\\\"分区并安装软件包\\\\\\",\\\\\\"progress\\\\\\":0.7,\\\\\\"log\\\\\\":\\\\\\"SW5zdGFsbCBPUwo=\\\\\\"}\\" http://osinstall.idcos.com/api/cloudboot/v1/devices/$_sn/installations/progress\\n%end\\n\\n%post\\n# dowload commonSetting.py for setting items in OS\\ncurl -o /tmp/commonSetting.py \\"http://osinstall.idcos.com/scripts/commonSetting.py\\"\\n\\n# config network\\npython /tmp/commonSetting.py --network=Y\\n\\n# change root passwd\\npython /tmp/commonSetting.py --osuser=Y\\n\\n# complete\\npython /tmp/commonSetting.py --complete=Y\\n%end"}', 'success', 'localhost:8083', '2023-12-07 09:19:23', 'admin', '操作成功', '{"id":25}', 0.001);
/*!40000 ALTER TABLE `api_log` ENABLE KEYS */;

-- 导出  表 cloudboot_cyclone.approval 结构
CREATE TABLE IF NOT EXISTS `approval` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `title` varchar(255) DEFAULT NULL COMMENT '审批标题',
  `type` varchar(255) DEFAULT NULL COMMENT '审批类型',
  `metadata` json DEFAULT NULL COMMENT '元数据(JSON对象)',
  `front_data` json DEFAULT NULL COMMENT '全量元数据',
  `initiator` varchar(255) DEFAULT NULL COMMENT '发起人ID',
  `approvers` json DEFAULT NULL COMMENT '审批人ID列表(JSON字符串数组)',
  `cc` json DEFAULT NULL COMMENT '抄送人ID列表(JSON字符串数组)',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注',
  `start_time` timestamp NULL DEFAULT NULL COMMENT '审批开始时间',
  `end_time` timestamp NULL DEFAULT NULL COMMENT '审批结束时间',
  `is_rejected` enum('yes','no') DEFAULT NULL COMMENT '审批单是否被拒绝',
  `status` enum('approval','completed','revoked','failure') DEFAULT NULL COMMENT '审批流水状态。approval-审批中; completed-已完成; revoked-已撤销;failure--失败',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='审批表';

-- 正在导出表  cloudboot_cyclone.approval 的数据：~0 rows (大约)
/*!40000 ALTER TABLE `approval` DISABLE KEYS */;
/*!40000 ALTER TABLE `approval` ENABLE KEYS */;

-- 导出  表 cloudboot_cyclone.approval_step 结构
CREATE TABLE IF NOT EXISTS `approval_step` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `approval_id` int(11) unsigned NOT NULL COMMENT '所属审批id',
  `approver` varchar(255) DEFAULT NULL COMMENT '审批步骤审批人ID',
  `next_approver` varchar(255) DEFAULT NULL COMMENT '下一审批步骤审批人ID',
  `title` varchar(255) DEFAULT NULL COMMENT '审批步骤标题',
  `action` enum('agree','reject') DEFAULT NULL COMMENT '审批步骤动作。agree-同意; reject-拒绝;',
  `remark` varchar(2048) DEFAULT NULL COMMENT '审批步骤备注',
  `start_time` timestamp NULL DEFAULT NULL COMMENT '审批步骤开始时间',
  `end_time` timestamp NULL DEFAULT NULL COMMENT '审批步骤结束时间',
  `hooks` json DEFAULT NULL COMMENT '当前步骤审批通过后的钩子字符串数组',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='审批步骤表';

-- 正在导出表  cloudboot_cyclone.approval_step 的数据：~0 rows (大约)
/*!40000 ALTER TABLE `approval_step` DISABLE KEYS */;
/*!40000 ALTER TABLE `approval_step` ENABLE KEYS */;

-- 导出  视图 cloudboot_cyclone.cabinet_power_info 结构
-- 创建临时表以解决视图依赖性错误
CREATE TABLE `cabinet_power_info` (
	`id` INT(11) UNSIGNED NOT NULL,
	`number` VARCHAR(255) NULL COMMENT '编号' COLLATE 'utf8_general_ci',
	`is_enabled` ENUM('yes','no') NULL COMMENT '是否启用' COLLATE 'utf8_general_ci',
	`is_powered` ENUM('yes','no') NULL COMMENT '是否开电' COLLATE 'utf8_general_ci',
	`idc` VARCHAR(255) NULL COMMENT '名称' COLLATE 'utf8_general_ci',
	`server_room` VARCHAR(255) NULL COMMENT '机房管理单元名称' COLLATE 'utf8_general_ci',
	`network_area` VARCHAR(255) NULL COMMENT '网络区域名称' COLLATE 'utf8_general_ci',
	`max_power` VARCHAR(255) NULL COMMENT '最大功率' COLLATE 'utf8_general_ci',
	`usite_total` BIGINT(21) NULL,
	`used_count` BIGINT(21) NULL,
	`free_count` BIGINT(21) NULL,
	`pre_occupied_count` BIGINT(21) NULL,
	`disabled_count` BIGINT(21) NULL,
	`known_used_power` DECIMAL(65,0) NULL,
	`free_power` VARCHAR(255) NULL COLLATE 'utf8_general_ci',
	`is_unknown_power_svr_count` DECIMAL(23,0) NULL
) ENGINE=MyISAM;

-- 导出  表 cloudboot_cyclone.component_log 结构
CREATE TABLE IF NOT EXISTS `component_log` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL COMMENT '记录创建时间',
  `updated_at` timestamp NULL DEFAULT NULL COMMENT '记录修改时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '记录删除时间',
  `sn` varchar(255) DEFAULT NULL COMMENT '设备序列号',
  `component` enum('agent','hw-server','peconfig','winconfig','peagent') DEFAULT NULL COMMENT '日志所属组件',
  `log` longtext COMMENT '日志内容',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='设备操作系统安装/配置过程中的文件日志';

-- 正在导出表  cloudboot_cyclone.component_log 的数据：~0 rows (大约)
/*!40000 ALTER TABLE `component_log` DISABLE KEYS */;
/*!40000 ALTER TABLE `component_log` ENABLE KEYS */;

-- 导出  表 cloudboot_cyclone.data_dict 结构
CREATE TABLE IF NOT EXISTS `data_dict` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL COMMENT '记录创建时间',
  `updated_at` timestamp NULL DEFAULT NULL COMMENT '记录修改时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '记录删除时间',
  `type` varchar(255) NOT NULL COMMENT '类型/标签',
  `name` varchar(255) NOT NULL COMMENT '数据字典键',
  `value` varchar(255) NOT NULL COMMENT '数据字典值',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注/说明',
  PRIMARY KEY (`id`),
  UNIQUE KEY `type` (`type`,`name`)
) ENGINE=InnoDB AUTO_INCREMENT=18 DEFAULT CHARSET=utf8 COMMENT='数据字典表';

-- 正在导出表  cloudboot_cyclone.data_dict 的数据：~17 rows (大约)
/*!40000 ALTER TABLE `data_dict` DISABLE KEYS */;
REPLACE INTO `data_dict` (`id`, `created_at`, `updated_at`, `deleted_at`, `type`, `name`, `value`, `remark`) VALUES
	(1, '2018-06-06 12:00:00', '2018-06-06 12:00:00', NULL, 'firmware', 'bios', 'BIOS', NULL),
	(2, '2018-06-06 12:00:00', '2018-06-06 12:00:00', NULL, 'firmware', 'ilo', 'HP iLO', NULL),
	(3, '2018-06-06 12:00:00', '2018-06-06 12:00:00', NULL, 'firmware', 'idrac', 'Dell iDRAC', NULL),
	(4, '2018-06-06 12:00:00', '2018-06-06 12:00:00', NULL, 'firmware', 'raid', 'RAID', NULL),
	(5, '2018-06-06 12:00:00', '2018-06-06 12:00:00', NULL, 'firmware', 'nic', '网卡', NULL),
	(6, '2018-06-06 12:00:00', '2018-06-06 12:00:00', NULL, 'firmware', 'hba', 'HBA', NULL),
	(7, '2018-06-06 12:00:00', '2018-06-06 12:00:00', NULL, 'firmware', 'backplane_dell', 'Dell Backplane背板', NULL),
	(8, '2018-06-06 12:00:00', '2018-06-06 12:00:00', NULL, 'os_family', 'RedHat', 'RedHat', NULL),
	(9, '2018-06-06 12:00:00', '2018-06-06 12:00:00', NULL, 'os_family', 'CentOS', 'CentOS', NULL),
	(10, '2018-06-06 12:00:00', '2018-06-06 12:00:00', NULL, 'os_family', 'SUSE', 'SUSE', NULL),
	(11, '2018-06-06 12:00:00', '2018-06-06 12:00:00', NULL, 'os_family', 'Ubuntu', 'Ubuntu', NULL),
	(12, '2018-06-06 12:00:00', '2018-06-06 12:00:00', NULL, 'os_family', 'Windows Server', 'Windows Server', NULL),
	(13, '2018-06-06 12:00:00', '2018-06-06 12:00:00', NULL, 'os_family', 'VMware ESXi', 'VMware ESXI', NULL),
	(14, '2018-06-06 12:00:00', '2018-06-06 12:00:00', NULL, 'os_family', 'XenServer', 'XenServer', NULL),
	(15, '2023-02-17 12:00:00', '2023-02-17 12:00:00', NULL, 'os_family', 'EulerOS', 'EulerOS', NULL),
	(16, '2023-02-17 12:00:00', '2023-02-17 12:00:00', NULL, 'os_family', 'BootOS', 'BootOS', NULL),
	(17, '2023-02-17 12:00:00', '2023-02-17 12:00:00', NULL, 'os_family', 'Custom', 'Custom', NULL);
/*!40000 ALTER TABLE `data_dict` ENABLE KEYS */;

-- 导出  表 cloudboot_cyclone.device 结构
CREATE TABLE IF NOT EXISTS `device` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL COMMENT '记录创建时间',
  `updated_at` timestamp NULL DEFAULT NULL COMMENT '记录修改时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '记录删除时间',
  `fixed_asset_number` varchar(255) DEFAULT NULL COMMENT '固定资产编号',
  `sn` varchar(255) DEFAULT NULL COMMENT '序列号',
  `vendor` varchar(255) DEFAULT NULL COMMENT '产品厂商',
  `model` varchar(255) DEFAULT NULL COMMENT '产品型号',
  `power_status` enum('power_on','power_off') DEFAULT 'power_off' COMMENT '电源状态。可选值: power_on-开电; power_off-关电;',
  `oob_accessible` enum('yes','no') DEFAULT NULL COMMENT '带外纳管状态',
  `arch` varchar(255) DEFAULT NULL COMMENT 'CPU硬件架构',
  `usage` varchar(255) DEFAULT NULL COMMENT '用途',
  `category` varchar(255) DEFAULT NULL COMMENT '设备类型',
  `idc_id` int(11) unsigned DEFAULT NULL COMMENT '数据中心ID',
  `server_room_id` int(11) unsigned DEFAULT NULL COMMENT '机房管理单元ID',
  `server_cabinet_id` int(11) unsigned DEFAULT NULL COMMENT '机架（柜）ID',
  `server_usite_id` int(11) unsigned DEFAULT NULL COMMENT '机位(U位)ID',
  `store_room_id` int(11) unsigned DEFAULT NULL COMMENT '库房管理单元ID',
  `virtual_cabinet_id` int(11) unsigned DEFAULT NULL COMMENT '虚拟货架ID',
  `hardware_remark` varchar(1024) DEFAULT NULL COMMENT '硬件说明',
  `raid_remark` varchar(1024) DEFAULT NULL COMMENT 'RAID说明',
  `oob_init` json DEFAULT NULL COMMENT 'OOB初始带外用户及密码',
  `started_at` datetime DEFAULT NULL COMMENT '启用时间',
  `onshelve_at` datetime DEFAULT NULL COMMENT '上架时间',
  `operation_status` enum('run_with_alarm','run_without_alarm','reinstalling','moving','pre_retire','retiring','retired','pre_deploy','on_shelve','recycling','maintaining','pre_move','in_store') DEFAULT NULL COMMENT '运营状态：运营中(需告警),运营中(无需告警),重装中,搬迁中,待退役,退役中,已退役,待部署,已上架,回收中,维护中,待搬迁',
  `oob_ip` varchar(256) DEFAULT NULL COMMENT '带外IP',
  `oob_user` varchar(256) DEFAULT NULL COMMENT '带外用户名',
  `oob_password` varchar(256) DEFAULT NULL COMMENT '带外密码',
  `cpu_sum` int(11) DEFAULT '0' COMMENT 'CPU总核数',
  `cpu` json DEFAULT NULL COMMENT 'CPU列表(JSON)',
  `memory_sum` int(11) DEFAULT '0' COMMENT '内存总容量(MB)',
  `memory` json DEFAULT NULL COMMENT '内存列表(JSON)',
  `disk_sum` int(11) DEFAULT '0' COMMENT '逻辑磁盘总容量(GB)',
  `disk` json DEFAULT NULL COMMENT '逻辑磁盘列表(JSON)',
  `disk_slot` json DEFAULT NULL COMMENT '物理磁盘列表(JSON)',
  `nic` json DEFAULT NULL COMMENT '网卡列表(JSON)',
  `nic_device` text COMMENT '网卡厂商',
  `bootos_ip` varchar(256) DEFAULT NULL COMMENT 'BootOS IP',
  `bootos_mac` varchar(256) DEFAULT NULL COMMENT 'BootOS 网口MAC地址',
  `motherboard` json DEFAULT NULL COMMENT '主板(JSON)',
  `raid` json DEFAULT NULL COMMENT 'RAID卡',
  `oob` json DEFAULT NULL COMMENT '带外',
  `bios` json DEFAULT NULL COMMENT 'BIOS',
  `fan` json DEFAULT NULL COMMENT '风扇',
  `power` json DEFAULT NULL COMMENT '电源',
  `hba` json DEFAULT NULL COMMENT 'HBA卡',
  `pci` json DEFAULT NULL COMMENT 'PCI插槽',
  `switch` json DEFAULT NULL COMMENT '交换机',
  `lldp` json DEFAULT NULL COMMENT '交换机信息(JSON)',
  `extra` json DEFAULT NULL COMMENT '用户自定义扩展(JSON)',
  `origin_node` varchar(255) DEFAULT NULL COMMENT '源节点',
  `origin_node_ip` varchar(255) DEFAULT NULL COMMENT '分布式部署下的源节点IP',
  `operation_user_id` varchar(255) DEFAULT NULL COMMENT '运维管理用户id',
  `creator` varchar(255) DEFAULT NULL COMMENT '记录创建者ID',
  `updater` varchar(255) DEFAULT NULL COMMENT '记录更新者ID',
  `remark` varchar(1024) DEFAULT NULL COMMENT '备注',
  `order_number` varchar(255) DEFAULT NULL COMMENT '订单编号',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_device_sn` (`sn`) USING BTREE,
  UNIQUE KEY `uk_device_idc_room_cabinet_usite` (`idc_id`,`server_room_id`,`server_cabinet_id`,`server_usite_id`),
  UNIQUE KEY `uk_device_fixed_asset_number` (`fixed_asset_number`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='物理机设备表';

-- 正在导出表  cloudboot_cyclone.device 的数据：~0 rows (大约)
/*!40000 ALTER TABLE `device` DISABLE KEYS */;
/*!40000 ALTER TABLE `device` ENABLE KEYS */;

-- 导出  表 cloudboot_cyclone.device_category 结构
CREATE TABLE IF NOT EXISTS `device_category` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL COMMENT '记录创建时间',
  `updated_at` timestamp NULL DEFAULT NULL COMMENT '记录修改时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '记录删除时间',
  `category` varchar(255) DEFAULT NULL COMMENT '设备类型',
  `hardware` varchar(255) DEFAULT NULL COMMENT '硬件配置',
  `central_processor_manufacture` varchar(255) DEFAULT NULL COMMENT '处理器生产商',
  `central_processor_arch` varchar(255) DEFAULT NULL COMMENT '处理器架构',
  `power` varchar(255) DEFAULT NULL COMMENT '功率',
  `unit` int(11) DEFAULT '2' COMMENT '设备U数',
  `is_fiti_eco_product` enum('yes','no') DEFAULT 'no' COMMENT '是否金融信创生态产品:yes-是;no-否;',
  `remark` varchar(1024) DEFAULT NULL COMMENT '备注',
  `creator` varchar(255) DEFAULT NULL COMMENT '记录创建者',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_category` (`category`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='设备类型表';

-- 正在导出表  cloudboot_cyclone.device_category 的数据：~0 rows (大约)
/*!40000 ALTER TABLE `device_category` DISABLE KEYS */;
/*!40000 ALTER TABLE `device_category` ENABLE KEYS */;

-- 导出  表 cloudboot_cyclone.device_lifecycle 结构
CREATE TABLE IF NOT EXISTS `device_lifecycle` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL COMMENT '记录创建时间',
  `updated_at` timestamp NULL DEFAULT NULL COMMENT '记录修改时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '记录删除时间',
  `fixed_asset_number` varchar(255) DEFAULT NULL COMMENT '固定资产编号',
  `sn` varchar(255) NOT NULL COMMENT '序列号',
  `asset_belongs` varchar(255) DEFAULT NULL COMMENT '资产归属',
  `owner` varchar(255) DEFAULT NULL COMMENT '负责人',
  `is_rental` enum('yes','no') DEFAULT 'no' COMMENT '是否租赁:yes-是;no-否;',
  `maintenance_service_provider` varchar(255) DEFAULT NULL COMMENT '维保服务供应商',
  `maintenance_service` text COMMENT '维保服务详细内容',
  `logistics_service` text COMMENT '物流服务详细内容',
  `maintenance_service_date_begin` datetime DEFAULT NULL COMMENT '维保服务起始日期',
  `maintenance_service_date_end` datetime DEFAULT NULL COMMENT '维保服务截止日期',
  `maintenance_service_status` enum('under_warranty','out_of_warranty','inactive') DEFAULT 'under_warranty' COMMENT '维保状态:在保-under_warranty;过保-out_of_warranty;未激活-inactive',
  `device_retired_date` date DEFAULT NULL COMMENT '设备退役时间',
  `lifecycle_log` json DEFAULT NULL COMMENT '变更记录(JSON)',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_sn` (`sn`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='设备生命周期表';

-- 正在导出表  cloudboot_cyclone.device_lifecycle 的数据：~0 rows (大约)
/*!40000 ALTER TABLE `device_lifecycle` DISABLE KEYS */;
/*!40000 ALTER TABLE `device_lifecycle` ENABLE KEYS */;

-- 导出  表 cloudboot_cyclone.device_log 结构
CREATE TABLE IF NOT EXISTS `device_log` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `created_at` timestamp NULL DEFAULT NULL COMMENT '记录创建时间',
  `updated_at` timestamp NULL DEFAULT NULL COMMENT '记录修改时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '记录删除时间',
  `sn` varchar(255) DEFAULT NULL COMMENT '设备序列号',
  `device_setting_id` int(11) unsigned NOT NULL COMMENT '设备装机参数ID',
  `type` enum('install','install_history') NOT NULL COMMENT '进度日志类型。install-安装中的进度日志; install_history-历史的装机进度日志;',
  `title` varchar(1024) NOT NULL COMMENT '进度日志标题',
  `content` longtext COMMENT '进度日志内容',
  PRIMARY KEY (`id`),
  KEY `idx_device_log_sn` (`sn`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='操作系统安装进度日志表';

-- 正在导出表  cloudboot_cyclone.device_log 的数据：~0 rows (大约)
/*!40000 ALTER TABLE `device_log` DISABLE KEYS */;
/*!40000 ALTER TABLE `device_log` ENABLE KEYS */;

-- 导出  表 cloudboot_cyclone.device_oob_history 结构
CREATE TABLE IF NOT EXISTS `device_oob_history` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL COMMENT '记录创建时间',
  `updated_at` timestamp NULL DEFAULT NULL COMMENT '记录修改时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '记录删除时间',
  `sn` varchar(255) DEFAULT NULL COMMENT '序列号',
  `username_old` varchar(255) NOT NULL COMMENT '修改前带外用户名',
  `username_new` varchar(255) NOT NULL COMMENT '修改后带外用户名',
  `password_old` varchar(255) NOT NULL COMMENT '修改前带外密码(加密)',
  `password_new` varchar(255) NOT NULL COMMENT '修改后带外密码(加密)',
  `creator` varchar(255) DEFAULT NULL COMMENT '记录创建者ID',
  `remark` varchar(1024) DEFAULT NULL COMMENT '备注',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='带外账户修改的操作历史记录表';

-- 正在导出表  cloudboot_cyclone.device_oob_history 的数据：~0 rows (大约)
/*!40000 ALTER TABLE `device_oob_history` DISABLE KEYS */;
/*!40000 ALTER TABLE `device_oob_history` ENABLE KEYS */;

-- 导出  表 cloudboot_cyclone.device_setting 结构
CREATE TABLE IF NOT EXISTS `device_setting` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `created_at` timestamp NULL DEFAULT NULL COMMENT '记录创建时间',
  `updated_at` timestamp NULL DEFAULT NULL COMMENT '记录修改时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '记录删除时间',
  `installation_start_time` datetime DEFAULT NULL COMMENT '安装开始时间',
  `installation_end_time` datetime DEFAULT NULL COMMENT '安装结束时间',
  `sn` varchar(255) NOT NULL COMMENT '设备序列号',
  `hardware_template_id` int(11) DEFAULT '0' COMMENT '硬件配置模板ID',
  `image_template_id` int(11) DEFAULT '0' COMMENT '镜像安装模板ID',
  `system_template_id` int(11) DEFAULT NULL COMMENT '系统安装模板ID',
  `need_extranet_ip` enum('yes','no') DEFAULT NULL COMMENT '是否需要外网IP',
  `extranet_ip_network_id` int(11) DEFAULT NULL COMMENT '外网IP所属网段ID',
  `extranet_ip` varchar(255) DEFAULT NULL COMMENT '外网IP',
  `intranet_ip_network_id` int(11) DEFAULT NULL COMMENT '内网IP所属网段ID',
  `intranet_ip` varchar(255) DEFAULT NULL COMMENT '内网IP',
  `install_type` enum('image','pxe') DEFAULT 'image' COMMENT '系统安装方式。image-镜像安装; pxe-PXE安装;',
  `install_progress` decimal(11,4) DEFAULT '0.0000' COMMENT '装机进度值',
  `status` enum('pre_install','installing','failure','success') DEFAULT NULL COMMENT '装机状态。pre_install-等待安装; installing-正在安装; failure-安装失败; success-安装成功;',
  `creator` varchar(255) DEFAULT NULL COMMENT '记录创建者ID',
  `updater` varchar(255) DEFAULT NULL COMMENT '记录更新者ID',
  `need_intranet_ipv6` enum('yes','no') DEFAULT 'no' COMMENT '是否需要内网IPv6',
  `need_extranet_ipv6` enum('yes','no') DEFAULT 'no' COMMENT '是否需要外网IPv6',
  `extranet_ipv6_network_id` int(11) DEFAULT '0' COMMENT '外网IPv6所属网段ID',
  `extranet_ipv6` varchar(255) DEFAULT NULL COMMENT '外网IPv6',
  `intranet_ipv6_network_id` int(11) DEFAULT '0' COMMENT '内网IPv6所属网段ID',
  `intranet_ipv6` varchar(255) DEFAULT NULL COMMENT '内网IPv6',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_device_setting_sn` (`sn`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='设备装机参数表';

-- 正在导出表  cloudboot_cyclone.device_setting 的数据：~0 rows (大约)
/*!40000 ALTER TABLE `device_setting` DISABLE KEYS */;
/*!40000 ALTER TABLE `device_setting` ENABLE KEYS */;

-- 导出  表 cloudboot_cyclone.device_setting_rule 结构
CREATE TABLE IF NOT EXISTS `device_setting_rule` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL COMMENT '记录创建时间',
  `updated_at` timestamp NULL DEFAULT NULL COMMENT '记录修改时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '记录删除时间',
  `condition` json NOT NULL COMMENT '前件',
  `action` text NOT NULL COMMENT '结论',
  `rule_category` enum('os','network','raid') NOT NULL COMMENT '规则分类',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=19 DEFAULT CHARSET=utf8 COMMENT='产生式部署参数规则表';

-- 正在导出表  cloudboot_cyclone.device_setting_rule 的数据：~16 rows (大约)
/*!40000 ALTER TABLE `device_setting_rule` DISABLE KEYS */;
REPLACE INTO `device_setting_rule` (`id`, `created_at`, `updated_at`, `deleted_at`, `condition`, `action`, `rule_category`) VALUES
	(1, NULL, NULL, NULL, '[{"value": ["DMZ"], "operator": "contains", "attribute": "physical_area"}]', 'la_wa', 'network'),
	(2, NULL, NULL, NULL, '[{"value": ["ECN"], "operator": "contains", "attribute": "physical_area"}]', 'la_wa', 'network'),
	(3, NULL, NULL, NULL, '[{"value": ["Management"], "operator": "contains", "attribute": "physical_area"}]', 'la', 'network'),
	(4, NULL, NULL, NULL, '[{"value": ["BDP"], "operator": "contains", "attribute": "physical_area"}]', 'la', 'network'),
	(5, NULL, NULL, NULL, '[{"value": ["ServerFarm"], "operator": "contains", "attribute": "physical_area"}]', 'la', 'network'),
	(6, NULL, NULL, NULL, '[{"value": ["Q-Zone"], "operator": "contains", "attribute": "physical_area"}]', 'la_wa', 'network'),
	(7, NULL, '2023-11-28 14:30:33', NULL, '[{"value": ["arm"], "operator": "equal", "attribute": "arch"}]', 'openEuler release 20.03 (LTS-SP3-aarch64)', 'os'),
	(9, NULL, '2023-12-06 21:35:46', NULL, '[{"value": ["x86"], "operator": "equal", "attribute": "arch"}]', 'CentOS 7.9', 'os'),
	(10, NULL, NULL, NULL, '[{"value": ["Y1-WBCN11-10G", "Y1-WBCG11-10G", "Y1-CN01-10G", "Y1-CN02-10G"], "operator": "in", "attribute": "category"}]', '[6+6]RAID50', 'raid'),
	(11, NULL, NULL, NULL, '[{"value": ["Y1-WBCG13-10G", "Y1-WBCG12-10G", "Y1-WBCG14-10G", "Y1-WBGI50-10G", "Y1-CG03-10G", "Y1-CG02-10G", "YI-GI30-10G"], "operator": "in", "attribute": "category"}]', 'RAID1', 'raid'),
	(12, NULL, NULL, NULL, '[{"value": ["Y1-WBSH15-10G", "Y1-WBSH13-10G", "Y1-WBSH14-10G", "Y1-WBGT1-100G", "SH2-10G"], "operator": "in", "attribute": "category"}]', 'NORAID', 'raid'),
	(13, NULL, NULL, NULL, '[{"value": ["Y1-WBBX12-10G"], "operator": "in", "attribute": "category"}, {"value": ["and"], "operator": "equal", "attribute": "logical_operator"}, {"value": ["h3c"], "operator": "equal", "attribute": "vendor"}]', '2RAID1+12RAID0', 'raid'),
	(14, NULL, NULL, NULL, '[{"value": ["Y1-WBBX12-10G", "Y1-BX02-10G"], "operator": "in", "attribute": "category"}, {"value": ["and"], "operator": "equal", "attribute": "logical_operator"}, {"value": ["lenovo", "inspur", "huawei", "dell"], "operator": "in", "attribute": "vendor"}]', '12RAID0+2RAID1', 'raid'),
	(15, NULL, NULL, NULL, '[{"value": ["GA-B30-10G"], "operator": "in", "attribute": "category"}]', '12RAID0+2RAID1', 'raid'),
	(17, NULL, '2021-09-02 17:44:36', NULL, '[{"value": ["GA-F10-10G"], "operator": "in", "attribute": "category"}]', '12RAID0+2RAID1+2RAID1', 'raid'),
	(18, '2023-11-28 14:31:54', '2023-11-28 14:31:54', NULL, '[{"value": ["GA-TG10-25G"], "operator": "in", "attribute": "category"}]', '[1-13]RAID0', 'raid');
/*!40000 ALTER TABLE `device_setting_rule` ENABLE KEYS */;

-- 导出  表 cloudboot_cyclone.dhcp_token_bucket 结构
CREATE TABLE IF NOT EXISTS `dhcp_token_bucket` (
  `token` varchar(255) NOT NULL DEFAULT '' COMMENT '令牌',
  `bucket` varchar(255) NOT NULL DEFAULT '' COMMENT '所属令牌桶(TOR)',
  `sn` varchar(255) DEFAULT NULL COMMENT '关联的设备序列号',
  `created_at` timestamp NULL DEFAULT NULL COMMENT '记录创建时间',
  `updated_at` timestamp NULL DEFAULT NULL COMMENT '记录修改时间',
  PRIMARY KEY (`token`),
  UNIQUE KEY `uk_sn_dhcp_token_bucket` (`sn`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='设备请求dhcp ip的令牌记录表';

-- 正在导出表  cloudboot_cyclone.dhcp_token_bucket 的数据：~0 rows (大约)
/*!40000 ALTER TABLE `dhcp_token_bucket` DISABLE KEYS */;
/*!40000 ALTER TABLE `dhcp_token_bucket` ENABLE KEYS */;

-- 导出  表 cloudboot_cyclone.hardware_template 结构
CREATE TABLE IF NOT EXISTS `hardware_template` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL COMMENT '记录创建时间',
  `updated_at` timestamp NULL DEFAULT NULL COMMENT '记录修改时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '记录删除时间',
  `name` varchar(255) DEFAULT NULL COMMENT '硬件模板名',
  `vendor` varchar(255) NOT NULL COMMENT '厂商',
  `model` varchar(255) NOT NULL COMMENT '型号',
  `builtin` enum('yes','no') NOT NULL DEFAULT 'yes' COMMENT '是否是系统内置模板',
  `data` json DEFAULT NULL COMMENT '硬件配置数据',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_hardware_template_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8 COMMENT='硬件配置模板表';

-- 正在导出表  cloudboot_cyclone.hardware_template 的数据：~10 rows (大约)
/*!40000 ALTER TABLE `hardware_template` DISABLE KEYS */;
REPLACE INTO `hardware_template` (`id`, `created_at`, `updated_at`, `deleted_at`, `name`, `vendor`, `model`, `builtin`, `data`) VALUES
	(1, '2018-11-27 09:51:35', '2018-11-27 09:51:35', NULL, '2RAID1+RAID5+RAID5', 'webank', 'webank', 'no', '[{"action": "clear_settings", "category": "raid", "metadata": {"clear": "ON", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid1", "drives": "1-2", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid5", "drives": "3-7", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid5", "drives": "8-", "controller_index": "0"}}, {"action": "init_disk", "category": "raid", "metadata": {"init": "ON", "controller_index": "0"}}]'),
	(2, '2018-11-27 09:51:35', '2018-11-27 09:51:35', NULL, '4NORAID', 'webank', 'webank', 'no', '[{"action": "clear_settings", "category": "raid", "metadata": {"clear": "ON", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid0", "drives": "1", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid0", "drives": "2", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid0", "drives": "3", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid0", "drives": "4", "controller_index": "0"}}, {"action": "init_disk", "category": "raid", "metadata": {"init": "ON", "controller_index": "0"}}]'),
	(3, '2018-11-27 09:51:35', '2018-11-27 09:51:35', NULL, '6NORAID', 'webank', 'webank', 'no', '[{"action": "clear_settings", "category": "raid", "metadata": {"clear": "ON", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid0", "drives": "1", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid0", "drives": "2", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid0", "drives": "3", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid0", "drives": "4", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid0", "drives": "5", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid0", "drives": "6", "controller_index": "0"}}, {"action": "init_disk", "category": "raid", "metadata": {"init": "ON", "controller_index": "0"}}]'),
	(4, '2018-11-27 09:51:35', '2018-11-27 09:51:35', NULL, '12NORAID', 'webank', 'webank', 'no', '[{"action": "clear_settings", "category": "raid", "metadata": {"clear": "ON", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid0", "drives": "1", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid0", "drives": "2", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid0", "drives": "3", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid0", "drives": "4", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid0", "drives": "5", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid0", "drives": "6", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid0", "drives": "7", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid0", "drives": "8", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid0", "drives": "9", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid0", "drives": "10", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid0", "drives": "11", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid0", "drives": "12", "controller_index": "0"}}, {"action": "init_disk", "category": "raid", "metadata": {"init": "ON", "controller_index": "0"}}]'),
	(5, '2018-11-27 09:51:35', '2018-11-27 09:51:35', NULL, '2+12NORAID', 'webank', 'webank', 'no', '[{"action": "clear_settings", "category": "raid", "metadata": {"clear": "ON", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid0", "drives": "1", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid0", "drives": "2", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid0", "drives": "3", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid0", "drives": "4", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid0", "drives": "5", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid0", "drives": "6", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid0", "drives": "7", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid0", "drives": "8", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid0", "drives": "9", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid0", "drives": "10", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid0", "drives": "11", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid0", "drives": "12", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid0", "drives": "13", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid0", "drives": "14", "controller_index": "0"}}, {"action": "init_disk", "category": "raid", "metadata": {"init": "ON", "controller_index": "0"}}]'),
	(6, '2018-11-27 09:51:35', '2018-11-27 09:51:35', NULL, '2RAID1+10RAID5', 'webank', 'webank', 'no', '[{"action": "clear_settings", "category": "raid", "metadata": {"clear": "ON", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid1", "drives": "1-2", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid5", "drives": "3-12", "controller_index": "0"}}, {"action": "init_disk", "category": "raid", "metadata": {"init": "ON", "controller_index": "0"}}]'),
	(7, '2018-11-27 09:51:35', '2018-11-27 09:51:35', NULL, '2RAID1+12单盘RAID0', 'webank', 'webank', 'no', '[{"action": "clear_settings", "category": "raid", "metadata": {"clear": "ON", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid1", "drives": "1-2", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid0", "drives": "3", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid0", "drives": "4", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid0", "drives": "5", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid0", "drives": "6", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid0", "drives": "7", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid0", "drives": "8", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid0", "drives": "9", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid0", "drives": "10", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid0", "drives": "11", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid0", "drives": "12", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid0", "drives": "13", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid0", "drives": "14", "controller_index": "0"}}, {"action": "init_disk", "category": "raid", "metadata": {"init": "ON", "controller_index": "0"}}]'),
	(8, '2018-11-27 09:51:35', '2018-11-27 09:51:35', NULL, 'RAID1', 'webank', 'webank', 'no', '[{"action": "clear_settings", "category": "raid", "metadata": {"clear": "ON", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid1", "drives": "all", "controller_index": "0"}}, {"action": "init_disk", "category": "raid", "metadata": {"init": "ON", "controller_index": "0"}}]'),
	(9, '2018-11-27 09:51:35', '2018-11-27 09:51:35', NULL, 'RAID5', 'webank', 'webank', 'no', '[{"action": "clear_settings", "category": "raid", "metadata": {"clear": "ON", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid5", "drives": "all", "controller_index": "0"}}, {"action": "init_disk", "category": "raid", "metadata": {"init": "ON", "controller_index": "0"}}]'),
	(10, '2018-11-27 09:51:35', '2018-11-27 09:51:35', NULL, 'RAID10', 'webank', 'webank', 'no', '[{"action": "clear_settings", "category": "raid", "metadata": {"clear": "ON", "controller_index": "0"}}, {"action": "create_array", "category": "raid", "metadata": {"level": "raid10", "drives": "all", "controller_index": "0"}}, {"action": "init_disk", "category": "raid", "metadata": {"init": "ON", "controller_index": "0"}}]');
/*!40000 ALTER TABLE `hardware_template` ENABLE KEYS */;

-- 导出  表 cloudboot_cyclone.idc 结构
CREATE TABLE IF NOT EXISTS `idc` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL COMMENT '记录创建时间',
  `updated_at` timestamp NULL DEFAULT NULL COMMENT '记录修改时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '记录删除时间',
  `name` varchar(255) DEFAULT NULL COMMENT '名称',
  `usage` enum('production','disaster_recovery','pre_production','testing') DEFAULT NULL COMMENT '用途。production-生产; disaster_recovery-容灾; pre_production-准生产; testing-测试;',
  `first_server_room` json DEFAULT NULL COMMENT '一级机房',
  `vendor` varchar(255) DEFAULT NULL COMMENT '供应商',
  `status` enum('under_construction','accepted','production','abolished') DEFAULT NULL COMMENT '状态。under_construction-建设中; accepted-已验收; production-已投产; abolished-已裁撤;',
  `creator` varchar(255) DEFAULT NULL COMMENT '记录创建者ID',
  `updater` varchar(255) DEFAULT NULL COMMENT '记录更新者ID',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_idc` (`name`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='数据中心信息表';

-- 正在导出表  cloudboot_cyclone.idc 的数据：~0 rows (大约)
/*!40000 ALTER TABLE `idc` DISABLE KEYS */;
/*!40000 ALTER TABLE `idc` ENABLE KEYS */;

-- 导出  表 cloudboot_cyclone.image_template 结构
CREATE TABLE IF NOT EXISTS `image_template` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `created_at` timestamp NULL DEFAULT NULL COMMENT '记录创建时间',
  `updated_at` timestamp NULL DEFAULT NULL COMMENT '记录修改时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '记录删除时间',
  `family` varchar(255) DEFAULT NULL COMMENT '族系',
  `name` varchar(255) NOT NULL COMMENT '模板名',
  `boot_mode` enum('legacy_bios','uefi') DEFAULT 'legacy_bios' COMMENT '启动模式。可选值: legacy_bios-传统BIOS模式; uefi-UEFI模式;',
  `image_url` varchar(1024) DEFAULT NULL COMMENT '镜像下载地址',
  `username` varchar(255) DEFAULT NULL COMMENT '操作系统用户名',
  `password` varchar(255) DEFAULT NULL COMMENT '操作系统用户密码',
  `partition` longtext COMMENT '分区配置',
  `post_script` longtext COMMENT 'post脚本',
  `pre_script` longtext COMMENT 'pre脚本',
  `os_lifecycle` enum('testing','active_default','active','containment','end_of_life') DEFAULT 'testing' COMMENT 'OS生命周期：Testing|Active(Default)|Active|Containment|EOL',
  `arch` enum('unknown','x86_64','aarch64') DEFAULT NULL COMMENT 'OS架构',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_image_template_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8 COMMENT='镜像安装模板表';

-- 正在导出表  cloudboot_cyclone.image_template 的数据：~2 rows (大约)
/*!40000 ALTER TABLE `image_template` DISABLE KEYS */;
REPLACE INTO `image_template` (`id`, `created_at`, `updated_at`, `deleted_at`, `family`, `name`, `boot_mode`, `image_url`, `username`, `password`, `partition`, `post_script`, `pre_script`, `os_lifecycle`, `arch`) VALUES
	(2, '2019-09-20 10:09:58', '2023-03-16 16:01:15', NULL, 'SUSE', 'SUSE Linux Enterprise Server 12 SP3', 'uefi', 'http://osinstall.idcos.com/images/sles12sp3.tar.gz', 'root', 'Cyclone@1234', '[{"name":"/dev/sda","partitions":[{"name":"/dev/sda1","size":"200","fstype":"vfat","mountpoint":"/boot/efi"},{"name":"/dev/sda2","size":"1024","fstype":"ext4","mountpoint":"/boot"},{"name":"/dev/sda3","size":"8192","fstype":"swap","mountpoint":"swap"},{"name":"/dev/sda4","size":"30720","fstype":"xfs","mountpoint":"/"},{"name":"/dev/sda5","size":"10240","fstype":"xfs","mountpoint":"/tmp"},{"name":"/dev/sda6","size":"5120","fstype":"xfs","mountpoint":"/home"},{"name":"/dev/sda7","size":"20480","fstype":"xfs","mountpoint":"/usr/local"},{"name":"/dev/sda8","size":"free","fstype":"xfs","mountpoint":"/data"}]}]', '# dowload commonSetting.py for setting items in OS\ncurl -o /tmp/commonSetting.py "http://osinstall.idcos.com/scripts/commonSetting.py"\n\n# config network\npython /tmp/commonSetting.py --network=Y \n\n# change root passwd\npython /tmp/commonSetting.py --osuser=Y\n\n# complete\npython /tmp/commonSetting.py --complete=Y\n', '', 'active', 'x86_64'),
	(6, '2021-03-22 17:35:27', '2023-12-06 21:38:04', NULL, 'CentOS', 'CentOS 7.9', 'uefi', 'http://osinstall.idcos.com/images/centos7u9.tar.gz', 'root', 'Cyclone@1234', '[{"name":"/dev/sda","partitions":[{"name":"/dev/sda1","size":"200","fstype":"vfat","mountpoint":"/boot/efi"},{"name":"/dev/sda2","size":"1024","fstype":"ext4","mountpoint":"/boot"},{"name":"/dev/sda3","size":"30720","fstype":"ext4","mountpoint":"/"},{"name":"/dev/sda4","size":"10240","fstype":"ext4","mountpoint":"/tmp"},{"name":"/dev/sda5","size":"5120","fstype":"ext4","mountpoint":"/home"},{"name":"/dev/sda6","size":"20480","fstype":"ext4","mountpoint":"/usr/local"},{"name":"/dev/sda7","size":"free","fstype":"ext4","mountpoint":"/data"}]}]', '# dowload commonSetting.py for setting items in OS\ncurl -o /tmp/commonSetting.py "http://osinstall.idcos.com/scripts/commonSetting.py"\n\n# config network\npython /tmp/commonSetting.py --network=Y \n\n# change root passwd\npython /tmp/commonSetting.py --osuser=Y\n\n# complete\npython /tmp/commonSetting.py --complete=Y\n', '', 'active', 'x86_64');
/*!40000 ALTER TABLE `image_template` ENABLE KEYS */;

-- 导出  表 cloudboot_cyclone.inspection 结构
CREATE TABLE IF NOT EXISTS `inspection` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL COMMENT '记录创建时间',
  `updated_at` timestamp NULL DEFAULT NULL COMMENT '记录修改时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '记录删除时间',
  `job_id` char(36) NOT NULL DEFAULT '' COMMENT '所属任务ID',
  `start_time` timestamp NULL DEFAULT NULL COMMENT '任务执行开始时间',
  `end_time` timestamp NULL DEFAULT NULL COMMENT '任务执行结束时间',
  `origin_node` varchar(255) DEFAULT NULL COMMENT '源节点',
  `sn` varchar(255) DEFAULT NULL COMMENT '设备序列号',
  `running_status` enum('running','done') DEFAULT NULL COMMENT '任务运行状态',
  `health_status` enum('nominal','warning','critical','unknown') DEFAULT 'unknown' COMMENT '设备健康状况: nominal-正常; warning-警告; critical-异常; unknown-未知;',
  `error` varchar(2048) DEFAULT NULL COMMENT '任务执行错误信息',
  `ipmi_result` json DEFAULT NULL COMMENT '硬件巡检结果(ipmi获取的传感器信息)。由id、name、type、state、reading、units、event等属性组成对象数组。',
  `ipmisel_result` json DEFAULT NULL COMMENT 'IPMI系统事件数据 ID,Date,Time,Name,Type,State,Event',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `insp_sn` (`sn`) USING BTREE,
  KEY `insp_start_time` (`start_time`) USING BTREE,
  KEY `insp_end_time` (`end_time`) USING BTREE,
  KEY `insp_job_id` (`job_id`) USING BTREE,
  KEY `insp_running_stat` (`running_status`) USING BTREE,
  KEY `insp_health_stat` (`health_status`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='硬件巡检执行明细表';

-- 正在导出表  cloudboot_cyclone.inspection 的数据：~0 rows (大约)
/*!40000 ALTER TABLE `inspection` DISABLE KEYS */;
/*!40000 ALTER TABLE `inspection` ENABLE KEYS */;

-- 导出  表 cloudboot_cyclone.ip 结构
CREATE TABLE IF NOT EXISTS `ip` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL COMMENT '记录创建时间',
  `updated_at` timestamp NULL DEFAULT NULL COMMENT '记录修改时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '记录删除时间',
  `ip_network_id` int(11) unsigned DEFAULT NULL COMMENT '网段ID',
  `ip` varchar(255) DEFAULT NULL COMMENT 'IP地址',
  `category` enum('business','pxe') DEFAULT NULL COMMENT 'IP类别。pxe-PXE用IP; business-业务用IP;',
  `is_used` enum('yes','no','disabled') DEFAULT 'no' COMMENT '是否已被使用。yes-是; no-否; disabled-不可用;',
  `sn` varchar(255) DEFAULT NULL COMMENT '使用IP的设备序列号',
  `scope` enum('intranet','extranet') DEFAULT NULL COMMENT 'IP作用范围(分配给内网/外网)',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注',
  `updater` varchar(255) DEFAULT NULL COMMENT '记录更新者ID',
  `release_date` date DEFAULT NULL COMMENT '释放日期',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_ip_ip_network_id_ip` (`ip_network_id`,`ip`),
  KEY `idx_ip_sn` (`sn`),
  KEY `idx_ip_network_id_category_is_used` (`ip_network_id`,`category`,`is_used`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='IP信息表';

-- 正在导出表  cloudboot_cyclone.ip 的数据：~0 rows (大约)
/*!40000 ALTER TABLE `ip` DISABLE KEYS */;
/*!40000 ALTER TABLE `ip` ENABLE KEYS */;

-- 导出  表 cloudboot_cyclone.ip_network 结构
CREATE TABLE IF NOT EXISTS `ip_network` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL COMMENT '记录创建时间',
  `updated_at` timestamp NULL DEFAULT NULL COMMENT '记录修改时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '记录删除时间',
  `idc_id` int(11) NOT NULL COMMENT '数据中心ID',
  `server_room_id` int(11) unsigned NOT NULL COMMENT '机房管理单元ID',
  `category` enum('ilo','tgw_intranet','tgw_extranet','intranet','extranet','v_intranet','v_extranet') DEFAULT NULL COMMENT '网段类别。ilo-服务器ILO; tgw_intranet-服务器TGW内网; tgw_extranet-服务器TGW外网; intranet-服务器普通内网; extranet-服务器普通外网; v_intranet-服务器虚拟化内网; v_extranet-服务器虚拟化外网;',
  `cidr` varchar(255) DEFAULT NULL COMMENT 'CIDR网段',
  `netmask` varchar(255) DEFAULT NULL COMMENT '掩码',
  `gateway` varchar(255) DEFAULT NULL COMMENT '网关',
  `ip_pool` varchar(1024) DEFAULT NULL COMMENT 'IP资源池(含PXE IP资源池)',
  `pxe_pool` varchar(1024) DEFAULT NULL COMMENT 'PXE IP资源池',
  `switches` json DEFAULT NULL COMMENT '网段作用范围内的交换机固定资产编号字符串数组',
  `vlan` varchar(255) DEFAULT NULL COMMENT 'VLAN',
  `creator` varchar(255) DEFAULT NULL COMMENT '记录创建者ID',
  `updater` varchar(255) DEFAULT NULL COMMENT '记录更新者ID',
  `version` enum('ipv4','ipv6') DEFAULT 'ipv4' COMMENT 'IP协议版本:v4v6',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_ip_network_room_category_cidr` (`server_room_id`,`category`,`cidr`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='网段信息表';

-- 正在导出表  cloudboot_cyclone.ip_network 的数据：~0 rows (大约)
/*!40000 ALTER TABLE `ip_network` DISABLE KEYS */;
/*!40000 ALTER TABLE `ip_network` ENABLE KEYS */;

-- 导出  表 cloudboot_cyclone.job 结构
CREATE TABLE IF NOT EXISTS `job` (
  `id` char(36) NOT NULL DEFAULT '' COMMENT '任务全局唯一ID',
  `created_at` timestamp NULL DEFAULT NULL COMMENT '记录创建时间',
  `updated_at` timestamp NULL DEFAULT NULL COMMENT '记录修改时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '记录删除时间',
  `builtin` enum('yes','no') DEFAULT 'no' COMMENT '是否是内建任务。yes-是; no-否;',
  `title` varchar(255) DEFAULT NULL COMMENT '任务描述',
  `category` enum('inspection','installation_timeout','release_ip','auto_deploy','mailsend','update_device_lifecycle') DEFAULT NULL,
  `rate` enum('immediately','fixed_delay','fixed_rate') DEFAULT NULL COMMENT '任务执行频率。immediately-立即执行; fixed_delay-延迟执行; fixed_rate-固定频率(周期性)执行',
  `cron` varchar(255) DEFAULT NULL COMMENT 'cron表达式（一次性任务该值为空）',
  `cron_render` text COMMENT 'cron表达式UI渲染信息',
  `next_run_time` timestamp NULL DEFAULT NULL COMMENT '任务的下一次运行时间',
  `target` json DEFAULT NULL COMMENT '任务的作用目标JSON对象（如设备列表）',
  `status` enum('running','paused','stoped','deleted') DEFAULT NULL COMMENT 'running-运行中; paused-已暂停; stoped-已停止; deleted-已删除;',
  `creator` varchar(255) DEFAULT NULL COMMENT '任务创建者',
  `updater` varchar(255) DEFAULT NULL COMMENT '任务更新者',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='一次性/周期性任务表';

-- 正在导出表  cloudboot_cyclone.job 的数据：~6 rows (大约)
/*!40000 ALTER TABLE `job` DISABLE KEYS */;
REPLACE INTO `job` (`id`, `created_at`, `updated_at`, `deleted_at`, `builtin`, `title`, `category`, `rate`, `cron`, `cron_render`, `next_run_time`, `target`, `status`, `creator`, `updater`) VALUES
	('04d0141e-aaed-4c66-8389-3eabb9450b0d', '2022-06-07 15:00:00', NULL, NULL, 'yes', '更新设备维保状态', 'update_device_lifecycle', 'fixed_rate', '0 0 13 * * ?', NULL, NULL, NULL, 'paused', NULL, NULL),
	('8467ac30808311eba709fd0db572a67f', '2021-03-15 15:00:00', NULL, NULL, 'yes', '邮件-待部署设备信息', 'mailsend', 'fixed_rate', '0 30 8 * * ?', NULL, NULL, NULL, 'paused', NULL, NULL),
	('bbd2b0c07c0511eb954efd0db572a67f', '2021-03-15 15:00:00', NULL, NULL, 'yes', '邮件-过保设备信息', 'mailsend', 'fixed_rate', '0 30 9 * * ?', NULL, NULL, NULL, 'paused', NULL, NULL),
	('d801f154c884437a805854412308f3ee', '2018-11-26 01:42:13', NULL, NULL, 'yes', '装机超时处理任务', 'installation_timeout', 'fixed_rate', '0 0/5 * * * ?', NULL, NULL, NULL, 'paused', NULL, NULL),
	('f0ee9638182a464087e6608c95093b27', '2020-07-10 15:00:00', NULL, NULL, 'yes', '释放IP', 'release_ip', 'fixed_rate', '0 0 3 * * ?', NULL, NULL, NULL, 'paused', NULL, NULL),
	('f6d66383cfcb452f9848b3fabf2221d1', '2020-09-08 15:00:00', NULL, NULL, 'yes', '自动部署', 'auto_deploy', 'fixed_rate', '0 0 3 * * ?', NULL, NULL, NULL, 'paused', NULL, NULL);
/*!40000 ALTER TABLE `job` ENABLE KEYS */;

-- 导出  表 cloudboot_cyclone.network_area 结构
CREATE TABLE IF NOT EXISTS `network_area` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL COMMENT '记录创建时间',
  `updated_at` timestamp NULL DEFAULT NULL COMMENT '记录修改时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '记录删除时间',
  `idc_id` int(11) DEFAULT NULL COMMENT '所属数据中心ID',
  `server_room_id` int(11) unsigned NOT NULL COMMENT '所属机房ID',
  `name` varchar(255) DEFAULT NULL COMMENT '网络区域名称',
  `physical_area` json DEFAULT NULL COMMENT '关联的物理区域列表',
  `status` enum('nonproduction','production','offline') DEFAULT NULL COMMENT '状态。nonproduction-未投产; production-已投产; offline-已下线(回收)',
  `creator` varchar(255) DEFAULT NULL COMMENT '记录创建者ID',
  `updater` varchar(255) DEFAULT NULL COMMENT '记录更新者ID',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_idc_room_area` (`name`,`server_room_id`,`idc_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='网络区域信息表';

-- 正在导出表  cloudboot_cyclone.network_area 的数据：~0 rows (大约)
/*!40000 ALTER TABLE `network_area` DISABLE KEYS */;
/*!40000 ALTER TABLE `network_area` ENABLE KEYS */;

-- 导出  表 cloudboot_cyclone.network_device 结构
CREATE TABLE IF NOT EXISTS `network_device` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL COMMENT '记录创建时间',
  `updated_at` timestamp NULL DEFAULT NULL COMMENT '记录修改时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '记录删除时间',
  `idc_id` int(11) unsigned NOT NULL COMMENT '所属数据中心ID',
  `server_room_id` int(11) NOT NULL COMMENT '机房管理单元ID',
  `server_cabinet_id` int(11) unsigned NOT NULL COMMENT '所属机架(柜)ID',
  `fixed_asset_number` varchar(255) DEFAULT NULL COMMENT '固定资产编号',
  `sn` varchar(255) DEFAULT NULL COMMENT '序列号',
  `type` enum('switch') DEFAULT NULL COMMENT '类型。switch-交换机;',
  `tor` varchar(255) DEFAULT NULL COMMENT 'TOR名称',
  `name` varchar(255) DEFAULT NULL COMMENT '设备名称',
  `model` varchar(255) DEFAULT NULL COMMENT '设备型号',
  `vendor` varchar(255) DEFAULT NULL COMMENT '厂商',
  `os` varchar(255) DEFAULT NULL COMMENT '操作系统',
  `usage` varchar(1024) DEFAULT NULL COMMENT '用途',
  `creator` varchar(255) DEFAULT NULL COMMENT '记录创建者ID',
  `updater` varchar(255) DEFAULT NULL COMMENT '记录更新者ID',
  `status` varchar(255) DEFAULT NULL COMMENT '状态',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_network_device_fixed_asset_number` (`fixed_asset_number`),
  UNIQUE KEY `uk_network_device_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='网络设备信息表';

-- 正在导出表  cloudboot_cyclone.network_device 的数据：~0 rows (大约)
/*!40000 ALTER TABLE `network_device` DISABLE KEYS */;
/*!40000 ALTER TABLE `network_device` ENABLE KEYS */;

-- 导出  表 cloudboot_cyclone.operation_log 结构
CREATE TABLE IF NOT EXISTS `operation_log` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `created_at` timestamp NULL DEFAULT NULL COMMENT '记录创建时间',
  `operator` varchar(255) NOT NULL COMMENT '操作人ID',
  `url` varchar(255) NOT NULL COMMENT '请求url',
  `http_method` varchar(8) NOT NULL COMMENT '请求方式',
  `category_name` varchar(64) NOT NULL COMMENT '操作类别名称',
  `category_code` varchar(64) NOT NULL COMMENT '操作类别编码',
  `source` json DEFAULT NULL COMMENT '操作前的源数据',
  `destination` json DEFAULT NULL COMMENT '操作后的数据',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户操作日志表';

-- 正在导出表  cloudboot_cyclone.operation_log 的数据：~0 rows (大约)
/*!40000 ALTER TABLE `operation_log` DISABLE KEYS */;
/*!40000 ALTER TABLE `operation_log` ENABLE KEYS */;

-- 导出  表 cloudboot_cyclone.order 结构
CREATE TABLE IF NOT EXISTS `order` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL COMMENT '记录创建时间',
  `updated_at` timestamp NULL DEFAULT NULL COMMENT '记录修改时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '记录删除时间',
  `number` varchar(255) DEFAULT NULL COMMENT '订单编号',
  `idc_id` int(11) unsigned NOT NULL COMMENT '所属数据中心ID',
  `server_room_id` int(11) unsigned NOT NULL COMMENT '所属机房ID',
  `physical_area` varchar(255) DEFAULT NULL COMMENT '物理区域',
  `usage` varchar(255) DEFAULT NULL COMMENT '用途',
  `category` varchar(255) DEFAULT NULL COMMENT '设备类型',
  `amount` int(11) DEFAULT NULL COMMENT '数量',
  `left_amount` int(11) DEFAULT NULL COMMENT '未到货数量数量',
  `expected_arrival_date` date DEFAULT NULL COMMENT '预计到货时间',
  `pre_occupied_cabinets` json DEFAULT NULL COMMENT '预占用机架信息',
  `pre_occupied_usites` json DEFAULT NULL COMMENT '预占用机位信息',
  `status` enum('purchasing','partly_arrived','all_arrived','canceled','finished') DEFAULT NULL COMMENT '状态。采购中|部分到货|全部到货|已取消|已完成',
  `remark` varchar(1024) DEFAULT NULL COMMENT '备注',
  `creator` varchar(255) DEFAULT NULL COMMENT '订单创建者',
  `asset_belongs` varchar(255) DEFAULT NULL COMMENT '资产归属',
  `owner` varchar(255) DEFAULT NULL COMMENT '负责人',
  `is_rental` enum('yes','no') DEFAULT 'no' COMMENT '是否租赁:yes-是;no-否;',
  `maintenance_service_provider` varchar(255) DEFAULT NULL COMMENT '维保服务供应商',
  `maintenance_service` text COMMENT '维保服务详细内容',
  `logistics_service` text COMMENT '物流服务详细内容',
  `maintenance_service_date_begin` datetime DEFAULT NULL COMMENT '维保服务起始日期',
  `maintenance_service_date_end` datetime DEFAULT NULL COMMENT '维保服务截止日期',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_number` (`number`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='订单表';

-- 正在导出表  cloudboot_cyclone.order 的数据：~0 rows (大约)
/*!40000 ALTER TABLE `order` DISABLE KEYS */;
/*!40000 ALTER TABLE `order` ENABLE KEYS */;

-- 导出  表 cloudboot_cyclone.permission_code 结构
CREATE TABLE IF NOT EXISTS `permission_code` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL COMMENT '记录创建时间',
  `updated_at` timestamp NULL DEFAULT NULL COMMENT '记录修改时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '记录删除时间',
  `code` varchar(255) DEFAULT NULL COMMENT '权限码',
  `title` varchar(255) DEFAULT NULL COMMENT '权限码标题',
  `note` varchar(255) DEFAULT NULL COMMENT '权限码备注',
  `pid` int(11) unsigned NOT NULL COMMENT '父级ID',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `perm_data_dict_code` (`code`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='权限数据字典';

-- 正在导出表  cloudboot_cyclone.permission_code 的数据：~0 rows (大约)
/*!40000 ALTER TABLE `permission_code` DISABLE KEYS */;
/*!40000 ALTER TABLE `permission_code` ENABLE KEYS */;

-- 导出  表 cloudboot_cyclone.server_cabinet 结构
CREATE TABLE IF NOT EXISTS `server_cabinet` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL COMMENT '记录创建时间',
  `updated_at` timestamp NULL DEFAULT NULL COMMENT '记录修改时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '记录删除时间',
  `idc_id` int(11) unsigned NOT NULL COMMENT '所属数据中心ID',
  `server_room_id` int(11) unsigned NOT NULL COMMENT '所属机房ID',
  `network_area_id` int(11) unsigned NOT NULL COMMENT '所属网络区域ID',
  `number` varchar(255) DEFAULT NULL COMMENT '编号',
  `height` int(11) unsigned NOT NULL COMMENT '高度(单位U)',
  `type` enum('server','kvm_server','network_device','reserved') DEFAULT NULL COMMENT '类型:server-通用服务器;kvm_server-虚拟化服务器; network_device-网络设备; reserved-预留;',
  `network_rate` varchar(255) DEFAULT NULL COMMENT '网络速率',
  `current` varchar(255) DEFAULT NULL COMMENT '电流',
  `available_power` varchar(255) DEFAULT NULL COMMENT '可用功率',
  `max_power` varchar(255) DEFAULT NULL COMMENT '最大功率',
  `is_enabled` enum('yes','no') DEFAULT NULL COMMENT '是否启用',
  `enable_time` datetime DEFAULT NULL COMMENT '启用时间',
  `is_powered` enum('yes','no') DEFAULT NULL COMMENT '是否开电',
  `power_on_time` datetime DEFAULT NULL COMMENT '开电时间',
  `power_off_time` datetime DEFAULT NULL COMMENT '关电时间',
  `status` enum('under_construction','not_enabled','enabled','offline','locked') DEFAULT NULL COMMENT '状态。under_construction-建设中; not_enabled-未启用; enabled-已启用; offline-已下线;',
  `remark` varchar(1024) DEFAULT NULL COMMENT '备注',
  `creator` varchar(255) DEFAULT NULL COMMENT '记录创建者ID',
  `updater` varchar(255) DEFAULT NULL COMMENT '记录更新者ID',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_idc_room_cabinet` (`idc_id`,`server_room_id`,`number`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='机架(柜)信息表';

-- 正在导出表  cloudboot_cyclone.server_cabinet 的数据：~0 rows (大约)
/*!40000 ALTER TABLE `server_cabinet` DISABLE KEYS */;
/*!40000 ALTER TABLE `server_cabinet` ENABLE KEYS */;

-- 导出  表 cloudboot_cyclone.server_room 结构
CREATE TABLE IF NOT EXISTS `server_room` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL COMMENT '记录创建时间',
  `updated_at` timestamp NULL DEFAULT NULL COMMENT '记录修改时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '记录删除时间',
  `idc_id` int(11) unsigned NOT NULL COMMENT '所属数据中心ID',
  `name` varchar(255) DEFAULT NULL COMMENT '机房管理单元名称',
  `first_server_room` int(11) DEFAULT NULL COMMENT '一级机房ID',
  `city` varchar(255) DEFAULT NULL COMMENT '所属城市',
  `address` varchar(255) DEFAULT NULL COMMENT '机房管理单元地址',
  `server_room_manager` varchar(255) DEFAULT NULL COMMENT '机房管理单元负责人',
  `vendor_manager` varchar(255) DEFAULT NULL COMMENT '供应商负责人',
  `network_asset_manager` varchar(255) DEFAULT NULL COMMENT '网络资产负责人',
  `support_phone_number` varchar(255) DEFAULT NULL COMMENT '7*24小时保障电话',
  `status` enum('under_construction','accepted','production','abolished') DEFAULT NULL COMMENT '状态。under_construction-建设中; accepted-已验收; production-已投产; abolished-已裁撤;',
  `creator` varchar(255) DEFAULT NULL COMMENT '记录创建者ID',
  `updater` varchar(255) DEFAULT NULL COMMENT '记录更新者ID',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_idc_room` (`name`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='机房信息表';

-- 正在导出表  cloudboot_cyclone.server_room 的数据：~0 rows (大约)
/*!40000 ALTER TABLE `server_room` DISABLE KEYS */;
/*!40000 ALTER TABLE `server_room` ENABLE KEYS */;

-- 导出  表 cloudboot_cyclone.server_usite 结构
CREATE TABLE IF NOT EXISTS `server_usite` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL COMMENT '记录创建时间',
  `updated_at` timestamp NULL DEFAULT NULL COMMENT '记录修改时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '记录删除时间',
  `idc_id` int(11) unsigned NOT NULL COMMENT '所属数据中心ID',
  `server_room_id` int(11) unsigned NOT NULL COMMENT '所属机房ID',
  `server_cabinet_id` int(11) unsigned NOT NULL COMMENT '所属机架ID',
  `number` varchar(255) DEFAULT NULL COMMENT '编号',
  `beginning` int(11) unsigned NOT NULL COMMENT '起始U数',
  `height` int(11) DEFAULT NULL COMMENT '高度(单位U)',
  `physical_area` varchar(255) DEFAULT NULL COMMENT '物理区域',
  `oobnet_switches` json DEFAULT NULL COMMENT '管理网交换机对象数组(对象含name、port属性)',
  `intranet_switches` json DEFAULT NULL COMMENT '内网交换机对象数组(对象含name、port属性)',
  `extranet_switches` json DEFAULT NULL COMMENT '外网交换机对象数组(对象含name、port属性)',
  `la_wa_port_rate` enum('GE','10GE','25GE','40GE') DEFAULT NULL COMMENT '内外网端口速率',
  `status` enum('free','pre_occupied','used','disabled') DEFAULT NULL COMMENT '状态。free-空闲; pre_occupied-预占用; used-已使用; disabled-不可用;',
  `remark` varchar(1024) DEFAULT NULL COMMENT '备注',
  `creator` varchar(255) DEFAULT NULL COMMENT '记录创建者ID',
  `updater` varchar(255) DEFAULT NULL COMMENT '记录更新者ID',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_idc_room_cabinet_usite` (`idc_id`,`server_room_id`,`server_cabinet_id`,`number`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='机位(U位)信息表';

-- 正在导出表  cloudboot_cyclone.server_usite 的数据：~0 rows (大约)
/*!40000 ALTER TABLE `server_usite` DISABLE KEYS */;
/*!40000 ALTER TABLE `server_usite` ENABLE KEYS */;

-- 导出  表 cloudboot_cyclone.store_room 结构
CREATE TABLE IF NOT EXISTS `store_room` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL COMMENT '记录创建时间',
  `updated_at` timestamp NULL DEFAULT NULL COMMENT '记录修改时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '记录删除时间',
  `idc_id` int(11) unsigned NOT NULL COMMENT '所属数据中心ID',
  `name` varchar(255) DEFAULT NULL COMMENT '库房管理单元名称',
  `first_server_room` varchar(255) DEFAULT NULL COMMENT '一级机房ID',
  `city` varchar(255) DEFAULT NULL COMMENT '所属城市',
  `address` varchar(255) DEFAULT NULL COMMENT '库房管理单元地址',
  `store_room_manager` varchar(255) DEFAULT NULL COMMENT '库房管理单元负责人',
  `vendor_manager` varchar(255) DEFAULT NULL COMMENT '供应商负责人',
  `status` enum('under_construction','accepted','production','abolished') DEFAULT NULL COMMENT '状态。under_construction-建设中; accepted-已验收; production-已投产; abolished-已裁撤;',
  `creator` varchar(255) DEFAULT NULL COMMENT '记录创建者ID',
  `updater` varchar(255) DEFAULT NULL COMMENT '记录更新者ID',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_store_room_name` (`name`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='库房管理单元表';

-- 正在导出表  cloudboot_cyclone.store_room 的数据：~0 rows (大约)
/*!40000 ALTER TABLE `store_room` DISABLE KEYS */;
/*!40000 ALTER TABLE `store_room` ENABLE KEYS */;

-- 导出  表 cloudboot_cyclone.system_setting 结构
CREATE TABLE IF NOT EXISTS `system_setting` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `created_at` timestamp NULL DEFAULT NULL COMMENT '记录创建时间',
  `updated_at` timestamp NULL DEFAULT NULL COMMENT '记录修改时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '记录删除时间',
  `key` varchar(255) NOT NULL COMMENT '系统配置名',
  `value` mediumtext NOT NULL COMMENT '系统配置值',
  `desc` varchar(512) DEFAULT '' COMMENT '系统配置描述',
  `updater` varchar(255) DEFAULT NULL COMMENT '记录更新者ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_system_setting_key` (`key`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8 COMMENT='系统设置信息表';

-- 正在导出表  cloudboot_cyclone.system_setting 的数据：~5 rows (大约)
/*!40000 ALTER TABLE `system_setting` DISABLE KEYS */;
REPLACE INTO `system_setting` (`id`, `created_at`, `updated_at`, `deleted_at`, `key`, `value`, `desc`, `updater`) VALUES
	(1, '2018-11-26 03:37:20', '2018-11-26 03:37:20', NULL, 'installation_timeout', '3600', '操作系统安装超时时间(单位秒)', NULL),
	(2, '2021-09-01 15:00:00', '2021-09-01 15:00:00', NULL, 'route_conf', '[{"category_code":"add_idcs","category_name":"新增数据中心","http_method":"Post","method":"GetIDCByID","url":"/api/cloudboot/v1/idcs"},{"category_code":"mod_idcs","category_name":"修改数据中心","http_method":"Put","method":"GetIDCByID","url":"/api/cloudboot/v1/idcs/{id}"},{"category_code":"mod_idcs_status","category_name":"修改数据中心状态","http_method":"Put","method":"","url":"/api/cloudboot/v1/idcs/status"},{"category_code":"rm_idcs","category_name":"删除数据中心","http_method":"Delete","method":"GetIDCByID","url":"/api/cloudboot/v1/idcs/{id}"},{"category_code":"add_server-rooms","category_name":"新增机房","http_method":"Post","method":"GetServerRoomByID","url":"/api/cloudboot/v1/server-rooms"},{"category_code":"mod_server-rooms","category_name":"修改机房","http_method":"Put","method":"GetServerRoomByID","url":"/api/cloudboot/v1/server-rooms/{id}"},{"category_code":"mod_server-rooms_status","category_name":"修改机房状态","http_method":"Put","method":"","url":"/api/cloudboot/v1/server-rooms/status"},{"category_code":"rm_server-rooms","category_name":"删除机房","http_method":"Delete","method":"GetServerRoomByID","url":"/api/cloudboot/v1/server-rooms/{id}"},{"category_code":"add_server-rooms","category_name":"加载导入机房文件","http_method":"Post","method":"GetServerRoomByID","url":"/api/cloudboot/v1/server-rooms/upload"},{"category_code":"add_server-rooms","category_name":"导入机房文件预览","http_method":"Post","url":"/api/cloudboot/v1/server-rooms/imports/previews"},{"category_code":"import_server-rooms","category_name":"导入机房文件","http_method":"Post","method":"GetServerRoomByID","url":"/api/cloudboot/v1/server-rooms/imports"},{"category_code":"add_network-areas","category_name":"新增网络区域","http_method":"Post","method":"GetNetworkAreaByID","url":"/api/cloudboot/v1/network-areas"},{"category_code":"mod_network-areas","category_name":"修改网络区域","http_method":"Put","method":"GetNetworkAreaByID","url":"/api/cloudboot/v1/network-areas/{id}"},{"category_code":"mod_network-areas_status","category_name":"修改网络区域状态","http_method":"Put","method":"","url":"/api/cloudboot/v1/network-areas/status"},{"category_code":"rm_network-areas","category_name":"删除网络区域","http_method":"Delete","method":"GetNetworkAreaByID","url":"/api/cloudboot/v1/network-areas/{id}"},{"category_code":"add_network-areas","category_name":"加载导入导入网络区域文件","http_method":"Post","method":"GetNetworkAreaByID","url":"/api/cloudboot/v1/network-areas/upload"},{"category_code":"add_network-areas","category_name":"导入网络区域文件预览","http_method":"Post","url":"/api/cloudboot/v1/network-areas/imports/previews"},{"category_code":"add_network-areas","category_name":"导入网络区域文件","http_method":"Post","url":"/api/cloudboot/v1/network-areas/imports"},{"category_code":"add_server-cabinets","category_name":"新增机架(柜)","http_method":"Post","method":"GetServerCabinetByID","url":"/api/cloudboot/v1/server-cabinets"},{"category_code":"mod_server-cabinets","category_name":"修改机架(柜)","http_method":"Put","method":"GetServerCabinetByID","url":"/api/cloudboot/v1/server-cabinets/{id}"},{"category_code":"add_server-cabinets_status","category_name":"修改机架(柜)状态","http_method":"Put","method":"","url":"/api/cloudboot/v1/server-cabinets/status"},{"category_code":"rm_server-cabinets","category_name":"删除机架(柜)","http_method":"Delete","method":"GetServerCabinetByID","url":"/api/cloudboot/v1/server-cabinets/{id}"},{"category_code":"power_on_server-cabinets","category_name":"机架开电","http_method":"Post","method":"GetServerCabinetByID","url":"/api/cloudboot/v1/server-cabinets/{id}/power"},{"category_code":"power_off_server-cabinets","category_name":"机架关电","http_method":"Delete","method":"GetServerCabinetByID","url":"/api/cloudboot/v1/server-cabinets/{id}/power"},{"category_code":"add_server-cabinets","category_name":"加载导入机架(柜)文件","http_method":"Post","method":"GetServerCabinetByID","url":"/api/cloudboot/v1/server-cabinets/upload"},{"category_code":"add_server-cabinets","category_name":"导入机架(柜)文件预览","http_method":"Post","url":"/api/cloudboot/v1/server-cabinets/imports/previews"},{"category_code":"add_server-cabinets","category_name":"导入机架(柜)文件","http_method":"Post","url":"/api/cloudboot/v1/server-cabinets/imports"},{"category_code":"add_server-usites","category_name":"新增机位(U位)","http_method":"Post","method":"GetServerUSiteByID","url":"/api/cloudboot/v1/server-usites"},{"category_code":"mod_server-usites","category_name":"修改机位(U位)","http_method":"Put","method":"GetServerUSiteByID","url":"/api/cloudboot/v1/server-usites/{id}"},{"category_code":"rm_server-usites","category_name":"删除机位(U位)","http_method":"Delete","method":"GetServerUSiteByID","url":"/api/cloudboot/v1/server-usites/{id}"},{"category_code":"rm_server-usites_ports","category_name":"删除机位(U位)端口号","http_method":"Delete","method":"GetServerUSiteByID","url":"/api/cloudboot/v1/server-usites/{id}/ports"},{"category_code":"mod_server-usites","category_name":"修改机位(U位)状态","http_method":"Put","method":"GetServerUSiteByID","url":"/api/cloudboot/v1/server-usites/status"},{"category_code":"add_server-usites","category_name":"加载导入机位(U位)文件","http_method":"Post","url":"/api/cloudboot/v1/server-usites/upload"},{"category_code":"add_server-usites","category_name":"导入机位(U位)文件预览","http_method":"Post","url":"/api/cloudboot/v1/server-usites/imports/previews"},{"category_code":"add_server-usites","category_name":"导入机位(U位)文件","http_method":"Post","url":"/api/cloudboot/v1/server-usites/imports"},{"category_code":"add_server-usites","category_name":"加载导入机位(U位)端口号文件","http_method":"Post","url":"/api/cloudboot/v1/server-usites/ports/upload"},{"category_code":"add_server-usites","category_name":"导入机位(U位)端口号文件预览","http_method":"Post","url":"/api/cloudboot/v1/server-usites/ports/imports/previews"},{"category_code":"add_server-usites","category_name":"导入机位(U位)端口号文件","http_method":"Post","url":"/api/cloudboot/v1/server-usites/ports/imports"},{"category_code":"add_devices","category_name":"新增物理机","http_method":"Post","url":"/api/cloudboot/v1/devices/{sn}/collections"},{"category_code":"add_devices","category_name":"加载导入物理机文件","http_method":"Post","url":"/api/cloudboot/v1/devices/upload"},{"category_code":"add_devices","category_name":"导入物理机文件预览","http_method":"Post","url":"/api/cloudboot/v1/devices/imports/previews"},{"category_code":"add_devices","category_name":"加入物理机文件","http_method":"Post","url":"/api/cloudboot/v1/devices/imports"},{"category_code":"mod_device","category_name":"修改物理机","http_method":"Put","url":"/api/cloudboot/v1/device"},{"category_code":"mod_device","category_name":"批量修改物理机状态/用途","http_method":"Put","url":"/api/cloudboot/v1/devices"},{"category_code":"rm_devices","category_name":"删除物理机","http_method":"Delete","url":"/api/cloudboot/v1/devices"},{"category_code":"add_ip-networks","category_name":"新增网段","http_method":"Post","url":"/api/cloudboot/v1/ip-networks"},{"category_code":"add_ip-networks","category_name":"修改网段","http_method":"Put","url":"/api/cloudboot/v1/ip-networks/{id}"},{"category_code":"rm_ip-networks","category_name":"删除网段","http_method":"Delete","url":"/api/cloudboot/v1/ip-networks/{id}"},{"category_code":"sync_ip-networks","category_name":"网段同步","http_method":"Post","method":"","url":"/api/cloudboot/v1/ip-networks/sync"},{"category_code":"mod_ips","category_name":"手动分配IP","http_method":"Put","url":"/api/cloudboot/v1/ips/assigns"},{"category_code":"mod_ips","category_name":"手动取消IP分配","http_method":"Put","url":"/api/cloudboot/v1/ips/unassigns"},{"category_code":"rm_network-device","category_name":"删除网络设备","http_method":"Delete","method":"GetNetworkDeviceByID","url":"/api/cloudboot/v1/network/devices/{id}"},{"category_code":"add_network-device","category_name":"新增网络设备","http_method":"Post","method":"GetNetworkDeviceByID","url":"/api/cloudboot/v1/network/devices"},{"category_code":"sync_network-device","category_name":"网段设备同步","http_method":"Post","method":"","url":"/api/cloudboot/v1/network/devices/sync"},{"category_code":"add_devices","category_name":"上报安装进度","http_method":"Post","url":"/api/cloudboot/v1/devices/{sn}/installations/progress"},{"category_code":"add_devices","category_name":"新增装机参数","http_method":"Post","url":"/api/cloudboot/v1/devices/settings"},{"category_code":"mod_devices","category_name":"批量重装设备","http_method":"Put","url":"/api/cloudboot/v1/devices/installations/reinstalls"},{"category_code":"mod_devices","category_name":"设备重装","http_method":"Post","url":"/api/cloudboot/v1/devices/installations/os-reinstallations"},{"category_code":"mod_devices","category_name":"批量取消安装设备","http_method":"Put","url":"/api/cloudboot/v1/devices/installations/cancels"},{"category_code":"rm_devices","category_name":"批量删除设备装机参数","http_method":"Delete","url":"/api/cloudboot/v1/devices/settings"},{"category_code":"rm_hardware-templates","category_name":"删除硬件模板","http_method":"Delete","url":"/api/cloudboot/v1/hardware-templates/{id}"},{"category_code":"mod_hardware-templates","category_name":"修改硬件模板","http_method":"Put","url":"/api/cloudboot/v1/hardware-templates/{id}"},{"category_code":"add_hardware-templates","category_name":"新增硬件模板","http_method":"Post","url":"/api/cloudboot/v1/hardware-templates"},{"category_code":"add_image-templates","category_name":"新增镜像模板","http_method":"Post","url":"/api/cloudboot/v1/image-templates"},{"category_code":"mod_image-templates","category_name":"修改镜像模板","http_method":"Put","url":"/api/cloudboot/v1/image-templates/{id}"},{"category_code":"rm_image-templates","category_name":"删除镜像模板","http_method":"Delete","url":"/api/cloudboot/v1/image-templates/{id}"},{"category_code":"add_system-templates","category_name":"新增系统模板","http_method":"Post","url":"/api/cloudboot/v1/system-templates"},{"category_code":"mod_system-templates","category_name":"修改系统模板","http_method":"Put","url":"/api/cloudboot/v1/system-templates/{id}"},{"category_code":"rm_system-templates","category_name":"删除系统模板","http_method":"Delete","url":"/api/cloudboot/v1/system-templates/{id}"},{"category_code":"add_devices","category_name":"新增带外管理","http_method":"Post","url":"/api/cloudboot/v1/devices/power"},{"category_code":"add_devices","category_name":"带外管理批量PXE重启","http_method":"Put","url":"/api/cloudboot/v1/devices/power/pxe/restart"},{"category_code":"mod_devices","category_name":"带外管理批量重启","http_method":"Put","url":"/api/cloudboot/v1/devices/power/restart"},{"category_code":"rm_devices","category_name":"带外管理批量关机","http_method":"Delete","url":"/api/cloudboot/v1/devices/power"},{"category_code":"mod_devices","category_name":"修改带外用户密码","http_method":"Put","req_param_desensitization":"oob_password_old,oob_password_new","url":"/api/cloudboot/v1/devices/{sn}/oob/password"},{"category_code":"add_jobs","category_name":"新增巡检任务","http_method":"Post","url":"/api/cloudboot/v1/jobs/inspections"},{"category_code":"mod_user_password","category_name":"修改用户密码","http_method":"Put","req_param_desensitization":"old_password,new_password","url":"/api/cloudboot/v1/users/password"},{"category_code":"mod_gen_pxe","category_name":"为目标设备生成PXE文件","http_method":"Post","url":"/api/cloudboot/v1/devices/{sn}/pxe"},{"category_code":"mod_centos6_uefi_pxe","category_name":"生成centos6pxe模板","http_method":"Post","url":"/api/cloudboot/v1/devices/{sn}/centos6/uefi/pxe"},{"category_code":"submit_cabinet_offline_approval","category_name":"发起机架下线审批","http_method":"Post","url":"/api/cloudboot/v1/approvals/server-cabinets/offlines"},{"category_code":"submit_cabinet_poweroff_approval","category_name":"发起机架关电审批","http_method":"Post","url":"/api/cloudboot/v1/approvals/server-cabinets/poweroffs"},{"category_code":"submit_device_migrate_approval","category_name":"发起物理机迁移审批","http_method":"Post","url":"/api/cloudboot/v1/approvals/devices/migrations"},{"category_code":"submit_device_retire_approval","category_name":"发起物理机退役审批","http_method":"Post","url":"/api/cloudboot/v1/approvals/devices/retirements"},{"category_code":"submit_device_reinstall_approval","category_name":"发起物理机重装审批","http_method":"Post","url":"/api/cloudboot/v1/approvals/devices/os-reinstallations"},{"category_code":"return_limiter_token","category_name":"归还设备持有的DHCPIP令牌","http_method":"Delete","url":"/api/cloudboot/v1/devices/{sn}/limiters/tokens"},{"http_method":"Post","url":"/api/cloudboot/v1/devices/store/upload","category_code":"upload","category_name":"上传设备导入到库房文件"},{"http_method":"Post","url":"/api/cloudboot/v1/devices/store/imports/previews","category_code":"add_devices","category_name":"物理机导入到库房预览"},{"http_method":"Post","url":"/api/cloudboot/v1/devices/store/imports","category_code":"add_devices","category_name":"物理机导入到库房"},{"http_method":"Post","url":"/api/cloudboot/v1/devices/stock/upload","category_code":"add_devices","category_name":"上传存量设备文件"},{"http_method":"Post","url":"/api/cloudboot/v1/devices/stock/imports/previews","category_code":"add_devices","category_name":"存量设备导入预览"},{"http_method":"Post","url":"/api/cloudboot/v1/devices/stock/imports","category_code":"add_devices","category_name":"存量设备导入"},{"http_method":"Put","url":"/api/cloudboot/v1/device/operation/status","category_code":"mod_device","category_name":"修改设备运营状态"},{"http_method":"Post","url":"/api/cloudboot/v1/approvals/devices/migrations/upload","category_code":"add_approvals","category_name":"上传物理机搬迁文件"},{"http_method":"Post","url":"/api/cloudboot/v1/approvals/devices/migrations/imports/previews","category_code":"add_approvals","category_name":"物理机搬迁导入预览"},{"http_method":"Post","url":"/api/cloudboot/v1/approvals/devices/migrations/imports","category_code":"add_approvals","category_name":"物理机搬迁导入"},{"http_method":"Post","url":"/api/cloudboot/v1/approvals/idc/abolish","category_code":"add_approvals","category_name":"数据中心裁撤审批"},{"http_method":"Post","url":"/api/cloudboot/v1/approvals/server-room/abolish","category_code":"add_approvals","category_name":"机房裁撤审批"},{"http_method":"Post","url":"/api/cloudboot/v1/approvals/network-area/offline","category_code":"add_approvals","category_name":"网络区域下线审批"},{"http_method":"Post","url":"/api/cloudboot/v1/approvals/ip/unassign","category_code":"add_approvals","category_name":"IP回收审批"},{"http_method":"Post","url":"/api/cloudboot/v1/approvals/device/poweroff","category_code":"add_approvals","category_name":"物理机关电审批"},{"http_method":"Post","url":"/api/cloudboot/v1/approvals/device/restart","category_code":"add_approvals","category_name":"物理机重启审批"},{"http_method":"Post","url":"/api/cloudboot/v1/approvals/devices/recycle","category_code":"add_approvals","category_name":"物理机回收审批"},{"http_method":"Post","url":"/api/cloudboot/v1/store-room","category_code":"add_store-room","category_name":"新增库房"},{"http_method":"Delete","url":"/api/cloudboot/v1/store-room/{id}","category_code":"rm_store-room","category_name":"删除库房"},{"http_method":"Put","url":"/api/cloudboot/v1/store-room","category_code":"mod_store-room","category_name":"修改库房"},{"http_method":"Post","url":"/api/cloudboot/v1/store-room/upload","category_code":"add_store-room","category_name":"上传库房导入文件"},{"http_method":"Post","url":"/api/cloudboot/v1/store-room/imports/previews","category_code":"add_store-room","category_name":"库房导入预览"},{"http_method":"Post","url":"/api/cloudboot/v1/store-room/imports","category_code":"add_store-room","category_name":"库房导入"},{"http_method":"Post","url":"/api/cloudboot/v1/virtual-cabinet","category_code":"add_virtual-cabinet","category_name":"新增虚拟货架"},{"http_method":"Delete","url":"/api/cloudboot/v1/virtual-cabinet/{id}","category_code":"rm_virtual-cabinet","category_name":"删除虚拟货架"},{"http_method":"Post","url":"/api/cloudboot/v1/network/devices/upload","category_code":"add_network","category_name":"上传网络设备导入文件"},{"http_method":"Post","url":"/api/cloudboot/v1/network/devices/imports/previews","category_code":"add_network","category_name":"网络设备导入预览"},{"http_method":"Post","url":"/api/cloudboot/v1/network/devices/imports","category_code":"add_network","category_name":"网络设备导入"},{"http_method":"Post","url":"/api/cloudboot/v1/ip-networks/upload","category_code":"add_ip-networks","category_name":"上传网段导入文件"},{"http_method":"Post","url":"/api/cloudboot/v1/ip-networks/imports/previews","category_code":"add_ip-networks","category_name":"网段导入预览"},{"http_method":"Post","url":"/api/cloudboot/v1/ip-networks/imports","category_code":"add_ip-networks","category_name":"网段导入"},{"http_method":"Put","url":"/api/cloudboot/v1/ips/status/disable","category_code":"mod_ips","category_name":"禁用IP"},{"http_method":"Post","url":"/api/cloudboot/v1/special-device","category_code":"add_special-device","category_name":"新增特殊设备"},{"http_method":"Post","url":"/api/cloudboot/v1/special-devices/upload","category_code":"add_special-devices","category_name":"上传特殊设备导入文件"},{"http_method":"Post","url":"/api/cloudboot/v1/special-devices/imports/previews","category_code":"add_special-devices","category_name":"特殊设备导入预览"},{"http_method":"Post","url":"/api/cloudboot/v1/special-devices/imports","category_code":"add_special-devices","category_name":"特殊设备导入"},{"http_method":"Put","url":"/api/cloudboot/v1/devices/oob/re-access","category_code":"mod_devices","category_name":"重新纳管带外"},{"http_method":"Delete","url":"/api/cloudboot/v1/devices/limiters/tokens","category_code":"rm_devices","category_name":"一键释放token"},{"http_method":"Post","url":"/api/cloudboot/v1/devices/{sn}/components/hw-server/logs","category_code":"add_devices","category_name":"保存hw-server组件日志"},{"http_method":"Post","url":"/api/cloudboot/v1/devices/{sn}/components/peconfig/logs","category_code":"add_devices","category_name":"保存peconfig组件日志"},{"http_method":"Post","url":"/api/cloudboot/v1/devices/{sn}/components/cloudboot-agent/logs","category_code":"add_devices","category_name":"保存cloudboot-agent组件日志"},{"http_method":"Post","url":"/api/cloudboot/v1/devices/{sn}/components/winconfig/logs","category_code":"add_devices","category_name":"保存winconfig组件日志"},{"http_method":"Post","url":"/api/cloudboot/v1/devices/{sn}/components/network-config/logs","category_code":"add_devices","category_name":"保存系统配置日志"},{"http_method":"Post","url":"/api/cloudboot/v1/devices/{sn}/components/image-clone/logs","category_code":"add_devices","category_name":"保存镜像制作日志"},{"http_method":"Post","url":"/api/cloudboot/v1/order","category_code":"add_order","category_name":"新增订单"},{"http_method":"Delete","url":"/api/cloudboot/v1/orders","category_code":"rm_orders","category_name":"删除订单"},{"http_method":"Put","url":"/api/cloudboot/v1/order","category_code":"mod_order","category_name":"修改订单"},{"http_method":"Put","url":"/api/cloudboot/v1/order/status","category_code":"mod_order","category_name":"修改订单状态"},{"http_method":"Post","url":"/api/cloudboot/v1/device-category","category_code":"add_device-category","category_name":"新增设备类型"},{"http_method":"Delete","url":"/api/cloudboot/v1/device-categories","category_code":"rm_device-categories","category_name":"删除设备类型"},{"http_method":"Put","url":"/api/cloudboot/v1/device-category","category_code":"mod_device-category","category_name":"修改设备类型"},{"http_method":"Delete","url":"/api/cloudboot/v1/jobs/{job_id}","category_code":"rm_jobs","category_name":"删除任务"},{"http_method":"Put","url":"/api/cloudboot/v1/jobs/{job_id}/pausing","category_code":"mod_jobs","category_name":"暂停任务"},{"http_method":"Put","url":"/api/cloudboot/v1/jobs/{job_id}/unpausing","category_code":"mod_jobs","category_name":"继续已暂停任务"},{"http_method":"Post","method":"","url":"/api/cloudboot/v1/device-setting-rules","category_code":"add_device-setting-rules","category_name":"新增装机参数规则"},{"http_method":"Delete","method":"","url":"/api/cloudboot/v1/device-setting-rules","category_code":"rm_device-setting-rules","category_name":"删除装机参数规则"},{"http_method":"Put","method":"","url":"/api/cloudboot/v1/device-setting-rules","category_code":"mod_device-setting-rules","category_name":"修改装机参数规则"},{"category_code":"mod_devices","category_name":"自动重装设备","http_method":"Put","method":"","url":"/api/cloudboot/v1/devices/installations/autoreinstalls"}]', '操作记录中间件配置', NULL),
	(3, '2021-09-01 15:00:00', '2021-09-01 15:00:00', NULL, 'cloudboot_menu_permission', '{"content":{"children":[{"children":[],"id":"menu_home","note":"菜单-概览","title":"概览"},{"children":[{"children":[],"id":"menu_idc","note":"菜单-数据中心管理-数据中心","title":"数据中心"},{"children":[],"id":"menu_server_room","note":"菜单-数据中心管理-机房","title":"机房"},{"children":[],"id":"menu_network_area","note":"菜单-数据中心管理-网络区域","title":"网络区域"},{"children":[],"id":"menu_server_cabinet","note":"菜单-数据中心管理-机架","title":"机架"},{"children":[],"id":"menu_server_usite","note":"菜单-数据中心管理-机位","title":"机位"},{"children":[],"id":"menu_store_room","note":"菜单-数据中心管理-库房","title":"库房"}],"id":"menu_idc_management","note":"菜单-数据中心管理","title":"数据中心管理"},{"children":[{"children":[],"id":"menu_network_device","note":"菜单-网络管理-网络设备","title":"网络设备"},{"children":[],"id":"menu_ip_network","note":"菜单-网络管理-IP网段","title":"IP网段"},{"children":[],"id":"menu_ip","note":"菜单-网络管理-IP分配","title":"IP分配"}],"id":"menu_network_management","note":"菜单-网络管理","title":"网络管理"},{"children":[{"children":[],"id":"menu_physical_machine","note":"菜单-物理机管理-物理机列表","title":"物理机列表"},{"children":[],"id":"menu_special_device","note":"菜单-物理机管理-特殊设备列表","title":"特殊设备列表"},{"children":[],"id":"menu_oob_info","note":"菜单-物理机管理-带外信息","title":"带外信息"},{"children":[],"id":"menu_predeploy_physical_machine","note":"菜单-物理机管理-待部署物理机","title":"待部署物理机"},{"children":[],"id":"menu_device_setting","note":"菜单-物理机管理-部署列表","title":"部署列表"},{"children":[],"id":"menu_device_setting_rule","note":"菜单-物理机管理-装机参数规则列表","title":"装机参数规则列表"},{"children":[],"id":"menu_inspection","note":"菜单-物理机管理-硬件巡检","title":"硬件巡检"}],"id":"menu_physical_machine_management","note":"菜单-物理机管理","title":"物理机管理"},{"children":[{"children":[],"id":"menu_order","note":"菜单-订单管理-订单列表","title":"订单列表"},{"children":[],"id":"menu_device_category","note":"菜单-订单管理-设备类型列表","title":"设备类型列表"}],"id":"menu_order_management","note":"菜单-订单管理","title":"订单管理"},{"children":[{"children":[],"id":"menu_system_template","note":"菜单-配置管理-装机配置","title":"装机配置"},{"children":[],"id":"menu_hardware_template","note":"菜单-配置管理-硬件配置","title":"硬件配置"}],"id":"menu_template_management","note":"菜单-配置管理","title":"配置管理"},{"children":[],"id":"menu_approval","note":"菜单-审批管理","title":"审批管理"},{"children":[],"id":"menu_task_management","note":"菜单-任务管理","title":"任务管理"},{"children":[{"children":[],"id":"menu_audit_api","note":"菜单-操作审计-接口调用记录","title":"接口调用记录"},{"children":[],"id":"menu_audit_log","note":"菜单-操作审计-操作记录","title":"操作记录"}],"id":"menu_audit","note":"菜单-操作审计","title":"操作审计"},{"children":[],"id":"menu_user_management","note":"菜单-用户管理","title":"用户管理"}],"id":"","note":"","title":"所有权限"},"message":"操作成功","status":"success"}', '菜单权限列表', NULL),
	(4, '2021-09-01 15:00:00', '2021-09-01 15:00:00', NULL, 'cloudboot_button_permission', '{"content": {"children": [{"children": [{"children": [{"children": [],"id": "button_idc_create","note": "按钮-数据中心管理-新增","title": "新增数据中心"},{"children": [],"id": "button_idc_update","note": "按钮-数据中心管理-修改","title": "修改数据中心"},{"children": [],"id": "button_idc_delete","note": "按钮-数据中心管理-删除","title": "删除数据中心"},{"children": [],"id": "button_idc_accepted","note": "按钮-数据中心管理-验收","title": "验收数据中心"},{"children": [],"id": "button_idc_production","note": "按钮-数据中心管理-投产","title": "投产数据中心"},{"children": [],"id": "button_idc_abolished","note": "按钮-数据中心管理-裁撤","title": "裁撤数据中心"}],"id": "menu_idc","note": "菜单-数据中心管理-数据中心","title": "数据中心"},{"children": [{"children": [],"id": "button_server_room_create","note": "按钮-机房信息管理-新增","title": "新增机房信息"},{"children": [],"id": "button_server_room_update","note": "按钮-机房信息管理-修改","title": "修改机房信息"},{"children": [],"id": "button_server_room_delete","note": "按钮-机房信息管理-删除","title": "删除机房信息"},{"children": [],"id": "button_server_room_accepted","note": "按钮-机房信息管理-验收","title": "验收机房信息"},{"children": [],"id": "button_server_room_production","note": "按钮-机房信息管理-投产","title": "投产机房信息"},{"children": [],"id": "button_server_room_abolished","note": "按钮-机房信息管理-裁撤","title": "裁撤机房信息"},{"children": [],"id": "button_server_room_import","note": "按钮-机房信息管理-导入","title": "导入机房信息"},{"children": [],"id": "button_server_room_download","note": "按钮-机房信息管理-下载模板","title": "下载机房信息模板"}],"id": "menu_server_room","note": "菜单-数据中心管理-机房","title": "机房"},{"children": [{"children": [],"id": "button_network_area_create","note": "按钮-网络区域-新增","title": "新增网络区域"},{"children": [],"id": "button_network_area_update","note": "按钮-网络区域-修改","title": "修改网络区域"},{"children": [],"id": "button_network_area_delete","note": "按钮-网络区域-删除","title": "删除网络区域"},{"children": [],"id": "button_network_area_production","note": "按钮-网络区域-投产","title": "投产网络区域"},{"children": [],"id": "button_network_area_offline","note": "按钮-网络区域-下线","title": "下线网络区域"},{"children": [],"id": "button_network_area_import","note": "按钮-网络区域-导入","title": "导入网络区域"},{"children": [],"id": "button_network_area_download","note": "按钮-网络区域-下载模板","title": "下载网络区域模板"}],"id": "menu_network_area","note": "菜单-数据中心管理-网络区域","title": "网络区域"},{"children": [{"children": [],"id": "button_server_cabinet_create","note": "按钮-机架-新增","title": "新增机架"},{"children": [],"id": "button_server_cabinet_update","note": "按钮-机架-修改","title": "修改机架"},{"children": [],"id": "button_server_cabinet_delete","note": "按钮-机架-删除","title": "删除机架"},{"children": [],"id": "button_server_cabinet_enabled","note": "按钮-机架-启用","title": "启用机架"},{"children": [],"id": "button_server_cabinet_powerOn","note": "按钮-机架-开电","title": "开电机架"},{"children": [],"id": "button_server_cabinet_powerOff","note": "按钮-机架-关电","title": "关电机架"},{"children": [],"id": "button_server_cabinet_offline","note": "按钮-机架-下线","title": "下线机架"},{"children": [],"id": "button_server_cabinet_locked","note": "按钮-机架-锁定","title": "锁定机架"},{"children": [],"id": "button_server_cabinet_type","note": "按钮-机架-更新机架类型","title": "更新机架类型"},{"children": [],"id": "button_server_cabinet_remark","note": "按钮-机架-备注","title": "更新机架备注"},{"children": [],"id": "button_server_cabinet_import","note": "按钮-机架-导入","title": "导入机架"},{"children": [],"id": "button_server_cabinet_download","note": "按钮-机架-下载模板","title": "下载机架模板"}],"id": "menu_server_cabinet","note": "菜单-数据中心管理-机架","title": "机架"},{"children": [{"children": [],"id": "button_server_usite_create","note": "按钮-机位-新增","title": "新增机位"},{"children": [],"id": "button_server_usite_update","note": "按钮-机位-修改","title": "修改机位"},{"children": [],"id": "button_server_usite_delete","note": "按钮-机位-删除","title": "删除机位"},{"children": [],"id": "button_server_usite_delete_port","note": "按钮-机位-删除端口","title": "删除机位端口"},{"children": [],"id": "button_server_usite_status","note": "按钮-机位-更新机位状态","title": "更新机位状态"},{"children": [],"id": "button_server_usite_remark","note": "按钮-机位-备注","title": "更新机位备注"},{"children": [],"id": "button_server_usite_import","note": "按钮-机位-导入","title": "导入机位"},{"children": [],"id": "button_server_usite_download","note": "按钮-机位-下载模板","title": "下载机位模板"},{"children": [],"id": "button_server_usite_import_port","note": "按钮-机位端口-导入","title": "导入机位端口"},{"children": [],"id": "button_server_usite_download_port","note": "按钮-机位端口-下载模板","title": "下载机位端口模板"}],"id": "menu_server_usite","note": "菜单-数据中心管理-机位","title": "机位"},{"children": [{"children": [],"id": "button_store_room_create","note": "按钮-库房-新增","title": "新增库房管理单元"},{"children": [],"id": "button_store_room_import","note": "按钮-库房-导入","title": "导入库房管理单元"},{"children": [],"id": "button_store_room_delete","note": "按钮-库房-删除","title": "删除库房管理单元"},{"children": [],"id": "button_store_room_update","note": "按钮-库房-修改","title": "修改库房管理单元"},{"children": [],"id": "button_virtual_cabinet_create","note": "按钮-虚拟货架-新增","title": "新增虚拟货架"},{"children": [],"id": "button_virtual_cabinet_delete","note": "按钮-虚拟货架-删除","title": "删除虚拟货架"}],"id": "menu_store_room","note": "菜单-数据中心管理-库房","title": "库房信息管理"}],"id": "menu_idc_management","note": "菜单-数据中心管理","title": "数据中心管理"},{"children": [{"children": [{"children": [],"id": "button_network_device_create","note": "按钮-网络设备-新增","title": "新增网络设备"},{"children": [],"id": "button_network_device_delete","note": "按钮-网络设备-删除","title": "删除网络设备"},{"children": [],"id": "button_network_device_sync","note": "按钮-网络设备-同步","title": "同步网络设备"},{"children": [],"id": "button_network_device_import_download","note": "按钮-网络设备-下载导入模板","title": "下载导入网络设备模板"},{"children": [],"id": "button_network_device_import","note": "按钮-网络设备-导入","title": "导入网络设备"}],"id": "menu_network_device","note": "菜单-网络管理-网络设备","title": "网络设备"},{"children": [{"children": [],"id": "button_ip_network_create","note": "按钮-IP网段管理-新增","title": "新增IP网段管理"},{"children": [],"id": "button_ip_network_update","note": "按钮-IP网段管理-修改","title": "修改IP网段管理"},{"children": [],"id": "button_ip_network_delete","note": "按钮-IP网段管理-删除","title": "删除IP网段管理"},{"children": [],"id": "button_ip_network_sync","note": "按钮-IP网段管理-同步","title": "同步IP网段"},{"children": [],"id": "button_ip_network_import_download","note": "按钮-IP网段管理-导入模板下载","title": "下载导入IP网段模板"},{"children": [],"id": "button_ip_network_import","note": "按钮-IP网段管理-导入","title": "导入IP网段"}],"id": "menu_ip_network","note": "菜单-网络管理-IP网段","title": "IP网段"},{"children": [{"children": [],"id": "button_ip_assign","note": "按钮-IP分配-分配","title": "IP分配"},{"children": [],"id": "button_ip_unassign","note": "按钮-IP分配-取消分配","title": "取消IP分配"}],"id": "menu_ip","note": "菜单-网络管理-IP分配","title": "IP分配"}],"id": "menu_network_management","note": "菜单-网络管理","title": "网络管理"},{"children": [{"children": [{"children": [],"id": "button_physical_machine_update_status","note": "按钮-物理机列表-修改","title": "批量修改物理机状态"},{"children": [],"id": "button_physical_machine_update_usage","note": "按钮-物理机列表-修改","title": "批量修改物理机用途"},{"children": [],"id": "button_physical_machine_update","note": "按钮-物理机列表-修改","title": "修改物理机"},{"children": [],"id": "button_physical_machine_import","note": "按钮-物理机列表-导入","title": "存量物理机导入"},{"children": [],"id": "button_physical_machine_download","note": "按钮-物理机列表-下载模板","title": "存量物理机模板下载"},{"children": [],"id": "button_physical_machine_export","note": "按钮-物理机列表-导出","title": "物理机导出"},{"children": [],"id": "button_device_store_import","note": "按钮-物理机列表-导入","title": "物理机导入到库房"},{"children": [],"id": "button_device_store_import_download","note": "按钮-物理机列表-下载模板","title": "物理机导入到库房模板下载"}],"id": "menu_physical_machine","note": "菜单-物理机管理-物理机列表","title": "物理机列表"},{"children": [{"children": [],"id": "button_physical_machine_update_oob","note": "按钮-带外管理-修改带外","title": "修改物理机带外"},{"children": [],"id": "button_physical_machine_powerOn","note": "按钮-带外管理-开电","title": "开电物理机"},{"children": [],"id": "button_physical_machine_networkBoot","note": "按钮-带外管理-从网卡启动","title": "从网卡启动物理机"},{"children": [],"id": "button_oob_re_access","note": "按钮-带外管理-重新纳管带外","title": "重新纳管带外"},{"children": [],"id": "button_physical_oob_export","note": "按钮-带外管理-导出带外","title": "导出带外"}],"id": "menu_oob_info","note": "菜单-物理机管理-带外管理","title": "带外管理"},{"children": [{"children": [],"id": "button_special_device_create","note": "按钮-特殊设备列表-新增","title": "新增特殊设备"},{"children": [],"id": "button_special_device_delete","note": "按钮-特殊设备列表-新增","title": "新增特殊设备"},{"children": [],"id": "button_special_device_import_download","note": "按钮-特殊设备列表-下载模板","title": "导入特殊设备模板下载"},{"children": [],"id": "button_special_device_import","note": "按钮-特殊设备列表-导入","title": "导入特殊设备"}],"id": "menu_special_device","note": "菜单-物理机管理-特殊设备列表","title": "特殊设备列表"},{"children": [{"children": [],"id": "button_predeploy_physical_machine_osInstall","note": "按钮-待部署物理机-申请上架部署","title": "待部署物理机申请上架部署"},{"children": [],"id": "button_predeploy_physical_machine_powerOn","note": "按钮-待部署物理机-开电","title": "待部署物理机开电"},{"children": [],"id": "button_predeploy_physical_machine_powerOff","note": "按钮-待部署物理机-关电","title": "待部署物理机关电"},{"children": [],"id": "button_predeploy_physical_machine_reBoot","note": "按钮-待部署物理机-重启","title": "待部署物理机重启"},{"children": [],"id": "button_predeploy_physical_machine_networkBoot","note": "按钮-待部署物理机-从网卡启动","title": "待部署物理机从网卡启动"},{"children": [],"id": "button_predeploy_physical_machine_download","note": "按钮-待部署物理机-下载","title": "下载待部署物理机模板"},{"children": [],"id": "button_predeploy_physical_machine_import","note": "按钮-待部署物理机-导入","title": "导入待部署物理机模板"},{"children": [],"id": "button_predeploy_physical_machine_delete","note": "按钮-待部署物理机-删除","title": "删除物理机"}],"id": "menu_predeploy_physical_machine","note": "菜单-物理机管理-待部署物理机","title": "待部署物理机"},{"children": [{"children": [],"id": "button_device_setting_reInstall","note": "按钮-部署列表-重新部署","title": "重新部署"},{"children": [],"id": "button_dhcp_token_release","note": "按钮-部署列表-释放DHCP令牌","title": "释放DHCP令牌"},{"children": [],"id": "button_dhcp_token_batch_release","note": "按钮-部署列表-批量释放DHCP令牌","title": "批量释放DHCP令牌"},{"children": [],"id": "button_device_setting_cancelInstall","note": "按钮-部署列表-取消部署","title": "取消部署"},{"children": [],"id": "button_device_setting_delete","note": "按钮-部署列表-删除","title": "删除部署物理机"}],"id": "menu_device_setting","note": "菜单-物理机管理-部署列表","title": "部署列表"},{"children": [{"children": [],"id": "button_device_setting_rule_create","note": "按钮-装机参数规则-新增","title": "新增装机参数规则"},{"children": [],"id": "button_device_setting_rule_update","note": "按钮-装机参数规则-修改","title": "修改装机参数规则"},{"children": [],"id": "button_device_setting_rule_delete","note": "按钮-装机参数规则-删除","title": "删除装机参数规则"}],"id": "menu_device_setting_rule","note": "菜单-物理机管理-装机参数规则列表","title": "装机参数规则列表"},{"children": [{"children": [],"id": "button_inspection_addTask","note": "按钮-硬件巡检-新建巡检任务","title": "新建巡检任务"},{"children": [],"id": "button_inspection_inspect","note": "按钮-硬件巡检-重新巡检","title": "重新巡检"},{"children": [],"id": "button_inspection_inspect_all","note": "按钮-硬件巡检-巡检全部","title": "巡检全部"}],"id": "menu_inspection","note": "菜单-物理机管理-硬件巡检","title": "硬件巡检"}],"id": "menu_physical_machine_management","note": "菜单-物理机管理","title": "物理机管理"},{"children": [{"children": [{"children": [],"id": "button_order_create","note": "按钮-订单列表-新增","title": "新增订单"},{"children": [],"id": "button_order_delete","note": "按钮-订单列表-删除","title": "删除订单"},{"children": [],"id": "button_order_cancel","note": "按钮-订单列表-取消","title": "取消订单"},{"children": [],"id": "button_order_confirm","note": "按钮-订单列表-确认","title": "确认订单"},{"children": [],"id": "button_order_export","note": "按钮-订单列表-导出","title": "导出订单"}],"id": "menu_order","note": "菜单-订单管理-订单列表","title": "订单列表"},{"children": [{"children": [],"id": "button_device_category_create","note": "按钮-设备类型列表-新增","title": "新增设备类型"},{"children": [],"id": "button_device_category_delete","note": "按钮-设备类型列表-删除","title": "删除设备类型"},{"children": [],"id": "button_device_category_update","note": "按钮-设备类型列表-修改","title": "修改设备类型"}],"id": "menu_device_category","note": "菜单-订单管理-设备类型列表","title": "设备类型列表"}],"id": "menu_order_management","note": "菜单-订单管理","title": "订单管理"},{"children": [{"children": [{"children": [],"id": "button_system_template_create","note": "按钮-PXE配置-新建","title": "新建PXE配置"},{"children": [],"id": "button_system_template_update","note": "按钮-PXE配置-修改","title": "修改PXE配置"},{"children": [],"id": "button_system_template_delete","note": "按钮-PXE配置-删除","title": "删除PXE配置"},{"children": [],"id": "button_mirror_template_create","note": "按钮-镜像配置-新建","title": "新建镜像配置"},{"children": [],"id": "button_mirror_template_update","note": "按钮-镜像配置-修改","title": "修改镜像配置"},{"children": [],"id": "button_mirror_template_delete","note": "按钮-镜像配置-删除","title": "删除镜像配置"}],"id": "menu_system_template","note": "菜单-装机管理-系统/镜像配置","title": "系统/镜像配置"},{"children": [{"children": [],"id": "button_hardware_template_create","note": "按钮-硬件配置-新建","title": "新建/克隆硬件配置"},{"children": [],"id": "button_hardware_template_update","note": "按钮-硬件配置-修改","title": "修改硬件配置"},{"children": [],"id": "button_hardware_template_delete","note": "按钮-硬件配置-删除","title": "删除硬件配置"}],"id": "menu_hardware_template","note": "菜单-配置管理-硬件配置","title": "硬件配置"}],"id": "menu_template_management","note": "菜单-配置管理","title": "配置管理"},{"children": [{"children": [],"id": "button_approval_idc_abolish","note": "按钮-审批管理-数据中心裁撤","title": "数据中心裁撤审批"},{"children": [],"id": "button_approval_server_room_abolish","note": "按钮-审批管理-机房裁撤","title": "机房裁撤审批"},{"children": [],"id": "button_approval_cabinet_offline","note": "按钮-审批管理-机架下线","title": "机架下线审批"},{"children": [],"id": "button_approval_cabinet_powerOff","note": "按钮-审批管理-机架关电","title": "机架关电审批"},{"children": [],"id": "button_approval_network_area_offline","note": "按钮-审批管理-网络区域下线","title": "网络区域下线审批"},{"children": [],"id": "button_approval_physical_machine_move","note": "按钮-审批管理-物理机搬迁","title": "物理机搬迁审批"},{"children": [],"id": "button_approval_physical_machine_move_download","note": "按钮-审批管理-物理机搬迁导入模板","title": "下载物理机搬迁审批导入模板"},{"children": [],"id": "button_approval_physical_machine_move_import","note": "按钮-审批管理-导入物理机搬迁","title": "导入物理机搬迁审批"},{"children": [],"id": "button_approval_physical_machine_retirement","note": "按钮-审批管理-物理机退役","title": "物理机退役审批"},{"children": [],"id": "button_approval_physical_machine_reInstall","note": "按钮-审批管理-物理机重装","title": "物理机重装审批"},{"children": [],"id": "button_approval_physical_machine_power_off","note": "按钮-审批管理-物理机关电","title": "物理机关电审批"},{"children": [],"id": "button_approval_physical_machine_restart","note": "按钮-审批管理-物理机重启","title": "物理机重启审批"},{"children": [],"id": "button_approval_physical_machine_recycle_retire","note": "按钮-审批管理-物理机回收退役","title": "物理机回收退役审批"},{"children": [],"id": "button_approval_physical_machine_recycle_move","note": "按钮-审批管理-物理机回收搬迁","title": "物理机回收搬迁审批"},{"children": [],"id": "button_approval_physical_machine_recycle_reInstall","note": "按钮-审批管理-物理机回收重装","title": "物理机回收重装审批"},{"children": [],"id": "button_ip_unassign","note": "按钮-审批管理-IP回收","title": "IP回收审批"},{"children": [],"id": "button_approval_agree","note": "按钮-审批管理-审批通过","title": "审批通过"},{"children": [],"id": "button_approval_disagree","note": "按钮-审批管理-审批拒绝","title": "审批拒绝"},{"children": [],"id": "button_approval_revoke","note": "按钮-审批管理-审批撤销","title": "审批撤销"}],"id": "menu_approval","note": "菜单-审批管理","title": "审批管理"},{"children": [{"children": [],"id": "button_task_pause","note": "按钮-任务列表-暂停","title": "暂停"},{"children": [],"id": "button_task_continue","note": "按钮-任务列表-继续","title": "继续"},{"children": [],"id": "button_task_delete","note": "按钮-任务列表-删除","title": "删除"}],"id": "menu_task_management","note": "菜单-任务管理","title": "任务管理"}],"id": "","note": "","title": "所有权限"},"message": "操作成功","status": "success"}', '按钮权限列表', NULL),
	(5, '2021-09-01 15:00:00', '2021-09-01 15:00:00', NULL, 'authorization', '[{"api": {"desc": "查看数据中心","method": "GET","uri_regexp": "^/api/cloudboot/v1/idcs$"},"codes": ["menu_idc","menu_idc_management"]},{"api": {"desc": "新增数据中心","method": "POST","uri_regexp": "^/api/cloudboot/v1/idcs$"},"codes": ["button_idc_create"]},{"api": {"desc": "修改数据中心","method": "PUT","uri_regexp": "^/api/cloudboot/v1/idcs/[0-9A-Za-z_]+$"},"codes": ["button_idc_update"]},{"api": {"desc": "修改数据中心状态","method": "PUT","uri_regexp": "^/api/cloudboot/v1/idcs/status$"},"codes": ["button_idc_accepted","button_idc_abolished","button_idc_production"]},{"api": {"desc": "删除数据中心","method": "DELETE","uri_regexp": "^/api/cloudboot/v1/idcs/[0-9A-Za-z_]+$"},"codes": ["button_idc_delete"]},{"api": {"desc": "查看机房","method": "GET","uri_regexp": "^/api/cloudboot/v1/server-rooms$"},"codes": ["menu_server_room"]},{"api": {"desc": "新增机房","method": "POST","uri_regexp": "^/api/cloudboot/v1/server-rooms$"},"codes": ["button_server_room_create"]},{"api": {"desc": "修改机房","method": "PUT","uri_regexp": "^/api/cloudboot/v1/server-rooms$"},"codes": ["button_server_room_update"]},{"api": {"desc": "修改机房状态","method": "PUT","uri_regexp": "^/api/cloudboot/v1/server-rooms/status$"},"codes": ["button_server_room_production","button_server_room_abolished","button_server_room_accepted"]},{"api": {"desc": "删除机房","method": "DELETE","uri_regexp": "^/api/cloudboot/v1/server-rooms/[0-9A-Za-z_]+$"},"codes": ["button_server_room_delete"]},{"api": {"desc": "加载导入机房文件","method": "POST","uri_regexp": "^/api/cloudboot/v1/server-rooms/upload$"},"codes": ["button_server_room_import"]},{"api": {"desc": "导入机房文件预览","method": "POST","uri_regexp": "^/api/cloudboot/v1/server-rooms/imports/previews$"},"codes": ["button_server_room_import"]},{"api": {"desc": "导入机房文件","method": "POST","uri_regexp": "^/api/cloudboot/v1/server-rooms/imports$"},"codes": ["button_server_room_import"]},{"api": {"desc": "查看网络区域","method": "GET","uri_regexp": "^/api/cloudboot/v1/network-areas$"},"codes": ["menu_network_area"]},{"api": {"desc": "新增网络区域","method": "POST","uri_regexp": "^/api/cloudboot/v1/network-areas$"},"codes": ["button_network_area_create"]},{"api": {"desc": "修改网络区域","method": "PUT","uri_regexp": "^/api/cloudboot/v1/network-areas$"},"codes": ["button_network_area_update"]},{"api": {"desc": "修改网络区域状态","method": "PUT","uri_regexp": "^/api/cloudboot/v1/network-areas/status$"},"codes": ["button_network_area_production","button_network_area_offline"]},{"api": {"desc": "删除网络区域","method": "DELETE","uri_regexp": "^/api/cloudboot/v1/network-areas/[0-9A-Za-z_]+$"},"codes": ["button_network_area_delete"]},{"api": {"desc": "加载导入导入网络区域文件","method": "POST","uri_regexp": "^/api/cloudboot/v1/network-areas/upload$"},"codes": ["button_network_area_import"]},{"api": {"desc": "导入网络区域文件预览","method": "POST","uri_regexp": "^/api/cloudboot/v1/network-areas/imports/previews$"},"codes": ["button_network_area_import"]},{"api": {"desc": "导入网络区域文件","method": "POST","uri_regexp": "^/api/cloudboot/v1/network-areas/imports$"},"codes": ["button_network_area_import"]},{"api": {"desc": "查看机架(柜)","method": "GET","uri_regexp": "^/api/cloudboot/v1/server-cabinets$"},"codes": ["menu_server_cabinet"]},{"api": {"desc": "新增机架(柜)","method": "POST","uri_regexp": "^/api/cloudboot/v1/server-cabinets$"},"codes": ["button_server_cabinet_create"]},{"api": {"desc": "修改机架(柜)","method": "PUT","uri_regexp": "^/api/cloudboot/v1/server-cabinets$"},"codes": ["button_server_cabinet_update"]},{"api": {"desc": "修改机架(柜)状态","method": "PUT","uri_regexp": "^/api/cloudboot/v1/server-cabinets/status$"},"codes": ["button_server_cabinet_powerOn","button_server_cabinet_powerOff","button_server_cabinet_offline","button_server_cabinet_enabled"]},{"api": {"desc": "删除机架(柜)","method": "DELETE","uri_regexp": "^/api/cloudboot/v1/server-cabinets/[0-9A-Za-z_]+$"},"codes": ["button_server_cabinet_delete"]},{"api": {"desc": "机架开电","method": "POST","uri_regexp": "^/api/cloudboot/v1/server-cabinets/[0-9A-Za-z_]+/power$"},"codes": ["button_server_cabinet_powerOn"]},{"api": {"desc": "机架关电","method": "DELETE","uri_regexp": "^/api/cloudboot/v1/server-cabinets/[0-9A-Za-z_]+/power$"},"codes": ["button_server_cabinet_powerOff"]},{"api": {"desc": "加载导入机架(柜)文件","method": "POST","uri_regexp": "^/api/cloudboot/v1/server-cabinets/upload$"},"codes": ["button_server_cabinet_import"]},{"api": {"desc": "导入机架(柜)文件预览","method": "POST","uri_regexp": "^/api/cloudboot/v1/server-cabinets/imports/previews$"},"codes": ["button_server_cabinet_import"]},{"api": {"desc": "导入机架(柜)文件","method": "POST","uri_regexp": "^/api/cloudboot/v1/server-cabinets/imports$"},"codes": ["button_server_cabinet_import"]},{"api": {"desc": "查看机位(U位)","method": "GET","uri_regexp": "^/api/cloudboot/v1/server-usites$"},"codes": ["menu_server_usite"]},{"api": {"desc": "新增机位(U位)","method": "POST","uri_regexp": "^/api/cloudboot/v1/server-usites$"},"codes": ["button_server_usite_create"]},{"api": {"desc": "修改机位(U位)","method": "PUT","uri_regexp": "^/api/cloudboot/v1/server-usites$"},"codes": ["button_server_usite_update"]},{"api": {"desc": "删除机位(U位)","method": "DELETE","uri_regexp": "^/api/cloudboot/v1/server-usites/[0-9A-Za-z_]+$"},"codes": ["button_server_usite_delete"]},{"api": {"desc": "删除机位(U位)端口号","method": "DELETE","uri_regexp": "^/api/cloudboot/v1/server-usites/[0-9A-Za-z_]+/ports$"},"codes": ["button_server_usite_delete_port"]},{"api": {"desc": "修改机位(U位)状态","method": "PUT","uri_regexp": "^/api/cloudboot/v1/server-usites/status$"},"codes": ["button_server_usite_status"]},{"api": {"desc": "加载导入机位(U位)文件","method": "POST","uri_regexp": "^/api/cloudboot/v1/server-usites/upload$"},"codes": ["button_server_usite_import"]},{"api": {"desc": "导入机位(U位)文件预览","method": "POST","uri_regexp": "^/api/cloudboot/v1/server-usites/imports/previews$"},"codes": ["button_server_usite_import"]},{"api": {"desc": "导入机位(U位)文件","method": "POST","uri_regexp": "^/api/cloudboot/v1/server-usites/imports$"},"codes": ["button_server_usite_import"]},{"api": {"desc": "加载导入机位(U位)端口号文件","method": "POST","uri_regexp": "^/api/cloudboot/v1/server-usites/ports/upload$"},"codes": ["button_server_usite_import_port"]},{"api": {"desc": "导入机位(U位)端口号文件预览","method": "POST","uri_regexp": "^/api/cloudboot/v1/server-usites/ports/imports/previews$"},"codes": ["button_server_usite_import_port"]},{"api": {"desc": "导入机位(U位)端口号文件","method": "POST","uri_regexp": "^/api/cloudboot/v1/server-usites/ports/imports$"},"codes": ["button_server_usite_import_port"]},{"api": {"desc": "查看库房管理单元","method": "GET","uri_regexp": "^/api/cloudboot/v1/store-room$"},"codes": ["menu_store_room"]},{"api": {"desc": "新增库房管理单元","method": "POST","uri_regexp": "^/api/cloudboot/v1/store-room$"},"codes": ["button_store_room_create"]},{"api": {"desc": "上传导入库房管理单元文件","method": "POST","uri_regexp": "^/api/cloudboot/v1/store-room/upload$"},"codes": ["button_store_room_import"]},{"api": {"desc": "导入库房管理单元文件预览","method": "POST","uri_regexp": "^/api/cloudboot/v1/store-room/imports/previews$"},"codes": ["button_store_room_import"]},{"api": {"desc": "导入库房管理单元","method": "POST","uri_regexp": "^/api/cloudboot/v1/store-room/imports$"},"codes": ["button_store_room_import"]},{"api": {"desc": "删除库房管理单元","method": "DELETE","uri_regexp": "^/api/cloudboot/v1/store-room/[0-9]+$"},"codes": ["button_store_room_delete"]},{"api": {"desc": "修改库房管理单元","method": "PUT","uri_regexp": "^/api/cloudboot/v1/store-room$"},"codes": ["button_store_room_update"]},{"api": {"desc": "新增虚拟货架","method": "POST","uri_regexp": "^/api/cloudboot/v1/virtual-cabinet$"},"codes": ["button_virtual_cabinet_create"]},{"api": {"desc": "删除虚拟货架","method": "DELETE","uri_regexp": "^/api/cloudboot/v1/virtual-cabinet/[0-9]+$"},"codes": ["button_virtual_cabinet_delete"]},{"api": {"desc": "新增带外管理","method": "POST","uri_regexp": "^/api/cloudboot/v1/devices/power$"},"codes": ["button_physical_machine_powerOn","button_predeploy_physical_machine_powerOn"]},{"api": {"desc": "带外管理批量PXE重启","method": "PUT","uri_regexp": "^/api/cloudboot/v1/devices/power/pxe/restart$"},"codes": ["button_physical_machine_networkBoot","button_predeploy_physical_machine_networkBoot"]},{"api": {"desc": "带外管理批量重启","method": "PUT","uri_regexp": "^/api/cloudboot/v1/devices/power/restart$"},"codes": ["button_physical_machine_reBoot","button_predeploy_physical_machine_reBoot"]},{"api": {"desc": "带外管理批量关机","method": "DELETE","uri_regexp": "^/api/cloudboot/v1/devices/power$"},"codes": ["button_physical_machine_powerOff","button_predeploy_physical_machine_powerOff"]},{"api": {"desc": "修改带外用户密码","method": "PUT","uri_regexp": "^/api/cloudboot/v1/devices/[0-9A-Za-z_]+/oob/password$"},"codes": ["button_physical_machine_update_oob"]},{"api": {"desc": "新增物理机","method": "POST","uri_regexp": "^/api/cloudboot/v1/devices/[0-9A-Za-z_]+/collections$"},"codes": []},{"api": {"desc": "加载导入物理机文件","method": "POST","uri_regexp": "^/api/cloudboot/v1/devices/upload$"},"codes": ["button_predeploy_physical_machine_import"]},{"api": {"desc": "导入物理机文件预览","method": "POST","uri_regexp": "^/api/cloudboot/v1/devices/imports/previews$"},"codes": ["button_predeploy_physical_machine_import"]},{"api": {"desc": "加入物理机文件","method": "POST","uri_regexp": "^/api/cloudboot/v1/devices/imports$"},"codes": ["button_predeploy_physical_machine_import"]},{"api": {"desc": "查看物理机","method": "GET","uri_regexp": "^/api/cloudboot/v1/devices$"},"codes": ["menu_physical_machine"]},{"api": {"desc": "修改物理机","method": "PUT","uri_regexp": "^/api/cloudboot/v1/device$"},"codes": ["button_physical_machine_update"]},{"api": {"desc": "存量物理机导入","method": "POST","uri_regexp": "^/api/cloudboot/v1/devices/stock/upload$"},"codes": ["button_physical_machine_import"]},{"api": {"desc": "存量物理机导入预览","method": "POST","uri_regexp": "^/api/cloudboot/v1/devices/stock/imports/previews$"},"codes": ["button_physical_machine_import"]},{"api": {"desc": "加载存量物理机导入文件","method": "POST","uri_regexp": "^/api/cloudboot/v1/devices/stock//imports$"},"codes": ["button_physical_machine_import"]},{"api": {"desc": "上传物理机导入到库房文件","method": "POST","uri_regexp": "^/api/cloudboot/v1/devices/store/upload$"},"codes": ["button_device_store_import"]},{"api": {"desc": "物理机导入到库房预览","method": "POST","uri_regexp": "^/api/cloudboot/v1/devices/store/imports/previews$"},"codes": ["button_device_store_import"]},{"api": {"desc": "物理机导入到库房","method": "POST","uri_regexp": "^/api/cloudboot/v1/devices/store/imports$"},"codes": ["button_device_store_import"]},{"api": {"desc": "重新批量纳管带外","method": "PUT","uri_regexp": "^/api/cloudboot/v1/devices/oob/re-access$"},"codes": ["button_oob_re_access"]},{"api": {"desc": "上传特殊设备文件","method": "POST","uri_regexp": "^/api/cloudboot/v1/special-devices/upload$"},"codes": ["button_special_device_import"]},{"api": {"desc": "特殊设备预览","method": "POST","uri_regexp": "^/api/cloudboot/v1/special-devices/imports/previews$"},"codes": ["button_special_device_import"]},{"api": {"desc": "导入特殊设备","method": "POST","uri_regexp": "^/api/cloudboot/v1/special-devices/imports$"},"codes": ["button_special_device_import"]},{"api": {"desc": "删除物理机","method": "DELETE","uri_regexp": "^/api/cloudboot/v1/devices/[0-9A-Za-z_]+$"},"codes": ["button_physical_machine_delete","button_predeploy_physical_machine_delete"]},{"api": {"desc": "查看订单","method": "GET","uri_regexp": "^/api/cloudboot/v1/order$"},"codes": ["menu_order"]},{"api": {"desc": "新增订单","method": "POST","uri_regexp": "^/api/cloudboot/v1/order$"},"codes": ["button_order_create"]},{"api": {"desc": "删除订单","method": "DELETE","uri_regexp": "^/api/cloudboot/v1/orders$"},"codes": ["button_order_delete"]},{"api": {"desc": "修改订单状态","method": "PUT","uri_regexp": "^/api/cloudboot/v1/order/status$"},"codes": ["button_order_cancel","button_order_confirm"]},{"api": {"desc": "导出订单","method": "GET","uri_regexp": "^/api/cloudboot/v1/orders/export$"},"codes": ["button_order_export"]},{"api": {"desc": "查看设备类型","method": "GET","uri_regexp": "^/api/cloudboot/v1/device-category$"},"codes": ["menu_device_category"]},{"api": {"desc": "新增设备类型","method": "POST","uri_regexp": "^/api/cloudboot/v1/device-category$"},"codes": ["button_device_category_create"]},{"api": {"desc": "修改设备类型","method": "POST","uri_regexp": "^/api/cloudboot/v1/device-category$"},"codes": ["button_device_category_update"]},{"api": {"desc": "删除设备类型","method": "DELETE","uri_regexp": "^/api/cloudboot/v1/device-categories$"},"codes": ["button_device_category_delete"]},{"api": {"desc": "查看网段","method": "GET","uri_regexp": "^/api/cloudboot/v1/ip-networks$"},"codes": ["menu_ip_network"]},{"api": {"desc": "新增网段","method": "POST","uri_regexp": "^/api/cloudboot/v1/ip-networks$"},"codes": ["button_ip_network_create"]},{"api": {"desc": "修改网段","method": "POST","uri_regexp": "^/api/cloudboot/v1/ip-networks$"},"codes": ["button_ip_network_update"]},{"api": {"desc": "删除网段","method": "DELETE","uri_regexp": "^/api/cloudboot/v1/ip-networks/[0-9A-Za-z_]+$"},"codes": ["button_ip_network_delete"]},{"api": {"desc": "上传导入IP网段文件","method": "POST","uri_regexp": "^/api/cloudboot/v1/ip-networks/upload$"},"codes": ["button_ip_network_import"]},{"api": {"desc": "导入IP网段文件预览","method": "POST","uri_regexp": "^/api/cloudboot/v1/ip-networks/imports/previews$"},"codes": ["button_ip_network_import"]},{"api": {"desc": "导入IP网段文件","method": "POST","uri_regexp": "^/api/cloudboot/v1/ip-networks/imports$"},"codes": ["button_ip_network_import"]},{"api": {"desc": "查看IP列表","method": "GET","uri_regexp": "^/api/cloudboot/v1/ips$"},"codes": ["menu_ip"]},{"api": {"desc": "手动分配IP","method": "PUT","uri_regexp": "^/api/cloudboot/v1/ips/assigns$"},"codes": ["button_ip_assign"]},{"api": {"desc": " 手动取消IP分配","method": "PUT","uri_regexp": "^/api/cloudboot/v1/ips/unassigns$"},"codes": ["button_ip_unassign"]},{"api": {"desc": "删除网络设备","method": "DELETE","uri_regexp": "^/api/cloudboot/v1/network/devices/[0-9A-Za-z_]+$"},"codes": ["button_network_device_delete"]},{"api": {"desc": "查看网络设备","method": "GET","uri_regexp": "^/api/cloudboot/v1/network/devices$"},"codes": ["menu_network_device"]},{"api": {"desc": "新增网络设备","method": "POST","uri_regexp": "^/api/cloudboot/v1/network/devices$"},"codes": ["button_network_device_create"]},{"api": {"desc": "网络设备同步","method": "POST","uri_regexp": "^/api/cloudboot/v1/network/devices/sync$"},"codes": ["button_network_device_sync"]},{"api": {"desc": "上传网络设备文件","method": "POST","uri_regexp": "^/api/cloudboot/v1/network/devices/upload$"},"codes": ["button_ip_network_import"]},{"api": {"desc": "导入网络设备文件预览","method": "POST","uri_regexp": "^/api/cloudboot/v1/network/devices/imports/previews$"},"codes": ["button_ip_network_import"]},{"api": {"desc": "导入网络设备文件","method": "POST","uri_regexp": "^/api/cloudboot/v1/network/devices/imports$"},"codes": ["button_ip_network_import"]},{"api": {"desc": "上报安装进度","method": "POST","uri_regexp": "^/api/cloudboot/v1/devices/[0-9A-Za-z_]+/installations/progress$"},"codes": []},{"api": {"desc": "查看装机参数","method": "GET","uri_regexp": "^/api/cloudboot/v1/devices/settings$"},"codes": ["menu_device_setting"]},{"api": {"desc": "新增装机参数","method": "POST","uri_regexp": "^/api/cloudboot/v1/devices/settings$"},"codes": ["button_predeploy_physical_machine_osInstall"]},{"api": {"desc": "批量重装设备","method": "PUT","uri_regexp": "^/api/cloudboot/v1/devices/installations/reinstalls$"},"codes": ["button_device_setting_reInstall"]},{"api": {"desc": "批量取消安装设备","method": "PUT","uri_regexp": "^/api/cloudboot/v1/devices/installations/cancels$"},"codes": ["button_device_setting_cancelInstall"]},{"api": {"desc": "批量删除设备装机参数","method": "DELETE","uri_regexp": "^/api/cloudboot/v1/devices/settings$"},"codes": ["button_device_setting_delete"]},{"api": {"desc": "删除硬件模板","method": "DELETE","uri_regexp": "^/api/cloudboot/v1/hardware-templates/[0-9A-Za-z_]+$"},"codes": ["button_hardware_template_delete"]},{"api": {"desc": "修改硬件模板","method": "PUT","uri_regexp": "^/api/cloudboot/v1/hardware-templates$"},"codes": ["button_hardware_template_update"]},{"api": {"desc": "查看硬件模板","method": "GET","uri_regexp": "^/api/cloudboot/v1/hardware-templates$"},"codes": ["menu_hardware_template"]},{"api": {"desc": "新增硬件模板","method": "POST","uri_regexp": "^/api/cloudboot/v1/hardware-templates$"},"codes": ["button_hardware_template_create"]},{"api": {"desc": "查看镜像模板","method": "GET","uri_regexp": "^/api/cloudboot/v1/image-templates$"},"codes": ["menu_template_management"]},{"api": {"desc": "新增镜像模板","method": "POST","uri_regexp": "^/api/cloudboot/v1/image-templates$"},"codes": ["button_mirror_template_create"]},{"api": {"desc": "修改镜像模板","method": "PUT","uri_regexp": "^/api/cloudboot/v1/image-templates$"},"codes": ["button_mirror_template_update"]},{"api": {"desc": "删除镜像模板","method": "DELETE","uri_regexp": "^/api/cloudboot/v1/image-templates/[0-9A-Za-z_]+$"},"codes": ["button_mirror_template_delete"]},{"api": {"desc": "查看系统模板","method": "GET","uri_regexp": "^/api/cloudboot/v1/system-templates$"},"codes": ["menu_template_management"]},{"api": {"desc": "新增系统模板","method": "POST","uri_regexp": "^/api/cloudboot/v1/system-templates$"},"codes": ["button_system_template_create"]},{"api": {"desc": "修改系统模板","method": "PUT","uri_regexp": "^/api/cloudboot/v1/system-templates$"},"codes": ["button_system_template_update"]},{"api": {"desc": "删除系统模板","method": "DELETE","uri_regexp": "^/api/cloudboot/v1/system-templates/[0-9A-Za-z_]+$"},"codes": ["button_system_template_delete"]},{"api": {"desc": "查看巡检任务","method": "GET","uri_regexp": "^/api/cloudboot/v1/jobs/inspections$"},"codes": ["menu_inspection"]},{"api": {"desc": "新增巡检任务","method": "POST","uri_regexp": "^/api/cloudboot/v1/jobs/inspections$"},"codes": ["button_inspection_addTask","button_inspection_inspect","button_inspection_inspect_all"]},{"api": {"desc": "查看任务列表","method": "GET","uri_regexp": "^/api/cloudboot/v1/jobs$"},"codes": ["menu_task_management"]},{"api": {"desc": "删除目标任务API","method": "DELETE","uri_regexp": "^/api/cloudboot/v1/jobs/[0-9A-Za-z_]+$"},"codes": ["button_task_delete"]},{"api": {"desc": "暂停目标任务API","method": "PUT","uri_regexp": "^/api/cloudboot/v1/jobs/[0-9A-Za-z_]+/pausing$"},"codes": ["button_task_pause"]},{"api": {"desc": "继续已暂停的目标任务API","method": "PUT","uri_regexp": "^/api/cloudboot/v1/jobs/[0-9A-Za-z_]+/unpausing$"},"codes": ["button_task_continue"]},{"api": {"desc": "为目标设备生成PXE文件","method": "POST","uri_regexp": "^/api/cloudboot/v1/devices/[0-9A-Za-z_]+/pxe$"},"codes": []},{"api": {"desc": "生成centos6 pxe模板","method": "POST","uri_regexp": "^/api/cloudboot/v1/devices/[0-9A-Za-z_]+/centos6/uefi/pxe$"},"codes": []},{"api": {"desc": "发起数据中心裁撤审批","method": "POST","uri_regexp": "^/api/cloudboot/v1/approvals/idc/abolish$"},"codes": ["button_approval_idc_abolish"]},{"api": {"desc": "发起机房裁撤审批","method": "POST","uri_regexp": "^/api/cloudboot/v1/approvals/server-room/abolish$"},"codes": ["button_approval_server_room_abolish"]},{"api": {"desc": "发起机架关电审批","method": "POST","uri_regexp": "^/api/cloudboot/v1/approvals/server-cabinets/poweroffs$"},"codes": ["button_approval_cabinet_powerOff"]},{"api": {"desc": "发起机架下线审批","method": "POST","uri_regexp": "^/api/cloudboot/v1/approvals/server-cabinets/offlines$"},"codes": ["button_approval_cabinet_offline"]},{"api": {"desc": "发起网络区域下线审批","method": "POST","uri_regexp": "^/api/cloudboot/v1/approvals/network-area/offline$"},"codes": ["button_approval_network_area_offline"]},{"api": {"desc": "发起物理机搬迁审批","method": "POST","uri_regexp": "^/api/cloudboot/v1/approvals/devices/migrations$"},"codes": ["button_approval_physical_machine_move"]},{"api": {"desc": "上传物理机搬迁审批文件","method": "POST","uri_regexp": "^/api/cloudboot/v1/approvals/devices/migrations/upload$"},"codes": ["button_approval_physical_machine_move_import"]},{"api": {"desc": "导入物理机搬迁审批预览","method": "POST","uri_regexp": "^/api/cloudboot/v1/approvals/devices/migrations/imports/previews$"},"codes": ["button_approval_physical_machine_move_import"]},{"api": {"desc": "导入物理机搬迁审批","method": "POST","uri_regexp": "^/api/cloudboot/v1/approvals/devices/migrations/imports$"},"codes": ["button_approval_physical_machine_move_import"]},{"api": {"desc": "发起物理机退役审批","method": "POST","uri_regexp": "^/api/cloudboot/v1/approvals/devices/retirements$"},"codes": ["button_approval_physical_machine_retirement"]},{"api": {"desc": "发起物理机重装审批","method": "POST","uri_regexp": "^/api/cloudboot/v1/approvals/devices/os-reinstallations$"},"codes": ["button_approval_physical_machine_reInstall"]},{"api": {"desc": "发起物理机关电审批","method": "POST","uri_regexp": "^/api/cloudboot/v1/approvals/device/poweroff$"},"codes": ["button_approval_physical_machine_power_off"]},{"api": {"desc": "发起物理机重启审批","method": "POST","uri_regexp": "^/api/cloudboot/v1/approvals/device/restart$"},"codes": ["button_approval_physical_machine_restart"]},{"api": {"desc": "发起物理机回收审批","method": "POST","uri_regexp": "^/api/cloudboot/v1/approvals/devices/recycle$"},"codes": ["button_approval_physical_machine_recycle_retire","button_approval_physical_machine_recycle_move","button_approval_physical_machine_recycle_reInstall"]},{"api": {"desc": "发起IP回收审批","method": "POST","uri_regexp": "^/api/cloudboot/v1/approvals/ip/unassign$"},"codes": ["button_ip_unassign"]},{"api": {"desc": "归还设备持有的DHCP IP令牌","method": "DELETE","uri_regexp": "^/api/cloudboot/v1/devices/[0-9A-Za-z_]+/limiters/tokens$"},"codes": ["menu_device_setting"]},{"api": {"desc": "一键归还所有设备持有的DHCP IP令牌","method": "DELETE","uri_regexp": "^/api/cloudboot/v1/devices/limiters/tokens$"},"codes": ["menu_device_setting"]},{"api": {"desc": "查看装机参数规则","method": "GET","uri_regexp": "^/api/cloudboot/v1/device-setting-rules$"},"codes": ["menu_device_setting_rule"]},{"api": {"desc": "新增装机参数规则","method": "POST","uri_regexp": "^/api/cloudboot/v1/device-setting-rules$"},"codes": ["button_device_setting_rule_create"]},{"api": {"desc": "修改装机参数规则","method": "PUT","uri_regexp": "^/api/cloudboot/v1/device-setting-rules$"},"codes": ["button_device_setting_rule_update"]},{"api": {"desc": "删除装机参数规则","method": "DELETE","uri_regexp": "^/api/cloudboot/v1/device-setting-rules$"},"codes": ["button_device_setting_rule_delete"]}]', 'API权限集合', NULL);
/*!40000 ALTER TABLE `system_setting` ENABLE KEYS */;

-- 导出  表 cloudboot_cyclone.system_template 结构
CREATE TABLE IF NOT EXISTS `system_template` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `created_at` timestamp NULL DEFAULT NULL COMMENT '记录创建时间',
  `updated_at` timestamp NULL DEFAULT NULL COMMENT '记录修改时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '记录删除时间',
  `family` varchar(255) DEFAULT NULL COMMENT '族系',
  `name` varchar(255) NOT NULL COMMENT '模板名',
  `boot_mode` enum('legacy_bios','uefi') DEFAULT 'legacy_bios' COMMENT '启动模式。可选值: legacy_bios-传统BIOS模式; uefi-UEFI模式;',
  `username` varchar(255) DEFAULT NULL COMMENT '操作系统用户名',
  `password` varchar(255) DEFAULT NULL COMMENT '操作系统用户密码',
  `content` longtext COMMENT '系统模板内容',
  `pxe` longtext COMMENT 'PXE引导配置',
  `os_lifecycle` enum('testing','active_default','active','containment','end_of_life') DEFAULT 'testing' COMMENT 'OS生命周期：Testing|Active(Default)|Active|Containment|EOL',
  `arch` enum('unknown','x86_64','aarch64') DEFAULT NULL COMMENT 'OS架构',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_system_template_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=37 DEFAULT CHARSET=utf8 COMMENT='系统安装模板表';

-- 正在导出表  cloudboot_cyclone.system_template 的数据：~12 rows (大约)
/*!40000 ALTER TABLE `system_template` DISABLE KEYS */;
REPLACE INTO `system_template` (`id`, `created_at`, `updated_at`, `deleted_at`, `family`, `name`, `boot_mode`, `username`, `password`, `content`, `pxe`, `os_lifecycle`, `arch`) VALUES
	(1, '2018-05-22 06:43:04', '2023-03-16 16:08:22', NULL, 'BootOS', 'bootos_x86_64', 'uefi', '', '', '#', '#!ipxe\nkernel http://osinstall.idcos.com/bootos/x86_64/vmlinuz initrd=initrd.img console=tty0 selinux=0 biosdevname=0 SERVER_ADDR=http://osinstall.idcos.com LOOP_INTERVAL=5 DEVELOPER=1 BOOTIF=01-${netX/mac:hexhyp}\ninitrd http://osinstall.idcos.com/bootos/x86_64/initrd.img\nboot', 'testing', 'x86_64'),
	(2, '2018-05-22 06:43:04', '2023-12-06 21:46:22', NULL, 'BootOS', 'bootos_arm64', 'uefi', '', '', '#', '#!ipxe\nkernel http://osinstall.idcos.com/bootos/aarch64/vmlinuz initrd=initrd.img console=tty0 selinux=0 biosdevname=0 SERVER_ADDR=http://osinstall.idcos.com LOOP_INTERVAL=5 DEVELOPER=1 BOOTIF=01-${netX/mac:hexhyp}\ninitrd http://osinstall.idcos.com/bootos/aarch64/initrd.img\nboot\n', 'testing', 'aarch64'),
	(4, '2018-05-22 06:53:22', '2023-03-16 16:09:35', NULL, 'BootOS', 'winpe2012_x86_64', 'legacy_bios', '', '', '#', '#!ipxe\nkernel http://osinstall.idcos.com/winpe/wimboot\ninitrd http://osinstall.idcos.com/winpe/2012/Boot/BCD BCD\ninitrd http://osinstall.idcos.com/winpe/2012/Boot/boot.sdi boot.sdi\ninitrd http://osinstall.idcos.com/winpe/2012/sources/boot.wim boot.wim\nboot', 'testing', 'x86_64'),
	(17, '2019-11-20 10:51:40', '2023-12-07 09:18:01', NULL, 'CentOS', 'CentOS_7.6_aarch64', 'uefi', 'root', 'Cyclone@1234', 'install\nurl --url=http://osinstall.idcos.com/centos/7.6/os/aarch64/\nlang en_US.UTF-8\nkeyboard us\nnetwork --onboot yes --device bootif --bootproto dhcp --noipv6\nrootpw  Cyclone@1234\nfirewall --disabled\nauthconfig --enableshadow --passalgo=sha512\nselinux --disabled\ntimezone Asia/Shanghai\ntext\nreboot\nzerombr\nbootloader --location=mbr\nclearpart --all --initlabel\npart /boot/efi --fstype=efi --size=200 --ondisk=sda \npart /boot --fstype=xfs --size=1024 --ondisk=sda\npart swap --size=8192 --ondisk=sda\npart / --fstype=xfs --size=30720 --ondisk=sda\npart /tmp --fstype=xfs --size=10240 --ondisk=sda\npart /home --fstype=xfs --size=5120 --ondisk=sda\npart /usr/local --fstype=xfs --size=20480 --ondisk=sda\npart /data --fstype=xfs --size=1 --grow --ondisk=sda\n\n%packages --ignoremissing\n@base\n@core\n@development\ndmidecode\n%end\n\n%pre\n_sn=$(dmidecode -s system-serial-number 2>/dev/null | awk \'/^[^#]/ { print $1 }\')\n\ncurl -H "Content-Type: application/json" -X POST -d "{\\"title\\":\\"启动OS安装程序\\",\\"progress\\":0.6,\\"log\\":\\"SW5zdGFsbCBPUwo=\\"}" http://osinstall.idcos.com/api/cloudboot/v1/devices/$_sn/installations/progress\ncurl -H "Content-Type: application/json" -X POST -d "{\\"title\\":\\"分区并安装软件包\\",\\"progress\\":0.7,\\"log\\":\\"SW5zdGFsbCBPUwo=\\"}" http://osinstall.idcos.com/api/cloudboot/v1/devices/$_sn/installations/progress\n%end\n\n%post\n# dowload commonSetting.py for setting items in OS\ncurl -o /tmp/commonSetting.py "http://osinstall.idcos.com/scripts/commonSetting.py"\n\n# config network\npython /tmp/commonSetting.py --network=Y \n\n# config osuser\npython /tmp/commonSetting.py --osuser=Y\n\n# complete\npython /tmp/commonSetting.py --complete=Y\n%end', '#!ipxe\nkernel http://osinstall.idcos.com/centos/7.6/os/aarch64/images/pxeboot/vmlinuz initrd=initrd.img ksdevice=bootif ks=http://osinstall.idcos.com/api/cloudboot/v1/devices/{sn}/settings/system-template console=tty0 selinux=0 net.ifnames=0 biosdevname=0 BOOTIF=01-${netX/mac:hexhyp}\ninitrd http://osinstall.idcos.com/centos/7.6/os/aarch64/images/pxeboot/initrd.img\nboot', 'active', 'aarch64'),
	(22, '2020-05-15 17:10:03', '2023-12-06 21:42:41', NULL, 'Windows Server', 'win2012r2', 'uefi', 'administrator', 'Cyclone@1234', '<?xml version="1.0" encoding="utf-8"?>\n<unattend xmlns="urn:schemas-microsoft-com:unattend">\n    <settings pass="generalize">\n        <component name="Microsoft-Windows-OutOfBoxExperience" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\n            <DoNotOpenInitialConfigurationTasksAtLogon>true</DoNotOpenInitialConfigurationTasksAtLogon>\n        </component>\n        <component name="Microsoft-Windows-ServerManager-SvrMgrNc" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\n            <DoNotOpenServerManagerAtLogon>true</DoNotOpenServerManagerAtLogon>\n        </component>\n    </settings>\n    <settings pass="specialize">\n        <component name="Microsoft-Windows-Shell-Setup" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\n            <ProductKey>26BKW-N43FC-2B86G-PT2FG-PDJF8</ProductKey>\n            <ShowWindowsLive>false</ShowWindowsLive>\n            <DisableAutoDaylightTimeSet>false</DisableAutoDaylightTimeSet>\n            <TimeZone>China Standard Time</TimeZone>\n			<ComputerName>WebankWinServer</ComputerName>\n            <RegisteredOwner>WeBank</RegisteredOwner>\n            <RegisteredOrganization>TCTP</RegisteredOrganization>\n        </component>\n        <component name="Microsoft-Windows-TerminalServices-LocalSessionManager" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\n            <fDenyTSConnections>false</fDenyTSConnections>\n        </component>\n        <component name="Networking-MPSSVC-Svc" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\n            <FirewallGroups>\n                <FirewallGroup wcm:action="add" wcm:keyValue="RemoteDesktop">\n                    <Active>true</Active>\n                    <Group>@FirewallAPI.dll,-28752</Group>\n                    <Profile>all</Profile>\n                </FirewallGroup>\n            </FirewallGroups>\n        </component>\n    </settings>\n    <settings pass="oobeSystem">\n        <component name="Microsoft-Windows-Shell-Setup" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\n            <AutoLogon>\n                <Password>\n                    <Value>Cyclone@1234</Value>\n                    <PlainText>true</PlainText>\n                </Password>\n                <Enabled>true</Enabled>\n                <LogonCount>1</LogonCount>\n                <Username>Administrator</Username>\n            </AutoLogon>\n            <OOBE>\n                <HideEULAPage>true</HideEULAPage>\n            </OOBE>\n            <UserAccounts>\n                <AdministratorPassword>\n                    <PlainText>true</PlainText>\n                    <Value>Cyclone@1234</Value>\n                </AdministratorPassword>\n            </UserAccounts>\n            <FirstLogonCommands>\n                <SynchronousCommand wcm:action="add">\n                    <Order>1</Order>\n                    <CommandLine>C:\\firstboot\\winconfig.exe</CommandLine>\n                    <Description></Description>\n                    <RequiresUserInput>false</RequiresUserInput>\n                </SynchronousCommand>\n            </FirstLogonCommands>\n        </component>\n        <component name="Microsoft-Windows-International-Core" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\n            <SystemLocale>zh-CN</SystemLocale>\n            <UILanguage>zh-CN</UILanguage>\n            <UILanguageFallback>zh-CN</UILanguageFallback>\n            <UserLocale>zh-CN</UserLocale>\n            <InputLocale>zh-CN</InputLocale>\n        </component>\n    </settings>\n    <settings pass="windowsPE">\n        <component name="Microsoft-Windows-Setup" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\n            <ImageInstall>\n                <OSImage>\n                    <InstallFrom>\n                        <Path>Z:\\windows\\2012r2\\sources\\install.wim</Path>\n                        <MetaData wcm:action="add">\n                            <Key>/IMAGE/NAME</Key>\n                            <Value>Windows Server 2012 R2 SERVERSTANDARD</Value>\n                        </MetaData>\n                    </InstallFrom>\n                    <InstallTo>\n                        <DiskID>0</DiskID>\n                        <PartitionID>3</PartitionID>\n                    </InstallTo>\n                    <WillShowUI>OnError</WillShowUI>\n                </OSImage>\n            </ImageInstall>\n            <UserData>\n                <ProductKey>\n                    <WillShowUI>OnError</WillShowUI>\n                    <Key>26BKW-N43FC-2B86G-PT2FG-PDJF8</Key>\n                </ProductKey>\n                <AcceptEula>true</AcceptEula>\n            </UserData>\n            <EnableFirewall>true</EnableFirewall>\n            <EnableNetwork>true</EnableNetwork>\n            <DiskConfiguration>\n                <Disk wcm:action="add">\n                    <CreatePartitions>\n                        <CreatePartition wcm:action="add">\n                            <Extend>false</Extend>\n                            <Order>1</Order>\n                            <Size>260</Size>\n                            <Type>EFI</Type>\n                        </CreatePartition>\n                        <CreatePartition wcm:action="add">\n                            <Extend>false</Extend>\n                            <Order>2</Order>\n                            <Size>16</Size>\n                            <Type>MSR</Type>\n                        </CreatePartition>\n                        <CreatePartition wcm:action="add">\n                            <Extend>false</Extend>\n                            <Order>3</Order>\n                            <Size>202400</Size>\n                            <Type>Primary</Type>\n                        </CreatePartition>\n                        <CreatePartition wcm:action="add">\n                            <Extend>false</Extend>\n                            <Order>4</Order>\n                            <Size>202400</Size>\n                            <Type>Primary</Type>\n                        </CreatePartition>\n                    </CreatePartitions>\n                    <ModifyPartitions>\n                        <ModifyPartition wcm:action="add">\n                            <Format>FAT32</Format>\n                            <Label>System</Label>\n                            <Order>1</Order>\n                            <PartitionID>1</PartitionID>\n                        </ModifyPartition>\n                        <ModifyPartition wcm:action="add">\n                            <Format>NTFS</Format>\n                            <Label>System</Label>\n                            <Letter>C</Letter>\n                            <Order>2</Order>\n                            <PartitionID>3</PartitionID>\n                        </ModifyPartition>\n                        <ModifyPartition wcm:action="add">\n                            <Extend>true</Extend>\n                            <Format>NTFS</Format>\n                            <Label>data</Label>\n                            <Letter>D</Letter>\n                            <Order>3</Order>\n                            <PartitionID>4</PartitionID>\n                        </ModifyPartition>\n                    </ModifyPartitions>\n                    <DiskID>0</DiskID>\n                    <WillWipeDisk>true</WillWipeDisk>\n                </Disk>\n                <WillShowUI>OnError</WillShowUI>\n            </DiskConfiguration>\n        </component>\n        <component name="Microsoft-Windows-International-Core-WinPE" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\n            <SetupUILanguage>\n                <UILanguage>zh-CN</UILanguage>\n            </SetupUILanguage>\n            <InputLocale>zh-CN</InputLocale>\n            <SystemLocale>zh-CN</SystemLocale>\n            <UILanguage>zh-CN</UILanguage>\n            <UserLocale>zh-CN</UserLocale>\n            <UILanguageFallback>zh-CN</UILanguageFallback>\n        </component>\n    </settings>\n    <settings pass="offlineServicing">\n        <component name="Microsoft-Windows-LUA-Settings" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\n            <EnableLUA>false</EnableLUA>\n        </component>\n    </settings>\n    <cpi:offlineImage cpi:source="catalog://osinstall/image/windows/clg/install_windows server 2012 r2 serverstandard.clg" xmlns:cpi="urn:schemas-microsoft-com:cpi" />\n</unattend>', '#!ipxe\nkernel http://osinstall.idcos.com/winpe/wimboot\ninitrd http://osinstall.idcos.com/winpe/2012/bootmgr bootmgr \ninitrd http://osinstall.idcos.com/winpe/2012/EFI/Boot/bootx64.efi bootx64.efi \ninitrd http://osinstall.idcos.com/winpe/2012/Boot/BCD BCD \ninitrd http://osinstall.idcos.com/winpe/2012/Boot/boot.sdi boot.sdi \ninitrd http://osinstall.idcos.com/winpe/2012/sources/boot.wim boot.wim \nboot', 'testing', 'x86_64'),
	(24, '2021-02-24 10:14:17', '2023-12-06 21:42:14', NULL, 'Windows Server', 'win2016', 'uefi', 'administrator', 'Cyclone@1234', '<?xml version="1.0" encoding="utf-8"?>\n<unattend xmlns="urn:schemas-microsoft-com:unattend">\n    <settings pass="generalize">\n        <component name="Microsoft-Windows-OutOfBoxExperience" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\n            <DoNotOpenInitialConfigurationTasksAtLogon>true</DoNotOpenInitialConfigurationTasksAtLogon>\n        </component>\n        <component name="Microsoft-Windows-ServerManager-SvrMgrNc" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\n            <DoNotOpenServerManagerAtLogon>true</DoNotOpenServerManagerAtLogon>\n        </component>\n    </settings>\n    <settings pass="specialize">\n        <component name="Microsoft-Windows-Shell-Setup" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\n            <ProductKey>NTYKT-H46BX-YH9D7-BJ7M3-GF47P</ProductKey>\n            <ShowWindowsLive>false</ShowWindowsLive>\n            <DisableAutoDaylightTimeSet>false</DisableAutoDaylightTimeSet>\n            <TimeZone>China Standard Time</TimeZone>\n			<ComputerName>WebankWinServer</ComputerName>\n            <RegisteredOwner>WeBank</RegisteredOwner>\n            <RegisteredOrganization>TCTP</RegisteredOrganization>\n        </component>\n        <component name="Microsoft-Windows-TerminalServices-LocalSessionManager" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\n            <fDenyTSConnections>false</fDenyTSConnections>\n        </component>\n        <component name="Networking-MPSSVC-Svc" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\n            <FirewallGroups>\n                <FirewallGroup wcm:action="add" wcm:keyValue="RemoteDesktop">\n                    <Active>true</Active>\n                    <Group>@FirewallAPI.dll,-28752</Group>\n                    <Profile>all</Profile>\n                </FirewallGroup>\n            </FirewallGroups>\n        </component>\n    </settings>\n    <settings pass="oobeSystem">\n        <component name="Microsoft-Windows-Shell-Setup" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\n            <AutoLogon>\n                <Password>\n                    <Value>Cyclone@1234</Value>\n                    <PlainText>true</PlainText>\n                </Password>\n                <Enabled>true</Enabled>\n                <LogonCount>1</LogonCount>\n                <Username>Administrator</Username>\n            </AutoLogon>\n            <OOBE>\n                <HideEULAPage>true</HideEULAPage>\n            </OOBE>\n            <UserAccounts>\n                <AdministratorPassword>\n                    <PlainText>true</PlainText>\n                    <Value>Cyclone@1234</Value>\n                </AdministratorPassword>\n            </UserAccounts>\n            <FirstLogonCommands>\n                <SynchronousCommand wcm:action="add">\n                    <Order>1</Order>\n                    <CommandLine>C:\\firstboot\\winconfig.exe</CommandLine>\n                    <Description></Description>\n                    <RequiresUserInput>false</RequiresUserInput>\n                </SynchronousCommand>\n            </FirstLogonCommands>\n        </component>\n        <component name="Microsoft-Windows-International-Core" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\n            <SystemLocale>zh-CN</SystemLocale>\n            <UILanguage>zh-CN</UILanguage>\n            <UILanguageFallback>zh-CN</UILanguageFallback>\n            <UserLocale>zh-CN</UserLocale>\n            <InputLocale>zh-CN</InputLocale>\n        </component>\n    </settings>\n    <settings pass="windowsPE">\n        <component name="Microsoft-Windows-Setup" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\n            <ImageInstall>\n                <OSImage>\n                    <InstallFrom>\n                        <Path>Z:\\windows\\2016\\sources\\install.wim</Path>\n                        <MetaData wcm:action="add">\n                            <Key>/IMAGE/NAME</Key>\n                            <Value>Windows Server 2016 SERVERSTANDARD</Value>\n                        </MetaData>\n                    </InstallFrom>\n                    <InstallTo>\n                        <DiskID>0</DiskID>\n                        <PartitionID>3</PartitionID>\n                    </InstallTo>\n                    <WillShowUI>OnError</WillShowUI>\n                </OSImage>\n            </ImageInstall>\n            <UserData>\n                <ProductKey>\n                    <WillShowUI>OnError</WillShowUI>\n                    <Key>NTYKT-H46BX-YH9D7-BJ7M3-GF47P</Key>\n                </ProductKey>\n                <AcceptEula>true</AcceptEula>\n            </UserData>\n            <EnableFirewall>true</EnableFirewall>\n            <EnableNetwork>true</EnableNetwork>\n            <DiskConfiguration>\n                <Disk wcm:action="add">\n                    <CreatePartitions>\n                        <CreatePartition wcm:action="add">\n                            <Extend>false</Extend>\n                            <Order>1</Order>\n                            <Size>260</Size>\n                            <Type>EFI</Type>\n                        </CreatePartition>\n                        <CreatePartition wcm:action="add">\n                            <Extend>false</Extend>\n                            <Order>2</Order>\n                            <Size>16</Size>\n                            <Type>MSR</Type>\n                        </CreatePartition>\n                        <CreatePartition wcm:action="add">\n                            <Extend>false</Extend>\n                            <Order>3</Order>\n                            <Size>202400</Size>\n                            <Type>Primary</Type>\n                        </CreatePartition>\n                        <CreatePartition wcm:action="add">\n                            <Extend>false</Extend>\n                            <Order>4</Order>\n                            <Size>202400</Size>\n                            <Type>Primary</Type>\n                        </CreatePartition>\n                    </CreatePartitions>\n                    <ModifyPartitions>\n                        <ModifyPartition wcm:action="add">\n                            <Format>FAT32</Format>\n                            <Label>System</Label>\n                            <Order>1</Order>\n                            <PartitionID>1</PartitionID>\n                        </ModifyPartition>\n                        <ModifyPartition wcm:action="add">\n                            <Format>NTFS</Format>\n                            <Label>System</Label>\n                            <Letter>C</Letter>\n                            <Order>2</Order>\n                            <PartitionID>3</PartitionID>\n                        </ModifyPartition>\n                        <ModifyPartition wcm:action="add">\n                            <Extend>true</Extend>\n                            <Format>NTFS</Format>\n                            <Label>data</Label>\n                            <Letter>D</Letter>\n                            <Order>3</Order>\n                            <PartitionID>4</PartitionID>\n                        </ModifyPartition>\n                    </ModifyPartitions>\n                    <DiskID>0</DiskID>\n                    <WillWipeDisk>true</WillWipeDisk>\n                </Disk>\n                <WillShowUI>OnError</WillShowUI>\n            </DiskConfiguration>\n        </component>\n        <component name="Microsoft-Windows-International-Core-WinPE" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\n            <SetupUILanguage>\n                <UILanguage>zh-CN</UILanguage>\n            </SetupUILanguage>\n            <InputLocale>zh-CN</InputLocale>\n            <SystemLocale>zh-CN</SystemLocale>\n            <UILanguage>zh-CN</UILanguage>\n            <UserLocale>zh-CN</UserLocale>\n            <UILanguageFallback>zh-CN</UILanguageFallback>\n        </component>\n    </settings>\n    <settings pass="offlineServicing">\n        <component name="Microsoft-Windows-LUA-Settings" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\n            <EnableLUA>false</EnableLUA>\n        </component>\n    </settings>\n    <cpi:offlineImage cpi:source="catalog://osinstall/image/windows/clg/install_windows server 2016 serverstandard.clg" xmlns:cpi="urn:schemas-microsoft-com:cpi" />\n</unattend>', '#!ipxe\nkernel http://osinstall.idcos.com/winpe/wimboot\ninitrd http://osinstall.idcos.com/winpe/2012/bootmgr bootmgr \ninitrd http://osinstall.idcos.com/winpe/2012/EFI/Boot/bootx64.efi bootx64.efi \ninitrd http://osinstall.idcos.com/winpe/2012/Boot/BCD BCD \ninitrd http://osinstall.idcos.com/winpe/2012/Boot/boot.sdi boot.sdi \ninitrd http://osinstall.idcos.com/winpe/2012/sources/boot.wim boot.wim \nboot', 'testing', 'x86_64'),
	(25, '2021-03-21 17:07:22', '2023-12-07 09:19:23', NULL, 'CentOS', 'CentOS 7.9', 'uefi', 'root', 'Cyclone@1234', 'install\nurl --url=http://osinstall.idcos.com/centos/7.9/os/x86_64/\nlang en_US.UTF-8\nkeyboard us\nnetwork --onboot yes --device bootif --bootproto dhcp --noipv6\nrootpw  Cyclone@1234\nfirewall --disabled\nauthconfig --enableshadow --passalgo=sha512\nselinux --disabled\ntimezone Asia/Shanghai\ntext\nreboot\nzerombr\nbootloader --location=mbr\nclearpart --all --initlabel\npart /boot/efi --fstype=efi --size=200 --ondisk=sda \npart /boot --fstype=ext4 --size=1024 --ondisk=sda\npart swap --size=8192 --ondisk=sda\npart / --fstype=ext4 --size=30720 --ondisk=sda\npart /tmp --fstype=ext4 --size=10240 --ondisk=sda\npart /home --fstype=ext4 --size=5120 --ondisk=sda\npart /usr/local --fstype=ext4 --size=20480 --ondisk=sda\npart /data --fstype=ext4 --size=1 --grow --ondisk=sda\n\n%packages --ignoremissing\n@base\n@core\n@development\ndmidecode\n%end\n\n%pre\n_sn=$(dmidecode -s system-serial-number 2>/dev/null | awk \'/^[^#]/ { print $1 }\')\n\ncurl -H "Content-Type: application/json" -X POST -d "{\\"title\\":\\"启动OS安装程序\\",\\"progress\\":0.6,\\"log\\":\\"SW5zdGFsbCBPUwo=\\"}" http://osinstall.idcos.com/api/cloudboot/v1/devices/$_sn/installations/progress\ncurl -H "Content-Type: application/json" -X POST -d "{\\"title\\":\\"分区并安装软件包\\",\\"progress\\":0.7,\\"log\\":\\"SW5zdGFsbCBPUwo=\\"}" http://osinstall.idcos.com/api/cloudboot/v1/devices/$_sn/installations/progress\n%end\n\n%post\n# dowload commonSetting.py for setting items in OS\ncurl -o /tmp/commonSetting.py "http://osinstall.idcos.com/scripts/commonSetting.py"\n\n# config network\npython /tmp/commonSetting.py --network=Y\n\n# change root passwd\npython /tmp/commonSetting.py --osuser=Y\n\n# complete\npython /tmp/commonSetting.py --complete=Y\n%end', '#!ipxe\nkernel http://osinstall.idcos.com/centos/7.9/os/x86_64/images/pxeboot/vmlinuz initrd=initrd.img ksdevice=bootif ks=http://osinstall.idcos.com/api/cloudboot/v1/devices/${serial}/settings/system-template console=tty0 selinux=0 net.ifnames=0 biosdevname=0 BOOTIF=01-${netX/mac:hexhyp}\ninitrd http://osinstall.idcos.com/centos/7.9/os/x86_64/images/pxeboot/initrd.img\nboot\n', 'active', 'x86_64'),
	(28, '2022-04-23 10:11:26', '2023-12-07 09:18:58', NULL, 'EulerOS', 'openEuler release 20.03 (LTS-SP3-aarch64)', 'uefi', 'root', 'Cyclone@1234', '#version=openEuler release 20.03 (LTS-SP3)\n# Reboot after installation\nreboot\n# Use text mode install\ntext\n\n\n%pre --log=/tmp/kickstart_pre.log\n\necho "which dmidecode:"\nwhich dmidecode\n_sn=$(/usr/sbin/dmidecode -s system-serial-number 2>/dev/null | awk \'/^[^#]/ { print $1 }\')\necho "system-serial-number：$_sn"\necho "curl：http://osinstall.idcos.com/api/cloudboot/v1/devices/$_sn/installations/progress"\ncurl -H "Content-Type: application/json" -X POST -d "{\\"title\\":\\"启动OS安装程序\\",\\"progress\\":0.6,\\"log\\":\\"SW5zdGFsbCBPUwo=\\"}" http://osinstall.idcos.com/api/cloudboot/v1/devices/$_sn/installations/progress\ncurl -H "Content-Type: application/json" -X POST -d "{\\"title\\":\\"分区并安装软件包\\",\\"progress\\":0.7,\\"log\\":\\"SW5zdGFsbCBPUwo=\\"}" http://osinstall.idcos.com/api/cloudboot/v1/devices/$_sn/installations/progress\n\n%end\n\n\n%post --nochroot --log=/mnt/sysimage/tmp/kickstart_post_nochroot.log\n\necho "Copying %pre stage log files"\n/usr/bin/cp -rv /tmp/kickstart_pre.log /mnt/sysimage/tmp/\n\n%end\n\n\n%post --log=/tmp/kickstart_post.log\n#enable kdump\n#sed  -i "s/ crashkernel=512M/ crashkernel=1024M,high /" /boot/efi/EFI/openEuler/grub.cfg\n\nvirsh net-destroy default\n\necho "dowload commonSetting.py for setting items in OS(chroot):"\ncurl -o /tmp/commonSetting.py "http://osinstall.idcos.com/scripts/commonSetting.py"\n#curl -o /tmp/driver_upgrade_aarch64.sh "http://osinstall.idcos.com/scripts/driver_upgrade_aarch64.sh"\n#bash /tmp/driver_upgrade_aarch64.sh\n# config network\npython /tmp/commonSetting.py --network=Y\n\n# config osuser\npython /tmp/commonSetting.py --osuser=Y\n\n# complete\npython /tmp/commonSetting.py --complete=Y\n%end\n\n%packages --ignoremissing\n#@^minimal-environment\n@^server-product-environment\n@standard\n@system-tools\n#@development\n#@performance\n\n%end\n\n# Keyboard layouts\nkeyboard --vckeymap=us --xlayouts=\'us\'\n# System language\nlang en_US.UTF-8\n\n# Firewall configuration\nfirewall --disabled\n# Network information\nnetwork  --bootproto=dhcp --device=bootif --ipv6=auto --activate\nnetwork  --hostname=openeuler.webank\n\n# Use network installation\nurl --url="http://osinstall.idcos.com/openEuler/20.03_LTS_SP3/os/aarch64/"\ndriverdisk --source="http://osinstall.idcos.com/openEuler/20.03_LTS_SP3/drivers/aarch64/openEuler-20.03-LTS-SP3.iso"\n# System authorization information\nauth --enableshadow --passalgo=sha512\n# SELinux configuration\nselinux --disabled\n\n# Do not configure the X Window System\nskipx\n# System services\nservices --disabled="chronyd"\n\nignoredisk --only-use=sda\n# Partition clearing information\nclearpart --all --initlabel --drives=sda\n# Disk partitioning information\npart /boot/efi --fstype="efi" --ondisk=sda --size=200\npart /boot --fstype="ext4" --ondisk=sda --size=1024\npart / --fstype="ext4" --ondisk=sda --size=30720\npart /tmp --fstype="ext4" --ondisk=sda --size=10240\npart /home --fstype="ext4" --ondisk=sda --size=5120\npart /usr/local --fstype="ext4" --ondisk=sda --size=20480\npart /data --fstype="ext4" --size=1 --grow --ondisk=sda\n\n\n# System timezone\ntimezone Asia/Shanghai --utc\n\n# Root password\nrootpw Cyclone@1234\n\n%addon com_redhat_kdump --disable --reserve-mb=\'128\'\n\n%end', '#!ipxe\nkernel http://osinstall.idcos.com/openEuler/20.03_LTS_SP3/os/aarch64/images/pxeboot/vmlinuz initrd=initrd.img ksdevice=bootif ks=http://osinstall.idcos.com/api/cloudboot/v1/devices/{sn}/settings/system-template console=tty0 selinux=0 net.ifnames=0 biosdevname=0 BOOTIF=01-${netX/mac:hexhyp}\ninitrd http://osinstall.idcos.com/openEuler/20.03_LTS_SP3/os/aarch64/images/pxeboot/initrd.img\nboot', 'active', 'aarch64'),
	(32, '2022-09-07 08:57:15', '2023-12-07 09:18:40', NULL, 'EulerOS', 'openEuler release 20.03 (LTS-SP3-x86_64)', 'uefi', 'root', 'Cyclone@1234', '#version=openEuler release 20.03 (LTS-SP3)\n# Reboot after installation\nreboot\n# Use text mode install\ntext\n\n\n%pre --log=/tmp/kickstart_pre.log\n\necho "which dmidecode:"\nwhich dmidecode\n_sn=$(/usr/sbin/dmidecode -s system-serial-number 2>/dev/null | awk \'/^[^#]/ { print $1 }\')\necho "system-serial-number：$_sn"\necho "curl：http://osinstall.idcos.com/api/cloudboot/v1/devices/$_sn/installations/progress"\ncurl -H "Content-Type: application/json" -X POST -d "{\\"title\\":\\"启动OS安装程序\\",\\"progress\\":0.6,\\"log\\":\\"SW5zdGFsbCBPUwo=\\"}" http://osinstall.idcos.com/api/cloudboot/v1/devices/$_sn/installations/progress\ncurl -H "Content-Type: application/json" -X POST -d "{\\"title\\":\\"分区并安装软件包\\",\\"progress\\":0.7,\\"log\\":\\"SW5zdGFsbCBPUwo=\\"}" http://osinstall.idcos.com/api/cloudboot/v1/devices/$_sn/installations/progress\n\n%end\n\n\n%post --nochroot --log=/mnt/sysimage/tmp/kickstart_post_nochroot.log\n\necho "Copying %pre stage log files"\n/usr/bin/cp -rv /tmp/kickstart_pre.log /mnt/sysimage/tmp/\n\n%end\n\n\n%post --log=/tmp/kickstart_post.log\n#enable kdump\n#sed  -i "s/ crashkernel=512M/ crashkernel=1024M,high /" /boot/efi/EFI/openEuler/grub.cfg\n\necho "dowload commonSetting.py for setting items in OS(chroot):"\ncurl -o /tmp/commonSetting.py "http://osinstall.idcos.com/scripts/commonSetting.py"\n\ncurl -o /tmp/driver_upgrade.sh "http://osinstall.idcos.com/scripts/driver_upgrade.sh"\nbash /tmp/driver_upgrade.sh\n\n# config network\npython /tmp/commonSetting.py --network=Y\n\n# config osuser\npython /tmp/commonSetting.py --osuser=Y\n\n# complete\npython /tmp/commonSetting.py --complete=Y\n%end\n\n%packages --ignoremissing\n#@^minimal-environment\n@^server-product-environment\n@standard\n@system-tools\n#@development\n#@performance\n\n%end\n\n# Keyboard layouts\nkeyboard --vckeymap=us --xlayouts=\'us\'\n# System language\nlang en_US.UTF-8\n\n# Firewall configuration\nfirewall --disabled\n# Network information\nnetwork  --bootproto=dhcp --device=bootif --ipv6=auto --activate\nnetwork  --hostname=openeuler.webank\n\n# Use network installation\nurl --url="http://osinstall.idcos.com/openEuler/20.03_LTS_SP3/os/x86_64/"\n\n# System authorization information\nauth --enableshadow --passalgo=sha512\n# SELinux configuration\nselinux --disabled\n\n# Do not configure the X Window System\nskipx\n# System services\nservices --disabled="chronyd"\n\nignoredisk --only-use=sda\n# Partition clearing information\nclearpart --all --initlabel --drives=sda\n# Disk partitioning information\npart /boot/efi --fstype="efi" --ondisk=sda --size=200\npart /boot --fstype="ext4" --ondisk=sda --size=1024\npart / --fstype="ext4" --ondisk=sda --size=30720\npart /tmp --fstype="ext4" --ondisk=sda --size=10240\npart /home --fstype="ext4" --ondisk=sda --size=5120\npart /usr/local --fstype="ext4" --ondisk=sda --size=20480\npart /data --fstype="ext4" --size=1 --grow --ondisk=sda\n\n\n# System timezone\ntimezone Asia/Shanghai --utc\n\n# Root password\nrootpw Cyclone@1234\n\n%addon com_redhat_kdump --disable --reserve-mb=\'128\'\n\n%end', '#!ipxe\nkernel http://osinstall.idcos.com/openEuler/20.03_LTS_SP3/os/x86_64/images/pxeboot/vmlinuz initrd=initrd.img ksdevice=bootif ks=http://osinstall.idcos.com/api/cloudboot/v1/devices/{sn}/settings/system-template console=tty0 selinux=0 net.ifnames=0 biosdevname=0 BOOTIF=01-${netX/mac:hexhyp}\ninitrd http://osinstall.idcos.com/openEuler/20.03_LTS_SP3/os/x86_64/images/pxeboot/initrd.img\nboot', 'active_default', 'x86_64'),
	(33, '2023-03-16 16:05:19', '2023-03-16 16:05:19', NULL, 'BootOS', 'bootos_x86_64_intel', 'uefi', '', '', '#', '#!ipxe\nkernel http://osinstall.idcos.com/bootos/x86_64_intel/vmlinuz initrd=initrd.img console=tty0 selinux=0 biosdevname=0 SERVER_ADDR=http://osinstall.idcos.com LOOP_INTERVAL=5 DEVELOPER=1 BOOTIF=01-${netX/mac:hexhyp}\ninitrd http://osinstall.idcos.com/bootos/x86_64_intel/initrd.img\nboot', 'testing', 'x86_64'),
	(34, '2023-03-16 16:06:02', '2023-03-16 16:06:02', NULL, 'BootOS', 'bootos_x86_64_hygon', 'uefi', '', '', '#', '#!ipxe\nkernel http://osinstall.idcos.com/bootos/x86_64_hygon/vmlinuz initrd=initrd.img console=tty0 selinux=0 biosdevname=0 SERVER_ADDR=http://osinstall.idcos.com LOOP_INTERVAL=5 DEVELOPER=1 BOOTIF=01-${netX/mac:hexhyp}\ninitrd http://osinstall.idcos.com/bootos/x86_64_hygon/initrd.img\nboot', 'testing', 'x86_64'),
	(36, '2023-03-16 16:07:34', '2023-03-16 16:08:15', NULL, 'BootOS', 'bootos_default', 'uefi', '', '', '#', '#!ipxe\nkernel http://osinstall.idcos.com/bootos/x86_64/vmlinuz initrd=initrd.img console=tty0 selinux=0 biosdevname=0 SERVER_ADDR=http://osinstall.idcos.com LOOP_INTERVAL=5 DEVELOPER=1 BOOTIF=01-${netX/mac:hexhyp}\ninitrd http://osinstall.idcos.com/bootos/x86_64/initrd.img\nboot', 'testing', 'x86_64');
/*!40000 ALTER TABLE `system_template` ENABLE KEYS */;

-- 导出  表 cloudboot_cyclone.virtual_cabinet 结构
CREATE TABLE IF NOT EXISTS `virtual_cabinet` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL COMMENT '记录创建时间',
  `updated_at` timestamp NULL DEFAULT NULL COMMENT '记录修改时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '记录删除时间',
  `store_room_id` int(11) unsigned NOT NULL COMMENT '所属库房ID',
  `number` varchar(255) DEFAULT NULL COMMENT '编号',
  `status` enum('under_construction','not_enabled','enabled','offline') DEFAULT NULL COMMENT '状态。under_construction-建设中; not_enabled-未启用; enabled-已启用; offline-已下线;',
  `remark` varchar(1024) DEFAULT NULL COMMENT '备注',
  `creator` varchar(255) DEFAULT NULL COMMENT '记录创建者ID',
  `updater` varchar(255) DEFAULT NULL COMMENT '记录更新者ID',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_room_number` (`store_room_id`,`number`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='虚拟货架表';

-- 正在导出表  cloudboot_cyclone.virtual_cabinet 的数据：~0 rows (大约)
/*!40000 ALTER TABLE `virtual_cabinet` DISABLE KEYS */;
/*!40000 ALTER TABLE `virtual_cabinet` ENABLE KEYS */;

-- 导出  触发器 cloudboot_cyclone.trigger_update_device_oob 结构
SET @OLDTMP_SQL_MODE=@@SQL_MODE, SQL_MODE='';
DELIMITER //
CREATE TRIGGER trigger_update_device_oob BEFORE UPDATE ON device
FOR EACH ROW
BEGIN
    IF NEW.oob_user != OLD.oob_user OR NEW.oob_password != OLD.oob_password THEN
      INSERT INTO device_oob_history(created_at,updated_at,sn,username_old,username_new,password_old,password_new,creator) VALUES(NOW(),NOW(),NEW.sn,OLD.oob_user,NEW.oob_user,OLD.oob_password,NEW.oob_password,NEW.creator);
    END IF;
END//
DELIMITER ;
SET SQL_MODE=@OLDTMP_SQL_MODE;

-- 导出  视图 cloudboot_cyclone.cabinet_power_info 结构
-- 移除临时表并创建最终视图结构
DROP TABLE IF EXISTS `cabinet_power_info`;
CREATE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `cabinet_power_info` AS select `sc`.`id` AS `id`,`sc`.`number` AS `number`,`sc`.`is_enabled` AS `is_enabled`,`sc`.`is_powered` AS `is_powered`,`i`.`name` AS `idc`,`sr`.`name` AS `server_room`,`na`.`name` AS `network_area`,`sc`.`max_power` AS `max_power`,(case when isnull(`su`.`usite_total`) then 0 else `su`.`usite_total` end) AS `usite_total`,(case when isnull(`su`.`used_count`) then 0 else `su`.`used_count` end) AS `used_count`,(case when isnull(`su`.`free_count`) then 0 else `su`.`free_count` end) AS `free_count`,(case when isnull(`su`.`pre_occupied_count`) then 0 else `su`.`pre_occupied_count` end) AS `pre_occupied_count`,(case when isnull(`su`.`disabled_count`) then 0 else `su`.`disabled_count` end) AS `disabled_count`,(case when isnull(`cp`.`known_power`) then 0 else `cp`.`known_power` end) AS `known_used_power`,(case when isnull(`cp`.`known_power`) then `sc`.`max_power` else (`sc`.`max_power` - `cp`.`known_power`) end) AS `free_power`,(case when isnull(`cp`.`is_unknown_count`) then 0 else `cp`.`is_unknown_count` end) AS `is_unknown_power_svr_count` from (((((`cloudboot_cyclone`.`server_cabinet` `sc` left join (select `d`.`server_cabinet_id` AS `server_cabinet_id`,sum((case when ((`dc`.`power` = 'unknown') or isnull(`dc`.`power`)) then 0 when (locate('W',`dc`.`power`) > 0) then cast(substring_index(trim(`dc`.`power`),'W',1) as signed) else cast(trim(`dc`.`power`) as signed) end)) AS `known_power`,sum((case when ((`dc`.`power` = 'unknown') or isnull(`dc`.`power`)) then 1 else 0 end)) AS `is_unknown_count` from (`cloudboot_cyclone`.`device` `d` left join `cloudboot_cyclone`.`device_category` `dc` on((`d`.`category` = `dc`.`category`))) group by `d`.`server_cabinet_id`) `cp` on((`sc`.`id` = `cp`.`server_cabinet_id`))) left join `cloudboot_cyclone`.`idc` `i` on((`sc`.`idc_id` = `i`.`id`))) left join `cloudboot_cyclone`.`server_room` `sr` on((`sc`.`server_room_id` = `sr`.`id`))) left join `cloudboot_cyclone`.`network_area` `na` on((`sc`.`network_area_id` = `na`.`id`))) left join (select `u`.`server_cabinet_id` AS `server_cabinet_id`,count(1) AS `usite_total`,count(if((`u`.`status` = 'used'),TRUE,NULL)) AS `used_count`,count(if((`u`.`status` = 'free'),TRUE,NULL)) AS `free_count`,count(if((`u`.`status` = 'pre_occupied'),TRUE,NULL)) AS `pre_occupied_count`,count(if((`u`.`status` = 'disabled'),TRUE,NULL)) AS `disabled_count` from `cloudboot_cyclone`.`server_usite` `u` group by `u`.`server_cabinet_id`) `su` on((`su`.`server_cabinet_id` = `sc`.`id`)));

/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IF(@OLD_FOREIGN_KEY_CHECKS IS NULL, 1, @OLD_FOREIGN_KEY_CHECKS) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
