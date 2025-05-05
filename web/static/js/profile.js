import { renderNavBar } from '/static/js/navbar.js';
import { navigateTo } from '/static/js/routes.js';
import { renderErrorPage } from '/static/js/error.js';
import { populateCategories } from '/static/js/populate_categories.js';
import { populatePosts, setToggleEventListeners } from '/static/js/populate_posts.js';
import { renderCreatePostButton, setCreatePostsButtonListeners } from '/static/js/create_post_button.js';
import { renderLogoutButton, setLogoutButtonListeners } from '/static/js/logout_button.js';

export function renderProfilePage() {
    // Append page elements to app
    let app = document.getElementById("app")
    app.innerHTML = ""
    renderNavBar();
    
    // Add container to hold left profile page content
    let profilePage = document.createElement('div');
    profilePage.classList.add('profile-page');
    profilePage.id = 'profile_page';
    app.appendChild(profilePage);
    
    // Add container to hold left cluster of page
    let leftCluster = document.createElement('div');
    leftCluster.classList.add('left-cluster');
    leftCluster.id = 'left_cluster';

    // Create categories and online users containers
    let categoriesCard = document.createElement('div');
    let categoriesCardTitle = document.createElement('h3');
    let categoryContentContainer = document.createElement('div');
    categoriesCard.id = 'category_container';
    categoriesCardTitle.id = 'categories_title';
    categoryContentContainer.id = 'category_content_container';
    categoriesCardTitle.textContent = 'Categories';

    categoriesCard.appendChild(categoriesCardTitle);
    categoriesCard.appendChild(categoryContentContainer);
    leftCluster.appendChild(categoriesCard);
    
    // Add card to hold online users
    let onlineUsersCard = document.createElement('div');
    onlineUsersCard.classList.add('online-users-card');
    onlineUsersCard.id = 'online_users_card';
    
    // Add title for online users card
    let onlineUsersTitle = document.createElement('h3');
    onlineUsersTitle.id = 'online_users_title';
    onlineUsersTitle.textContent = 'who\'s online?';

    // Append ellements of left luster to left cluster and eventually to profile page
    onlineUsersCard.appendChild(onlineUsersTitle);
    leftCluster.appendChild(onlineUsersCard);
    profilePage.appendChild(leftCluster);
    
    // Add container for center cluster
    let centerCluster = document.createElement('div');
    centerCluster.classList.add('center-cluster');
    centerCluster.id = 'center_cluster';
    
    // Add container for right cluster
    let rightCluster = document.createElement('div');
    rightCluster.classList.add('right-cluster');
    rightCluster.id = 'right_cluster';
    
    // Add bio card to hold user's basic information
    let bioCard = document.createElement('div');
    bioCard.classList.add('bio-card');
    bioCard.id = 'bio_card';

    // Add container to hold user's image information
    let imageContainer = document.createElement('div');
    imageContainer.classList.add('image-container');
    imageContainer.id = 'image_container';

    // Add container to hold user's text information
    let profileInfoContainer = document.createElement('div');
    profileInfoContainer.classList.add('profile-info-container');
    profileInfoContainer.id = 'profile_info_container';
    
    // Add user image for bio card
    let profilePic = document.createElement('img');
    profilePic.id = 'profile_pic';
    profilePic.alt = 'profile image';
    
    // Add title for bio card
    let bioTitle = document.createElement('h3');
    bioTitle.id = 'bio_title';
    
    // Add user nickname for bio card
    let nickname = document.createElement('p');
    nickname.id = 'nickname';
    
    // Add user email for bio card
    let email = document.createElement('p');
    email.id = 'email';
    
    // Add user gender for bio card
    let gender = document.createElement('p');
    gender.id = 'gender_text';
    
    // Add user age for bio card
    let age = document.createElement('p');
    age.id = 'age';

    // Add container for bio paragraph
    let bioContainer = document.createElement('div');
    bioContainer.id = 'bio_container';

    // Add title for bio container
    let bioContainerTitle = document.createElement('h3');
    bioContainerTitle.id = 'bio_container_title';
    
    // Add bio paragraph
    let bioParagraph = document.createElement('p');
    bioParagraph.id = 'bio_paragraph';

    // Append user details
    bioContainer.appendChild(bioContainerTitle);
    bioContainer.appendChild(bioParagraph);
    profileInfoContainer.appendChild(bioTitle);
    profileInfoContainer.appendChild(nickname);
    profileInfoContainer.appendChild(email);
    profileInfoContainer.appendChild(gender);
    profileInfoContainer.appendChild(age);
    imageContainer.appendChild(profilePic);
    bioCard.appendChild(imageContainer);
    bioCard.appendChild(bioContainer);
    imageContainer.appendChild(profileInfoContainer);
    rightCluster.appendChild(bioCard);
    
    // Add posts card to hold user's recent posts
    let postsCard = document.createElement('div');
    postsCard.classList.add('posts-card');
    postsCard.id = 'posts_card';
    
    // Add posts card to hold user's recent posts
    let postsButtonContainer = document.createElement('div');
    postsButtonContainer.classList.add('posts-button-container');
    postsButtonContainer.id = 'posts_button_container';
    
    // Create buttons to toggle prefered posts    
    let myPostsButton = document.createElement('button');
    myPostsButton.classList.add('active');
    myPostsButton.id = 'my_posts_button';
    myPostsButton.textContent = 'My Posts';
    
    let likedPostsButton = document.createElement('button');
    likedPostsButton.textContent = 'Liked Posts'
    likedPostsButton.id = 'liked_posts_button';
    
    postsButtonContainer.appendChild(myPostsButton);
    postsButtonContainer.appendChild(likedPostsButton);
    
    // Add Posts Section            
    let postsContainer = document.createElement('div');
    postsContainer.id = 'posts_container';
    postsCard.appendChild(postsButtonContainer);    
    postsCard.appendChild(postsContainer);
    centerCluster.appendChild(postsCard);

    // Add container for user options
    let userOptions = document.createElement('div');
    userOptions.id = 'user_options';

    // Render create post buttom and related dependancies
    let createPostButton = renderCreatePostButton();
    setCreatePostsButtonListeners(createPostButton);
    userOptions.appendChild(createPostButton);

    // Render logout button and related dependancies
    let logoutButton = renderLogoutButton();
    setLogoutButtonListeners(logoutButton);
    userOptions.appendChild(logoutButton);

    profilePage.appendChild(centerCluster);
    profilePage.appendChild(rightCluster)
    app.appendChild(userOptions);
    
    
    fetch('/api/profile-data', {
        method: 'GET',
        credentials: 'include',
        headers: {
            'Accept': 'application/json'
        }
    })
    
    .then(response => {
        if (!response.ok) { 
            // If the response is not OK, check if it's a 401 (Unauthorized)
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
            
            // Extract data
            const user = data.data.user;
            const userPosts = data.data.user_posts;
            const likedPosts = data.data.liked_posts;
            const categories = data.data.categories;

            // Update document title to reflect user's profile
            document.title = `Profile: ${user.first_name.charAt(0).toUpperCase()}${user.first_name.slice(1)} ${user.last_name.charAt(0).toUpperCase()}${user.last_name.slice(1)}`;

            // Populate user bio
            bioTitle.textContent = `${user.first_name.charAt(0).toUpperCase()}${user.first_name.slice(1)} ${user.last_name.charAt(0).toUpperCase()}${user.last_name.slice(1)}`;
            nickname.textContent = `@${user.username}`;
            email.textContent = `${user.email}`;
            gender.textContent = user.gender;
            age.textContent = `${user.age} years old`;
            bioParagraph.textContent = user.bio;
            profilePic.src = user?.image ? `/static/images/${user.image}` : '/static/user-circle-svgrepo-com.svg';
            
            // Load 'My Posts' section by default
            populateCategories(categories);
            populatePosts(userPosts);

            // Listen for client activity, sets posts to client preference
            setToggleEventListeners(userPosts, likedPosts);

        } else {
            console.error('Error fetching home page data:', data.message);
        }
    })
    .catch(error => {
        console.error('Error:', error);
    });   
}
