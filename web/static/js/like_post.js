export function handleLikePost(postId, likeButton) {
    fetch('/posts/like', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
            'Accept': 'application/json',
        },
        body: `post-id=${postId}`,
        credentials: 'include'
    })

    .then(response => response.json())

    .then(data => {
        if (data.success) {
            // Update the likes count in the button
            const countSpan = likeButton.querySelector('.likes-count');
            countSpan.textContent = data.likesCount;
            
            // Optional: Add visual feedback
            likeButton.classList.add('liked');
            setTimeout(() => likeButton.classList.remove('liked'), 200);
        } else if (data.message === "Please login first") {
            window.location.href = '/login';
        } else {
            console.error('Like failed:', data.message);
        }
    })

    .catch(error => {
        console.error('Error:', error);
    });
}