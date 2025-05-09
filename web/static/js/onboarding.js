export function setupImageUpload () {
    const imageInput = document.getElementById("image");
    const fileNameDisplay = document.getElementById("file_name");

    // Create remove button dynamically
    let removeButton = document.createElement("button");
    removeButton.id = "remove_image";
    removeButton.classList.add("remove-btn");
    removeButton.textContent = "Remove Image";
    removeButton.style.display = "none"; // Hide initially

    // Append button after file name display
    fileNameDisplay.insertAdjacentElement("afterend", removeButton);

    // Handle file selection
    imageInput.addEventListener("change", function () {
        if (this.files.length > 0) {
            fileNameDisplay.textContent = this.files[0].name;
            removeButton.style.display = "inline-block"; // Show remove button
        }
    });

    // Handle image removal
    removeButton.addEventListener("click", function (event) {
        event.preventDefault(); // Prevent form submission
        imageInput.value = ""; // Clear input
        fileNameDisplay.textContent = "No file chosen";
        removeButton.style.display = "none"; // Hide button again
    });
}