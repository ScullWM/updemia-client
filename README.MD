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
Start the bin and define the screencapture destination path
```
chmod +x ./updemia-client
./updemia-client /tmp/updemia
```