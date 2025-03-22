window.searchPosts = function() {
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