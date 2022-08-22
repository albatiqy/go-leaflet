import './style.css'
import "leaflet/dist/leaflet.css"
import L from "leaflet"

const map = L.map('map', {
    zoomControl: true, maxZoom: 28, minZoom: 1
}).fitBounds([[-13.692431032626045, 94.3990624407764], [9.741550497819643, 143.1067589294119]]);
// L.Hash(map);
const bounds_group = new L.featureGroup([]);
L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '© OpenStreetMap'
}).addTo(map);
map.createPane('pane_Simplified_0');
map.getPane('pane_Simplified_0').style.zIndex = 400;
map.getPane('pane_Simplified_0').style['mix-blend-mode'] = 'normal';
const myLayer = L.geoJSON(null, {
    attribution: 'Albatiqy',
    interactive: true,
    // dataVar: 'json_Simplified_0',
    // layerName: 'layer_Simplified_0',
    pane: 'pane_Simplified_0',
    onEachFeature(feature, layer) {
        layer.bindPopup(feature.properties['wadmpr']);
    },
    style(feature) {
        return {
            pane: 'pane_Simplified_0',
            opacity: 1,
            color: 'rgba(35,35,35,1.0)',
            dashArray: '',
            lineCap: 'butt',
            lineJoin: 'miter',
            weight: 1.0,
            fill: true,
            fillOpacity: 1,
            fillColor: 'rgba(164,113,88,1.0)',
            interactive: true,
        }
    },
}).addTo(map);
bounds_group.addLayer(myLayer);
const req = await fetch("/propinsi.geojson");
myLayer.addData(await req.json())