import { renderNavBar } from '/static/js/navbar.js'
import { renderCreatePostButton, setCreatePostsButtonListeners } from '/static/js/create_post_button.js';
import { navigateTo } from '/static/js/routes.js';
import { populatePosts, setToggleEventListeners } from '/static/js/populate_posts.js';
import { renderErrorPage } from '/static/js/error.js';
import { populateCategories } from '/static/js/populate_categories.js';

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
    let profileTitle = document.createElement('h2');
    let profilePic = document.createElement('img');
    let bioContainer = document.createElement('div');
    let bioTitle = document.createElement('h3');
    let bioText = document.createElement('p');
    let profileActionContainer = document.createElement('div');
    let profileLink = document.createElement('a');
    let messangerLink = document.createElement('a');
    let logoutLink = document.createElement('form');
    let viewProfileButton = document.createElement('button');
    let messangerButton = document.createElement('button');
    let logoutButton = document.createElement('button');
    profileTitle.id = 'profile_title';
    profilePic.id = 'index_profile_pic';
    bioContainer.id = 'index_bio_container';
    bioTitle.id = 'profile_bio_title';
    bioText.id = 'profile_bio_text';
    profileActionContainer.id = 'profile_action_container';
    viewProfileButton.id = 'view_profile_button';
    messangerButton.id ='messanger_button';
    logoutButton.id = 'profile_logout_button';
    logoutButton.type = 'submit';
    profileLink.href = '/profile';
    messangerLink.href = '/chat';
    logoutLink.action = '/logout';
    logoutLink.method = 'POST';
    viewProfileButton.textContent = 'View Profile';
    messangerButton.textContent = 'Messager';
    logoutButton.textContent = 'Logout';

    // Append profile components to container
    bioContainer.appendChild(bioTitle);
    bioContainer.appendChild(bioText);
    profileLink.appendChild(viewProfileButton);
    messangerLink.appendChild(messangerButton);
    logoutLink.appendChild(logoutButton);
    profileActionContainer.appendChild(profileLink);
    profileActionContainer.appendChild(messangerLink);
    profileActionContainer.appendChild(logoutLink);
    profileCard.appendChild(profileTitle);
    profileCard.appendChild(profilePic);
    profileCard.appendChild(bioContainer);
    profileCard.appendChild(profileActionContainer);
    rightCluster.appendChild(profileCard);
    
    // fetch data from response
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
                // throw new Error(`HTTP error! status: ${response.status}`);
                return response.json().then(errorData => {
                    renderErrorPage(errorData.Issue || response.statusText, response.status);
                    return Promise.reject(errorData.Issue || 'Error occurred');
                });
            }
        }
        return response.json();
    })
    
    .then(data => {
        if (data.code === 200) {

            // Extract data for rendering
            const categories = data.data.categories
            const posts = data.data.posts
            const likedPosts = data.data.liked_posts
            const user = data.data.user

            // Update page with user infomation
            profileTitle.textContent = `${user.first_name.charAt(0).toUpperCase()}${user.first_name.slice(1)} ${user.last_name.charAt(0).toUpperCase()}${user.last_name.slice(1)}`;
            profilePic.src = `/static/images/${user.image}`;
            profilePic.alt = `${user.first_name.charAt(0).toUpperCase()}${user.first_name.slice(1)} ${user.last_name.charAt(0).toUpperCase()}${user.last_name.slice(1)} image`
            bioTitle.textContent = `About ${user.first_name.charAt(0).toUpperCase()}${user.first_name.slice(1)}`;
            bioText.textContent = `${user.bio.charAt(0).toUpperCase()}${user.bio.slice(1)}`;
            viewProfileButton.href = `/profile`;
            messangerButton.href = `/messages`;

            // Render categories and populate posts
            populateCategories(categories);
            populatePosts(posts);
            setToggleEventListeners(posts, likedPosts, []);

        } else if (data.code === 401) {
            navigateTo('/login');
        } else {
            renderErrorPage(data.message, data.code);
            console.error('Error fetching home page data:', data.message);
        }
    })

    .catch(error => {
        if (error !== 'Unauthorized') {
            console.error('Error fetching home page:', error);
        }
    });
}
