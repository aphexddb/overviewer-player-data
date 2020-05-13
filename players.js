window.PlayerLocations = {
  dataURL: "players.json",
  debug: false,
  playerMarkers: {},
  refreshSeconds: 60,

  createPlayerMarkers: (list) => {
    for (let playerName in window.PlayerLocations.playerMarkers) {
      if (!list.hasOwnProperty(playerName)) {
        window.PlayerLocations.playerMarkers[playerName].remove();
        delete window.PlayerLocations.playerMarkers[playerName];
      }
    }

    for (let playerName in list) {
      let playerData = list[playerName];
      let multi = playerData.dimension === "-1" ? 8 : 1;
      let latlng = overviewer.util.fromWorldToLatLng(
        playerData.x * multi,
        playerData.y,
        playerData.z * multi,
        window.PlayerLocations.getCurrentTileSet()
      );

      if (window.PlayerLocations.playerMarkers[playerName]) {
        window.PlayerLocations.playerMarkers[playerName].setLatLng(latlng);
      } else {
        let icon = L.icon({
          iconUrl: "https://overviewer.org/avatar/" + playerName,
          iconSize: [16, 32],
          iconAnchor: [15, 33],
        });

        let marker = L.marker(latlng, {
          icon: icon,
          title: `${playerName} @ ${Math.round(playerData.x)} , ${Math.round(
            playerData.y
          )} , ${Math.round(playerData.z)}`,
        });

        marker.addTo(overviewer.map);

        window.PlayerLocations.playerMarkers[playerName] = marker;
      }
    }
  },

  getCurrentTileSet: () => {
    let name = overviewer.current_world;
    for (let index in overviewerConfig.tilesets) {
      let tileset = overviewerConfig.tilesets[index];
      if (tileset.world === name) {
        return tileset;
      }
    }
  },

  initialize: () => {
    setTimeout(
      window.PlayerLocations.load,
      window.PlayerLocations.refreshSeconds * 1000
    );
    window.PlayerLocations.load();
  },

  load: () => {
    if (window.PlayerLocations.debug) {
      console.log("Loading player data from", window.PlayerLocations.dataURL);
    }

    try {
      var request = new XMLHttpRequest();
      request.open("GET", window.PlayerLocations.dataURL, true);

      request.onload = function () {
        if (request.status >= 200 && request.status < 400) {
          const data = JSON.parse(request.responseText);
          window.PlayerLocations.createPlayerMarkers(data);
        } else {
          console.error("Error reading player data", request.responseText);
        }
      };

      request.onerror = function () {
        console.error("Error reading player data");
      };

      request.send();
    } catch (error) {
      console.error("Error reading player data", error);
    }

    if (window.PlayerLocations.debug) {
      console.log("Sleeping for", window.PlayerLocations.refreshSeconds, "s");
    }

    setTimeout(
      window.PlayerLocations.load,
      window.PlayerLocations.refreshSeconds * 1000
    );
  },
};

overviewer.util.ready(window.PlayerLocations.initialize);
