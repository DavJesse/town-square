window.renderPosts = function(data) {
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