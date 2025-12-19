import fs from "fs";

class FileReader {
  constructor() {
    // These properties exist in the browser API
    this.onload = null;
    this.onerror = null;
    this.result = null;
  }

  readAsText(filename, encoding = "utf8") {
    fs.readFile(filename, encoding, (err, data) => {
      if (err) {
        if (this.onerror) {
          this.onerror(err);
        }
        return;
      }

      this.result = data;

      if (this.onload) {
        // Browser FileReader passes an event object
        this.onload({ target: this });
      }
    });
  }
}

// Make it globally available (like window.FileReader in browsers)
globalThis.FileReader = FileReader;
