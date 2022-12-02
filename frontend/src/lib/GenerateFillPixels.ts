interface PixelCtx {
  pixels: Uint16Array;
  paletteSel: number;
  width: number;
  height: number;
}

export function GenerateFillPixels(width: number, height: number, pos: {x: number; y: number}, data: ImageData, paletteSel: number) {
  let pixels = new Uint16Array(width * height);
  let ctx: PixelCtx = {pixels, paletteSel, width};
  console.log(ctx);

  PutPixel(ctx, i, sy);
  return pixels;
}

function PutPixel(ctx: PixelCtx, x: number, y: number) {
  if (x < 0 || y < 0 || x >= ctx.width || y >= ctx.height) return;
  ctx.pixels[y * ctx.width + x] = ctx.paletteSel;
}
