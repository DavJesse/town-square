export function renderProfileMenu(userData) {
    // Grab navbar right cluster
    let navbarRightCluster = document.getElementById('navbar_right_custer');

    // Create profile menu button
    let profileMenuButton = document.createElement('button');
    profileMenuButton.id = 'profile_menu_button';
    profileMenuButton.classList.add('inactive');

    // Add event listener for profile menu button
    profileMenuButton.addEventListener('click', () => {
        profileMenuButton.classList.add('active');
        profileMenuButton.classList.remove('inactive');
    });

    // Render profile menu
    if (profileMenuButton.classList.contains('active')) {
        
    }

    // Create profile menu image
    let profileMenuImage = document.createElement('image');
    profileMenuImage.id = 'profile_menu_image';
    profileMenuImage.src = userData.image;
    profileMenuImage.alt = userData.name;

    // Append image to button
    profileMenuButton.appendChild(profileMenuImage);

    // Append button to navbar right cluster
    navbarRightCluster.appendChild(profileMenuButton);

}

