#!/bin/bash
set -ex

if [[ `id -nu` != "minecraft" ]];then
  echo "Not minecraft user, exiting.."
  exit 1
fi

odp -file /opt/msm/html/players.json -host localhost -password YOUR_RCON_SECRET_HERE

