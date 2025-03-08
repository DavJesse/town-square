function fetchCategoryPosts(categoryID) {
    fetch(`/categories/${categoryID}`, {
        method: 'GET',
        credentials: 'include' // Ensure cookies are sent with the request
    })
    .then(response => response.json())
    .then(data => {
        if (data.code === 200) {
            const posts = data.data.Posts;

            // Clear existing content
            const mainContent = document.getElementById('main-content');
            mainContent.innerHTML = '';

            // Create category heading dynamically
            const categoryHeading = document.createElement('h3');
            categoryHeading.classList.add('main-content__heading');
            categoryHeading.textContent = `Posts in Category ${data.data.CategoryName}`;
            mainContent.appendChild(categoryHeading);

            const postsContainer = document.createElement('div');
            postsContainer.id = 'posts-container';
            mainContent.appendChild(postsContainer);

            // Update posts
            if (posts?.length > 0) {
                posts.forEach(post => {
                    const postElement = document.createElement('div');
                    postElement.classList.add('card');
                    postElement.innerHTML = `
                        <p class="card__title"><a href="/posts/display?pid=${post.UUID}" style="color: #0172eb;">${post.Title}</a></p>
                        <p class="card__date">${post.CreatedAt}</p>
                        <p class="card__description">${post.Content}</p>
                        ${post.Media ? `<img src="/static/media/${post.Media}" alt="Post Image" style="width: 70px; aspect-ratio: 1/1;">` : ''}
                        <div class="actions">
                            <form action="/posts/like" method="post">
                                <input type="hidden" name="post-id" value="${post.UUID}">
                                <button type="submit">${post.LikesCount} Like</button>
                            </form>
                            <form action="/posts/dislike" method="post">
                                <input type="hidden" name="post-id" value="${post.UUID}">
                                <button type="submit">${post.DislikesCount} Dislike</button>
                            </form>
                        </div>
                    `;
                    postsContainer.appendChild(postElement);
                });
            } else {
                postsContainer.innerHTML = '<p>No posts available in this category</p>';
            }
        } else {
            console.error('Error fetching category posts:', data.message);
        }
    })
    .catch(error => {
        console.error('Error:', error);
    });
}


// Fetch categories
fetch('/categories', {
    method: 'GET',
    credentials: 'include' // Ensure cookies are sent with the request
})
.then(response => response.json())
.then(data => {
    const categoriesList = document.getElementById('categories-list');
    categoriesList.innerHTML = '';

    data.forEach(category => {
        const categoryButton = document.createElement('button');
        categoryButton.classList.add('sidebar__list-item');
        categoryButton.textContent = category.name;
        categoryButton.onclick = () => fetchCategoryPosts(category.id);
        categoriesList.appendChild(categoryButton);
    });
})
.catch(error => {
    console.error('Error fetching categories:', error);
});

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
                    <p class="card__title"><a href="/posts/display?pid=${post.UUID}" style="color: #0172eb;">${post.Title}</a></p>
                    <p class="card__date">${post.CreatedAt}</p>
                    <p class="card__description">${post.Content}</p>
                    ${post.Media ? `<img src="/static/media/${post.Media}" alt="Post Image" style="width: 70px; aspect-ratio: 1/1;">` : ''}
                    <div class="actions">
                        <form action="/posts/like" method="post">
                            <input type="hidden" name="post-id" value="${post.UUID}">
                            <button type="submit">${post.LikesCount} Like</button>
                        </form>
                        <form action="/posts/dislike" method="post">
                            <input type="hidden" name="post-id" value="${post.UUID}">
                            <button type="submit">${post.DislikesCount} Dislike</button>
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
