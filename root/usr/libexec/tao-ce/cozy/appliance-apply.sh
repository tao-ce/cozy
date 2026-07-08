#!/usr/bin/bash
#

COZY_CONFIG_PATH=/etc/tao-ce/cozy/appliance.json
TAO_CE_CONFIG_PATH=/etc/tao-ce/config/tao.yaml

[ -f $COZY_CONFIG_PATH ] || {
    echo "Appliance configuration file not found"
    exit 1
}

HOTSPOT_NETWORK_IF=wlan0
HOTSPOT_NETWORK_NAME=wifi-${HOTSPOT_NETWORK_IF}
HOTSPOT_NMCONNECTION_PATH=/etc/NetworkManager/system-connections/${HOTSPOT_NETWORK_NAME}.nmconnection

setup_hotspot(){
    HOTSPOT_SSID=$(jq -r .hotspot.ssid $COZY_CONFIG_PATH)
    HOTSPOT_PASSWORD=$(jq -r .hotspot.password $COZY_CONFIG_PATH)
    HOTSPOT_SECURITY=$(jq -r .hotspot.security $COZY_CONFIG_PATH)
    [ "$HOTSPOT_SECURITY" = "open" ] && SECURITY=none
    [ "$HOTSPOT_SECURITY" = "open" ] && PSK=""
    [ "$HOTSPOT_SECURITY" = "wpa2" ] && SECURITY=wpa-psk
    [ "$HOTSPOT_SECURITY" = "wpa2" ] && PSK="psk=${HOTSPOT_PASSWORD}"
    cat <<EOF > $HOTSPOT_NMCONNECTION_PATH
[connection]
id=${HOTSPOT_NETWORK_NAME}
interface-name=${HOTSPOT_NETWORK_IF}
type=wifi
autoconnect=yes

[wifi]
mode=ap
ssid=${HOTSPOT_SSID}
[wifi-sec]
pmf=1
key-mgmt=${SECURITY}
${PSK}

[ipv4]
address1=10.74.0.254/24
method=shared

[ipv6]
method=ignore
EOF
    chmod 600 $HOTSPOT_NMCONNECTION_PATH
    nmcli connection reload
    nmcli connection up ${HOTSPOT_NETWORK_NAME}
    systemctl restart NetworkManager.service
}

echo "Enabling/disabling DDNS"
jq -e .features.ddns $COZY_CONFIG_PATH \
  && systemctl enable avahi-daemon.service \
  || systemctl disable avahi-daemon.service

echo "Setting FQDN"
sed -i \
"s/tao-community-edition.local/$(jq -r .taoCe.fqdn $COZY_CONFIG_PATH)/g" \
  /etc/hosts \
  /etc/firefox/policies/policies.json \
  /etc/tao-ce/config/tao.yaml \
  /etc/skel/Desktop/tao-portal.desktop

hostnamectl set-hostname $(jq -r .taoCe.fqdn $COZY_CONFIG_PATH)

echo "Setting up containers"
mkdir -p /etc/containers/systemd/tao-ce.container.d
jq \
  --slurpfile manifest $COZY_CONFIG_PATH \
  -sRr 'gsub("[$][{](?<p>.*)[}]"; .p as $p| $manifest[0]|getpath($p|split(".")) )' \
  /usr/libexec/tao-ce/cozy/templates/tao-ce.container.override.conf \
  > /etc/containers/systemd/tao-ce.container.d/override.conf

echo "Reloading systemd"
systemctl daemon-reload

echo "Setting timezone"
timedatectl set-timezone $(jq -r '.locale.timezoneRegion + "/" + .locale.timezoneCity' $COZY_CONFIG_PATH)

echo "Setting locale"
localectl set-locale $(jq -r '.locale.language' $COZY_CONFIG_PATH)

echo "Setting up hotspot"
jq -e .features.hotspot $COZY_CONFIG_PATH && setup_hotspot

echo "Starting TAO CE"
systemctl start tao-ce --no-block

echo "Done"

exit 0

