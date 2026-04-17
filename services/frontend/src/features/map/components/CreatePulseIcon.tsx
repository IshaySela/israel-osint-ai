import L from 'leaflet';

const createPulseIcon = (isSelected: boolean) => {
  return L.divIcon({
    className: 'custom-pulse-icon',
    html: `
      <div class="relative flex items-center justify-center">
        <div class="absolute w-4 h-4 rounded-full bg-cyan-500 ${isSelected ? 'animate-ping scale-150' : 'animate-pulse opacity-50'}"></div>
        <div class="relative w-3 h-3 rounded-full bg-cyan-400 border border-white shadow-[0_0_10px_rgba(34,211,238,0.8)]"></div>
      </div>
    `,
    iconSize: [20, 20],
    iconAnchor: [10, 10],
  });
};

export { createPulseIcon };