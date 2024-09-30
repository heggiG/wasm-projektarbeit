const fileInput = document.getElementById("file-input");

const filterPicker = document.getElementById("filter-picker");

const colorPicker = document.getElementById("color-shift-color");
const shiftAmount = document.getElementById("shift-amount");

const vignetteCenterX = document.getElementById("vignette-center-x");
const vignetteCenterY = document.getElementById("vignette-center-y");
const vignetteStrength = document.getElementById("vignette-strength");
const vignetteRadius = document.getElementById("vignette-radius");

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

document.getElementById("apply-button").addEventListener('click', () => {
    let result;
    switch (filterPicker.elements['filter'].value) {
        case "sobel":
            result = window.applySobel(inputFile);
            break;
        case "gaussean":
            result = window.applyGaussean(inputFile);
            break;
        case "shift":
            result = window.applyShift(inputFile, colorPicker.value, +shiftAmount.value)
            break;
        case "vignette":
            result = window.applyVignette(inputFile, +vignetteRadius.value, +vignetteCenterX.value, +vignetteCenterY.value, +vignetteStrength.value)
            break;
    }
    document.getElementById("targetImage").src = URL.createObjectURL(new Blob([result]));
});
