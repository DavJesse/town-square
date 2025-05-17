export function renderProfileMenu(userData) {
    // Grab navbar right cluster
    let navbarRightCluster = document.getElementById('navbar_right_cluster');

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
        renderNavbarBio(userData);
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

function renderNavbarBio(userData) {
    // Grab app
    let app = document.getElementById('app');

    // Create smoke screen
    let smokeScreen = document.createElement('div');
    smokeScreen.id = 'smoke_screen';
    app.appendChild(smokeScreen);

    // Create bio container
    let bioMenuContainer = document.createElement('div');
    bioMenuContainer.id = 'navbar_bio_container';
    smokeScreen.appendChild(bioMenuContainer);

    // Create close button and add related event listener
    let closeButton = document.createElement('button');
    closeButton.id = 'navbar_bio_close_button';
    closeButton.textContent = 'X';
    closeButton.addEventListener('click', () => {
        bioContainer.remove();
    });
    bioContainer.appendChild(closeButton);

    // Create bio title
    let profileHead = document.createElement('div');
    let profileTitle = document.createElement('h2');
    let profileSubtitle = document.createElement('h2');
    profileTitle.id = 'profile_title';
    profileSubtitle.id = 'profile_subtitle';
    profileHead.classList.add('profile-user-info-container');
    profileTitle.textContent = userData.name;
    profileSubtitle.textContent = `@${userData.username}`;
    profileHead.appendChild(profileTitle);
    profileHead.appendChild(profileSubtitle);
    bioMenuContainer.appendChild(profileHead);

    // Render user image
    let profilePic = document.createElement('img');
    profilePic.id = 'index_profile_pic';
    profilePic.src = userData.image;
    profilePic.alt = userData.name;
    bioMenuContainer.appendChild(profilePic);

    // Render user info
    let profileUserInfo = document.createElement('div');
    let profileContact = document.createElement('h3');
    let profileAgeGender = document.createElement('div');
    let profileGender = document.createElement('h3');
    let profileAgeGenderSeparator = document.createElement('h3');
    let profileAge = document.createElement('h3');
    profileUserInfo.classList.add('profile-user-info-container');
    profileGender.classList.add('profile-age-gender');
    profileAge.classList.add('profile-age-gender');
    profileContact.id = 'profile_contact';
    profileAgeGender.id = 'profile_age_gender_container';
    profileGender.id = 'profile_gender';
    profileAgeGenderSeparator.id = 'profile_age_gender_separator';
    profileAge.id = 'profile_age';
    profileContact.textContent = `mail: ${userData.email}`;
    profileGender.textContent = userData.gender;
    profileAgeGenderSeparator.textContent = '|';
    profileAge.textContent = userData.age;
    profileAgeGender.appendChild(profileGender);
    profileAgeGender.appendChild(profileAgeGenderSeparator);
    profileAgeGender.appendChild(profileAge);
    profileUserInfo.appendChild(profileContact);
    profileUserInfo.appendChild(profileAgeGender);
    bioMenuContainer.appendChild(profileUserInfo);

    // Render Bio section
    let bioContainer = document.createElement('div');
    let bioTitle = document.createElement('h3');
    let bioText = document.createElement('p');
    bioContainer.id = 'index_bio_container';
    bioTitle.id = 'profile_bio_title';
    bioText.id = 'profile_bio_text';
    bioTitle.textContent = `About ${userData.first_name}`;
    bioText.textContent = userData.bio;
    bioContainer.appendChild(bioTitle);
    bioContainer.appendChild(bioText);
    bioMenuContainer.appendChild(bioContainer);
    
    // Render action buttons
    let profileActionContainer = document.createElement('div');
    let logoutLink = document.createElement('form');
    let logoutButton = document.createElement('button');
    profileActionContainer.id = 'profile_action_container';
    logoutButton.id = 'profile_logout_button';
    logoutLink.method = 'POST';
    logoutButton.type = 'submit';
    logoutLink.action = '/logout';
    logoutButton.textContent = 'Logout';
    logoutLink.appendChild(logoutButton);
    profileActionContainer.appendChild(logoutLink);
    bioMenuContainer.appendChild(profileActionContainer);    

}
