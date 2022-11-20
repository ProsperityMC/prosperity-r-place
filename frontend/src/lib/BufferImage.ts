export class BufferImage {
  url: string;
  main: HTMLImageElement;

  constructor(url) {
    this.url = url;
    this.main = new Image();
    this.main.src = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8/5+hHgAHggJ/PchI7wAAAABJRU5ErkJggg==";
  }

  update() {
    fetch(this.url)
      .then(res => res.blob())
      .then(blob => {
        this.main.src = URL.createObjectURL(blob);
      });
  }
}
