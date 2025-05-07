import { populateComments } from '/static/js/populate_comments.js';
import { navigateTo } from '/static/js/routes.js';
import { renderErrorPage } from '/static/js/error.js';

export function fetchComments(postCard, postID) {
    fetch('/comment', {
        method: 'GET',
        credentials: 'include',
        headers: {
            'Accept': 'application/json'
        }
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
        postCard.innerHTML = '';
        populateComments(postCard, postID, data);
    })
    
    .catch(error => {
        console.error('Error fetching post data:', error);
    });
}