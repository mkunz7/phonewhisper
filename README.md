# PhoneWhisper
Voice chat for a tour guide using a phone and a travel router without needing internet.

A modern take on offline voice chat in 2025.

You no longer need to rent a $500 set of whisper devices that can be lost or pay daily fees to LiveTours.

Connections are established with websockets and voice is transmitted via WebRTC all locally. All phones can be in airplane mode and receive voice from the wifi on an travel router without internet.

I host this on a `glinet slate ax` you can likely use the cheaper `glinet opal` as well if you don't need to support 50+ users. These are armv7 openwrt devices, others will likley work if you don't mind recompiling.

I searched the internet for a while and couldn't find any software that does this. I originally went down the path of using icecast but ran into 5s delays. I tried mumble as well, but mid compiling the mumble webrtc proxy from 5 years and running into dependency issues I figured it would be easier to just write something purpose built myself.

# Installation

```
ssh root@192.168.8.1
curl https://ku.nz/blog/files/voice-server/install.sh | sh
```

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
- This should work when the screen locks (YMMV)
- Headphones / Bluetooth devices are recommended to prevent feedback loops
- This should run for 12+ hours on a 20,000mah battery brick with usbc
- I used a Thousover USBC to 3.5mm+charger cable https://www.amazon.com/Headphone-Charger-Adapter-Charging-Compatible/dp/B09QJYVF68 with a wired headphone and got charging and voice to work at the same time. YMMV with other branded OTG cables.
  
# Tested Devices

| Device              | OS Version | Capabilities       |
|---------------------|------------|--------------------|
| iPod Touch 5th Gen  | 12.5.5     | Listen only        |
| iPhone 5s           | 12.5.7     | Listen only        |
| iPhone 6            | 12.5.5     | Listen only        |
| iPhone 7 Plus       | 14.2       | Listen only        |
| iPhone 7 Plus       | 15.0.1     | Listen and Broadcast |
| iPhone 11           | 18.1.1     | Listen and Broadcast |
| iPhone 14           | 18.1.1     | Listen and Broadcast |
| iPhone 16           | 18.1.1     | Listen and Broadcast |
| Blu View 2          | Android 10 | Listen and Broadcast |
| Windows Chrome      | 10         | Listen and Broadcast |
| M1 Macbook Air Chrome     | 15         | Listen and Broadcast |

# Bandwidth Usage
iftop registers the wlan1 interface using 2mb for 5 devices during a voice broadcast, extrapolating 50 should only use 100mbps

# Screenshots
![IMG_3717](https://github.com/user-attachments/assets/644e9416-6e4b-432e-aecf-8a7b7f875ebe)
![IMG_0013](https://github.com/user-attachments/assets/19358e61-1377-4015-819d-404e2d4b458d)


