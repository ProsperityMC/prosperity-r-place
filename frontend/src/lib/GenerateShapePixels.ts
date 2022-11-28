export function GenerateShapePixels(
  width: number,
  height: number,
  shapeArea: {x1: number; y1: number; x2: number; y2: number},
  shapeSel: string,
  paletteSel: number,
) {
  let cx = shapeArea.x1 < shapeArea.x2;
  let sx = cx ? shapeArea.x1 : shapeArea.x2;
  let lx = cx ? shapeArea.x2 : shapeArea.x1;

  let cy = shapeArea.y1 < shapeArea.y2;
  let sy = cy ? shapeArea.y1 : shapeArea.y2;
  let ly = cy ? shapeArea.y2 : shapeArea.y1;

  let pixels = new Uint16Array(width * height);
  switch (shapeSel) {
    case "square":
      for (let i = sx; i < lx; i++) {
        pixels[sy * width + i] = paletteSel;
        pixels[ly * width + i + 1] = paletteSel;
      }

      for (let i = sy; i < ly; i++) {
        pixels[i * width + lx] = paletteSel;
        pixels[(i + 1) * width + sx] = paletteSel;
      }
      break;
  }
  return pixels;
}
