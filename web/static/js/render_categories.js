
window.renderCategories = function(data) {
    const sidebar = document.querySelector('.sidebar');
    const categoriesList = document.getElementById('categories-list');
    categoriesList.innerHTML = ''; // Clear existing content

    data.Categories.forEach(category => {
        const categoryButton = document.createElement('button');
        categoryButton.classList.add('sidebar__list-item');
        categoryButton.textContent = category.name;
        categoryButton.onclick = () => fetchCategoryPosts(category.id);
        categoriesList.appendChild(categoryButton);
    });
    sidebar.append(categoriesList);
}
