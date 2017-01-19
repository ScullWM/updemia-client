## UPDEMIA CLIENT
Tired of dropbox blocking preview on slack or facebook messenger ? Send your screenshot directly on updemia.com (no limitation, no ads) and preview on IM !
http://www.updemia.com


## HOW TO INSTALL UPDEMIA CLIENT

### First step:
Define to macosx screencapture destination path. In our case, it's "/tmp/updemia".

Set path and restart macosx screenshot tools
```
defaults write com.apple.screencapture location /tmp/updemia
killall SystemUIServer
```

### Second step:
Get and install client binary:
```
curl -L https://github.com/ScullWM/updemia-client/releases/download/v0.2/updemia-client -o updemia-client
chmod a+x updemia-client
mv updemia-client /usr/local/bin/updemia-client
```
(feel free to change your binaries dedicated path)

Start the bin and define the screencapture destination path
```
./updemia-client /tmp/updemia
```

## One-line installation (with `curl`)

Just run this installation script from your favorite Terminal:

```
curl -sS https://raw.github.com/ScullWM/updemia-client/master/install.sh | sh
```
