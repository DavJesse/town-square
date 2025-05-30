export function renderLogoutButton() {
    // Create button element
    let logoutButton = document.createElement('button');
    logoutButton.id = 'create_post_button';

    // Create logout icon
    let logoutIcon = document.createElement('span');
    logoutIcon.id = 'logout_icon';
    logoutIcon.classList.add('material-symbols-outlined');
    logoutIcon.textContent = 'logout';

    // Create "Create Post" text element
    let logoutText = document.createElement('span');
    logoutText.id = 'create_post_text';
    logoutText.textContent = 'Create Post';

    // Append elements
    logoutButton.appendChild(logoutIcon);
    logoutButton.appendChild(logoutText);
    setLogoutButtonListeners(logoutButton);

    // Append button to app
    // let app = document.getElementById('app');
    document.body.appendChild(logoutButton);

    // return logoutButton;
}

export function setLogoutButtonListeners(logoutButton) {
    // Add Event listener for clicking create posts button
    logoutButton.addEventListener('click', () => {
        fetch('/logout', {
            method: 'POST',
            credentials: 'include',
        })
    });
}