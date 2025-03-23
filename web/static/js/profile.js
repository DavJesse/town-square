import { renderNavBar } from '/static/js/navbar.js';
export function renderProfilePage() {
    renderNavBar();
    fetch('/profile', {
        method: 'GET',
        credentials: 'include' // ensure cookies are sent with the request
    })
    .then(response => response.json())
    .then(data => {
            if (data.code === 200) {
                const user = data.data.User;
                const posts = data.data.Posts;

                // Clear existing content
                const mainContent = document.getElementById('main-content');
                mainContent.innerHTML = '';

                // Create profile section dynamically
                const profileHeading = document.createElement('h3');
                profileHeading.classList.add('main-content__heading');
                profileHeading.id = 'profile-heading';
                profileHeading.textContent = `${user.username}'s Profile`;
                mainContent.appendChild(profileHeading);

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
                const emailStrong = document.createElement('strong');
                emailStrong.textContent = 'Email:';
                const userEmailSpan = document.createElement('span');
                userEmailSpan.id = 'user-email';
                userEmailSpan.textContent = user.email;
                emailParagraph.appendChild(emailStrong);
                emailParagraph.appendChild(document.createTextNode(' '));
                emailParagraph.appendChild(userEmailSpan);
                profileDetails.appendChild(emailParagraph);

                const bioParagraph = document.createElement('p');
                const bioStrong = document.createElement('strong');
                bioStrong.textContent = 'Bio:';
                const userBioSpan = document.createElement('span');
                userBioSpan.id = 'user-bio';
                userBioSpan.textContent = user?.bio || 'No Bio';
                bioParagraph.appendChild(bioStrong);
                bioParagraph.appendChild(document.createTextNode(' '));
                bioParagraph.appendChild(userBioSpan);
                profileDetails.appendChild(bioParagraph);

                const likedPostsLink = document.createElement('a');
                likedPostsLink.href = '/liked-posts';
                const likedPostsButton = document.createElement('button');
                likedPostsButton.classList.add('navbar__button');
                likedPostsButton.textContent = 'view posts i\'ve liked';
                likedPostsLink.appendChild(likedPostsButton);
                profileDetails.appendChild(likedPostsLink);

                profileInfo.appendChild(profileDetails);
                profileCard.appendChild(profileInfo);
                mainContent.appendChild(profileCard);

                const postsHeading = document.createElement('h3');
                postsHeading.classList.add('main-content__heading');
                postsHeading.textContent = 'Posts';
                mainContent.appendChild(postsHeading);

                const postsContainer = document.createElement('div');
                postsContainer.id = 'posts-container';
                mainContent.appendChild(postsContainer);

                // Update posts
                if (posts?.length > 0) {
                    posts.forEach(post => {
                        const postElement = document.createElement('div');
                        postElement.classList.add('card');
                        postElement.innerHTML = `
                            <p class="card__title"><a href="#" onclick="fetchPostDisplay('${post.uuid}')" style="color: #0172eb;">${post.title}</a></p>
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

        document.title = `${data.data.User.first_name}'s profile`
}
