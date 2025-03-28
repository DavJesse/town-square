export function renderCreatePostButton() {
    // Create button element
    let createPostButton = document.createElement('button');
    createPostButton.id = 'create_post_button';

    // Create plus sign element
    let plusSign = document.createElement('span');
    plusSign.id = 'plus_sign';
    plusSign.textContent = '+';

    // Create "Create Post" text element
    let createText = document.createElement('span');
    createText.id = 'create_post_text';
    createText.textContent = 'Create Post';

    // Append elements
    createPostButton.appendChild(plusSign);
    createPostButton.appendChild(createText);

    // Append button to app
    let app = document.getElementById('app');
    app.appendChild(createPostButton);

    return createPostButton;
}

export function setCreatePostsButtonListeners(createPostButton) {
    // Add Event listener for clicking create posts button
    createPostButton.addEventListener('click', () => {
        alert('Create post button clicked');
    });
}