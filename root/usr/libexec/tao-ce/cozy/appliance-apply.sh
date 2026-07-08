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
    SSID=$(jq -r .hotspot.ssid $COZY_CONFIG_PATH)
    PASSWORD=$(jq -r .hotspot.password $COZY_CONFIG_PATH)
    SECURITY=$(jq -r .hotspot.security $COZY_CONFIG_PATH)
    [ "$SECURITY" = "open" ] && SECURITY=none
    [ "$SECURITY" = "open" ] && PSK=""
    [ "$SECURITY" = "wpa2" ] && SECURITY=wpa-psk
    [ "$SECURITY" = "wpa2" ] && PSK="psk=${PASSWORD}"
    cat <<EOF > $HOTSPOT_NMCONNECTION_PATH
[connection]
id=${HOTSPOT_NETWORK_NAME}
interface-name=${HOTSPOT_NETWORK_IF}
type=wifi
autoconnect=yes

[wifi]
mode=ap
ssid=${SSID}
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

jq -e .features.hotspot $COZY_CONFIG_PATH && setup_hotspot
jq -e .features.ddns $COZY_CONFIG_PATH \
  && systemctl enable --now avahi-daemon.service \
  || systemctl disable --now avahi-daemon.service

sed -i \
"s/tao-community-edition.local/$(jq -r .taoCe.fqdn $COZY_CONFIG_PATH)/g" \
  /etc/hosts \
  /etc/firefox/policies/policies.json \
  /etc/tao-ce/config/tao.yaml \
  /etc/skel/Desktop/tao-portal.desktop

hostnamectl set-hostname $(jq -r .taoCe.fqdn $COZY_CONFIG_PATH)

mkdir -p /etc/containers/systemd/tao-ce.container.d
jq \
  --slurpfile manifest $COZY_CONFIG_PATH \
  -sRr 'gsub("[$][{](?<p>.*)[}]"; .p as $p| $manifest[0]|getpath($p|split(".")) )' \
  /usr/libexec/tao-ce/cozy/templates/tao-ce.container.override.conf \
  > /etc/containers/systemd/tao-ce.container.d/override.conf

timedatectl set-timezone $(jq -r '.locale.timezoneRegion + "/" + .locale.timezoneCity' $COZY_CONFIG_PATH)
localectl set-locale $(jq -r '.locale.language' $COZY_CONFIG_PATH)

systemctl daemon-reload

systemctl restart tao-ce.service --no-block

exit 0

