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

    indexPageContainer.appendChild(rightCluster);
    indexPageContainer.appendChild(centerCluster);
    indexPageContainer.appendChild(leftCluster);

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
    indexPageContainer.appendChild(leftCluster);

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
