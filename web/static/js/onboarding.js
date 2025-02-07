document.addEventListener("DOMContentLoaded", function() {
    const imageInput = document.getElementById("image");
    const fileNameDisplay = document.getElementById("file-name");
    const removeButton = document.getElementById("remove-image");

    imageInput.addEventListener("change", function() {
        if (this.files.length > 0) {
            fileNameDisplay.textContent = this.files[0].name;
            removeButton.style.display = "inline-block"; // Show remove button
        } else {
            fileNameDisplay.textContent = "No file chosen";
            removeButton.style.display = "none"; // Hide remove button
        }
    });

    removeButton.addEventListener("click", function(event) {
        event.preventDefault(); // Prevent form submission
        imageInput.value = ""; // Clear file input
        fileNameDisplay.textContent = "No file chosen";
        removeButton.style.display = "none"; // Hide remove button
    });
});

