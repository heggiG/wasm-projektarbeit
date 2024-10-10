const fileInput = document.getElementById("file-input");

const srcImage = document.getElementById("source-image");

const filterPicker = document.getElementById("filter-picker");

const colorPicker = document.getElementById("color-shift-color");
const shiftAmount = document.getElementById("shift-amount");

const vignetteCenterX = document.getElementById("vignette-center-x");
const vignetteCenterY = document.getElementById("vignette-center-y");
const vignetteStrength = document.getElementById("vignette-strength");
const vignetteRadius = document.getElementById("vignette-radius");

fileInput.addEventListener("change", handleFiles);

/**
 * Function to map clicks on the input image (which is shown cropped to a square) to coordinates on
 * the uncropped image.
 */
srcImage.addEventListener('click', function (event) {
    let woh = imgDimensions.x > imgDimensions.y ? imgDimensions.y : imgDimensions.x
    let bounds= this.getBoundingClientRect();
    let x = event.pageX - bounds.left;
    let y = event.pageY - bounds.top;
    let imageX= x / this.clientWidth * woh
    let imageY= y / this.clientHeight * woh
    let diff = Math.abs(imgDimensions.x - imgDimensions.y);
    if (imgDimensions.x > imgDimensions.y) {
        imageX += diff / 2
    } else {
        imageY += diff / 2
    }
    vignetteCenterX.value = Math.floor(imageX);
    vignetteCenterY.value = Math.floor(imageY);

});

let inputFile;
let imgDimensions;

async function handleFiles() {
    if (this.files.length) {
        const file = this.files[0];
        let url = URL.createObjectURL(new Blob([file]));
        srcImage.src = url;
        const arrayBuffer = await file.arrayBuffer();
        inputFile = new Uint8Array(arrayBuffer);
        const img = new Image();
        img.src = url;
        imgDimensions =  {x: img.width, y: img.height};
        fileInput.disabled = false;
    }
}

document.getElementById("apply-button").addEventListener('click', () => {
    let result;
    switch (filterPicker.elements['filter'].value) {
        case "sobel":
            result = window.applySobel(inputFile);
            break;
        case "gaussean":
            result = window.applyGaussian(inputFile);
            break;
        case "shift":
            result = window.applyShift(inputFile, colorPicker.value, +shiftAmount.value)
            break;
        case "vignette":
            result = window.applyVignette(inputFile, +vignetteRadius.value, +vignetteCenterX.value, +vignetteCenterY.value, +vignetteStrength.value)
            break;
    }
    document.getElementById("target-image").src = URL.createObjectURL(new Blob([result]));
});
