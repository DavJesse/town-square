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