import { Monitor, Smartphone, Bot } from "@lucide/svelte";

const getDeviceIcon = (device: string) => {
  switch (device) {
    case "desktop":
      return Monitor;
    case "mobile":
      return Smartphone;
    case "bot":
      return Bot;
    default:
      return Monitor;
  }
};

export { getDeviceIcon };
