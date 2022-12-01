export class BufferImage {
  url: string;
  main: HTMLCanvasElement;

  constructor(url, width, height) {
    this.url = url;
    this.main = document.createElement("canvas");
    this.update("data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8/5+hHgAHggJ/PchI7wAAAABJRU5ErkJggg==");
  }

  update(url = "") {
    // PNG is imported... trust me...
    PNG.load(url, this.main);
  }
}
