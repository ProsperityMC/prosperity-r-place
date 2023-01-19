export class BufferImage {
  url: string;
  main: HTMLCanvasElement;

  constructor(url, width, height) {
    this.main = document.createElement("canvas");
    this.update("data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8/5+hHgAHggJ/PchI7wAAAABJRU5ErkJggg==");
  }

  update(url = "") {
    if (url == this.url) return;
    // PNG is imported... trust me...
    console.log("image changed");
    this.url = url;
    PNG.load(url, this.main);
  }
}
