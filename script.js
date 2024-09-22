const fileInput = document.getElementById("fileInput");
fileInput.addEventListener("change", handleFiles);

async function handleFiles() {
    if (this.files.length) {
        const file = this.files[0];
        document.getElementById("sourceImage").src = URL.createObjectURL(new Blob([file]));
        const arrayBuffer = await file.arrayBuffer();
        const input = new Uint8Array(arrayBuffer);
        const result = window.applySobel(input);
        targetImage.src = URL.createObjectURL(new Blob([result]));
    }
}
