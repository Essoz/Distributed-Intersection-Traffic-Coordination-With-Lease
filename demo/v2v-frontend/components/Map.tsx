// components/Map.tsx

import React from 'react';
import { MapContainer, TileLayer, Marker, Popup } from 'react-leaflet';
import { CarData } from '../interfaces/CarData';

interface MapProps {
  carData: CarData[];
}

const Map: React.FC<MapProps> = ({ carData }) => {
  return (
    <MapContainer center={[0, 0]} zoom={13} style={{ height: '100vh', width: '100%' }}>
      <TileLayer
        attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
        url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
      />
      {carData.map(car => (
        <Marker key={car.metadata.name} position={car.dynamics.location}>
          <Popup>
            <strong>ID:</strong> {car.metadata.name} <br/>
          </Popup>
        </Marker>
      ))}
    </MapContainer>
  );
};

export default Map;
