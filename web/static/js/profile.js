import { renderNavBar } from '/static/js/navbar.js';

export function renderProfilePage() {
    renderNavBar();
    let profileBody = document.createElement('div');
    profileBody.classList.add('container');

    let userBody = document.createElement('div');
    userBody.classList.add('main-content');
    userBody.id = 'main_content';
    profileBody.appendChild(userBody);

    let app = document.getElementById("app")
    app.appendChild(profileBody);
    
    fetch('/profile', {
        method: 'GET',
        credentials: 'include',
        headers: {
            'Accept': 'application/json'
        }
    })
    .then(response => response.json())
    .then(data => {
        console.log('Profile data:', data);
        if (data.code === 200) {
            const user = data.data.User;
            const posts = data.data.Posts;
            document.title = user.username;

            const mainContent = document.getElementById('main_content');
            mainContent.innerHTML = ''; // Clear existing content

            const profileHeading = document.createElement('h3');
            profileHeading.classList.add('main-content__heading');
            profileHeading.id = 'profile-heading';
            profileHeading.textContent = `${user.username}'s Profile`;
            mainContent.appendChild(profileHeading);

            // Add Profile Details
            const profileCard = document.createElement('div');
            profileCard.classList.add('card');
            profileCard.id = 'profile-card';

            const profileInfo = document.createElement('div');
            profileInfo.classList.add('profile-info');

            const profileImage = document.createElement('img');
            profileImage.id = 'profile-image';
            profileImage.src = user?.image ? `/static/images/${user.image}` : '/static/user-circle-svgrepo-com.svg';
            profileInfo.appendChild(profileImage);

            const profileDetails = document.createElement('div');
            profileDetails.classList.add('profile-details');

            const emailParagraph = document.createElement('p');
            emailParagraph.innerHTML = `<strong>Email:</strong> <span id="user-email">${user.email}</span>`;
            profileDetails.appendChild(emailParagraph);

            const bioParagraph = document.createElement('p');
            bioParagraph.innerHTML = `<strong>Bio:</strong> <span id="user-bio">${user?.bio || 'No Bio'}</span>`;
            profileDetails.appendChild(bioParagraph);

            const likedPostsLink = document.createElement('a');
            likedPostsLink.href = '/liked-posts';
            likedPostsLink.innerHTML = `<button class="navbar__button">View posts I've liked</button>`;
            profileDetails.appendChild(likedPostsLink);

            profileInfo.appendChild(profileDetails);
            profileCard.appendChild(profileInfo);
            mainContent.appendChild(profileCard);

            // Add Posts Section
            const postsHeading = document.createElement('h3');
            postsHeading.classList.add('main-content__heading');
            postsHeading.textContent = 'Posts';
            mainContent.appendChild(postsHeading);

            const postsContainer = document.createElement('div');
            postsContainer.id = 'posts-container';
            mainContent.appendChild(postsContainer);

            // Add Posts
            if (posts?.length > 0) {
                posts.forEach(post => {
                    const postElement = document.createElement('div');
                    postElement.classList.add('card');
                    postElement.innerHTML = `
                        <p class="card__title">
                            <a href="#" onclick="fetchPostDisplay('${post.uuid}')" style="color: #0172eb;">
                                ${post.title}
                            </a>
                        </p>
                        <p class="card__description">${post.content}</p>
                        ${post.Media ? `<img src="/static/media/${post.media}" alt="Post Image" style="width: 70px; aspect-ratio: 1/1;">` : ''}
                    `;
                    postsContainer.appendChild(postElement);
                });
            } else {
                postsContainer.innerHTML = '<p>No posts available</p>';
            }
        } else {
            console.error('Error fetching profile data:', data.message);
        }
    })
    .catch(error => {
        console.error('Error:', error);
    });
}