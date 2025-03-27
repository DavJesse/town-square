export function renderCreatePostButton() {
    // Create button element
    let createPostButton = document.createElement('button');
    createPostButton.id = 'create_post_button';

    // Create plus sign element
    let plusSign = document.createElement('span');
    plusSign.id = 'plus_sign';
    plusSign.textContent = '+';
    createPostButton.appendChild(plusSign);

    // Append button to app
    let app = document.getElementById('app');
    app.appendChild(createPostButton);

    return createPostButton
}

export function setCreatePostsButtonListeners(createPostButton) {
    // Add Event listener for clicking create posts button
    createPostButton.addEventListener('click', () => {
        alert('Create post button clicked');
    });
}