function toggleSidebar() {
    document.getElementById('sidebar').classList.toggle('active');
}

function renderPage(data) {
    renderNavbar(data);
    renderCategories(data);
    renderPosts(data);
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
