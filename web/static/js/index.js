fetchPosts();

function fetchPosts() {
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
