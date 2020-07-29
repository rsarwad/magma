#!/bin/bash
# Setting up env variable, user and project path
ERROR=""
INFO=""
SUCCESS_MESSAGE="ok"

addError() {
    ERROR="$ERROR\n$1  to fix it: $2"
}

addInfo() {
    INFO="$INFO $1 \n"
}

if ! grep -q 'Debian' /etc/issue; then
  addError "Debian is not installed" "Restart installation following agw_install.sh, agw has to run on Debian"
fi

if ! grep -q "$MAGMA_USER ALL=(ALL) NOPASSWD:ALL" /etc/sudoers; then
    addError "Debian is not installed" "Restart installation following agw_install.sh, magma has to be sudoer"
fi

KVERS=$(uname -r)
if [ "$KVERS" != "4.9.0-9-amd64" ]; then
    addError "Kernel version is not 4.9.0-9-amd64" "Restart installation following agw_install.sh, KVERS has to be 4.9.0-9-amd64"
fi

interfaces=("eth1" "eth0" "gtp_br0")
for interface in "${interfaces[@]}"; do
    OPERSTATE_LOCATION="/sys/class/net/$interface/operstate"
    if test -f "$OPERSTATE_LOCATION"; then
        OPERSTATE=$(cat "$OPERSTATE_LOCATION")
        if [[ $OPERSTATE == 'down'  ]]; then
            addError "$interface is not configured" "Try to ifup $interface"
        fi
    else
        addError "$interface is not configured" "Check that /etc/network/interfaces.d/$interface has been set up"
    fi
done

PING_RESULT=$(ping -c 1 -I eth0 8.8.8.8 > /dev/null 2>&1 && echo "$SUCCESS_MESSAGE")
if [ "$PING_RESULT" != "$SUCCESS_MESSAGE" ]; then
    addError "eth0 is connected to the internet" "Make sure the hardware has been properly plugged in (eth0 to internet)"
fi

allServices=("control_proxy" "directoryd" "dnsd" "enodebd" "magmad" "mme" "mobilityd" "pipelined" "policydb" "redis" "sctpd" "sessiond" "state" "subscriberdb")
for service in "${allServices[@]}"; do
    if ! systemctl is-active --quiet "magma@$service"; then
        addError "$service is not running" "Please check our faq"
    fi
done

packages=("magma" "magma-cpp-redis" "magma-libfluid" "oai-gtp" "libopenvswitch" "openvswitch-datapath-dkms" "openvswitch-datapath-source" "openvswitch-common" "openvswitch-switch")
for package in "${packages[@]}"; do
    PACKAGE_INSTALLED=$(dpkg-query -W -f='${Status}' $package  > /dev/null 2>&1 && echo "$SUCCESS_MESSAGE")
    if [ "$PACKAGE_INSTALLED" != "$SUCCESS_MESSAGE" ]; then
        addError "$package hasn't been installed" "Rerun the agw_install.sh"
    fi
done

if [ -z "$ERROR" ]; then
    echo "Installation went smoothly, please let us know what went wrong/good on github"
else
    echo "There was a few errors during installation"
    printf "%s" "$ERROR"
fi

apt-get update > /dev/null
addInfo "$(apt list -qq --upgradable 2> /dev/null)"

if [ -n "$INFO" ]; then
    echo "INFO:"
    printf "%s" "$INFO"
fi
