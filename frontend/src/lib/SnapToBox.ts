interface Obj {
  x: number;
  y: number;
  width: number;
  height: number;
}

interface Box {
  width: number;
  height: number;
}

export function snapInsideBox(obj: Obj, box: Box): Obj {
  if (obj.x + obj.width > box.width) obj.x = box.width - obj.width;
  if (obj.y + obj.height > box.height) obj.y = box.height - obj.height;
  if (obj.x < 0) obj.x = 0;
  if (obj.y < 0) obj.y = 0;
  return obj;
}

export function snapAroundBox(obj: Obj, box: Box): Obj {
  if (obj.x + obj.width < box.width) obj.x = box.width - obj.width;
  if (obj.y + obj.height < box.height) obj.y = box.height - obj.height;
  if (obj.x > 0) obj.x = 0;
  if (obj.y > 0) obj.y = 0;
  return obj;
}

export function snapInsideLine(a: number, width: number, line: number): number {
  if (a + width > line) a = line - width;
  if (a < 0) a = 0;
  return a;
}

export function snapOutsideLine(a: number, width: number, line: number): number {
  if (a + width < line) a = line - width;
  if (a > 0) a = 0;
  return a;
}
