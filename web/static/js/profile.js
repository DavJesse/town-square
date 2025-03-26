import { renderNavBar } from '/static/js/navbar.js';
import { renderLoginPage } from '/static/js/login.js';

export function renderProfilePage() {    
    // Add container to hold left profile page content
    let profilePage = document.createElement('div');
    profilePage.classList.add('profile-page');
    profilePage.id = 'profile_page';
    
    // Add container to hold left cluster of page
    let leftCluster = document.createElement('div');
    leftCluster.classList.add('left-cluster');
    leftCluster.id = 'left_cluster';
    
    // Add card to hold online users
    let onlineUsersCard = document.createElement('div');
    onlineUsersCard.classList.add('online-users-card');
    onlineUsersCard.id = 'online_users_card';
    
    // Add title for online users card
    let onlineUsersTitle = document.createElement('h3');
    onlineUsersTitle.id = 'online_users_title';
    onlineUsersTitle.textContent = 'who\'s online?';
    onlineUsersCard.appendChild(onlineUsersTitle);
    
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
    gender.id = 'gender';
    
    // Add user age for bio card
    let age = document.createElement('p');
    age.id = 'age';
    
    // Add bio paragraph
    let bioParagraph = document.createElement('p');
    bioParagraph.id = 'bio_paragraph';
    
    // Add posts card to hold user's recent posts
    let postsCard = document.createElement('div');
    postsCard.classList.add('posts-card');
    postsCard.id = 'posts_card';
    
    // Add posts card to hold user's recent posts
    let postsButtonContainer = document.createElement('div');
    postsButtonContainer.classList.add('posts-button-container');
    postsButtonContainer.id = 'posts_button_container';
    
    // Create buttons to toggle prefered posts
    let myPostsLink = document.createElement('a');
    myPostsLink.href = '/my-posts';
    myPostsLink.textContent = 'my posts';
    myPostsLink.id = 'my_posts_link';
    
    let myPostsButton = document.createElement('button');
    myPostsButton.classList.add('posts-button');
    myPostsButton.id = 'my_posts_button';
    myPostsButton.appendChild(myPostsLink);
    postsButtonContainer.appendChild(myPostsButton);
    
    let likedPostsLink = document.createElement('a');
    likedPostsLink.href = '/liked-posts';
    likedPostsLink.textContent = 'posts I\'ve liked';
    likedPostsLink.id = 'liked_posts_link';
    
    let likedPostsButton = document.createElement('button');
    likedPostsButton.classList.add('posts-button');
    likedPostsButton.id = 'liked_posts_button';
    likedPostsButton.appendChild(likedPostsLink);
    postsButtonContainer.appendChild(likedPostsButton);    
    postsCard.appendChild(postsButtonContainer);    
    
    // Add Posts Section            
    let postsContainer = document.createElement('div');
    postsContainer.id = 'posts_container';
    postsCard.appendChild(postsContainer);
    
    fetch('/profile', {
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
            const posts = data.data.posts;
            document.title = user.username;
            
            // Populate bioCard with user data
            bioTitle.textContent = `${user.first_name} ${user.last_name}`;
            nickname.textContent = `@${user.username}`;
            email.textContent = user.email;
            gender.textContent = user.gender;
            age.textContent = user.age;
            bioParagraph.textContent = user.bio;
            profilePic.src = user?.image?`/static/images/${user.image}` : '/static/user-circle-svgrepo-com.svg';
            
            // Append user details to bioCard
            bioCard.appendChild(bioTitle);
            bioCard.appendChild(nickname);
            bioCard.appendChild(email);
            bioCard.appendChild(gender);
            bioCard.appendChild(age);
            bioCard.appendChild(bioParagraph);
            bioCard.appendChild(profilePic);            
            
            // Add Posts
            if (posts?.length > 0) {
                posts.forEach(post => {
                    let postElement = document.createElement('div');
                    postElement.classList.add('card');
                    
                    let postTitle = document.createElement('p');
                    postTitle.classList.add('card__title');
                    postTitle.id = 'post_title';
                    postTitle.textContent = post.title;
                    postElement.appendChild(postTitle);
                    
                    let postContent = document.createElement('p');
                    postContent.classList.add('card__description');
                    postContent.id = 'post_content';
                    postContent.textContent = post.content;
                    postElement.appendChild(postContent);
                    
                    if (post.media) {
                        let postMedia = document.createElement('img');
                        postMedia.id = 'post_media';
                        postMedia.src = `/static/media/${post.media}`;
                        postElement.appendChild(postMedia);
                    }
                    postsContainer.appendChild(postElement);
                });
            } else {
                let noPostsAlerts = document.createElement('p');
                noPostsAlerts.id = 'no_posts';
                noPostsAlerts.textContent = 'No posts available';
                postsContainer.appendChild(noPostsAlerts);
            }
        } else {
            console.error('Error fetching profile data:', data.message);
        }
    })
    .catch(error => {
        console.error('Error:', error);
    });

    // Append elements to left cluster
    leftCluster.appendChild(onlineUsersCard);
    
    // Append elements to center cluster
    centerCluster.appendChild(bioCard);
    centerCluster.appendChild(postsCard);
    
    // Append page elements to app
    let app = document.getElementById("app")
    app.innerHTML = ""
    renderNavBar();
    profilePage.appendChild(leftCluster);
    profilePage.appendChild(centerCluster);
    profilePage.appendChild(rightCluster)
    app.appendChild(profilePage);
    
    // Add script tag and link to js
    let script = document.createElement('script');
    script.src = '/static/js/index.js';
    script.type = 'module';
    script.defer = true;
    document.body.appendChild(script);
}