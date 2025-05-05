import { navigateTo } from "/static/js/routes.js";
import { renderErrorPage } from "/static/js/error.js";

export function handleLikeComment(commentID, postID, likeContainer) {
    let formData = new URLSearchParams();
    formData.append(`comment-id=${commentID}`);
    formData.append(`post-id=${postID}`);

    fetch('/comments/like', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
            'Accept': 'application/json',
        },
        body: formData,
        credentials: 'include'
    })
    .then(response => {
        if (!response.ok) {
            if (response.status === 401) {
                navigateTo('/login');
                return;
            } else {
                renderErrorPage(response.statusText, response.status);
            }
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        return response.json();
    })
    .then(data => {
        if (data.success) {
            // Update the likes count in the container
            const countSpan = likeContainer.querySelector('#comment_engagement_count');
            if (countSpan) {
                countSpan.textContent = data.likesCount;
            }
            
        } else {
            console.error('Like failed:', data.message);
        }
    })
    .catch(error => {
        console.error('Error:', error);
    });
}
    