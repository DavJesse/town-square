
window.fetchPostDisplay = function(postID) {
    fetch(`/posts/display?pid=${postID}`, {
        method: 'GET',
        credentials: 'include', // Ensure cookies are sent with the request
        headers: {
            'Accept': 'application/json'
        }
    })
    .then(response => response.json())
    .then(data => {
        console.log("POST: ", data)
        const postContentContainer = document.getElementById('main-content');
        postContentContainer.innerHTML = '';

        // Create post header dynamically
        const postHeader = document.createElement('header');
        postHeader.classList.add('post-head');
        postHeader.innerHTML = `
            <h1>${data.PostData.title}</h1>
            <p class="author">Posted by:<strong> ${data.PostData.username}</strong></p>
            <p class="date">Posted on:<strong> ${data.PostData.created_at}</strong></p>
            <div class="metrics">
                <form action="/posts/like" method="post">
                    <input type="hidden" name="post-id" value="${data.PostData.uuid}">
                    <button type="submit" class="like-btn">👍 ${data.PostData.likes_count}</button>
                </form>
                <form action="/posts/dislike" method="post">
                    <input type="hidden" name="post-id" value="${data.PostData.uuid}">
                    <button type="submit" class="dislike-btn">👎 ${data.PostData.dislikes_count}</button>
                </form>
            </div>
        `;
        postContentContainer.appendChild(postHeader);

        // Create article section dynamically
        const articleSection = document.createElement('article');
        articleSection.innerHTML = `
            <div>
                ${data.PostData.media ? `<img src="/static/media/${data.PostData.media}" alt="Post Media" />` : '<p>No media available for this post.</p>'}
            </div>
            <div class="post-content">
                <p>${data.PostData.content.replace(/\n/g, '<br>')}</p>
            </div>
        `;
        postContentContainer.appendChild(articleSection);

        // Create category section dynamically
        const categorySection = document.createElement('section');
        categorySection.classList.add('category-section');
        const categoryList = document.createElement('ul');
        data.PostData.categories.forEach(category => {
            const categoryPill = document.createElement('li');
            categoryPill.classList.add('category-pill');
            categoryPill.textContent = category;
            categoryList.appendChild(categoryPill);
        });
        categorySection.appendChild(categoryList);
        postContentContainer.appendChild(categorySection);

        // Create comment section dynamically
        const commentSection = document.createElement('section');
        if (data.PostData.comments && data.PostData.comments.length > 0) {
            const commentsContainer = document.createElement('div');
            commentsContainer.classList.add('comment-section');
            data.PostData.comments.forEach(comment => {
                const commentElement = document.createElement('div');
                commentElement.classList.add('comment');
                commentElement.innerHTML = `
                    <p class="comment-author"><strong>${comment.creator}</strong> commented:</p>
                    <p>${comment.content}</p>
                    <p class="comment-date"><em>${comment.created_at}</em></p>
                    <div class="actions">
                        <form action="/comments/like" method="post">
                            <input type="hidden" name="comment-id" value="${comment.uuid}">
                            <input type="hidden" name="post-id" value="${data.PostData.uuid}">
                            <button type="submit" class="like-btn">👍 ${comment.likes_count}</button>
                        </form>
                        <form action="/comments/dislike" method="post">
                            <input type="hidden" name="comment-id" value="${comment.uuid}">
                            <input type="hidden" name="post-id" value="${data.PostData.uuid}">
                            <button type="submit" class="dislike-btn">👎 ${comment.dislikes_count}</button>
                        </form>
                    </div>
                `;
                commentsContainer.appendChild(commentElement);
            });
            commentSection.appendChild(commentsContainer);
        } else {
            commentSection.innerHTML = '<div class="comment-section">No comments yet. Be the first to comment!</div>';
        }

        // Add comment form if logged in
        if (data.IsLogged) {
            const commentForm = document.createElement('div');
            commentForm.classList.add('form-container');
            commentForm.innerHTML = `
                <form action="/comment" method="POST">
                    <input type="hidden" name="postUUID" value="${data.PostData.uuid}">
                    <input type="text" name="comment" placeholder="Write a comment..." required>
                    <button type="submit">Post</button>
                </form>
            `;
            commentSection.appendChild(commentForm);
        }

        postContentContainer.appendChild(commentSection);
    })
    .catch(error => {
        console.error('Error fetching post data:', error);
    });
}

function renderPage(data) {
    renderNavbar(data);
    renderCategories(data);
    renderPosts(data);
}

function renderNavbar(data) {
    const navbarHtml = `
        <div class="navbar">
            <a href="/">
                <h1 class="navbar__title">forum</h1>
            </a>
            
            <div class="navbar__actions">
                <form id="search-form" class="navbar__search-form">
                    <input type="text" id="search-query" name="q" class="navbar__search" placeholder="Search...">
                    <button type="button" class="navbar__button" onclick="searchPosts()">Search</button>
                </form>
                
                <button class="navbar__button" id="create-post-button">Create Post</button>

                    <form action="/logout" method="POST" style="display: inline;">
                        <button class="navbar__button" id="login-button">Log Out</button>
                    </form>
                    <div class="navbar__Profile__Pic">
                        ${data.ProfPic ? `
                            <button onclick="profile">
                                <img src="/static/images/${data.ProfPic}" alt="user">
                            </button>
                        ` : `
                            <a href="#" onclick="profile()">
                                <img src="/static/user-circle-svgrepo-com.svg" alt="user">
                            </a>
                        `}
                    </div>
            </div>
        </div>
    `;
    document.querySelector('.navbar').outerHTML = navbarHtml;
    document.getElementById('create-post-button').addEventListener('click', renderCreatePostForm);
}

function renderCategories(data) {
    const sidebar = document.querySelector('.sidebar');
    const categoriesList = document.getElementById('categories-list');
    categoriesList.innerHTML = ''; // Clear existing content

    data.Categories.forEach(category => {
        const categoryButton = document.createElement('button');
        categoryButton.classList.add('sidebar__list-item');
        categoryButton.textContent = category.name;
        categoryButton.onclick = () => fetchCategoryPosts(category.id);
        categoriesList.appendChild(categoryButton);
    });
    sidebar.append(categoriesList);
}


function renderPosts(data) {
    let postsHtml = '';
    data.Posts.forEach(post => {
        postsHtml += `
            <div class="card">
                <p class="card__title">
                    <a href="#" onclick="fetchPostDisplay('${post.uuid}')" style="color: #0172eb;">${post.title}</a>
                </p>
                <p class="card__subject">${post.creator}</p>
                <p class="card__description">${post.content}</p>
                ${post.media ? `<img src="/static/media/${post.media}" alt="Post Image" style="width: 70px; aspect-ratio: 1/1;">` : ''}
                
                <div class="actions">
                    <form action="/posts/like" method="post">
                        <input type="hidden" name="post-id" value="${post.uuid}">
                        <button type="submit">${post.likes_count} Like</button>
                    </form>
                    <form action="/posts/dislike" method="post">
                        <input type="hidden" name="post-id" value="${post.uuid}">
                        <button type="submit">${post.dislikes_count} Dislike</button>
                    </form>
                </div>
            </div>
        `;
    });
    document.getElementById('main-content').innerHTML = postsHtml;
}

function toggleSidebar() {
    document.getElementById('sidebar').classList.toggle('active');
}


function renderCreatePostForm() {
    // Clear the main content area
    const mainContent = document.getElementById('main-content');
    mainContent.innerHTML = '';

    // Fetch categories and render the form
    fetch('/posts/create', {
        method: 'GET',
        credentials: 'include'
    })
    .then(response => response.json())
    .then(data => {
        // Create the form dynamically using the fetched data
        const formHTML = `
            <div class="create-post-form">
                <div class="form-container">
                    <h1>Create a New Post</h1>
                    <form id="create-post-form" enctype="multipart/form-data">
                        <label for="title">Title</label>
                        <input type="text" id="title" name="title" placeholder="Subject of your post" required>

                        <label for="content">Content</label>
                        <textarea id="content" name="content" placeholder="Share your thoughts" required></textarea>

                        <div class="custom-file-upload">
                            <label for="image" class="upload-btn">Upload Image (PNG, GIF or JPG, max 20MB)</label>
                            <input type="file" id="image" name="image" accept="image/gif, image/jpeg, image/png">
                            <p id="file-name">No file chosen</p>
                        </div>

                        <div class="categories-container">
                            <label for="categories" class="select-categories">Select Categories</label>
                            ${data.categories.map(category => `
                                <div class="category-item">
                                    <input type="checkbox" id="category-${category.id}" name="categories" value="${category.id}">
                                    <label for="category-${category.id}">${category.name}</label>
                                </div>
                            `).join('')}
                        </div>

                        <button type="submit">Create Post</button>
                    </form>
                </div>
            </div>
        `;

        mainContent.innerHTML = formHTML;

        // Handle form submission
        const form = document.getElementById('create-post-form');
        form.addEventListener('submit', function(e) {
            e.preventDefault();

            const formData = new FormData(form);

            fetch('/posts/create', {
                method: 'POST',
                body: formData,
                credentials: 'include'
            })
            .then(response => response.json())
            .then(data => {
                if (data.message) {
                    alert(data.message);  // Show the message (success or error)

                    /*---------------------------------------------
                    * TODO
                    ---------------------------------------------*/
                    // fetchAllPosts();
                    
                }
            })
            .catch(error => console.error('Error creating post:', error));
        });
    })
    .catch(error => {
        console.error('Error fetching categories:', error);
        alert('Failed to load categories. Please try again.');
    });
}

// Function to fetch posts by category
function fetchCategoryPosts(categoryID) {
    fetch(`/categories/${categoryID}`, {
        method: 'GET',
        credentials: 'include' // Ensure cookies are sent with the request
    })
    .then(response => response.json())
    .then(data => {
        const mainContent = document.getElementById('main-content');
        mainContent.innerHTML = '';

        const categoryHeading = document.createElement('h3');
        categoryHeading.classList.add('main-content__heading');
        categoryHeading.textContent = `Posts in Category ${data.CategoryName}`;
        mainContent.appendChild(categoryHeading);

        const postsContainer = document.createElement('div');
        postsContainer.id = 'posts-container';
        mainContent.appendChild(postsContainer);

        // Update posts
        if (data.Posts?.length > 0) {
            data.Posts.forEach(post => {
                const postElement = document.createElement('div');
                postElement.classList.add('card');
                postElement.innerHTML = `
                    <p class="card__title">
                        <a href="#" onclick="fetchPostDisplay('${post.uuid}')" style="color: #0172eb;">${post.title}</a>
                    </p>
                    <p class="card__date">${post.created_at}</p>
                    <p class="card__description">${post.Content}</p>
                    ${post.media ? `<img src="/static/media/${post.media}" alt="Post Image" style="width: 70px; aspect-ratio: 1/1;">` : ''}
                    <div class="actions">
                        <form action="/posts/like" method="post">
                            <input type="hidden" name="post-id" value="${post.uuid}">
                            <button type="submit">${post.likes_count} Like</button>
                        </form>
                        <form action="/posts/dislike" method="post">
                            <input type="hidden" name="post-id" value="${post.uuid}">
                            <button type="submit">${post.dislikes_count} Dislike</button>
                        </form>
                    </div>
                `;
                postsContainer.appendChild(postElement);
            });
        } else {
            postsContainer.innerHTML = '<p>No posts available in this category</p>';
        }
    })
    .catch(error => {
        console.error('Error fetching category posts:', error);
    });
}

function searchPosts() {
    const query = document.getElementById('search-query').value;

    if (!query) {
        alert('Please enter a search query');
        return;
    }

    // Send search query via GET request using fetch
    fetch(`/search?q=${encodeURIComponent(query)}`, {
        method: 'GET',
        credentials: 'include' // Ensure cookies are sent with the request
    })
    .then(response => response.json())
    .then(data => {
        // Assuming the response has 'posts', 'is_logged', and 'prof_pic' properties
        const postsContainer = document.getElementById('main-content');
        postsContainer.innerHTML = '';  // Clear previous results

        if (data.posts.length > 0) {
            data.posts.forEach(post => {
                const postElement = document.createElement('div');
                postElement.classList.add('card');
                postElement.innerHTML = `
                    <p class="card__title">
                        <a href="#" onclick="fetchPostDisplay('${post.uuid}')" style="color: #0172eb;">${post.title}</a>
                    </p>
                    <p class="card__date">${post.created_at}</p>
                    <p class="card__description">${post.content}</p>
                    ${post.media ? `<img src="/static/media/${post.media}" alt="Post Image" style="width: 70px; aspect-ratio: 1/1;">` : ''}
                `;
                postsContainer.appendChild(postElement);
            });
        } else {
            postsContainer.innerHTML = '<p>No posts found</p>';
        }
    })
    .catch(error => {
        console.error('Error fetching posts:', error);
    });
}

function profile() {
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
}

(() => {
    fetch('/posts')  // Assuming '/posts' returns a JSON array of posts
        .then(response => response.json())
        .then(data => {
            if (data.code === 200 && data.data) {
                renderPage(data.data);
            } else if (data.redirect) {
                window.location.href = data.redirect;
            } else {
                console.error('Error fetching posts:', data);
                document.getElementById('main-content').innerHTML = '<p>Error loading posts.</p>';
            }
        })
        .catch(error => {
            console.error('Error fetching posts:', error);
            document.getElementById('main-content').innerHTML = '<p>Error loading posts.</p>';
        });
})();
