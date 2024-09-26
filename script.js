const fileInput = document.getElementById("file-input");
const filterPicker = document.getElementById("filter-picker")
fileInput.addEventListener("change", handleFiles);

let inputFile;

async function handleFiles() {
    if (this.files.length) {
        const file = this.files[0];
        document.getElementById("sourceImage").src = URL.createObjectURL(new Blob([file]));
        const arrayBuffer = await file.arrayBuffer();
        inputFile = new Uint8Array(arrayBuffer);
    }
}

document.getElementById("apply-button").addEventListener('click', async () => {
    let result;
    switch (filterPicker.elements['filter'].value) {
        case "sobel":
            result = await window.applySobel(inputFile);
            break;
        case "gaussean":
            result = await window.applyGaussean(inputFile);
            break;
    }
    targetImage.src = URL.createObjectURL(new Blob([result]));
});
