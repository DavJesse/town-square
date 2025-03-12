export function setupImageUpload () {
    console.log("✅ DOM fully loaded!"); // Check if this prints

    const imageInput = document.getElementById("image");
    if (!imageInput) {
        console.error("⚠️ ERROR: #image input field not found!");
        return;
    } else {
        console.log("✅ Image input field found!", imageInput);
    }

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
        console.log("📸 File selected:", this.files.length > 0 ? this.files[0].name : "No file");
        if (this.files.length > 0) {
            fileNameDisplay.textContent = this.files[0].name;
            removeButton.style.display = "inline-block"; // Show remove button
        }
    });

    // Handle image removal
    removeButton.addEventListener("click", function (event) {
        event.preventDefault(); // Prevent form submission
        console.log("🗑 Removing file...");
        imageInput.value = ""; // Clear input
        fileNameDisplay.textContent = "No file chosen";
        removeButton.style.display = "none"; // Hide button again
    });
}
