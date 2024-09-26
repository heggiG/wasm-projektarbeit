const fileInput = document.getElementById("file-input");
const filterPicker = document.getElementById("filter-picker")
fileInput.addEventListener("change", handleFiles);

async function handleFiles() {
    if (this.files.length) {
        const file = this.files[0];
        document.getElementById("sourceImage").src = URL.createObjectURL(new Blob([file]));
        const arrayBuffer = await file.arrayBuffer();
        const input = new Uint8Array(arrayBuffer);
        let result;
        switch (filterPicker.elements['filter'].value) {
            case "sobel":
                result = window.applySobel(input);
                break;
            case "gaussean":
                result = window.applyGaussean(input);
                break;
        }
        targetImage.src = URL.createObjectURL(new Blob([result]));
    }
}
