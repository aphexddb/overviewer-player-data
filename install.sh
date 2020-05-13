#!/bin/bash
set -e

echo "Installing packages"

sudo add-apt-repository ppa:longsleep/golang-backports
sudo apt update
sudo apt install -y golang-go git

echo "Installing ODP"

cd /usr/local/src

folder="/usr/local/src/overviewer-player-data"
if ! git clone https://github.com/aphexddb/overviewer-player-data.git "${folder}" 2>/dev/null && [ -d "${folder}" ] ; then
    echo "Clone failed because the folder ${folder} exists"
fi

cd overviewer-player-data
sudo chmod +x /usr/local/src/overviewer-player-data/odp.sh
sudo make
rm /usr/local/bin/odp || true
sudo ln -s /usr/local/src/overviewer-player-data/odp /usr/local/bin/odp

mkdir -p /opt/msm/html
cp /usr/local/src/overviewer-player-data/players.js /opt/msm/html/players.js

echo "Configuring ODP service"

sudo cp /usr/local/src/overviewer-player-data/odp.service /etc/systemd/system/odp.service
sudo chmod 644 /etc/systemd/system/odp.service

sudo systemctl daemon-reload
sudo systemctl enable odp.service

sudo systemctl status odp

echo "--------------------------------------"
echo " To complete installation:"
echo ""
echo " 1. Edit config in:"
echo " /usr/local/src/overviewer-player-data/odp.sh"
echo ""
echo " 2. Start the service"
echo "sudo systemctl start odp"
echo "--------------------------------------"
