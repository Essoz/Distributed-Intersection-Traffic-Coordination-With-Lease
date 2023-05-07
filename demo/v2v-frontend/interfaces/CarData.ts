// interfaces/CarData.ts

// interfaces/CarData.ts

export interface CarDynamicsPassingBlocksElem {
    name: string;
    position: [number, number];
    size: [number, number];
  }
  
export interface CarData {
    dynamics: {
        acceleration: number;
        destination: [number, number];
        heading: number;
        location: [number, number];
        passingBlocks: CarDynamicsPassingBlocksElem[];
        speed: [number, number];
        stage: string;
    };
    metadata: {
        name: string;
        isV2v: boolean;
        type: string;
    };
}
