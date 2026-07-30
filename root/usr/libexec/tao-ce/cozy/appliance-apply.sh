#!/usr/bin/bash -ex
#

mkdir -p $(dirname $COZY_LOG_PATH)
exec >$COZY_LOG_PATH 2>&1

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
    mkdir -p $(dirname $HOTSPOT_NMCONNECTION_PATH)
    cat <<EOF > $HOTSPOT_NMCONNECTION_PATH
[connection]
id=${HOTSPOT_NETWORK_NAME}
interface-name=${HOTSPOT_NETWORK_IF}
type=802-11-wireless
autoconnect=yes

[802-11-wireless]
mode=ap
ssid=${HOTSPOT_SSID}
band=bg

[802-11-wireless-security]
pmf=2
key-mgmt=${SECURITY}
${PSK}

[ipv4]
address1=10.74.0.254/24
method=shared

[ipv6]
method=ignore
EOF
    chmod 600 $HOTSPOT_NMCONNECTION_PATH
}

echo "Setting FQDN"
#TODO: Set FQDN in /etc/firefox/policies/policies.json

# Update /etc/hostname
jq -r .taoCe.fqdn $COZY_CONFIG_PATH >/etc/hostname

# Update /etc/hosts
sed -Ei \
  -e '/^0.0.0.0/d' \
  -e "1i0.0.0.0 $(jq -r .taoCe.fqdn $COZY_CONFIG_PATH)" \
  /etc/hosts
sed -Ei \
  -e '/^host-name=/d' \
  -e "/\[server\]/ahost-name=$(jq -r '.taoCe.fqdn|split(".")|.[0]' $COZY_CONFIG_PATH)" \
  /etc/avahi/avahi-daemon.conf

# Update /etc/skel/Desktop/tao-portal.desktop
sed -Ei \
  's@URL=.*@URL=https://'$(jq -r .taoCe.fqdn $COZY_CONFIG_PATH)'@' \
  /etc/skel/Desktop/tao-portal.desktop

tao_ce_tmp=$(mktemp)

cat $TAO_CE_CONFIG_PATH \
  | python3 -c 'import sys, yaml, json; json.dump(yaml.safe_load(sys.stdin), sys.stdout, indent=2)' \
  | jq --slurpfile manifest $COZY_CONFIG_PATH \
    '.
      | .spec.publicDomain = ($manifest[0].taoCe.fqdn)
      | .spec.defaultLocale = ($manifest[0].locale.language|split(".")|.[0]|gsub("[_]";"-"))
      ' \
  > $tao_ce_tmp

python3 -c 'import sys, yaml, json; yaml.dump(json.load(sys.stdin), sys.stdout, indent=2)' \
  < $tao_ce_tmp \
  > $TAO_CE_CONFIG_PATH

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
 ln -sf \
   /usr/share/zoneinfo/$(jq -r '.locale.timezoneRegion + "/" + .locale.timezoneCity' $COZY_CONFIG_PATH) \
   /etc/localtime 

echo "Setting locale"

jq -r \
  '"LANG=" + .locale.language' \
  < $COZY_CONFIG_PATH \
  >/etc/locale.conf

echo "Setting up hotspot"
jq -e .features.hotspot $COZY_CONFIG_PATH && setup_hotspot

echo "Starting TAO CE"
systemctl start tao-ce --no-block

echo "Enabling/disabling MulticastDNS"
jq -e .features.mdns $COZY_CONFIG_PATH \
  && { systemctl enable avahi-daemon.service && systemctl restart avahi-daemon.service --no-block ;} \
  || systemctl disable avahi-daemon.service

echo "Done"

exit 0

