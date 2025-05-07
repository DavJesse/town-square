import { fetchPostsPerCategory } from '/static/js/fetch_category_posts.js';

export function populateCategories(categories) {
    let categoriesContentContainer = document.getElementById('category_content_container');

    categories = categories ?? [];

    // Create buttons for each category
    categories.forEach(category => {
        let categoryButton = document.createElement('button');
        categoryButton.id = `category_button_${category.id}`;
        categoryButton.classList.add('inactive');
        categoryButton.href = `/categories/${category.id}`;
        categoryButton.textContent = category.name;

        // Update posts when category is clicked
        categoryButton.addEventListener('click', (e) => {
            e.preventDefault;
            fetchPostsPerCategory(category.id);
        })
        categoriesContentContainer.appendChild(categoryButton);
    });
}
