export class BufferImage {
  main: HTMLImageElement;

  constructor() {
    this.main = new Image();
    this.main.src = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8/5+hHgAHggJ/PchI7wAAAABJRU5ErkJggg==";
  }

  update(url) {
    let that = this;
    let a = new Image();
    a.src = url;
    a.onload = () => {
      that.main = a;
    };
  }
}
