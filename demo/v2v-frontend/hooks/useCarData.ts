// hooks/useCarData.ts

import { useState, useEffect } from 'react';
import axios from 'axios';
import { CarData } from '../interfaces/CarData';

export const useCarData = (refreshInterval: number = 5000) => {
  const [carData, setCarData] = useState<CarData[]>([]);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await axios.get('http://localhost:11002/car/getAll');
        // Convert the returned map to an array of CarData objects
        const carDataArray: CarData[] = Object.values(response.data);
        setCarData(carDataArray);
      } catch (error) {
        console.error(error);
      }
    };

    fetchData();
    const interval = setInterval(fetchData, refreshInterval);
    return () => clearInterval(interval);
  }, [refreshInterval]);

  return carData;
};
