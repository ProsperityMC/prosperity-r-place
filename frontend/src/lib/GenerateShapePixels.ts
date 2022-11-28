export function GenerateShapePixels(
  width: number,
  height: number,
  shapeArea: {x1: number; y1: number; x2: number; y2: number},
  shapeSel: string,
  paletteSel: number,
) {
  let pixels = new Uint16Array(width * height);
  switch (shapeSel) {
    case "square":
      for (let i = shapeArea.x1; i < shapeArea.x2; i++) {
        pixels[shapeArea.y1 * width + i] = paletteSel;
        pixels[shapeArea.y2 * width + i + 1] = paletteSel;
      }

      for (let i = shapeArea.y1; i < shapeArea.y2; i++) {
        pixels[i * width + shapeArea.x2] = paletteSel;
        pixels[(i + 1) * width + shapeArea.x1] = paletteSel;
      }
      break;
  }
  return pixels;
}
