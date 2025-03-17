        // When the "Create Post" button is clicked, render the form
        document.getElementById('create-post-button').addEventListener('click', function(event) {
            event.preventDefault(); // Prevent default action (navigation)

            // Clear the main content area
            const mainContent = document.getElementById('main-content');
            mainContent.innerHTML = '';

            // Fetch categories and render the form
            fetch('/posts/create', {
                method: 'GET',
                credentials: 'include'
            })
            .then(response => response.json())
            .then(data => {
                // Create the form dynamically using the fetched data
                const formHTML = `
                    <div class="create-post-form">
                        <div class="form-container">
                            <h1>Create a New Post</h1>
                            <form id="create-post-form" enctype="multipart/form-data">
                                <label for="title">Title</label>
                                <input type="text" id="title" name="title" placeholder="Subject of your post" required>

                                <label for="content">Content</label>
                                <textarea id="content" name="content" placeholder="Share your thoughts" required></textarea>

                                <div class="custom-file-upload">
                                    <label for="image" class="upload-btn">Upload Image (PNG, GIF or JPG, max 20MB)</label>
                                    <input type="file" id="image" name="image" accept="image/gif, image/jpeg, image/png">
                                    <p id="file-name">No file chosen</p>
                                </div>

                                <div class="categories-container">
                                    <label for="categories" class="select-categories">Select Categories</label>
                                    ${data.categories.map(category => `
                                        <div class="category-item">
                                            <input type="checkbox" id="category-${category.id}" name="categories" value="${category.id}">
                                            <label for="category-${category.id}">${category.name}</label>
                                        </div>
                                    `).join('')}
                                </div>

                                <button type="submit">Create Post</button>
                            </form>
                        </div>
                    </div>
                `;

                mainContent.innerHTML = formHTML;

                // Handle form submission
                const form = document.getElementById('create-post-form');
                form.addEventListener('submit', function(e) {
                    e.preventDefault();

                    const formData = new FormData(form);

                    fetch('/posts/create', {
                        method: 'POST',
                        body: formData,
                        credentials: 'include'
                    })
                    .then(response => response.json())
                    .then(data => {
                        if (data.message) {
                            alert(data.message);  // Show the message (success or error)

                            /*---------------------------------------------
                            * TODO
                            ---------------------------------------------*/
                            // fetchAllPosts();
                            
                        }
                    })
                    .catch(error => console.error('Error creating post:', error));
                });
            })
            .catch(error => {
                console.error('Error fetching categories:', error);
                alert('Failed to load categories. Please try again.');
            });
        });