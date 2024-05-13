# somfy-rts-gateway

A Somfy RTS Gateway using Signalduino

### Find USB device
```shell
ls -l /dev/serial/by-id 

total 0
lrwxrwxrwx 1 root root 13 Oct 24 07:49 usb-1a86_USB2.0-Serial-if00-port0 -> ../../ttyUSB1
lrwxrwxrwx 1 root root 13 Dec 22 13:02 usb-SHK_SIGNALduino_868-if00-port0 -> ../../ttyUSB0
```

If needed remove the modemmanager
```shell
sudo apt-get purge modemmanager
```


### compile for raspberry pi

##### linux
```shell 
GOOS=linux GOARCH=arm GOARM=5 CGO_ENABLED=0 go build -mod=vendor -o somfy-rts-gateway
```

##### windows
```cmd
set GOOS=linux
set GOARCH=arm
set GOARM=6
set CGO_ENABLED=0
go build -mod=vendor -o somfy-rts-gateway
```


### copy with scp

```shell
scp -r ./arm/* pi@192.168.66.139:/opt/somfy
```

### install as service

Copy this file into `/etc/systemd/system` as root, for example:

```shell
sudo cp myscript.service /etc/systemd/system/myscript.service
```
Once this has been copied, you can attempt to start the service using the following command:

```shell
sudo systemctl start myscript.service
```
Stop it using following command:

```shell
sudo systemctl stop myscript.service
```
When you are happy that this starts and stops your app, you can have it start automatically on reboot by using this command:

```shell
sudo systemctl enable myscript.service
```