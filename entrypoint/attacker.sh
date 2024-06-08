#!/bin/sh

echo "bootstart attacker"

if [ "$BOOT_DELAY" != "" ] &&  [ "$BOOT_DELAY" != "0" ];then
	        sleep $BOOT_DELAY
fi

/usr/bin/attacker > /root/attackerdata/d.log 2>&1
