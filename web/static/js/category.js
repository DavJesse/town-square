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
                    <p class="card__title"><a href="/posts/display?pid=${post.uuid}" style="color: #0172eb;">${post.title}</a></p>
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
