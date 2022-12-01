interface PixelCtx {
  pixels: Uint16Array;
  paletteSel: number;
  width: number;
  height: number;
}

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

  let dx = lx - sx;
  let dy = ly - sy;

  let mx = dx / 2 + sx;
  let mx1 = Math.floor(mx);
  let mx2 = Math.ceil(mx);

  let my = dy / 2 + sy;
  let my1 = Math.floor(my);
  let my2 = Math.ceil(my);

  let pixels = new Uint16Array(width * height);
  let ctx: PixelCtx = {pixels, paletteSel, width};
  switch (shapeSel) {
    case "circle":
      //console.log(shapeArea.x1, shapeArea.y1, dx / 2, dy / 2);
      BresenhamEllipse(ctx, shapeArea.x1, shapeArea.y1, Math.abs(dx), Math.abs(dy));
      break;
    case "triangle":
      //   mx
      //  /  \
      // /____\
      BresenhamLine(ctx, mx2, sy, lx, ly);
      BresenhamLine(ctx, lx, ly, sx, ly);
      BresenhamLine(ctx, sx, ly, mx1, sy);
      break;
    case "square":
      for (let i = sx; i < lx; i++) {
        PutPixel(ctx, i, sy);
        PutPixel(ctx, i + 1, ly);
      }

      for (let i = sy; i < ly; i++) {
        PutPixel(ctx, lx, i);
        PutPixel(ctx, sx, i + 1);
      }
      break;
    case "diamond":
      //   mx
      //  /  \
      // m    m
      // y    y
      //  \  /
      //   mx
      BresenhamLine(ctx, mx2, sy, lx, my1);
      BresenhamLine(ctx, mx2, ly, lx, my2);
      BresenhamLine(ctx, mx1, ly, sx, my2);
      BresenhamLine(ctx, mx1, sy, sx, my1);
      break;
  }
  return pixels;
}

function BresenhamLine(ctx: PixelCtx, sx: number, sy: number, lx: number, ly: number) {
  let dx = Math.abs(lx - sx);
  let dy = Math.abs(ly - sy);
  let fx = sx < lx ? 1 : -1;
  let fy = sy < ly ? 1 : -1;
  let e = dx - dy;

  while (true) {
    PutPixel(ctx, sx, sy);
    if (sx == lx && sy == ly) break;
    let e2 = 2 * e;
    if (e2 > -dy) {
      e -= dy;
      sx += fx;
    }
    if (e2 < dx) {
      e += dx;
      sy += fy;
    }
  }
}

// Function for circle-generation
// using Bresenham's algorithm
function BresenhamEllipse(ctx: PixelCtx, cx: number, cy: number, rx: number, ry: number) {
  let dx, dy, d1, d2;
  let x = 0;
  let y = ry;

  // Initial decision parameter of region 1
  d1 = ry * ry - rx * rx * ry + 0.25 * rx * rx;
  dx = 2 * ry * ry * x;
  dy = 2 * rx * rx * y;

  // For region 1
  while (dx < dy) {
    PutPixel(ctx, cx + x, cy + y);
    PutPixel(ctx, cx - x, cy + y);
    PutPixel(ctx, cx + x, cy - y);
    PutPixel(ctx, cx - x, cy - y);

    // Checking and updating value of
    // decision parameter based on algorithm
    if (d1 < 0) {
      x++;
      dx = dx + 2 * ry * ry;
      d1 = d1 + dx + ry * ry;
    } else {
      x++;
      y--;
      dx = dx + 2 * ry * ry;
      dy = dy - 2 * rx * rx;
      d1 = d1 + dx - dy + ry * ry;
    }
  }

  // Decision parameter of region 2
  d2 = ry * ry * ((x + 0.5) * (x + 0.5)) + rx * rx * ((y - 1) * (y - 1)) - rx * rx * ry * ry;

  // Plotting points of region 2
  while (y >= 0) {
    PutPixel(ctx, cx + x, cy + y);
    PutPixel(ctx, cx - x, cy + y);
    PutPixel(ctx, cx + x, cy - y);
    PutPixel(ctx, cx - x, cy - y);

    // Checking and updating parameter
    // value based on algorithm
    if (d2 > 0) {
      y--;
      dy = dy - 2 * rx * rx;
      d2 = d2 + rx * rx - dy;
    } else {
      y--;
      x++;
      dx = dx + 2 * ry * ry;
      dy = dy - 2 * rx * rx;
      d2 = d2 + dx - dy + rx * rx;
    }
  }
}

function PutPixel(ctx: PixelCtx, x: number, y: number) {
  if (x < 0 || y < 0 || x >= ctx.width || y >= ctx.height) return;
  ctx.pixels[y * ctx.width + x] = ctx.paletteSel;
}
