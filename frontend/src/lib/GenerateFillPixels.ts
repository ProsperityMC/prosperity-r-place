interface PixelCtx {
  pixels: Uint16Array;
  paletteSel: number;
  width: number;
  height: number;
}

export function GenerateFillPixels(width: number, height: number, pos: {x: number; y: number}, data: ImageData, paletteSel: number) {
  let pixels = new Uint16Array(width * height);
  let ctx: PixelCtx = {pixels, paletteSel, width, height};

  floodFill(ctx, data, width, height, pos.x, pos.y);
  return pixels;
}

function floodFill(ctx: PixelCtx, data: ImageData, width: number, height: number, x: number, y: number) {
  let store = new Uint16Array(width * height);
  let first = GetHexIDP(data, x, y);
  let queue: Array<{x: number; y: number}> = [{x, y}];
  store[y * width + x] = 1;
  while (queue.length > 0) {
    let last = queue.pop();
    PutPixel(ctx, last.x, last.y);
    floodPixelNeighbours(data, last.x + 1, last.y, store, queue, first);
    floodPixelNeighbours(data, last.x - 1, last.y, store, queue, first);
    floodPixelNeighbours(data, last.x, last.y + 1, store, queue, first);
    floodPixelNeighbours(data, last.x, last.y - 1, store, queue, first);
  }
}

function floodPixelNeighbours(data: ImageData, x: number, y: number, store: Uint16Array, queue: Array<{x: number; y: number}>, check: number) {
  if (x < 0 || y < 0 || x >= data.width || y >= data.height) return;
  if (store[y * data.width + x] === 1) return;
  store[y * data.width + x] = 1;
  if (GetHexIDP(data, x, y) === check) queue.push({x, y});
}

function PutPixel(ctx: PixelCtx, x: number, y: number) {
  if (x < 0 || y < 0 || x >= ctx.width || y >= ctx.height) return;
  ctx.pixels[y * ctx.width + x] = ctx.paletteSel;
}

function GetIDP(data: ImageData, x: number, y: number) {
  let v = data.width * y * 4 + x * 4;
  return data.data.slice(v, v + 3); // use 3 to only get RGB
}

function GetHexIDP(data: ImageData, x: number, y: number) {
  return [...GetIDP(data, x, y).values()].reduce((a, b) => (a << 8) | b);
}
