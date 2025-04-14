import { navigateTo } from '/static/js/routes.js';
import { renderErrorPage } from '/static/js/error.js';

export function renderCreatePostForm() {
    document.title = 'Create Post';
    let app = document.getElementById('app');
    
    // Create a smoke screen over page comtent
    let smokeScreen = document.createElement('div');
    smokeScreen.id = 'smoke_screen';
    app.appendChild(smokeScreen);

    // Create form container
    let formContainer = document.createElement('div');
    formContainer.id = 'create_post_form_container';
    smokeScreen.appendChild(formContainer);

    // Add close button
    let closeButton = document.createElement('span');
    closeButton.textContent = '✖';
    closeButton.id = 'close_button';
    closeButton.onclick = function () {
    smokeScreen.remove();
    document.body.style.overflow = 'auto'; // Re-enable scroll
    };
    formContainer.appendChild(closeButton);


    // Create forum logo
    let logo = document.createElement('h1');
    logo.textContent = 'Create New Post';
    formContainer.appendChild(logo);

    // Create text to capture user-generated errors
    let errorMessage = document.createElement('p');
    errorMessage.id = 'error_text';
    errorMessage.style.color = 'red';
    formContainer.appendChild(errorMessage);

    // Create form to capture user input
    let form = document.createElement('form');
    form.method = 'POST';
    form.action = '/posts/create';
    form.enctype = 'multipart/form-data';
    formContainer.appendChild(form);

    // Create title field
    let postTitle = document.createElement('input');
    postTitle.type = 'text';
    postTitle.id = 'post_title';
    postTitle.name = 'title';
    postTitle.placeholder = 'subject of your post';
    postTitle.required = true;
    form.appendChild(postTitle);

    // Create content field
    let postContent = document.createElement('textarea');
    postContent.id = 'post_content';
    postContent.name = 'content';
    postContent.placeholder = 'share your thoughts';
    postContent.required = true;
    form.appendChild(postContent);

    // Create image upload container
    let imageUpload = document.createElement('div');
    imageUpload.classList.add('custom-file-upload');
    form.appendChild(imageUpload);

    // Add image upload label
    let uploadButton = document.createElement('label');
    uploadButton.htmlFor = 'image';
    uploadButton.textContent = 'Upload Image (PNG, GIF or JPG, max 20MB)';
    uploadButton.classList.add('upload-btn');
    imageUpload.appendChild(uploadButton);

    // Create image upload field
    let postImage = document.createElement('input');
    postImage.type = 'file';
    postImage.id = 'image';
    postImage.name = 'image';
    postImage.accept = 'image/gif, image/jpeg, image/png';
    imageUpload.appendChild(postImage);

    // Create upload status section
    let fileName = document.createElement('p');
    fileName.id = 'file_name';
    fileName.textContent = 'No file chosen';
    imageUpload.appendChild(fileName);

    // Create submit button
    let removeButton = document.createElement('button');
    removeButton.id = "remove-image";
    removeButton.classList.add("remove-btn");
    removeButton.style.display = "none";
    removeButton.textContent = "Remove Image";
    imageUpload.appendChild(removeButton);

    // Create categories container
    let categoriesContainer = document.createElement('div');
    categoriesContainer.id = 'categories_container';
    form.appendChild(categoriesContainer);

    // Create categories title
    let categoriesTitle = document.createElement('h2');
    categoriesTitle.textContent = 'Select Categories';
    categoriesTitle.id = 'categories_title';
    categoriesContainer.appendChild(categoriesTitle);

    fetchCategories(); // fetch and render categories

    form.addEventListener('submit', async (e) => {
        e.preventDefault(); // stop normal form submission
    
        const formData = new FormData(form);
    
        try {
            const res = await fetch('/posts/create', {
                method: 'POST',
                body: formData,
                credentials: 'include' // if using sessions/cookies
            });
    
            const data = await res.json();
    
            if (res.ok && data.message === "Post created successfully") {
                // ✅ Redirect to index using your client-side router
                navigateTo('/');
            } else {
                // Show error from server
                errorMessage.textContent = data.message || 'Something went wrong.';
            }
        } catch (error) {
            console.error("Post submission error:", error);
            errorMessage.textContent = "An error occurred while submitting your post.";
        }
    });
    
    // Create submit button
    let submitButton = document.createElement('button');
    submitButton.type = 'submit';
    submitButton.textContent = 'Create Post';
    form.appendChild(submitButton);

    // Link js
    let scriptTag = document.createElement('script');
    scriptTag.src = '/static/js/onboarding.js';
    scriptTag.defer = true;
    app.appendChild(scriptTag);
}

function fetchCategories() {
    fetch('/posts/create', {
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
            }
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        return response.json();
    })

    .then(data => {
        if (data.code === 200) {
            if (data.message === 'Post created successfully') {
                alert(data.message);
                navigateTo('/');
                return
            }

            const categories = data.categories;

            let categoriesContainer = document.getElementById('categories_container');
            let combinedCluster = document.createElement('div');
            combinedCluster.id = 'combined_cluster';
            categoriesContainer.appendChild(combinedCluster);

            let leftCategoriesCluster = document.createElement('div');
            leftCategoriesCluster.id = 'left_categories_cluster';
            combinedCluster.appendChild(leftCategoriesCluster);

            let rightCategoriesCluster = document.createElement('div');
            rightCategoriesCluster.id = 'right_categories_cluster';
            combinedCluster.appendChild(rightCategoriesCluster);

            categories.forEach((category, index) => {
                // Create category container
                let categoryPill = document.createElement('div');
                categoryPill.id = 'category_pill';

                // Create category checkbox
                let categoryCheckbox = document.createElement('input');
                categoryCheckbox.type = 'checkbox';
                categoryCheckbox.id = `category_checkbox`;
                categoryCheckbox.name = 'categories';
                categoryCheckbox.value = category.id;
                categoryPill.appendChild(categoryCheckbox);

                // Create category label
                let categoryLabel = document.createElement('label');
                categoryLabel.htmlFor = `category_checkbox`;
                categoryLabel.textContent = category.name;
                categoryPill.appendChild(categoryLabel);

                // Append container category to appropriate cluster
                if (index % 2 === 0) {
                    leftCategoriesCluster.appendChild(categoryPill);
                } else {
                    rightCategoriesCluster.appendChild(categoryPill);
                }
            })            
        } else  {
            renderErrorPage(data.message, data.code);
        }
    })

    .catch(error => {        
        console.error('Error fetching categories:', error);
        renderErrorPage("Internal Server Error", 500);
    });
}

function setPostSubmissionListener() {
    
}