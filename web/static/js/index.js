import { renderNavBar } from '/static/js/navbar.js'
import { renderCreatePostButton, setCreatePostsButtonListeners } from '/static/js/create_post_button.js';
import { navigateTo } from '/static/js/routes.js';
import { populatePosts, setToggleEventListeners } from '/static/js/populate_posts.js';
import { renderErrorPage } from '/static/js/error.js';
import { populateCategories } from '/static/js/populate_categories.js';
import { renderMobileCategoriesMenu } from '/static/js/mobile_categories_menu.js';
import { initChat, populateOnlineUsersList, fetchAllUsers } from '/static/js/chat.js';

export function renderIndexPage() {
    // Extract app from dom
    let app = document.getElementById('app');
    app.innerHTML = "";

    // Set document title
    document.title = 'real-time-forum';

    // Create index page container
    let indexPageContainer = document.createElement('div');
    indexPageContainer.id = 'index_page'
    app.appendChild(indexPageContainer);

    // Render create post button
    let userOptions = document.createElement('div');
    userOptions.id = 'user_options';

    // Render create post button and append to app
    let createPostButton = renderCreatePostButton();
    setCreatePostsButtonListeners(createPostButton);
    userOptions.appendChild(createPostButton);
    app.appendChild(userOptions);

    // Create clusters of index page
    let centerCluster = document.createElement('div');
    let rightCluster = document.createElement('div');
    centerCluster.id = 'center_cluster';
    rightCluster.id = 'right_cluster';

    indexPageContainer.appendChild(centerCluster);
    indexPageContainer.appendChild(rightCluster);

    // Only create left cluster if screen is larger than 540px
    // if (window.innerWidth > 540) {
    let leftCluster = document.createElement('div');
    leftCluster.id = 'left_cluster';
    indexPageContainer.appendChild(leftCluster);

    // Create categories and online users containers
    let categoriesCard = document.createElement('div');
    let onlineUsersCard = document.createElement('div');
    let onlineUsersTitle = document.createElement('h3');
    let categoriesCardTitle = document.createElement('h3');
    let categoryContentContainer = document.createElement('div');
    categoriesCard.id = 'category_container';
    onlineUsersCard.id = 'online_users_card';
    onlineUsersTitle.id = 'online_users_title';
    categoriesCardTitle.id = 'categories_title';
    categoryContentContainer.id = 'category_content_container';
    onlineUsersTitle.textContent = 'Who\'s online?';
    categoriesCardTitle.textContent = 'Categories';

    onlineUsersCard.appendChild(onlineUsersTitle);
    categoriesCard.appendChild(categoriesCardTitle);
    categoriesCard.appendChild(categoryContentContainer);
    leftCluster.appendChild(categoriesCard);
    leftCluster.appendChild(onlineUsersCard);

    // Function to handle responsive layout
    function handleResponsiveLayout() {
        const categoryContainer = document.getElementById('category_container');
        const onlineUsersCard = document.getElementById('online_users_card');
        const profileCard = document.getElementById('profile_card');
        let leftCuster = document.getElementById('left_cluster');
        if (!document.getElementById('hamburger_menu')) {
        let hamburgerMenu = document.createElement('button');
        hamburgerMenu.id = 'hamburger_menu';
        hamburgerMenu.textContent = '☰';
        hamburgerMenu.addEventListener('click', () => {
            hamburgerMenu.style.display = 'none';
            renderMobileCategoriesMenu();
        });
        leftCluster.appendChild(hamburgerMenu);

        if (window.innerWidth <= 540) {
            hamburgerMenu.style.display = 'block';
        } else {
            hamburgerMenu.style.display = 'none';
        }
        }

        if (window.innerWidth <= 540) {
            // Hide elements on small screens
            if (categoryContainer) {
                categoryContainer.style.display = 'none';
            }
            if (onlineUsersCard) {
                onlineUsersCard.style.display = 'none';
            }
            if (profileCard) {
                profileCard.style.display = 'none';
            }
        } else {
            // Show elements on larger screens
            if (categoryContainer) {
                categoryContainer.style.display = '';
            }
            if (onlineUsersCard) {
                onlineUsersCard.style.display = '';
            }
            if (profileCard) {
                profileCard.style.display = '';
            }
        }
    }

    // Add resize event listener
    window.addEventListener('resize', handleResponsiveLayout);

    // Call initially to set correct state
    handleResponsiveLayout();

    // Create container for post toggling button
    let postsButtonContainer = document.createElement('div');
    postsButtonContainer.id = 'posts_button_container';
    postsButtonContainer.classList.add('posts-button-container');

    // Create buttons to toggle prefered posts
    let allPostsButton = document.createElement('button');
    let likedPostsButton = document.createElement('button');
    let myPostsButton = document.createElement('button');
    allPostsButton.classList.add('active');
    allPostsButton.id = 'all_posts_button';
    likedPostsButton.id = 'liked_posts_button';
    myPostsButton.id = 'my_posts_button';
    allPostsButton.textContent = 'All Posts';
    likedPostsButton.textContent = 'Liked Posts';
    myPostsButton.textContent = 'My Posts';

    // Append post buttons to button container
    postsButtonContainer.appendChild(allPostsButton);
    postsButtonContainer.appendChild(likedPostsButton);
    postsButtonContainer.appendChild(myPostsButton);

    // Add posts card to hold user's recent posts
    let postsCard = document.createElement('div');
    postsCard.classList.add('posts-card');
    postsCard.id = 'posts_card';

    // Add Posts Section
    let postsContainer = document.createElement('div');
    postsContainer.id = 'posts_container';
    postsCard.appendChild(postsButtonContainer);
    postsCard.appendChild(postsContainer);
    centerCluster.appendChild(postsCard);

    // Add profile container
    if (window.innerWidth > 540) {
    let profileCard = document.createElement('div');
    profileCard.id = 'profile_card';

    // Add contents of profile container
    let profileHead = document.createElement('div');
    let profileTitle = document.createElement('h2');
    let profileSubtitle = document.createElement('h2');
    let profilePic = document.createElement('img');
    let profileUserInfo = document.createElement('div');
    let profileContact = document.createElement('h3');
    let profileAgeGender = document.createElement('div');
    let profileAgeGenderSeparator = document.createElement('h3');
    let profileAge = document.createElement('h3');
    let profileGender = document.createElement('h3');
    let bioContainer = document.createElement('div');
    let bioTitle = document.createElement('h3');
    let bioText = document.createElement('p');
    let profileActionContainer = document.createElement('div');
    let messangerLink = document.createElement('a');
    let logoutLink = document.createElement('form');
    let messangerButton = document.createElement('button');
    let logoutButton = document.createElement('button');
    profileTitle.id = 'profile_title';
    profileSubtitle.id = 'profile_subtitle';
    profileContact.id = 'profile_contact';
    profileAgeGender.id = 'profile_age_gender_container';
    profileAgeGenderSeparator.id = 'profile_age_gender_separator';
    profileAge.id = 'profile_age';
    profileGender.id = 'profile_gender';
    profilePic.id = 'index_profile_pic';
    bioContainer.id = 'index_bio_container';
    bioTitle.id = 'profile_bio_title';
    bioText.id = 'profile_bio_text';
    profileAge.classList.add('profile-age-gender');
    profileGender.classList.add('profile-age-gender');
    profileHead.classList.add('profile-user-info-container');
    profileUserInfo.classList.add('profile-user-info-container');
    profileActionContainer.id = 'profile_action_container';
    messangerButton.id ='messanger_button';
    logoutButton.id = 'profile_logout_button';
    logoutButton.type = 'submit';
    messangerLink.href = '/chat';
    logoutLink.action = '/logout';
    logoutLink.method = 'POST';
    profileAgeGenderSeparator.textContent = '|';
    messangerButton.textContent = 'Messager';
    logoutButton.textContent = 'Logout';

    // Append profile components to container
    profileAgeGender.appendChild(profileGender);
    profileAgeGender.appendChild(profileAgeGenderSeparator);
    profileAgeGender.appendChild(profileAge);
    bioContainer.appendChild(bioTitle);
    bioContainer.appendChild(bioText);
    messangerLink.appendChild(messangerButton);
    logoutLink.appendChild(logoutButton);
    profileActionContainer.appendChild(messangerLink);
    profileActionContainer.appendChild(logoutLink);
    profileHead.appendChild(profileTitle);
    profileHead.appendChild(profileSubtitle);
    profileUserInfo.appendChild(profileContact);
    profileUserInfo.appendChild(profileAgeGender);
    profileCard.appendChild(profileHead);
    profileCard.appendChild(profilePic);
    profileCard.appendChild(profileUserInfo);
    profileCard.appendChild(bioContainer);
    profileCard.appendChild(profileActionContainer);
    rightCluster.appendChild(profileCard);
    }



    // Fetch home page data from server
    fetchIndexData();
}

// fetch data from response
function fetchIndexData() {
    fetch('/api/index-data', {
        method: 'GET',
        credentials: 'include',
        headers: {
            'Accept': 'application/json'
        }
    })

    .then(response => {
        if (!response.ok) {
            // Check for unauthorized resposes
            if (response.status === 401) {
                navigateTo('/login');
                return Promise.reject('Unauthorized');
            } else {
                // Render error page
                renderErrorPage(response.statusText, response.status);
                throw new Error(`HTTP error! status: ${response.status}`);
            }
        }
        return response.json();
    })

    .then(data => {
            // Extract data for rendering
            const categories = data.categories
            const allPosts = data.all_posts
            const likedPosts = data.liked_posts
            const userPosts = data.user_posts
            let user = data.user
            const userData = {
                name: `${user.first_name.charAt(0).toUpperCase()}${user.first_name.slice(1)} ${user.last_name.charAt(0).toUpperCase()}${user.last_name.slice(1)}`,
                username: user.username,
                first_name: `${user.first_name.charAt(0).toUpperCase()}${user.first_name.slice(1)}`,
                email: user.email,
                age: `${user.age} years`,
                gender: `${user.gender.charAt(0).toUpperCase()}${user.gender.slice(1)}`,
                bio: `${user.bio.charAt(0).toUpperCase()}${user.bio.slice(1)}`,
                image: `/static/images/${user.image}`
            }

            // Update profile section if profile card exists
            if (document.querySelector('#profile_card')) {
                let profileTitle = document.getElementById('profile_title');
                let profileSubtitle = document.getElementById('profile_subtitle');
                let profilePic = document.getElementById('index_profile_pic');
                let profileContact = document.getElementById('profile_contact');
                let profileGender = document.getElementById('profile_gender');
                let profileAge = document.getElementById('profile_age');
                let bioTitle = document.getElementById('profile_bio_title');
                let bioText = document.getElementById('profile_bio_text');

                // Update page with user infomation
                profileTitle.textContent = userData.name;
                profileSubtitle.textContent = `@${userData.username}`;
                profilePic.src = userData.image;
                profilePic.alt = userData.name;
                profileContact.textContent = `mail: ${userData.email}`;
                profileGender.textContent = userData.gender;
                profileAge.textContent = userData.age;
                bioTitle.textContent = `About ${userData.first_name}`;
                bioText.textContent = userData.bio;
            }

            // Render navbar
            renderNavBar(userData);

            // Render categories and populate posts
            populateCategories(categories, 'category_content_container');
            populatePosts(allPosts);
            setToggleEventListeners(allPosts, likedPosts, userPosts);

            // Initialize chat and fetch online users
            console.log("USERID: ", user.id);
            initChat(user.id);

            // Get the online users container and populate it
            const onlineUsersCard = document.querySelector('#online_users_card');
            if (onlineUsersCard) {
                // Create a container for the online users list
                const onlineUsersContent = document.createElement('div');
                onlineUsersContent.id = 'online_users_content';
                onlineUsersCard.appendChild(onlineUsersContent);

                // Fetch users and populate the list
                fetchAllUsers();
            }
    })

    .catch(error => {
            console.error('Error fetching home page:', error);
    });
}