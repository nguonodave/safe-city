let map = L.map('map').setView([-1.2921, 36.8219], 13); // Default to Nairobi

L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    maxZoom: 19,
}).addTo(map);

let marker = L.marker([-1.2921, 36.8219], { draggable: true }).addTo(map);

function updateCoordinates(lat, lng) {
    document.getElementById('coordinates').textContent = `Selected Coordinates: ${lat.toFixed(6)}, ${lng.toFixed(6)}`;
    document.getElementById('location').value = `${lat.toFixed(6)},${lng.toFixed(6)}`;
}

marker.on('dragend', function (e) {
    const { lat, lng } = e.target.getLatLng();
    updateCoordinates(lat, lng);
});

map.on('click', function (e) {
    const { lat, lng } = e.latlng;
    marker.setLatLng([lat, lng]);
    updateCoordinates(lat, lng);
});

// Initialize coordinates field with default marker position
updateCoordinates(marker.getLatLng().lat, marker.getLatLng().lng);
