export class BufferImage {
  url: string;
  main: HTMLImageElement;

  constructor(url) {
    this.url = url;
    this.main = new Image();
    this.main.src = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8/5+hHgAHggJ/PchI7wAAAABJRU5ErkJggg==";
  }

  update(url = "") {
    if (url.startsWith("data:")) this.main.src = url;
    else
      fetch(url ? url : this.url)
        .then(res => res.blob())
        .then(blob => {
          this.main.src = URL.createObjectURL(blob);
        });
  }
}
