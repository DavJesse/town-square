export function renderLogoutButton() {
    // Create button element
    let logoutButton = document.createElement('button');
    logoutButton.id = 'logout_button';

    // Create plus sign element
    let logoutIcon = document.createElement('span');
    logoutIcon.classList.add('material-symbols-outlined');
    logoutIcon.id = 'logout_icon';
    logoutIcon.textContent = 'logout';

    // Create "Create Post" text element
    let logoutText = document.createElement('span');
    logoutText.id = 'logout_text';
    logoutText.textContent = 'logout';

    // Append elements
    logoutButton.appendChild(logoutIcon);
    logoutButton.appendChild(logoutText);

    return logoutButton;
}

export function setLogoutButtonListeners(logoutButton) {
    // Add Event listener for clicking create posts button
    logoutButton.addEventListener('click', () => {
        alert('Create post button clicked');
    });
}