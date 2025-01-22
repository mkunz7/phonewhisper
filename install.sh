#!/bin/sh
# cleanup
if [ -f /usr/share/voice-server/voiceserver ]; then
    echo "cleaning up previous install"
    /etc/init.d/voice-server stop 2>&1
    /etc/init.d/voice-server disable 2>&1
    rm -rf /usr/share/voice-server 
    rm -f /etc/init.d/voice-server
fi

mkdir -p /usr/share/voice-server
mkdir -p /usr/share/voice-server/public

ARCH=$(uname -m)
if echo "$ARCH" | grep -q "mips"; then
    echo "System is running on MIPS."
    curl -L https://github.com/mkunz7/phonewhisper/releases/download/main/voice-server-linux-mips -o /usr/share/voice-server/voiceserver

# Check for ARMv7
elif echo "$ARCH" | grep -q "armv7"; then
    echo "System is running on ARMv7."
    curl -L https://github.com/mkunz7/phonewhisper/releases/download/main/voice-server-linux-arm7 -o /usr/share/voice-server/voiceserver

else
    echo "Architecture not currently supported by installer: $ARCH"
    exit 1
fi

curl https://raw.githubusercontent.com/mkunz7/phonewhisper/refs/heads/main/voice-server -o /etc/init.d/voice-server
curl https://raw.githubusercontent.com/mkunz7/phonewhisper/refs/heads/main/public/listen.html -o /usr/share/voice-server/public/listen.html
curl https://raw.githubusercontent.com/mkunz7/phonewhisper/refs/heads/main/public/broadcast.html -o /usr/share/voice-server/public/broadcast.html

#openssl req -x509 -newkey rsa:4096 -keyout /usr/share/voice-server/key.pem -out /usr/share/voice-server/cert.pem -days 36500 -nodes -subj "/CN=localhost"
curl https://raw.githubusercontent.com/mkunz7/phonewhisper/refs/heads/main/cert.pem -o /usr/share/voice-server/cert.pem
curl https://raw.githubusercontent.com/mkunz7/phonewhisper/refs/heads/main/key.pem -o /usr/share/voice-server/key.pem

chmod +x /usr/share/voice-server/voiceserver
chmod +x /etc/init.d/voice-server

/etc/init.d/voice-server start
/etc/init.d/voice-server enable
netstat -tulpnw | grep 300
echo "done"
