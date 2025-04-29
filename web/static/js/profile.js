import { renderNavBar } from '/static/js/navbar.js';
import { renderLoginPage } from '/static/js/login.js';
import { handleLikePost } from '/static/js/like_post.js';
import { handleDislikePost } from '/static/js/dislike_post.js';
import { populateCategories } from '/static/js/index.js';
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
                renderLoginPage();  // Redirect to login page
                return;
            }
            throw new Error(`HTTP error! Status: ${response.status}`);
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

        } else if (data.code === 401) {
            navigateTo('/login');
        } else {
            renderErrorPage(data.message, data.code);
            console.error('Error fetching home page data:', data.message);
        }
    })
    .catch(error => {
        console.error('Error:', error);
    });   
}

export function populatePosts(posts) {
    let postsContainer = document.getElementById('posts_container');
    postsContainer.innerHTML = ""; // clear content for fresh population

    // Set null values to empty array
    posts = posts ?? [];

    if (posts && posts.length > 0) {
        posts.forEach(post => {
            let postElement = document.createElement('div');
            postElement.id = 'post_card';

            let postHead = document.createElement('div');
            postHead.id ='post_head';
            postElement.appendChild(postHead);

            let postCreatorPic = document.createElement('img');
            postCreatorPic.id = 'post_creator_pic';
            postCreatorPic.src = post.creator_image? `/static/images/${post.creator_image}` : '/static/user-circle-svgrepo-com.svg';
            postHead.appendChild(postCreatorPic);

            let postCreatorInfoContainer = document.createElement('div');
            postCreatorInfoContainer.id = 'post_creator_info_container';
            postHead.appendChild(postCreatorInfoContainer);

            let postCreatorName = document.createElement('h3');
            postCreatorName.id = 'post_creator_name';
            postCreatorName.textContent = post.creator_first_name? `${post.creator_first_name.charAt(0).toUpperCase()}${post.creator_first_name.slice(1)} ${post.creator_last_name.charAt(0).toUpperCase()}${post.creator_last_name.slice(1)}` : post.creator_username;
            postCreatorInfoContainer.appendChild(postCreatorName);

            let postCreatorUsername = document.createElement('h4');
            postCreatorUsername.id = 'post_creator_username';
            postCreatorUsername.textContent = `@${post.creator_username}`;
            postCreatorInfoContainer.appendChild(postCreatorUsername);

            let postContentContainer = document.createElement('div');
            postContentContainer.id = 'post_content_container';
            postElement.appendChild(postContentContainer);

            let postTitle = document.createElement('h4');
            postTitle.id = 'post_title';
            postTitle.textContent = post.title;
            postContentContainer.appendChild(postTitle);

            let postContent = document.createElement('p');
            postContent.id = 'post_content';
            postContent.textContent = post.content;
            postContentContainer.appendChild(postContent);

            if (post.media) {
                let postMedia = document.createElement('img');
                postMedia.id = 'post_media';
                postMedia.src = `/static/media/${post.media}`;
                postElement.appendChild(postMedia);
            }

            let postCreationDate = document.createElement('p');
            let options = { year: 'numeric', month: 'long', day: 'numeric', hour: '2-digit', minute: '2-digit' };
            postCreationDate.id = 'post_creation_date';
            postCreationDate.textContent = `Posted on ${new Date(post.created_at).toLocaleString('en-GB', options)}`;
            postElement.appendChild(postCreationDate);

            let postEngagement = document.createElement('div');
            postEngagement.id = 'post_engagement';
            postElement.appendChild(postEngagement);

            // Create engagement section
            let likeContainer = document.createElement('div');
            let dislikeContainer = document.createElement('div');
            let commentContainer = document.createElement('div');
            let likeLink = document.createElement('a');
            let dislikeLink = document.createElement('a');
            let commentLink = document.createElement('a');
            let likeCount = document.createElement('p');
            let dislikeCount = document.createElement('p');
            let commentCount = document.createElement('p');
            let likeIcon = document.createElement('span');
            let dislikeIcon = document.createElement('span');
            let commentIcon = document.createElement('span');

            // Set IDs and classes
            likeContainer.id = 'like_container';
            dislikeContainer.id = 'dislike_container';
            commentContainer.id = 'comment_container';
            likeLink.id = 'like_link';
            likeLink.href = '#';
            likeLink.dataset.postId = post.uuid;  // Add post ID as data attribute
            dislikeLink.id = 'dislike_link';
            dislikeLink.href = '#';
            dislikeLink.dataset.postId = post.uuid;  // Change this to prevent page refresh
            commentLink.href = '#';   // Change this to prevent page refresh
            likeCount.id = 'engagement_count';
            dislikeCount.id = 'engagement_count';
            commentCount.id = 'engagement_count';
            likeIcon.classList.add('material-symbols-outlined');
            dislikeIcon.classList.add('material-symbols-outlined');
            commentIcon.classList.add('material-symbols-outlined');
            likeIcon.id = 'engagement_icon';
            dislikeIcon.id = 'engagement_icon';
            commentIcon.id = 'engagement_icon';

            // Set content
            likeCount.textContent = post.likes_count;
            dislikeCount.textContent = post.dislikes_count;
            commentCount.textContent = post.comments.length;
            likeIcon.textContent = 'thumb_up';
            dislikeIcon.textContent = 'thumb_down';
            commentIcon.textContent = 'comment';

            // Add event listeners
            likeLink.addEventListener('click', (e) => {
                e.preventDefault();
                handleLikePost(post.uuid, likeContainer);
            });

            dislikeLink.addEventListener('click', (e) => {
                e.preventDefault();
                handleDislikePost(post.uuid, dislikeContainer);
            });

            commentLink.addEventListener('click', (e) => {
                e.preventDefault();
                // Add your comment handling here
            });

            // Assemble the structure
            postEngagement.appendChild(likeContainer);
            postEngagement.appendChild(dislikeContainer);
            postEngagement.appendChild(commentContainer);
            likeContainer.appendChild(likeCount);
            likeContainer.appendChild(likeLink);
            likeLink.appendChild(likeIcon);
            dislikeContainer.appendChild(dislikeCount);
            dislikeContainer.appendChild(dislikeLink);
            dislikeLink.appendChild(dislikeIcon);
            commentContainer.appendChild(commentCount);
            commentContainer.appendChild(commentLink);
            commentLink.appendChild(commentIcon);

            postsContainer.appendChild(postElement);
        });
    } else {
        let noPostsContainer = document.createElement('div');
        noPostsContainer.id = 'no_posts_container';

        let noPostsAlerts = document.createElement('p');
        noPostsAlerts.id = 'no_posts';
        noPostsAlerts.textContent = 'No posts available';
        noPostsContainer.appendChild(noPostsAlerts);

        postsContainer.appendChild(noPostsContainer);
    }
}

export function setToggleEventListeners(userPosts, likedPosts) {
    // Extract page elements
    let myPostsButton = document.getElementById('my_posts_button');
    let likedPostsButton = document.getElementById('liked_posts_button');

    // Set null values to empty array
    userPosts = userPosts ?? [];
    likedPosts = likedPosts ?? [];

    myPostsButton.addEventListener('click', () => {
        myPostsButton.classList.add('active');
        likedPostsButton.classList.remove('active');
        populatePosts(userPosts);
    });

    likedPostsButton.addEventListener('click', () => {
        likedPostsButton.classList.add('active');
        myPostsButton.classList.remove('active');
        populatePosts(likedPosts);
    });
}
