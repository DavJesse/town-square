import { renderNavBar } from '/static/js/navbar.js'
import { renderCreatePostButton, setCreatePostsButtonListeners } from '/static/js/create_post_button.js';
import { navigateTo } from '/static/js/routes.js';
import { populatePosts, setToggleEventListeners } from '/static/js/populate_posts.js';
import { renderErrorPage } from '/static/js/error.js';
import { populateCategories } from '/static/js/populate_categories.js';
import { initChat, populateOnlineUsersList, fetchAllUsers } from '/static/js/chat.js';

export function renderIndexPage() {
    // Extract app from dom
    let app = document.getElementById('app');
    app.innerHTML = "";

    // Set document title
    document.title = 'real-time-forum';

    // Render navbar
    renderNavBar();

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
    let leftCluster = document.createElement('div');
    let centerCluster = document.createElement('div');
    let rightCluster = document.createElement('div');
    leftCluster.id = 'left_cluster';
    centerCluster.id = 'center_cluster';
    rightCluster.id = 'right_cluster';

    indexPageContainer.appendChild(leftCluster);
    indexPageContainer.appendChild(centerCluster);
    indexPageContainer.appendChild(rightCluster);

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

    // Create container for post toggling button
    let postsButtonContainer = document.createElement('div');
    postsButtonContainer.id = 'posts_button_container';

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

    // Add prifile container
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
            const user = data.user
            let profileTitle = document.getElementById('profile_title');
            let profileSubtitle = document.getElementById('profile_subtitle');
            let profilePic = document.getElementById('index_profile_pic');
            let profileContact = document.getElementById('profile_contact');
            let profileGender = document.getElementById('profile_gender');
            let profileAge = document.getElementById('profile_age');
            let bioTitle = document.getElementById('profile_bio_title');
            let bioText = document.getElementById('profile_bio_text');

            // Update page with user infomation
            profileTitle.textContent = `${user.first_name.charAt(0).toUpperCase()}${user.first_name.slice(1)} ${user.last_name.charAt(0).toUpperCase()}${user.last_name.slice(1)}`;
            profileSubtitle.textContent = `@${user.username}`;
            profilePic.src = `/static/images/${user.image}`;
            profilePic.alt = `${user.first_name.charAt(0).toUpperCase()}${user.first_name.slice(1)} ${user.last_name.charAt(0).toUpperCase()}${user.last_name.slice(1)} image`;
            profileContact.textContent = `contact: ${user.email}`;
            profileGender.textContent = `${user.gender.charAt(0).toUpperCase()}${user.gender.slice(1)}`;
            profileAge.textContent = `${user.age} years old`;
            bioTitle.textContent = `About ${user.first_name.charAt(0).toUpperCase()}${user.first_name.slice(1)}`;
            bioText.textContent = `${user.bio.charAt(0).toUpperCase()}${user.bio.slice(1)}`;

            // Render categories and populate posts
            populateCategories(categories);
            populatePosts(allPosts);
            setToggleEventListeners(allPosts, likedPosts, userPosts);

            // Initialize chat and fetch online users
            console.log("USERID: ", user.id);
            initChat(user.id);

            // Get the online users container and populate it
            const onlineUsersCard = document.getElementById('online_users_card');
            if (onlineUsersCard) {
                // Create a container for the online users list
                const onlineUsersContent = document.createElement('div');
                onlineUsersContent.id = 'online_users_content';
                onlineUsersCard.appendChild(onlineUsersContent);

                // Add some CSS styles for the online users list
                const style = document.createElement('style');
                style.textContent = `
                    #online_users_content {
                        padding: 10px;
                        max-height: 300px;
                        overflow-y: auto;
                    }
                    .compact-user-item {
                        display: flex;
                        align-items: center;
                        padding: 8px;
                        cursor: pointer;
                        border-radius: 4px;
                        margin-bottom: 4px;
                        transition: background-color 0.2s;
                    }
                    .compact-user-item:hover {
                        background-color: #f0f0f0;
                    }
                    .user-section-header {
                        font-weight: bold;
                        margin-top: 8px;
                        margin-bottom: 4px;
                        color: #555;
                        font-size: 0.9em;
                    }
                    .status-indicator {
                        width: 8px;
                        height: 8px;
                        border-radius: 50%;
                        margin-right: 8px;
                    }
                    .status-online {
                        background-color: #4CAF50;
                    }
                    .status-offline {
                        background-color: #9e9e9e;
                    }
                    .user-nickname {
                        font-size: 0.9em;
                    }
                    .notification-badge {
                        background-color: #f44336;
                        color: white;
                        border-radius: 50%;
                        padding: 2px 6px;
                        font-size: 0.7em;
                        margin-left: auto;
                    }
                    .show-more-btn {
                        text-align: center;
                        color: #2196F3;
                        font-size: 0.8em;
                        padding: 5px;
                        cursor: pointer;
                        margin-top: 5px;
                    }
                    .show-more-btn:hover {
                        text-decoration: underline;
                    }
                    .empty-user-item {
                        color: #757575;
                        font-style: italic;
                        padding: 8px;
                        font-size: 0.9em;
                    }
                `;
                document.head.appendChild(style);

                // Fetch users and populate the list
                fetchAllUsers();
            }
    })

    .catch(error => {
            console.error('Error fetching home page:', error);
    });
}