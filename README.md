# PhoneWhisper
Voice chat for a tour group using phones and a travel router without needing internet.

This is a modern take on offline voice chat in 2025 that does not require any apps, has no account signups, and can be instantly setup by scanning a qr code to connect to wifi and another to open a website.

You no longer need to rent a $500 set of whisper devices that can be lost or pay daily fees to LiveTours. 

You can use your own wireless headphones like airpods and don't need to worry about hygiene issues.

Connections are established with websockets and voice is transmitted via encrypted WebRTC all locally. 

All phones can be in airplane mode and receive voice from the wifi on an travel router without any internet.

I host this on a $120 `glinet slate ax`. You can also use the cheaper $30 `glinet opal`. I built and tested this on a faster router because I needed to support 50+ users. GliNet advertises the opal will work with a max of 52 users and the slate will handle 120. Advertised speeds aside. When calculating pi, a fairly memory intensive operation as well, the opal took about 68x longer, for phonewhisper I don't think this will matter much. The more important features are likely range and network bandwidth. I was able to push 500mbps over wifi on the opal and line of sight you won't have range issues. While untested in a group setting, it'll likely functionally work the same for large groups at a much cheaper price. 

| Product|  Operating System | Memory| Processor            | Cores | Clock Speed | Time to Calculate 10,000,000 Digits of Pi |
|---------|-----|----|---|-------|-------------|-------------------------------|
| GLiNet Opal | OpenWRT 18.06 | 118,784 kB / 17,844 kB free | sf19a28  MIPSLE      | 4     | 800 MHz      | 7m 0.9355s        |
| GLiNet Slate AX | OpenWRT 21.02 | 402,300 kB / 148,752kB free| Qualcomm IPQ6000 ARMv7| 4     | 1.2 GHz      | 0m 6.1206s        |

You can likely use other openwrt devices as well. The program is written in golang, if you don't mind recompiling or modifying the installer it should be able to run just about anywhere.

# Other Uses

This software really opens the door to a lot of things and is not just limited to tour guides. You can make you own silent disco. Broadcast audio at at church. Broadcast audio at a tailgate. Stream audio from a tv at the gym. Stream audio around your house. Build your own intercom for your house. Make your own walkie talkies for a roadtrip (granted installing umurmur would work better). You could run this software on a pair of ubiquiti nanostation locos and for $100 get a 6 mile line of sight audio link. There's a certain lost beauty in today's commercialized and internet connected day and age that this software is somehow still immune and able to embody.

# Development
I searched the internet for a while and couldn't find any software that does what I needed. I originally went down the path of using icecast but ran into large delays even after removing the buffer it was still around 5s. I tried mumble as well, but mid compiling the mumble webrtc proxy from 5 years ago and running into dependency issues I figured it would be easier to just write something purpose built myself. I almost modified umurmur to only enable broadcasting by certain users, but decided it's going to be a pain to get users to install this unfortuantely outdated and clumsy mumble app.

# Installation
## Connect your travel router via ethernet
Connect your travel router to the internet connecting the wan port to a port on your home router

## ssh into your device
from another computer
```
ssh root@192.168.8.1
```
if you have an opal
```
ssh -oHostKeyAlgorithms=+ssh-rsa root@192.168.8.1
```
password is whatever you used during setup, default is goodlife

## installer
```
curl https://raw.githubusercontent.com/mkunz7/phonewhisper/refs/heads/main/install.sh | sh
```
that's it, this only takes seconds to install, if you see the word done it worked

You can now disconnect the router from the internet it won't be needed anymore
# Usage
Connect to the wifi of your travel router. I recommend making a qr code using https://qifi.org.

Connect to http://192.168.8.1:3000 to listen to a stream, qr codes can be made using https://qr.15c.me/qr.html

![listen](https://github.com/user-attachments/assets/ff8ea596-4f23-4b58-be04-8627b151dc41)


Connect to https://192.168.8.1:3001 to broadcast, qr codes can be made using https://qr.15c.me/qr.html

![broadcast](https://github.com/user-attachments/assets/e0fc41ab-1da0-40a3-888b-e165d10d25fb)

- You can save the website on your home screen using the share button to look like an app
- Safari is recommended when using an iphone
- Broadcasting does not currently support multiple simultaneous broadcasters, the latest takes over, one person should disconnect by refreshing to prevent issues
- Certificates are required for microphone access on modern browsers, hence the invalid warning
- The broadcaster will need to accept the invalid certificate, the one I've included is good for 100 years. It was generated using `openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 36500 -nodes`
- This should work when the screen locks (YMMV), the broadcaster should probably have their screen set to never sleep, do some testing on your own.
- Headphones / Bluetooth devices are recommended to prevent feedback loops
- This should run for 12+ hours on a 20,000mah battery brick with usbc
- I used a Thousover USBC to 3.5mm+charger cable https://www.amazon.com/Headphone-Charger-Adapter-Charging-Compatible/dp/B09QJYVF68 on a cheap android burner phone with a wired headphone and got charging and voice to work at the same time. YMMV with other branded OTG cables.
  
# Tested Devices

| Device              | OS Version | Capabilities       |
|---------------------|------------|--------------------|
| iPod Touch 5th Gen  | 12.5.5     | Listen only        |
| iPhone 5s           | 12.5.7     | Listen only        |
| iPhone 6            | 12.5.5     | Listen only        |
| iPhone 7 Plus       | 14.2       | Listen only        |
| iPhone 7 Plus       | 15.0.1     | Listen and Broadcast |
| iPhone 8           | 16.7.8     | Listen and Broadcast |
| iPhone 11           | 18.1.1     | Listen and Broadcast |
| iPhone 14 Pro          | 18.1.1     | Listen and Broadcast |
| iPhone 16 Pro          | 18.2.1     | Listen and Broadcast |
| Blu View 2          | Android 10 | Listen and Broadcast |
| Windows Chrome      | 10         | Listen and Broadcast |
| M1 Macbook Air Chrome     | Sequoia         | Listen and Broadcast |

Generally speaking if you have an iphone <10 years old it should be fine with updates. WebRTC saw apple adoption in Sept 2021. WebRTC reached recommendation status Jan 2021. Older androids will work with updated chrome browsers.

Albeit too may older devices may impact how many users you can have connect at once. 

Do a PR if you have a testimony you would like to share.

# Bandwidth Usage
I setup a broadcaster listening to a mr beast video and opened 50 listeners on my laptop connected via wifi to simulate 50 connected phones to the opal.

The command line utility iftop for wlan1 is saying the network is only using 5.53mbit (691kB/s) for 50 users which is next to nothing.

I didn't believe the numbers so I setup a wireshark capture to record all network traffic and I also spot checked with menumeters on my mac.

Wireshark:
- 10 users generated 41.7k packets in 60s, Wireshark Capture size 7mB
- 50 users generated 205.9k packets in 60s, Wireshark Capture size 35mB

Most the packets are UDP for webrtc.

Menu meters:
- 1 user TX Rate = 1.9kB/s, RX Rate = 8.0 kB/s
- 5 users TX Rate = 8.1 kB/s, RX Rate = 39.2 kB/s
- 10 users TX Rate = 15.3 kB/s, RX Rate = 80 kB/s
- 50 users TX Rate = 72.3 kB/s, RX Rate = 386.4 kB/s

iftop isn't lying, bandwidth wise the opal should be able to support 5,000 streams no problem. 

In the real world having that many phones connect to one wifi network is a different story.

# Screenshots
![IMG_3717](https://github.com/user-attachments/assets/644e9416-6e4b-432e-aecf-8a7b7f875ebe)
![IMG_0013](https://github.com/user-attachments/assets/19358e61-1377-4015-819d-404e2d4b458d)


