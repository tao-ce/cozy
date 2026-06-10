---
icon: fontawesome/brands/raspberry-pi
---

# Raspberry Pi 4

!!! abstract inline end "About this guide"
    In this guide, we will explain how to download and start an Appliance to run TAO Community Edition on Raspberry Pi 4.

## Requirements

To achieve this deployment, we will need:

* a Raspberry Pi 4, with at least 8GB of RAM
* a microSD card, at least 64GB, and supporting class U1 or higher (class U3 is prefered)
* a microSD reader
* appropriate USB-C power supply (5V, ≥3A)
* (optional) an USB keyboard, and a HDMI screen with micro HDMI cable
* reliable Ethernet network, with high-speed internet connection
* a computer running recent OS, with at least 20GB of free space
* image etcher tool (e.g. `dd` on Linux, [Rufus](https://rufus.ie) on Windows)
* compression software supporting `xz` format (`xz` on Linux, [7-Zip](https://www.7-zip.org/) on Windows)

## Download Appliance image

<!-- TODO get more definitive image download link -->

From your computer, open the following link and start downloading [`disk.raw.xz`](https://github.com/tao-ce/cozy/actions/runs/27224360057/artifacts/7515755796) (~1.2GB), and uncompress it.

=== "On Linux"
    Download and uncompress in a single line:
    ```bash
    curl -Ls https://github.com/tao-ce/cozy/actions/runs/27224360057/artifacts/7515755796 | xz -d
    ```


=== "On Windows"
    * Download [`disk.raw.xz`](https://github.com/tao-ce/cozy/actions/runs/27224360057/artifacts/7515755796)
    * Uncompress it, you should obtain `disk.raw` file


## Etch microSD card 

We will now load this image on your microSD card.

!!! danger inline end "Risk of data loss"
    Please note this operation will completly remove all data on microSD card.
    
    Also, depending on your system, pay attention to carefully identify target device, to avoid damaging other storage on your computer.

=== "On Linux"
    1. First, plug microSD reader with the card in it. Your system may prompt to mount and explore it. Ignore it.
    2. Ensure to identify your microSD drive, exemple with `lsblk`, we can find it is called `/dev/mmcblk0` on this computer.
    ```
    # lsblk -pd \
      -o NAME,SIZE,VENDOR,STATE,MODEL,SERIAL \
      --filter 'TYPE=="disk" and HOTPLUG==1'
    NAME          SIZE VENDOR STATE MODEL SERIAL
    /dev/mmcblk0 58.2G                    0x000006fd
    ```
    3. Find `disk.raw` image and run the following command (ensure to change paths to match your device):
    ```bash
    sudo dd if=path/to/disk.raw of=/dev/mmcblk0 bs=1M status=progress
    ```

=== "On Windows"
    1. First, plug microSD reader with the card in it. Your system may prompt to mount and explore it. Ignore it.
    2. Open your image etcher tool (e.g. Rufus)
    3. In `Device`, look for the microSD card
    4. In `Boot selection`, open `disk.raw`
    5. Proceed

Depending your hardware, this operation can take 10 to 30 minutes.

Once done, ensure to unmount your microSD, and unplug it.

## Start Raspberry Pi 4 with Appliance image

Now that the image has been copied on microSD card, you can insert it in Raspberry Pi 4.

If you have keyboard and screen, you may plug them to Raspberry Pi 4.

Prepare USB power supply, and plug it.

In few minutes, you should see some logs on your screen.

!!! notes "Long process"
    Depending on your hardware, it can takes several dozens of minutes before Raspberry Pi 4 can be reachable.

    You may monitor logs on console, and log in a terminal (default credentials are `tao`/`tao`), however most of the installation is happening in background.

!!! tip "Under construction"
    Appliance image is still in active development, and lack some readiness probe to let know the users that system can be used.

    However, after few minutes, you may be able to access Cockpit web console on [`https://tao-community-edition.local:9090`](https://tao-community-edition.local:9090). Use username `tao` and password `tao` to log in.

!!! warning "Change password"
    Whenever you can, log in on local console or web console to change default password.

## Access TAO Community Edition


??? note inline end "About TAO CE address"
    Note the address is different than usual container deployment (`community.tao.internal`), as we rely on mDNS protocol to propagate appliance addreses on your network.

Once ready, you should be able to connect to [`https://tao-community-edition.local`](https://tao-community-edition.local). 

You may now follow TAO Community Edition [first steps](https://tao-ce.github.io/tao-ce/start/first/) and [user guide](https://userguide.taotesting.com/).