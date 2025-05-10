import { handleLikePost } from '/static/js/like_post.js';
import { handleDislikePost } from '/static/js/dislike_post.js';
import { populateComments } from '/static/js/populate_comments.js';
import { formatTextWithLineBreaks } from '/static/js/posts_utils.js';

export function populatePosts(posts) {
    let postsContainer = document.getElementById('posts_container');
    postsContainer.innerHTML = ""; // clear content for fresh population

    if (posts && posts.length > 0) {
        posts.forEach(post => {
            let postElement = document.createElement('div');
            postElement.id = 'post_card';

            let postHead = document.createElement('div');
            postHead.id ='post_head';
            postElement.appendChild(postHead);

            let postCreatorPic = document.createElement('img');
            postCreatorPic.id = 'post_creator_pic';
            postCreatorPic.src = post.creator_image? `/static/images/${post.creator_image}` : '/static/user-circle-svgrepo-com.svg';
            postHead.appendChild(postCreatorPic);

            let postCreatorInfoContainer = document.createElement('div');
            postCreatorInfoContainer.id = 'post_creator_info_container';
            postHead.appendChild(postCreatorInfoContainer);

            let postCreatorName = document.createElement('h3');
            postCreatorName.id = 'post_creator_name';
            postCreatorName.textContent = post.creator_first_name? `${post.creator_first_name.charAt(0).toUpperCase()}${post.creator_first_name.slice(1)} ${post.creator_last_name.charAt(0).toUpperCase()}${post.creator_last_name.slice(1)}` : post.creator_username;
            postCreatorInfoContainer.appendChild(postCreatorName);

            let postCreatorUsername = document.createElement('h4');
            postCreatorUsername.id = 'post_creator_username';
            postCreatorUsername.textContent = `@${post.creator_username}`;
            postCreatorInfoContainer.appendChild(postCreatorUsername);

            let postContentContainer = document.createElement('div');
            postContentContainer.id = 'post_content_container';
            postElement.appendChild(postContentContainer);

            let postTitle = document.createElement('h4');
            postTitle.id = 'post_title';
            postTitle.textContent = post.title;
            postContentContainer.appendChild(postTitle);

            let postContent = document.createElement('p');
            postContent.id = 'post_content';
            formatTextWithLineBreaks(post.content, postContent);
            postContentContainer.appendChild(postContent);

            if (post.media) {
                let postMedia = document.createElement('img');
                postMedia.id = 'post_media';
                postMedia.src = `/static/media/${post.media}`;
                postElement.appendChild(postMedia);
            }

            let postCreationDate = document.createElement('p');
            let options = { year: 'numeric', month: 'long', day: 'numeric', hour: '2-digit', minute: '2-digit' };
            postCreationDate.id = 'post_creation_date';
            postCreationDate.textContent = `Posted on ${new Date(post.created_at).toLocaleString('en-GB', options)}`;
            postElement.appendChild(postCreationDate);

            let postEngagement = document.createElement('div');
            postEngagement.id = 'post_engagement';
            postElement.appendChild(postEngagement);

            // Create engagement section
            let likeContainer = document.createElement('div');
            let dislikeContainer = document.createElement('div');
            let commentContainer = document.createElement('div');
            let likeLink = document.createElement('button');
            let dislikeLink = document.createElement('button');
            let commentLink = document.createElement('button');
            let likeCount = document.createElement('p');
            let dislikeCount = document.createElement('p');
            let commentCount = document.createElement('p');
            let likeIcon = document.createElement('span');
            let dislikeIcon = document.createElement('span');
            let commentIcon = document.createElement('span');

            // Set IDs and classes
            likeContainer.id = 'like_container';
            dislikeContainer.id = 'dislike_container';
            commentContainer.id = 'comment_container';
            likeLink.id = 'like_link';
            likeLink.dataset.postId = post.uuid;  // Add post ID as data attribute
            dislikeLink.id = 'dislike_link';
            dislikeLink.dataset.postId = post.uuid;  // Change this to prevent page refresh
            likeCount.id = 'post_like_count';
            dislikeCount.id = 'post_dislike_count';
            commentLink.id = 'comment_link';
            commentCount.id = 'engagement_count';
            likeIcon.classList.add('material-symbols-outlined');
            dislikeIcon.classList.add('material-symbols-outlined');
            commentIcon.classList.add('material-symbols-outlined');
            likeIcon.id = 'engagement_icon';
            dislikeIcon.id = 'engagement_icon';
            commentIcon.id = 'engagement_icon';

            // Set content
            likeCount.textContent = post.likes_count;
            dislikeCount.textContent = post.dislikes_count;
            commentCount.textContent = post.comments.length;
            likeIcon.textContent = 'thumb_up';
            dislikeIcon.textContent = 'thumb_down';
            commentIcon.textContent = 'comment';

            // Add event listeners
            likeLink.addEventListener('click', (e) => {
                e.preventDefault();
                handleLikePost(post.uuid, likeContainer, dislikeContainer);
            });

            dislikeLink.addEventListener('click', (e) => {
                e.preventDefault();
                handleDislikePost(post.uuid, likeContainer, dislikeContainer);
            });

            // Add event listener to fetch comments only if comments section is not rendered
            commentLink.addEventListener('click', (e) => {
                e.preventDefault();
                let commentsSection = postElement.querySelector('#comments_section');
                if (commentsSection) {
                    let commentInput = commentsSection.querySelector('#comment_input_field');
                    if (!commentInput) {
                        populateComments(postElement, post.uuid, post.comments);                    
                    }
                } else {
                    populateComments(postElement, post.uuid, post.comments);
                }
            });

            // Assemble the structure
            postEngagement.appendChild(likeContainer);
            postEngagement.appendChild(dislikeContainer);
            postEngagement.appendChild(commentContainer);
            likeLink.appendChild(likeCount);
            likeContainer.appendChild(likeLink);
            likeLink.appendChild(likeIcon);
            dislikeLink.appendChild(dislikeCount);
            dislikeContainer.appendChild(dislikeLink);
            dislikeLink.appendChild(dislikeIcon);
            commentLink.appendChild(commentCount);
            commentLink.appendChild(commentIcon);
            commentContainer.appendChild(commentLink);

            postsContainer.appendChild(postElement);
        });
    } else {
        let noPostsContainer = document.createElement('div');
        noPostsContainer.id = 'no_posts_container';

        let noPostsAlerts = document.createElement('p');
        noPostsAlerts.id = 'no_posts';
        noPostsAlerts.textContent = 'No posts available';
        noPostsContainer.appendChild(noPostsAlerts);

        postsContainer.appendChild(noPostsContainer);
    }
}

export function setToggleEventListeners(allPosts, likedPosts, userPosts) {
    // Extract page elements
    let allPostsButton = document.getElementById('all_posts_button');
    let likedPostsButton = document.getElementById('liked_posts_button');
    let myPostsButton = document.getElementById('my_posts_button');

    // Set null values to empty array
    allPosts = allPosts ?? [];
    likedPosts = likedPosts ?? [];
    userPosts = userPosts ?? [];

    allPostsButton.addEventListener('click', () => {
        allPostsButton.classList.add('active');
        likedPostsButton.classList.remove('active');
        myPostsButton.classList.remove('active');
        populatePosts(allPosts);
    });

    likedPostsButton.addEventListener('click', () => {
        likedPostsButton.classList.add('active');
        allPostsButton.classList.remove('active');
        myPostsButton.classList.remove('active');
        populatePosts(likedPosts);
    });

    myPostsButton.addEventListener('click', () => {
        myPostsButton.classList.add('active');
        allPostsButton.classList.remove('active');
        likedPostsButton.classList.remove('active');
        populatePosts(userPosts);
    });
}
