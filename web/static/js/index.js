import { renderNavBar } from '/static/js/navbar.js'
import { renderCreatePostButton } from '/static/js/create_post.js'
import { setCreatePostsButtonListeners } from '/static/js/create_post.js';

export function renderIndexPage() {
    // Extract app from dom
    let app = document.getElementById('app');
    app.innerHTML = "";

    // Render navbar
    renderNavBar();

    // Render create post button
    let userOptions = document.createElement('div');
    userOptions.id = 'user_options';

    // Render create post button and append to app
    let createPostButton = renderCreatePostButton();
    setCreatePostsButtonListeners(createPostButton);
    userOptions.appendChild(createPostButton);
    app.appendChild(userOptions);

    // Create index page container
    let indexPageContainer = document.createElement('div');
    indexPageContainer.id = 'index_page'
    app.appendChild(indexPageContainer);

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
    categoriesCard.id = 'category_container';
    onlineUsersCard.id = 'online_users_card';
    onlineUsersTitle.id = 'online_users_title';
    onlineUsersTitle.textContent = 'Who\'s online?';

    onlineUsersCard.appendChild(onlineUsersTitle);
    leftCluster.appendChild(categoriesCard);
    leftCluster.appendChild(onlineUsersCard);
   
    // Create container for post toggling button
    let postsButtonContainer = document.createElement('div');
    postsButtonContainer.id = 'posts_button_container';

    // Create buttons to toggle prefered posts    
    let myPostsButton = document.createElement('button');
    let likedPostsButton = document.createElement('button');
    myPostsButton.classList.add('active');
    myPostsButton.id = 'my_posts_button';
    likedPostsButton.id = 'liked_posts_button';
    myPostsButton.textContent = 'My Posts';
    likedPostsButton.textContent = 'Liked Posts'

    // Append post buttons to button container
    postsButtonContainer.appendChild(myPostsButton);
    postsButtonContainer.appendChild(likedPostsButton);

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
    let viewProfileButton = document.createElement('button');
    let messangerButton = document.createElement('button');
    let logoutButton = document.createElement('button');
    profileTitle.id = 'profile_title';
    profilePic.id = 'profile_pic';
    bioContainer.id = 'bio_container';
    bioTitle.id = 'profile_bio_title';
    bioText.id = 'profile_bio_text';
    profileActionContainer.id = 'profile_action_container';
    viewProfileButton.id = 'view_profile_button';
    messangerButton.id ='messanger_button';
    logoutButton.id = 'profile_logout_button';
    viewProfileButton.textContent = 'View Profile';
    messangerButton.textContent = 'Messager';
    logoutButton.textContent = 'Logout';

    // Append profile components to container
    bioContainer.appendChild(bioTitle);
    bioContainer.appendChild(bioText);
    profileActionContainer.appendChild(viewProfileButton);
    profileActionContainer.appendChild(messangerButton);
    profileActionContainer.appendChild(logoutButton);
    profileCard.appendChild(profileTitle);
    profileCard.appendChild(profilePic);
    profileCard.appendChild(bioContainer);
    profileCard.appendChild(profileActionContainer);
    rightCluster.appendChild(profileCard);
    
}

// window.fetchAllPosts = function(){
//     fetch('/posts')  // Assuming '/posts' returns a JSON array of posts
//         .then(response => response.json())
//         .then(data => {
//             if (data.code === 200 && data.data) {
//                 renderPage(data.data);
//             } else if (data.redirect) {
//                 window.location.href = data.redirect;
//             } else {
//                 console.error('Error fetching posts:', data);
//                 document.getElementById('main-content').innerHTML = '<p>Error loading posts.</p>';
//             }
//         })
//         .catch(error => {
//             console.error('Error fetching posts:', error);
//             document.getElementById('main-content').innerHTML = '<p>Error loading posts.</p>';
//         });
// }
